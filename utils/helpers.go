package utils

import (
	"reflect"
)

func UpdateStruct(dst interface{}, src interface{}) {
	dstVal := reflect.ValueOf(dst)
	if dstVal.Kind() != reflect.Ptr || dstVal.IsNil() {
		return
	}

	dstVal = dstVal.Elem()
	srcVal := reflect.ValueOf(src).Elem()

	for i := 0; i < srcVal.NumField(); i++ {
		srcField := srcVal.Field(i)

		if srcField.Kind() == reflect.Ptr && !srcField.IsNil() {
			fieldName := srcVal.Type().Field(i).Name
			dstField := dstVal.FieldByName(fieldName)

			if dstField.IsValid() && dstField.CanSet() {
				dstField.Set(srcField.Elem())
			}
		}
	}
}
