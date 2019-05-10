package struql

import (
	"errors"
	"reflect"
	"strings"
)

// Field represents struct field in the table
type Field struct {
	Name  string
	Value interface{}

	idx  int
	kind reflect.Kind
}

// Index returns field index in the row
func (f Field) Index() int {
	return f.idx
}

func (f Field) passModifier(mod ValueModifier) interface{} {
	if mod != nil {
		return mod(f.Value)
	}
	return f.Value
}

func (f Field) compareGreater(filter *Filter) (bool, error) {
	switch f.kind {
	case reflect.String:
		return f.passModifier(filter.Modifier).(string) > filter.Value.(string), nil
	case reflect.Int:
		return f.passModifier(filter.Modifier).(int) > filter.Value.(int), nil
	case reflect.Float32:
		return f.passModifier(filter.Modifier).(float32) > filter.Value.(float32), nil
	case reflect.Float64:
		return f.passModifier(filter.Modifier).(float64) > filter.Value.(float64), nil
	case reflect.Int32:
		return f.passModifier(filter.Modifier).(int32) > filter.Value.(int32), nil
	case reflect.Int64:
		return f.passModifier(filter.Modifier).(int64) > filter.Value.(int64), nil
	}
	return false, errors.New(errUnsuppotredCompare)
}

func (f Field) compareLesser(filter *Filter) (bool, error) {
	switch f.kind {
	case reflect.String:
		return f.passModifier(filter.Modifier).(string) < filter.Value.(string), nil
	case reflect.Int:
		return f.passModifier(filter.Modifier).(int) < filter.Value.(int), nil
	case reflect.Float32:
		return f.passModifier(filter.Modifier).(float32) < filter.Value.(float32), nil
	case reflect.Float64:
		return f.passModifier(filter.Modifier).(float64) < filter.Value.(float64), nil
	case reflect.Int32:
		return f.passModifier(filter.Modifier).(int32) < filter.Value.(int32), nil
	case reflect.Int64:
		return f.passModifier(filter.Modifier).(int64) < filter.Value.(int64), nil
	}
	return false, errors.New(errUnsuppotredCompare)
}

func (f Field) compareEqual(filter *Filter) (bool, error) {
	return f.passModifier(filter.Modifier) == filter.Value, nil
}

func (f Field) compareBeginWith(filter *Filter) (bool, error) {
	if f.kind == reflect.String {
		return strings.HasPrefix(f.passModifier(filter.Modifier).(string), filter.Value.(string)), nil
	}
	return false, errors.New(errUnsuppotredCompare)
}

func (f Field) compareEndWith(filter *Filter) (bool, error) {
	if f.kind == reflect.String {
		return strings.HasSuffix(f.passModifier(filter.Modifier).(string), filter.Value.(string)), nil
	}
	return false, errors.New(errUnsuppotredCompare)
}

func (f Field) compareIn(filter *Filter) (bool, error) {
	if f.kind == reflect.Slice {
		fieldValue := reflect.ValueOf(f.Value)
		for j := 0; j < fieldValue.Len(); j++ {
			if fieldValue.Index(j).Interface() == filter.Value {
				return true, nil
			}
		}
	} else {
		return false, errors.New(errUnsuppotredCompare)
	}
	return false, nil
}

func (f Field) compare(filter *Filter) (bool, error) {
	switch filter.Operation {
	case ComparisonEqual:
		return f.compareEqual(filter)
	case ComparisonGreater:
		return f.compareGreater(filter)
	case ComparisonLesser:
		return f.compareLesser(filter)
	case ComparisonBeginWith:
		return f.compareBeginWith(filter)
	case ComparisonEndWith:
		return f.compareEndWith(filter)
	case ComparisonIn:
		return f.compareIn(filter)
	}
	return false, errors.New(errUnsuppotredCompare)
}
