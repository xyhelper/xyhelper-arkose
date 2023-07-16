package main

import (
	// _ "xyhelper-arkose/ja3proxy"

	"context"
	"strings"
	"time"
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
	// 每小时清理一次
	_, err := gcron.AddSingleton(ctx, "0 0 * * * *", func(ctx context.Context) {
		tokenURI := "http://127.0.0.1:" + gconv.String(config.Port) + "/token"
		g.Log().Print(ctx, "Every hour", tokenURI)
		g.Client().Get(ctx, tokenURI)
	}, "clean")
	if err != nil {
		panic(err)
	}

	s.SetPort(config.Port)
	s.SetServerRoot("resource/public")
	s.BindHandler("/ping", func(r *ghttp.Request) {
		total := config.TokenQueue.Size()
		r.Response.WriteJson(g.Map{
			"code":  1,
			"msg":   "pong",
			"total": total,
		})

	})

	s.BindHandler("/token", func(r *ghttp.Request) {
		ctx := r.Context()

		var token interface{}
		if config.TokenQueue.Size() == 0 {
			r.Response.WriteJson(g.Map{
				"code": 0,
				"msg":  "token is empty",
			})
			return

		} else {
			for config.TokenQueue.Size() > 0 {
				token = config.TokenQueue.Pop()
				var tokenStuct config.Token
				gconv.Struct(token, &tokenStuct)
				if time.Now().Unix()-tokenStuct.Created < int64(config.TokenExpire) {
					break
				} else {
					g.Log().Info(ctx, "token is expired,will pop one ", config.TokenQueue.Size(), tokenStuct.Created, config.TokenExpire)
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
