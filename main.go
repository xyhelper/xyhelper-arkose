package main

import (
	// _ "xyhelper-arkose/ja3proxy"

	"context"
	"strings"
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
		if config.PayloadQueue.Size() >= gconv.Int64(config.PayloadQueueSize) {
			g.Log().Info(ctx, "PayloadQueue is full,will pop one ", config.PayloadQueue.Size(), config.PayloadQueueSize)
			config.PayloadQueue.Pop()
		}
		if config.TokenQueue.Size() >= gconv.Int64(config.TokenQueueSize) {
			g.Log().Info(ctx, "TokenQueue is full,will pop one ", config.TokenQueue.Size(), config.TokenQueueSize)
			config.TokenQueue.Pop()
		}

		g.Log().Print(ctx, "Every second", config.PayloadQueue.Size(), config.TokenQueue.Size())
		playload, token := autoclick.AutoClick()
		if !strings.Contains(token, "sup=1|rid=") {
			g.Log().Error(ctx, "token error", token)
			return
		}
		config.PayloadQueue.Push(playload)
		config.TokenQueue.Push(token)

	}, "get")
	if err != nil {
		panic(err)
	}
	s.SetPort(8199)
	s.SetServerRoot("resource/public")
	s.BindHandler("/arkose", func(r *ghttp.Request) {
		payload, token := autoclick.AutoClick()
		r.Response.WriteJson(g.Map{
			"payload": payload,
			"token":   token,
		})
	})
	s.BindHandler("/payload", func(r *ghttp.Request) {
		payload := config.PayloadQueue.Pop()
		if payload == nil {
			payload, _ = autoclick.AutoClick()
		}
		r.Response.WriteJson(g.Map{
			"payload": payload,
		})
	})
	s.BindHandler("/token", func(r *ghttp.Request) {
		token := config.TokenQueue.Pop()
		if token == nil {
			_, token = autoclick.AutoClick()
		}
		r.Response.WriteJson(g.Map{
			"token": token,
		})
	})
	s.Run()
}
