package main

import (
	"fmt"
	"maoguo/henan/misc/config"
	"maoguo/henan/model/model_mysql"
	"maoguo/henan/model/mongo"
	"maoguo/henan/result"
	"maoguo/henan/router/chatrouter"
	"maoguo/henan/router/favoriteRouter"
	"maoguo/henan/router/feedRouter"
	"maoguo/henan/router/georouter"
	"net/http"
	"strings"

	"github.com/wonderivan/logger"
)

//入口函数
func main() {

	logStart()
	config.InitConfig()
	if logJsonAddr, ok := config.CONFIG["LogJsonAddr"]; ok {
		logger.SetLogger(logJsonAddr)
	}
	model_mysql.Init_db()

	mongo.Init_Mongo()

	InitRouter()

}

//需要拦截的接口请直接调用此函数
func HTTPInterceptor(url, method string, h http.HandlerFunc) (string, http.HandlerFunc) {
	return url, http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			//拦截请求方式
			if method != "" && !strings.EqualFold(strings.ToUpper(method), r.Method) {
				AssertMethod(w, r)
				return
			}
			//登陆拦截
			// if !inteceptor.AssertLogin(r.Header["Token"], r.Method) {
			// 	AssertLogin(w, r)
			// 	return
			// }
			//放行
			logger.Info("Request Url:", url)
			h(w, r)
		})
}

//未登陆的操作被拦截，将由此函数统一返回
func AssertLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, result.ResponseWrite(-1, "请重新登陆"))
}

//请求方式不对，返回404
func AssertMethod(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 page not found"))
}

