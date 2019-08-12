package chatrouter

import (
	"fmt"
	"io"
	"maoguo/henan/constants"
	"maoguo/henan/constants/imConstants"
	"maoguo/henan/controller"
	"maoguo/henan/misc/config"
	"maoguo/henan/misc/redis"
	"maoguo/henan/misc/sms"
	"maoguo/henan/misc/utils"
	"maoguo/henan/misc/utils/goMap"
	"maoguo/henan/misc/utils/parse"
	"maoguo/henan/model/model_mysql"
	"maoguo/henan/model/model_mysql/imGroupMemberModel"
	"maoguo/henan/model/model_mysql/imGroupModel"
	"maoguo/henan/model/model_mysql/imTopModel"
	"maoguo/henan/model/model_mysql/imUserModel"
	"maoguo/henan/result"
	"maoguo/henan/services/impl"
	"net/http"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/wonderivan/logger"
)

func DoLogin(w http.ResponseWriter, r *http.Request) {
	device := r.PostFormValue("device")
	if strings.EqualFold("0", device) {
		fmt.Fprint(w, result.ResponseWrite(-1, "拒绝老用户登陆"))
		return
	}
	pwd := r.PostFormValue("pwd")
	mobile := r.PostFormValue("mobile")
	code := r.PostFormValue("code")
	UUID := r.PostFormValue("UUID")
	result := controller.DoLogin(pwd, mobile, code, UUID, device)
	fmt.Fprint(w, parse.ParseJson(result))
}

func Ready(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	token := query.Get("token")
	result := controller.Ready(token)
	fmt.Fprint(w, parse.ParseJson(result))
}

//修改个人信息
func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userId := r.PostFormValue("userId")
	if userId == "" {
		userId = "0"
	}
	nickName := r.PostFormValue("nickName")
	city := r.PostFormValue("city")
	province := r.PostFormValue("province")
	name := r.PostFormValue("name")
	sign := r.PostFormValue("sign")
	sex := r.PostFormValue("sex")
	district := r.PostFormValue("district")
	user := new(impl.UserServiceImpl).UpdateProfile(userId, nickName, city, province, name, sign, sex, district)
	if &user != nil {
		user.Pwd = ""
		fmt.Fprintf(w, result.ResponseWriteData(1, user))
	} else {
		fmt.Fprintf(w, result.ResponseWrite(-1, "更新失败"))
	}
}

func FindPwd(w http.ResponseWriter, r *http.Request) {
	validateNum := r.PostFormValue("validateNum")
	mobile := r.PostFormValue("mobile")
	pwd := r.PostFormValue("pwd")
	if !sms.EqualValidate(mobile, validateNum) {
		fmt.Fprintf(w, result.ResponseWrite(-1, "验证码不正确"))
		return
	}
	user := imUserModel.GetUser(map[string]interface{}{"mobile": mobile})
	if &user != nil {
		user.Pwd = utils.MD5(pwd)
		imUserModel.UpdateUser(user)
		fmt.Fprintf(w, result.ResponseWriteData(1, user))
		return
	}
	fmt.Fprintf(w, result.ResponseWrite(-1, "操作失败"))
}

//上传文件
func DoUploads(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	r.ParseForm()
	headUrl := "/upload/" + parse.Int64ToString(time.Now().Unix()) + "." + query.Get("fileExt") //+ r.PostFormValue("fileExt")
	fmt.Println(headUrl)
	realPath := config.CONFIG["FILE_UPLOAD_ADDR"]
	buf := make([]byte, 100000)
	file, err := os.Create(realPath + headUrl)
	defer file.Close()
	if err != nil {
		fmt.Println(err)
		fmt.Fprintf(w, result.ResponseWrite(-1, "上传失败"))
		return
	}
	defer r.Body.Close()
	for {
		n, err := r.Body.Read(buf)
		if err != nil {
			break
		}
		file.Write(buf[:n])
	}
	url := config.CONFIG["FILE_REQUEST_ADDR"] + headUrl
	fmt.Fprintf(w, result.ResponseWrite(1, url))
}

type Sizer interface {
	Size() int64
}

