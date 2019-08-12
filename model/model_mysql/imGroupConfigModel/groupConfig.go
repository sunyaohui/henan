package imGroupConfigModel

import . "maoguo/henan/model/model_mysql"

type ImGroupConfig struct {
	Id     int64   `json:"id"`
	Level  int16   `json:"level"`
	Expire int64   `json:"expire"`
	Fee    float64 `json:"fee"`
}

func (ImGroupConfig) TableName() string {
	return "im_group_config"
}

func Raws(sql string, args ...interface{}) (configs []ImGroupConfig) {
	Db.Raw(sql, args...).Scan(&configs)
	return
}
