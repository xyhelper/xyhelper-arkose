package api

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
)

func Index(r *ghttp.Request) {
	ctx := r.Context()
	trustHost := r.Host == "localhost:3000"
	g.Log().Info(ctx, r.Host, trustHost)
	r.Cookie.Set("uid", gtime.Now().String())
	if trustHost {
		r.Response.RedirectTo("/xyhelper/")
		return
	}

	r.Response.Write("Hello Xyhelper-Arkose!")
}
