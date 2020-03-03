package repositories

import (
	"database/sql"
	"rabbitmq/product/datamodels"
	"rabbitmq/product/dbcon"
)

//第一步，先开发对应的接口
//第二步，实现定义的接口

type IProduct interface {
	Conn() (err error)
	Insert(*datamodels.Product) (int64, error)
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
		mysql, err := dbcon.NewMysqlConn()
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
	sql := "INSERT product SET productName=?,productNum=?,productImage=?,productUrl=?"
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
