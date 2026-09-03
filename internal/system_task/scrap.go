package systemtask

import (
	"context"
	"encoding/json"
	"izumi/constant"
	"izumi/db"
	"izumi/internal/handler"
	"izumi/internal/service"
	"izumi/model"
	"izumi/scraper"
	"izumi/scraper/dlsite"
	"izumi/scraper/twodfan"
	"izumi/scraper/vndb"
	"izumi/utils"
	"strings"
	"sync/atomic"
	"time"

	"github.com/abadojack/whatlanggo"
	"github.com/dokidokikoi/go-common/gopool"
	zaplog "github.com/dokidokikoi/go-common/log/zap"
	meta "github.com/dokidokikoi/go-common/meta/option"
	"github.com/dokidokikoi/go-common/notice"
	"github.com/dokidokikoi/go-common/tools"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// 任务整体进度按 10000 刻度计算，划分三个阶段：
//
//	抓取 (scrap) 5000 -> 合并 (merge) 1000 -> 存图/入库 (save) 4000
//
// 这些权重集中在此定义，避免 10000/5000/4000 等魔数散落在多处。
const (
	progressTotal int32 = 10000
	weightScrap   int32 = 5000
	weightMerge   int32 = 1000
	weightSave    int32 = 4000
	tickInterval        = 500 * time.Millisecond

	// defaultGameVersion 是未指定版本号时的占位值，
	// 避免 GameInstance 的 uk_game_version 唯一约束在多版本场景下因空串冲突。
	defaultGameVersion = "default"
)

// defaultVersion 在 v 为空串时返回占位版本号 defaultGameVersion。
func defaultVersion(v string) string {
	if v == "" {
		return defaultGameVersion
	}
	return v
}

// progress 负责采集任务的进度推送：按固定间隔广播当前进度，
// 任务结束时广播一次完成（100%）进度并停止 ticker。
//
// 进度通过「阶段」推进：每个阶段声明自己的权重与子项总数，
// 子项完成时调用 phase.done()，phase 内部按比例折算并 clamp，
// 保证单阶段累加不超过其权重、整体不超过 progressTotal。
type progress struct {
	taskID uint
	total  int32
	curr   int32 // 已推进的刻度，始终 <= total
	event  string
	stop   chan struct{}
}

func newProgress(taskID uint) *progress {
	return &progress{
		taskID: taskID,
		total:  progressTotal,
		event:  utils.ConcatEvent(constant.TOPIC_SCRAPER, constant.EVENT_SCRAPER_AUTOSCRAPING),
		stop:   make(chan struct{}),
	}
}

// broadcast 发送一次进度广播。
func (p *progress) broadcast() {
	notice.HubIns.SendBroadcast("", notice.NoticeResponse{
		Rid:   uuid.NewString(),
		Event: p.event,
		Data: map[string]any{
			"task_id":  p.taskID,
			"proccess": atomic.LoadInt32(&p.curr),
			"total":    p.total,
		},
	})
}

// start 启动定时进度推送 goroutine。
func (p *progress) start() {
	go func() {
		ticker := time.NewTicker(tickInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				p.broadcast()
			case <-p.stop:
				return
			}
		}
	}()
}

// add 原子增加刻度并 clamp 到 [0, total]，防止回调异常导致进度溢出。
func (p *progress) add(step int32) {
	for {
		old := atomic.LoadInt32(&p.curr)
		next := old + step
		if next > p.total {
			next = p.total
		}
		if next < 0 {
			next = 0
		}
		if atomic.CompareAndSwapInt32(&p.curr, old, next) {
			return
		}
	}
}

// set 原子设置当前刻度（仍 clamp 到 total）。
func (p *progress) set(v int32) {
	if v > p.total {
		v = p.total
	}
	if v < 0 {
		v = 0
	}
	atomic.StoreInt32(&p.curr, v)
}

