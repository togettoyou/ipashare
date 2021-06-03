package tools

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// StructToMap 利用反射将结构体转化为map
func StructToMap(obj interface{}) map[string]interface{} {
	obj1 := reflect.TypeOf(obj)
	obj2 := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < obj1.NumField(); i++ {
		data[obj1.Field(i).Name] = obj2.Field(i).Interface()
	}
	return data
}

// ArrayToString 将数组格式化为字符串
func ArrayToString(array []interface{}) string {
	return strings.Replace(strings.Trim(fmt.Sprint(array), "[]"), " ", ",", -1)
}

// StrToFloat64 支持指定精度
func StrToFloat64(str string, len int) (float64, error) {
	lenstr := "%." + strconv.Itoa(len) + "f"
	value, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0, err
	}
	nstr := fmt.Sprintf(lenstr, value) //指定精度
	val, err := strconv.ParseFloat(nstr, 64)
	if err != nil {
		return 0, err
	}
	return val, nil
}
