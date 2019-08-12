package vo

type ResultVo struct {
	Code int         `json:code`
	Data interface{} `json:data`
}