//上传图片
func UploadImage(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(64 << 20)
	var fileList []string
	path := config.CONFIG["FILE_UPLOAD_ADDR"]
	for _, fileHeader := range r.MultipartForm.File["uploadfile"] {
		srcFile, err := fileHeader.Open()
		defer srcFile.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
		srcName := strings.ToLower(fileHeader.Filename)
		// isImage := false
		// if strings.Contains(srcName, ".png") || strings.Contains(srcName, ".jpg") || strings.Contains(srcName, ".jpeg") {
		// 	isImage = true
		// }
		fileName := parse.Int64ToString(time.Now().Unix()) + "_" + utils.GetRand(6) + srcName
		if err != nil {
			fmt.Println(err)
			return
		}
		// var f float64 = 1024
		// size := (parse.Int64ToFloat64(srcFile.(Sizer).Size())) / f
		if srcName != "" {
			file, err := os.Create(path + "/" + fileName)
			if err != nil {
				fmt.Println(err)
			}
			io.Copy(file, srcFile)
			// if size > 1000 && isImage {
			// 	image.Scale(path+"/fileName", path+"/small_"+fileName, 400, 0)
			// 	fileName = "small_" + fileName
			// }
			fileList = append(fileList, config.CONFIG["FILE_REQUEST_ADDR"]+"/upload/"+fileName)
		}
	}
	fmt.Fprint(w, result.ResponseWriteData(1, fileList))
}

/**
修改备注
*/
func UpdateRemark(w http.ResponseWriter, r *http.Request) {
	destId := r.PostFormValue("destId")
	userId := r.PostFormValue("userId")
	remark := r.PostFormValue("remark")
	if parse.StringToInt64(destId) > 0 && parse.StringToInt64(userId) > 0 {
		model_mysql.Exec("update im_friend set remark=? where userId=? and friendId=?", remark, userId, destId)
	}
	//此段代码任何
	new(impl.ImFriendServiceImpl).CacheFriendInfo(parse.StringToInt64(userId))
	fmt.Fprint(w, result.ResponseWrite(1, "设置成功"))
}

//修改群名称
func UpdateGroup(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	userId := query.Get("userId")
	groupId := query.Get("groupId")
	name := query.Get("name")
	if new(impl.GroupServiceImpl).UpdateGroup(parse.StringToInt64(userId), parse.StringToInt64(groupId), name) > 0 {
		fmt.Fprintf(w, result.ResponseWrite(1, "更新成功"))
		return
	}
	fmt.Fprintf(w, result.ResponseWrite(-1, "各参数不能为空"))
}

/**
修改群成员备注
*/
func UpdateGroupMemberMark(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	userId := parse.StringToInt64(query.Get("userId"))
	groupId := parse.StringToInt64(query.Get("groupId"))
	markName := query.Get("markName")
	if new(impl.GroupServiceImpl).UpdateGroupMemberMark(userId, groupId, markName) > 0 {
		fmt.Fprintf(w, result.ResponseWrite(1, "更新成功"))
		return
	}
	fmt.Fprintf(w, result.ResponseWrite(-1, "各参数不能为空"))
}

/**
修改隐私设置
*/
func UpdatePrivateSetting(w http.ResponseWriter, r *http.Request) {
	needAuth := parse.StringToInt(r.PostFormValue("needAuth"))
	newNotification := parse.StringToInt(r.PostFormValue("newNotification"))
	searchMobile := parse.StringToInt(r.PostFormValue("searchMobile"))
	userId := parse.StringToInt64(r.PostFormValue("userId"))
	setting := new(impl.UserServiceImpl).GetImUser(userId)
	setting.NeedAuth = needAuth
	setting.NewNotification = newNotification
	setting.SearchMobile = searchMobile
	imUserModel.UpdateUser(setting)
	fmt.Fprintf(w, result.ResponseWrite(1, "更新成功"))
}

/**
免打扰设置
*/
func SetIgonre(w http.ResponseWriter, r *http.Request) {
	receiveTip := parse.StringToInt(r.PostFormValue("receiveTip"))
	//receiveTip := parse.StringToInt(r.PostFormValue("receiveTip"))
	destId := parse.StringToInt64(r.PostFormValue("destId"))
	destType := parse.StringToInt(r.PostFormValue("destType"))
	userId := parse.StringToInt64(r.PostFormValue("userId"))
	if destType == 1 {
		model_mysql.Exec("update im_friend set receiveTip=? where userId=? and friendId=?", receiveTip, userId, destId)
	} else {
		model_mysql.Exec("update im_group_member set receiveTip=? where userId=? and groupId=?", receiveTip, userId, destId)
	}
	new(impl.ImFriendServiceImpl).CacheFriendInfo(userId)
	new(impl.GroupServiceImpl).CacheGroupsInfo(userId)
	new(impl.ImFriendServiceImpl).CacheFriendInfo(destId)
	new(impl.GroupServiceImpl).CacheGroupsInfo(destId)
	fmt.Fprintf(w, result.ResponseWrite(1, "设置成功"))
}

