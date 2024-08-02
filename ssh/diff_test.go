package ssh

import (
	"reflect"
	"testing"
)

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
