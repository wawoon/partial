package partial_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wawoon/partial"
)

type dummyStruct struct {
	A int
	B string
	C bool
}

type dummyStructUpdate struct {
	A *int
	B *string
	C *bool
}

func TestNewUpdater(t *testing.T) {
	successCases := []interface{}{
		&dummyStruct{},
		&dummyStruct{A: 1},
		&dummyStruct{A: 1, B: "b"},
		&dummyStruct{A: 1, B: "b", C: true},
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

	err = updater.Update(dummyStruct{
		A: 2,
		B: "b2",
		C: false,
	})

	assert.NoError(t, err, "Error updating struct")
	assert.NotEmpty(t, updater.UpdatedFields, "Updated fields should not be empty")
	assert.Empty(t, updater.NotFoundFields, "Not found fields should be empty")
	assert.Empty(t, updater.NotAssignableFields, "Not assignable fields should be empty")

	assert.Equal(t, 2, val.A, "a should be 2")
	assert.Equal(t, "b2", val.B, "b should be b2")
	assert.Equal(t, false, val.C, "c should be false")
}

func TestUpdateWithPtrFields(t *testing.T) {
	val := dummyStruct{A: 1, B: "b", C: true}
	updater, err := partial.NewUpdater(&val)
	assert.NoError(t, err, "Error creating updater")
	assert.NotNil(t, updater, "Updater should not be nil")

	a := 2
	b := "b2"
	c := false

	err = updater.Update(dummyStructUpdate{
		A: &a,
		B: &b,
		C: &c,
	})

	assert.NoError(t, err, "Error updating struct")
	assert.NotEmpty(t, updater.UpdatedFields, "Updated fields should not be empty")
	assert.Empty(t, updater.NotFoundFields, "Not found fields should be empty")
	assert.Empty(t, updater.NotAssignableFields, "Not assignable fields should be empty")

	assert.Equal(t, 2, val.A, "a should be 2")
	assert.Equal(t, "b2", val.B, "b should be b2")
	assert.Equal(t, false, val.C, "c should be false")
}
