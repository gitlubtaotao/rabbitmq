package controllers

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"net/http"
	"path/filepath"
	"rabbitmq/product/datamodels"
	"rabbitmq/product/services"
	"rabbitmq/product/util"
	"strconv"
)

type ProductController struct {
	Ctx            iris.Context
	ProductService services.IProductService
	OrderService   services.IOrderService
	Session        sessions.Sessions
}

func (p *ProductController) GetDetail() mvc.View {
	product, err := p.ProductService.GetProductByID(1)
	if err != nil {
		p.Ctx.Application().Logger().Error(err)
	}
	return mvc.View{
		Layout: "shared/productLayout.html",
		Name:   "product/view.html",
		Data: iris.Map{
			"product": product,
		},
	}
}

func (p *ProductController) GetOrder() mvc.View {
	userString := p.Ctx.GetCookie("uid")
	productID, err := strconv.Atoi(p.Ctx.URLParam("productID"))
	if err != nil {
		p.Ctx.Application().Logger().Error(err)
		p.getOrderError(0)
	}
	product, err := p.ProductService.GetProductByID(int64(productID))
	if err != nil {
		p.Ctx.Application().Logger().Error(err)
		return p.getOrderError(0)
	}
	if product.ProductNum > 0 {
		return p.getOrderSuccess(product, userString)
	} else {
		return p.getOrderError(0)
	}
}

func (p *ProductController) getOrderSuccess(product *datamodels.Product, userString string) mvc.View {
	product.ProductNum -= 1
	err := p.ProductService.UpdateProduct(product)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
		return p.getOrderError(0)
	}
	//创建订单
	userID, err := strconv.Atoi(userString)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
		return p.getOrderError(0)
	}
	
	order := &datamodels.Order{
		UserId:      int64(userID),
		ProductId:   int64(product.ID),
		OrderStatus: datamodels.OrderSuccess,
	}
	//新建订单
	orderID, err := p.OrderService.InsertOrder(order)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
		return p.getOrderError(0)
	}
	return mvc.View{
		Layout: "shared/productLayout.html",
		Name:   "product/result.html",
		Data: iris.Map{
			"orderID":     orderID,
			"showMessage": "抢购成功！",
		},
	}
}
func (p *ProductController) getOrderError(orderID int) mvc.View {
	return mvc.View{
		Layout: "shared/productLayout.html",
		Name:   "product/result.html",
		Data: iris.Map{
			"orderID":     orderID,
			"showMessage": "抢购失败！",
		},
	}
}

var (
	//生成的Html保存目录
	htmlOutPath = "./web/htmlProductShow/"
	//静态文件模版目录
	templatePath = "./web/views/template/"
)

func (p *ProductController) GetGenerateHtml() (response *mvc.Response) {
	productString := p.Ctx.URLParam("productID")
	productID, err := strconv.Atoi(productString)
	if err != nil {
		p.Ctx.Application().Logger().Error(err)
		return responseStatus(false)
	}
	contentsTmp, fileName, err := util.GenerateFilerName(filepath.Join(templatePath, "product.html"), filepath.Join(htmlOutPath, "htmlProduct.html"))
	if err != nil {
		p.Ctx.Application().Logger().Error(err)
		return responseStatus(false)
	}
	defer fileName.Close()
	product, err := p.ProductService.GetProductByID(int64(productID))
	if err != nil {
		p.Ctx.Application().Logger().Error(err)
		return responseStatus(false)
	}
	err = contentsTmp.Execute(fileName, &product)
	if err != nil {
		p.Ctx.Application().Logger().Error(err)
		return responseStatus(false)
	}
	return responseStatus(true)
}

func responseStatus(status bool) *mvc.Response {
	var object map[string]string
	object = make(map[string]string)
	if status {
		object["status"] = strconv.FormatInt(http.StatusOK, 10)
		object["content"] = "生成静态文件成功"
	} else {
		object["status"] = strconv.FormatInt(http.StatusInternalServerError, 10)
		object["content"] = "生成静态文件失败"
	}
	return &mvc.Response{
		Object:      object,
		ContentType: "application/json",
	}
}
