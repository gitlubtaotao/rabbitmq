package util

import (
	"net/http"
	"strconv"
	"sync"
)

//用来存放控制信息，
type AccessControl struct {
	//用来存放用户想要存放的信息
	sourcesArray map[int]interface{}
	sync.RWMutex
}

func NewAccessControl() *AccessControl {
	return &AccessControl{sourcesArray: make(map[int]interface{})}
}

func (a *AccessControl) GetNewRecord(uid int) interface{} {
	a.RLock()
	defer a.RUnlock()
	return a.sourcesArray[uid]
}

//对map的读写需要进行加锁处理
func (a *AccessControl) SetNewRecord(uid int) {
	a.Lock()
	defer a.Unlock()
	a.sourcesArray[uid] = "hello world"
}

func (a *AccessControl) GetDistributedRight(req *http.Request, hashConsistent *Consistent,
	localHost string, getOtherData func(host string, request *http.Request) bool) bool {
	uid, err := req.Cookie("uid")
	if err != nil {
		return false
	}
	//采用一致性hash算法，根据用户ID，判断获取具体机器
	hostRequest, err := hashConsistent.Get(uid.Value)
	if err != nil {
		return false
	}
	//判断是否为本机
	if hostRequest == localHost {
		//执行本机数据读取和校验
		return a.GetDataFromMap(uid.Value)
	} else {
		//不是本机充当代理访问数据返回结果
		return getOtherData(hostRequest, req)
	}
}

func (a *AccessControl) GetDataFromMap(uid string) (isOk bool) {
	uidInt, err := strconv.Atoi(uid)
	if err != nil {
		return false
	}
	data := a.GetNewRecord(uidInt)
	//执行逻辑判断
	if data != nil {
		return true
	}
	return
}
