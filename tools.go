package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
)

func writeError(err []byte, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/json")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(err)))
	w.Header().Set("X-Zc-Content-Length", fmt.Sprintf("%d", len(err)))
	w.Header().Set("X-Zc-Msg-Name", ZC_MSG_NAME_ERR)
	if len(err) <= 0 {
		w.Write([]byte(`{"error":"internal error", "errorCode":3000}`))
	} else {
		w.Write(err)
	}
}

func ParseError(oriErr []byte) (errMsg string, code int64, err error) {
	if len(oriErr) <= 0 {
		return
	}
	var errMap map[string]interface{}
	err = json.Unmarshal(oriErr, &errMap)
	if err != nil {
		return
	}
	errMsg = GetString(errMap["error"])
	code = GetInt(errMap["errorCode"])
	return
}

func GetInt(v interface{}) int64 {
	if v == nil {
		return 0
	}
	kind := reflect.TypeOf(v).Kind()
	switch kind {
	case reflect.Int:
		return int64(v.(int))
	case reflect.Int64:
		return int64(v.(int64))
	case reflect.Int8:
		return int64(v.(int8))
	case reflect.Int16:
		return int64(v.(int16))
	case reflect.Int32:
		return int64(v.(int32))
	case reflect.Uint:
		return int64(v.(uint))
	case reflect.Uint8:
		return int64(v.(uint8))
	case reflect.Uint16:
		return int64(v.(uint16))
	case reflect.Uint32:
		return int64(v.(uint32))
	case reflect.Uint64:
		return int64(v.(uint64))
	case reflect.Float32:
		return int64(v.(float32))
	case reflect.Float64:
		return int64(v.(float64))
	}
	return 0
}
func GetFloat(v interface{}) float64 {
	if v == nil {
		return 0
	}
	switch reflect.TypeOf(v).Kind() {
	case reflect.Float64, reflect.Float32:
		return v.(float64)
	default:
		return 0
	}
}
func GetString(v interface{}) string {
	s, ok := v.(string)
	if !ok {
		return ""
	}
	return s
}
