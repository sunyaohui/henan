package getParameter

import (
	"encoding/json"
	"net/http"
)

//接收request中的参数，并转话为map
func GetParameter(request *http.Request) map[string]interface{} {
	var params map[string]interface{}
	decoder := json.NewDecoder(request.Body)
	decoder.Decode(&params)
	return params
}
