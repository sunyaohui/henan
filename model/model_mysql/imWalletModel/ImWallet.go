package imWalletModel

import (
	. "maoguo/henan/model/model_mysql"
)

type ImWallet struct {
	Id     int64   `json:"id"`
	UserId int64   `json:"userId" gorm:"column:userId"`
	Money  float64 `json:"money"`
	Status int     `json:"status"`
}

func (ImWallet) TableName() string {
	return "im_wallet"
}

func QueryImWallet(maps map[string]interface{}) *ImWallet {
	var wallet ImWallet
	Db.Table("im_wallet").Where(maps).Scan(&wallet)
	return &wallet
}

func QueryImWalletList(maps map[string]interface{}) (imWallet []ImWallet) {
	Db.Table("im_wallet").Where(maps).Scan(&imWallet)
	return
}

func (wallet *ImWallet) Create() {
	Db.Create(&wallet)
}
