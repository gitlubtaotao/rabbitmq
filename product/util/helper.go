package util

import (
	"errors"
	"html/template"
	"net"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func DataToStructByTagSql(data map[string]string, obj interface{}) {
	objValue := reflect.ValueOf(obj).Elem()
	
	for i := 0; i < objValue.NumField(); i++ {
		//获取sql对应的值
		value := data[objValue.Type().Field(i).Tag.Get("sql")]
		//获取对应字段的名称
		name := objValue.Type().Field(i).Name
		//获取对应字段类型
		structFieldType := objValue.Field(i).Type()
		//获取变量类型，也可以直接写"string类型"
		val := reflect.ValueOf(value)
		var err error
		if structFieldType != val.Type() {
			//类型转换
			val, err = TypeConversion(value, structFieldType.Name()) //类型转换
			if err != nil {
			
			}
		}
		//设置类型值
		objValue.FieldByName(name).Set(val)
	}
}

//类型转换
func TypeConversion(value string, ntype string) (reflect.Value, error) {
	if ntype == "string" {
		return reflect.ValueOf(value), nil
	} else if ntype == "time.Time" {
		t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
		return reflect.ValueOf(t), err
	} else if ntype == "Time" {
		t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
		return reflect.ValueOf(t), err
	} else if ntype == "int" {
		i, err := strconv.Atoi(value)
		return reflect.ValueOf(i), err
	} else if ntype == "int8" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(int8(i)), err
	} else if ntype == "int32" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(int64(i)), err
	} else if ntype == "int64" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(i), err
	} else if ntype == "float32" {
		i, err := strconv.ParseFloat(value, 64)
		return reflect.ValueOf(float32(i)), err
	} else if ntype == "float64" {
		i, err := strconv.ParseFloat(value, 64)
		return reflect.ValueOf(i), err
	}
	
	//else if .......增加其他一些类型的转换
	
	return reflect.ValueOf(value), errors.New("未知的类型：" + ntype)
}

func ClientPublicIP(r *http.Request) string {
	var ip string
	for _, ip = range strings.Split(r.Header.Get("X-Forwarded-For"), ",") {
		ip = strings.TrimSpace(ip)
		if ip != "" {
			return ip
		}
	}
	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if ip != "" {
		return ip
	}
	
	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}
	return ""
}

func GenerateFilerName(templatePath, htmlOutPath string) (*template.Template, *os.File, error) {
	contentsTap, err := template.ParseFiles(templatePath)
	if err != nil {
		return &template.Template{}, &os.File{}, err
	}
	filer, err := GenerateStaticHtml(htmlOutPath)
	return contentsTap, filer, err
}
func GenerateStaticHtml(fileName string) (filer *os.File, err error) {
	if staticHtmlExist(fileName) {
		err := os.Remove(fileName)
		if err != nil {
			return &os.File{}, err
		}
	}
	//2.生成静态文件
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return &os.File{}, err
	}
	return file, err
}

func staticHtmlExist(fileName string) bool {
	_, err := os.Stat(fileName)
	return err == nil || os.IsExist(err)
}
