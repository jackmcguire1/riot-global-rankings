package utils

import "encoding/json"

func ToJSON(i interface{}) string {
	data, _ := json.Marshal(i)
	return string(data)
}

func ToJsonBytes(i interface{}) json.RawMessage {
	data, _ := json.Marshal(i)
	return data
}

func PrettyPrintJson(i interface{}) json.RawMessage {
	data, _ := json.MarshalIndent(i, "", "\t")
	return data
}
