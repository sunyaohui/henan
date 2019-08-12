package imGroupModel

import (
	. "maoguo/henan/model/model_mysql"
	"time"
)

type ImGroup struct {
	Id             int64      `grom:"column:id" json:"id"`
	Name           string     `grom:"column:name;default:null" json:"name"`     //名称
	Descriptions   string     `grom:"column:descriptions" json:"descriptions"`  //描述
	Detail         string     `grom:"column:detail;default:null" json:"detail"` //  备注
	HeadUrl        string     `gorm:"column:headUrl" json:"headUrl"`            //头像class或者url
	CreaterId      int64      `gorm:"column:createrId" json:"createrId"`
	CreateTime     int64      `gorm:"column:createTime" json:"createTime"`
	JoinStatus     int        `gorm:"column:joinStatus" json:"joinStatus"` //加群 验证 0-无需验证 1 需要验证消息  2 需要回答问题并由管理员审核 3 需要正确回答问题 4 只允许群成员邀请 5 不允许任何人加入 6 付费加群
	Question       string     `gorm:"column:question" json:"question"`
	Answer         string     `gorm:"column:answer" json:"answer"`
	Fee            float64    `gorm:"column:fee" json:"fee"`
	Level          int16      `gorm:"column:level" json:"level"`
	Expire         int64      `gorm:"column:expire" json:"expire"`                   //到期时间
	LookMemberInfo int        `gorm:"column:look_member_info" json:"lookMemberInfo"` //查看群成员信息
	MemberCount    int        `gorm:"column:memberCount" json:"memberCount"`
	MaxMemberCount int        `gorm:"column:maxMemberCount" json:"maxMemberCount"`
	Useenddate     *time.Time `gorm:"column:useenddate" json:"useenddate"`
	Stopstatus     int        `gorm:"column:stopstatus" json:"stopstatus"`
	SpeakStatus    int        `gorm:"column:speakStatus" json:"speakStatus"` //0 标识未禁言 -1标识禁言
}

func (ImGroup) TableName() string {
	return "im_group"
}

func GetImGroupById(id int64) *ImGroup {
	var group ImGroup
	Db.Table("im_group").Where("id = ?", id).Scan(&group)
	return &group
}

func Save(this *ImGroup) {
	Db.Create(this)
}

func Raw(sql string, args ...interface{}) (group ImGroup) {
	Db.Raw(sql, args...).Scan(&group)
	return
}

func Update(user ImGroup) {
	Db.Debug().Model(&user).Updates(user)
}
