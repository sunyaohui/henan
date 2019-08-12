package impl

import (
	"fmt"
	"maoguo/henan/constants/imConstants"
	"maoguo/henan/misc/redis"
	"maoguo/henan/misc/utils"
	"maoguo/henan/misc/utils/parse"
	"maoguo/henan/model/model_mysql/imFeedCommentModel"
	"maoguo/henan/model/model_mysql/imFeedModel"
	"maoguo/henan/model/model_mysql/imFeedPraiseModel"
	"maoguo/henan/model/model_mysql/imMessageModel"
	"maoguo/henan/model/model_mysql/imUserModel"
	"strings"
	"time"
)

type MessageServiceImpl struct {
}

func (messageService *MessageServiceImpl) ModifyGroupNoticeCommon(fromId, destId int64, content string, msgType int) *imMessageModel.ImMessage {
	user := new(UserServiceImpl).GetImUser(fromId)
	return &imMessageModel.ImMessage{
		FromName:     user.Name,
		ImageIconUrl: user.HeadUrl,
		DestId:       destId,
		FromId:       fromId,
		FromType:     imConstants.MSG_FROM_SYS,
		Content:      content,
		MessageType:  msgType,
		MsgId:        utils.GetUUID(),
		SendTime:     time.Now().UnixNano() / 1e6,
	}
}

func (messageService *MessageServiceImpl) OtherLoginNotice(fromId int64, content string) *imMessageModel.ImMessage {
	return &imMessageModel.ImMessage{
		DestId:   fromId,
		FromId:   fromId,
		Content:  content,
		FromType: imConstants.MSG_TYPE_OTHER_LOGIN,
		MsgId:    utils.GetUUID(),
		SendTime: time.Now().UnixNano() / 1e6,
	}
}

func (messageService *MessageServiceImpl) NoticeOffline(id string) {
	ip := redis.HGet("user_ip", id)
	if strings.EqualFold(ip, "") {
		fmt.Println("暂时关闭其他地方登陆")
	}
}

func (this *MessageServiceImpl) UserModifyProfileNotice(fromId, destId int64, user *imUserModel.ImUser) *imMessageModel.ImMessage {
	imUser := make(map[string]interface{}, 0)
	imUser["nickName"] = user.NickName
	imUser["district"] = user.District
	imUser["city"] = user.City
	imUser["province"] = user.Province
	imUser["name"] = user.Name
	imUser["sex"] = user.Sex
	imUser["sign"] = user.Sign
	imUser["headUrl"] = user.HeadUrl
	imUser["id"] = user.Id
	imUser["feedBackImage"] = user.FeedBackImage
	return &imMessageModel.ImMessage{
		DestId:       destId,
		FromId:       fromId,
		FromName:     user.Name,
		ImageIconUrl: user.HeadUrl,
		Content:      parse.ParseJson(imUser),
		FromType:     imConstants.MSG_FROM_SYS,
		MessageType:  imConstants.MSG_TYPE_MODIFY_PROFILE,
		MsgId:        utils.GetUUID(),
		SendTime:     time.Now().UnixNano() / 1e6,
	}
}

func (this *MessageServiceImpl) ModifyGroupNotice(fromId, destId int64, name string) *imMessageModel.ImMessage {
	user := new(UserServiceImpl).GetImUser(fromId)
	return &imMessageModel.ImMessage{
		FromName:     name,
		ImageIconUrl: user.HeadUrl,
		DestId:       destId,
		FromId:       fromId,
		FromType:     imConstants.MSG_FROM_SYS,
		MessageType:  imConstants.MSG_TYPE_MODIFY_GROUP,
		MsgId:        utils.GetUUID(),
		Content:      name,
		SendTime:     time.Now().UnixNano() / 1e6,
	}
}

func (this *MessageServiceImpl) ModifyGroupCommentNotice(fromId, destId int64, comment string) *imMessageModel.ImMessage {
	user := new(UserServiceImpl).GetImUser(fromId)
	return &imMessageModel.ImMessage{
		FromName:     user.Name,
		ImageIconUrl: user.HeadUrl,
		DestId:       destId,
		FromId:       fromId,
		FromType:     imConstants.MSG_FROM_SYS,
		MessageType:  imConstants.MSG_TYPE_MODIFY_GROUP_COMMENT,
		MsgId:        utils.GetUUID(),
		Content:      comment,
		SendTime:     time.Now().UnixNano() / 1e6,
	}
}