//路由注册
//需要拦截器请直接调用HTTPInterceptor 函数，参数未： url：路由，method:请求方式，h:请求函数
//不需要拦截器可直接按照http最原始方式调用函数接口；
func InitRouter() {
	//文件访问服务，输入路由加文件名可直接访问文件，直接访问路由将返回文件列表
	http.Handle("/upload/", http.StripPrefix("/upload/", http.FileServer(http.Dir(config.CONFIG["FILE_UPLOAD_ADDR"]))))

	//ready接口
	http.HandleFunc(HTTPInterceptor("/chat/ready", "GET", chatrouter.Ready))

	//登陆
	http.HandleFunc(HTTPInterceptor("/chat/doLogin", "POST", chatrouter.DoLogin))

	//上传图片
	http.HandleFunc(HTTPInterceptor("/chat/uploadImage", "POST", chatrouter.UploadImage))

	//s上传文件接口
	http.HandleFunc(HTTPInterceptor("/chat/doUploads", "POST", chatrouter.DoUploads))

	//修改个人信息
	http.HandleFunc(HTTPInterceptor("/chat/updateProfile", "POST", chatrouter.UpdateProfile))

	//修改备注
	http.HandleFunc(HTTPInterceptor("/chat/updateRemark", "POST", chatrouter.UpdateRemark))

	//修改群名称
	http.HandleFunc(HTTPInterceptor("/chat/updateGroup", "GET", chatrouter.UpdateGroup))

	//修改群成员备注
	http.HandleFunc(HTTPInterceptor("/chat/updateGroupMemberMark", "GET", chatrouter.UpdateGroupMemberMark))

	//修改隐私设置
	http.HandleFunc(HTTPInterceptor("/chat/updatePrivateSetting", "POST", chatrouter.UpdatePrivateSetting))

	//免打扰设置
	http.HandleFunc(HTTPInterceptor("/chat/setIgonre", "POST", chatrouter.SetIgonre))

	//创建群
	http.HandleFunc(HTTPInterceptor("/chat/createGroup", "POST", chatrouter.CreateGroup))

	//删除好友
	http.HandleFunc(HTTPInterceptor("/chat/delFriend", "POST", chatrouter.DelFriend))

	//加入黑名单
	http.HandleFunc(HTTPInterceptor("/chat/setBlock", "POST", chatrouter.SetBlock))

	//移除黑名单
	http.HandleFunc(HTTPInterceptor("/chat/removeBlack", "GET", chatrouter.RemoveBlock))

	//同意加好友或者加群
	http.HandleFunc(HTTPInterceptor("/chat/acceptRequest", "POST", chatrouter.AcceptRequest))

	//找回密码
	http.HandleFunc(HTTPInterceptor("/chat/findPwd", "POST", chatrouter.FindPwd))

	//摇一摇
	http.HandleFunc(HTTPInterceptor("/geo/findRandom", "GET", georouter.FindRandom))

	//添加收藏
	http.HandleFunc(HTTPInterceptor("/favorite/addFavorite", "GET", favoriteRouter.AddFavorite))

	//收藏列表(未完成)
	http.HandleFunc(HTTPInterceptor("/favorite/getFavorite", "GET", favoriteRouter.GetFavorite))

	//更新头像
	http.HandleFunc(HTTPInterceptor("/chat/updateHead", "POST", chatrouter.UpdateHead))

	//更新客户端所在位置
	http.HandleFunc(HTTPInterceptor("/geo/updateGeo", "GET", georouter.UpdateGeo))

	//朋友圈发布动态
	http.HandleFunc(HTTPInterceptor("/feed/addFeed", "POST", feedRouter.AddFeed))

	//朋友圈评论
	http.HandleFunc(HTTPInterceptor("/feed/addFeedComment", "POST", feedRouter.AddFeedComment))

	//朋友圈，发布赞
	http.HandleFunc(HTTPInterceptor("/feed/addFeedPraise", "GET", feedRouter.AddFeedPraise))

	//朋友圈，取消赞
	http.HandleFunc(HTTPInterceptor("/feed/canclePraise", "GET", feedRouter.CanclePraise))

	//朋友圈，我自己的相册或好友相册
	http.HandleFunc(HTTPInterceptor("/feed/getMyFeed", "GET", feedRouter.GetMyFeed))

	//朋友圈，查看好友或群成员个人资料时，里面的相册预览
	http.HandleFunc(HTTPInterceptor("/feed/getFeedAlbum", "GET", feedRouter.GetFeedAlbum))

	//朋友圈，点发现时，获取朋友圈列表
	http.HandleFunc(HTTPInterceptor("/feed/getFeed", "GET", feedRouter.GetFeed))

	//朋友圈，获取陌生人相册
	http.HandleFunc(HTTPInterceptor("/feed/getFriendFeed", "GET", feedRouter.GetFriendFeed))

	//通过手机号查询用户
	http.HandleFunc(HTTPInterceptor("/chat/queryUser", "GET", chatrouter.QueryUser))

	//查找群
	http.HandleFunc(HTTPInterceptor("/chat/queryGroup", "POST", chatrouter.QueryGroup))

	//查找附近的人
	http.HandleFunc(HTTPInterceptor("/geo/findNearby", "GET", georouter.FindNearby))

	//查询聊天记录
	//http.HandleFunc("/chat/getMsgHistory", chatrouter.Get)

	//注册
	http.HandleFunc(HTTPInterceptor("/chat/register", "POST", chatrouter.Register))

	//用户认证
	http.HandleFunc(HTTPInterceptor("/chat/updateAuth", "GET", chatrouter.UpdateAuth))

	//申请加好友，或者申请加群
	http.HandleFunc(HTTPInterceptor("/chat/requestFriend", "GET", chatrouter.RequestFriend))

	//取消置顶
	http.HandleFunc(HTTPInterceptor("/chat/cancleTop", "GET", chatrouter.CancleTop))

	//置顶
	http.HandleFunc(HTTPInterceptor("/chat/setTop", "POST", chatrouter.SetTop))

	//获取群信息
	http.HandleFunc(HTTPInterceptor("/chat/getGroupById", "POST", chatrouter.GetGroupById))

	//获取群成�����
	http.HandleFunc(HTTPInterceptor("/chat/getGroupMember", "POST", chatrouter.GetGroupMember))

	//获取验证码
	http.HandleFunc(HTTPInterceptor("/chat/getValidateNum", "GET", chatrouter.GetValidateNum))

	//退群
	http.HandleFunc(HTTPInterceptor("/chat/quitGroup", "GET", chatrouter.QuitGroup))

	//邀请好友加群
	http.HandleFunc(HTTPInterceptor("/chat/addGroupMember", "GET", chatrouter.AddGroupMember))

	//通过用户id查询资料
	http.HandleFunc(HTTPInterceptor("/chat/getUserById", "GET", chatrouter.GetUserById))

	//
	http.HandleFunc(HTTPInterceptor("/chat/inviteContact", "POST", chatrouter.InviteContact))

	http.HandleFunc(HTTPInterceptor("/feed/getBackImage", "GET", feedRouter.GetBackImage))

	//删除朋友圈
	http.HandleFunc(HTTPInterceptor("/feed/deleteFeed", "GET", feedRouter.DeleteFeed))

	//朋友圈获取单条动态
	http.HandleFunc(HTTPInterceptor("/feed/getOneFeed", "GET", feedRouter.GetOneFeed))

	//更换项目背景
	http.HandleFunc(HTTPInterceptor("/feed/changeBackImage", "GET", feedRouter.ChangeBackImage))

	//转让群
	http.HandleFunc(HTTPInterceptor("/chat/transGroup", "POST", chatrouter.TransGroup))

	//设置群头像
	http.HandleFunc(HTTPInterceptor("/chat/setGroupHeader", "POST", chatrouter.SetGroupHeader))

	//获取群公告
	http.HandleFunc(HTTPInterceptor("/chat/getNoteList", "POST", chatrouter.GetNoteList))

	//发布群公告
	http.HandleFunc(HTTPInterceptor("/chat/updateGroupNote", "GET", chatrouter.UpdateGroupNote))

	//解散群
	http.HandleFunc(HTTPInterceptor("/chat/dismissGroup", "POST", chatrouter.DismissGroup))

	//修改群介绍
	http.HandleFunc(HTTPInterceptor("/chat/updateGroupDesc", "POST", chatrouter.UpdateGroupDesc))

	//获取群配置
	http.HandleFunc(HTTPInterceptor("/chat/queryGroupConfig", "POST", chatrouter.QueryGroupConfig))

	//	申请加入群组
	http.HandleFunc(HTTPInterceptor("/chat/requestGroupJoin", "POST", chatrouter.RequestGroupJoin))

	//		删除收藏
	http.HandleFunc(HTTPInterceptor("/favorite/delFavorite", "GET", favoriteRouter.DelFavorite))

	//		设置朋友圈权限
	http.HandleFunc(HTTPInterceptor("/feed/setFeedAuth", "POST", feedRouter.SetFeedAuth))

	http.HandleFunc(HTTPInterceptor("/chat/transGroupConfirm", "POST", chatrouter.TransGroupConfirm))

	http.HandleFunc(HTTPInterceptor("/chat/updateJoinstatus", "POST", chatrouter.UpdateJoinstatus))

	http.ListenAndServe(":"+config.CONFIG["PORT"], nil)
}