// phase 返回阶段化进度的构造器。weight 是该阶段占用的整体刻度，
// total 是该阶段的子项总数（用于按比例推进）；阶段之间串行使用。
// base 锁定为构造时的当前刻度，该阶段只在 [base, base+weight] 区间内推进。
func (p *progress) phase(weight int32, total int) *phaseProgress {
	return &phaseProgress{
		p:      p,
		weight: weight,
		total:  max(int32(total), 1),
		base:   atomic.LoadInt32(&p.curr),
	}
}

// finish 停止 ticker 并广播一次完成进度（100%），应放在 defer 中调用。
func (p *progress) finish() {
	close(p.stop)
	p.set(p.total)
	p.broadcast()
}

// phaseProgress 描述一个阶段：把该阶段的权重 weight 按子项总数 total 均分，
// 子项完成时调用 done 推进进度。该阶段在 [base, base+weight] 区间内推进，
// 不触碰其它阶段，且线程安全（子项回调可能并发触发）。
type phaseProgress struct {
	p      *progress
	weight int32 // 该阶段占用的整体刻度
	total  int32 // 该阶段的子项总数
	base   int32 // 该阶段起点刻度（构造时锁定）
	doneN  int32 // 已完成子项数
}

// done 标记完成一个子项，按已完成比例把进度推进到 [base, base+weight] 内。
func (ph *phaseProgress) done() {
	n := atomic.AddInt32(&ph.doneN, 1)
	if n > ph.total {
		n = ph.total
	}
	// 该阶段应达到的刻度 = 起点 + 权重 * (已完成/总数)
	target := ph.base + ph.weight*n/ph.total
	// 只允许在该阶段区间内前进，避免回退或越界
	for {
		cur := atomic.LoadInt32(&ph.p.curr)
		if target <= cur || target > ph.base+ph.weight {
			return
		}
		if atomic.CompareAndSwapInt32(&ph.p.curr, cur, target) {
			return
		}
	}
}

// doneAll 直接把该阶段进度置满（用于阶段总量未知、走不到 total 时的兜底）。
func (ph *phaseProgress) doneAll() {
	for {
		cur := atomic.LoadInt32(&ph.p.curr)
		target := ph.base + ph.weight
		if target <= cur {
			return
		}
		if atomic.CompareAndSwapInt32(&ph.p.curr, cur, target) {
			return
		}
	}
}

func StartAutoScrap() {
	ts, err := db.GetStore().SystemTask().List(context.Background(), &model.SystemTask{
		Type:  model.SystemTaskTypeScrap,
		State: model.SystemTaskStateRunning,
	}, &meta.ListOption{Order: "id desc"})
	if err != nil {
		zaplog.L().Error("system download task error", zap.Error(err))
		return
	}
	if len(ts) == 0 {
		return
	}

	for _, t := range ts {
		AutoScrap(t)
	}
}