//创建群组
func CreateGroup(w http.ResponseWriter, r *http.Request) {
	ids := r.PostFormValue("ids")
	memberIds := parse.StringToInt64Arr(ids)
	name := r.PostFormValue("name")
	userId := parse.StringToInt64(r.PostFormValue("userId"))
	if name != "" && userId != 0 {

		//创建��组
		group := &imGroupModel.ImGroup{
			Name:       name,
			CreaterId:  userId,
			CreateTime: time.Now().Unix(),
		}
		imGroupModel.Save(group)
		//添加群主到群成员
		member := &imGroupMemberModel.ImGroupMember{
			CreateTime: time.Now().Unix(),
			CreatorId:  userId,
			GroupId:    group.Id,
			Role:       1,
			UserId:     userId,
			IsAccept:   1,
			ReceiveTip: 1,
		}
		imGroupMemberModel.Save(member)
		//此句���码���������做了查���，无任何其他��作，无意义
		new(impl.GroupServiceImpl).CacheGroupsInfo(userId)

		for _, memberId := range memberIds {
			msg := new(impl.MessageServiceImpl).InviteGroupNotice(memberId, group.Id)
			new(impl.MsgStoreAndNotifyServiceImpl).StoreMsgAndNotifyImServer(msg, memberId)
			//���������������������������������������������������员
			members := &imGroupMemberModel.ImGroupMember{
				CreateTime: time.Now().Unix(),
				CreatorId:  userId,
				GroupId:    group.Id,
				Role:       3,
				UserId:     memberId,
				IsAccept:   1,
				ReceiveTip: 1,
			}
			imGroupMemberModel.Save(members)

			new(impl.GroupServiceImpl).CacheGroupsInfo(memberId)
		}

		obj := model_mysql.Query("select b.id,b.name,a.receiveTip,headUrl from (select groupId,userId,receiveTip from im_group_member where userId=? and groupId=? and isaccept=1 group by groupId,userId,receiveTip) a  LEFT JOIN im_group b on(a.groupId=b.id)", userId, group.Id)
		fmt.Fprintf(w, result.ResponseWriteData(1, obj))
	} else {
		fmt.Fprintf(w, result.ResponseWrite(-1, "各参数不能为空"))
	}
}

//删除好友
func DelFriend(w http.ResponseWriter, r *http.Request) {
	userId := parse.StringToInt64(r.PostFormValue("userId"))
	if userId > 0 {
		destId := parse.StringToInt64(r.PostFormValue("destId"))
		new(impl.UserServiceImpl).DelUser(userId, destId)
		fmt.Fprintf(w, result.ResponseWrite(1, "操作成功"))
	} else {
		fmt.Fprintf(w, result.ResponseWrite(-1, "获取失败"))
	}
}

//加入黑名单
func SetBlock(w http.ResponseWriter, r *http.Request) {
	userId := parse.StringToInt64(r.PostFormValue("userId"))
	blockId := parse.StringToInt64(r.PostFormValue("blockId"))
	model_mysql.Exec("update im_friend set isBlack=1 where userId=? and friendId=?", userId, blockId)
	new(impl.ImFriendServiceImpl).CacheFriendInfo(userId)
	fmt.Fprintf(w, result.ResponseWrite(1, "设置成功"))
}

func RemoveBlock(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	userId := parse.StringToInt64(query.Get("userId"))
	blockId := parse.StringToInt64(query.Get("destId"))
	if userId > 0 && blockId > 0 {
		new(impl.UserServiceImpl).RemoveBlock(userId, blockId)
	}
	fmt.Fprintf(w, result.ResponseWrite(1, "移出成功"))
}

//同意添加为好友
func AcceptRequest(w http.ResponseWriter, r *http.Request) {
	destType := parse.StringToInt(r.PostFormValue("destType"))
	userId := parse.StringToInt64(r.PostFormValue("userId"))
	destId := parse.StringToInt64(r.PostFormValue("destId"))
	//groupId := this.PostFormInt64("groupId")
	if destType == 1 {
		msg := new(impl.MessageServiceImpl).AcceptFriendNotice(userId, destId)
		new(impl.MsgStoreAndNotifyServiceImpl).StoreMsgAndNotifyImServer(msg, destId)

		model_mysql.Exec("update im_friend set isFriend=1 where userId=? and friendId=?", userId, destId)
		model_mysql.Exec("update im_friend set isFriend=1 where userId=? and friendId=?", destId, userId)

		user := new(impl.UserServiceImpl).GetImUser(destId)
		user.Pwd = ""

		new(impl.ImFriendServiceImpl).CacheFriendInfo(userId)
		new(impl.ImFriendServiceImpl).CacheFriendInfo(destId)

		redis.Sadd([]byte("friends_"+parse.Int64ToString(userId)), parse.Int64ToBytes(destId))
		redis.Sadd([]byte("friends_"+parse.Int64ToString(destId)), parse.Int64ToBytes(userId))
		fmt.Fprintf(w, result.ResponseWriteData(1, user))
	} else {
		fmt.Fprintf(w, result.ResponseWrite(-1, "加好友失败"))
	}
}

