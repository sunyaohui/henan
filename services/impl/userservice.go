package impl

import (
	"encoding/json"
	"maoguo/henan/misc/redis"
	"maoguo/henan/misc/sms"
	"maoguo/henan/misc/utils/parse"
	"maoguo/henan/misc/utils/serializ"
	"maoguo/henan/misc/utils/slice"
	"maoguo/henan/model/model_mysql"
	"maoguo/henan/model/model_mysql/imRockRecordModel"
	"maoguo/henan/model/model_mysql/imUserContactModel"
	"maoguo/henan/model/model_mysql/imUserModel"
	"maoguo/henan/paramter"
	"reflect"
	"strings"

	"github.com/wonderivan/logger"
)

type UserServiceImpl struct {
}

func (this *UserServiceImpl) GetImUser(userId int64) imUserModel.ImUser {
	key := paramter.GetImUserKey(userId)
	data := redis.GetBytes(key)
	if data != nil {
		var imUser imUserModel.ImUser
		serializ.BinaryRead(data, &imUser)
		return imUser
	}

	imUser := imUserModel.GetUserById(userId)
	if &imUser != nil && imUser.Id > 0 {
		this.CacheUser(&imUser)
	}
	return imUser
}

func (userService *UserServiceImpl) CacheUser(user *imUserModel.ImUser) {
	redis.SetkeyExPrire(paramter.GetImUserKey(user.Id), 60*60*24)
}

func (this *UserServiceImpl) DelUser(userId, destId int64) {
	//开始事务
	tx := model_mysql.Db.Begin()
	//删除好友
	if err := model_mysql.Execer("delete from im_friend where userId=? and	friendId=?", userId, destId); err != nil {
		tx.Rollback()
		return
	}
	if err := model_mysql.Execer("delete from im_friend where friendId=? and userId=?", userId, destId); err != nil {
		tx.Rollback()
		return
	}
	//删除置顶
	if err := model_mysql.Execer("delete from im_top where userId=? and destId=? and destType=1", userId, destId); err != nil {
		tx.Rollback()
		return
	}
	if err := model_mysql.Execer("delete from im_top where destId=? and userId=? and destType=1", userId, destId); err != nil {
		tx.Rollback()
		return
	}
	//删除朋友圈
	new(FeedServiceImpl).DeleteFriendFeed(userId, destId)

	//清除好友缓存
	key := []byte("friends_" + parse.Int64ToString(userId))
	redis.Delete(key)
	key = []byte("fiends_" + parse.Int64ToString(destId))
	redis.Delete(key)

	msg := new(MessageServiceImpl).DelFriendNotice(userId, destId)
	new(MsgStoreAndNotifyServiceImpl).StoreMsgAndNotifyImServer(msg, destId)
	tx.Commit()
}

func (this *UserServiceImpl) UpdateProfile(userId, nickName, city, province, name, sign, sex, district string) imUserModel.ImUser {

	user := this.GetImUser(parse.StringToInt64(userId))

	if !reflect.DeepEqual(user, imUserModel.ImUser{}) {
		if nickName != "" {
			user.NickName = nickName
			user.Name = nickName

			new(FeedServiceImpl).UpdateFeed(parse.StringToInt64(userId), nickName, "")
		}
		if city != "" {
			user.City = city
		}
		if province != "" {
			user.Province = province
		}
		if name != "" {
			user.Name = name
			user.NickName = name

			new(FeedServiceImpl).UpdateFeed(parse.StringToInt64(userId), name, "")

		}
		if sign != "" {
			user.Sign = sign
		}
		if sex != "" {
			user.Sex = sex
		}
		if district != "" {
			user.District = district
		}
		imUserModel.UpdateUser(user)
		this.CacheUser(&user)
	}
	//通知好友
	friends := new(ImFriendServiceImpl).GetFriends(user.Id)
	for _, v := range friends {
		new(ImFriendServiceImpl).CacheFriendInfo(v.FriendId)
		msg := new(MessageServiceImpl).UserModifyProfileNotice(user.Id, v.FriendId, &user)
		new(MsgStoreAndNotifyServiceImpl).StoreMsgAndNotifyImServer(msg, parse.StringToInt64(userId))
	}
	return user
}

