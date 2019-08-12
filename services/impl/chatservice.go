package impl

import (
	"maoguo/henan/misc/redis"
	"maoguo/henan/misc/utils/parse"
	"maoguo/henan/model/model_mysql/imMessageModel"

	"gopkg.in/fatih/set.v0"
)

type MsgStoreAndNotifyServiceImpl struct {
}

func (msgStoreAndNotifyService *MsgStoreAndNotifyServiceImpl) GetImServiceUrl() string {
	return "123.57.47.19"
}

var notifyIpSet = set.New(set.ThreadSafe)

func (msgStoreAndNotifyService *MsgStoreAndNotifyServiceImpl) StoreMsgAndNotifyImServer(message *imMessageModel.ImMessage, id int64) {
	redis.LPush("msg_"+parse.Int64ToString(id), parse.ParseJson(message))
	addToSet(id)
}

func addToSet(id int64) {
	ip := redis.HGet("user_ip", parse.Int64ToString(id))
	notifyIpSet.Add(ip)
}
