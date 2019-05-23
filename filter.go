package struql

import (
	"errors"
)

const (
	// ComparisonEqual - =
	ComparisonEqual = iota
	// ComparisonIn - in [values]
	ComparisonIn
	// ComparisonGreater - >
	ComparisonGreater
	// ComparisonLesser - <
	ComparisonLesser
	// ComparisonBeginWith - like 'value%'
	ComparisonBeginWith
	// ComparisonEndWith - like '%value'
	ComparisonEndWith
	// ComparisonExists - value in[]
	ComparisonExists
)

// ValueModifier - callback for transforming field value before compare
type ValueModifier func(Value interface{}) interface{}

// Filter sets the search condition
type Filter struct {
	FieldName string
	Value     interface{}
	Operation int
	Modifier  ValueModifier

	// fieldIndex - to improve performance
	fieldIndex int
}

// Validate check field and get fieldIndex from row
func (f *Filter) validate(r RowCollection) error {
	if len(r) == 0 {
		return errors.New("empty dataset")
	}
	field := r[0].FieldByName(f.FieldName)
	if field == nil {
		return errors.New("no such field in dataset")
	}
	f.fieldIndex = field.Index()
	return nil
}
