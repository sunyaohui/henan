package impl

import "maoguo/henan/model/model_mysql/imUserPayModel"

type ImUserPayServiceImpl struct {
}

func (this *ImUserPayServiceImpl) GetImUserPay(userId int64) *imUserPayModel.ImUserPay {
	userPay := this.cacheImUserPay(userId)
	return userPay
}

func (this *ImUserPayServiceImpl) cacheImUserPay(userId int64) *imUserPayModel.ImUserPay {
	userPay := imUserPayModel.QueryImUserPayByUserId(userId)
	return userPay
}
