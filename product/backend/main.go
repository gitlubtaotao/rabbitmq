package main

import (
	"context"
	"github.com/kataras/iris/v12"
	
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/kataras/iris/v12/mvc"
	"log"
	"rabbitmq/product/backend/web/controllers"
	"rabbitmq/product/util"
	
	"rabbitmq/product/repositories"
	"rabbitmq/product/services"
)

func main() {
	//1.创建iris 实例
	app := iris.New()
	//2.设置错误模式，在mvc模式下提示错误
	app.Logger().SetLevel("debug")
	app.Use(recover.New())
	app.Use(logger.New())
	//3.注册模板
	template := iris.HTML("./web/views", ".html").Layout("shared/layout.html").Reload(true)
	app.RegisterView(template)
	//	注册静态资源
	//app.Favicon("./web/assets/favicon.ico")
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
	db, err := util.NewMysqlConn()
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	//5.注册控制器
	productRepository := repositories.NewProductManager("product", db)
	productService := services.NewProductService(productRepository)
	mvc.New(app.Party("/product")).Register(ctx, productService).Handle(new(controllers.ProductController))
	
	orderRepository := repositories.NewOrderMangerRepository("order", db)
	orderService := services.NewOrderService(orderRepository)
	mvc.New(app.Party("/order")).Register(ctx, orderService).Handle(new(controllers.OrderController))
	
	//运行iris
	config := iris.WithConfiguration(iris.YAML("./config/iris.yml"))
	err = app.Listen(":8080", config)
	if err != nil {
		log.Fatal(err)
	}
}
