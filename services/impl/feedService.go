package impl

import (
	"fmt"
	"maoguo/henan/misc/page"
	"maoguo/henan/misc/utils"
	"maoguo/henan/misc/utils/goMap"
	"maoguo/henan/misc/utils/parse"
	"maoguo/henan/model/model_mysql"
	"maoguo/henan/model/model_mysql/imFeedCommentModel"
	"maoguo/henan/model/model_mysql/imFeedModel"
	"maoguo/henan/model/model_mysql/imFeedPraiseModel"
	"maoguo/henan/model/model_mysql/imUserModel"
	"maoguo/henan/model/mongo"
	"reflect"
	"strings"
	"time"

	"github.com/wonderivan/logger"
	"gopkg.in/mgo.v2/bson"
)

type FeedServiceImpl struct {
}

func (this *FeedServiceImpl) UpdateFeed(userId int64, userName, headUrl string) {
	if userId != 0 {
		new(mongo.ImFeed).UpdateAll(bson.M{"userId": userId}, mongo.UpdateData("userName", userName))
	}
	if headUrl != "" {
		new(mongo.ImFeed).UpdateAll(bson.M{"userId": userId}, mongo.UpdateData("userHeadUrl", headUrl))
	}
}

func (this *FeedServiceImpl) DeleteFriendFeed(userId, destId int64) {
	q := bson.M{"belongUserId": userId, "userId": destId}
	new(mongo.ImFeed).Remove(q)

	q = bson.M{"userId": userId, "belongUserId": destId}
	new(mongo.ImFeed).Remove(q)
}

func (this *FeedServiceImpl) AddFeed(images, text string, userId int64, lat string, lng, address string, priv int,
	at, uids, video, ext string) imFeedModel.ImFeed {
	user := imUserModel.GetUserById(userId)

	imFeed := imFeedModel.ImFeed{
		CreateTime:  time.Now().UnixNano() / 1e6,
		FeedImgs:    images,
		FeedText:    text,
		Status:      1,
		UserId:      userId,
		UserName:    user.Name,
		UserHeadUrl: user.HeadUrl,
		Lng:         lng,
		Lat:         lat,
		Address:     address,
		Priv:        priv,
		At:          at,
		Uids:        uids,
		FeedVideo:   video,
		Ext:         ext,
	}
	imFeedModel.Save(&imFeed)
	feed := &mongo.ImFeed{
		CreateTime:  imFeed.CreateTime,
		FeedImgs:    imFeed.FeedImgs,
		FeedText:    imFeed.FeedText,
		Status:      imFeed.Status,
		UserId:      imFeed.UserId,
		UserName:    imFeed.UserName,
		UserHeadUrl: imFeed.UserHeadUrl,
		FeedId:      imFeed.Id,
		Lng:         lng,
		Lat:         lat,
		Address:     address,
		Priv:        priv,
		At:          at,
		Uids:        uids,
		FeedVideo:   video,
		Ext:         ext,
		Id:          bson.NewObjectId(),
	}
	feed.BelongUserId = userId
	feed.Insert()
	friends := new(ImFriendServiceImpl).GetFriends(userId)
	except := this.GetFeedExcept(userId)
	for _, friend := range friends {
		if utils.IsExistItem(except, friend.FriendId) {
			continue
		}
		flag := isAppend(friend.FriendId, uids, priv)
		bool := isAppend(friend.FriendId, at, 2)

		if flag {
			feed.Id = bson.NewObjectId()
			feed.BelongUserId = friend.FriendId
			feed.Insert()
			message := new(MessageServiceImpl).NewFeedNotice(&imFeed, friend.FriendId)
			new(MsgStoreAndNotifyServiceImpl).StoreMsgAndNotifyImServer(&message, friend.FriendId)
		}
		if bool {
			fuser := new(UserServiceImpl).GetImUser(friend.FriendId)
			atuser := &mongo.ImFeedUser{
				UserName:    fuser.NickName,
				UserHeadUrl: fuser.HeadUrl,
				UserId:      fuser.Id,
				FeedId:      imFeed.Id,
			}
			if !strings.EqualFold(friend.Remark, "") {
				atuser.UserName = friend.Remark
			}
			atuser.Save()

			in := []int64{fuser.Id, userId}
			query := bson.M{
				"feedId": imFeed.Id,
				"belongUserId": bson.M{
					"$in": in,
				},
			}
			params := bson.M{
				"atUsers": atuser,
			}
			new(mongo.ImFeed).UpdateAll(query, params)
			message := new(MessageServiceImpl).RemindFeedNotice(&imFeed, friend.FriendId)
			new(MsgStoreAndNotifyServiceImpl).StoreMsgAndNotifyImServer(&message, friend.FriendId)
		}
	}
	return imFeed
}