func UpdateHead(w http.ResponseWriter, r *http.Request) {
	userId := parse.StringToInt64(r.PostFormValue("userId"))
	headUrl := r.PostFormValue("headUrl")
	user := new(impl.UserServiceImpl).GetImUser(userId)
	if &user != nil {
		user.HeadUrl = headUrl
		imUserModel.UpdateUser(user)
		new(impl.UserServiceImpl).CacheUser(&user)
		new(impl.FeedServiceImpl).UpdateFeed(userId, "", headUrl)
	}
	fmt.Fprintf(w, result.ResponseWrite(1, "头像更新完成"))
}

func QueryUser(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	mobile := query.Get("mobile")
	logger.Info("queryUser params:", mobile)
	// userId := r.PostFormValue("userId")
	list := model_mysql.QueryList("select id,name,nickName,mobile,mail,sex,birthday,sign,province,isOnline,city,district,createTime,status,detail,headUrl,imNumber as IMNo from im_user where (mobile like '%" + mobile + "%' or id = '" + mobile + "') and searchMobile=1 limit 0,20")
	fmt.Fprintf(w, result.ResponseWriteData(1, map[string]interface{}{"info": list}))
}

func QueryGroup(w http.ResponseWriter, r *http.Request) {
	groupName := r.PostFormValue("groupName")
	list := model_mysql.QueryList("select * from im_group where name like '%" + groupName + "%' or id= '" + groupName + "' limit 0,20")
	fmt.Fprintf(w, result.ResponseWriteData(1, map[string]interface{}{"info": list}))
}

func UpdateAuth(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	userId := parse.StringToInt64(query.Get("userId"))
	idNo := query.Get("idNo")
	name := query.Get("name")
	validateNum := query.Get("validateNum")
	mobile := query.Get("mobile")
	status := new(impl.UserServiceImpl).CheckIdNo(userId, idNo)
	if status == 2 {
		fmt.Fprintf(w, result.ResponseWrite(2, "身份证已被使用"))
		return
	}
	if status == 0 {
		fmt.Fprintf(w, result.ResponseWrite(0, "认证审核中"))
		return
	}
	if status == 1 {
		fmt.Fprintf(w, result.ResponseWrite(1, "认证成功"))
		return
	}
	if new(impl.UserServiceImpl).UpdateAuth(userId, name, validateNum, idNo, mobile) > 0 {
		fmt.Fprintf(w, result.ResponseWrite(1, "提交成功，请等待审核"))
		return
	}
	fmt.Fprintf(w, result.ResponseWrite(-1, "各参数不能为空"))
}

func RequestFriend(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	destType := parse.StringToInt(query.Get("destType"))
	userId := parse.StringToInt64(query.Get("userId"))
	destId := parse.StringToInt64(query.Get("destId"))
	if destType == imConstants.MSG_FROM_P2P {
		code := new(impl.ImFriendServiceImpl).RequestFriend(destType, userId, destId)
		fmt.Fprintf(w, result.ResponseWrite(code, "完成"))
		return
	}
	if imConstants.MSG_FROM_GROUP == destType {
		group := new(impl.GroupServiceImpl).RequestGroupJion(userId, destId, destType)
		if reflect.DeepEqual(&group, &imGroupModel.ImGroup{}) {
			fmt.Fprintf(w, result.ResponseWrite(-1, "加入群失败，可能超过人数限制700"))
			return
		}
		fmt.Fprintf(w, result.ResponseWriteData(1, group))
		return
	}
	fmt.Fprintf(w, result.ResponseWrite(-1, "可能参数不对"))
}

func SetTop(w http.ResponseWriter, r *http.Request) {
	destId := parse.StringToInt64(r.PostFormValue("destId"))
	userId := parse.StringToInt64(r.PostFormValue("userId"))
	destType := parse.StringToInt(r.PostFormValue("destType"))
	if destId > 0 && userId > 0 && destType > 0 {
		model_mysql.Exec("delete from im_top where userId=? and destId=? and destType=?", userId, destId, destType)
		top := imTopModel.ImTop{
			DestId:   destId,
			DestType: destType,
			UserId:   userId,
		}
		imTopModel.Save(&top)
		new(impl.ImTopServiceImpl).CacheTopList(userId)
	}
	fmt.Fprintf(w, result.ResponseWrite(1, "设置成功"))
}

func CancleTop(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	destId := parse.StringToInt64(query.Get("destId"))
	userId := parse.StringToInt64(query.Get("userId"))
	destType := parse.StringToInt(query.Get("destType"))
	if destId > 0 && userId > 0 && destType > 0 {
		model_mysql.Exec("delete from im_top where userId=? and destId=? and destType=?", userId, destId, destType)
	}
	fmt.Fprintf(w, result.ResponseWrite(1, "设置成功"))
}

