package controllers

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"rabbitmq/product/services"
)

type OrderController struct {
	Ctx          iris.Context
	OrderService services.OrderService
}

func (o *OrderController) Get() mvc.View {
	orderArray, err := o.OrderService.GetAllInfo()
	if err != nil {
		o.Ctx.Application().Logger().Debug("查询订单信息失败")
	}
	return mvc.View{
		Name: "order/view.html",
		Data: iris.Map{
			"order": orderArray,
		},
	}
}

