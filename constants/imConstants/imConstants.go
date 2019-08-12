package imConstants

const (
	DEV_PC                        = 3
	DEV_IOS                       = 2
	DEV_ANDROID                   = 1
	IM_EVENT_CHAT                 = "chat"
	IM_EVENT_HEART                = "heart"
	MSG_FROM_P2P                  = 1
	MSG_FROM_GROUP                = 2
	MSG_FROM_SYS                  = 3
	MSG_FROM_THIRD                = 4
	MSG_TYPE_READY                = 1
	MSG_TYPE_TEXT                 = 2
	MSG_TYPE_IMG                  = 3
	MSG_TYPE_FILE                 = 4
	MSG_TYPE_ONLINE               = 5
	MSG_TYPE_OFFLINE              = 6
	MSG_TYPE_DEL_FRIEND           = 7
	MSG_TYPE_JOIN_GROUP           = 8
	MSG_TYPE_QUIT_GROUP           = 9
	MSG_TYPE_INVITE_GROUP         = 10
	MSG_TYPE_FRIEND_REQ           = 11
	MSG_TYPE_ACCEPT_FRIEND        = 12
	MSG_TYPE_GROUP_REQ            = 13
	MSG_TYPE_ACCEPT_GROUP         = 14
	MSG_TYPE_OTHER_LOGIN          = 15
	MSG_TYPE_VOICE                = 16
	MSG_TYPE_RED_PACKET           = 17
	MSG_TYPE_TRANSFER             = 18
	MSG_TYPE_RECEIVE_RED_NOTICE   = 19
	MSG_TYPE_TRANSFER_REC         = 20
	MSG_TYPE_TRANSFER_BACK        = 21
	MSG_TYPE_REDPACKET_BACK       = 22
	MSG_TYPE_AUTO_FRIEND          = 23
	MSG_TYPE_RED_FINISHED         = 24
	MSG_TYPE_MODIFY_GROUP         = 25
	MSG_TYPE_MODIFY_PROFILE       = 26
	MSG_TYPE_MODIFY_GROUP_COMMENT = 27
	MSG_TYPE_SEND_CARD            = 28
	MSG_TYPE_SEND_LOCATION        = 29
	MSG_TYPE_SEND_VIDEO           = 30
	MSG_TYPE_HEARTBEAT            = 31
	MSG_TYPE_REBACK               = 32
	MSG_TYPE_READED               = 33
	MSG_TYPE_EMOJI_YUN            = 34
	MSG_TYPE_AT                   = 35
	MSG_TYPE_NEW_FEED             = 36
	MSG_TYPE_REFER_FEED           = 37
	MSG_TYPE_REFER_PRAISE         = 38
	MSG_TYPE_REC_MONEY            = 39
	MSG_TYPE_AA_RECEIVE           = 40

	RED_PACKET_TYPE_COMMAND = 3 // 口令红包

	MSG_TYPE_GROUP_HEAD = 46 //	  群头像修改

	MSG_TYPE_GROUP_NOTE_NEW = 47 //		新群公告

	MSG_TYPE_GROUP_NOTE_UPDATE = 48 //  群公告修改

	MSG_TYPE_FEED_AT = 49 //  朋友圈提到了谁

	MSG_TYPE_COMMAND_RED_NOTICE = 50 //  口令红包抢光通知

	NOTICE_OF_GROUP_DISMISS = 51 // 群解散通知

	NOTICE_OF_GROUP_TRANS = 52 // 群转让邀请

	NOTICE_OF_GROUP_TRANS_YES = 53 // 群转让邀请同意

	NOTICE_OF_GROUP_TRANS_NO = 54 // 群转让邀请拒绝

	NOTICE_OF_GROUP_ADMIN_SET = 55 // 群成员成为管理员

	NOTICE_OF_GROUP_DESCRIPTIONS = 56 // 群简介更新

	NOTICE_OF_GROUP_ADMIN_CANCEL = 57 // 群成员取消管理员

	NOTICE_OF_GROUP_EXPIRE = 58 // 群到期提醒

	NOTICE_OF_GROUP_REQUEST = 59 // 请求加群

	NOTICE_OF_GROUP_REQUEST_YES = 60 // 同意加群

	NOTICE_OF_GROUP_REQUEST_NO = 61 // 拒绝加群

	NOTICE_OF_GROUP_SILENCE_YES = 62 // 禁言成功

	NOTICE_OF_GROUP_SILENCE_NO = 63 // 取消禁言

	NOTICE_OF_GROUP_MEMBER_REMOVE = 64 // 群成员被删除

	NOTICE_OF_GROUP_UPDATE_EXPIRE = 65 // 群到期时间更新

	MSG_OF_NO_FRIEND = 66 // 验证非通过消息 拉黑消息 或者非好友

	NOTICE_OF_HELLO = 67 // 打招呼消息

	//设置APPID/AK/SK
	BAIDU_APP_ID = "10019750"

	BAIDU_API_KEY = "rUHQRermbFl91SU16Rs8q9jc"

	BAIDU_SECRET_KEY = "4a1a8527f73ddea6a6e5d67a575ab393"

	AIPAY_ORDER_KEY = "alipayOrder"

	AIPAY_ORDER_SUCCUSS_KEY = "alipaySuccess"

	WXPAY_ORDER_KEY = "wxpayOrder"

	WXPAY_ORDER_SUCCUSS_KEY = "wxpaySuccess"

	// 群返回code码
	RESPONSE_GROUP_2000 = 2000 // 群组不存在或已删除

	RESPONSE_GROUP_2001 = 2001 // 无权操作

	RESPONSE_GROUP_2003 = 2003 // 用户不是群成员

	//  钱包类型
	IM_NUMBER_KEY = "SamNumber"

	MONEY_TYPE_9 = 9 // 购买靓号

	// 分组设置
	MAX_FRIEND_GROUP = 20
)
