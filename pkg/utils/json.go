package utils

import (
	"encoding/json"
)

func JsonToString(v interface{}) string {
	jsonMsg, _ := json.MarshalIndent(v, "", "  ")
	return string(jsonMsg)
}
