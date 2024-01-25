package utils

import (
	"reflect"
	"time"
)

// ReflectiveUpdate
//   - This function updates a model with the values of another model
//     if the value of the other model is not the zero value of its type
func ReflectiveUpdate(target interface{}, source interface{}) {
	targetValue := reflect.ValueOf(target).Elem()
	sourceValue := reflect.ValueOf(source).Elem()

	for i := 0; i < targetValue.NumField(); i++ {
		targetField := targetValue.Field(i)
		sourceField := sourceValue.Field(i)

		if !reflect.DeepEqual(sourceField.Interface(), reflect.Zero(sourceField.Type()).Interface()) ||
			(sourceField.Type() == reflect.TypeOf(time.Time{}) && !sourceField.Interface().(time.Time).IsZero()) {
			targetField.Set(sourceField)
		}
	}
}
