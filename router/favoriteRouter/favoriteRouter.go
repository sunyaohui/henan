package favoriteRouter

import (
	"fmt"
	"maoguo/henan/misc/page"
	"maoguo/henan/misc/utils/parse"
	"maoguo/henan/model/model_mysql"
	"maoguo/henan/model/model_mysql/imFavoriteModel"
	"maoguo/henan/result"
	"maoguo/henan/services/impl"
	"net/http"
)

//添加收藏
func AddFavorite(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	userId := parse.StringToInt64(query.Get("userId"))
	category := parse.StringToInt(query.Get("category"))
	fromId := parse.StringToInt64(query.Get("fromId"))
	content := query.Get("content")
	fromName := query.Get("fromName")
	fromHeadUrl := query.Get("fromHeadUrl")
	if userId <= 0 || category <= 0 || fromId <= 0 {
		fmt.Fprint(w, result.ResponseWrite(-1, "收藏失败，参数不正确"))
	}
	res := new(impl.FavoriteServiceImpl).AddFavorite(userId, fromId, category, content, fromName, fromHeadUrl)
	if res {
		fmt.Fprintf(w, result.ResponseWrite(1, "收藏成功"))
		return
	}
	fmt.Fprintf(w, result.ResponseWrite(-1, "收藏失败，重复收藏"))
}

func GetFavorite(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	userId := parse.StringToInt64(query.Get("userId"))
	pageNo := parse.StringToInt(query.Get("pageNo"))
	if pageNo == 0 {
		pageNo = 1
	}
	pg := page.Page{
		PageNo:   pageNo,
		PageSize: 20,
	}
	imFavoriteModel.GetFavoritePage(map[string]interface{}{"user_id": userId}, &pg)
	fmt.Fprintf(w, result.ResponseWriteData(1, pg))
}

func DelFavorite(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	// userId := parse.StringToInt64(query.Get("userId"))
	favoriteId := parse.StringToInt64(query.Get("favoriteId"))
	model_mysql.Exec("delete from im_favorite where id=?", favoriteId)
	fmt.Fprintf(w, result.ResponseWrite(1, "删除成功"))

}