func AutoScrap(t *model.SystemTask) {
	gopool.CtxGo(context.Background(), func() {
		var result string
		var e error
		defer func() {
			state := model.SystemTaskStateDone
			if e != nil {
				zaplog.L().Error("system scrap task error", zap.Error(e))
				state = model.SystemTaskStateFailed
			}
			if err := db.GetStore().SystemTask().Update(context.Background(), &model.SystemTask{
				ID:     t.ID,
				State:  state,
				Result: result,
			}, nil); err != nil {
				zaplog.L().Error("system update error", zap.Error(err))
			}

			if err := notice.HubIns.SendBroadcast("", notice.NoticeResponse{
				Rid:     uuid.NewString(),
				Event:   utils.ConcatEvent(constant.TOPIC_SCRAPER, constant.EVENT_SCRAPER_AUTOSCRAP),
				Success: e == nil,
				Message: result,
			}); err != nil {
				zaplog.L().Error("send notify error", zap.Error(err))
			}
		}()

		// vndb 是数据合并的主源，必须存在
		if !containsScraper(t.Param.ScrapObjs, vndb.Name) {
			result = "auto scrap need vndb item"
			e = errors.New(result)
			return
		}

		prog := newProgress(t.ID)
		prog.start()
		defer prog.finish()

		// 阶段 1/3：抓取。每个 scraper 完成回调一次。
		scrapPhase := prog.phase(weightScrap, len(t.Param.ScrapObjs))
		requestID := uuid.NewString()
		wait, err := service.NewScrap(db.GetStore()).Scrap(context.Background(), requestID, t.Param.ScrapObjs, func(scraperName string, success bool) {
			scrapPhase.done()
		})
		if err != nil {
			e = err
			result = err.Error()
			return
		}
		wait.Wait()
		scrapPhase.doneAll() // 抓取阶段结束，补满该阶段权重

		// 阶段 2/3：合并 vndb 与其它来源
		list, err := db.GetStore().Task().List(context.Background(), &model.Task{
			RequestID: requestID,
			Status:    model.TaskStatusSucceed,
		}, nil)
		if err != nil {
			e = err
			result = err.Error()
			return
		}

		itemM, err := groupScrapItems(list)
		if err != nil {
			e = err
			result = err.Error()
			return
		}
		gVo, ok := mergeScrapResults(itemM)
		if !ok {
			result = "vndb item scrap failed"
			return
		}
		prog.add(weightMerge) // 合并阶段无子项回调，整体完成后一次性推进

		// 阶段 3/3：存图 + 入库。UpsertFull 内部按文件数回调 step 刻度（权重 weightSave）。
		srv := service.NewGame(db.GetStore())
		if err := srv.UpsertFull(context.Background(), &gVo, &model.GameInstance{
			Path:    t.Param.Path,
			Version: defaultVersion(t.Param.Version),
		}, func(step int) {
			prog.add(int32(step))
		}); err != nil {
			e = err
			result = err.Error()
			return
		}
		result = gVo.Name

		// 下载补充信息，失败仅记录，不影响任务结果
		go func() {
			if derr := srv.DownloadInfo(context.Background(), gVo.ID, time.Now()); derr != nil {
				zaplog.L().With(zap.Uint("game id", gVo.ID)).Error("downloadInfo", zap.Error(derr))
			}
		}()
	})
}

// containsScraper 判断抓取对象列表中是否包含指定 scraper。
func containsScraper(objs []model.ScrapObj, name string) bool {
	for _, obj := range objs {
		if obj.Name == name {
			return true
		}
	}
	return false
}

// groupScrapItems 把抓取任务的结果按 scraper 名分组返回。
func groupScrapItems(list []*model.Task) (map[string][]scraper.GameItem, error) {
	itemM := map[string][]scraper.GameItem{}
	for _, l := range list {
		var item scraper.GameItem
		if err := json.Unmarshal([]byte(l.Result), &item); err != nil {
			zaplog.L().Error("system scrap task error", zap.Error(err))
			continue
		}
		itemM[item.ScraperName] = append(itemM[item.ScraperName], item)
	}
	return itemM, nil
}

// mergeScrapResults 以 vndb 为主源，合并其它来源的别名、标签、图片、角色、制作人员等信息。
// 第二个返回值表示是否拿到有效的 vndb 数据。
func mergeScrapResults(itemM map[string][]scraper.GameItem) (handler.GameVo, bool) {
	items, ok := itemM[vndb.Name]
	if !ok || len(items) == 0 {
		return handler.GameVo{}, false
	}
	delete(itemM, vndb.Name)

	gVo := items[0].GameVo

	// 建立名称（小写去空格）到下标的索引，用于跨来源匹配角色 / 制作人员
	aliasM := map[string]struct{}{}
	for _, a := range gVo.Alias {
		aliasM[a] = struct{}{}
	}
	tagM := map[string]struct{}{}
	for _, t := range gVo.Tags {
		tagM[t.Name] = struct{}{}
	}
	nameIdx := buildMatchIndex(gVo.Characters, func(c handler.CharacterVo) []string { return append(c.Alias, c.Name) })
	staffIdx := buildMatchIndex(gVo.Staff, func(s handler.StaffVo) []string { return append(s.Alias, s.Name) })

	for _, item := range itemM {
		for _, i := range item {
			mergeItem(&gVo, i, aliasM, tagM, nameIdx, staffIdx)
		}
	}
	gVo.Alias = tools.Keys(aliasM)

	// story 来源优先级：twodfan > dlsite（这两个来源的简介更完整）
	pickStory := func(name string) (string, bool) {
		if vs, ok := itemM[name]; ok {
			for _, v := range vs {
				if v.Story != "" {
					return v.Story, true
				}
			}
		}
		return "", false
	}
	if s, ok := pickStory(dlsite.Name); ok {
		gVo.Story = s
	}
	if s, ok := pickStory(twodfan.Name); ok {
		gVo.Story = s
	}

	// 仅保留有明确职责关系的制作人员
	gVo.Staff = filterRelatedStaff(gVo.Staff)
	gVo.Images = tools.NewSet(gVo.Images...).Slice()
	return gVo, true
}

