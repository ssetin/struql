package struql

type StringModifier func(Value string) string

// Filter ...
type Filter struct {
	FieldName string
	Value     interface{}
	Operation int
	Modifier  StringModifier
}
