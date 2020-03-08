package repositories

import (
	"database/sql"
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

func NewOrderMangerRepository(table string, sql *sql.DB) IOrderRepository {
	return &OrderMangerRepository{table: table, mysqlConn: sql}
}

type OrderMangerRepository struct {
	table     string
	mysqlConn *sql.DB
}

func (o *OrderMangerRepository) Conn() error {
	if o.mysqlConn == nil {
		mysql, err := util.NewMysqlConn()
		if err != nil {
			return err
		}
		o.mysqlConn = mysql
	}
	if o.table == "" {
		o.table = "order"
	}
	return nil
}

func (o *OrderMangerRepository) Insert(order *datamodels.Order) (productID int64, err error) {
	if err = o.Conn(); err != nil {
		return
	}
	sql := "INSERT " + o.table + " set userID=?,productID=?,orderStatus=?"
	stmt, errStmt := o.mysqlConn.Prepare(sql)
	if errStmt != nil {
		return productID, errStmt
	}
	defer stmt.Close()
	result, errResult := stmt.Exec(order.UserId, order.ProductId, order.OrderStatus)
	if errResult != nil {
		return productID, errResult
	}
	return result.LastInsertId()
}

func (o *OrderMangerRepository) Delete(orderID int64) (isOk bool) {
	if err := o.Conn(); err != nil {
		return
	}
	sql := "delete from " + o.table + " where ID =?"
	stmt, errStmt := o.mysqlConn.Prepare(sql)
	if errStmt != nil {
		return
	}
	defer stmt.Close()
	_, err := stmt.Exec(orderID)
	if err != nil {
		return
	}
	return true
}

func (o *OrderMangerRepository) Update(order *datamodels.Order) (err error) {
	if errConn := o.Conn(); errConn != nil {
		return errConn
	}
	
	sql := "Update " + o.table + " set userID=?,productID=?,orderStatus=? Where ID=" + strconv.FormatInt(order.ID, 10)
	stmt, errStmt := o.mysqlConn.Prepare(sql)
	if errStmt != nil {
		return errStmt
	}
	defer stmt.Close()
	_, err = stmt.Exec(order.UserId, order.ProductId, order.OrderStatus)
	return
}

func (o *OrderMangerRepository) SelectByKey(orderID int64) (order *datamodels.Order, err error) {
	if errConn := o.Conn(); errConn != nil {
		return &datamodels.Order{}, errConn
	}
	
	sql := "Select * From " + o.table + " where ID=" + strconv.FormatInt(orderID, 10)
	row, errRow := o.mysqlConn.Query(sql)
	
	if errRow != nil {
		return &datamodels.Order{}, errRow
	}
	defer row.Close()
	result := util.GetResultRow(row)
	if len(result) == 0 {
		return &datamodels.Order{}, err
	}
	
	order = &datamodels.Order{}
	util.DataToStructByTagSql(result, order)
	return
}

func (o *OrderMangerRepository) SelectAll() (orderArray []*datamodels.Order, err error) {
	if errConn := o.Conn(); errConn != nil {
		return nil, errConn
	}
	sql := "Select * from " + o.table
	rows, errRows := o.mysqlConn.Query(sql)
	
	if errRows != nil {
		return nil, errRows
	}
	defer rows.Close()
	result := util.GetResultRows(rows)
	if len(result) == 0 {
		return nil, err
	}
	
	for _, v := range result {
		order := &datamodels.Order{}
		util.DataToStructByTagSql(v, order)
		orderArray = append(orderArray, order)
	}
	return
}

func (o *OrderMangerRepository) SelectAllWithInfo() (OrderMap map[int]map[string]string, err error) {
	if errConn := o.Conn(); errConn != nil {
		return nil, errConn
	}
	sql := "Select o.ID,p.productName,o.orderStatus From imooc.order as o left join product as p on o.productID=p.ID"
	rows, errRows := o.mysqlConn.Query(sql)
	
	if errRows != nil {
		return nil, errRows
	}
	defer rows.Close()
	return util.GetResultRows(rows), err
}
