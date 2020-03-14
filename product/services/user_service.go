package services

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"rabbitmq/product/datamodels"
	"rabbitmq/product/repositories"
)

type IUserService interface {
	IsPwdSuccess(userName string, pwd string) (user *datamodels.User, isOk bool)
	AddUser(user *datamodels.User) (userId int64, err error)
	UpdateUser(user *datamodels.User) (err error)
	GetUserById(id int64) (*datamodels.User, error)
}

func NewUserService(repository repositories.IUserRepository) IUserService {
	return &UserService{repository}
}

func NewUserServiceNew() IUserService {
	repository := &repositories.UserManagerRepository{}
	return &UserService{UserRepository: repository}
}

type UserService struct {
	UserRepository repositories.IUserRepository
}

func (u *UserService) IsPwdSuccess(userName string, pwd string) (user *datamodels.User, isOk bool) {
	
	user, err := u.UserRepository.Select(userName)
	
	if err != nil {
		return
	}
	isOk, _ = ValidatePassword(pwd, user.HashPassword)
	
	if !isOk {
		return &datamodels.User{}, false
	}
	
	return
}

func (u *UserService) AddUser(user *datamodels.User) (userId int64, err error) {
	pwdByte, errPwd := GeneratePassword(user.HashPassword)
	if errPwd != nil {
		return userId, errPwd
	}
	user.HashPassword = string(pwdByte)
	return u.UserRepository.Insert(user)
}

func (u *UserService) UpdateUser(user *datamodels.User) (err error) {
	return u.UserRepository.UpdateUser(user)
}

func (u *UserService) GetUserById(id int64) (*datamodels.User, error) {
	user, err := u.UserRepository.SelectByID(id)
	return user, err
}

func GeneratePassword(userPassword string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
}

func ValidatePassword(userPassword string, hashed string) (isOK bool, err error) {
	if err = bcrypt.CompareHashAndPassword([]byte(hashed), []byte(userPassword)); err != nil {
		return false, errors.New("密码比对错误！")
	}
	return true, nil
}
