package struql

import (
	"errors"
	"fmt"
	"reflect"
)

const (
	errEmptyDataSet       = "empty dataset"
	errNoSuchField        = "no such field in dataset"
	errUnsuppotredCompare = "unsupported comparison"
)

// Row - collection of fields / table row
type Row struct {
	fieldMap  map[string]*Field
	fieldList []*Field
}

// RowCollection - collection of rows / table
type RowCollection []*Row

// NewRow creates new initialized row
func newRow() *Row {
	row := &Row{}
	row.init()
	return row
}

// Init initializes row, allocating memory for fields
func (r *Row) init() {
	r.fieldMap = make(map[string]*Field)
	r.fieldList = make([]*Field, 0, 2)
}

// FieldByName returns field by it's name or nil
func (r *Row) FieldByName(name string) *Field {
	if f, ok := r.fieldMap[name]; ok {
		return f
	}
	return nil
}

// FieldByIndex returns field by it's index or nil
func (r *Row) FieldByIndex(index int) *Field {
	if index > len(r.fieldList) {
		return nil
	}
	return r.fieldList[index]
}

// AddField add new field to row
func (r *Row) AddField(name string, value interface{}) {
	if _, ok := r.fieldMap[name]; ok {
		return
	}

	newField := &Field{
		Name:  name,
		Value: value,
		idx:   len(r.fieldList),
		kind:  reflect.ValueOf(value).Kind(),
	}
	r.fieldMap[name] = newField
	r.fieldList = append(r.fieldList, newField)
}

// Where collects data in the rows according to filters
func (r RowCollection) Where(result RowCollection, filters ...Filter) (RowCollection, error) {
	var (
		ok         int
		filtersLen int
		field      *Field
		err        error
		compareOk  bool
	)
	filtersLen = len(filters)

	for idx := 0; idx < len(filters); idx++ {
		err = filters[idx].validate(r)
		if err != nil {
			return nil, err
		}
	}

	for _, row := range r {
		ok = 0
		for idx := 0; idx < len(filters); idx++ {
			field = row.FieldByIndex(filters[idx].fieldIndex)
			if compareOk, err = field.compare(&filters[idx]); err == nil {
				if compareOk {
					ok++
				}
			} else {
				return nil, err
			}
		}
		if ok == filtersLen {
			result = append(result, row)
		}
	}
	return result, nil
}

// AddField add field to all rows
func (r RowCollection) AddField(name string, value interface{}) {
	for _, row := range r {
		row.AddField(name, value)
	}
}

// FieldIndex returns field index in row collection
func (r RowCollection) fieldIndex(fieldName string) (int, error) {
	if len(r) == 0 {
		return -1, errors.New(errEmptyDataSet)
	}
	fld := r[0].FieldByName(fieldName)
	if fld == nil {
		return -1, errors.New(errNoSuchField)
	}
	return fld.Index(), nil
}

// CollectValues returns array of values of according field
func (r RowCollection) CollectValues(fieldName string) ([]interface{}, error) {
	idx, err := r.fieldIndex(fieldName)
	if err != nil {
		return nil, err
	}
	result := make([]interface{}, len(r))

	for i := 0; i < len(r); i++ {
		result[i] = r[i].FieldByIndex(idx).Value
	}
	return result, nil
}

// TODO: Select returns new row collection with specified fields
//func (r RowCollection) Select(fieldName ...string) RowCollection {
//	result := make(RowCollection, len(r))
//	return result
//}

// String represents rows as string
func (r RowCollection) String() string {
	var result string
	for i := 0; i < len(r); i++ {
		result += r[i].String() + "\n"
	}
	return result
}

// Count of rows
func (r RowCollection) Count() int {
	return len(r)
}

// String represents row as string
func (r Row) String() string {
	var result string
	for i := 0; i < len(r.fieldList); i++ {
		result += fmt.Sprintf("%s: [%v]\t", r.fieldList[i].Name, r.fieldList[i].Value)
	}
	return result
}
