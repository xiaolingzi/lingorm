package typehelper

import "reflect"

func InterfaceToSlice(slice interface{}) []interface{} {
	val := reflect.ValueOf(slice)
	if val.Kind() != reflect.Slice {
		panic("Invalid argument for slice type")
	}

	result := make([]interface{}, val.Len())

	for i := 0; i < val.Len(); i++ {
		result[i] = val.Index(i).Interface()
	}

	return result
}
