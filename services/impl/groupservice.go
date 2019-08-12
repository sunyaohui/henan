package impl

import (
	"maoguo/henan/constants/imConstants"
	"maoguo/henan/misc/page"
	"maoguo/henan/misc/utils/parse"
	"maoguo/henan/model/model_mysql"
	"maoguo/henan/model/model_mysql/imGroupConfigModel"
	"maoguo/henan/model/model_mysql/imGroupMemberModel"
	"maoguo/henan/model/model_mysql/imGroupModel"
	"maoguo/henan/model/model_mysql/imGroupNoteModel"
	"maoguo/henan/result"
	"reflect"
	"strings"
	"time"
)

type GroupServiceImpl struct {
}

func (groupService *GroupServiceImpl) CheckUserGroup(userId int64) {
	list := model_mysql.QueryList("select * from im_group where createrId=? and level>0 and expire<?", userId, (time.Now().UnixNano() / 1e6))

	if len(list) > 0 {
		for _, value := range list {
			var maps map[string]interface{}
			maps["groupId"] = value["id"]
			maps["groupName"] = value["name"]
			maps["expire"] = value["expire"]

			msg := new(MessageServiceImpl).ModifyGroupNoticeCommon(userId, parse.StringToInt64(value["id"].(string)), parse.ParseJson(maps), imConstants.NOTICE_OF_GROUP_EXPIRE)
			if parse.StringToInt64(value["expire"].(string)) > (time.Now().UnixNano() / 1e6) {
				new(MsgStoreAndNotifyServiceImpl).StoreMsgAndNotifyImServer(msg, userId)
			}
		}
	}
}

func (this *GroupServiceImpl) UpdateGroup(userId, groupId int64, name string) int {
	row := 0
	group := imGroupModel.Raw("select * from im_group where id=? ", groupId)

	if !reflect.DeepEqual(group, imGroupModel.ImGroup{}) {
		group.Name = name
		imGroupModel.Update(group)
		row = 1
		members := this.GetAllMemberFromCache(groupId)
		for _, member := range members {
			msg := new(MessageServiceImpl).ModifyGroupNotice(userId, groupId, name)
			new(MsgStoreAndNotifyServiceImpl).StoreMsgAndNotifyImServer(msg, parse.StringToInt64(member["userId"].(string)))
		}
	}
	return row
}

func (this *GroupServiceImpl) GetAllMemberFromCache(groupId int64) []map[string]interface{} {
	list := model_mysql.QueryList("select a.userId,a.role,a.silence,IFNULL(a.markName,b.name) name,b.id,b.headUrl,b.province,b.city,b.district,b.sign from im_group_member a LEFT JOIN im_user b on(a.userId=b.id) where groupId=? and isAccept=1 group by userId ORDER BY rand()", groupId)
	return list
}

func (this *GroupServiceImpl) UpdateGroupMemberMark(userId, groupId int64, markName string) int {
	row := 0
	row = model_mysql.ExecInt("update im_group_member set markName=? where groupId=? and userId=?", markName, groupId, userId)
	members := this.GetAllMemberFromCache(groupId)
	for _, member := range members {
		memberId := parse.StringToInt64(member["userId"].(string))
		if memberId == userId {
			msg := new(MessageServiceImpl).ModifyGroupCommentNotice(userId, groupId, markName)
			new(MsgStoreAndNotifyServiceImpl).StoreMsgAndNotifyImServer(msg, parse.StringToInt64(member["userId"].(string)))
		}
	}
	this.CleanGroupMemberCache(groupId)

	return row
}

func (this *GroupServiceImpl) CleanGroupMemberCache(groupId int64) {

}

func (this *GroupServiceImpl) CacheGroupsInfo(userId int64) []map[string]interface{} {
	groupsInfo := model_mysql.QueryList("select b.*,a.receiveTip,a.markName,a.bgurl,a.role,a.silence from (select markName,groupId,userId,receiveTip,bgurl,role,silence from im_group_member where userId=? and isAccept=1 group by groupId,userId,receiveTip) a LEFT JOIN im_group b on(a.groupId=b.id)", userId)
	return groupsInfo
}

func (this *GroupServiceImpl) GroupInfo(groupId int64) imGroupModel.ImGroup {
	group := imGroupModel.Raw("select * from im_group where id=?", groupId)
	return group
}

