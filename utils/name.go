package utils

import (
	"regexp"
	"strings"
)

// 清洗游戏名时用到的正则：
//   - 括号内容：[...] (...) （...）【...】
//   - 圈号序号：①②③...⑳
var (
	bracketRe = regexp.MustCompile(`[\[\(（【].*?[\]\)）】]`)
	circleRe  = regexp.MustCompile(`[①-⑳]`)
)

// versionNumRe 匹配版本号（v1.0 / Ver2 / ver.1.2 等），命名组 version 捕获纯数字部分。
// CleanGameName（删除）与 ExtractVersion（提取）共用同一正则，避免规则漂移。
var versionNumRe = regexp.MustCompile(`(?i)v(?:er)?\.?\s*(?P<version>\d+(\.\d+)*)`)

// ExtractVersion 从游戏名/文件名中提取版本号，返回去掉 v/ver 前缀的纯数字部分
// （如 "游戏 v1.2" -> "1.2"；"Title Ver3" -> "3"）。提取不到返回空串。
// 应在 CleanGameName 之前调用（清洗会删除版本号）。
func ExtractVersion(s string) string {
	m := versionNumRe.FindStringSubmatch(s)
	if len(m) < 2 {
		return ""
	}
	return m[1]
}

// fullwidthToHalfwidth 把全角 ASCII / 数字 / 标点转成半角，
// 便于与 vndb 等数据源匹配。不含全角假名/汉字。
func fullwidthToHalfwidth(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		switch {
		case r >= '！' && r <= '～': // U+FF01 ~ U+FF5E
			r = r - '！' + '!'
		case r == '　': // 全角空格
			r = ' '
		}
		b.WriteRune(r)
	}
	return b.String()
}

// CleanGameName 清洗从文件名/目录名提取出的游戏名，便于在数据源（如 vndb）搜索：
//   - 去除括号及其内容（资源站标记，如 [DL版]、(同人)、【汉化】）
//   - 去除版本号（v1.0、Ver2 等）
//   - 去除圈号序号（①②③）
//   - 全角字母/数字/标点转半角
//   - 折叠多余空白并 trim
func CleanGameName(s string) string {
	s = fullwidthToHalfwidth(s)
	s = bracketRe.ReplaceAllString(s, "")
	s = versionNumRe.ReplaceAllString(s, "")
	s = circleRe.ReplaceAllString(s, "")
	s = strings.TrimSpace(s)
	// 折叠连续空白为单个空格
	fields := strings.Fields(s)
	return strings.Join(fields, " ")
}
