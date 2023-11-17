# GoFrame Template For SingleRepo

Quick Start: 

- https://goframe.org/pages/viewpage.action?pageId=1114399

## init

```bash
gf init mygf
```

## Create Table

```sql
CREATE TABLE `user`
(
    `id`        int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'User ID',
    `passport`  varchar(45) NOT NULL COMMENT 'User Passport',
    `password`  varchar(45) NOT NULL COMMENT 'User Password',
    `nickname`  varchar(45) NOT NULL COMMENT 'User Nickname',
    `create_at` datetime DEFAULT NULL COMMENT 'Created Time',
    `update_at` datetime DEFAULT NULL COMMENT 'Updated Time',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
```

## Gen Dao

add the following config in hack\config.yaml

```yaml
# 增加数据库连接
gfcli:
  gen:
    dao:
      - link: "mysql:root:12345@tcp(127.0.0.1:3306)/test"
        tables: "user"
        removePrefix: "gf_"
        descriptionTag: true
        noModelComment: true
```

execute the gen command

```yaml
(base) PS D:\workspace\mygf> gf gen dao
generated: internal\dao\user.go
generated: internal\dao\internal\user.go
generated: internal\model/do\user.go
generated: internal\model/entity\user.go
done!
```

## Write Logic Code

### Logic Code Input Params(model/user.go)

```go
package model

type UserCreateInput struct {
	Passport string
	Password string
	Nickname string
}
```

### Logic Code(logic/user/user.go)

```go
package user

import (
	"context"

	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"

	"mygf/internal/dao"
	"mygf/internal/model"
	"mygf/internal/model/do"
)

type (
	sUser struct{}
)

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
}
// 判读用户是否已存在
func (user *sUser) isPassportAvailable(ctx context.Context, passport string) (bool, error) {

	count, err := dao.User.Ctx(ctx).Where(do.User{
		Passport: passport,
	}).Count()

	if err != nil {
		return false, err
	}

	return count == 0, nil
}
```

## Gen Service Interface

```powershell
(base) PS D:\workspace\mygf> gf gen service
generating service go file: internal/service\user.go
generating init go file: internal/logic\logic.go
gofmt go files in "internal/service"
update main.go
done!
```

### 给Logic Code增加依赖注入代码(logic/user/user.go)

```go
type (
	sUser struct{}
)
...
// 增加以下代码
func init() {
	service.RegisterUser(New())
}

func New() service.IUser {
	return &sUser{}
}
...
```

## Gen Controller

mkdir api/user/v1,新建文件user.go

```go
package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)
// 前端传参
type SignupReq struct {
	g.Meta    `path:"/user/sign-up" tags:"User" method:"post" summary:"You first user api"`
	Nickname  string `v:"required"`
	Passport  string `v:"required"`
	Password  string `v:"required"`
	Password2 string `v:"required|length:6,16|same:Password"`
}
// 返回值
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
```

```powershell
(base) PS D:\workspace\mygf> gf gen ctrl
generated: D:\workspace\mygf\api\user\user.go
generated: internal/controller\user\user.go
generated: internal/controller\user\user_new.go
generated: internal/controller\user\user_v1_signup.go
generated: internal/controller\user\user_v1_signin.go
done!
```

## controller调用service服务

修改internal/controller/user/user_v1_signup.go

```go
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
```

## 增加路由(cmd.go)

```go
package cmd

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
	// 增加
	"mygf/internal/controller/user"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			s := g.Server()
			s.Group("/", func(group *ghttp.RouterGroup) {
				group.Middleware(ghttp.MiddlewareHandlerResponse)
				group.Bind(
					//增加
					user.NewV1(),
				)
			})
			s.Run()
			return nil
		},
	}
)
```

## 启动服务前的配置

- locig/user/user.go引入文件修改

```go
import "github.com/gogf/gf/database/gdb"
->
import "github.com/gogf/gf/v2/database/gdb"
```

- main.go增加

```go
import _ "github.com/gogf/gf/contrib/drivers/mysql/v2"
```

- 安装包

```powershell
go get github.com/gogf/gf/v2/database/gdb
go get github.com/gogf/gf/contrib/drivers/mysql/v2
```

- 配置数据库连接 manifest/config/config.yaml

```yaml
# Database.
database:
  logger:
    level: "all"
    stdout: true

  default:
    link: "mysql:root:12345@tcp(127.0.0.1:3306)/test"
    debug: true
```

- docker启动数据库

## 启动应用

```yaml
gf run main.go
```

```coq
(base) PS D:\workspace\mygf> gf run main.go
build: main.go
go build -o ./\main.exe  main.go
./\main.exe
build running pid: 5340
2023-11-17 14:53:26.712 [INFO] pid[5340]: http server started listening on [:8000]
2023-11-17 14:53:26.712 [INFO] {e0c4536e83569817d652071f6b5268e0} swagger ui is serving at address: http://127.0.0.1:8000/swagger/
2023-11-17 14:53:26.713 [INFO] {e0c4536e83569817d652071f6b5268e0} openapi specification is serving at address: http://127.0.0.1:8000/api.json

  ADDRESS | METHOD |     ROUTE     |                             HANDLER                             |           MIDDLEWARE
----------|--------|---------------|-----------------------------------------------------------------|----------------------------------
  :8000   | ALL    | /*            | github.com/gogf/gf/v2/net/ghttp.internalMiddlewareServerTracing | GLOBAL MIDDLEWARE
----------|--------|---------------|-----------------------------------------------------------------|----------------------------------
  :8000   | ALL    | /api.json     | github.com/gogf/gf/v2/net/ghttp.(*Server).openapiSpec           |
----------|--------|---------------|-----------------------------------------------------------------|----------------------------------
  :8000   | ALL    | /swagger/*    | github.com/gogf/gf/v2/net/ghttp.(*Server).swaggerUI             | HOOK_BEFORE_SERVE
----------|--------|---------------|-----------------------------------------------------------------|----------------------------------
  :8000   | POST   | /user/sign-in | mygf/internal/controller/user.(*ControllerV1).Signin            | ghttp.MiddlewareHandlerResponse
----------|--------|---------------|-----------------------------------------------------------------|----------------------------------
  :8000   | POST   | /user/sign-up | mygf/internal/controller/user.(*ControllerV1).Signup            | ghttp.MiddlewareHandlerResponse
----------|--------|---------------|-----------------------------------------------------------------|----------------------------------
```

## 调用服务测试

http://localhost:8000/user/sign-up

![Untitled](https://prod-files-secure.s3.us-west-2.amazonaws.com/e71874ab-0dc0-4662-9caf-5af461baee23/e7037d57-4cd7-460f-9a56-81571f575cf3/Untitled.png)

**查询数据库**

```coq
mysql> use test;
Reading table information for completion of table and column names
You can turn off this feature to get a quicker startup with -A

Database changed
mysql> select * from user;
+----+----------+----------+----------+---------------------+---------------------+
| id | passport | password | nickname | create_at           | update_at           |
+----+----------+----------+----------+---------------------+---------------------+
|  1 | JKs      | 123abc   | jasper   | 2023-11-17 08:51:58 | 2023-11-17 08:51:58 |
|  2 | Clark    | 123abc   | Clark    | 2023-11-17 09:08:52 | 2023-11-17 09:08:52 |
+----+----------+----------+----------+---------------------+---------------------+
2 rows in set (0.09 sec)
```
