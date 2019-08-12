package impl

import (
	"maoguo/henan/model/model_mysql"
	"maoguo/henan/model/model_mysql/imFriendGroupModel"
	"time"
)

type CommonServiceImpl struct {
}

func (this *CommonServiceImpl) GetUserFriends(userId int64) map[string]interface{} {
	this.SetDefaultFriendGroup(userId)
	list := model_mysql.QueryList("select * from im_friend_group where userId=? order by sort desc limit 1", userId)
	result := make(map[string]interface{}, 0)
	if len(list) > 0 {
		for _, value := range list {
			result["gid"] = value["id"]
			result["gname"] = value["name"]
			friendsInfo := model_mysql.QueryList("select a.remark,a.isBlack,a.receiveTip,b.id,name,nickName,mobile,mail,sex,birthday,sign,province,isOnline,city,district,b.createTime,status,detail,headUrl,a.bgurl,b.imNumber as IMNo from im_friend a left join im_user b on(a.friendId=b.id) where a.userId=? and a.isFriend=1 and a.groupid=?", userId, value["id"])
			if friendsInfo == nil {
				friendsInfo = make([]map[string]interface{}, 0)
			}
			result["friends"] = friendsInfo
		}
	}
	return result
}

func (this *CommonServiceImpl) SetDefaultFriendGroup(userId int64) {
	list := model_mysql.QueryList("select * from im_friend_group where userId=? and isdefault=1", userId)
	if len(list) > 0 {
		friendGroup := &imFriendGroupModel.ImFriendGroup{
			UserId:     userId,
			Createtime: time.Now().Unix(),
			Name:       "我的好友",
			Sort:       0,
			Isdefault:  0,
		}
		friendGroup.Create()
		model_mysql.Exec("update im_friend set groupid=? where userId=?;", friendGroup.Id, userId)
	} else {
		group := imFriendGroupModel.QueryImFriendGroupByUserId(userId)
		model_mysql.Exec("update im_friend set groupid=? where userId=? and groupid=0", group.Id, userId)
	}
}
