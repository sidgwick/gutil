package json

import (
	"bytes"
	"encoding/json"
	"errors"
	"reflect"
)

func Json(i interface{}) string {
	marshal, _ := json.Marshal(i)
	return string(marshal)
}

func IsJSON(s string) bool {
	var js interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}

func LoadData(dst interface{}, src interface{}) error {
	// 字符串/字符序列可以直接反序列化
	var err error
	var data []byte
	srcType := reflect.TypeOf(src).String()

	if srcType == "string" {
		data = []byte(src.(string))
	} else if srcType == "[]uint8" {
		data = src.([]uint8)
	} else {
		data, err = json.Marshal(src)
		if err != nil {
			return err
		}
	}

	if !IsJSON(string(data)) {
		return errors.New("not a valid JSON string")
	}

	d := json.NewDecoder(bytes.NewReader(data))
	d.UseNumber()
	err = d.Decode(dst)
	return err
}