func GetGroupById(w http.ResponseWriter, r *http.Request) {
	groupId := parse.StringToInt64(r.PostFormValue("groupId"))
	userId := parse.StringToInt64(r.PostFormValue("userId"))
	if userId > 0 {
		group := new(impl.GroupServiceImpl).GroupInfo(groupId)
		if reflect.DeepEqual(group, imGroupModel.ImGroup{}) {
			fmt.Fprintf(w, result.ResponseWrite(-1, "群不存在"))
			return
		}
		fmt.Fprintf(w, result.ResponseWriteData(1, group))
		return
	}
	fmt.Fprintf(w, result.ResponseWrite(-1, "	参数错误"))
}

func GetGroupMember(w http.ResponseWriter, r *http.Request) {
	logger.Info("-----------get group member start:----------------")
	groupId := parse.StringToInt64(r.PostFormValue("groupId"))
	logger.Info("---------params:", groupId)
	// userId := parse.StringToInt64(r.PostFormValue("userId"))
	list := new(impl.GroupServiceImpl).GetAllMemberFromCache(groupId)
	logger.Info("--------result:", parse.ParseJson(list))
	logger.Info("--------------------get group member end!!!!---------------------")
	fmt.Fprintf(w, result.ResponseWriteData(1, map[string]interface{}{"info": list}))
}

// func GetValidateNum(w http.ResponseWriter, r *http.Request) {
// 	mobile := r.PostFormValue("mobile")
// 	app := r.PostFormValue("app")
// 	validateNum := utils.GetRand(4)
// 	redis.SetkeyExPrire([]byte(constants.SMS_KEY), 60*30)
// 	redis.HSet(constants.SMS_KEY, mobile, validateNum)
// 	var integer int
// 	if strings.EqualFold(app, "samim") {
// 		integer = impl.SendSms(mobile, validateNum)
// 	} else {
// 		integer = impl.SendSmsApp(mobile, validateNum, app)
// 	}
// 	mString := "短信发送失败，请稍后重试"
// 	switch integer {
// 	case 2:
// 		mString = "手机号格式不正确"
// 		break
// 	case 22:
// 		mString = "1小时内只能获取3次验证码"
// 		break
// 	case 33:
// 		mString = "30秒内只能获取1次验证码"
// 		break
// 	case 20:
// 		mString = "不支持该地区"
// 		break
// 	case 43:
// 		mString = "今日验证码次数已达到上限"
// 		break
// 	case 3:
// 		mString = "发送失败，请联系客服ErrCode=3"
// 		break
// 	default:
// 		break
// 	}
// 	if integer == 0 {
// 		fmt.Fprintf(w, result.ResponseWriteData(1, validateNum))
// 	} else {
// 		fmt.Fprintf(w, result.ResponseWrite(-1, mString))
// 	}
// }
func GetValidateNum(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	mobile := query.Get("mobile")
	valiType := query.Get("type")
	if mobile == "" {
		fmt.Fprintf(w, result.ResponseWrite(-1, "手机号不能为空"))
		return
	}
	if strings.Contains(mobile, "+86") {
		mobile = strings.ReplaceAll(mobile, "+86", "")
	}
	validateNum := utils.GetSmsValcationCode()
	redis.SetkeyExPrire([]byte(constants.SMS_KEY), 60*5)
	redis.HSet(constants.SMS_KEY, mobile, validateNum)
	res := sms.SendValiCode(mobile, validateNum, valiType)
	if res["success"].(bool) {
		r := map[string]interface{}{
			"info": "发送成功",
		}
		fmt.Fprintf(w, result.ResponseWriteData(1, r))
	} else {
		r := map[string]interface{}{
			"info": res["desc"],
		}
		fmt.Fprintf(w, result.ResponseWriteData(-1, r))
	}
}

func QuitGroup(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	groupId := parse.StringToInt64(query.Get("groupId"))
	userId := parse.StringToInt64(query.Get("userId"))
	if groupId > 0 && userId > 0 {
		new(impl.GroupServiceImpl).UpdateQuitGroup(groupId, userId)
		fmt.Fprintf(w, result.ResponseWrite(1, "退出成功"))
	} else {
		fmt.Fprintf(w, result.ResponseWrite(-1, "退出失败"))
	}
}

func AddGroupMember(w http.ResponseWriter, r *http.Request) {
	groupId := parse.StringToInt64(r.PostFormValue("groupId"))
	userId := parse.StringToInt64(r.PostFormValue("userId"))
	inviteIds := r.PostFormValue("inviteIds")
	new(impl.GroupServiceImpl).InviteJoinGroup(inviteIds, groupId, userId)
	fmt.Fprintf(w, result.ResponseWrite(1, "邀请完成"))
}

