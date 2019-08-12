package utils

import (
	"bytes"
	"crypto/md5"
	"database/sql"
	"fmt"
	"math/big"
	"math/rand"
	"reflect"
	"time"

	. "maoguo/henan/misc/mysql"
	"maoguo/henan/misc/utils/parse"

	uuid "github.com/satori/go.uuid"
	"github.com/wonderivan/logger"
)

/*
获取UUID
*/
func GetUUID() string {
	UUID, err := uuid.NewV4()
	if err != nil {
		logger.Error("GET UUID failed,err :", err)
		return ""
	}
	return UUID.String()
}

/*
token生成
*/
func Encrypt(str string) string {
	bi := new(big.Int)
	bi, _ = bi.SetString("01213910847463829232312312", 10)
	a := new(big.Int)
	a = a.SetBytes([]byte(str))
	a.Xor(a, bi)
	return a.Text(16)
}

/*
token解析
*/
func Decrypt(str string) string {
	bi := new(big.Int)
	bi, _ = bi.SetString("01213910847463829232312312", 10)
	a := new(big.Int)
	a, _ = a.SetString(str, 16)
	a.Xor(a, bi)
	return string(a.Bytes())
}

func GetValue(maps map[string]string, key string) string {
	if maps == nil {
		return ""
	}
	if a, ok := maps[key]; ok {
		return a
	}
	return ""
}

/*
	字符串MD5
*/
func MD5(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)
	return md5str
}

//查询，返回map集合()
func QueryList(sqlStr string) []map[string]string {
	rows, err := DB.Query(sqlStr)
	if err != nil {
		logger.Error("query List failed,err:", err)
		return nil
	}
	columns, err := rows.Columns()
	if err != nil {
		logger.Error("query List failed,err:", err)
		return nil
	}
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	var result []map[string]string
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			logger.Error("query List failed,err:", err)
			return nil
		}
		var value string
		obj := map[string]string{}
		for i, col := range values {
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			obj[columns[i]] = value
			result = append(result, obj)
		}
	}
	if err = rows.Err(); err != nil {
		logger.Error("query List failed,err:", err)
		return nil
	}
	return result
}

func Query(sqlStr string) (map[string]string, error) {
	list := QueryList(sqlStr)
	if len(list) > 0 {
		return list[0], nil
	}
	return make(map[string]string), fmt.Errorf("数据不存在")
}

/*
	生成随机字符串
*/
func GetRand(length int) string {
	rand.Seed(time.Now().Unix())
	var bt bytes.Buffer
	for i := 0; i < length; i++ {
		index := rand.Intn(3)
		switch index {
		case 0:
			data := rand.Intn(10)
			bt.WriteString(parse.IntToString(data))
		case 1:
			data := rand.Intn(26) + 65
			bt.WriteString(string(data))
		case 2:
			data := rand.Intn(26) + 97
			bt.WriteString(string(data))
		}
	}
	return bt.String()
}

func GetSmsValcationCode() string {
	rand.Seed(time.Now().Unix())
	var bt bytes.Buffer
	for i := 0; i < 4; i++ {
		data := rand.Intn(10)
		bt.WriteString(parse.IntToString(data))
	}
	return bt.String()
}

func GetRandInt(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max)%(max-min+1) + min
}

//判断元素是否存在与slice中
func IsExistItem(array interface{}, value interface{}) bool {
	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)
		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(value, s.Index(i).Interface()) {
				return true
			}
		}
	}
	return false
}
