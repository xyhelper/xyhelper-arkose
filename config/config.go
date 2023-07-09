package config

import (
	"github.com/gogf/gf/v2/container/gqueue"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

var (
	PayloadExpire = 60 * 60 * 3
	TokenExpire   = 60 * 20
	PayloadQueue  = gqueue.New()
	TokenQueue    = gqueue.New()
	Port          = 8199
)

type Payload struct {
	Payload string `json:"payload"`
	Created int64  `json:"created"`
}

type Token struct {
	Token   string `json:"token"`
	Created int64  `json:"created"`
}

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
	payLoadExpire := g.Cfg().MustGetWithEnv(ctx, "PAYLOAD_EXPIRE").Int()
	if payLoadExpire != 0 {
		PayloadExpire = payLoadExpire
	}
	tokenExpire := g.Cfg().MustGetWithEnv(ctx, "TOKEN_EXPIRE").Int()
	if tokenExpire != 0 {
		TokenExpire = tokenExpire
	}
}