func isAppend(userId int64, uids string, priv int) bool {
	var users []string
	if !strings.EqualFold(uids, "") {
		users = strings.Split(uids, ",")
	}
	flag := false
	switch priv {
	case 0:
		flag = true
		break
	case 1:
		flag = false
		break
	case 2:
		flag = utils.IsExistItem(users, parse.Int64ToString(userId))
		break
	case 3:
		flag = !utils.IsExistItem(users, parse.Int64ToString(userId))
		break
	default:
		break
	}
	return flag
}

func (this *FeedServiceImpl) GetFeedExcept(userId int64) []int64 {
	list := model_mysql.QueryColumn("select friendId from im_friend where userId=? and isFriend=1 and isBlack=1", userId)
	flist := model_mysql.QueryColumn("select userId from im_friend where friendId=? and isFriend=1 and isBlack=1", userId)
	// 朋友圈设置了不给好友看的
	ufeedList := model_mysql.QueryColumn("select friendId from im_friend where userId=? and isFriend=1 and isBlack=0 and ufeedPriv=0", userId)
	// 朋友圈设置了不给我看的
	feedlist := model_mysql.QueryColumn("select userId from im_friend where friendId=? and isFriend=1 and isBlack=0 and feedPriv=0", userId)

	var lists []interface{}
	lists = append(lists, list...)
	lists = append(lists, flist...)
	lists = append(lists, ufeedList...)
	lists = append(lists, feedlist...)

	data := goMap.NewMap()
	for _, item := range lists {
		data.Put(item.(string), item)
	}
	var result []int64
	for item := range data.Map {
		result = append(result, parse.StringToInt64(item))
	}
	return result
}

func (this *FeedServiceImpl) GetOneFeed(feedId, userId int64) mongo.ImFeed {
	query := bson.M{
		"feedId":       feedId,
		"belongUserId": userId,
	}
	feed := new(mongo.ImFeed).FindOne(query)
	return feed
}

func (this *FeedServiceImpl) AddFeedComment(feedId int64, text string, userId, replyId int64) {
	dbFeed := imFeedModel.QueryById(feedId)
	user := imUserModel.GetUserById(userId)
	replyUser := imUserModel.GetUserById(replyId)
	//  三方有一方设置了权限或者
	expcets := this.GetFeedExcept(dbFeed.UserId)
	if utils.IsExistItem(expcets, parse.Int64ToString(userId)) {
		return
	}
	dbComment := &imFeedCommentModel.ImFeedComment{
		CommentText: text,
		CreateTime:  time.Now().UnixNano() / 1e6,
		FeedId:      dbFeed.Id,
		ReplyUserId: replyId,
		UserId:      userId,
		UserName:    user.Name,
		UserHeadUrl: user.HeadUrl,
		ImFeed:      &dbFeed,
	}
	if !reflect.DeepEqual(replyUser, &imUserModel.ImUser{}) && !strings.EqualFold(replyUser.Name, "") {
		dbComment.ReplyUserName = replyUser.Name
		dbComment.ReplyUserHeadUrl = replyUser.HeadUrl
	}
	dbComment.Create()
	comment := &mongo.ImFeedComment{
		Id:               dbComment.Id,
		CreateTime:       dbComment.CreateTime,
		CommentText:      dbComment.CommentText,
		ReplyUserId:      replyId,
		UserId:           userId,
		UserHeadUrl:      dbComment.UserHeadUrl,
		UserName:         dbComment.UserName,
		ReplyUserHeadUrl: dbComment.ReplyUserHeadUrl,
		ReplyUserName:    dbComment.ReplyUserName,
	}
	comment.Insert()
	feedCreatorId := dbFeed.UserId
	ids := new(ImFriendServiceImpl).GetMixFriend(feedCreatorId, userId, replyId)

	// params := bson.M{
	// 	"imFeedComments": []*mongo.ImFeedComment{comment},
	// }

	params := bson.M{"imFeedComments": comment}

	for _, item := range ids {
		id := parse.StringToInt64(item)
		query := bson.M{
			"feedId": feedId,
			"belongUserId": bson.M{
				"$in": []int64{id},
			},
		}

		new(mongo.ImFeed).UpdatePush(query, params)
	}

	push := []int64{userId, feedCreatorId}
	if replyId > 0 {
		push = append(push, replyId)
	}
	querys := bson.M{
		"feedId": feedId,
		"belongUserId": bson.M{
			"$in": push,
		},
	}
	new(mongo.ImFeed).UpdatePush(querys, params)

	findOneQuery := bson.M{
		"feedId": feedId,
		"belongUserId": bson.M{
			"$in": []int64{userId},
		},
	}
	feed := new(mongo.ImFeed).FindOne(findOneQuery)

	friendsMap := make(map[int64]int64)
	friendsMap[feedCreatorId] = feedCreatorId
	if feed.ImFeedComments != nil && len(feed.ImFeedComments) > 0 {
		for _, item := range feed.ImFeedComments {
			friendsMap[item.UserId] = item.UserId
		}
	}

	if feed.ImFeedPraises != nil && len(feed.ImFeedPraises) > 0 {
		for _, item := range feed.ImFeedPraises {
			friendsMap[item.UserId] = item.UserId
		}
	}
	delete(friendsMap, userId)
	for _, id := range friendsMap {
		msg := new(MessageServiceImpl).ReferFeedCommentNotice(dbComment, id)
		new(MsgStoreAndNotifyServiceImpl).StoreMsgAndNotifyImServer(&msg, id)
	}

}

