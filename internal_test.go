package partial

import (
    "reflect"
    "testing"

    "github.com/stretchr/testify/assert"
)

// TestFindFieldByInsensitiveName verifies that the helper correctly finds fields
// regardless of the case of the supplied name.
func TestFindFieldByInsensitiveName(t *testing.T) {
    type sample struct {
        Field1       int
        AnotherField string
    }

    s := sample{Field1: 10, AnotherField: "hello"}

    v := reflect.ValueOf(s)

    tests := []struct {
        name       string
        searchName string
        wantValid  bool
    }{
        {"exact match", "Field1", true},
        {"lowercase", "field1", true},
        {"uppercase", "FIELD1", true},
        {"mixed case", "AnOtHeRfIeLd", true},
        {"not found", "NotExists", false},
    }

    for _, tt := range tests {
        got := findFieldByInsensitiveName(v, tt.searchName)
        assert.Equalf(t, tt.wantValid, got.IsValid(), "case %q", tt.name)
        if tt.wantValid {
            // Additional sanity check – make sure we locate the correct field by comparing its kind.
            switch tt.searchName {
            case "Field1", "field1", "FIELD1":
                assert.Equal(t, reflect.Int, got.Kind())
            case "AnotherField", "AnOtHeRfIeLd":
                assert.Equal(t, reflect.String, got.Kind())
            }
        }
    }
}

// TestShouldSkipFieldImpl ensures that pointer fields that are nil are skipped
// and others are processed.
func TestShouldSkipFieldImpl(t *testing.T) {
    type ptrStruct struct {
        P *int
        V int
    }

    var p *int // nil pointer
    s := ptrStruct{P: p, V: 0}

    // nil pointer case
    structFieldP, _ := reflect.TypeOf(s).FieldByName("P")
    valP := reflect.ValueOf(s).FieldByName("P")
    assert.True(t, shouldSkipFieldImpl(structFieldP, valP))

    // populated pointer case
    intVal := 5
    s2 := ptrStruct{P: &intVal, V: 0}
    structFieldP2, _ := reflect.TypeOf(s2).FieldByName("P")
    valP2 := reflect.ValueOf(s2).FieldByName("P")
    assert.False(t, shouldSkipFieldImpl(structFieldP2, valP2))

    // non-pointer value should not be skipped (IsValid is true)
    structFieldV, _ := reflect.TypeOf(s2).FieldByName("V")
    valV := reflect.ValueOf(s2).FieldByName("V")
    assert.False(t, shouldSkipFieldImpl(structFieldV, valV))
}
