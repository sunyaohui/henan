package slice

import (
	"reflect"
)

//去重
func SliceRemoveDuplicate(a interface{}) (ret []interface{}) {
	if reflect.TypeOf(a).Kind() != reflect.Slice {
		return ret
	}

	va := reflect.ValueOf(a)
	for i := 0; i < va.Len(); i++ {
		if i > 0 && reflect.DeepEqual(va.Index(i-1).Interface(), va.Index(i).Interface()) {
			continue
		}
		ret = append(ret, va.Index(i).Interface())
	}
	return ret
}