func (this *GroupServiceImpl) GetGroupConfig() []imGroupConfigModel.ImGroupConfig {
	list := imGroupConfigModel.Raws("select id as lid,level,expire,fee from im_group_config where level>? order by level asc", 0)
	return list
}

func (this *GroupServiceImpl) RequestGroupJion(userId, destId int64, destType int) imGroupModel.ImGroup {
	groupMembers := this.GetAllMemberFromCache(destId)
	if len(groupMembers) > 700 {
		return imGroupModel.ImGroup{}
	}
	group := imGroupModel.Raw("select * from im_group where id=?", destId)
	model_mysql.Exec("delete from im_group_member where userId=? and groupId=?", userId, destId)
	member := imGroupMemberModel.ImGroupMember{
		CreateTime: time.Now().UnixNano() / 1e6,
		CreatorId:  userId,
		Role:       3,
		IsAccept:   1,
		GroupId:    destId,
		ReceiveTip: 1,
	}
	imGroupMemberModel.Save(&member)
	this.CacheGroupsInfo(userId)
	for _, users := range groupMembers {
		noticeUserId := parse.StringToInt64(users["userId"].(string))
		if noticeUserId != userId {
			msg := new(MessageServiceImpl).JoinedGroupNotice(userId, destId)
			new(MsgStoreAndNotifyServiceImpl).StoreMsgAndNotifyImServer(&msg, noticeUserId)
		}
	}
	this.CleanGroupMemberCache(destId)
	return group
}

func (this *GroupServiceImpl) UpdateQuitGroup(groupId, userId int64) {
	message := new(MessageServiceImpl).QuitGroupNotice(userId, groupId)
	list := this.GetAllMemberFromCache(groupId)
	for _, item := range list {
		memberId := parse.StringToInt64(item["userId"].(string))
		if memberId != userId {
			new(MsgStoreAndNotifyServiceImpl).StoreMsgAndNotifyImServer(&message, memberId)
		}
	}
	model_mysql.Exec("delete from im_group_member where groupId=? and userId=?", groupId, userId)
	this.CacheGroupsInfo(userId)
	this.CleanGroupMemberCache(groupId)
}

func (this *GroupServiceImpl) InviteJoinGroup(inviteIds string, groupId, userId int64) {
	userIds := parse.StringToInt64Arr(inviteIds)
	allMember := this.GetAllMemberFromCache(groupId)
	for _, memberId := range userIds {
		model_mysql.Exec("delete from im_group_member where groupId=? and userId=?", groupId, memberId)
		member := imGroupMemberModel.ImGroupMember{
			CreateTime: time.Now().UnixNano() / 1e6,
			CreatorId:  userId,
			UserId:     memberId,
			Role:       3,
			IsAccept:   1,
			GroupId:    groupId,
			ReceiveTip: 1,
		}
		imGroupMemberModel.Save(&member)
		msg := new(MessageServiceImpl).InviteGroupNotice(memberId, groupId)
		new(MsgStoreAndNotifyServiceImpl).StoreMsgAndNotifyImServer(msg, memberId)
		for _, otherMember := range allMember {
			otherMemberId := parse.StringToInt64(otherMember["userId"].(string))
			if otherMemberId != userId {
				joinMsg := new(MessageServiceImpl).JoinedGroupNotice(memberId, groupId)
				new(MsgStoreAndNotifyServiceImpl).StoreMsgAndNotifyImServer(&joinMsg, otherMemberId)
			}
		}
	}
	this.CleanGroupMemberCache(groupId)
}

func (this *GroupServiceImpl) UpdateGroupNote(userId, noteId int64, content, title string) int {
	row := 0
	note := imGroupNoteModel.Raw("select * from im_group_note where id=?", noteId)
	if reflect.DeepEqual(note, imGroupNoteModel.ImGroupNote{}) {
		return row
	}
	if strings.EqualFold(title, "") {
		title = note.Title
	} else {
		note.Title = title
	}
	if strings.EqualFold(content, "") {
		content = note.Content
	} else {
		note.Content = content
	}
	row = model_mysql.ExecInt("update im_group_note set content=?,title =? where id=?", content, title, noteId)
	if row > 0 {
		this.SendNoticeToMembers(userId, note.GroupId, parse.ParseJson(note), imConstants.MSG_TYPE_GROUP_NOTE_UPDATE)
	}
	return row
}

