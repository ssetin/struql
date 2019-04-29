package struql

// Field ...
type Field struct {
	idx   int
	Name  string
	IsSet bool
	Value interface{}
}
