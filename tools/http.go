package tools

import (
	"context"
	"fmt"
	"hermes/config"
	"io"
	"net"
	"net/http"
	"time"

	zaplog "github.com/dokidokikoi/go-common/log/zap"
	"github.com/dokidokikoi/go-common/tools"
	"go.uber.org/zap"
	"golang.org/x/net/proxy"
)

func MakeRequest(
	method, uri string,
	proxy config.ProxyConfig,
	body io.Reader,
	header map[string]string,
	cookies []*http.Cookie,
	retryCnt int) (data []byte, status int, err error) {

	client := &http.Client{}
	// 构建请求客户端
	if proxy.Host != "" && proxy.Port != 0 {
		p, err := tools.Socks5Proxy(fmt.Sprintf("%s:%d", proxy.Host, proxy.Port), proxy.Username, proxy.Password)
		if err != nil {
			zaplog.L().Error("proxy error", zap.Error(err))
		}
		client = createHTTPClient(p)
	}

	// 创建请求对象
	req, err := createRequest(method, uri, body, header, cookies)
	// 检查错误
	if err != nil {
		return nil, 0, err
	}

	var retry int = 0
	var res *http.Response
	for retry <= retryCnt {
		// 执行请求
		res, err = client.Do(req)
		// 检查错误
		if err == nil {
			break
		}
		zaplog.L().Warn("request failed", zap.String("url", uri), zap.Error(err), zap.Int("retry", retry))
		retry++
	}
	if err != nil {
		return nil, 0, fmt.Errorf("%s [Request]: %s", uri, err)
	}

	// 获取请求状态码
	status = res.StatusCode
	// 读取请求内容
	data, err = io.ReadAll(res.Body)
	// 关闭请求连接
	_ = res.Body.Close()

	return data, status, err
}

// 创建http客户端
func createHTTPClient(dialer proxy.Dialer) *http.Client {
	var transport http.RoundTripper
	if dialer != nil {
		transport = &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return dialer.Dial(network, addr)
			},
		}
	}

	// 返回客户端
	return &http.Client{
		Transport: transport,
		Timeout:   60 * time.Second,
	}
}

// 创建请求对象
func createRequest(method, uri string, body io.Reader, header map[string]string, cookies []*http.Cookie) (*http.Request, error) {
	// 新建请求
	req, err := http.NewRequest(method, uri, body)
	// 检查错误
	if err != nil {
		return nil, fmt.Errorf("%s [Request]: %s", uri, err)
	}

	// 循环头部信息
	for k, v := range header {
		// 设置头部
		req.Header.Set(k, v)
	}

	// 设置了cookie
	if len(cookies) > 0 {
		// 循环cookie
		for _, cookie := range cookies {
			// 加入cookie
			req.AddCookie(cookie)
		}
	}

	return req, err
}
