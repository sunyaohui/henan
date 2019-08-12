package imUserContactModel

import (
	. "maoguo/henan/model/model_mysql"
	"time"
)

type ImUserContact struct {
	Id         int64      `gorm:"column:id" json:"id"`
	UserId     int64      `gorm:"column:userId;default:null" json:"userId"`
	Phone      string     `gorm:"column:phone;default:null" json:"phone"`
	Header     string     `gorm:"column:header;default:null" json:"header"`
	MarkName   string     `gorm:"column:markName;default:null" json:"markName"`
	CreateTime *time.Time `gorm:"column:createTime" json:"createTime"`
	UpdateTime *time.Time `gorm:"column:updateTime" json:"updateTime"`
}

func (this ImUserContact) Inserts(contacts []ImUserContact) {
	for _, constact := range contacts {
		Db.Create(&constact)
	}
}
