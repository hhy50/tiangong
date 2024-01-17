package common

import (
	"reflect"
)

// GetAllFields
func GetAllFields(o interface{}) map[string]reflect.StructField {
	fields := make(map[string]reflect.StructField)
	t := reflect.TypeOf(o)
	switch t.Kind() {
	case reflect.Ptr:
		t = t.Elem()
	default:

	}

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		fields[f.Name] = f
	}
	return fields
}

// GetPtr
func GetPtr(o interface{}) (*reflect.Value, bool) {
	v := reflect.ValueOf(o)
	switch v.Kind() {
	case reflect.Ptr:
		return &v, true
	default:
		return nil, false
	}
}

// GetFieldsTag Returns a field that declares the specified tag
// key is field Name,value is tagValue
func GetTags(tagName string, o interface{}) map[string]string {
	fields := make(map[string]string)
	t := reflect.TypeOf(o)
	switch t.Kind() {
	case reflect.Ptr:
		t = t.Elem()
	default:

	}

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		val := f.Tag.Get(tagName)
		if IsNotEmpty(val) {
			fields[f.Name] = val
		}
	}
	return fields
}
