package imRockRecordModel

import (
	. "maoguo/henan/model/model_mysql"
	"time"
)

type ImRockRecord struct {
	Id         int64     `json:"id"`
	UserId     int64     `json:"userId" gorm:"column:userId"` //用户id
	DestId     int64     `json:"destId" gorm:"column:destId"` //摇到的用户id
	Name       string    `json:"name"`
	HeadUrl    string    `json:"headUrl" gorm:"column:headUrl;default:null"` //摇一摇的用户头像连接
	Sex        string    `json:"sex"`
	Distance   int64     `json:"distance"`
	Sign       string    `json:"sign"`
	CreateTime time.Time `json:"createTime" gorm:"column:createTime"`
	UpdateTime time.Time `json:"updateTime" gorm:"column:updateTime"`
	Isdelete   int16     `json:"isdelete"` //删除标记  0- 未删除 1- 已经删除
}

func (ImRockRecord) TableName() string {
	return "im_rock_record"
}

func Save(r *ImRockRecord) {
	Db.Create(r)
}
