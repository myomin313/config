package utils

import (
	"reflect"
	"strings"
)

var ignoreAuditFields = map[string]bool{
	"id":         true,
	"created_at": true,
	"updated_at": true,
	"created_by": true,
	"updated_by": true,
}

func StructChanges(oldObj, newObj interface{}) (map[string]interface{}, map[string]interface{}) {
	oldMap := make(map[string]interface{})
	newMap := make(map[string]interface{})

	diffValue(
		reflect.ValueOf(oldObj),
		reflect.ValueOf(newObj),
		oldMap,
		newMap,
	)

	return oldMap, newMap
}

func diffValue(oldVal, newVal reflect.Value, oldOut, newOut map[string]interface{}) {

	if !oldVal.IsValid() || !newVal.IsValid() {
		return
	}

	if oldVal.Kind() == reflect.Ptr {
		if oldVal.IsNil() {
			return
		}
		oldVal = oldVal.Elem()
	}

	if newVal.Kind() == reflect.Ptr {
		if newVal.IsNil() {
			return
		}
		newVal = newVal.Elem()
	}

	if oldVal.Kind() != newVal.Kind() {
		return
	}

	switch oldVal.Kind() {

	case reflect.Struct:

		t := oldVal.Type()

		for i := 0; i < oldVal.NumField(); i++ {

			field := t.Field(i)

			if !field.IsExported() {
				continue
			}

			name := field.Tag.Get("json")
			if name == "" {
				name = field.Name
			} else {
				name = strings.Split(name, ",")[0]
				if name == "" || name == "-" {
					name = field.Name
				}
			}

			if ignoreAuditFields[name] {
				continue
			}

			oldField := oldVal.Field(i)
			newField := newVal.Field(i)

			switch oldField.Kind() {

			case reflect.Struct:

				childOld := map[string]interface{}{}
				childNew := map[string]interface{}{}

				diffValue(oldField, newField, childOld, childNew)

				if len(childOld) > 0 {
					oldOut[name] = childOld
					newOut[name] = childNew
				}

			case reflect.Map:

				if !reflect.DeepEqual(oldField.Interface(), newField.Interface()) {
					childOld, childNew := diffMap(oldField, newField)

					oldOut[name] = childOld
					newOut[name] = childNew
				}

			case reflect.Slice, reflect.Array:

				if !reflect.DeepEqual(oldField.Interface(), newField.Interface()) {
					oldOut[name] = oldField.Interface()
					newOut[name] = newField.Interface()
				}

			default:

				if !reflect.DeepEqual(oldField.Interface(), newField.Interface()) {
					oldOut[name] = oldField.Interface()
					newOut[name] = newField.Interface()
				}
			}
		}

	case reflect.Map:

		childOld, childNew := diffMap(oldVal, newVal)

		for k, v := range childOld {
			oldOut[k] = v
		}

		for k, v := range childNew {
			newOut[k] = v
		}
	}
}

func diffMap(oldMapVal, newMapVal reflect.Value) (map[string]interface{}, map[string]interface{}) {

	oldResult := make(map[string]interface{})
	newResult := make(map[string]interface{})

	keys := map[string]struct{}{}

	for _, k := range oldMapVal.MapKeys() {
		keys[k.String()] = struct{}{}
	}

	for _, k := range newMapVal.MapKeys() {
		keys[k.String()] = struct{}{}
	}

	for key := range keys {

		oldV := oldMapVal.MapIndex(reflect.ValueOf(key))
		newV := newMapVal.MapIndex(reflect.ValueOf(key))

		if !oldV.IsValid() {
			newResult[key] = newV.Interface()
			continue
		}

		if !newV.IsValid() {
			oldResult[key] = oldV.Interface()
			continue
		}

		if reflect.DeepEqual(oldV.Interface(), newV.Interface()) {
			continue
		}

		if oldV.Kind() == reflect.Map && newV.Kind() == reflect.Map {

			childOld, childNew := diffMap(oldV, newV)

			oldResult[key] = childOld
			newResult[key] = childNew

			continue
		}

		oldResult[key] = oldV.Interface()
		newResult[key] = newV.Interface()
	}

	return oldResult, newResult
}
