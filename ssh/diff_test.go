package ssh_test

import (
	"fmt"
	"reflect"
	"testing"
)

type FieldDiff struct {
	FieldName string
	LhsValue  string
	RhsValue  string
}

type TypeMismatchError struct {
	Text string
}

func (err *TypeMismatchError) Error() string {
	return err.Text
}

type diff struct{}

var Diff diff

func (diff) Structs(lhs any, rhs any) ([]FieldDiff, error) {
	lhsUnwrapped := reflect.ValueOf(lhs).Elem()
	rhsUnwrapped := reflect.ValueOf(rhs).Elem()

	lhsType := reflect.TypeOf(lhs)
	rhsType := reflect.TypeOf(rhs)

	if lhsType != rhsType {
		return []FieldDiff{}, &TypeMismatchError{
			Text: fmt.Sprintf("Wrong types - comparing %v to %v", lhsType, rhsType),
		}
	}

	fieldDiffs := []FieldDiff{}

	for _, field := range reflect.VisibleFields(lhsUnwrapped.Type()) {
		lhsField := lhsUnwrapped.FieldByName(field.Name)
		rhsField := rhsUnwrapped.FieldByName(field.Name)

		if lhsField.String() != rhsField.String() {
			fieldDiffs = append(fieldDiffs, FieldDiff{
				FieldName: field.Name,
				LhsValue:  lhsField.String(),
				RhsValue:  rhsField.String(),
			})
		}
	}

	return fieldDiffs, nil
}

type TestStruct1 struct {
	Field1 string
	Field2 string
	Field3 string
}

type TestStruct2 struct {
	Field4 string
	Field5 string
	Field6 string
}

func TestDiffStructs(t *testing.T) {
	t.Run("Type mismatch", func(t *testing.T) {
		struct1 := &TestStruct1{}
		struct2 := &TestStruct2{}

		_, err := Diff.Structs(struct1, struct2)
		if err == nil {
			t.Fatal("Expected an error, got nothing")
		}
	})

	t.Run("Two structs match", func(t *testing.T) {
		struct1 := &TestStruct1{Field1: "a", Field2: "b", Field3: "c"}
		struct2 := &TestStruct1{Field1: "a", Field2: "b", Field3: "c"}

		diffs, _ := Diff.Structs(struct1, struct2)
		if len(diffs) != 0 {
			t.Fatalf("Expected no diffs, got %v", diffs)
		}
	})

	t.Run("Two structs don't match", func(t *testing.T) {
		struct1 := &TestStruct1{Field1: "a", Field2: "b", Field3: "c"}
		struct2 := &TestStruct1{Field1: "a", Field2: "d", Field3: "c"}

		diffsGot, _ := Diff.Structs(struct1, struct2)
		diffsWant := []FieldDiff{{
			FieldName: "Field2",
			LhsValue:  "b",
			RhsValue:  "d",
		}}

		if !reflect.DeepEqual(diffsGot, diffsWant) {
			t.Fatalf("\nwant\t%v\ngot\t\t%v", diffsWant, diffsGot)
		}
	})
}
