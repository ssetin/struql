package struql

type StringModifier func(Value string) string

//type Modifier func(Value interface{}) interface{}

// Filter ...
type Filter struct {
	FieldName string
	Value     interface{}
	Modifier  StringModifier
}
