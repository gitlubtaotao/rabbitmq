package main

import (
	"errors"
	"fmt"
	"net/http"
	"rabbitmq/product/encrypt"
	"rabbitmq/product/util"
)

func Check(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("执行check函数")
}

//统一验证拦截器，每个接口都需要提前验证
func Auth(w http.ResponseWriter, r *http.Request) error {
	fmt.Println("执行验证！")
	//添加基于cookie的权限验证
	err := CheckUserInfo(r)
	if err != nil {
		return err
	}
	return nil
}

func CheckUserInfo(r *http.Request) error {
	uidCookie, err := r.Cookie("uid")
	if err != nil {
		return errors.New("用户UID Cookie 获取失败！")
	}
	signCookie, err := r.Cookie("sign")
	if err != nil {
		return errors.New("用户加密串 Cookie 获取失败！")
	}
	//对信息进行解密
	signByte, err := encrypt.DePwdCode(signCookie.Value)
	if err != nil {
		return errors.New("加密串已被篡改！")
	}
	if checkInfo(uidCookie.Value, string(signByte)) {
		return nil
	}
	return errors.New("身份校验失败！")
}

//自定义逻辑判断
func checkInfo(checkStr string, signStr string) bool {
	if checkStr == signStr {
		return true
	}
	return false
}
func main() {
	//1、过滤器
	filter := util.NewFilter()
	//注册拦截器
	filter.RegisterFilterUri("/check", Auth)
	http.HandleFunc("/check", filter.Handler(Check))
	err := http.ListenAndServe(":8083", nil)
	if err != nil {
		panic(err)
	}
}
