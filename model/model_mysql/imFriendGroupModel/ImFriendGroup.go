package imFriendGroupModel

import (
	. "maoguo/henan/model/model_mysql"
)

type ImFriendGroup struct {
	Id         int    `json:"id"`
	UserId     int64  `gorm:"column:userId" json:"userId"`
	Name       string `gorm:"column:name" json:"name"`
	Sort       int16  `gorm:"column:sort" json:"sort"`
	Createtime int64  `gorm:"column:createTime" json:"createTime"`
	//CreateTimestamp time.Time `grom:"column:createTimestamp;default Now()"`
	Isdefault int16 `gorm:"column:isdefault" json:"isdefault"`
}

//创建ImFriendGroup
func (group *ImFriendGroup) Create() {
	Db.Create(&group)
}

func QueryImFriendGroupByUserId(userId int64) ImFriendGroup {
	var group ImFriendGroup
	Db.Table("im_friend_group").Where("userId = ?", userId).Scan(&group)
	return group
}
