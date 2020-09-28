package object

import (
	"bytes"
	"fmt"

	"github.com/ebiiim/monkey/ast"
)

// Type represents type (in Monkey language) of the Object.
type Type string

// Object types.
const (
	NULL_OBJ         = "NULL"
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ        = "ERROR"
	FUNCTION_OBJ     = "FUNCTION"
)

type Object interface {
	Type() Type
	Inspect() string
}

// Null contains a NULL type value.
type Null struct{}

var _ Object = (*Null)(nil)

func (o *Null) Type() Type      { return NULL_OBJ }
func (o *Null) Inspect() string { return "null" }

// Integer contains an INTEGER type value.
type Integer struct{ Value int64 }

var _ Object = (*Integer)(nil)

func (o *Integer) Type() Type      { return INTEGER_OBJ }
func (o *Integer) Inspect() string { return fmt.Sprint(o.Value) }

// Boolean contains a BOOLEAN type value.
type Boolean struct{ Value bool }

var _ Object = (*Boolean)(nil)

func (o *Boolean) Type() Type      { return BOOLEAN_OBJ }
func (o *Boolean) Inspect() string { return fmt.Sprint(o.Value) }

// ReturnValue wraps an Object as the return value used by evaluator.
type ReturnValue struct{ Value Object }

var _ Object = (*ReturnValue)(nil)

func (o *ReturnValue) Type() Type      { return RETURN_VALUE_OBJ }
func (o *ReturnValue) Inspect() string { return o.Value.Inspect() }

// Error contains an error that is used by evaluator.
type Error struct{ Message error }

var _ Object = (*Error)(nil)

func (o *Error) Type() Type      { return ERROR_OBJ }
func (o *Error) Inspect() string { return fmt.Sprintf("ERROR: %s", o.Message) }

type Function struct {
	Env        *Environment
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
}

var _ Object = (*Function)(nil)

func (o *Function) Type() Type { return FUNCTION_OBJ }
func (o *Function) Inspect() string {
	var out bytes.Buffer
	fmt.Fprint(&out, "fn(")
	for i, param := range o.Parameters {
		fmt.Fprint(&out, param)
		if i+1 != len(o.Parameters) {
			fmt.Fprintf(&out, ", ")
		}
	}
	fmt.Fprint(&out, ") {\n")
	fmt.Fprint(&out, o.Body.String())
	fmt.Fprint(&out, "\n}")
	return out.String()
}