// mergeItem 把单个其它来源的数据合并进主 GameVo。
func mergeItem(
	gVo *handler.GameVo, i scraper.GameItem,
	aliasM map[string]struct{}, tagM map[string]struct{},
	characterM, staffM *matchIndex,
) {
	gVo.RelIDs = append(gVo.RelIDs, i.RelIDs...)
	for _, a := range i.Alias {
		aliasM[a] = struct{}{}
	}
	for _, t := range i.Tags {
		if _, ok := tagM[t.Name]; !ok {
			tagM[t.Name] = struct{}{}
			gVo.Tags = append(gVo.Tags, t)
		}
	}
	gVo.Images = append(gVo.Images, i.Images...)
	gVo.Links = append(gVo.Links, i.Links...)
	if len(gVo.Brands) == 0 {
		gVo.Brands = append(gVo.Brands, i.Brands...)
	}
	if (gVo.Category == nil || gVo.Category.Name == "") && i.Category != nil && i.Category.Name != "" {
		gVo.Category = &model.Category{Name: lettersOnly(strings.ToUpper(i.Category.Name))}
	}
	if gVo.IssueDate.IsZero() {
		gVo.IssueDate = i.IssueDate
	}
	if gVo.DlCode == "" {
		gVo.DlCode = i.DlCode
	}
	if gVo.JanCode == "" {
		gVo.JanCode = i.JanCode
	}
	if gVo.Price == "" {
		gVo.Price = i.Price
	}
	mergeStory(&gVo.Story, i.Story)
	for _, c := range i.Characters {
		mergeCharacter(gVo, c, characterM)
	}
	for _, s := range i.Staff {
		mergeStaff(gVo, s, staffM)
	}
}

// mergeCharacter 把单个角色并入主 GameVo：能唯一匹配则补全字段，否则作为新角色追加。
func mergeCharacter(gVo *handler.GameVo, c handler.CharacterVo, idx *matchIndex) {
	name := utils.ToLowerNoSpace(c.Name)
	if m, ok := idx.get(name); ok && len(m) == 1 {
		for k := range m {
			gVo.Characters[k].RelIDs = append(gVo.Characters[k].RelIDs, c.RelIDs...)
			gVo.Characters[k].Images = append(gVo.Characters[k].Images, c.Images...)
			if c.Cover != "" {
				gVo.Characters[k].Images = append(gVo.Characters[k].Images, c.Cover)
			}
			if gVo.Characters[k].Gender == model.UnKnown {
				gVo.Characters[k].Gender = c.Gender
			}
			if gVo.Characters[k].Rlation == "" {
				gVo.Characters[k].Rlation = c.Rlation
			}
			mergeStory(&gVo.Characters[k].Summary, c.Summary)
		}
		return
	}
	idx.addNew(name, len(gVo.Characters))
	gVo.Characters = append(gVo.Characters, c)
}

