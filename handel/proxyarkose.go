package handel

import (
	"bytes"
	"crypto/tls"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
	"xyhelper-arkose/config"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

var (
	UpStream = "https://client-api.arkoselabs.com/"
)

func init() {

}

func Proxy(r *ghttp.Request) {
	ctx := r.Context()
	trutHost := r.Host == "localhost:3000"
	payload := &config.Payload{
		Payload: "",
		Created: time.Now().Unix(),
	}
	u, _ := url.Parse(UpStream)
	var proxy *httputil.ReverseProxy

	proxy = &httputil.ReverseProxy{}
	// g.Dump(config.PROXY(ctx))
	if config.PROXY(ctx).String() != "" {
		proxy.Transport = &http.Transport{
			Proxy: http.ProxyURL(config.PROXY(ctx)),
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}
	r.Header.Set("Origin", "https://client-api.arkoselabs.com")
	r.Header.Set("Referer", "https://client-api.arkoselabs.com/v2/1.5.4/enforcement.cd12da708fe6cbe6e068918c38de2ad9.html")
	r.Header.Del("Cf-Connecting-Ip")
	r.Header.Del("Cf-Ipcountry")
	r.Header.Del("Cf-Ray")
	r.Header.Del("Cf-Request-Id")
	r.Header.Del("Cf-Visitor")
	r.Header.Del("Cf-Warp-Tag-Id")
	r.Header.Del("Cf-Worker")
	r.Header.Del("Cf-Device-Type")
	r.Header.Del("Cf-Request-Id")
	r.Header.Del("X-Forwarded-Host")
	r.Header.Del("X-Forwarded-Proto")
	r.Header.Del("X-Forwarded-For")
	r.Header.Del("X-Forwarded-Port")
	r.Header.Del("X-Forwarded-Server")
	r.Header.Del("X-Real-Ip")
	r.Header.Del("Accept-Encoding")
	// requrl := r.Request.URL.Path

	// if requrl == "/fc/gt2/public_key/35536E1E-65B4-4D96-9D97-6ADB7EFF8147" {
	// 	body := r.GetBodyString()
	// 	bodyArray := gstr.Split(body, "&")
	// 	g.Dump(bodyArray)
	// 	// 遍历数组 当数组元素以 "site=http" 开头时，将其替换为 "site=http%3A%2F%2Flocalhost%3A3000"
	// 	for i, v := range bodyArray {
	// 		if gstr.HasPrefix(v, "site=http") {
	// 			bodyArray[i] = "site=http%3A%2F%2Flocalhost%3A3000"
	// 		}
	// 	}
	// 	body = gstr.Join(bodyArray, "&")

	// 	payload.Payload = body
	// 	payload.UserAgent = r.Header.Get("User-Agent")
	// 	r.Body = io.NopCloser(bytes.NewReader(gconv.Bytes(body)))
	// 	r.ContentLength = int64(len(body))
	// }

	proxy.Rewrite = func(proxyRequest *httputil.ProxyRequest) {
		// g.Dump(proxyRequest)
		proxyRequest.SetURL(u)
	}

	proxy.ModifyResponse = func(resp *http.Response) error {
		cookieStr := resp.Header.Get("Set-Cookie")
		// 移除域名限制
		cookieStr = strings.Replace(cookieStr, "Domain=.arkoselabs.com;", "", -1)
		// 重写cookie
		resp.Header.Set("Set-Cookie", cookieStr)

		// 解码 url
		if resp.StatusCode <= 400 {
			g.Log().Info(r.Context(), resp.StatusCode, resp.Request.URL.Path)
			// 获取返回的body
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			g.Log().Info(r.Context(), "body", string(body))
			// 解压缩body

			// g.Dump(string(unzipbody))
			token := gjson.New(body).Get("token").String()
			g.Log().Info(r.Context(), "token", token)
			if strings.Contains(token, "sup=1|rid=") && trutHost {
				// 获取请求的body
				err := config.Cache.Set(r.Context(), "payload", payload, 0)
				if err != nil {
					return err
				}
				g.Log().Info(r.Context(), "refresh payload cache", payload)

			}
			// 将原始body 返回
			resp.Body = io.NopCloser(bytes.NewReader(body))
		} else {
			g.Log().Warning(r.Context(), resp.StatusCode, resp.Request.URL.Path)

		}
		return nil
	}
	g.Dump(r.Header)

	proxy.ServeHTTP(r.Response.RawWriter(), r.Request)

}
