package autoclick

import (
	"context"
	"net/url"
	"strings"
	"xyhelper-arkose/config"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

func AutoClick() (payload, token string) {
	ctx := gctx.New()
	g.Log().Info(ctx, "开始获取 payload 和 token")
	// 获取 wsurl
	browserURL := config.BROWSERURL(ctx)
	g.Log().Info(ctx, "browserURL: ", browserURL)
	versionURL := browserURL + "/json/version"
	g.Log().Info(ctx, "versionURL: ", versionURL)
	versionResp := g.Client().GetVar(ctx, versionURL)
	g.Dump(versionResp)
	WSURL := gjson.New(versionResp).Get("webSocketDebuggerUrl").String()

	if WSURL == "" {
		g.Log().Error(ctx, "WSURL 为空")
		return
	}
	g.Log().Info(ctx, "WSURL: ", WSURL)
	allocCtx, _ := chromedp.NewRemoteAllocator(ctx, WSURL)
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	// ensure that the browser process is started
	if err := chromedp.Run(ctx); err != nil {
		panic(err)
	}
	// 捕获 POST 请求的内容
	chromedp.ListenTarget(ctx, func(ev interface{}) {
		if msg, ok := ev.(*network.EventRequestWillBeSent); ok {
			// g.Log().Info(ctx, msg)
			if msg.Type == "XHR" && msg.Request.Method == "POST" {
				if msg.Request.URL == "https://tcr9i.chat.openai.com/fc/gt2/public_key/35536E1E-65B4-4D96-9D97-6ADB7EFF8147" {
					payload = msg.Request.PostData
					// g.Log().Info(ctx, msg.Request.URL, msg.Request.PostData, msg.Request.Headers)
					// postMap, err := convertToMap(msg.Request.PostData)
					// if err != nil {
					// 	g.Log().Error(ctx, err)
					// } else {
					// 	g.Dump(postMap)
					// 	postJson := gjson.New(postMap)
					// 	g.Dump(postJson)
					// 	g.Dump(msg.Request.Headers)

					// }
				}
				if msg.Request.URL == "https://client-api.arkoselabs.com/fc/gt2/public_key/35536E1E-65B4-4D96-9D97-6ADB7EFF8147" {
					payload = msg.Request.PostData
				}

			}
		}
	})
	// 设置代理为 http://127.0.0.1:7006
	// chromedp.ProxyServer("http://127.0.0.1:7006")
	chromedp.Run(ctx, network.Enable(), chromedp.ActionFunc(func(ctx context.Context) error {
		// 打开 http://localhost:8199
		siteUrl := "https://chat.openai.com/" + config.PageName
		g.Log().Info(ctx, "打开 ", siteUrl)
		chromedp.Navigate(siteUrl).Do(ctx)

		// 等待 <button id="enforcement-trigger">...</button> 出现
		g.Log().Info(ctx, "等待 <button id=\"enforcement-trigger\">...</button> 出现")
		chromedp.WaitVisible(`#enforcement-trigger`, chromedp.ByID).Do(ctx)
		g.Log().Info(ctx, "等待 <button id=\"enforcement-trigger\">...</button> 出现 完成")
		// 等待 1 秒
		chromedp.Sleep(1 * 1000 * 1000 * 1000).Do(ctx)
		// 点击 <button id="enforcement-trigger">...</button>
		chromedp.Click(`#enforcement-trigger`, chromedp.ByID).Do(ctx)
		// 获取 <div id="token"></div> 的文本内容
		// 等待 1 秒
		i := 0
		for token == "" && i < 10 {
			i++
			g.Log().Info(ctx, "等待 <div id=\"token\"></div> 的文本内容", i)
			chromedp.Sleep(3 * 1000 * 1000 * 1000).Do(ctx)
			chromedp.Text(`#token`, &token, chromedp.ByID).Do(ctx)
		}

		// g.Log().Info(ctx, token)

		// select {}
		return nil
	}))
	// g.Log().Info(ctx, "payload", payload)
	// g.Log().Info(ctx, "token", token)
	g.Log().Info(ctx, "获取 payload 和 token 完成")
	return
}

func convertToMap(data string) (*gmap.StrStrMap, error) {
	m := gmap.NewStrStrMap(true)

	pairs := strings.Split(data, "&")
	for _, pair := range pairs {
		keyValue := strings.Split(pair, "=")
		key := keyValue[0]
		value, err := url.PathUnescape(keyValue[1])
		if err != nil {
			// Handle error if necessary
			return nil, err
		}
		// g.Log().Info(context.Background(), "key-value", key, value)
		m.Set(key, value)
	}

	return m, nil
}
