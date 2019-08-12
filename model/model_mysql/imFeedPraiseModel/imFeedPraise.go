package imFeedPraiseModel

import "maoguo/henan/model/model_mysql/imFeedModel"
import . "maoguo/henan/model/model_mysql"

type ImFeedPraise struct {
	Id          int64              `json:"id"`
	FeedId      int64              `json:"feedId"`
	UserId      int64              `json:"userId"`
	UserName    string             `json:"userName" gorm:"default:null"`
	UserHeadUrl string             `json:"userHeadUrl" gorm:"default:null"`
	CreateTime  int64              `json:"createTime"`
	ImFeed      imFeedModel.ImFeed `gorm:"-"`
}

func (ImFeedPraise) TableName() string {
	return "im_feed_praise"
}

func Create(fp *ImFeedPraise) {
	Db.Create(&fp)
}
