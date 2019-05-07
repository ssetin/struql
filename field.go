package struql

import (
	"errors"
	"reflect"
	"strings"
)

const (
	// ComparsionEqual - =
	ComparsionEqual = iota
	// ComparsionGreater - >
	ComparsionGreater
	// ComparsionLesser - <
	ComparsionLesser
	// ComparsionBeginWith - like 'value%'
	ComparsionBeginWith
	// ComparsionEndWith - like '%value'
	ComparsionEndWith
	// ComparsionIn - value in[]
	ComparsionIn
)

// Field ...
type Field struct {
	Name  string
	Value interface{}

	idx  int
	kind reflect.Kind
}

// Index returns field index in the row
func (f *Field) Index() int {
	return f.idx
}

func (f *Field) passModifier(mod ValueModifier) interface{} {
	if mod != nil {
		return mod(f.Value)
	}
	return f.Value
}

func (f *Field) compare(filter *Filter) (bool, error) {
	if filter.Operation == ComparsionEqual {
		return f.passModifier(filter.Modifier) == filter.Value, nil
	}

	switch filter.Operation {

	case ComparsionGreater:
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
		default:
			return false, errors.New("unsupported comparsion")
		}

	case ComparsionLesser:
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
		default:
			return false, errors.New("unsupported comparsion")
		}

	case ComparsionBeginWith:
		switch f.kind {
		case reflect.String:
			return strings.HasPrefix(f.passModifier(filter.Modifier).(string), filter.Value.(string)), nil
		default:
			return false, errors.New("unsupported comparsion")
		}

	case ComparsionEndWith:
		switch f.kind {
		case reflect.String:
			return strings.HasSuffix(f.passModifier(filter.Modifier).(string), filter.Value.(string)), nil
		default:
			return false, errors.New("unsupported comparsion")
		}

	case ComparsionIn:
		switch f.kind {
		case reflect.Slice:
			fieldValue := reflect.ValueOf(f.Value)
			for j := 0; j < fieldValue.Len(); j++ {
				if fieldValue.Index(j).Interface() == filter.Value {
					return true, nil
				}
			}
		default:
			return false, errors.New("unsupported comparsion")
		}

	}
	return false, nil
}
