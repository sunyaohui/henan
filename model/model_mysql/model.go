package model_mysql

import (
	"fmt"
	"maoguo/henan/misc/config"

	"github.com/jinzhu/gorm"
	"github.com/wonderivan/logger"
)

var Db *gorm.DB

func Init_db() {
	var (
		err                                  error
		dbType, dbName, user, password, host string
	)
	dbType = "mysql"
	dbName = "imdb"
	user = "root"
	password = config.CONFIG["MysqlPwd"]
	// host = "42.236.74.134"
	// password = "123456"
	host = "127.0.0.1"

	Db, err = gorm.Open(dbType, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		user,
		password,
		host,
		dbName))

	if err != nil {
		logger.Error("mysql start failed,", err)
	}

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return defaultTableName
	}

	//禁用表名复数
	Db.SingularTable(true)
	Db.DB().SetMaxIdleConns(10)
	Db.DB().SetMaxOpenConns(100)

	logger.Info("mysql start success")
}

func CloseDB() {
	defer Db.Close()
}
