package sms

import (
	"fmt"
	"io/ioutil"
	"maoguo/henan/constants"
	"maoguo/henan/misc/config"
	"maoguo/henan/misc/redis"
	"maoguo/henan/misc/utils/parse"
	"net/http"
	"strings"

	"github.com/axgle/mahonia"
	"github.com/wonderivan/logger"
)

func SendValiCode(mobile, code, dowhat string) map[string]interface{} {
	url := config.CONFIG["SMS_URL"]
	user := config.CONFIG["SMS_USER"]
	pwd := config.CONFIG["SMS_PWD"]
	msg := strings.ReplaceAll(config.CONFIG[dowhat], "#code#", code)
	enc := mahonia.NewEncoder("gbk")
	resp, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader(fmt.Sprintf("user="+user+"&pwd="+pwd+"&mobile="+mobile+"&msg="+enc.ConvertString(msg))))
	if err != nil {
		logger.Error("send sms valicode failed", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error("send sms valicode failed", err)
	}
	str := string(body)
	jsonMap := parse.JsonToMap(str)
	if result, ok := jsonMap["result"]; ok {
		switch result.(float64) {
		case 100:
			return map[string]interface{}{
				"success": true,
				"desc":    "发送成功",
			}
		case 2:
			return map[string]interface{}{
				"success": false,
				"desc":    "手机号格式不正确",
			}
		case 22:
			return map[string]interface{}{
				"success": false,
				"desc":    "1小时内只能获取3次验证码",
			}
		case 33:
			return map[string]interface{}{
				"success": false,
				"desc":    "30秒内只能获取1次验证码",
			}
		case 20:
			return map[string]interface{}{
				"success": false,
				"desc":    "不支持该地区",
			}
		case 43:
			return map[string]interface{}{
				"success": false,
				"desc":    "今日验证码次数已达到上限",
			}
		case 3:
			return map[string]interface{}{
				"success": false,
				"desc":    "发送失败，请联系客服",
			}
		default:
			break
		}
	}
	return map[string]interface{}{
		"success": false,
		"desc":    "短信发送失败，请稍后重试",
	}
}

func EqualValidate(mobile, validateNum string) bool {
	num := redis.HGet(constants.SMS_KEY, mobile)
	if strings.EqualFold(num, validateNum) {
		return true
	}
	return false
}
