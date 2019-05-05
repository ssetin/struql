package struql

import (
	"fmt"
	"reflect"
)

// Row ...
type Row struct {
	fieldMap  map[string]*Field
	fieldList []*Field
}

// RowCollection ...
type RowCollection []*Row

// NewRow ...
func NewRow() *Row {
	row := &Row{}
	row.Init()
	return row
}

// Init ...
func (r *Row) Init() {
	r.fieldMap = make(map[string]*Field)
	r.fieldList = make([]*Field, 0, 2)
}

// FieldByName ...
func (r *Row) FieldByName(name string) *Field {
	if f, ok := r.fieldMap[name]; ok {
		return f
	}
	return nil
}

// FieldByIndex ...
func (r *Row) FieldByIndex(index int) *Field {
	if index > len(r.fieldList) {
		return nil
	}
	return r.fieldList[index]
}

// AddField ...
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

// PrintValues ...
func (r *Row) PrintValues() {
	for _, f := range r.fieldList {
		fmt.Printf("%s: [%v]\t", f.Name, f.Value)
	}
	fmt.Println("")
}

// Where ...
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
		err = filters[idx].Validate(r)
		if err != nil {
			return nil, err
		}
	}

	for _, row := range r {
		ok = 0
		for idx := 0; idx < len(filters); idx++ {
			field = row.FieldByIndex(filters[idx].fieldIndex)
			if compareOk, err = field.compare(&filters[idx]); compareOk {
				if err != nil {
					return nil, err
				}
				ok++
			}
		}
		if ok == filtersLen {
			result = append(result, row)
		}
	}
	return result, nil
}

// AddField ...
func (r RowCollection) AddField(name string, value interface{}) {
	for _, row := range r {
		row.AddField(name, value)
	}
}

// Print ...
func (r RowCollection) Print() {
	for _, row := range r {
		row.PrintValues()
	}
}