func logStart() {
	logger.Info("")
	logger.Info("")
	logger.Info("                       _oo0oo_")
	logger.Info("                      o8888888o")
	logger.Info("                      88\" . \"88")
	logger.Info("                      (| -_- |)")
	logger.Info("                      0\\  =  /0")
	logger.Info("                    ___/`---'\\___")
	logger.Info("                  .' \\\\|     |// '.")
	logger.Info("                 / \\\\|||  :  |||// \\")
	logger.Info("                / _||||| -:- |||||- \\")
	logger.Info("               |   | \\\\\\  -  /// |   |")
	logger.Info("               | \\_|  ''\\---/''  |_/ |")
	logger.Info("               \\  .-\\__  '-'  ___/-. /")
	logger.Info("             ___'. .'  /--.--\\  `. .'___")
	logger.Info("          .\"\" '<  `.___\\_<|>_/___.' >' \"\".")
	logger.Info("         | | :  `- \\`.;`\\ _ /`;.`/ - ` : | |")
	logger.Info("         \\  \\ `_.   \\_ __\\ /__ _/   .-` /  /")
	logger.Info("     =====`-.____`.___ \\_____/___.-`___.-'=====")
	logger.Info("                       `=---='")
	logger.Info("")
	logger.Info("")
	logger.Info("  ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
	logger.Info("")
	logger.Info("               佛祖保佑         永无BUG")
	logger.Info("")
}
