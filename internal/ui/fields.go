package ui

type FieldType int

const (
	FieldBool FieldType = iota
	FieldString
	FieldInt
)

type Field struct {
	Label string
	Type  FieldType

	// pointers into cfg
	Bool   *bool
	String *string
	Int    *int
}
