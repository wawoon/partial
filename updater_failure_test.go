package partial_test

import (
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/wawoon/partial"
)

// dummyStruct already defined in updater_test.go, re-declare here for clarity.
type dummyStruct struct {
    A int
    B string
}

// updateFailureStruct is crafted so that:
//   1. Field "B" exists in the destination but has an incompatible type -> NotAssignableFields.
//   2. Field "Z" does not exist in the destination -> NotFoundFields.
type updateFailureStruct struct {
    B time.Time
    Z int
}

func TestUpdateFailureCases(t *testing.T) {
    original := dummyStruct{A: 1, B: "b"}

    updater, err := partial.NewUpdater(&original)
    assert.NoError(t, err)

    input := updateFailureStruct{
        B: time.Now(),
        Z: 100,
    }

    err = updater.Update(input)

    // We expect an error because there are not-found and not-assignable fields.
    assert.ErrorIs(t, err, partial.ErrUpdateFieldsFailure)

    // Validate internal bookkeeping.
    assert.Contains(t, updater.NotAssignableFields, "B")
    assert.Contains(t, updater.NotFoundFields, "Z")

    // Ensure the value was NOT updated.
    assert.Equal(t, "b", original.B)
    assert.NotContains(t, updater.UpdatedFields, "B")
}
