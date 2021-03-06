package repositories

import (
	"database/sql"
	"errors"
	"rabbitmq/product/datamodels"
	"rabbitmq/product/util"
	
	"strconv"
)

type IUserRepository interface {
	Conn() error
	Select(userName string) (user *datamodels.User, err error)
	Insert(user *datamodels.User) (userId int64, err error)
	UpdateUser(user *datamodels.User) (err error)
	SelectByID(userId int64) (*datamodels.User, error)
}

func NewUserRepository(table string, db *sql.DB) IUserRepository {
	return &UserManagerRepository{table, db}
}

type UserManagerRepository struct {
	table     string
	mysqlConn *sql.DB
}

func (u *UserManagerRepository) Conn() (err error) {
	if u.mysqlConn == nil {
		mysql, errMysql := util.NewMysqlConn()
		if errMysql != nil {
			return errMysql
		}
		u.mysqlConn = mysql
	}
	if u.table == "" {
		u.table = "user"
	}
	return
}

func (u *UserManagerRepository) Select(userName string) (user *datamodels.User, err error) {
	if userName == "" {
		return &datamodels.User{}, errors.New("条件不能为空！")
	}
	if err = u.Conn(); err != nil {
		return &datamodels.User{}, err
	}
	
	sql := "Select * from " + u.table + " where userName=?"
	rows, errRows := u.mysqlConn.Query(sql, userName)
	defer rows.Close()
	if errRows != nil {
		return &datamodels.User{}, errRows
	}
	
	result := util.GetResultRow(rows)
	if len(result) == 0 {
		return &datamodels.User{}, errors.New("用户不存在！")
	}
	
	user = &datamodels.User{}
	util.DataToStructByTagSql(result, user)
	return
}

func (u *UserManagerRepository) Insert(user *datamodels.User) (userId int64, err error) {
	if err = u.Conn(); err != nil {
		return
	}
	
	sql := "INSERT " + u.table + " SET nickName=?,userName=?,passWord=?"
	stmt, errStmt := u.mysqlConn.Prepare(sql)
	defer stmt.Close()
	if errStmt != nil {
		return userId, errStmt
	}
	result, errResult := stmt.Exec(user.NickName, user.UserName, user.HashPassword)
	if errResult != nil {
		return userId, errResult
	}
	return result.LastInsertId()
}

func (u *UserManagerRepository) SelectByID(userId int64) (user *datamodels.User, err error) {
	if err = u.Conn(); err != nil {
		return &datamodels.User{}, err
	}
	sql := "select * from " + u.table + " where ID=" + strconv.FormatInt(userId, 10)
	row, errRow := u.mysqlConn.Query(sql)
	defer row.Close()
	if errRow != nil {
		return &datamodels.User{}, errRow
	}
	result := util.GetResultRow(row)
	if len(result) == 0 {
		return &datamodels.User{}, errors.New("用户不存在！")
	}
	user = &datamodels.User{}
	util.DataToStructByTagSql(result, user)
	return
}

func (u *UserManagerRepository) UpdateUser(user *datamodels.User) (err error) {
	if err = u.Conn(); err != nil {
		return err
	}
	stmt, err := u.mysqlConn.Prepare("update " + u.table + " set nickName=?,userName=?,ipAddress=? Where ID=" + strconv.FormatInt(user.ID, 10))
	if err != nil {
		return err
	}
	_, err = stmt.Exec(user.NickName, user.UserName, user.IpAddress)
	return err
}