func (this *FeedServiceImpl) AddFeedPraise(feedId, userId int64) {
	dbFeed := imFeedModel.QueryById(feedId)
	user := imUserModel.GetUserById(userId)

	expcets := this.GetFeedExcept(dbFeed.UserId)
	if utils.IsExistItem(expcets, parse.Int64ToString(userId)) {
		return
	}
	dbPraise := imFeedPraiseModel.ImFeedPraise{
		CreateTime:  time.Now().UnixNano() / 1e6,
		ImFeed:      dbFeed,
		FeedId:      dbFeed.Id,
		UserId:      userId,
		UserName:    user.Name,
		UserHeadUrl: user.HeadUrl,
	}
	imFeedPraiseModel.Create(&dbPraise)

	praise := &mongo.ImFeedPraise{
		CreateTime:  dbPraise.CreateTime,
		Id:          dbPraise.Id,
		UserId:      dbPraise.UserId,
		UserName:    dbPraise.UserName,
		UserHeadUrl: dbPraise.UserHeadUrl,
	}
	praise.Insert()
	feedCreatorId := dbFeed.UserId
	//赞我的朋友的动态,推送给我与好友交集的好友。
	ids := new(ImFriendServiceImpl).GetMixFriend(feedCreatorId, userId, 0)

	params := bson.M{"imFeedPraises": praise}

	for _, id := range ids {
		id := parse.StringToInt64(id)
		query := bson.M{
			"feedId": feedId,
			"belongUserId": bson.M{
				"$in": []int64{id},
			},
		}
		new(mongo.ImFeed).UpdatePush(query, params)
	}
	query := bson.M{
		"feedId": feedId,
		"belongUserId": bson.M{
			"$in": []int64{feedCreatorId, userId},
		},
	}
	new(mongo.ImFeed).UpdatePush(query, params)

	findOneQuery := bson.M{
		"feedId": feedId,
		"belongUserId": bson.M{
			"$in": []int64{userId},
		},
	}
	feed := new(mongo.ImFeed).FindOne(findOneQuery)
	friendIds := map[int64]int64{feedCreatorId: feedCreatorId}
	if feed.ImFeedComments != nil && len(feed.ImFeedComments) > 0 {
		for _, item := range feed.ImFeedComments {
			friendIds[item.UserId] = item.UserId
		}
	}

	if feed.ImFeedPraises != nil && len(feed.ImFeedPraises) > 0 {
		for _, item := range feed.ImFeedPraises {
			friendIds[item.UserId] = item.UserId
		}
	}
	delete(friendIds, userId)
	for _, id := range friendIds {
		msg := new(MessageServiceImpl).ReferFeedPraiseNotice(&dbPraise, id)
		new(MsgStoreAndNotifyServiceImpl).StoreMsgAndNotifyImServer(&msg, id)
	}
}

func (this *FeedServiceImpl) CanclePraise(feedId, userId int64) {
	model_mysql.Exec("delete from im_feed_praise where feed_id=? and user_id=?", feedId, userId)
	query := bson.M{"feedId": feedId}
	feeds := new(mongo.ImFeed).Find(query)
	if feeds != nil && len(feeds) > 0 {
		for _, feed := range feeds {
			praises := feed.ImFeedPraises
			if praises != nil && len(praises) > 0 {
				for _, praise := range praises {
					if praise.UserId == userId {
						new(mongo.ImFeed).UpdatePull(query, bson.M{"imFeedPraises": praise})
					}
				}
			}
		}
	}
}

