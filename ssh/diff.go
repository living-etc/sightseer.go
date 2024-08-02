package ssh

import (
	"fmt"
	"reflect"
)

type FieldDiff struct {
	FieldName string
	LhsValue  string
	RhsValue  string
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
