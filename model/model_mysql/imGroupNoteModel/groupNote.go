package imGroupNoteModel

import (
	. "maoguo/henan/model/model_mysql"
	"time"
)

type ImGroupNote struct {
	Id         int64      `gorm:column:"id" json:"id"`
	Content    string     `gorm:"column:content;default:null" json:"content"`
	GroupId    int64      `gorm:"column:groupId" json:"groupId"`
	UserId     int64      `gorm:"column:userId" json:"userId"`
	CreateTime *time.Time `gorm:"column:createTime" json:"createTime"`
	UpdateTime *time.Time `gorm:"column:updateTime" json:"updateTime"`
	Title      string     `gorm:"column:title" json:"title"`
}

func (ImGroupNote) TableName() string {
	return "im_group_name"
}

func Raw(sql string, args ...interface{}) (note ImGroupNote) {
	Db.Raw(sql, args...).Scan(&note)
	return
}