func (this *FeedServiceImpl) GetMyFeed(userId int64, pg *page.Page) {
	query := bson.M{
		"userId":       userId,
		"belongUserId": userId,
	}
	//查询总数
	count := new(mongo.ImFeed).Count(query)
	pg.TotalCount = count
	var datas []mongo.ImFeed
	mongo.GetImFeed().Find(query).Skip((pg.PageNo - 1) * pg.PageSize).Limit(pg.PageSize).Sort("-createTime").All(&datas)
	if datas != nil && len(datas) > 0 {
		for i, _ := range datas {
			if len(datas[i].ImFeedPraises) == 0 {
				datas[i].ImFeedPraises = make([]mongo.ImFeedPraise, 0)
			}
			if len(datas[i].ImFeedComments) == 0 {
				datas[i].ImFeedComments = make([]mongo.ImFeedComment, 0)
			}
			if len(datas[i].AtUsers) == 0 {
				datas[i].AtUsers = make([]mongo.ImFeedUser, 0)
			}
		}
		pg.List = datas
	} else {
		pg.List = make([]mongo.ImFeed, 0)
	}
}

func (this *FeedServiceImpl) GetFeed(userId int64, pg *page.Page) {
	blackid := new(UserServiceImpl).GetUserBlackList(userId)
	query := bson.M{
		"belongUserId": userId,
		"userId": bson.M{
			"$nin": blackid,
		},
	}
	//查询总数
	count := new(mongo.ImFeed).Count(query)
	pg.TotalCount = count
	var datas []mongo.ImFeed
	mongo.GetImFeed().Find(query).Skip((pg.PageNo - 1) * pg.PageSize).Limit(pg.PageSize).Sort("-createTime").All(&datas)
	if datas != nil && len(datas) > 0 {
		for i, _ := range datas {
			if len(datas[i].ImFeedPraises) == 0 {
				datas[i].ImFeedPraises = make([]mongo.ImFeedPraise, 0)
			}
			if len(datas[i].ImFeedComments) == 0 {
				datas[i].ImFeedComments = make([]mongo.ImFeedComment, 0)
			}
			if len(datas[i].AtUsers) == 0 {
				datas[i].AtUsers = make([]mongo.ImFeedUser, 0)
			}
		}
		pg.List = datas
	} else {
		pg.List = make([]mongo.ImFeed, 0)
	}

}

func (this *FeedServiceImpl) GetFeedAuth(userId, friendId int64) map[string]interface{} {
	res := model_mysql.Query("select ufeedPriv as feedPrivUser,feedPriv as feedPrivFriend,isBlack,remark from im_friend where userId=? and friendId=?", userId, friendId)
	r := goMap.NewMap()
	if res == nil {
		r.Put("feedPrivUser", "0")
		r.Put("feedPrivFriend", "0")
		r.Put("isBlack", "0")
		r.Put("remark", "")
		return r.Map
	}
	logger.Info(parse.ParseJson(res))
	r.Put("feedPrivUser", getMapDefault("feedPrivUser", "0", res))
	r.Put("feedPrivFriend", getMapDefault("feedPrivFriend", "0", res))
	r.Put("isBlack", getMapDefault("isBlack", "0", res))
	r.Put("remark", getMapDefault("remark", "0", res))
	return r.Map
}

func getMapDefault(k, defaultValue string, maps map[string]interface{}) string {
	if v, ok := maps[k]; ok {
		switch v.(type) {
		case string:
			return v.(string)
		default:
			return defaultValue
		}
	}
	return defaultValue
}

func (this *FeedServiceImpl) DeleteFeed(feedId, userId int64) {
	model_mysql.Exec("delete from im_feed where id=? and user_id=?", feedId, userId)
	query := bson.M{"feedId": feedId}
	new(mongo.ImFeed).Remove(query)
}

func (this *FeedServiceImpl) SetFeedAuth(userId, friendId int64, priv, direct int) int {
	field := ""
	if direct == 0 {
		field = "ufeedPriv"
	} else {
		field = "feedPriv"
	}
	if model_mysql.ExecInt(fmt.Sprintf("update im_friend set %s = %d where userId = ? and friendId = ?", field, priv), userId, friendId) > 0 {
		if priv == 0 {
			var query bson.M
			if direct == 0 {
				query = bson.M{
					"belongUserId": friendId,
					"userId":       userId,
				}
			} else {
				query = bson.M{
					"belongUserId": userId,
					"userId":       friendId,
				}
			}
			mongo.GetImFeed().Remove(query)
		}
		return 1
	}
	return 0
}