// mergeStaff 把单个制作人员并入主 GameVo：能唯一匹配则补全字段，否则作为新制作人员追加。
func mergeStaff(gVo *handler.GameVo, s handler.StaffVo, idx *matchIndex) {
	name := utils.ToLowerNoSpace(s.Name)
	if m, ok := idx.get(name); ok && len(m) == 1 {
		for k := range m {
			gVo.Staff[k].RelIDs = append(gVo.Staff[k].RelIDs, s.RelIDs...)
			if gVo.Staff[k].Cover == "" {
				gVo.Staff[k].Cover = s.Cover
			} else if s.Cover != "" {
				gVo.Staff[k].Images = append(gVo.Staff[k].Images, s.Cover)
			}
			gVo.Staff[k].Images = append(gVo.Staff[k].Images, s.Images...)
			gVo.Staff[k].Links = append(gVo.Staff[k].Links, s.Links...)
			if gVo.Staff[k].Gender == model.UnKnown {
				gVo.Staff[k].Gender = s.Gender
			}
			if len(gVo.Staff[k].Relation) == 0 {
				gVo.Staff[k].Relation = s.Relation
			}
			mergeStory(&gVo.Staff[k].Summary, s.Summary)
		}
		return
	}
	idx.addNew(name, len(gVo.Staff))
	gVo.Staff = append(gVo.Staff, s)
}

// mergeStory 用 dst 补全/替换 src：src 为空时直接采用 dst；两者都有时按 needReplace 决定。
func mergeStory(src *string, dst string) {
	if *src == "" {
		*src = dst
		return
	}
	if dst != "" && needReplace(*src, dst) {
		*src = dst
	}
}

// filterRelatedStaff 仅保留至少有一个明确职责关系（非 PRelationUnknown）的制作人员。
func filterRelatedStaff(staff []handler.StaffVo) []handler.StaffVo {
	staffs := make([]handler.StaffVo, 0, len(staff))
	for _, s := range staff {
		for _, r := range s.Relation {
			if r != model.PRelationUnknown {
				staffs = append(staffs, s)
				break
			}
		}
	}
	return staffs
}

// lettersOnly 仅保留字符串中的大写字母（输入应已 ToUpper）。
func lettersOnly(s string) string {
	b := strings.Builder{}
	for _, r := range s {
		if r >= 'A' && r <= 'Z' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// matchIndex 维护「名称（小写去空格）-> 该名称命中的下标集合」。
// 用于跨来源匹配角色 / 制作人员：一个名称可能对应多个下标（重名）。
type matchIndex struct {
	m map[string]map[int]struct{}
}

func newMatchIndex() *matchIndex {
	return &matchIndex{m: map[string]map[int]struct{}{}}
}

func (mi *matchIndex) add(name string, idx int) {
	if mi.m[name] == nil {
		mi.m[name] = map[int]struct{}{}
	}
	mi.m[name][idx] = struct{}{}
}

func (mi *matchIndex) get(name string) (map[int]struct{}, bool) {
	m, ok := mi.m[name]
	return m, ok
}

// addNew 为一个新追加的下标建立索引（add 的语义别名，提升调用处可读性）。
func (mi *matchIndex) addNew(name string, idx int) { mi.add(name, idx) }

// buildFrom 用一组元素及其别名/主名初始化索引。
func buildMatchIndex[T any](items []T, namesOf func(T) []string) *matchIndex {
	mi := newMatchIndex()
	for i, it := range items {
		for _, n := range namesOf(it) {
			mi.add(utils.ToLowerNoSpace(n), i)
		}
	}
	return mi
}

func needReplace(src, dst string) bool {
	srcLang := whatlanggo.DetectLang(src).Iso6391()
	dstlang := whatlanggo.DetectLang(dst).Iso6391()
	if srcLang == "zh" {
		if dstlang == "zh" {
			return len(src) < len(dst)
		}
		return false
	}
	if dstlang == "zh" {
		return true
	}
	if srcLang == "ja" {
		if dstlang == "ja" {
			return len(src) < len(dst)
		}
		return false
	}
	if dstlang == "ja" {
		return true
	}
	return false
}
