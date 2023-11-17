// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. 
// =================================================================================

package user

import (
	"context"
	
	"mygf/api/user/v1"
)

type IUserV1 interface {
	Signup(ctx context.Context, req *v1.SignupReq) (res *v1.SignupRes, err error)
	Signin(ctx context.Context, req *v1.SigninReq) (res *v1.SigninRes, err error)
}


