package tools

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/net/proxy"
	"resty.dev/v3"
)

var client *resty.Client = resty.New().SetRetryCount(3)

type Option func(*resty.Request) error

func Req(method, url string, body any, options ...Option) (*resty.Response, error) {
	req := client.R().SetBody(body)
	for _, o := range options {
		o(req)
	}
	return req.Execute(method, url)
}

func ReqWithProxy(method, url string, body any, proxy string, options ...Option) (*resty.Response, error) {
	proxy = strings.TrimSpace(proxy)
	if proxy == "" {
		return Req(method, url, body, options...)
	}
	client := resty.New().SetRetryCount(3).SetProxy(proxy)
	req := client.R().SetBody(body)
	for _, o := range options {
		o(req)
	}
	return req.Execute(method, url)
}

func SetHeadersWithOption(headers map[string]string) Option {
	return func(r *resty.Request) error {
		r.SetHeaders(headers)
		return nil
	}
}

func SetQueryParamsWithOption(params map[string]string) Option {
	return func(r *resty.Request) error {
		r.SetQueryParams(params)
		return nil
	}
}

func SetCookiesWithOption(cookies ...*http.Cookie) Option {
	return func(r *resty.Request) error {
		r.SetCookies(cookies)
		return nil
	}
}

func SetFromWithOption(data map[string]string) Option {
	return func(r *resty.Request) error {
		r.SetFormData(data)
		return nil
	}
}

func SetMultipartWithOption(fields ...*resty.MultipartField) Option {
	return func(r *resty.Request) error {
		r.SetMultipartFields(fields...)
		return nil
	}
}

func SetMultiFileWithOption(params map[string]string, files map[string][]string) Option {
	return func(r *resty.Request) error {
		buf := bytes.NewBuffer([]byte{})
		writer := multipart.NewWriter(buf)
		for k, v := range params {
			w, err := writer.CreateFormField(k)
			if err != nil {
				return err
			}
			_, err = w.Write([]byte(v))
			if err != nil {
				return err
			}
		}
		for k, v := range files {
			for _, vv := range v {
				_, filename := filepath.Split(vv)
				h := make(textproto.MIMEHeader)
				h.Set("Content-Disposition",
					fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
						escapeQuotes(k), escapeQuotes(filename)))
				h.Set("Content-Type", "application/octet-stream")
				w, err := writer.CreatePart(h)
				if err != nil {
					return err
				}
				f, err := os.Open(vv)
				if err != nil {
					return err
				}
				_, err = io.Copy(w, f)
				f.Close()
				if err != nil {
					return err
				}
			}
		}
		r.SetBody(buf)
		r.SetHeader("Content-Type", writer.FormDataContentType())
		return nil
	}
}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

func SetMultipartFormWithOption(params map[string]string) Option {
	return func(r *resty.Request) error {
		r.SetMultipartFormData(params)
		return nil
	}
}

// 创建http客户端
func createHTTPClient(dialer proxy.Dialer) *http.Client {
	transport := http.DefaultTransport
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
