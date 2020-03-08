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
}

type UserManagerRepository struct {
	table     string
	mysqlConn *sql.DB
}

func NewUserRepository(table string, db *sql.DB) IUserRepository {
	return &UserManagerRepository{table, db}
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
	user = &datamodels.User{}
	if userName == "" {
		return user, errors.New("条件不能为空！")
	}
	if err = u.Conn(); err != nil {
		return user, err
	}
	rows, errRows := u.mysqlConn.Query("Select * from "+u.table+" where userName=?", userName)
	if errRows != nil {
		return user, errRows
	}
	defer rows.Close()
	result := util.GetResultRow(rows)
	if len(result) == 0 {
		return user, errors.New("用户不存在！")
	}
	util.DataToStructByTagSql(result, user)
	return
}

func (u *UserManagerRepository) Insert(user *datamodels.User) (userId int64, err error) {
	if err = u.Conn(); err != nil {
		return
	}
	stmt, errStmt := u.mysqlConn.Prepare("INSERT " + u.table + " SET nickName=?,userName=?,passWord=?")
	if errStmt != nil {
		return userId, errStmt
	}
	defer stmt.Close()
	result, errResult := stmt.Exec(user.NickName, user.UserName, user.HashPassword)
	if errResult != nil {
		return userId, errResult
	}
	return result.LastInsertId()
}

func (u *UserManagerRepository) SelectByID(userId int64) (user *datamodels.User, err error) {
	user = &datamodels.User{}
	if err = u.Conn(); err != nil {
		return user, err
	}
	row, errRow := u.mysqlConn.Query("select * from " + u.table + " where ID=" + strconv.FormatInt(userId, 10))
	if errRow != nil {
		return user, errRow
	}
	defer row.Close()
	result := util.GetResultRow(row)
	if len(result) == 0 {
		return user, errors.New("用户不存在！")
	}
	util.DataToStructByTagSql(result, user)
	return
}
