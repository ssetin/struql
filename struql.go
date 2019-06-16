package struql

import (
	"errors"
	"reflect"
)

// StruQL implements struct to table transforming algorithm
type StruQL struct {
	Rows       RowCollection
	currentRow int
}

// Init - initialize struql from presented object (struct)
func (s *StruQL) Init(object interface{}) error {
	reflObjectValue := reflect.ValueOf(object)
	dataKind := reflObjectValue.Kind()

	if dataKind != reflect.Struct {
		return errors.New("object should be a struct type")
	}

	s.Rows = make(RowCollection, 0, 5)
	row := newRow()
	s.Rows = append(s.Rows, row)
	s.currentRow = 0

	err := s.object2table(object)
	return err
}

// Where collects data in the rows according to filters
func (s *StruQL) Where(result RowCollection, filters ...Filter) (RowCollection, error) {
	return s.Rows.Where(result, filters...)
}

func (s *StruQL) copyRow(row Row) {
	newRow := &Row{}
	newRow.init()
	for _, f := range row.fieldList {
		newRow.fieldList = append(newRow.fieldList, f)
		newRow.fieldMap[f.Name] = f
	}
	s.Rows = append(s.Rows, newRow)
	s.currentRow++
}

// String represents struql rows as string
func (s StruQL) String() string {
	return s.Rows.String()
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
			if !fieldValue.CanInterface() {
				continue
			}
			fieldKind := fieldValue.Kind()

			if fieldKind == reflect.Slice && (fieldValue.IsNil() || fieldValue.Len() == 0) {
				fieldValue = reflect.MakeSlice(fieldValue.Type(), 1, 1)
			}

			switch fieldKind {
			case reflect.Struct:
				s.object2table(fieldValue.Interface(), objPrefix+reflObjectValue.Type().Field(i).Name)
			case reflect.Slice:
				if fieldValue.Len() > 0 {
					rowToCopy := *s.Rows[s.currentRow]

					for j := 0; j < fieldValue.Len(); j++ {
						elem := fieldValue.Index(j)
						elemKind := elem.Kind()

						if elemKind == reflect.Ptr {
							if elem.IsNil() {
								elem = reflect.Zero(elem.Type().Elem())
							} else {
								elem = reflect.Indirect(elem)
							}
							elemKind = elem.Kind()
						}

						if elemKind == reflect.Struct {
							s.object2table(elem.Interface(), objPrefix+reflObjectValue.Type().Field(i).Name)
							if j < fieldValue.Len()-1 {
								s.copyRow(rowToCopy)
							}
						} else {
							s.Rows.AddField(objPrefix+reflObjectValue.Type().Field(i).Name, fieldValue.Interface())
						}
					}
				}
			default:
				s.Rows.AddField(objPrefix+reflObjectValue.Type().Field(i).Name, fieldValue.Interface())
			}
		}
	}

	return nil
}
