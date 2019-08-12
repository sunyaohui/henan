package mongo

import (
	"github.com/wonderivan/logger"
)

type ImFeedComment struct {
	Id               int64  `bson:"_id" json:"fcid"`
	UserId           int64  `bson:"userId" json:"userId"`
	UserName         string `bson:"userName" json:"userName"`
	UserHeadUrl      string `bson:"userHeaderUrl" json:"userHeaderUrl"`
	ReplyUserId      int64  `bson:"replyUserId" json:"replyUserId"`
	ReplyUserName    string `bson:"replyUserName" json:"replyUserName"`
	ReplyUserHeadUrl string `bson:"replyUserHeadUrl" json:"replyUserHeadurl"`
	CommentText      string `bson:"commentText" json:"commentText"`
	CreateTime       int64  `bson:"createTime" json:"createTime"`
}

func (feed *ImFeedComment) Insert() {
	err := GetImFeedComment().Insert(feed)
	if err != nil {
		logger.Error("mongo ImFeedComment insert failed", err)
	}
}