func (this *GroupServiceImpl) SendNoticeToMembers(userId, groupId int64, content string, msgType int) int {
	row := 0
	group := imGroupModel.Raw("select * from im_group where id=? ", groupId)
	if !reflect.DeepEqual(group, imGroupModel.ImGroup{}) {
		row = 1
		members := this.GetAllMemberFromCache(groupId)
		for _, member := range members {
			msg := new(MessageServiceImpl).ModifyGroupNoticeCommon(userId, groupId, content, msgType)
			new(MsgStoreAndNotifyServiceImpl).StoreMsgAndNotifyImServer(msg, parse.StringToInt64(member["userId"].(string)))
		}
	}
	return row
}

func (this *GroupServiceImpl) GetNoteList(userId, groupId int64, pageNo int) page.Page {
	if pageNo == 0 {
		pageNo = 1
	}
	pg := page.Page{
		PageNo:   1,
		PageSize: 20,
	}
	pg = model_mysql.QueryPage("select gn.* from im_group_note gn inner join im_group_member gm on gm.groupId=gn.groupId where gn.groupId = ? and gm.userId = ? order by gn.updateTime desc", pg, groupId, userId)
	return pg
}

func (this *GroupServiceImpl) TransGroupRequest(userId, groupId, destId int64) string {
	group := imGroupModel.Raw("select * from im_group where id=? ", groupId)
	if !reflect.DeepEqual(group, imGroupModel.ImGroup{}) {
		if userId == destId {
			return result.ResponseWrite(-1, "无法转让给自己")
		}
		member := imGroupMemberModel.Raw("select * from im_group_member where userId=? and groupId=?", userId, groupId)
		memberDest := imGroupMemberModel.Raw("select * from im_group_member where userId=? and groupId=?", destId, groupId)
		if !reflect.DeepEqual(member, imGroupMemberModel.ImGroupMember{}) {
			role := member.Role
			if role != 1 {
				return result.ResponseWrite(-1, "无权操作")
			}
			if !reflect.DeepEqual(memberDest, imGroupMemberModel.ImGroupMember{}) {
				return result.ResponseWrite(-1, "请选择群组成员")
			}
			row := model_mysql.QueryList("select * from im_group_trans_log where status=0 and userId=" + parse.Int64ToString(userId) + " and groupId=" + parse.Int64ToString(groupId) + " and destId=" + parse.Int64ToString(destId))
			if len(row) > 0 {
				content := map[string]interface{}{
					"groupId":   groupId,
					"groupName": group.Name,
					"destId":    userId,
				}
				msg := new(MessageServiceImpl).ModifyGroupNoticeCommon(userId, groupId, parse.ParseJson(content), imConstants.NOTICE_OF_GROUP_TRANS)
				new(MsgStoreAndNotifyServiceImpl).StoreMsgAndNotifyImServer(msg, destId)
				return result.ResponseWrite(1, "发送请求成功")
			}
			query := model_mysql.ExecInt("insert into im_group_trans_log (userId,groupId,destId) values(?,?,?)", userId, groupId, destId)
			if query > 0 {
				content := map[string]interface{}{
					"groupId":   groupId,
					"groupName": group.Name,
					"destId":    userId,
				}
				msg := new(MessageServiceImpl).ModifyGroupNoticeCommon(userId, groupId, parse.ParseJson(content), imConstants.NOTICE_OF_GROUP_TRANS)
				new(MsgStoreAndNotifyServiceImpl).StoreMsgAndNotifyImServer(msg, destId)
				return result.ResponseWrite(1, "发送请求成功")
			}
			return result.ResponseWrite(-1, "操作失败")
		}
		return result.ResponseWrite(-1, "您不在该群组，无法操作")
	}
	return result.ResponseWrite(-1, "群资料获取失败")
}

