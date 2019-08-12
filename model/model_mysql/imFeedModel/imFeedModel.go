package imFeedModel

import (
	. "maoguo/henan/model/model_mysql"
)

type ImFeed struct {
	Id           int64  `json:"id"`
	FeedText     string `json:"feedText"`
	FeedImgs     string `json:"feedImgs"`
	FeedVideo    string `json:"feedVideo"`
	UserId       int64  `json:"userId"`
	UserName     string `json:"userName" gorm:"default:null"`
	UserHeadUrl  string `json:"UserHeadUrl" gorm:"default:null"`
	CreateTime   int64  `json:createTime`
	Status       int    `json:"status"`
	Lat          string `json:"lat"`
	Lng          string `json:"lng"`
	Address      string `json:"address"`
	Priv         int    `json:priv`
	At           string `json:"at"`
	Uids         string `json:uids`
	Ext          string `json:"ext"`
	BelongUserId int64  `gorm:"-" json:"belongUserId"`
}

func (ImFeed) TableName() string {
	return "im_feed"
}

func Save(feed *ImFeed) {
	Db.Create(&feed)
}

func QueryById(id int64) (feed ImFeed) {
	Db.Table("im_feed").Where("id= ?", id).Scan(&feed)
	return
}

func Raws(sql string, args ...interface{}) (feed []ImFeed) {
	Db.Raw(sql, args...).Scan(&feed)
	return
}

func Raw(sql string, args ...interface{}) (feed ImFeed) {
	Db.Raw(sql, args...).Scan(&feed)
	return
}

// func QueryPageFeed() page.Page {
//
//
// }
