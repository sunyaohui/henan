package mongo

import (
	"github.com/wonderivan/logger"
)

type ImFeedPraise struct {
	Id          int64  `bson:"_id" json:"fpid"`
	UserId      int64  `bson:"userId" json:"userId"`
	UserName    string `bson:"userName" json:"userName"`
	UserHeadUrl string `bson:"userHeadUrl" json:"userHeadUrl"`
	CreateTime  int64  `bson:"createTime" json:"createTime"`
}

func (feed *ImFeedPraise) Insert() {
	err := GetImFeedPraise().Insert(feed)
	if err != nil {
		logger.Info("im feed praise insert failed", err)
	}
}
