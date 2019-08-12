package mongo

import (
	mgo "gopkg.in/mgo.v2"
)

func GetImFeed() *mgo.Collection {
	return Database.C("ImFeed")
}

func GetImFeedComment() *mgo.Collection {
	return Database.C("ImFeedComment")
}

func GetImFeedPraise() *mgo.Collection {
	return Database.C("ImFeedPraise")
}

func GetImFeedUser() *mgo.Collection {
	return Database.C("ImFeedUser")
}
