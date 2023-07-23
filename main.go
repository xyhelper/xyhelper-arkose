package main

import (
	// _ "xyhelper-arkose/ja3proxy"

	"context"
	"strings"
	"time"
	"xyhelper-arkose/config"

	"github.com/gogf/gf/v2/encoding/gjson"
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
	s.EnableHTTPS("./resource/certs/server.crt", "./resource/certs/server.key")
	s.SetHTTPSPort(443)

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
			// bypass https://ai.fakeopen.com/api/arkose/token
			result := g.Client().Proxy(config.Proxy).GetVar(ctx, "https://ai.fakeopen.com/api/arkose/token")
			token = gjson.New(result).Get("token").String()
			if token == "" {
				r.Response.WriteJson(g.Map{
					"code": 0,
					"msg":  "token is empty",
				})
				return
			} else {

				r.Response.WriteJson(g.Map{
					"token":   token,
					"created": time.Now().Unix(),
				})
				g.Log().Info(ctx, "get token from ai.fakeopen.com", token)
				return
			}

		} else {
			for config.TokenQueue.Size() > 0 {
				token = config.TokenQueue.Pop()
				var tokenStuct config.Token
				gconv.Struct(token, &tokenStuct)
				if time.Now().Unix()-tokenStuct.Created < int64(config.TokenExpire) {
					break
				} else {
					g.Log().Info(ctx, "token is expired,will pop one ", config.TokenQueue.Size(), tokenStuct.Created, config.TokenExpire)
					token = nil
				}
			}
		}
		if token == nil {
			r.Response.WriteJson(g.Map{
				"code": 0,
				"msg":  "token is empty",
			})
			return
		}
		r.Response.WriteJson(token)
	})
	s.BindHandler("/pushtoken", func(r *ghttp.Request) {
		// g.Dump(r.Header)
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
		forwordURL := g.Cfg().MustGetWithEnv(ctx, "FORWORD_URL").String()
		g.Log().Info(ctx, "forwordURL", forwordURL)

		if forwordURL != "" {
			result := g.Client().Proxy(config.Proxy).PostVar(ctx, forwordURL, g.Map{
				"token": token,
			})
			g.Log().Info(ctx, getRealIP(r), "forwordURL", forwordURL, result)
			r.Response.WriteJson(g.Map{
				"code":       1,
				"msg":        "success",
				"forwordURL": forwordURL,
			})
			return
		}
		Token := config.Token{
			Token:   token,
			Created: time.Now().Unix(),
		}
		config.TokenQueue.Push(Token)
		g.Log().Info(r.Context(), getRealIP(r), "pushtoken", token)
		r.Response.WriteJson(g.Map{
			"code": 1,
			"msg":  "success",
		})
	})
	s.Run()
}

func getRealIP(req *ghttp.Request) string {
	// 优先获取Cf-Connecting-Ip
	if ip := req.Header.Get("Cf-Connecting-Ip"); ip != "" {
		return ip
	}

	// 优先获取X-Real-IP
	if ip := req.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	// 其次获取X-Forwarded-For
	if ip := req.Header.Get("X-Forwarded-For"); ip != "" {
		return ip
	}
	// 最后获取RemoteAddr
	ip := req.RemoteAddr
	// 处理端口
	if index := strings.Index(ip, ":"); index != -1 {
		ip = ip[0:index]
	}
	return ip
}
