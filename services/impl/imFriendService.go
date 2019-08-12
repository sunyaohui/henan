package impl

import (
	"maoguo/henan/misc/redis"
	"maoguo/henan/misc/utils/parse"
	"maoguo/henan/misc/utils/serializ"
	"maoguo/henan/model/model_mysql"
	"maoguo/henan/model/model_mysql/imFriendModel"
	"maoguo/henan/model/model_mysql/imUserModel"
	"time"
)

type ImFriendServiceImpl struct {
}

func (this *ImFriendServiceImpl) GetFriends(userId int64) []imFriendModel.ImFriend {
	friends := imFriendModel.QueryImFriendByUserId(userId)
	for _, v := range friends {
		key := []byte("friends_" + parse.Int64ToString(userId))
		redis.Sadd(key, serializ.BinaryWrite(v))
	}
	return friends
}

func (this *ImFriendServiceImpl) CacheFriendInfo(userId int64) []map[string]interface{} {
	result := model_mysql.QueryList("select a.remark,a.isBlack,a.receiveTip,b.id,name,nickName,mobile,mail,sex,birthday,sign,province,isOnline,city,district,b.createTime,status,detail,headUrl,a.bgurl,b.imNumber as IMNo from im_friend a left join im_user b on(a.friendId=b.id) where a.userId=? and a.isFriend=1", userId)
	if result == nil {
		return make([]map[string]interface{}, 0)
	}
	return result
}

func (this *ImFriendServiceImpl) GetFriendInfo(userId int64) []map[string]interface{} {
	return this.CacheFriendInfo(userId)
}

// func (this *ImFriendServiceImpl) cacheFriendInfo(userId string) []map[string]interface{} {
// 	return model_mysql.QueryList("select a.remark,a.isBlack,a.receiveTip,b.id,name,nickName,mobile,mail,sex,birthday,sign,province,isOnline,city,district,b.createTime,status,detail,headUrl,a.bgurl,b.imNumber as IMNo from im_friend a left join im_user b on(a.friendId=b.id) where a.isFriend=1 and a.userId=?", userId)
// }

func (this *ImFriendServiceImpl) GetMixFriend(userId, friendId, replyUserId int64) []string {
	userkey := []byte("friends_" + parse.Int64ToString(userId))
	friendKey := []byte("friends_" + parse.Int64ToString(friendId))
	replyKey := []byte("friends_" + parse.Int64ToString(replyUserId))
	if replyUserId > 0 {
		this.GetFriends(replyUserId)
		return redis.Sinter(userkey, friendKey, replyKey)
	} else {
		return redis.Sinter(userkey, friendKey)
	}
}

func (this *ImFriendServiceImpl) RequestFriend(destType int, userId, destId int64) int {
	if destType == 1 {
		destUser := imUserModel.GetUserById(destId)
		if &destUser == nil {
			return -1
		}
		if destUser.NeedAuth == 0 {
			redis.Sadd([]byte("friends_"+parse.Int64ToString(userId)), []byte(parse.Int64ToString(destId)))
			redis.Sadd([]byte("friends_"+parse.Int64ToString(destId)), []byte(parse.Int64ToString(userId)))
			//对方不需要进行认证时，直接添加为好友
			this.AddFriend(userId, destId, true)
			msg := new(MessageServiceImpl).FriendNotice(userId, destId)
			new(MsgStoreAndNotifyServiceImpl).StoreMsgAndNotifyImServer(&msg, destId)
			return 2
		} else {
			this.AddFriend(userId, destId, false)
			msg := new(MessageServiceImpl).FriendNotice(userId, destId)
			new(MsgStoreAndNotifyServiceImpl).StoreMsgAndNotifyImServer(&msg, destId)
			//已发送，等待对方同意
			return 1
		}
	}
	return -1
}

func (this *ImFriendServiceImpl) AddFriend(userId, destId int64, isDouble bool) {
	if isDouble {
		friend := imFriendModel.ImFriend{
			CreaterId:  userId,
			CreateTime: time.Now().Unix(),
			FriendId:   destId,
			IsBlack:    0,
			IsFriend:   1,
			ReceiveTip: 1,
			UserId:     userId,
		}
		model_mysql.Exec("delete from im_friend where userId=? and friendId=?", userId, destId)
		imFriendModel.Save(friend)

		friend2 := imFriendModel.ImFriend{
			CreaterId:  userId,
			CreateTime: time.Now().Unix(),
			FriendId:   userId,
			IsBlack:    0,
			IsFriend:   1,
			ReceiveTip: 1,
			UserId:     destId,
		}
		model_mysql.Exec("delete from im_friend where userId=? and friendId=?", destId, userId)
		imFriendModel.Save(friend2)
	} else {
		friend := imFriendModel.ImFriend{
			CreaterId:  userId,
			CreateTime: time.Now().Unix(),
			FriendId:   destId,
			IsBlack:    0,
			IsFriend:   0,
			ReceiveTip: 1,
			UserId:     userId,
		}
		model_mysql.Exec("delete from im_friend where userId=? and friendId=?", userId, destId)
		imFriendModel.Save(friend)

		friend2 := imFriendModel.ImFriend{
			CreaterId:  userId,
			CreateTime: time.Now().Unix(),
			FriendId:   userId,
			IsBlack:    0,
			IsFriend:   0,
			ReceiveTip: 1,
			UserId:     destId,
		}
		model_mysql.Exec("delete from im_friend where userId=? and friendId=?", destId, userId)
		imFriendModel.Save(friend2)
	}
	this.CacheFriendInfo(userId)
	this.CacheFriendInfo(destId)
}
