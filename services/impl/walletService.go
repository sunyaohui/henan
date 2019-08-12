package impl

import "maoguo/henan/model/model_mysql/imWalletModel"

//import "maoguo/henan/model/imWalletModel"

type ImWalletServiceImpl struct {
}

func (walletService *ImWalletServiceImpl) InitWallet(userId int64) {
	wallet := imWalletModel.QueryImWallet(map[string]interface{}{"userId": userId})
	if wallet != nil {
		var userWallet *imWalletModel.ImWallet
		userWallet.UserId = userId
		userWallet.Money = 1000
		userWallet.Status = 1
		userWallet.Create()
	}
}
