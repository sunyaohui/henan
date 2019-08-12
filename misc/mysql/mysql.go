package mysql

import (
	"database/sql"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

//数据库配置
const (
	userName = "root"
	password = "7uyrhAQ3!41X"
	ip       = "42.236.74.134"
	port     = "3306"
	dbName   = "imdb"
)

//Db数据库连接池
var DB *sql.DB

//mysql
func InitDB() {
	path := strings.Join([]string{userName, ":", password, "@tcp(", ip, ":", port, ")/", dbName, "?charset=utf8"}, "")
	DB, _ = sql.Open("mysql", path)
	DB.SetConnMaxLifetime(100)
	DB.SetMaxIdleConns(10)
	//验证连接
	if err := DB.Ping(); err != nil {
		log.Fatalln("connect database failed", err)
		return
	}
	log.Println("connect database success")
}
