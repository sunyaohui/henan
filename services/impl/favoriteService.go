package impl

import (
	"maoguo/henan/model/model_mysql/imFavoriteModel"
	"time"
)

type FavoriteServiceImpl struct{}

func (this *FavoriteServiceImpl) AddFavorite(userId, fromId int64, category int, content, fromName, fromHeadUrl string) bool {
	params := map[string]interface{}{"fromId": fromId, "content": content, "category": category, "userId": userId}
	dbFavorite := imFavoriteModel.GetImFavorite(params)
	if &dbFavorite != nil && dbFavorite.Id > 0 {
		return false
	}
	favorite := imFavoriteModel.ImFavorite{
		Category:    category,
		Content:     content,
		CreateTime:  time.Now().Unix(),
		FromHeadUrl: fromHeadUrl,
		FromName:    fromName,
		UserId:      userId,
		FromId:      fromId,
	}
	imFavoriteModel.Save(&favorite)
	return true
}
