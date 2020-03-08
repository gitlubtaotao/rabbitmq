package repositories

import (
	"database/sql"
	"log"
	"rabbitmq/product/datamodels"
	"rabbitmq/product/util"
	"strconv"
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
	if err = o.Conn(); err != nil {
		return
	}
	s := "INSERT " + o.table + "set userID=?,ProductId=?,OrderStatus=?"
	stmt, err := o.mysqlConn.Prepare(s)
	if err != nil {
		return
	}
	defer stmt.Close()
	result, err := stmt.Exec(order.UserId, order.ProductId, order.OrderStatus)
	if err != nil {
		return
	}
	return result.LastInsertId()
}

func (o *OrderMangerRepository) Delete(orderID int64) bool {
	if err := o.Conn(); err != nil {
		log.Fatal(err)
		return false
	}
	if stmt, err := o.mysqlConn.Prepare(" DELETE from" + o.table + "WHERE ID= ?"); err != nil {
		log.Fatal(err)
		return false
	} else {
		defer stmt.Close()
		_, err := stmt.Exec(orderID)
		if err != nil {
			log.Fatal(err)
			return false
		}
	}
	return true
}

func (o *OrderMangerRepository) Update(order *datamodels.Order) (err error) {
	if err = o.Conn(); err != nil {
		return err
	}
	stmt, err := o.mysqlConn.Prepare("UPDATE " + o.table + "SET UserId=?,ProductId=?,OrderStatus=?" + "WHERE ID=" + strconv.FormatInt(order.ID, 10))
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(order.UserId, order.ProductId, order.OrderStatus)
	return
}

func (o *OrderMangerRepository) SelectByKey(orderID int64) (order *datamodels.Order, err error) {
	order = &datamodels.Order{}
	if err = o.Conn(); err != nil {
		return order, err
	}
	s := "SELECT * from " + o.table + "WHERE ID=" + strconv.FormatInt(orderID, 10)
	stmt, err := o.mysqlConn.Prepare(s)
	if err != nil {
		return order, err
	}
	defer stmt.Close()
	row, err := stmt.Query(s)
	if err != nil {
		return order, err
	}
	defer row.Close()
	result := util.GetResultRow(row)
	if len(result) == 0 {
		return order, err
	}
	util.DataToStructByTagSql(result, order)
	return
}

func (o *OrderMangerRepository) SelectAll() (orders []*datamodels.Order, err error) {
	if err = o.Conn(); err != nil {
		return nil, err
	}
	s := "Select * from " + o.table
	rows, err := o.mysqlConn.Query(s)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	result := util.GetResultRows(rows)
	if len(result) == 0 {
		return nil, err
	}
	for _, v := range result {
		order := &datamodels.Order{}
		util.DataToStructByTagSql(v, order)
		orders = append(orders, order)
	}
	return
}

func (o *OrderMangerRepository) SelectAllWithInfo() (orders map[int]map[string]string, err error) {
	return
}
