package services

import (
	"rabbitmq/product/datamodels"
	"rabbitmq/product/repositories"
)

type IOrderService interface {
	Get(int64) (*datamodels.Order, error)
	GetAll() ([]*datamodels.Order, error)
	GetAllInfo() (map[int]map[string]string, error)
	Delete(int64) bool
	Insert(order datamodels.Order) (int64, error)
	Update(order *datamodels.Order) error
}

type OrderService struct {
	OrderRepository repositories.IOrderRepository
}

func NewOrderService(OrderRepository repositories.IOrderRepository) IOrderService {
	return &OrderService{OrderRepository}
}
func (o *OrderService) Get(orderID int64) (order *datamodels.Order, err error) {
	return
}

func (o *OrderService) GetAll()(orders []*datamodels.Order,err error) {
	return
}

func (o *OrderService) GetAllInfo()(orderMaps map[int]map[string]string,err error)  {
	return
}

func (o *OrderService) Delete(orderID int64) bool  {
	return true
}

func (o *OrderService) Insert(order datamodels.Order)(orderID int64,err error)  {
	return
}

func (o *OrderService) Update(order *datamodels.Order)(err error)  {
	return
}
