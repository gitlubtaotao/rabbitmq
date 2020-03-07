package main

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"rabbitmq/product/backend/web/controllers"
	"rabbitmq/product/dbcon"
	"rabbitmq/product/repositories"
	"rabbitmq/product/services"
)

func main() {
	//1.创建iris 实例
	app := iris.New()
	//2.设置错误模式，在mvc模式下提示错误
	app.Logger().SetLevel("debug")
	//3.注册模板
	template := iris.HTML("./web/views", ".html").Layout("shared/layout.html").Reload(true)
	app.RegisterView(template)
	//	注册静态资源
	app.Favicon("./web/assets/favicon.ico")
	app.HandleDir("/assets", "./web/assets", iris.DirOptions{
		Gzip: false,
		// List the files inside the current requested directory if `IndexName` not found.
		ShowList: true,
	})
	
	//错误页面处理
	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("message", ctx.Values().GetStringDefault("message", "访问的页面出错！"))
		ctx.ViewLayout("")
		_ = ctx.View("shared/error.html")
	})
	//连接sql服务器
	db, err := dbcon.NewMysqlConn()
	if err != nil {
		panic(err)
	}
	ctx,cancel := context.WithCancel(context.Background())
	defer cancel()
	
	//5.注册控制器
	productRepository := repositories.NewProductManager("product", db)
	productSerivce := services.NewProductService(productRepository)
	productParty := app.Party("/product")
	product := mvc.New(productParty)
	product.Register(ctx, productSerivce)
	product.Handle(new(controllers.ProductController))
	
	//运行iris
	config := iris.WithConfiguration(iris.YAML("./config/iris.yml"))
	_ = app.Run(
		iris.Addr("localhost:8080"),
		config,
	)
}
