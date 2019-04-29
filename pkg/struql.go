package struql

import (
	"errors"
	"reflect"
	"sync"
)

// StruQL ...
type StruQL struct {
	Rows RowCollection

	currentRow int
	// TODO:
	pool sync.Pool
}

// Init ...
func (s *StruQL) Init(object interface{}) error {
	reflObjectValue := reflect.ValueOf(object)
	dataKind := reflObjectValue.Kind()

	if dataKind != reflect.Struct {
		return errors.New("object should be a struct type")
	}

	s.Rows = make(RowCollection, 0, 5)
	row := &Row{}
	row.Init()
	s.Rows = append(s.Rows, row)
	s.currentRow = 0

	err := s.object2table(object)
	return err
}

// Where ...
func (s *StruQL) Where(filters ...Filter) RowCollection {
	return s.Rows.Where(filters...)
}

// expandRow ...
func (s *StruQL) expandRow(fc map[string]Field) {
	newRow := &Row{}
	newRow.copyFields(fc)

	s.Rows = append(s.Rows, newRow)
	s.currentRow++
}

// Print ...
func (s *StruQL) Print() {
	for _, row := range s.Rows {
		row.PrintValues()
	}
}

func (s *StruQL) object2table(object interface{}, prefix ...string) error {
	reflObjectValue := reflect.ValueOf(object)
	objectKind := reflObjectValue.Kind()
	objPrefix := ""
	for _, pref := range prefix {
		objPrefix += pref + "."
	}

	switch objectKind {
	case reflect.Struct:
		for i := 0; i < reflObjectValue.NumField(); i++ {
			fieldValue := reflObjectValue.Field(i)
			fieldKind := fieldValue.Kind()

			switch fieldKind {
			case reflect.Struct:
				s.object2table(fieldValue.Interface(), objPrefix+reflObjectValue.Type().Field(i).Name)
			case reflect.Slice:
				if fieldValue.Len() > 0 {
					elemKind := fieldValue.Index(0).Kind()

					if elemKind == reflect.Struct {
						fieldsToCopy := make(map[string]Field)
						for k, v := range s.Rows[s.currentRow].Fields {
							fieldsToCopy[k] = v
						}

						for j := 0; j < fieldValue.Len(); j++ {
							s.object2table(fieldValue.Index(j).Interface(), objPrefix+reflObjectValue.Type().Field(i).Name)
							if j < fieldValue.Len()-1 {
								s.expandRow(fieldsToCopy)
							}
						}
					} else {
						s.Rows.AddField(objPrefix+reflObjectValue.Type().Field(i).Name, fieldValue.Interface())
					}
				}
			default:
				s.Rows.AddField(objPrefix+reflObjectValue.Type().Field(i).Name, fieldValue.Interface())
			}
		}
	default:
	}

	return nil
}
