package paramter

import "strconv"

func GetImUserKey(userId int64) []byte {
	string := strconv.FormatInt(userId, 10)
	return []byte(string)

}
