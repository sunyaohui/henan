//类型转换
package parse

import (
	"encoding/binary"
	"encoding/json"
	"reflect"
	"strconv"
	"strings"

	"github.com/wonderivan/logger"
)

/*
object转换为json
*/
func ParseJson(params interface{}) string {
	jsonStr, err := json.MarshalIndent(params, "", " ")
	if err != nil {
		logger.Error("parse json failed", err)
		return ""
	}
	return string(jsonStr)
}

/*
json转换为map
*/
func JsonToMap(jsonStr string) map[string]interface{} {
	result := make(map[string]interface{}, 0)

	err := json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		logger.Error("json parse to map failed", err)
		return nil
	}
	return result
}

func ByteToObj(byte []byte, res interface{}) interface{} {
	err := json.Unmarshal(byte, &res)
	if err != nil {
		logger.Error("json parse to map failed", err)
		return nil
	}
	return res
}

//结构体首字母大写，转换成首字母小写的map
func StructToJsonMap(obj interface{}) map[string]interface{} {
	str := ParseJson(obj)
	return JsonToMap(str)
}

/*
结构体转换为map
*/
func StructToMap(obj interface{}) map[string]interface{} {
	if obj == nil {
		return nil
	}
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}

/*
int转string
*/
func IntToString(params int) string {
	return strconv.Itoa(params)
}

func IntToInt64(params int) int64 {
	str := IntToString(params)
	return StringToInt64(str)
}

/*
string转int
*/
func StringToInt(params string) int {
	if params == "" {
		return 0
	}
	result, err := strconv.Atoi(params)
	if err != nil {
		logger.Error("string parse to int failed", err)
		return 0
	}
	return result
}

func StringToIntDefault(params string, value int) int {
	if params == "" {
		return value
	}
	result, err := strconv.Atoi(params)
	if err != nil {
		logger.Error("string parse to int failed", err)
		return 0
	}
	return result
}

/*
string转int64(mysql：bigint,java:long)
*/
func StringToInt64(params string) int64 {
	int64, err := strconv.ParseInt(params, 10, 64)
	if err != nil {
		logger.Error("string parse to int64 failed,params:%d", params, err)
		return 0
	}
	return int64
}

func StringToInt64Default(params string, value int64) int64 {
	if params == "" {
		return value
	}
	int64, err := strconv.ParseInt(params, 10, 64)
	if err != nil {
		logger.Error("string parse to int64 failed,params:%d", params, err)
		return 0
	}
	return int64
}

/*
int64转string
*/
func Int64ToString(params int64) string {
	string := strconv.FormatInt(params, 10)
	return string
}

/*
int64转float64(mysql:double)
*/
func Int64ToFloat64(params int64) float64 {
	str := Int64ToString(params)
	result, err := strconv.ParseFloat(str, 64)
	if err != nil {
		logger.Error("int64 parse to float64 failed", err)
		return 0
	}
	return result
}

/*
int64转int
*/
func Int64ToInt(params int64) int {
	str := Int64ToString(params)
	return StringToInt(str)
}

/*
string转int64数组
*/
func StringToInt64Arr(params string) []int64 {
	if params == "" {
		return nil
	}
	var result []int64
	arr := strings.Split(params, ",")
	for _, value := range arr {
		if value == "" {
			continue
		}
		result = append(result, StringToInt64(value))
	}
	return result
}

func StringToFloat64(params string) float64 {
	if params == "" {
		return 0
	}
	result, err := strconv.ParseFloat(params, 64)
	if err != nil {
		logger.Error("string parse to Float64 failed", err)
		return 0
	}
	return result
}

func IntToFloat64(params int) float64 {
	str := IntToString(params)
	return StringToFloat64(str)
}

/**
float64转string
'b' (-ddddp±ddd，二进制指数)
'e' (-d.dddde±dd，十进制指数)
'E' (-d.ddddE±dd，十进制指数)
'f' (-ddd.dddd，没有指数)
'g' ('e':大指数，'f':其它情况)
'G' ('E':大指数，'f':其它情况)
*/
func Float64ToString(params float64) string {
	str := strconv.FormatFloat(params, 'f', -1, 64)
	return str
}

func Float64ToInt(v float64) int {
	str := Float64ToString(v)
	return StringToInt(str)
}

func Int64ToBytes(i int64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}

func BytesToInt64(buf []byte) int64 {
	return int64(binary.BigEndian.Uint64(buf))
}

func Float64ToInt64(params float64) int64 {
	str := Float64ToString(params)
	return StringToInt64(str)
}

func Int64ToUnit64(params int64) uint64 {
	str := Int64ToString(params)
	r, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		logger.Error("parse uint failed", err)
		return 0
	}
	return r
}
