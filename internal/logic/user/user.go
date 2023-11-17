package user

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"

	"mygf/internal/dao"
	"mygf/internal/model"
	"mygf/internal/model/do"
	"mygf/internal/service"
)

type (
	sUser struct{}
)

func init() {
	service.RegisterUser(New())
}

func New() service.IUser {
	return &sUser{}
}

func (user *sUser) Create(ctx context.Context, in model.UserCreateInput) (err error) {

	// 如果nickname为空，就默认为passport
	if in.Nickname == "" {
		in.Nickname = in.Passport
	}

	// 以下判断passport是否已经存在
	var (
		available bool
	)

	available, err = user.isPassportAvailable(ctx, in.Passport)
	if err != nil {
		return err
	}
	if !available {
		return gerror.Newf("Passport %s has ben taken by others.", in.Passport)
	}

	return dao.User.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {

		_, err = dao.User.Ctx(ctx).Data(do.User{
			Passport: in.Passport,
			Password: in.Password,
			Nickname: in.Nickname,
		}).Insert()

		return err
	})

	//return dao.User.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
	//		_, err = dao.User.Ctx(ctx).Data(do.User{
	//			Passport: in.Passport,
	//			Password: in.Password,
	//			Nickname: in.Nickname,
	//		}).Insert()
	//		return err
	//	})

}

func (user *sUser) isPassportAvailable(ctx context.Context, passport string) (bool, error) {

	count, err := dao.User.Ctx(ctx).Where(do.User{
		Passport: passport,
	}).Count()

	if err != nil {
		return false, err
	}

	return count == 0, nil
}
