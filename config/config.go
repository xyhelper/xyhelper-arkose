package config

import (
	"github.com/gogf/gf/v2/container/gqueue"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

var (
	PayloadQueueSize = 1000
	TokenQueueSize   = 180
	PayloadQueue     = gqueue.New(PayloadQueueSize)
	TokenQueue       = gqueue.New(TokenQueueSize)
	Port             = 8199
)

func BROWSERURL(ctx g.Ctx) string {
	BROWSERURL := g.Cfg().MustGetWithEnv(ctx, "BROWSERURL").String()
	// g.Log().Infof(ctx, "BROWSERURL: %s", BROWSERURL)

	return BROWSERURL
}

func INTERVAL(ctx g.Ctx) int {
	INTERVAL := g.Cfg().MustGetWithEnv(ctx, "INTERVAL").Int()
	g.Log().Infof(ctx, "INTERVAL: %d", INTERVAL)

	return INTERVAL
}

func init() {
	ctx := gctx.GetInitCtx()
	port := g.Cfg().MustGetWithEnv(ctx, "PORT").Int()
	if port != 0 {
		Port = port
	}
}
