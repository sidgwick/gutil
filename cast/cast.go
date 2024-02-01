package cast

import (
	"reflect"

	"github.com/sidgwick/gutil/json"
	"github.com/spf13/cast"
)

// ToXStringMap casts an interface to a map[string]interface{} type.
func ToXStringMap(i interface{}) map[string]interface{} {
	v, _ := ToXStringMapE(i)
	return v
}

// ToXStringMapE casts an interface to a map[string]interface{} type.
func ToXStringMapE(i interface{}) (map[string]interface{}, error) {
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

	err := json.LoadData(&m, i)
	return m, err
}
