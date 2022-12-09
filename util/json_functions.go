package util

import (
	"bytes"
	"encoding/json"
)

func JSONCompactStringify(data interface{}) string {
	bytesOfObj, _ := json.Marshal(data)
	bytesBuffer := new(bytes.Buffer)
	json.Compact(bytesBuffer, bytesOfObj)
	return bytesBuffer.String()
}
