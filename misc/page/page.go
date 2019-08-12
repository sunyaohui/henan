package page

type Page struct {
	PageNo     int         `json:"pageNo"`
	PageSize   int         `json:"pageSize"`
	TotalCount int         `json:"totalCount"`
	OrderBy    string      `json:"orderBy"`
	OrderType  string      `json:"orderType"`
	List       interface{} `json:"list"`
	Skip       int         `json:"skip"`
	TotalPage  int         `json:"totalPage"`
}
