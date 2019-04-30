package struql

import "errors"

// StringModifier - callback for transforming field value before compare
type StringModifier func(Value string) string

// Filter ...
type Filter struct {
	FieldName string
	Value     interface{}
	Operation int
	Modifier  StringModifier

	// fieldIndex - to improve performance
	fieldIndex int
	validated  bool
}

// Validate check field and get fieldIndex from row
func (f *Filter) Validate(r RowCollection) error {
	if f.validated {
		return nil
	}
	if len(r) == 0 {
		return errors.New("empty dataset")
	}
	field := r[0].FieldByName(f.FieldName)
	if field == nil {
		return errors.New("no such field in dataset")
	}
	f.fieldIndex = field.Index()
	f.validated = true
	return nil
}