/**
 * 有人被邀请加入了群
 * @param fromId 被邀请人
 * @param destId 群id
 * @return imMessageModel.ImMessage
 */
func (this *MessageServiceImpl) InviteGroupNotice(fromId, destId int64) *imMessageModel.ImMessage {
	user := new(UserServiceImpl).GetImUser(fromId)
	return &imMessageModel.ImMessage{
		FromName:     user.Name,
		ImageIconUrl: user.HeadUrl,
		DestId:       destId,
		FromId:       fromId,
		FromType:     imConstants.MSG_FROM_SYS,
		MessageType:  imConstants.MSG_TYPE_INVITE_GROUP,
		MsgId:        utils.GetUUID(),
		Content:      "{\"name\":\"" + user.Name + "\",\"headUrl\":\"" + "\"}",
		SendTime:     time.Now().UnixNano() / 1e6,
	}
}

func (this *MessageServiceImpl) DelFriendNotice(fromId, destId int64) *imMessageModel.ImMessage {
	return &imMessageModel.ImMessage{
		DestId:      destId,
		FromId:      fromId,
		FromType:    imConstants.MSG_FROM_SYS,
		MessageType: imConstants.MSG_TYPE_DEL_FRIEND,
		MsgId:       utils.GetUUID(),
		SendTime:    time.Now().UnixNano() / 1e6,
	}
}

func (this *MessageServiceImpl) AcceptFriendNotice(fromId, destId int64) *imMessageModel.ImMessage {
	user := new(UserServiceImpl).GetImUser(fromId)
	return &imMessageModel.ImMessage{
		FromName:     user.Name,
		ImageIconUrl: user.HeadUrl,
		DestId:       destId,
		FromId:       fromId,
		FromType:     imConstants.MSG_FROM_SYS,
		MessageType:  imConstants.MSG_TYPE_ACCEPT_FRIEND,
		MsgId:        utils.GetUUID(),
		Content:      "{\"name\":\"" + user.Name + "\",\"headUrl\":\"" + user.HeadUrl + "\"}",
		SendTime:     time.Now().UnixNano() / 1e6,
	}
}

func (this *MessageServiceImpl) NewFeedNotice(feed *imFeedModel.ImFeed, destId int64) imMessageModel.ImMessage {
	user := new(UserServiceImpl).GetImUser(feed.UserId)
	return imMessageModel.ImMessage{
		DestId:       destId,
		FromId:       feed.UserId,
		FromName:     user.Name,
		ImageIconUrl: user.HeadUrl,
		Content:      "{\"feedId\":" + parse.Int64ToString(feed.Id) + ",\"msg\":\"" + feed.FeedText + "\",\"imgs\":\"" + feed.FeedImgs + "\"}",
		FromType:     imConstants.MSG_TYPE_NEW_FEED,
		MsgId:        utils.GetUUID(),
		SendTime:     time.Now().UnixNano() / 1e6,
	}
}

func (this *MessageServiceImpl) RemindFeedNotice(feed *imFeedModel.ImFeed, destId int64) imMessageModel.ImMessage {
	user := new(UserServiceImpl).GetImUser(feed.UserId)
	return imMessageModel.ImMessage{
		DestId:       destId,
		FromId:       feed.UserId,
		FromName:     user.Name,
		ImageIconUrl: user.HeadUrl,
		Content:      "{\"feedId\":" + parse.Int64ToString(feed.Id) + ",\"msg\":\"\",\"feedmsg\":\"" + feed.FeedText + "\",\"imgs\":\"" + feed.FeedImgs + "\"}",
		FromType:     imConstants.MSG_FROM_P2P,
		MessageType:  imConstants.MSG_TYPE_FEED_AT,
		MsgId:        utils.GetUUID(),
		SendTime:     time.Now().UnixNano() / 1e6,
	}
}

