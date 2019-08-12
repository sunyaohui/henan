package services

import (
	"maoguo/henan/model/model_mysql/imFriendModel"
	"maoguo/henan/model/model_mysql/imMessageModel"
	"maoguo/henan/model/model_mysql/imUserModel"
)

type MessageService interface {
	ModifyGroupNoticeCommon(fromId, destId int64, content string, msgType int) *imMessageModel.ImMessage
	OtherLoginNotice(fromId int64, content string) *imMessageModel.ImMessage
	NoticeOffline(id string)
}

type UserService interface {
	GetImUser(userId int64) *imUserModel.ImUser
	CacheUser(user *imUserModel.ImUser)
}

type GroupService interface {
	CheckUserGroup(userId int64)
}

type ImWalletService interface {
	InitWallet(userId int64)
}

type CommonService interface {
	GetUserFriends(userId string) map[string]interface{}
	SetDefaultFriendGroup(userId int)
}

type MsgStoreAndNotify interface {
	StoreMsgAndNotifyImServer(message *imMessageModel.ImMessage, id string)
	GetImServiceUrl() string
}

type ImFriendService interface {
	GetFriends(userId int64) *[]imFriendModel.ImFriend
}

type FeedService interface {
	UpdateFeed(userId int64, userName, headUrl string)
}
