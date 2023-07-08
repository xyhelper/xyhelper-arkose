package config

import (
	"github.com/gogf/gf/v2/container/gqueue"
	"github.com/gogf/gf/v2/frame/g"
)

var (
	PayloadQueueSize = 1000
	TokenQueueSize   = 180
	PayloadQueue     = gqueue.New(PayloadQueueSize)
	TokenQueue       = gqueue.New(TokenQueueSize)
)

func BROWSERURL(ctx g.Ctx) string {
	BROWSERURL := g.Cfg().MustGetWithEnv(ctx, "BROWSERURL").String()
	// g.Log().Infof(ctx, "BROWSERURL: %s", BROWSERURL)

	return BROWSERURL
}