func Register(w http.ResponseWriter, r *http.Request) {
	mobile := r.PostFormValue("mobile")
	pwd := r.PostFormValue("pwd")
	validateNum := r.PostFormValue("validateNum")
	headUrl := r.PostFormValue("headUrl")
	name := r.PostFormValue("name")
	recommandUserId := parse.StringToInt64(r.PostFormValue("recommandUserId"))
	if !sms.EqualValidate(mobile, validateNum) {
		fmt.Fprintf(w, result.ResponseWrite(3, "验证码不正确"))
		return
	}
	if pwd == "" && validateNum == "" {
		fmt.Fprintf(w, result.ResponseWrite(-1, "有参数不正确"))
		return
	}
	list := model_mysql.QueryList("select * from im_user where mobile=?", mobile)
	if len(list) > 0 {
		fmt.Fprintf(w, result.ResponseWrite(2, "手机号码已存在"))
		return
	}

	var user imUserModel.ImUser
	if name == "" {
		user.Name = mobile
	} else {
		user.Name = name
	}
	user.Pwd = utils.MD5(pwd)
	user.Mobile = mobile
	user.HeadUrl = headUrl
	user.IsOnline = 0
	user.CreateTime = time.Now().Unix()
	user.Status = 1
	user.NeedAuth = 1
	user.NewNotification = 0
	user.SearchMobile = 1
	user.NickName = name
	user.IsAuth = 0
	user.RecommandUserId = recommandUserId
	imUserModel.Save(&user)
	new(impl.UserServiceImpl).CacheUser(&user)
	fmt.Fprintf(w, result.ResponseWriteData(1, user))
}

func GetUserById(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	userId := parse.StringToInt64(query.Get("userId"))
	destId := parse.StringToInt64(query.Get("destId"))
	if userId > 0 && destId > 0 {
		user := new(impl.UserServiceImpl).GetImUser(destId)
		maps := goMap.NewMap()
		userJson := parse.StructToJsonMap(user)
		for k, v := range userJson {
			maps.Put(k, v)
		}
		maps.Put("pwd", "")
		feedAuths := new(impl.FeedServiceImpl).GetFeedAuth(userId, destId)
		for k, v := range feedAuths {
			maps.Put(k, v)
		}
		isNumbers := model_mysql.Query("select imNumber as IMNo from im_user where id=?", destId)
		for k, v := range isNumbers {
			logger.Info(k, v)
			maps.Put(k, v)
		}
		fmt.Fprintf(w, result.ResponseWriteData(1, maps.Map))
	} else {
		fmt.Fprintf(w, "发生异常")
	}
}

func InviteContact(w http.ResponseWriter, r *http.Request) {
	userId := parse.StringToInt64(r.PostFormValue("userId"))
	contact := r.PostFormValue("contact")
	data := make(map[string]interface{})
	if contact != "" {
		new(impl.UserServiceImpl).ImportContact(userId, contact)
		// data["info"] = new(impl.UserServiceImpl).MyContact(userId)
		// fmt.Fprintf(w, result.ResponseWriteData(1, data))
		// return
	}
	list := new(impl.UserServiceImpl).MyContact(userId)
	data["info"] = list
	if list == nil || len(list) == 0 {
		data["info"] = make([]map[string]interface{}, 0)
	}
	fmt.Fprintf(w, result.ResponseWriteData(1, data))
}

//发布群公告
func UpdateGroupNote(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	userId := parse.StringToInt64(query.Get("userId"))
	groupId := parse.StringToInt64(query.Get("groupId"))
	noteId := parse.StringToInt64(query.Get("noteId"))
	content := query.Get("content")
	title := query.Get("title")
	member := model_mysql.Query("select * from im_group_member where userId=? and groupId=? and isAccept=1", userId, groupId)
	if member != nil {
		role := member["role"]
		if !strings.EqualFold(role.(string), "1") && !strings.EqualFold(role.(string), "2") {
			fmt.Fprintf(w, result.ResponseWrite(-1, "无权操作"))
			return
		}
		if new(impl.GroupServiceImpl).UpdateGroupNote(userId, noteId, content, title) > 0 {
			fmt.Fprintf(w, result.ResponseWrite(1, "更新成功"))
			return
		}
		fmt.Fprintf(w, result.ResponseWrite(-1, "各参数不能为空"))
		return
	}
	fmt.Fprintf(w, result.ResponseWrite(-1, "您不在该群组，无法操作"))
}