//保存摇一摇历史记录
func (this *UserServiceImpl) SaveRockRecord(userId, destId int64, distance float64) {
	imUser := this.GetImUser(userId)
	if &imUser != nil {
		record := imRockRecordModel.ImRockRecord{
			UserId:   userId,
			DestId:   destId,
			Name:     imUser.Name,
			HeadUrl:  imUser.HeadUrl,
			Sex:      imUser.Sex,
			Sign:     imUser.Sign,
			Isdelete: 0,
			Distance: parse.Float64ToInt64(distance),
		}
		imRockRecordModel.Save(&record)
	}
}

func (this *UserServiceImpl) GetUserBlackList(userId int64) []interface{} {
	list := model_mysql.QueryColumn("select friendId from ImFriend where userId=? and isFriend=1 and isBlack=1", userId)
	fList := model_mysql.QueryColumn("select userId from ImFriend where friendId=? and isFriend=1 and isBlack=1", userId)
	feedlist := model_mysql.QueryColumn("select userId from im_friend where friendId=? and isFriend=1 and ufeedPriv=0", userId)
	ufeedlist := model_mysql.QueryColumn("select friendId from im_friend where userId=? and isFriend=1 and feedPriv=0", userId)
	var datas []interface{}
	datas = append(datas, list...)
	datas = append(datas, fList...)
	datas = append(datas, feedlist...)
	datas = append(datas, ufeedlist...)
	return slice.SliceRemoveDuplicate(datas)
}

func (this *UserServiceImpl) CheckIdNo(userId int64, idNo string) int {
	user := model_mysql.Query("select * from im_user where idNo=? limit 1", idNo)
	if user != nil {
		if parse.StringToInt64(user["id"].(string)) == userId {
			return 2 //已经被使用
		} else {
			return parse.StringToInt(user["isAuth"].(string))
		}
	}
	return -1
}

func (this *UserServiceImpl) UpdateAuth(userId int64, realname, validateNum, idNo, mobile string) int {
	row := 0
	user := this.GetImUser(userId)
	if &user != nil && sms.EqualValidate(mobile, validateNum) {
		user.IdNo = idNo
		user.IsAuth = 1
		user.RealName = realname
		imUserModel.UpdateUserNull(user)
		row = 1
	}
	return row
}

func (this *UserServiceImpl) RemoveBlock(userId, destId int64) {
	model_mysql.Exec("update im_friend set isBlack=0 where userId=? and friendId=?", userId, destId)
	new(ImFriendServiceImpl).CacheFriendInfo(userId)
}

func (this *UserServiceImpl) MyContact(userId int64) []map[string]interface{} {

	list := model_mysql.QueryList("SELECT uc.userId AS userId,uc.phone AS phone,uc.header AS header,uc.markName AS markName,IF(IFNULL(u.id,0)=0,0,1) AS isExists,u.id AS fuserId,u.nickName from `im_user_contact` uc LEFT JOIN im_user u on `u`.`mobile` like concat('%',`uc`.`phone`,'%') where uc.phone <>'' and uc.phone is not NULL and uc.userId=?", userId)

	return list
}

//导入手机通讯录
func (this *UserServiceImpl) ImportContact(userId int64, contact string) int {
	contact = strings.ReplaceAll(contact, "\\n", "")
	contact = strings.ReplaceAll(contact, "\\\"", "\"")
	var userContact []imUserContactModel.ImUserContact
	if err := json.Unmarshal([]byte(contact), &userContact); err != nil {
		logger.Error("UserService ImPortContact parse json array failed", err)
		return 0
	}
	model_mysql.Exec("delete from im_user_contact where userId=?", userId)
	new(imUserContactModel.ImUserContact).Inserts(userContact)
	return len(userContact)
}
