package object

import "fmt"

type Type string

const (
	NULL_OBJ    = "NULL"
	INTEGER_OBJ = "INTEGER"
	BOOLEAN_OBJ = "BOOLEAN"
)

type Object interface {
	Type() Type
	Inspect() string
}

type Null struct{}

var _ Object = (*Null)(nil)

func (o *Null) Type() Type      { return NULL_OBJ }
func (o *Null) Inspect() string { return "null" }

type Integer struct {
	Value int64
}

var _ Object = (*Integer)(nil)

func (o *Integer) Type() Type      { return INTEGER_OBJ }
func (o *Integer) Inspect() string { return fmt.Sprint(o.Value) }

type Boolean struct {
	Value bool
}

var _ Object = (*Boolean)(nil)

func (o *Boolean) Type() Type      { return BOOLEAN_OBJ }
func (o *Boolean) Inspect() string { return fmt.Sprint(o.Value) }
