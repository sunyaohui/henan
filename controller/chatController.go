package controller

import (
	"maoguo/henan/misc/redis"
	"maoguo/henan/misc/utils"
	"maoguo/henan/misc/utils/parse"
	"maoguo/henan/model/model_mysql/imUserModel"
	"maoguo/henan/model/model_mysql/imUserPayModel"
	"maoguo/henan/result"
	"maoguo/henan/services/impl"
	"maoguo/henan/vo"
	"reflect"
	"strings"

	"github.com/wonderivan/logger"
)

func Ready(token string) vo.ResultVo {
	data := make(map[string]interface{}, 0)
	userId := parse.StringToInt64(utils.Decrypt(token))
	imUser := new(impl.UserServiceImpl).GetImUser(userId)
	if &imUser != nil {
		imUser.Pwd = ""
	}
	payInfo := new(impl.ImUserPayServiceImpl).GetImUserPay(userId)
	if !reflect.DeepEqual(payInfo, &imUserPayModel.ImUserPay{}) {
		data["payInfo"] = payInfo.PayPwd
	} else {
		data["payInfo"] = ""
	}
	data["topList"] = new(impl.ImTopServiceImpl).GetTopList(userId)
	user := imUserModel.GetUserById(userId)
	info := parse.StructToJsonMap(imUser)
	info["IMNo"] = user.ImNumber
	data["myInfo"] = info
	data["friendsInfo"] = new(impl.ImFriendServiceImpl).GetFriendInfo(userId)
	//data["groupsInfo"] = new(impl.GroupServiceImpl).GetGroupsInfo(userId)
	data["friends"] = new(impl.CommonServiceImpl).GetUserFriends(userId)
	new(impl.GroupServiceImpl).CheckUserGroup(userId)
	return result.ResponseData(1, data)

}

func DoLogin(pwd, mobile, code, UUID, device string) vo.ResultVo {
	logger.Info("params>>>>  pwd:%d,mobile:%d,code:%d,UUID:%d,device:%d", pwd, mobile, code, UUID, device)
	if !strings.Contains(mobile, "+86") {
		mobile = code + mobile
	}
	mo := strings.ReplaceAll(mobile, "+86", "")
	user := imUserModel.Raw("select * from im_user where mobile=? or imNumber=?", mobile, mo)
	logger.Info("get user By mobile or imNumber, params>>> mobile:%d,imNumber:%d", mobile, mo, parse.ParseJson(user))
	if &user == nil {
		return result.ResponseMessage(2, "用户不存在")
	}
	if !strings.EqualFold(pwd, user.Pwd) {
		return result.ResponseMessage(3, "密码错误")
	}
	if UUID != "" {
		bytes := redis.HGetBytes([]byte("loginInfo"), []byte(parse.Int64ToString(user.Id)))
		var rUUID string
		if bytes == nil {
			rUUID = ""
		} else {
			rUUID = string(bytes[:])
		}
		if strings.EqualFold("", rUUID) && strings.EqualFold(UUID, rUUID) {
			redis.HSetBytes([]byte("loginInfo"), []byte(parse.Int64ToString(user.Id)), []byte(UUID))
			//同一设备不通知下线
		} else {
			redis.HSetBytes([]byte("loginInfo"), []byte(parse.Int64ToString(user.Id)), []byte(UUID))
			content := map[string]interface{}{"UUID": UUID}
			mess := new(impl.MessageServiceImpl).OtherLoginNotice(user.Id, parse.ParseJson(content))
			new(impl.MsgStoreAndNotifyServiceImpl).StoreMsgAndNotifyImServer(mess, user.Id)
		}
	} else {
		new(impl.MessageServiceImpl).NoticeOffline(parse.Int64ToString(user.Id))
	}
	imServerUrl := new(impl.MsgStoreAndNotifyServiceImpl).GetImServiceUrl()
	token := utils.Encrypt(parse.Int64ToString(user.Id))
	logger.Info("token:%d,   id:%d", token, parse.Int64ToString(user.Id))
	return result.ResponseData(1, map[string]interface{}{"imServerUrl": imServerUrl, "token": token})
}
