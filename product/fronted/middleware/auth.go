package middleware

import (
	"github.com/kataras/iris/v12"
	"rabbitmq/product/encrypt"
	"rabbitmq/product/services"
	"strconv"
)

func AuthConProduct(ctx iris.Context) {
	uid := ctx.GetCookie("sign")
	if uid == "" {
		ctx.Application().Logger().Debug("必须先登录!")
		ctx.Redirect("/user/login")
		return
	}
	ID, err := encrypt.DePwdCode(uid)
	if err != nil {
		ctx.Application().Logger().Error(err)
		ctx.Redirect("/user/login")
		return
	}
	userService := services.NewUserServiceNew()
	id, _ := strconv.Atoi(string(ID))
	_, err = userService.GetUserById(int64(id))
	if err != nil {
		ctx.Application().Logger().Error(err)
		ctx.Redirect("/user/login")
		return
	}
	ctx.Application().Logger().Info("已经登陆")
	ctx.Next()
}
