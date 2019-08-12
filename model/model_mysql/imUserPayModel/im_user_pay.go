package imUserPayModel

import (
	. "maoguo/henan/model/model_mysql"
)

type ImUserPay struct {
	Id         int64  `json:"id"`
	UserId     int64  `json:"userId" gorm:"column:userPwd;default:null"`
	PayPwd     string `json:"payPwd" gorm:"column:payPwd;default:null"`
	OldPwd     string `json:"oldPwd" gorm:"column:oldPwd;default:null"`
	CreateTime int64  `json:"createTime" gorm:"column:createTime"`
	UpdateTime int64  `json:"updateTime" gorm:"column:updateTime"`
}

func (ImUserPay) TableName() string {
	return "im_user_pay"
}

func QueryImUserPayByUserId(userId int64) *ImUserPay {
	var userPay ImUserPay
	Db.Table("im_user_pay").Where("userId=?", userId).Scan(&userPay)
	return &userPay
}