//获取群公告
func GetNoteList(w http.ResponseWriter, r *http.Request) {
	userId := parse.StringToInt64(r.PostFormValue("userId"))
	groupId := parse.StringToInt64(r.PostFormValue("groupId"))
	pageNo := parse.StringToInt(r.PostFormValue("pageNo"))
	if pageNo == 0 {
		pageNo = 1
	}
	member := model_mysql.Query("select * from im_group_member where userId=? and groupId=? and isAccept=1", userId, groupId)
	if member != nil {
		role := member["role"]
		editable := 0
		if strings.EqualFold(role.(string), "1") || strings.EqualFold(role.(string), "2") {
			editable = 1
		}
		pg := new(impl.GroupServiceImpl).GetNoteList(userId, groupId, pageNo)
		data := map[string]interface{}{
			"list":       pg.List,
			"pageNo":     pageNo,
			"pageSize":   20,
			"totalCount": pg.TotalCount,
			"orderBy":    pg.OrderBy,
			"orderType":  pg.OrderType,
			"totalPage":  pg.TotalPage,
			"editable":   editable,
		}
		fmt.Fprintf(w, result.ResponseWriteData(1, data))
		return
	}
	fmt.Fprintf(w, result.ResponseWrite(-1, "您不在该群组，无法操作"))
}

//���������������
func SetGroupHeader(w http.ResponseWriter, r *http.Request) {
	userId := parse.StringToInt64(r.PostFormValue("userId"))
	groupId := parse.StringToInt64(r.PostFormValue("groupId"))
	url := r.PostFormValue("url")
	member := model_mysql.Query("select * from im_group_member where userId=? and groupId=? and isAccept=1", userId, groupId)
	if member != nil {
		role := member["role"]
		if !strings.EqualFold(role.(string), "1") && !strings.EqualFold(role.(string), "2") {
			fmt.Fprintf(w, result.ResponseWrite(-1, "无权操作"))
			return
		}
		model_mysql.Exec("update im_group set headUrl=? where id=?", url, groupId)
		new(impl.GroupServiceImpl).CacheGroupsInfo(userId)
		new(impl.GroupServiceImpl).SendNoticeToMembers(userId, groupId, url, imConstants.MSG_TYPE_GROUP_HEAD)
		fmt.Fprintf(w, result.ResponseWrite(1, "操作成功"))
		return
	}
	fmt.Fprintf(w, result.ResponseWrite(-1, "设置失败"))
}

//转让群
func TransGroup(w http.ResponseWriter, r *http.Request) {
	userId := parse.StringToInt64(r.PostFormValue("userId"))
	groupId := parse.StringToInt64(r.PostFormValue("groupId"))
	destId := parse.StringToInt64(r.PostFormValue("destId"))
	if userId > 0 && groupId > 0 && destId > 0 {
		fmt.Fprintf(w, new(impl.GroupServiceImpl).TransGroupRequest(userId, groupId, destId))
		return
	}
	fmt.Fprintf(w, result.ResponseWrite(-1, "操作失败"))
}

//解散群
func DismissGroup(w http.ResponseWriter, r *http.Request) {
	userId := parse.StringToInt64(r.PostFormValue("userId"))
	groupId := parse.StringToInt64(r.PostFormValue("groupId"))
	if userId > 0 && groupId > 0 {
		fmt.Fprintf(w, new(impl.GroupServiceImpl).DismissGroup(userId, groupId))
		return
	}
	fmt.Fprintf(w, result.ResponseWrite(-1, "操作失败"))
}

//修改群介绍
func UpdateGroupDesc(w http.ResponseWriter, r *http.Request) {
	userId := parse.StringToInt64(r.PostFormValue("userId"))
	groupId := parse.StringToInt64(r.PostFormValue("groupId"))
	descriptions := r.PostFormValue("descriptions")
	if userId > 0 && groupId > 0 {
		res := new(impl.GroupServiceImpl).UpdateGroupDescriptions(userId, groupId, descriptions)
		fmt.Fprintf(w, res)
		return
	}
	fmt.Fprintf(w, result.ResponseWrite(-1, "操作失败"))
}

func QueryGroupConfig(w http.ResponseWriter, r *http.Request) {
	list := new(impl.GroupServiceImpl).GetGroupConfig()
	if list != nil {
		fmt.Fprintf(w, result.ResponseWriteData(1, map[string]interface{}{"info": list}))
		return
	}
	fmt.Fprintf(w, result.ResponseWrite(-1, "获取配置信息失败，请稍后重试"))
}

