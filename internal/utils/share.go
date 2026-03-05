package utils

import "encoding/json"

func BytesToString(data []byte) string {
	if data == nil {
		return ""
	}
	return string(data[:])
}

func BytesToStruct(data []byte, v interface{}) error {
	if data == nil {
		return nil
	}
	return json.Unmarshal(data, v)
}
