package impl

import (
	"maoguo/henan/model/model_mysql/imTopModel"
)

type ImTopServiceImpl struct{}

func (this *ImTopServiceImpl) GetTopList(userId int64) []imTopModel.ImTop {
	return this.CacheTopList(userId)
}

func (this *ImTopServiceImpl) CacheTopList(userId int64) []imTopModel.ImTop {
	return imTopModel.QueryImTopByUserId(userId)
}
