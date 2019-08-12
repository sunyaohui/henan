package mongo

import (
	"fmt"

	"github.com/wonderivan/logger"
	"gopkg.in/mgo.v2/bson"
)

type ImFeed struct {
	Id             bson.ObjectId   `bson:"_id" json:"idNo"`
	FeedText       string          `bson:"feedText" json:"feedText"`
	FeedImgs       string          `bson:"feedImgs" json:"feedImgs"`
	FeedVideo      string          `bson:"feedVideo" json:"feedVideo"`
	FeedId         int64           `bson:"feedId" json:"feedId"`
	UserId         int64           `bson:"userId" json:"userId"`
	UserName       string          `bson:"userName" json:"userName"`
	UserHeadUrl    string          `bson:"userHeadUrl" json:"userHeadUrl"`
	CreateTime     int64           `bson:"createTime" json:"createTime"`
	Status         int             `bson:"status" json:"status"`
	Lat            string          `bson:"lat" json:"lat"`
	Lng            string          `bson:"lng" json:"lng"`
	Address        string          `bson:"address" json:"address"`
	Priv           int             `bson:"priv" json:"priv"`
	At             string          `bson:"at" json:"at"`
	Uids           string          `bson:"uids" json:"uids"`
	Ext            string          `bson:"ext" json:"ext"`
	BelongUserId   int64           `bson:"belongUserId" json:"belongUserId"`
	ImFeedPraises  []ImFeedPraise  `bson:"imFeedPraises" json:"imFeedPraises"`
	ImFeedComments []ImFeedComment `bson:"imFeedComments" json:"imFeedComments"`
	AtUsers        []ImFeedUser    `bson:"atUsers" json:"atUsers"`
}

func (feed *ImFeed) UpdateAll(where, data bson.M) bool {
	_, err := GetImFeed().UpdateAll(where, bson.M{"$set": data})
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func (feed *ImFeed) UpdatePush(where, data bson.M) bool {
	_, err := GetImFeed().UpdateAll(where, bson.M{"$push": data})
	if err != nil {
		logger.Error("mongo imFeed update push failed", err)
		return false
	}
	return true
}

func (feed *ImFeed) UpdatePull(where, data bson.M) bool {
	_, err := GetImFeed().UpdateAll(where, bson.M{"$pull": data})
	if err != nil {
		logger.Error("mongo imFeed update pull failed", err)
		return false
	}
	return true
}

func (feed *ImFeed) Remove(where bson.M) bool {
	err := GetImFeed().Remove(where)
	if err != nil {
		return false
	}
	return true
}

func (feed *ImFeed) Insert() {
	err := GetImFeed().Insert(&feed)
	if err != nil {
		logger.Error("mongo insert imfeed failed,", err)
	}
}

func (feed *ImFeed) FindOne(where bson.M) (f ImFeed) {
	err := GetImFeed().Find(where).One(&f)
	if err != nil {
		logger.Error("mongo get imFeed failed,params:%d", where, err)
		f = ImFeed{}
	}
	return f
}

func (feed *ImFeed) Find(where bson.M) (f []ImFeed) {
	err := GetImFeed().Find(where).All(&f)
	if err != nil {
		return nil
	}
	return f
}

func (feed *ImFeed) Count(where bson.M) int {
	var f []ImFeed
	err := GetImFeed().Find(where).All(&f)
	if err != nil {
		return 0
	}
	return len(f)
}