func (this *MessageServiceImpl) ReferFeedCommentNotice(comment *imFeedCommentModel.ImFeedComment, id int64) imMessageModel.ImMessage {
	user := new(UserServiceImpl).GetImUser(comment.UserId)
	return imMessageModel.ImMessage{
		DestId:       id,
		FromId:       comment.UserId,
		FromName:     user.Name,
		ImageIconUrl: user.HeadUrl,
		Content:      "{\"feedId\":" + parse.Int64ToString(comment.ImFeed.Id) + ",\"msg\":\"" + comment.CommentText + "\",\"feedmsg\":\"" + comment.ImFeed.FeedText + "\",\"imgs\":\"" + comment.ImFeed.FeedImgs + "\"}",
		FromType:     imConstants.MSG_FROM_P2P,
		MessageType:  imConstants.MSG_TYPE_REFER_FEED,
		MsgId:        utils.GetUUID(),
		SendTime:     time.Now().UnixNano() / 1e6,
	}
}

func (this *MessageServiceImpl) ReferFeedPraiseNotice(praise *imFeedPraiseModel.ImFeedPraise, destId int64) imMessageModel.ImMessage {
	user := new(UserServiceImpl).GetImUser(praise.UserId)
	return imMessageModel.ImMessage{
		DestId:       destId,
		FromId:       praise.UserId,
		FromName:     user.Name,
		ImageIconUrl: user.HeadUrl,
		Content:      "{\"feedId\":" + parse.Int64ToString(praise.ImFeed.Id) + ",\"feedmsg\":\"" + praise.ImFeed.FeedText + "\",\"imgs\":\"" + praise.ImFeed.FeedImgs + "\"}",
		FromType:     imConstants.MSG_FROM_P2P,
		MessageType:  imConstants.MSG_TYPE_REFER_PRAISE,
		MsgId:        utils.GetUUID(),
		SendTime:     time.Now().UnixNano() / 1e6,
	}
}

func (this *MessageServiceImpl) FriendNotice(fromId, destId int64) imMessageModel.ImMessage {
	user := new(UserServiceImpl).GetImUser(fromId)
	return imMessageModel.ImMessage{
		FromName:     user.Name,
		ImageIconUrl: user.HeadUrl,
		Content:      "{\"name\":\"" + user.Name + "\",\"headUrl\":\"" + user.HeadUrl + "\"}",
		DestId:       destId,
		FromId:       fromId,
		FromType:     imConstants.MSG_FROM_SYS,
		MessageType:  imConstants.MSG_TYPE_AUTO_FRIEND,
		MsgId:        utils.GetUUID(),
		SendTime:     time.Now().UnixNano() / 1e6,
	}
}

//有人加入了群
func (this *MessageServiceImpl) JoinedGroupNotice(fromId, destId int64) imMessageModel.ImMessage {
	user := new(UserServiceImpl).GetImUser(fromId)
	return imMessageModel.ImMessage{
		DestId:       destId,
		FromName:     user.Name,
		ImageIconUrl: user.HeadUrl,
		FromId:       fromId,
		Content:      "{\"name\":\"" + user.Name + "\",\"headUrl\":\"" + user.HeadUrl + "\"}",
		FromType:     imConstants.MSG_FROM_SYS,
		MessageType:  imConstants.MSG_TYPE_JOIN_GROUP,
		MsgId:        utils.GetUUID(),
		SendTime:     time.Now().UnixNano() / 1e6,
	}
}

//有人退出了群
func (this *MessageServiceImpl) QuitGroupNotice(fromId, destId int64) imMessageModel.ImMessage {
	user := new(UserServiceImpl).GetImUser(fromId)
	return imMessageModel.ImMessage{
		FromName:     user.Name,
		ImageIconUrl: user.HeadUrl,
		DestId:       destId,
		FromId:       fromId,
		Content:      "{\"name\":\"" + user.Name + "\",\"headUrl\":\"" + user.HeadUrl + "\"}",
		FromType:     imConstants.MSG_FROM_SYS,
		MessageType:  imConstants.MSG_TYPE_QUIT_GROUP,
		MsgId:        utils.GetUUID(),
		SendTime:     time.Now().UnixNano() / 1e6,
	}
}