func RequestGroupJoin(w http.ResponseWriter, r *http.Request) {
	userId := parse.StringToInt64(r.PostFormValue("userId"))
	groupId := parse.StringToInt64(r.PostFormValue("groupId"))
	answer := r.PostFormValue("answer")
	member := model_mysql.Query("select * from im_group_member where userId=? and groupId=? and isAccept=1", userId, groupId)
	if member != nil {
		fmt.Fprintf(w, result.ResponseWrite(-1, "您已经是群组成员"))
		return
	}
	flag := new(impl.GroupServiceImpl).RequestGroupJionWithQuestion(userId, groupId, answer)
	if flag != 0 {
		r := "网络繁忙，请稍后重试"
		switch flag {
		case 1:
			r = "验证消息不能为空"
			break
		case 2:
			r = "加群问题答案不正确"
			break
		case 3:
			r = "不允许加入"
			break
		case 4:
			r = "请先支付"
			break
		case 5:
			r = "群人数已经到达上限"
			break
		default:
			r = "网络繁忙，请稍后重试"
			break
		}
		fmt.Fprintf(w, result.ResponseWrite(-flag, r))
		return
	}
	fmt.Fprintf(w, result.ResponseWrite(1, "申请成功"))
}

func UpdateJoinstatus(w http.ResponseWriter, r *http.Request) {
	userId := parse.StringToInt64(r.PostFormValue("userId"))
	groupId := parse.StringToInt64(r.PostFormValue("groupId"))
	status := parse.StringToInt(r.PostFormValue("status"))
	question := r.PostFormValue("question")
	answer := r.PostFormValue("answer")
	fee := parse.StringToFloat64(r.PostFormValue("fee"))
	member := model_mysql.Query("select * from im_group_member where userId=? and groupId=? and isAccept=1", userId, groupId)
	if member != nil {
		if status < 0 || status > 0 {
			fmt.Fprintf(w, result.ResponseWrite(-1, "网络繁忙，请稍后重试"))
			return
		}
		if status == 2 && strings.EqualFold(question, "") {
			fmt.Fprintf(w, result.ResponseWrite(-1, "请配置正确的验证问题"))
			return
		}
		if status == 3 && (strings.EqualFold(question, "") || strings.EqualFold(answer, "")) {
			fmt.Fprintf(w, result.ResponseWrite(-1, "请配置正确的验证问题及答案"))
			return
		}

		role := parse.StringToInt(member["role"].(string))
		if role != 1 && role != 2 {
			fmt.Fprintf(w, result.ResponseWrite(-1, "无权操作"))
			return
		}
		if new(impl.GroupServiceImpl).UpdateGroupJoinStatus(groupId, status, question, answer, fee) > 0 {
			fmt.Fprintf(w, result.ResponseWrite(1, "设置成功"))
			return
		}

		fmt.Fprintf(w, result.ResponseWrite(-1, "设置失败"))
		return
	}
	fmt.Fprintf(w, result.ResponseWrite(-1, "您不在该群组，无法操作"))
}

func TransGroupConfirm(w http.ResponseWriter, r *http.Request) {
	userId := parse.StringToInt64(r.PostFormValue("userId"))
	groupId := parse.StringToInt64(r.PostFormValue("groupId"))
	accept := parse.StringToIntDefault(r.PostFormValue("accept"), 2)
	if userId > 0 && groupId > 0 {
		fmt.Fprintf(w, new(impl.GroupServiceImpl).TransGroupConfirm(userId, groupId, accept))
		return
	}
	fmt.Fprintf(w, result.ResponseWrite(-1, "操作失败"))
}

func SetGroupAdmin(w http.ResponseWriter, r *http.Request) {
	userId := parse.StringToInt64(r.PostFormValue("userId"))
	groupId := parse.StringToInt64(r.PostFormValue("groupId"))
	destId := parse.StringToInt64(r.PostFormValue("destId"))
	isSet := parse.StringToIntDefault(r.PostFormValue("isSet"), 1)
	if userId > 0 && groupId > 0 {
		fmt.Fprintf(w, new(impl.GroupServiceImpl).SetGroupAdmin(userId, groupId, destId, isSet))
		return
	}
	fmt.Fprintf(w, result.ResponseWrite(-1, "操作失败"))
}

func UpdateGroupMember(w http.ResponseWriter, r *http.Request) {
	userId := parse.StringToInt64(r.PostFormValue("userId"))
	groupId := parse.StringToInt64(r.PostFormValue("groupId"))
	destId := parse.StringToInt64(r.PostFormValue("destId"))
	types := parse.StringToInt(r.PostFormValue("type"))
	ope := parse.StringToInt(r.PostFormValue("ope"))
	group := new(impl.GroupServiceImpl).GroupInfo(groupId)
	if reflect.DeepEqual(group, imGroupModel.ImGroup{}) {
		fmt.Fprintf(w, result.ResponseWrite(-1, "群不存在"))
		return
	}
	fmt.Fprintf(w, new(impl.GroupServiceImpl).UpdateGroupMember(userId, groupId, destId, types, ope))
}
