package repositories

import (
	"database/sql"
	"rabbitmq/product/datamodels"
	
	"rabbitmq/product/util"
	"strconv"
)

//第一步，先开发对应的接口
//第二步，实现定义的接口

type IProduct interface {
	Conn() (err error)
	Insert(datamodels.Product) (int64, error)
	Delete(int64) bool
	Update(*datamodels.Product) error
	SelectByKey(int64) (*datamodels.Product, error)
	SelectAll() ([]*datamodels.Product, error)
}

type ProductManager struct {
	table     string
	mysqlConn *sql.DB
}

func NewProductManager(table string, db *sql.DB) IProduct {
	return &ProductManager{table: table, mysqlConn: db}
}

func (p *ProductManager) Conn() (err error) {
	if p.mysqlConn == nil {
		mysql, err := util.NewMysqlConn()
		if err != nil {
			return err
		}
		p.mysqlConn = mysql
	}
	if p.table == "" {
		p.table = "product"
	}
	return
}

func (p *ProductManager) Insert(product datamodels.Product) (id int64, err error) {
	if err := p.Conn(); err != nil {
		return 0, err
	}
	//2.准备sql
	sql := "INSERT " + p.table + " SET productName=?,productNum=?,productImage=?,productUrl=?"
	stmt, errSql := p.mysqlConn.Prepare(sql)
	defer stmt.Close()
	if errSql != nil {
		return 0, errSql
	}
	//3.传入参数,可以使用gorm
	result, errStmt := stmt.Exec(product.ProductName, product.ProductNum, product.ProductImage, product.ProductUrl)
	if errStmt != nil {
		return 0, errStmt
	}
	return result.LastInsertId()
}

func (p *ProductManager) Delete(productId int64) bool {
	//1.判断连接是否存在
	if err := p.Conn(); err != nil {
		return false
	}
	sql := "delete from" + p.table + "where ID=?"
	stmt, err := p.mysqlConn.Prepare(sql)
	defer stmt.Close()
	if err != nil {
		return false
	}
	_, err = stmt.Exec(strconv.FormatInt(productId, 10))
	if err != nil {
		return false
	}
	return true
}
func (p *ProductManager) Update(product *datamodels.Product) (err error) {
	if err := p.Conn(); err != nil {
		return err
	}
	sql := "Update " + p.table + "set productName=?,productNum=?,productImage=?,productUrl=? where ID=" + strconv.FormatInt(product.ID, 10)
	
	stmt, err := p.mysqlConn.Prepare(sql)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(product.ProductName, product.ProductNum, product.ProductImage, product.ProductUrl)
	if err != nil {
		return err
	}
	return nil
}
func (p *ProductManager) SelectByKey(productID int64) (product *datamodels.Product, err error) {
	if err = p.Conn(); err != nil {
		return &datamodels.Product{}, err
	}
	sql := "Select * from " + p.table + " where ID =" + strconv.FormatInt(productID, 10)
	row, err := p.mysqlConn.Query(sql)
	if err != nil {
		return nil, err
	}
	defer row.Close()
	result := util.GetResultRow(row)
	if len(result) == 0 {
		return &datamodels.Product{}, nil
	}
	util.DataToStructByTagSql(result, product)
	return
}

//获取所有商品
func (p *ProductManager) SelectAll() (productArray []*datamodels.Product, errProduct error) {
	//1.判断连接是否存在
	if err := p.Conn(); err != nil {
		return nil, err
	}
	sql := "Select * from " + p.table
	rows, err := p.mysqlConn.Query(sql)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	
	result := util.GetResultRows(rows)
	if len(result) == 0 {
		return nil, nil
	}
	
	for _, v := range result {
		product := &datamodels.Product{}
		util.DataToStructByTagSql(v, product)
		productArray = append(productArray, product)
	}
	return
}
