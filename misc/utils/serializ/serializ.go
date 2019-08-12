package serializ

import (
	"encoding/json"
	"maoguo/henan/misc/utils/parse"

	"github.com/wonderivan/logger"
)

/*

 */
func BinaryWrite(info interface{}) []byte {
	// buf := new(bytes.Buffer)
	// err := binary.Write(buf, binary.LittleEndian, info)
	// if err != nil {
	// 	logger.Error("binary.Write failed:", err)
	// }
	// return buf.Bytes()
	r, err := json.Marshal(info)
	if err != nil {
		logger.Error("binary.Write failed ", err)
		return nil
	}
	return r
}

func BinaryRead(b []byte, info interface{}) interface{} {
	// buf := bytes.NewBuffer(b)
	//
	// err := binary.Read(buf, binary.LittleEndian, &info)
	// if err != nil {
	// 	logger.Error("binary.Read failed:", err)
	// }
	// return info
	return parse.ByteToObj(b, info)

}
