package mongo

import (
	"maoguo/henan/misc/config"
	"time"

	"github.com/wonderivan/logger"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var Session *mgo.Session
var Database *mgo.Database

func Init_Mongo() {

	var err error

	mgo.SetDebug(true)

	dail_info := &mgo.DialInfo{
		Addrs:    []string{"127.0.0.1"},
		Direct:   false,
		Timeout:  time.Second * 1,
		Database: "admin",
		Source:   "admin",
		Username: "roam",
		Password: config.CONFIG["MongoPassword"],
		// Password:  "123456",
		PoolLimit: 4096,
	}
	Session, err = mgo.DialWithInfo(dail_info)

	mgo.SetDebug(true)

	if err != nil {
		logger.Error("mgo dail error[%s]\n", err.Error())
		// err_handler(err)
	}
	Session.SetMode(mgo.Monotonic, true)
	Database = Session.DB("admin")

	logger.Info("mongo start success")
}

// 设置DEBUG模式  mgo.SetLogger(new(MongoLog)) // 设置日志.

// func err_handler(err error) {
// 	logger.Error("err_handler, error:%s\n", err.Error())
// }

func UpdateData(k, v string) bson.M {
	return bson.M{k: v}
}
