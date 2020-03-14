package main

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/kataras/iris/v12/mvc"
	"log"
	"rabbitmq/product/fronted/middleware"
	"rabbitmq/product/fronted/web/controllers"
	"rabbitmq/product/repositories"
	"rabbitmq/product/services"
	"rabbitmq/product/util"
)

func main() {
	//1.创建iris 实例
	app := iris.New()
	//2.设置错误模式，在mvc模式下提示错误
	app.Logger().SetLevel("debug")
	app.Use(recover.New())
	app.Use(logger.New())
	//3.注册模板
	tmplate := iris.HTML("./web/views", ".html").Layout("shared/layout.html").Reload(true)
	app.RegisterView(tmplate)
	//4.设置模板
	app.HandleDir("/public", "./web/public", iris.DirOptions{
		Gzip: false,
		// List the files inside the current requested directory if `IndexName` not found.
		ShowList: true,
	})
	app.HandleDir("/html", "./web/htmlProductShow", iris.DirOptions{})
	//访问生成好的html静态文件
	
	//出现异常跳转到指定页面
	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("message", ctx.Values().GetStringDefault("message", "访问的页面出错！"))
		ctx.ViewLayout("")
		ctx.View("shared/error.html")
	})
	//连接数据库
	db, err := util.NewMysqlConn()
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	//TODO 需要简化注册register code
	user := repositories.NewUserRepository("user", db)
	userService := services.NewUserService(user)
	
	mvc.New(app.Party("/user")).Register(userService, ctx).Handle(new(controllers.UserController))
	
	product := repositories.NewProductManager("product", db)
	productService := services.NewProductService(product)
	order := repositories.NewOrderMangerRepository("order", db)
	orderService := services.NewOrderService(order)
	productParty := app.Party("/product")
	productParty.Use(middleware.AuthConProduct)
	mvc.New(productParty).Register(productService, orderService, ctx).Handle(new(controllers.ProductController))
	
	
	config := iris.WithConfiguration(iris.YAML("./config/iris.yml"))
	err = app.Listen(":8082", config)
	if err != nil {
		log.Fatal(err)
	}
	
}
