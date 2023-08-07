package autoclick_test

import (
	"testing"
	"xyhelper-arkose/autoclick"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

func TestAutoClick(t *testing.T) {
	ctx := gctx.New()
	payload, token := autoclick.AutoClick()
	g.Log().Info(ctx, "payload: ", payload)
	g.Log().Info(ctx, "token: ", token)
}
