package feedRouter

import (
	"fmt"
	"maoguo/henan/misc/page"
	"maoguo/henan/misc/utils/parse"
	"maoguo/henan/model/model_mysql"
	"maoguo/henan/model/model_mysql/imFeedModel"
	"maoguo/henan/model/model_mysql/imUserModel"
	"maoguo/henan/model/mongo"
	"maoguo/henan/result"
	"maoguo/henan/services/impl"
	"net/http"
	"reflect"
)

func AddFeed(w http.ResponseWriter, r *http.Request) {
	images := r.PostFormValue("images")
	text := r.PostFormValue("text")
	userId := parse.StringToInt64(r.PostFormValue("userId"))
	lat := r.PostFormValue("lat")
	lng := r.PostFormValue("lng")
	address := r.PostFormValue("address")
	priv := parse.StringToInt(r.PostFormValue("priv"))
	at := r.PostFormValue("at")
	uids := r.PostFormValue("uids")
	video := r.PostFormValue("video")
	ext := r.PostFormValue("ext")
	feed := new(impl.FeedServiceImpl).AddFeed(images, text, userId, lat, lng, address, priv, at, uids, video, ext)
	if !reflect.DeepEqual(feed, &imFeedModel.ImFeed{}) && userId > 0 && feed.Id > 0 {
		feed2 := new(impl.FeedServiceImpl).GetOneFeed(feed.Id, userId)
		fmt.Fprintf(w, result.ResponseWriteData(1, map[string]interface{}{"list": []mongo.ImFeed{feed2}}))
	} else {
		fmt.Fprintf(w, result.ResponseWrite(-1, "发布失败"))
	}
}

func AddFeedComment(w http.ResponseWriter, r *http.Request) {
	feedId := parse.StringToInt64(r.PostFormValue("feedId"))
	text := r.PostFormValue("text")
	userId := parse.StringToInt64(r.PostFormValue("userId"))
	replyId := parse.StringToInt64(r.PostFormValue("replyId"))
	if feedId == 0 || userId == 0 || text == "" {
		fmt.Fprintf(w, result.ResponseWrite(-1, "评论发布失败"))
	} else {
		new(impl.FeedServiceImpl).AddFeedComment(feedId, text, userId, replyId)
		fmt.Fprintf(w, result.ResponseWrite(1, "评论发布成功"))
	}
}

func AddFeedPraise(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	feedId := parse.StringToInt64(query.Get("feedId"))
	userId := parse.StringToInt64(query.Get("userId"))

	if feedId == 0 || userId == 0 {
		fmt.Fprintf(w, result.ResponseWrite(-1, "发布失败"))
	} else {
		new(impl.FeedServiceImpl).AddFeedPraise(feedId, userId)
		fmt.Fprintf(w, result.ResponseWrite(1, "赞发布成功"))
	}
}

func CanclePraise(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	feedId := parse.StringToInt64(query.Get("feedId"))
	userId := parse.StringToInt64(query.Get("userId"))
	if feedId == 0 || userId == 0 {
		fmt.Fprintf(w, result.ResponseWrite(-1, "取消赞失败"))
	} else {
		new(impl.FeedServiceImpl).CanclePraise(feedId, userId)
		fmt.Fprintf(w, result.ResponseWrite(1, "取消赞成功"))
	}
}

func GetMyFeed(w http.ResponseWriter, r *http.Request) {
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
	new(impl.FeedServiceImpl).GetMyFeed(userId, &pg)
	fmt.Fprintf(w, result.ResponseWriteData(1, pg))
}

func GetFeedAlbum(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	// userId := parse.StringToInt64(r.PostFormValue("userId"))
	destId := parse.StringToInt64(query.Get("destId"))
	o := model_mysql.Query("select GROUP_CONCAT(feed_imgs) album from(select feed_imgs  from im_feed where user_id=? and feed_imgs <>'' ORDER BY create_time desc LIMIT 3)a", destId)
	if o != nil {
		fmt.Fprint(w, result.ResponseWriteData(1, o))
	} else {
		fmt.Fprintf(w, result.ResponseWrite(0, "没有相册数据"))
	}
}

