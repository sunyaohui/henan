package result

import (
	//"encoding/json"
	//"fmt"
	"encoding/json"
	"maoguo/henan/misc/utils/parse"
	"maoguo/henan/vo"
)

func ResponseMessage(code int, message string) vo.ResultVo {
	v := vo.ResultVo{
		Code: code,
		Data: map[string]string{"info": message},
	}
	return v
}

func ResponseData(code int, data interface{}) vo.ResultVo {
	v := vo.ResultVo{
		Code: code,
		Data: data,
	}
	return v
}

func ResponseWrite(code int, message string) string {
	v := ResponseMessage(code, message)
	result, err := json.Marshal(&v)
	if err != nil {
		return "系统异常"
	}
	return string(result)
}

func ResponseWriteData(code int, data interface{}) string {
	v := ResponseData(code, data)
	result := parse.ParseJson(&v)
	return string(result)
}
