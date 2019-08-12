package impl

// func GetImUserPay(userId string) map[string]string {
// 	return cacheImUserPay(userId)
//
// }
//
// func cacheImUserPay(userId string) map[string]string {
// 	result, _ := utils.Query("select * from ImUserPay where userId = " + userId)
// 	return result
// }
//
// func GetTopList(userId string) []map[string]string {
// 	return cacheTopList(userId)
// }
//
// func cacheTopList(userId string) []map[string]string {
// 	list := utils.QueryList("select userId,destType,destId from im_top where userId=" + userId)
// 	return list
// }
//
// func GetFriendInfo(userId string) []map[string]string {
// 	return cacheFriendInfo(userId)
// }
//
// func cacheFriendInfo(userId string) []map[string]string {
// 	return utils.QueryList("select a.remark,a.isBlack,a.receiveTip,b.id,name,nickName,mobile,mail,sex,birthday,sign,province,isOnline,city,district,b.createTime,status,detail,headUrl,a.bgurl,b.imNumber as IMNo from im_friend a left join im_user b on(a.friendId=b.id) where a.isFriend=1 and a.userId=" + userId)
// }

//
// func GetGroupsInfo(userId string) []map[string]string {
// 	return cacheGroupsInfo(userId)
// }
//
// func cacheGroupsInfo(userId string) []map[string]string {
// 	// select b.*,a.receiveTip,a.markName,a.bgurl,a.role,a.silence from (select markName,groupId,userId,receiveTip,bgurl,role,silence from im_group_member where userId=? and isAccept=1 group by groupId,userId,receiveTip) a LEFT JOIN im_group b on(a.groupId=b.id)
// 	return utils.QueryList("select b.*,a.receiveTip,a.markName,a.bgurl,a.role,a.silence from (select markName,groupId,userId,receiveTip,bgurl,role,silence from im_group_member where userId=" + userId + " and isAccept=1 group by groupId,userId,receiveTip) a LEFT JOIN im_group b on(a.groupId=b.id)")
// }