func GetFeed(w http.ResponseWriter, r *http.Request) {
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
	new(impl.FeedServiceImpl).GetFeed(userId, &pg)
	fmt.Fprintf(w, result.ResponseWriteData(1, pg))
}

func GetFriendFeed(w http.ResponseWriter, r *http.Request) {
	friendId := parse.StringToInt64(r.PostFormValue("friendId"))
	pageNo := parse.StringToInt(r.PostFormValue("pageNo"))
	if pageNo == 0 {
		pageNo = 1
	}
	pg := page.Page{
		PageNo:   pageNo,
		PageSize: 20,
	}
	list := imFeedModel.Raws("select * from im_feed where userId=? order by createTime desc limit ?,?", friendId, (pageNo-1)*20, 20)
	pg.List = list
	fmt.Fprintf(w, result.ResponseWriteData(1, pg))
}

//获取用户背景
func GetBackImage(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	userId := parse.StringToInt64(query.Get("userId"))
	user := imUserModel.Raw("select IFNULL(feedBackImage,'') feedBackImage from im_user where id=? ", userId)
	if &user != nil {
		fmt.Fprintf(w, result.ResponseWriteData(1, user.FeedBackImage))
		return
	}
	fmt.Fprintf(w, result.ResponseWrite(-1, "未找到用户"))
}

func DeleteFeed(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	feedId := parse.StringToInt64(query.Get("feedId"))
	userId := parse.StringToInt64(query.Get("userId"))
	if feedId == 0 || userId == 0 {
		fmt.Fprintf(w, result.ResponseWrite(-1, "删除失败，参数不正确"))
		return
	}
	new(impl.FeedServiceImpl).DeleteFeed(feedId, userId)
	fmt.Fprintf(w, result.ResponseWrite(1, "删除成功"))
}

//朋友圈获取单条动态
func GetOneFeed(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	feedId := parse.StringToInt64(query.Get("feedId"))
	userId := parse.StringToInt64(query.Get("userId"))
	if feedId == 0 || userId == 0 {
		fmt.Fprintf(w, result.ResponseWrite(-1, "失败，参数不正确"))
		return
	}
	feed := new(impl.FeedServiceImpl).GetOneFeed(feedId, userId)
	if reflect.DeepEqual(feed, imFeedModel.ImFeed{}) {
		fmt.Fprintf(w, result.ResponseWrite(-1, "该动态已被删除"))
		return
	}
	fmt.Fprintf(w, result.ResponseWriteData(1, feed))

}

//修改相册背景
func ChangeBackImage(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	userId := parse.StringToInt64(query.Get("userId"))
	imgUrl := query.Get("imgUrl")
	rows := model_mysql.ExecInt("update im_user set feedBackImage=? where id=?", imgUrl, userId)
	if rows > 0 {
		fmt.Fprintf(w, result.ResponseWrite(1, "更新背景成功"))
		return
	}
	fmt.Fprintf(w, result.ResponseWrite(-1, "更新背景失败"))
}

func SetFeedAuth(w http.ResponseWriter, r *http.Request) {
	userId := parse.StringToInt64(r.PostFormValue("userId"))
	friendId := parse.StringToInt64(r.PostFormValue("friendId"))
	priv := parse.StringToInt(r.PostFormValue("priv"))
	direct := parse.StringToInt(r.PostFormValue("direct"))
	if new(impl.FeedServiceImpl).SetFeedAuth(userId, friendId, priv, direct) > 0 {
		fmt.Fprintf(w, result.ResponseWrite(1, "设置成功"))
		return
	}
	fmt.Fprintf(w, result.ResponseWrite(-1, "设置失败"))
}