func (this *GroupServiceImpl) DismissGroup(userId, groupId int64) string {
	group := imGroupModel.Raw("select * from im_group where id=? ", groupId)
	if !reflect.DeepEqual(group, imGroupModel.ImGroup{}) {
		member := model_mysql.Query("select * from im_group_member where userId=? and groupId=? and isAccept=1", userId, groupId)
		if !reflect.DeepEqual(member, imGroupMemberModel.ImGroupMember{}) {
			role := parse.StringToIntDefault(member["role"].(string), 3)
			if role != 1 {
				return result.ResponseWrite(-1, "无权操作")
			}
			members := this.GetAllMemberFromCache(groupId)
			model_mysql.Exec("delete from im_group where id=? and createrId=?", groupId, userId)
			model_mysql.Exec("delete from im_group_member where groupId=?", groupId)
			content := map[string]interface{}{
				"groupId":   groupId,
				"groupName": group.Name,
				"destId":    userId,
			}
			for _, mem := range members {
				if parse.StringToInt64(mem["userId"].(string)) == userId {
					continue
				}
				msg := new(MessageServiceImpl).ModifyGroupNoticeCommon(userId, groupId, parse.ParseJson(content), imConstants.NOTICE_OF_GROUP_DISMISS)
				new(MsgStoreAndNotifyServiceImpl).StoreMsgAndNotifyImServer(msg, parse.StringToInt64(mem["userId"].(string)))
			}
			this.CacheGroupsInfo(userId)
			this.CleanGroupMemberCache(groupId)
			return result.ResponseWrite(1, "操作成功")
		}
		return result.ResponseWrite(-1, "您不在该群组，无法操作")
	}
	return result.ResponseWrite(-1, "群已经解散")
}

func (this *GroupServiceImpl) UpdateGroupDescriptions(userId, groupId int64, descriptions string) string {
	group := imGroupModel.Raw("select * from im_group where id=? ", groupId)
	if !reflect.DeepEqual(group, imGroupModel.ImGroup{}) {
		mem := imGroupMemberModel.Raw("select * from im_group_member where userId=? and groupId=?", userId, groupId)
		role := mem.Role
		if role != 1 {
			return result.ResponseWrite(imConstants.RESPONSE_GROUP_2001, "无权操作")
		}
		group.Descriptions = descriptions
		imGroupModel.Update(group)
		members := this.GetAllMemberFromCache(groupId)
		content := map[string]interface{}{"descriptions": descriptions}
		for _, member := range members {
			msg := new(MessageServiceImpl).ModifyGroupNoticeCommon(userId, groupId, parse.ParseJson(content), imConstants.NOTICE_OF_GROUP_DESCRIPTIONS)
			new(MsgStoreAndNotifyServiceImpl).StoreMsgAndNotifyImServer(msg, parse.StringToInt64(member["userId"].(string)))
		}
		return result.ResponseWrite(1, "设置成功")
	}
	return result.ResponseWrite(imConstants.RESPONSE_GROUP_2000, "群资料获取失败")
}

func (this *GroupServiceImpl) TransGroupConfirm(userId, groupId int64, accept int) string {
	group := imGroupModel.Raw("select * from im_group where id=? ", groupId)
	if !reflect.DeepEqual(group, imGroupModel.ImGroup{}) {
		destId := group.CreaterId
		member := imGroupMemberModel.Raw("select * from im_group_member where userId=? and groupId=?", userId, groupId)
		memberDest := imGroupMemberModel.Raw("select * from im_group_member where userId=? and groupId=?", destId, groupId)
		if !reflect.DeepEqual(member, imGroupMemberModel.ImGroupMember{}) {
			if reflect.DeepEqual(memberDest, imGroupMemberModel.ImGroupMember{}) {
				return result.ResponseWrite(-1, "群主信息获取失败")
			}
			row := model_mysql.QueryList("select * from im_group_trans_log where status=0 and userId=? and groupId=? and destId=?", destId, groupId, userId)
			if len(row) == 0 {
				return result.ResponseWrite(-1, "操作已失效")
			}
			query := model_mysql.ExecInt("update im_group_trans_log set status=? where userId=? and destId=? and groupId=?", accept, destId, userId, groupId)
			if query > 0 {
				content := map[string]interface{}{
					"groupId":   groupId,
					"groupName": group.Name,
					"destId":    userId,
				}
				msg := new(MessageServiceImpl).ModifyGroupNoticeCommon(userId, groupId, parse.ParseJson(content), imConstants.NOTICE_OF_GROUP_TRANS_NO)
				if accept == 1 {
					group.CreaterId = userId
					member.Role = 1
					memberDest.Role = 3
					msg = new(MessageServiceImpl).ModifyGroupNoticeCommon(userId, groupId, parse.ParseJson(content), imConstants.NOTICE_OF_GROUP_TRANS_YES)
				}
				new(MsgStoreAndNotifyServiceImpl).StoreMsgAndNotifyImServer(msg, destId)
				return result.ResponseWrite(1, "操作成功")
			}
			return result.ResponseWrite(-1, "操作失败")
		}
		return result.ResponseWrite(-1, "您不在该群组，无法操作")
	}
	return result.ResponseWrite(-1, "群资料获取失败")
}

