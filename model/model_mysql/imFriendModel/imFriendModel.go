package imFriendModel

import (
	. "maoguo/henan/model/model_mysql"
)

type ImFriend struct {
	Id         int64  `gorm:"column:id" json:"id"`
	UserId     int64  `gorm:"column:userId" json:"userId"`     //用户id
	FriendId   int64  `gorm:"column:friendId" json:"friendId"` //好友id
	Remark     string `json:"remark" gorm:"default:null"`      //备注名称
	CreaterId  int64  `gorm:"column:createrId" json:"createrId"`
	CreateTime int64  `gorm:"column:createTime" json:"createTime"`         //创建时间
	IsBlack    int    `gorm:"column:isBlack" json:"isBlack"`               //0不是黑名单1黑名单
	IsFriend   int    `gorm:"column:isFriend;default:0" json:"isFriend"`   //0单向1双向
	ReceiveTip int    `gorm:"column:receiveTip" json:"receiveTip"`         //1接收提示0不接收提示
	Bgurl      string `json:"bgurl"`                                       //好友聊天背景
	FeedPriv   int    `gorm:"column:feedPriv;default:1" json:"feedPriv"`   //好友的朋友圈 1看 0 不看
	UfeedPriv  int    `gorm:"column:ufeedPriv;default:1" json:"ufeedPriv"` //我的朋友圈对好友权限 0 不可看 1可看
	Groupid    int64  `json:"groupid"`                                     //分组id
}

func (ImFriend) Table() string {
	return "im_friend"
}

func QueryImFriendByUserId(userId int64) []ImFriend {
	var imFriends []ImFriend
	Db.Table("im_friend").Where("isFriend=1 and userId=?", userId).Scan(&imFriends)
	return imFriends
}

func Save(friend ImFriend) {
	Db.Create(&friend)
}
