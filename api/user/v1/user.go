package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

type SignupReq struct {
	g.Meta    `path:"/user/sign-up" tags:"User" method:"post" summary:"You first user api"`
	Nickname  string `v:"required"`
	Passport  string `v:"required"`
	Password  string `v:"required"`
	Password2 string `v:"required|length:6,16|same:Password"`
}

type SignupRes struct {
	//g.Meta `mime:"text/html" example:"string"`
	Nickname string
	Passport string
}

type SigninReq struct {
	g.Meta   `path:"/user/sign-in" tags:"User" method:"post"`
	Passport string `v:"required"`
	Password string `v:"required"`
}

type SigninRes struct {
}
