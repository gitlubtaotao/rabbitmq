package services

import (
	"golang.org/x/crypto/bcrypt"
	"rabbitmq/product/datamodels"
	"rabbitmq/product/repositories"
)

type IUserService interface {
	IsPwdSuccess(userName string, pwd string) (user *datamodels.User, isOk bool)
	AddUser(user *datamodels.User) (userId int64, err error)
}

type UserService struct {
	UserRepository repositories.IUserRepository
}

func NewUserService(repository repositories.IUserRepository) IUserService {
	return &UserService{UserRepository: repository}
}

func (u *UserService) IsPwdSuccess(userName string, pwd string) (user *datamodels.User, isOk bool) {
	user = &datamodels.User{}
	var err error
	user, err = u.UserRepository.Select(userName)
	if err != nil {
		return user, false
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

func ValidatePassword(userPassword string, hashed string) (isOk bool, err error) {
	err = bcrypt.CompareHashAndPassword([]byte(hashed), []byte(hashed))
	if err != nil {
		return
	}
	return true, nil
}

func GeneratePassword(userPassword string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
}
