package user

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"mygf/api/user/v1"
)

func (c *ControllerV1) Signin(ctx context.Context, req *v1.SigninReq) (res *v1.SigninRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
