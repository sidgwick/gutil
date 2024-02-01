package generic

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	gjson "github.com/sidgwick/gutil/json"
	"github.com/spf13/cast"
)

var NotFoundError = errors.New("not found")

func GetByPath(data interface{}, path string, _default interface{}) (interface{}, error) {
	xPath := strings.Split(path, ".")

	result, err := _getByPath(data, xPath, _default)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func _getByPath(data interface{}, path []string, _default interface{}) (interface{}, error) {
	if len(path) == 0 {
		return data, nil
	}

	ps := path[0]
	path = path[1:]

	child, err := _getValueByKeyName(data, ps)
	if err != nil {
		return nil, err
	}

	if len(path) == 0 && child == nil {
		return _default, nil
	}

	if len(path) > 0 && child == nil {
		return nil, NotFoundError
	}

	return _getByPath(child, path, _default)
}

func _getValueByKeyName(data interface{}, key string) (interface{}, error) {
	isJsonStr, index, key := parsePathString(key)

	dataValue := reflect.ValueOf(data)
	dataKind := reflect.TypeOf(data).Kind()

	found := false
	if key != "" {
		for _, k := range dataValue.MapKeys() {
			if k.String() != key {
				continue
			}

			dataValue = dataValue.MapIndex(k)
			dataKind = dataValue.Kind()
			found = true
			break
		}
	} else if index > 0 {
		// 尝试在数组里面找
		found = true
	}

	if !found {
		return nil, nil
	}

	if isJsonStr {
		err := gjson.LoadData(&data, dataValue.Interface())
		if err != nil {
			return nil, err
		}

		dataValue = reflect.ValueOf(data)
		dataKind = reflect.TypeOf(data).Kind()
	}

	if index >= 0 {
		data = dataValue.Interface()
		dataValue = reflect.ValueOf(data)
		dataKind = reflect.TypeOf(data).Kind()

		if dataKind != reflect.Array && dataKind != reflect.Slice {
			msg := fmt.Sprintf("path=%v object data is not a array or slice", key)
			return nil, errors.New(msg)
		}

		maxIndex := dataValue.Len()
		if index >= maxIndex {
			msg := fmt.Sprintf("index=%v out of slice length=%v", key, maxIndex)
			return nil, errors.New(msg)
		}

		dataValue = dataValue.Index(index)
	}

	return dataValue.Interface(), nil
}

// key could be:
// 1. `[0]` --  数组 index
// 2. `abc` --  哈希表的键等于 `abc`
// 3. `JSON:abc` -- json 字符串形式编码的哈希表里面的键等于 `abc`
func parsePathString(ps string) (bool, int, string) {
	isJsonStr := false // 希望找的数据中有 JSON 编码的数据
	index := -1
	leftSquare := strings.Index(ps, "[")

	if leftSquare != -1 {
		index = cast.ToInt(ps[leftSquare+1 : len(ps)-1])
		ps = ps[0:leftSquare]
	}

	if strings.HasPrefix(ps, "JSON:") {
		ps = ps[5:]
		isJsonStr = true
	}

	return isJsonStr, index, ps
}
