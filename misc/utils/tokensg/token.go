package tokensg

import (
	"encoding/base64"
	"fmt"
	"maoguo/henan/misc/redis"
	"maoguo/henan/misc/utils"
	"maoguo/henan/misc/utils/parse"
	"strings"

	"github.com/wonderivan/logger"
)

func TokenCoding(userid int64) string {
	uuid := utils.GetRand(8)
	fmt.Println(uuid)
	uid := utils.Encrypt(parse.Int64ToString(userid))
	return base64.StdEncoding.EncodeToString([]byte(uid + "," + uuid))
}

func TokenEncoding(token string) string {
	decodeBytes, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		logger.Error("token encoding failed", err)
	}
	return string(decodeBytes)
}

//验证token是否正确
func AssertToken(encoding string) bool {
	t := TokenEncoding(encoding)
	tos := strings.Split(t, ",")
	token := redis.HGet("token", tos[0])
	if strings.EqualFold(token, tos[1]) {
		return true
	}
	return false
}
