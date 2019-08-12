package inteceptor

//登陆拦截器
func AssertLogin(tokens []string, url string) bool {

	//根绝url判断此请求是否过滤

	if len(tokens) == 0 {
		return false
	}
	if tokens[0] == "" {
		return false
	}

	//验证token是否正确
	// return tokensg.AssertToken(tokens[0])

	return true
}
