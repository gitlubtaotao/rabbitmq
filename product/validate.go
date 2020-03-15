package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"rabbitmq/product/encrypt"
	"rabbitmq/product/util"
)

//设置集群地址，最好内外IP
var hostArray = []string{"127.0.0.1", "127.0.0.1"}
var localHost = "127.0.0.1"
var port = "8081"

var hashConsistent *util.Consistent

//创建全局变量
var accessControl = util.NewAccessControl()

func GetDataFromOtherMap(host string, request *http.Request) bool {
	uidPre, err := request.Cookie("uid")
	if err != nil {
		return false
	}
	//获取sign
	uidSign, err := request.Cookie("sign")
	if err != nil {
		return false
	}
	
	//模拟接口访问，
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://"+host+":"+port+"/check", nil)
	if err != nil {
		return false
	}
	
	//手动指定，排查多余cookies
	cookieUid := &http.Cookie{Name: "uid", Value: uidPre.Value, Path: "/"}
	cookieSign := &http.Cookie{Name: "sign", Value: uidSign.Value, Path: "/"}
	//添加cookie到模拟的请求中
	req.AddCookie(cookieUid)
	req.AddCookie(cookieSign)
	
	//获取返回结果
	response, err := client.Do(req)
	if err != nil {
		return false
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return false
	}
	
	//判断状态
	if response.StatusCode == 200 {
		if string(body) == "true" {
			return true
		} else {
			return false
		}
	}
	return false
}
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
	
	//负载均衡器设置
	//采用一致性哈希算法
	hashConsistent = util.NewConsistent()
	//采用一致性hash算法，添加节点
	for _, v := range hostArray {
		hashConsistent.Add(v)
	}
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
