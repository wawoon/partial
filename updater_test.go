package partial_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/wawoon/partial"
)

type dummyStruct struct {
	A int
	B string
	C bool
	D *int
	E *string
	F *bool
}

type dummyStructUpdate struct {
	A *int
	B *string
	C *bool
	D *int
	E *string
	F *bool
}

func TestNewUpdater(t *testing.T) {
	successCases := []interface{}{
		&dummyStruct{},
		&dummyStruct{A: 1},
		&dummyStruct{A: 1, B: "b"},
		&dummyStruct{A: 1, B: "b", C: true},
		&dummyStruct{A: 1, B: "b", C: true, D: new(int)},
		&dummyStruct{A: 1, B: "b", C: true, D: new(int), E: new(string)},
		&dummyStruct{A: 1, B: "b", C: true, D: new(int), E: new(string), F: new(bool)},
	}

	for _, successCase := range successCases {
		updater, err := partial.NewUpdater(successCase)
		assert.NoError(t, err, "Error creating updater")
		assert.NotNil(t, updater, "Updater should not be nil")
	}

	failureCases := []interface{}{
		nil,
		1,
		"",
		[]int{1, 2, 3},
		[]string{"a", "b", "c"},
		[]bool{true, false},
		[]dummyStruct{},
		dummyStruct{},
	}

	for _, failureCase := range failureCases {
		updater, err := partial.NewUpdater(failureCase)
		assert.Error(t, err, "Error creating updater")
		assert.Nil(t, updater, "Updater should be nil")
	}
}

func TestUpdateWithStruct(t *testing.T) {
	val := dummyStruct{A: 1, B: "b", C: true}
	updater, err := partial.NewUpdater(&val)
	assert.NoError(t, err, "Error creating updater")
	assert.NotNil(t, updater, "Updater should not be nil")

	d := 3
	e := "e"
	f := true

	input := dummyStruct{
		A: 2,
		B: "b2",
		C: false,
		D: &d,
		E: &e,
		F: &f,
	}
	err = updater.Update(input)

	assert.NoError(t, err, "Error updating struct")
	assert.NotEmpty(t, updater.UpdatedFields, "Updated fields should not be empty")
	assert.Empty(t, updater.NotFoundFields, "Not found fields should be empty")
	assert.Empty(t, updater.NotAssignableFields, "Not assignable fields should be empty")

	assert.Equal(t, val, input, "Struct should be equal to input")
}

func TestUpdateWithPtrFields(t *testing.T) {
	val := dummyStruct{A: 1, B: "b", C: true}
	updater, err := partial.NewUpdater(&val)
	assert.NoError(t, err, "Error creating updater")
	assert.NotNil(t, updater, "Updater should not be nil")

	a := 2
	b := "b2"
	c := false
	d := 3
	e := "e"
	f := true

	input := dummyStructUpdate{
		A: &a,
		B: &b,
		C: &c,
		D: &d,
		E: &e,
		F: &f,
	}
	err = updater.Update(input)

	assert.NoError(t, err, "Error updating struct")
	assert.NotEmpty(t, updater.UpdatedFields, "Updated fields should not be empty")
	assert.Empty(t, updater.NotFoundFields, "Not found fields should be empty")
	assert.Empty(t, updater.NotAssignableFields, "Not assignable fields should be empty")

	assert.Equal(t, a, val.A, "a should be 2")
	assert.Equal(t, b, val.B, "b should be b2")
	assert.Equal(t, c, val.C, "c should be false")
	assert.Equal(t, d, *val.D, "d should be 3")
	assert.Equal(t, e, *val.E, "e should be e")
	assert.Equal(t, f, *val.F, "f should be true")

	newA := 5
	newE := "e3"
	err = updater.Update(dummyStructUpdate{
		A: &newA,
		E: &newE,
	})

	assert.NotEmpty(t, updater.UpdatedFields, "Updated fields should not be empty")
	assert.NotEmpty(t, updater.SkippedFields, "Skipped fields should not be empty")
	assert.Empty(t, updater.NotFoundFields, "Not found fields should be empty")
	assert.Empty(t, updater.NotAssignableFields, "Not assignable fields should be empty")

	assert.NoError(t, err, "Error updating struct")
	assert.Equal(t, newA, val.A, "a should be 0")
	assert.Equal(t, b, val.B, "b should be b2")
	assert.Equal(t, c, val.C, "c should be false")
	assert.Equal(t, d, *val.D, "d should be 3")
	assert.Equal(t, newE, *val.E, "e should be e3")
	assert.Equal(t, f, *val.F, "f should be true")
}

type StructWithAdvancedField struct {
	Time time.Time
	Strs []string
	Map  map[string]string
}

type StructWithAdvancedFieldUpdate struct {
	Time *time.Time
	Strs *[]string
	Map  *map[string]string
}

func TestUpdateWithAdvancedField(t *testing.T) {
	val := StructWithAdvancedField{}
	updater, err := partial.NewUpdater(&val)
	assert.NoError(t, err, "Error creating updater")

	input := StructWithAdvancedField{
		Time: time.Now(),
		Strs: []string{"a", "b", "c"},
		Map:  map[string]string{"a": "b", "c": "d"},
	}
	err = updater.Update(input)
	assert.NoError(t, err, "Error updating struct")
	assert.Equal(t, input, val, "Struct should be equal to input")

	timeVal := time.Now()
	strs := []string{"e", "f", "g"}
	mapVal := map[string]string{"e": "f", "g": "h"}
	input2 := StructWithAdvancedFieldUpdate{
		Time: &timeVal,
		Strs: &strs,
		Map:  &mapVal,
	}

	err = updater.Update(input2)
	assert.NoError(t, err, "Error updating struct")
	assert.Equal(t, *input2.Time, val.Time, "Time should be equal")
	assert.Equal(t, *input2.Strs, val.Strs, "Time should be equal")
	assert.Equal(t, *input2.Map, val.Map, "Time should be equal")
}
