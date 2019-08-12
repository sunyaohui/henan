package imMessageModel

type ImMessage struct {
	Id           int64  `json:"id"`
	DevType      int    `json:"devType" gorm:"column:devType"`
	GeoId        int    `json:"geoId" gorm:"column:geoId"`
	MsgId        string `json:"msgId" gorm:"column:msgId;default:null"`
	FromId       int64  `json:"fromId" gorm:"column:fromId"`
	FromType     int    `json:"fromType" gorm:"column:fromType"`
	ImageIconUrl string `json:"imageIconUrl" gorm:"column:imageIconUrl"`
	DestId       int64  `json:"destId" gorm:"column:destId"`
	FromName     string `json:"fromName" gorm:"column:fromName;default:null"`
	Content      string `json:"content" gorm:"default:null"`
	MessageType  int    `json:"messageType" gorm:"column:messageType;default:null"`
	SendTime     int64  `json:"sendTime" gorm:"column:sendTime"`
	ReceiveTime  int64  `json:"receiveTime"`
	Version      int    `json:"version"`
	Status       int    `json:"status"`
}

func (ImMessage) TableName() string {
	return "im_message"
}
