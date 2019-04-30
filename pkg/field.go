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

func (f *Field) compare(Value interface{}, Operation int) (bool, error) {
	if Operation == ComparsionEqual {
		return f.Value == Value, nil
	}

	switch Operation {

	case ComparsionGreater:
		switch f.kind {
		case reflect.String:
			return f.Value.(string) > Value.(string), nil
		case reflect.Int:
			return f.Value.(int) > Value.(int), nil
		case reflect.Float32:
			return f.Value.(float32) > Value.(float32), nil
		case reflect.Float64:
			return f.Value.(float64) > Value.(float64), nil
		case reflect.Int32:
			return f.Value.(int32) > Value.(int32), nil
		case reflect.Int64:
			return f.Value.(int64) > Value.(int64), nil
		default:
			return false, errors.New("unsupported comparsion")
		}

	case ComparsionLesser:
		switch f.kind {
		case reflect.String:
			return f.Value.(string) < Value.(string), nil
		case reflect.Int:
			return f.Value.(int) < Value.(int), nil
		case reflect.Float32:
			return f.Value.(float32) < Value.(float32), nil
		case reflect.Float64:
			return f.Value.(float64) < Value.(float64), nil
		case reflect.Int32:
			return f.Value.(int32) < Value.(int32), nil
		case reflect.Int64:
			return f.Value.(int64) < Value.(int64), nil
		default:
			return false, errors.New("unsupported comparsion")
		}

	case ComparsionBeginWith:
		switch f.kind {
		case reflect.String:
			return strings.HasPrefix(f.Value.(string), Value.(string)), nil
		default:
			return false, errors.New("unsupported comparsion")
		}

	case ComparsionEndWith:
		switch f.kind {
		case reflect.String:
			return strings.HasSuffix(f.Value.(string), Value.(string)), nil
		default:
			return false, errors.New("unsupported comparsion")
		}

	case ComparsionIn:
		switch f.kind {
		case reflect.Slice:
			fieldValue := reflect.ValueOf(f.Value)
			for j := 0; j < fieldValue.Len(); j++ {
				if fieldValue.Index(j).Interface() == Value {
					return true, nil
				}
			}
		default:
			return false, errors.New("unsupported comparsion")
		}

	}
	return false, nil
}
