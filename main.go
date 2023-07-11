package main

import (
	// _ "xyhelper-arkose/ja3proxy"

	"context"
	"strings"
	"time"
	"xyhelper-arkose/autoclick"
	"xyhelper-arkose/config"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcron"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/gconv"
)

func main() {
	ctx := gctx.New()
	s := g.Server()
	_, err := gcron.AddSingleton(ctx, "* * * * * *", func(ctx context.Context) {
		// 开发模式跳过
		if g.Cfg().MustGetWithEnv(ctx, "MODE").String() == "dev" {
			return
		}
		g.Log().Print(ctx, "Every second", config.PayloadQueue.Size(), config.TokenQueue.Size())
		plaload, token := autoclick.AutoClick()
		if !strings.Contains(token, "sup=1|rid=") {
			g.Log().Error(ctx, "token error", token)
			return
		}
		Payload := config.Payload{
			Payload: plaload,
			Created: time.Now().Unix(),
		}
		Token := config.Token{
			Token:   token,
			Created: time.Now().Unix(),
		}
		config.PayloadQueue.Push(Payload)
		config.TokenQueue.Push(Token)
		// 生成延时
		time.Sleep(time.Duration(config.INTERVAL(ctx)) * time.Second)

	}, "get")
	if err != nil {
		panic(err)
	}
	// 每小时清理一次
	_, err = gcron.AddSingleton(ctx, "0 0 * * * *", func(ctx context.Context) {
		tokenURI := "http://127.0.0.1:" + gconv.String(config.Port) + "/token"
		payloadURI := "http://127.0.0.1:" + gconv.String(config.Port) + "/payload"
		g.Log().Print(ctx, "Every hour", tokenURI, payloadURI)
		g.Client().Get(ctx, tokenURI)
		g.Client().Get(ctx, payloadURI)
	}, "clean")
	if err != nil {
		panic(err)
	}

	s.SetPort(config.Port)
	s.SetServerRoot("resource/public/dist")
	s.BindHandler("/arkose", func(r *ghttp.Request) {
		payload, token := autoclick.AutoClick()
		r.Response.WriteJson(g.Map{
			"payload": payload,
			"token":   token,
		})
	})
	s.BindHandler("/payload", func(r *ghttp.Request) {
		ctx := r.Context()
		if r.Get("key").String() == "" {
			r.Response.WriteJson(g.Map{
				"code": 0,
				"msg":  "key is empty",
			})
			return
		}
		if !g.Cfg().MustGet(ctx, "key."+r.Get("key").String()).Bool() {
			r.Response.WriteJson(g.Map{
				"code": 0,
				"msg":  "key is error",
			})
			return
		}
		var payload interface{}
		if config.PayloadQueue.Size() == 0 {
			payloadStr, _ := autoclick.AutoClick()
			payload = config.Payload{
				Payload: payloadStr,
				Created: time.Now().Unix(),
			}
		} else {
			for config.PayloadQueue.Size() > 0 {
				payload = config.PayloadQueue.Pop()
				var payloadStuct config.Payload
				gconv.Struct(payload, &payloadStuct)
				if time.Now().Unix()-payloadStuct.Created < int64(config.PayloadExpire) {
					break
				} else {
					g.Log().Info(r.Context(), "payload is expired,will pop one ", config.PayloadQueue.Size(), payloadStuct.Created, config.PayloadExpire)
				}
			}
		}

		r.Response.WriteJson(payload)
	})
	s.BindHandler("/token", func(r *ghttp.Request) {
		ctx := r.Context()
		if r.Get("key").String() == "" {
			r.Response.WriteJson(g.Map{
				"code": 0,
				"msg":  "key is empty",
			})
			return
		}
		if !g.Cfg().MustGet(ctx, "key."+r.Get("key").String()).Bool() {
			r.Response.WriteJson(g.Map{
				"code": 0,
				"msg":  "key is error",
			})
			return
		}
		var token interface{}
		if config.TokenQueue.Size() == 0 {
			_, tokenStr := autoclick.AutoClick()
			token = config.Token{
				Token:   tokenStr,
				Created: time.Now().Unix(),
			}
		} else {
			for config.TokenQueue.Size() > 0 {
				token = config.TokenQueue.Pop()
				var tokenStuct config.Token
				gconv.Struct(token, &tokenStuct)
				if time.Now().Unix()-tokenStuct.Created < int64(config.TokenExpire) {
					break
				} else {
					g.Log().Info(r.Context(), "token is expired,will pop one ", config.TokenQueue.Size(), tokenStuct.Created, config.TokenExpire)
				}
			}
		}

		r.Response.WriteJson(token)
	})
	s.BindHandler("/pushtoken", func(r *ghttp.Request) {
		token := r.Get("token").String()
		if token == "" {
			r.Response.WriteJson(g.Map{
				"code": 0,
				"msg":  "token is empty",
			})
			return
		}
		if !strings.Contains(token, "sup=1|rid=") {
			g.Log().Error(ctx, "token error", token)
			r.Response.WriteJson(g.Map{
				"code": 0,
				"msg":  "token error",
			})
			return
		}
		Token := config.Token{
			Token:   token,
			Created: time.Now().Unix(),
		}
		config.TokenQueue.Push(Token)
		g.Log().Info(r.Context(), "pushtoken", token)
		r.Response.WriteJson(g.Map{
			"code": 1,
			"msg":  "success",
		})
	})
	s.Run()
}
