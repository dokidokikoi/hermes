package lua

import (
	"context"
	"izumi/scraper"
	"maps"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/yuin/gopher-lua"
	"go.uber.org/zap"
)

var (
	// ErrTimeout is returned when Lua execution times out
	ErrTimeout = errors.New("lua execution timeout")

	// DefaultExecutionTimeout is the default timeout for Lua script execution
	DefaultExecutionTimeout = 30 * time.Second
)

// LuaScraper implements scraper.IGameScraper interface
type LuaScraper struct {
	name     string
	script   string
	vmPool   *Pool
	headers  map[string]string
	proxy    string
	logger   *zap.Logger
	timeout  time.Duration

	mu sync.RWMutex
}

// NewLuaScraper creates a new Lua-based scraper
func NewLuaScraper(name, script string, vmPool *Pool, logger *zap.Logger, timeout time.Duration) *LuaScraper {
	return &LuaScraper{
		name:    name,
		script:  script,
		vmPool:  vmPool,
		headers: make(map[string]string),
		logger:  logger.Named(name),
		timeout: timeout,
	}
}

// GetItem retrieves game details from a URL
func (l *LuaScraper) GetItem(url string) (*scraper.GameItem, error) {
	ctx, cancel := context.WithTimeout(context.Background(), l.timeout)
	defer cancel()

	resultChan := make(chan *scraper.GameItem, 1)
	errChan := make(chan error, 1)

	go func() {
		result, err := l.executeGetItem(url)
		if err != nil {
			errChan <- err
			return
		}
		resultChan <- result
	}()

	select {
	case result := <-resultChan:
		return result, nil
	case err := <-errChan:
		return nil, err
	case <-ctx.Done():
		return nil, ErrTimeout
	}
}

func (l *LuaScraper) executeGetItem(url string) (*scraper.GameItem, error) {
	L := l.vmPool.Get()
	defer l.vmPool.Put(L)

	// Load and execute the script
	if err := L.DoString(l.script); err != nil {
		return nil, errors.Wrap(err, "failed to load script")
	}

	// Get the get_item function
	getItemFn := L.GetGlobal("get_item")
	if getItemFn.Type() != lua.LTFunction {
		return nil, errors.New("script must export a get_item function")
	}

	// Set script globals for headers and proxy
	l.setScriptGlobals(L)

	// Call get_item(url)
	L.Push(getItemFn)
	L.Push(lua.LString(url))

	if err := L.PCall(1, 1, nil); err != nil {
		return nil, errors.Wrap(err, "failed to call get_item")
	}

	// Get the return value
	ret := L.Get(-1)
	L.Pop(1)

	if ret.Type() == lua.LTNil {
		return nil, nil
	}

	tbl, ok := ret.(*lua.LTable)
	if !ok {
		return nil, errors.New("get_item must return a table")
	}

	luaItem := luaGameItemFromTable(L, tbl)
	return luaItem.ToGameItem(l.name), nil
}

// Search searches for games by keyword
func (l *LuaScraper) Search(keyword string, page int) ([]*scraper.SearchItem, error) {
	ctx, cancel := context.WithTimeout(context.Background(), l.timeout)
	defer cancel()

	resultChan := make(chan []*scraper.SearchItem, 1)
	errChan := make(chan error, 1)

	go func() {
		result, err := l.executeSearch(keyword, page)
		if err != nil {
			errChan <- err
			return
		}
		resultChan <- result
	}()

	select {
	case result := <-resultChan:
		return result, nil
	case err := <-errChan:
		return nil, err
	case <-ctx.Done():
		return nil, ErrTimeout
	}
}

