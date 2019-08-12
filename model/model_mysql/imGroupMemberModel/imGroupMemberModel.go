package imGroupMemberModel

import (
	. "maoguo/henan/model/model_mysql"
)

type ImGroupMember struct {
	Id         int64  `json:"id"`
	GroupId    int64  `gorm:"column:groupId" json:"groupId"` //群组id
	UserId     int64  `gorm:"column:UserId" json:"userId"`   //成员id
	MarkName   string `gorm:"column:markName;default:null" json:"markName"`
	Role       int    `gorm:"column:role" json:"role"`
	CreatorId  int64  `gorm:"column:creatorId" json:"creatorId"`
	CreateTime int64  `gorm:"column:createTime" json:"createTime"`
	ReceiveTip int    `gorm:"column:receiveTip" json:"receiveTip"`
	IsAccept   int    `gorm:"column:isAccept" json:"isAccept"`
	Bgurl      string `json:"bgurl"`
	Silence    int64  `json:"silence"`
}

func (ImGroupMember) TableName() string {
	return "im_group_member"
}

func Save(this *ImGroupMember) {
	Db.Create(&this)
}

func Raw(sql string, args ...interface{}) (groupMember ImGroupMember) {
	Db.Raw(sql, args...).Scan(&groupMember)
	return
}