func (this *GroupServiceImpl) UpdateGroupJoinStatus(groupId int64, status int, question, answer string, fee float64) int {
	row := 0
	row = model_mysql.ExecInt("update im_group set joinStatus=?,question=?,answer=?,fee=? where id=?", status, question, answer, fee, groupId)
	return row
}

func (this *GroupServiceImpl) RequestGroupJionWithQuestion(userId, destId int64, answer string) int {
	group := imGroupModel.Raw("select * from im_group where id=?", destId)
	groupMembers := this.GetAllMemberFromCache(destId)
	if (len(groupMembers) >= 100 && group.Level == 0) || (len(groupMembers) >= 200 && group.Level == 1) ||
		(len(groupMembers) >= 500 && group.Level == 2) || (len(groupMembers) >= 1000 && group.Level == 3) ||
		(len(groupMembers) >= 2000 && group.Level == 4) {
		return 5 //	群人数打到上限
	}
	model_mysql.Exec("delete from im_group_member where userId=? and groupId=?", userId, destId)
	member := imGroupMemberModel.ImGroupMember{
		CreateTime: time.Now().UnixNano() / 1e6,
		CreatorId:  userId,
		UserId:     userId,
		Role:       3,
		IsAccept:   1,
		GroupId:    destId,
		ReceiveTip: 1,
	}
	if group.JoinStatus == 1 || group.JoinStatus == 2 {
		member.IsAccept = 0
	}
	if (group.JoinStatus == 1 || group.JoinStatus == 2) && strings.EqualFold(answer, "") {
		return 1 //	 验证消息不能为空
	}
	if group.JoinStatus == 3 && !strings.EqualFold(answer, group.Answer) {
		return 2 //	 加群答案不能为空
	}
	if group.JoinStatus == 4 || group.JoinStatus == 5 {
		return 3 //	 不允许主动加入群
	}
	if group.JoinStatus == 6 && this.RequestGroupCheck(userId, destId) {
		return 4 //	没有付费
	}
	imGroupMemberModel.Save(&member)
	this.CacheGroupsInfo(userId)
	if member.IsAccept == 1 {
		for _, users := range groupMembers {
			noticeUserId := parse.StringToInt64(users["userId"].(string))
			if noticeUserId != userId {
				msg := new(MessageServiceImpl).JoinedGroupNotice(userId, destId)
				new(MsgStoreAndNotifyServiceImpl).StoreMsgAndNotifyImServer(&msg, noticeUserId)
			}
		}
	}
	if member.IsAccept == 0 {
		for _, users := range groupMembers {
			noticeUserId := parse.StringToInt64(users["userId"].(string))
			role := parse.StringToInt(users["role"].(string))
			if role == 1 || role == 2 {
				content := map[string]interface{}{
					"groupId":    group.Id,
					"groupName":  group.Name,
					"question":   group.Question,
					"answer":     group.Answer,
					"userAnswer": answer,
				}
				msg := new(MessageServiceImpl).ModifyGroupNoticeCommon(userId, group.Id, parse.ParseJson(content), imConstants.NOTICE_OF_GROUP_REQUEST)
				new(MsgStoreAndNotifyServiceImpl).StoreMsgAndNotifyImServer(msg, noticeUserId)
			}
		}
	}
	this.CleanGroupMemberCache(destId)
	return 0
}

func (this *GroupServiceImpl) RequestGroupCheck(userId, destId int64) bool {
	return true
}

