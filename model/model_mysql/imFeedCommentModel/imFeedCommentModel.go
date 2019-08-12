package imFeedCommentModel

import (
	. "maoguo/henan/model/model_mysql"
	"maoguo/henan/model/model_mysql/imFeedModel"
)

type ImFeedComment struct {
	Id               int64               `json:"id"`
	FeedId           int64               `json:"feedId"`
	UserId           int64               `json:"userId"`
	UserName         string              `json:"userName" gorm:"default:null"`
	UserHeadUrl      string              `json:"userHeadUrl" gorm:"default:null"`
	ReplyUserId      int64               `json:"replyUserId"`
	ReplyUserName    string              `json:"replyUserName" gorm:"default:null"`
	ReplyUserHeadUrl string              `json:"replyUserHeadUrl" gorm:"deafault:null"`
	CommentText      string              `json:"commentText"`
	CreateTime       int64               `json:"createTime"`
	ImFeed           *imFeedModel.ImFeed `gorm:"-"`
}

func (ImFeedComment) TableName() string {
	return "im_feed_comment"
}

func (feed *ImFeedComment) Create() {
	Db.Create(&feed)
}
