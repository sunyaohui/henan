package mongo

import "fmt"

type ImFeedUser struct {
	UserId      int64  `bson:"userId" json:"userName"`
	UserName    string `bson:"userName" json:"userName"`
	UserHeadUrl string `bson:"userHeadUrl" json:"userHeadUrl"`
	FeedId      int64  `bson:"feedId" json:"feedId"`
}

func (fu *ImFeedUser) Save() {
	err := GetImFeedUser().Insert(&fu)
	if err != nil {
		fmt.Println("插入失败")
	}
}
