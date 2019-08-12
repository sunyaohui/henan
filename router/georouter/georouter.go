package georouter

import (
	"fmt"
	"maoguo/henan/misc/page"
	"maoguo/henan/misc/redis"
	"maoguo/henan/misc/utils"
	"maoguo/henan/misc/utils/parse"
	"maoguo/henan/model/model_mysql"
	"maoguo/henan/result"
	"maoguo/henan/services/impl"
	"net/http"
	"strings"
)

//摇一摇
func FindRandom(w http.ResponseWriter, r *http.Request) {
	userId := parse.StringToInt64(r.PostFormValue("userId"))
	longitude := parse.StringToFloat64(r.PostFormValue("longitude"))
	latitude := parse.StringToFloat64(r.PostFormValue("latitude"))

	yaoKey := []byte("user_geo_yy")
	redis.SetkeyExPrire(yaoKey, 10)
	list := redis.Georadius(yaoKey, longitude, latitude, 9000, "km")
	redis.Geoadd(yaoKey, longitude, latitude, []byte(r.PostFormValue("userId")))

	if list == nil || len(list) == 0 {
		fmt.Fprint(w, result.ResponseWrite(-1, "没有找到同一时刻的人"))
		return
	}
	index := 0
	if len(list) > 1 {
		index = utils.GetRandInt(0, len(list)-1)
	}
	el := list[index]

	if len(list) == 1 && parse.BytesToInt64(el.Member) == userId {
		fmt.Fprint(w, result.ResponseWrite(-1, "没有找到同一时刻的人"))
		return
	}

	if len(list) > 1 && parse.BytesToInt64(el.Member) == userId {
		if index == 0 {
			index = len(list) - 1
		} else {
			index = 0
		}
		el = list[index]
	}

	distance := redis.Geodist(yaoKey, parse.Int64ToBytes(userId), el.Member, "m")
	member := model_mysql.Query("select id,name,sex,headUrl,sign from im_user where id=?", parse.BytesToInt64(el.Member))
	res := make(map[string]interface{}, 0)
	res["distance"] = distance
	res["user"] = member
	new(impl.UserServiceImpl).SaveRockRecord(userId, parse.BytesToInt64(el.Member), distance)
	fmt.Fprintf(w, result.ResponseWriteData(1, res))
}

func UpdateGeo(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	userId := query.Get("userId")
	longitude := parse.StringToFloat64(query.Get("longitude"))
	latitude := parse.StringToFloat64(query.Get("latitude"))
	redis.Geoadd([]byte("user_geo"), longitude, latitude, []byte(userId))
	model_mysql.Exec("update im_user set longitude=?,latitude=? where id=?", longitude, latitude, userId)
	fmt.Fprintf(w, result.ResponseWrite(1, "update finished"))
}

func FindNearby(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	longitude := parse.StringToFloat64(query.Get("longitude"))
	latitude := parse.StringToFloat64(query.Get("latitude"))
	userId := parse.StringToInt64(query.Get("userId"))
	sex := query.Get("sex")
	pageNo := parse.StringToInt(query.Get("pageNo"))
	memList := redis.Georadius([]byte("user_geo"), longitude, latitude, 100, "km")
	var memberIds []string
	for _, item := range memList {
		memberIds = append(memberIds, string(item.Member))
	}
	idStr := strings.Join(memberIds, ",")
	if idStr == "" {
		idStr = query.Get("userId")
	}

	sexStr := "and sex =\"" + sex + "\""
	if sex == "" {
		sexStr = "and 1=1"
	}
	if strings.EqualFold(sex, "女") {
		sexStr = "and (sex =\"" + sex + "\" or sex is null)"
	}
	if pageNo == 0 {
		pageNo = 1
	}
	pg := page.Page{
		PageNo:   pageNo,
		PageSize: 20,
	}
	list := model_mysql.QueryList("select id,name,sex,longitude,latitude,sign,headUrl from im_user where id in("+idStr+") and id<>? "+sexStr+"limit ?,?", userId, pageNo*20, 20)
	pg.List = list
	model_mysql.Exec("update im_user set longitude=?,latitude=? where id=?", longitude, latitude, userId)
	fmt.Fprintf(w, result.ResponseWriteData(1, pg))
}

func GetMsgHistory(w http.ResponseWriter, r *http.Request) {
	destType := parse.StringToInt64(r.PostFormValue("destType"))
	size := parse.StringToInt(r.PostFormValue("size"))
	sendTime := parse.StringToInt64(r.PostFormValue("sendTime"))
	userId := parse.StringToInt64(r.PostFormValue("userId"))
	if size == 0 {
		size = 20
	}
	if sendTime > 0 {
		data := model_mysql.QueryList("select * from im_message_view where belongId=? and sendTime<? and fromType=? order by id desc limit 0,?", userId, sendTime, destType, size)
		fmt.Fprintf(w, result.ResponseWriteData(1, data))
	} else {
		data := model_mysql.QueryList("select * from im_message_view where belongId=? and fromType=? order by id desc limit 0,?", userId, destType, size)
		fmt.Fprintf(w, result.ResponseWriteData(1, data))
	}

}
