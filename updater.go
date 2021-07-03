package partial

import (
	"errors"
	"fmt"
	"reflect"
)

type Updater struct {
	original            interface{}
	UpdatedFields       map[string]Field
	SkippedFields       map[string]Field
	NotFoundFields      map[string]Field
	NotAssignableFields map[string]Field
}

type Field struct {
	StructField reflect.StructField
	Value       reflect.Value
}

var ErrNonStructPtr error = errors.New("the given object is not a struct pointer")
var ErrNilPtr error = errors.New("the given object is nil")
var ErrUpdateFieldsFailure error = errors.New("update fields failure")

func NewUpdater(original interface{}) (*Updater, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Unexpected panic in NewUpdater: args %+v\n", original)
			return
		}
	}()

	if original == nil {
		return nil, ErrNilPtr
	}

	if reflect.TypeOf(original).Kind() != reflect.Ptr {
		return nil, ErrNonStructPtr
	}

	ptr := reflect.ValueOf(original).Elem()
	if ptr.Type().Kind() != reflect.Struct {
		return nil, ErrNonStructPtr
	}

	return &Updater{original: original}, nil
}

func ShouldSkip(structField reflect.StructField, value reflect.Value) bool {
	return !value.IsValid()
}

func (u *Updater) Update(newValue interface{}) error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Unexpected panic in Update: args %+v, panic: %+v\n", newValue, r)
			return
		}
	}()

	if newValue == nil {
		return ErrNilPtr
	}

	inputValue := reflect.ValueOf(newValue)

	if inputValue.Kind() == reflect.Ptr {
		inputValue = inputValue.Elem()
	}

	if inputValue.Type().Kind() != reflect.Struct {
		return ErrNonStructPtr
	}

	if reflect.TypeOf(u.original).Kind() != reflect.Ptr {
		return ErrNonStructPtr
	}

	originalValue := reflect.ValueOf(u.original).Elem()
	if !originalValue.IsValid() {
		return ErrNilPtr
	}

	updatedFields := make(map[string]Field)
	skippedFields := make(map[string]Field)
	notFoundFields := make(map[string]Field)
	notAssignableFields := make(map[string]Field)

	for i := 0; i < inputValue.NumField(); i++ {
		inputFieldType := inputValue.Type().Field(i)
		inputFieldValue := inputValue.Field(i)

		if ShouldSkip(inputFieldType, inputFieldValue) {
			skippedFields[inputFieldType.Name] = Field{
				StructField: inputFieldType,
				Value:       inputFieldValue,
			}
			continue
		}

		outputFieldValue := originalValue.FieldByName(inputFieldType.Name)

		if !outputFieldValue.IsValid() {
			notFoundFields[inputFieldType.Name] = Field{
				StructField: inputFieldType,
				Value:       inputFieldValue,
			}

			continue
		}

		// If the field is assignable without a type conversion,
		// we just assign the value.
		if inputFieldValue.Type().AssignableTo(outputFieldValue.Type()) {
			outputFieldValue.Set(inputFieldValue)

			updatedFields[inputFieldType.Name] = Field{
				StructField: inputFieldType,
				Value:       inputFieldValue,
			}
			continue
		}

		// if the input field is a pointer, we can't assign it to the output field
		// so we need to check if the output field is a pointer and if the input field
		// is assignable to the output field
		if inputFieldValue.Type().Kind() == reflect.Ptr &&
			inputFieldValue.Elem().Type().AssignableTo(outputFieldValue.Type()) {

			outputFieldValue.Set(inputFieldValue.Elem())
			updatedFields[inputFieldType.Name] = Field{
				StructField: inputFieldType,
				Value:       inputFieldValue.Elem(),
			}
			continue
		}

		// If we cannot assign the input field to the output field,
		// we keep this field in the notAssignableFields map
		notAssignableFields[inputFieldType.Name] = Field{
			StructField: inputFieldType,
			Value:       inputFieldValue,
		}
	}

	u.UpdatedFields = updatedFields
	u.NotFoundFields = notFoundFields
	u.NotAssignableFields = notAssignableFields

	if len(u.NotFoundFields) > 0 || len(u.NotAssignableFields) > 0 {
		return ErrUpdateFieldsFailure
	}

	return nil
}
