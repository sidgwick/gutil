package ghttp

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/url"
	"reflect"

	"github.com/spf13/cast"
)

func BuildGetQueryString(input interface{}) string {
	data := url.Values{}

	xInput := ToStringMap(input)
	for k, v := range xInput {
		data.Set(k, cast.ToString(v))
	}

	return data.Encode()
}

// ToStringMap casts an interface to a map[string]interface{} type.
func ToStringMap(i interface{}) map[string]interface{} {
	v, _ := ToStringMapE(i)
	return v
}

// ToStringMapE casts an interface to a map[string]interface{} type.
func ToStringMapE(i interface{}) (map[string]interface{}, error) {
	var m = make(map[string]interface{})
	if i == nil {
		return m, nil
	}

	ty := reflect.TypeOf(i).Kind()
	if ty == reflect.Map {
		vofi := reflect.ValueOf(i)
		for _, k := range vofi.MapKeys() {
			key := cast.ToString(k.Interface())
			value := vofi.MapIndex(k)

			m[key] = value.Interface()
		}

		return m, nil
	}

	err := LoadData(&m, i)
	return m, err
}

func MergeQueryString(a string, b string) string {
	qa, _ := url.ParseQuery(a)
	qb, _ := url.ParseQuery(b)

	for kb, _ := range qb {
		if qa.Get(kb) != "" {
			qa.Add(kb, qb.Get(kb))
		}

		qa.Set(kb, qb.Get(kb))
	}

	return qa.Encode()
}

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
