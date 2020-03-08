package repositories

import (
	"database/sql"
	"log"
	"rabbitmq/product/datamodels"
	"rabbitmq/product/util"
)

type IOrderRepository interface {
	Conn() error
	Insert(*datamodels.Order) (int64, error)
	Delete(int64) bool
	Update(*datamodels.Order) error
	SelectByKey(int64) (*datamodels.Order, error)
	SelectAll() ([]*datamodels.Order, error)
	SelectAllWithInfo() (map[int]map[string]string, error)
}

type OrderMangerRepository struct {
	table     string
	mysqlConn *sql.DB
}

func NewOrderMangerRepository(table string, sqlConn *sql.DB) IOrderRepository {
	return &OrderMangerRepository{table: table, mysqlConn: sqlConn}
}

func (o *OrderMangerRepository) Conn() (err error) {
	if o.mysqlConn == nil {
		if sqlConn, err := util.NewMysqlConn(); err != nil {
			log.Fatal(err)
			return err
		} else {
			o.mysqlConn = sqlConn
		}
	}
	if o.table == "" {
		o.table = "order"
	}
	return nil
}

func (o *OrderMangerRepository) Insert(order *datamodels.Order) (OrderId int64, err error) {
	return
}

func (o *OrderMangerRepository) Delete(orderID int64) bool {
	return true
}

func (o *OrderMangerRepository) Update(order *datamodels.Order) (err error) {
	return
}

func (o *OrderMangerRepository) SelectByKey(orderID int64) (order *datamodels.Order, err error) {
	return
}

func (o *OrderMangerRepository) SelectAll() (orders []*datamodels.Order, err error) {
	return
}
func (o *OrderMangerRepository) SelectAllWithInfo() (orders map[int]map[string]string, err error) {
	return
}
