package utils

import (
	"sync"
	"time"
)

// 全局按 host 串行的下载队列：
// 多个游戏刮削任务并发时，同一站点（host）的图片下载在此排队串行执行，
// 避免同站并发请求被限流。不同 host 之间互不影响。
var (
	dlMu        sync.Mutex
	dlQueues    = map[string]chan func(){}
	dlInterval  = 50 * time.Millisecond
	dlQueueSize = 1024
)

// DownloadDo 将 task 提交到 host 对应的串行队列并阻塞等待执行完成。
// host 应为 url.Parse(url).Host（如 "t.vndb.org"）。
func DownloadDo(host string, task func()) {
	if host == "" || task == nil {
		return
	}
	q := dlQueue(host)
	done := make(chan struct{})
	q <- func() {
		defer close(done)
		task()
	}
	<-done
}

// dlQueue 返回 host 专属的任务队列，首次使用时创建并启动消费 goroutine，
// 每个任务之间间隔 dlInterval。
func dlQueue(host string) chan func() {
	dlMu.Lock()
	defer dlMu.Unlock()
	if q, ok := dlQueues[host]; ok {
		return q
	}
	q := make(chan func(), dlQueueSize)
	dlQueues[host] = q
	go func() {
		for t := range q {
			t()
			time.Sleep(dlInterval)
		}
	}()
	return q
}
