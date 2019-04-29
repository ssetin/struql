package struql

import (
	"fmt"
)

// Row ...
type Row struct {
	Fields map[string]*Field
}

// RowCollection ...
type RowCollection []Row

// NewRow ...
func NewRow() Row {
	row := Row{}
	row.Init()
	return row
}

// Init ...
func (r *Row) Init() {
	r.Fields = make(map[string]*Field)
}

// Init ...
func (r *Row) copyFields(fc map[string]*Field) {
	r.Fields = make(map[string]*Field)
	for k, v := range fc {
		r.Fields[k] = v
	}
}

// FieldByName ...
func (r *Row) FieldByName(name string) *Field {
	if f, ok := r.Fields[name]; ok {
		return f
	}
	return nil
}

// AddField ...
func (r *Row) AddField(name string, value interface{}) {
	if _, ok := r.Fields[name]; ok {
		return
	}

	newField := &Field{
		Value: value,
	}
	r.Fields[name] = newField
}

// PrintValues ...
func (r *Row) PrintValues() {
	for fname, f := range r.Fields {
		fmt.Printf("%s: [%v]\t", fname, f.Value)
	}
	fmt.Println("")
}

// PrintHeaders ...
func (r *Row) PrintHeaders() {
	for names := range r.Fields {
		fmt.Printf("%s\t", names)
	}
	fmt.Println("")
}

// Where ...
func (r RowCollection) Where(result RowCollection, filters ...Filter) RowCollection {
	var (
		ok         int
		filtersLen int
		field      *Field
	)
	filtersLen = len(filters)

	for _, row := range r {
		ok = 0
		for _, filter := range filters {
			field = row.FieldByName(filter.FieldName)
			if field != nil && field.Value == filter.Value {
				ok++
			}
		}
		if ok == filtersLen {
			result = append(result, row)
		}
	}
	return result
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