func (l *LuaScraper) executeSearch(keyword string, page int) ([]*scraper.SearchItem, error) {
	L := l.vmPool.Get()
	defer l.vmPool.Put(L)

	// Load and execute the script
	if err := L.DoString(l.script); err != nil {
		return nil, errors.Wrap(err, "failed to load script")
	}

	// Get the search function
	searchFn := L.GetGlobal("search")
	if searchFn.Type() != lua.LTFunction {
		return nil, errors.New("script must export a search function")
	}

	// Set script globals for headers and proxy
	l.setScriptGlobals(L)

	// Call search(keyword, page)
	L.Push(searchFn)
	L.Push(lua.LString(keyword))
	L.Push(lua.LNumber(page))

	if err := L.PCall(2, 1, nil); err != nil {
		return nil, errors.Wrap(err, "failed to call search")
	}

	// Get the return value
	ret := L.Get(-1)
	L.Pop(1)

	if ret.Type() == lua.LTNil {
		return []*scraper.SearchItem{}, nil
	}

	tbl, ok := ret.(*lua.LTable)
	if !ok {
		return nil, errors.New("search must return a table")
	}

	var items []*scraper.SearchItem
	tbl.ForEach(func(_, value lua.LValue) {
		if itemTbl, ok := value.(*lua.LTable); ok {
			luaItem := luaSearchItemFromTable(L, itemTbl)
			items = append(items, luaItem.ToSearchItem(l.name))
		}
	})

	return items, nil
}

// GetName returns the scraper name
func (l *LuaScraper) GetName() string {
	return l.name
}

// SetHeader sets HTTP headers for requests
func (l *LuaScraper) SetHeader(header map[string]string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	maps.Copy(l.headers, header)
}

// SetProxy sets the proxy for requests
func (l *LuaScraper) SetProxy(proxy string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.proxy = proxy
}

// setScriptGlobals sets headers and proxy as global variables in the Lua script
func (l *LuaScraper) setScriptGlobals(L *lua.LState) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	// Set headers table
	headersTbl := L.NewTable()
	for k, v := range l.headers {
		L.SetField(headersTbl, k, lua.LString(v))
	}
	L.SetGlobal("scraper_headers", headersTbl)

	// Set proxy
	L.SetGlobal("scraper_proxy", lua.LString(l.proxy))
}

// GetInfo retrieves plugin info from the script
func (l *LuaScraper) GetInfo() (name string, version string, err error) {
	L := l.vmPool.Get()
	defer l.vmPool.Put(L)

	// Load and execute the script
	if err := L.DoString(l.script); err != nil {
		return "", "", errors.Wrap(err, "failed to load script")
	}

	// Get the get_info function
	getInfoFn := L.GetGlobal("get_info")
	if getInfoFn.Type() != lua.LTFunction {
		return l.name, "", nil
	}

	// Call get_info()
	L.Push(getInfoFn)
	if err := L.PCall(0, 1, nil); err != nil {
		return l.name, "", errors.Wrap(err, "failed to call get_info")
	}

	// Get the return value
	ret := L.Get(-1)
	L.Pop(1)

	tbl, ok := ret.(*lua.LTable)
	if !ok {
		return l.name, "", nil
	}

	name = lua.LVAsString(tbl.RawGetString("name"))
	version = lua.LVAsString(tbl.RawGetString("version"))

	if name == "" {
		name = l.name
	}

	return name, version, nil
}

// UpdateScript updates the Lua script (for hot reload)
func (l *LuaScraper) UpdateScript(script string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Verify the script is valid by loading it
	L := l.vmPool.Get()
	defer l.vmPool.Put(L)

	if err := L.DoString(script); err != nil {
		return errors.Wrap(err, "script validation failed")
	}

	// Verify get_info returns valid info
	getInfoFn := L.GetGlobal("get_info")
	if getInfoFn.Type() == lua.LTFunction {
		L.Push(getInfoFn)
		if err := L.PCall(0, 1, nil); err != nil {
			return errors.Wrap(err, "get_info validation failed")
		}
		L.Pop(1)
	}

	// Script is valid, update it
	l.script = script
	l.logger.Info("script updated successfully")
	return nil
}