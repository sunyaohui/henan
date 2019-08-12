package imFavoriteModel

import (
	"maoguo/henan/misc/page"
	. "maoguo/henan/model/model_mysql"

	_ "github.com/jinzhu/gorm"
)

type ImFavorite struct {
	Id          int64  `json:"id"`
	Category    int    `json:"category"`
	FromId      int64  `json:"fromId"`
	FromName    string `json:"fromName" gorm:"column:from_name;default:null"`
	FromHeadUrl string `json:"fromHeadUrl" gorm:"column:from_head_url;default:null"`
	CreateTime  int64  `json:"createTime" gorm:"column:create_time"`
	Content     string `json:"content"`
	UserId      int64  `json:"userId"`
}

func (ImFavorite) TableName() string {
	return "im_favorite"
}

func GetImFavorite(params map[string]interface{}) (favorite ImFavorite) {
	Db.Table("im_favorite").Where(params).Scan(&favorite)
	return
}

func Save(favorite *ImFavorite) {
	Db.Create(favorite)
}

func GetFavoritePage(params map[string]interface{}, pg *page.Page) {
	var count int
	var favorites []ImFavorite
	var totalCount int
	Db.Table("im_favorite").Where(params).Count(&count)
	Db.Table("im_favorite").Where(params).Offset((pg.PageNo - 1) * pg.PageSize).Limit(pg.PageSize).Order("create_time desc").Find(&favorites).Count(&totalCount)
	pg.List = favorites
	if favorites == nil || len(favorites) == 0 {
		pg.List = make([]ImFavorite, 0)
	}
	pg.TotalCount = count
}
