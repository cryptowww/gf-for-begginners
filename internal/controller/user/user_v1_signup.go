package user

import (
	"context"

	v1 "mygf/api/user/v1"
	"mygf/internal/model"
	"mygf/internal/service"
)

func (c *ControllerV1) Signup(ctx context.Context, req *v1.SignupReq) (res *v1.SignupRes, err error) {
	//return nil, gerror.NewCode(gcode.CodeNotImplemented)
	err = service.User().Create(ctx, model.UserCreateInput{
		Passport: req.Passport,
		Password: req.Password,
		Nickname: req.Nickname,
	})
	res = &v1.SignupRes{
		Nickname: req.Nickname,
		Passport: req.Passport,
	}

	return
}
