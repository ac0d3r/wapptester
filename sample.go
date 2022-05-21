package wapptester

import (
	"bytes"
	"context"
	"crypto/tls"
	"io"
	"net/http"
	"net/http/httputil"
	"strings"
	"sync"

	"github.com/gokitx/pkgs/bytesconv"
	"golang.org/x/net/html"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/transform"
)

const (
	userAgent = "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.0 Safari/537.36"
)

type Sample struct {
	URL        string            `wapptester:"url"`
	StatusCode float64           `wapptester:"status"`
	Header     string            `wapptester:"header"`
	Headers    map[string]string `wapptester:"headers"`
	Cookie     string            `wapptester:"cookie"`
	Cookies    map[string]string `wapptester:"cookies"`
	Server     string            `wapptester:"server"`
	Title      string            `wapptester:"title"`
	Body       string            `wapptester:"body"`
	Meta       map[string]string `wapptester:"meta"`
	Hash       string            `wapptester:"hash"`
	HashMMH3   string            `wapptester:"hashmmh3"`
}

var (
	httpclient *http.Client = &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives:   true,
			DisableCompression:  true,
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
			MaxIdleConns:        100,
			MaxConnsPerHost:     100,
			MaxIdleConnsPerHost: 100,
		},
	}
	bufferPool = sync.Pool{
		New: func() any {
			return bytes.NewBuffer(make([]byte, 1024))
		},
	}
)

func MakeSample(ctx context.Context, URL string) (*Sample, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, URL, nil)
	if err != nil {
		return nil, err
	}
	req.Close = true
	req.Header.Set("User-Agent", userAgent)

	resp, err := httpclient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var (
		transformFlag bool
		reader        io.Reader = resp.Body
	)
	contentType := resp.Header.Get("Content-Type")
	if ct := strings.ToUpper(contentType); strings.Contains(ct, "GB2312") || strings.Contains(ct, "GBK") {
		transformFlag = true
		reader = transform.NewReader(resp.Body, simplifiedchinese.GBK.NewDecoder())
	} else if strings.Contains(ct, "BIG5") {
		transformFlag = true
		reader = transform.NewReader(resp.Body, traditionalchinese.Big5.NewDecoder())
	}

	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer func() {
		buf.Reset()
		bufferPool.Put(buf)
	}()
	if _, err := io.Copy(buf, reader); err != nil {
		return nil, err
	}
	body := buf.String()

	if !transformFlag &&
		!strings.Contains(strings.ToLower(contentType), "utf-8") {
		if lowerBody := strings.ToLower(body); strings.Contains(lowerBody, "gb2312\"") || strings.Contains(lowerBody, "gbk\"") {
			if nBody, _, err := transform.String(simplifiedchinese.GBK.NewDecoder(), body); err == nil {
				body = nBody
			}
		}
	}

	headers := make(map[string]string)
	for key := range resp.Header {
		v := resp.Header.Get(key)
		headers[key] = v
	}
	cookie := make([]string, 0)
	cookies := make(map[string]string)
	for _, v := range resp.Cookies() {
		cookie = append(cookie, v.String())
		cookies[v.Name] = v.Value
	}
	headersText, _ := httputil.DumpResponse(resp, false)
	return &Sample{
		URL:        URL,
		StatusCode: float64(resp.StatusCode),
		Header:     bytesconv.BytesToString(headersText),
		Headers:    headers,
		Cookie:     strings.Join(cookie, "\n"),
		Cookies:    cookies,
		Server:     resp.Header.Get("Server"),
		Title:      extraTitle(body),
		Body:       body,
		Meta:       getMeta(body),
		Hash:       Md5(body),
		HashMMH3:   MMH3(body),
	}, nil
}

const (
	titleTagBegin = `<title>`
	titleTagEnd   = `</title>`
)

func extraTitle(body string) string {
	low := strings.ToLower(body)
	if len(low) != len(body) {
		return ""
	}
	begin := strings.Index(low, titleTagBegin)
	if begin < 0 {
		return ""
	}
	begin += len(titleTagBegin)
	end := strings.Index(low, titleTagEnd)
	if end < 0 {
		return ""
	}
	if begin >= end {
		return ""
	}
	return body[begin:end]
}

func getMeta(body string) map[string]string {
	result := make(map[string]string)
	doc, err := html.Parse(strings.NewReader(body))
	if err != nil {
		return nil
	}

	var f func(*html.Node)
	f = func(n *html.Node) {
		switch n.Type {
		case html.ElementNode:
			switch n.Data {
			case "meta":
				k := ""
				for _, a := range n.Attr {
					if a.Key == "name" {
						k = a.Val
					}
					if a.Key == "content" {
						if k == "" {
							continue
						}
						result[k] = a.Val
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(doc)

	return result
}
