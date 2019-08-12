package imTopModel

import (
	. "maoguo/henan/model/model_mysql"
)

type ImTop struct {
	Id       int64 `json:"id"`
	UserId   int64 `gorm:"column:userId" json:"userId"`
	DestType int   `gorm:"column:destType" json:"destType"`
	DestId   int64 `gorm:"column:destId" json:"destId"`
}

func (ImTop) TableName() string {
	return "im_top"
}

func QueryImTopByUserId(userId int64) (top []ImTop) {
	Db.Table("im_top").Where("userId=?", userId).Scan(&top)
	return top
}

func Save(top *ImTop) {
	Db.Create(&top)
}
