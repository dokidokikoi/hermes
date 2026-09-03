// Package ehtag 拉取并解析 EhTagTranslation/Database 的标签数据库
// (db.text.json)，作为手动匹配 decided_tags 时的规范 tag 来源。
package ehtag

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const Name = "ehtag"

// DatabaseURL 指向 DatabaseReleases 仓库的 db.text.json 发布文件。
const DatabaseURL = "https://raw.githubusercontent.com/EhTagTranslation/DatabaseReleases/master/db.text.json"

const userAgent = "dokidokikoi/izumi (https://github.com/dokidokikoi/izumi)"

// Entry 是单个标签条目：NS + Key 唯一，Name 为中文译名。
type Entry struct {
	NS    string
	Key   string
	Name  string
	Intro string
}

// Database 为 namespace -> key -> Entry。
type Database map[string]map[string]Entry

type entryJSON struct {
	Name  string `json:"name"`
	Intro string `json:"intro"`
}

type databaseFile struct {
	Data []struct {
		Namespace string               `json:"namespace"`
		Data      map[string]entryJSON `json:"data"`
	} `json:"data"`
}

var client = &http.Client{Timeout: 5 * time.Minute}

// Fetch 下载并解析 db.text.json。"rows" 块是命名空间索引表，跳过。
func Fetch(ctx context.Context) (Database, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, DatabaseURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch ehtag database: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch ehtag database: unexpected status %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read ehtag database: %w", err)
	}

	var f databaseFile
	if err := json.Unmarshal(body, &f); err != nil {
		return nil, fmt.Errorf("parse ehtag database: %w", err)
	}

	db := Database{}
	for _, block := range f.Data {
		if block.Namespace == "rows" {
			continue
		}
		entries := make(map[string]Entry, len(block.Data))
		for key, e := range block.Data {
			entries[key] = Entry{
				NS:    block.Namespace,
				Key:   key,
				Name:  e.Name,
				Intro: e.Intro,
			}
		}
		db[block.Namespace] = entries
	}
	return db, nil
}
