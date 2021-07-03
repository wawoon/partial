package partial

import (
	"errors"
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
	if structField.Type.Kind() == reflect.Ptr {
		if value.IsNil() {
			return true
		}
	}
	return !value.IsValid()
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
