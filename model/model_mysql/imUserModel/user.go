package imUserModel

import (
	. "maoguo/henan/model/model_mysql"

	_ "github.com/jinzhu/gorm"
)

type ImUser struct {
	Id              int64  `gorm:"column:id" json:"id"`
	QxNo            string `gorm:"column:qxNo;default:null" json:"qxNo"`
	AuthNo          int    `gorm:"column:autoNo" json:"authNo"`
	ImNumber        string `gorm:"column:imNumber;default:null" json:"imNumber"`
	Name            string `gorm:"column:name;default:null" json:"name"`
	IdNo            string `gorm:"column:idNo;default:null" json:"idNo"`
	NickName        string `gorm:"column:nickName;default:null" json:"nickName"`
	RealName        string `gorm:"column:realName;default:null" json:"realName"`
	Pwd             string `gorm:"column:pwd;default:null" json:"pwd"`
	Mobile          string `gorm:"column:mobile;default:null" json:"mobile"`
	Mail            string `gorm:"column:mail;default:null" json:"mail"`
	Sex             string `gorm:"column:sex;default:男" json:"sex"`
	Birthday        int64  `gorm:"column:birthday;default:null" json:"birthday"`
	Sign            string `gorm:"column:sign;default:null" json:"sign"`
	Province        string `gorm:"column:province;default:null" json:"province"`
	IsOnline        int    `gorm:"column:isOnline" json:"isOnline"`
	NeedAuth        int    `gorm:"column:needAuth" json:"needAuth"`
	SearchMobile    int    `gorm:"column:searchMobile" json:"searchMobile"`
	NewNotification int    `gorm:"column:newNotification" json:"newNotification"`
	City            string `gorm:"column:city;default:null" json:"city"`
	District        string `gorm:"column:district;default:null" json:"district"`
	CreateTime      int64  `gorm:"column:createTime" json:"createTime"`
	Status          int    `gorm:"column:status" json:"status"`
	Detail          string `gorm:"column:detail;default:null" json:"detail"`
	IsAuth          int    `gorm:"column:isAuth" json:"isAuth"`
	HeadUrl         string `gorm:"column:headUrl" json:"headUrl"`
	RecommandUserId int64  `gorm:"column:recommandUserId" json:"recommandUserId"`
	Longitude       string `gorm:"column:longitude;default:null" json:"longitude"`
	FeedBackImage   string `gorm:"column:feedBackImage;default:null" json:"feedBackImage"`
	Latitude        string `gorm:"column:latitude;default:null" json:"latitude"`
	Bgurl           string `gorm:"column:bgurl;default:null" json:"bgurl"`
	Isdelete        int8   `gorm:"column:isdelete" json:"isdelete"`
	UserGroupId     int    `gorm:"column:user_group_id" json:"userGroupId"`
}

// func (ImUser) TableName() string {
// 	return "im_user"
// }

func GetUserList(params map[string]interface{}) (user []ImUser) {
	Db.Table("im_user").Where(params).Scan(user)
	return
}

func GetUserById(userId int64) (user ImUser) {
	Db.Table("im_user").Where("id = ?", userId).Scan(&user)
	return
}

func GetUser(params map[string]interface{}) (user ImUser) {
	Db.Table("im_user").Where(params).Scan(&user)
	return
}

func UpdateUserNull(user ImUser) {
	Db.Debug().Model(&user).Save(user)
}

func UpdateUser(user ImUser) {
	Db.Debug().Model(&user).Save(user)
}

func Save(user *ImUser) {
	Db.Create(&user)
}

func Raw(sql string, args ...interface{}) (user ImUser) {
	Db.Raw(sql, args...).Scan(&user)
	return
}