func (this *GroupServiceImpl) SetGroupAdmin(userId, groupId, destId int64, isSet int) string {
	group := imGroupModel.Raw("select * from im_group where id=? ", groupId)
	if reflect.DeepEqual(group, imGroupModel.ImGroup{}) {
		member := imGroupMemberModel.Raw("select * from im_group_member where userId=? and groupId=?", userId, groupId)
		memberDest := imGroupMemberModel.Raw("select * from im_group_member where userId=? and groupId=?", destId, groupId)
		if !reflect.DeepEqual(member, imGroupMemberModel.ImGroupMember{}) {
			role := member.Role
			if role != 1 {
				return result.ResponseWrite(imConstants.RESPONSE_GROUP_2001, "无权操作")
			}
			if reflect.DeepEqual(memberDest, imGroupMemberModel.ImGroupMember{}) {
				return result.ResponseWrite(imConstants.RESPONSE_GROUP_2003, "不是群组成员")
			}
			r := 3
			if isSet == 1 {
				r = 2
			}
			memberDest.Role = r
			imGroupMemberModel.Save(&memberDest)
			members := this.GetAllMemberFromCache(groupId)
			destUser := new(UserServiceImpl).GetImUser(destId)
			content := map[string]interface{}{
				"name":    destUser.Name,
				"headUrl": destUser.HeadUrl,
				"role":    r,
			}
			for _, mem := range members {
				msg := new(MessageServiceImpl).ModifyGroupNoticeCommon(destId, groupId, parse.ParseJson(content), imConstants.NOTICE_OF_GROUP_ADMIN_SET)
				new(MsgStoreAndNotifyServiceImpl).StoreMsgAndNotifyImServer(msg, parse.StringToInt64(mem["userId"].(string)))
			}
			return result.ResponseWrite(1, "设置成功")
		}
		return result.ResponseWrite(-1, "设置失败")
	}
	return result.ResponseWrite(imConstants.RESPONSE_GROUP_2000, "群资料获取失败")
}

func (this *GroupServiceImpl) UpdateGroupMember(userId, groupId, destId int64, types, ope int) string {
	if userId > 0 && groupId > 0 && destId > 0 && types > 0 {
		group := imGroupModel.Raw("select * from im_group where id=?", groupId)
		user := new(UserServiceImpl).GetImUser(userId)
		if types == 1 {
			row := 0
			if ope == 1 {
				row = this.UpdateGroupMemberField(destId, groupId, "isAccept", "1")
			} else {
				row = model_mysql.ExecInt("delete from im_group_member where userId=? and groupId=?", destId, groupId)
			}
			if row > 0 {
				content := map[string]interface{}{
					"groupId":   group.Id,
					"groupName": group.Name,
					"adminId":   userId,
					"adminName": user.Name,
				}
				msgType := imConstants.NOTICE_OF_GROUP_REQUEST_NO
				if ope == 1 {
					msgType = imConstants.NOTICE_OF_GROUP_REQUEST_YES
				}
				msg := new(MessageServiceImpl).ModifyGroupNoticeCommon(userId, group.Id, parse.ParseJson(content), msgType)
				new(MsgStoreAndNotifyServiceImpl).StoreMsgAndNotifyImServer(msg, destId)
				return result.ResponseWrite(1, "操作成功")
			}
			return result.ResponseWrite(-1, "操作失败，请稍后重试")
		} else if types == 2 {
			//	禁言
			silenceTime := parse.IntToInt64(ope)
			if silenceTime != 0 {
				silenceTime = silenceTime*1000 + (time.Now().UnixNano() / 1e6)
			}

			row := model_mysql.ExecInt("update im_group_member set silence=? where groupId=? and userId=?", silenceTime, groupId, destId)
			if row > 0 {
				content := map[string]interface{}{
					"groupId":     group.Id,
					"groupName":   group.Name,
					"adminId":     userId,
					"adminName":   user.Name,
					"silenceTime": silenceTime,
				}
				if silenceTime == 0 {
					msg := new(MessageServiceImpl).ModifyGroupNoticeCommon(userId, group.Id, parse.ParseJson(content), imConstants.NOTICE_OF_GROUP_SILENCE_NO)
					new(MsgStoreAndNotifyServiceImpl).StoreMsgAndNotifyImServer(msg, destId)
					return result.ResponseWrite(1, "取消禁言成功")
				} else {
					msg := new(MessageServiceImpl).ModifyGroupNoticeCommon(userId, group.Id, parse.ParseJson(content), imConstants.NOTICE_OF_GROUP_SILENCE_YES)
					new(MsgStoreAndNotifyServiceImpl).StoreMsgAndNotifyImServer(msg, destId)
					return result.ResponseWrite(1, "禁言成功")
				}
			}
			return result.ResponseWrite(-1, "操作失败，请稍后重试")
		} else {
			return result.ResponseWrite(-1, "网络繁忙，请稍后重试")
		}
	}
	return result.ResponseWrite(-1, "网络繁忙，请稍后重试")
}

func (this *GroupServiceImpl) UpdateGroupMemberField(userId, groupId int64, field, value string) int {
	row := model_mysql.ExecInt("update im_group_member set "+field+"='"+value+"' where groupId=? and userId=?", groupId, userId)
	return row
}
