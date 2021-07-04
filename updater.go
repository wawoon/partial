package partial

import (
	"errors"
	"reflect"
	"strings"
)

type Updater struct {
	output          interface{}
	FindFieldByName func(reflect.Value, string) reflect.Value
	ShouldSkipField func(structField reflect.StructField, value reflect.Value) bool

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

// NewUpdater will returns Updater object
func NewUpdater(output interface{}) (*Updater, error) {
	if output == nil {
		return nil, ErrNilPtr
	}

	if reflect.TypeOf(output).Kind() != reflect.Ptr {
		return nil, ErrNonStructPtr
	}

	ptr := reflect.ValueOf(output).Elem()
	if ptr.Type().Kind() != reflect.Struct {
		return nil, ErrNonStructPtr
	}

	return &Updater{
		output:          output,
		FindFieldByName: findFieldByInsensitiveName,
		ShouldSkipField: shouldSkipFieldImpl,
	}, nil
}

func shouldSkipFieldImpl(structField reflect.StructField, value reflect.Value) bool {
	if structField.Type.Kind() == reflect.Ptr {
		if value.IsNil() {
			return true
		}
	}
	return !value.IsValid()
}

// findFieldByInsensitiveName finds a reflect.Value field by name in a non-case-sensitive way
func findFieldByInsensitiveName(field reflect.Value, name string) reflect.Value {
	inputFieldLowerCaseName := strings.ToLower(name)
	foundField := field.FieldByNameFunc(func(s string) bool {
		return strings.ToLower(s) == inputFieldLowerCaseName
	})
	return foundField
}

func (u *Updater) Update(newValue interface{}) error {
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

	if reflect.TypeOf(u.output).Kind() != reflect.Ptr {
		return ErrNonStructPtr
	}

	outputValue := reflect.ValueOf(u.output).Elem()
	if !outputValue.IsValid() {
		return ErrNilPtr
	}

	updatedFields := make(map[string]Field)
	skippedFields := make(map[string]Field)
	notFoundFields := make(map[string]Field)
	notAssignableFields := make(map[string]Field)

	for i := 0; i < inputValue.NumField(); i++ {
		inputFieldType := inputValue.Type().Field(i)
		inputFieldValue := inputValue.Field(i)

		if u.ShouldSkipField(inputFieldType, inputFieldValue) {
			skippedFields[inputFieldType.Name] = Field{
				StructField: inputFieldType,
				Value:       inputFieldValue,
			}
			continue
		}

		// The assign-target field is searched in a non-case-sensitive way.
		outputFieldValue := findFieldByInsensitiveName(outputValue, inputFieldType.Name)

		// If the field is not found, it is added to the not-found-fields list.
		if !outputFieldValue.IsValid() {
			notFoundFields[inputFieldType.Name] = Field{
				StructField: inputFieldType,
				Value:       inputFieldValue,
			}

			continue
		}

		// There are 3 possible cases:
		// 1. inputFieldValue is non-pointer (e.g. int), and outputFieldValue is non-pointer
		// 2. inputFieldValue is pointer (e.g. *int), and outputFieldValue is non-pointer
		// 3. inputFieldValue is pointer (e.g. *int), and outputFieldValue is pointer (e.g. *int)

		// input:non-pointer -> output:non-pointer case
		if inputFieldValue.Kind() != reflect.Ptr {
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
		}

		if inputFieldValue.Type().Kind() == reflect.Ptr {
			if outputFieldValue.Type().Kind() != reflect.Ptr {
				// input:pointer -> output:non-pointer case
				if !inputFieldValue.IsNil() && inputFieldValue.Elem().Type().AssignableTo(outputFieldValue.Type()) {
					outputFieldValue.Set(inputFieldValue.Elem())
					updatedFields[inputFieldType.Name] = Field{
						StructField: inputFieldType,
						Value:       inputFieldValue.Elem(),
					}
					continue
				}
			} else {
				// input:pointer -> ouput:pointer case
				if inputFieldType.Type.AssignableTo(outputFieldValue.Type()) {
					outputFieldValue.Set(inputFieldValue)
					updatedFields[inputFieldType.Name] = Field{
						StructField: inputFieldType,
						Value:       inputFieldValue.Elem(),
					}
					continue
				}
			}
		}

		// If we cannot assign the input field to the output field,
		// we keep this field in the notAssignableFields map
		notAssignableFields[inputFieldType.Name] = Field{
			StructField: inputFieldType,
			Value:       inputFieldValue,
		}
	}

	u.UpdatedFields = updatedFields
	u.SkippedFields = skippedFields
	u.NotFoundFields = notFoundFields
	u.NotAssignableFields = notAssignableFields

	if len(u.NotFoundFields) > 0 || len(u.NotAssignableFields) > 0 {
		return ErrUpdateFieldsFailure
	}

	return nil
}
