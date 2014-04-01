package factorcfg

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type KeyMapper func([]string) string

// ErrInvalidSpecification indicates that a specification is of the wrong type.
var ErrInvalidSpecification = errors.New("invalid specification: must be a struct")

// A ParseError occurs when an environment variable cannot be converted to
// the type required by a struct field during assignment.
type ParseError struct {
	FieldName string
	TypeName  string
	Value     interface{}
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("assigning to %s: converting '%s' to an %s", e.FieldName, e.Value, e.TypeName)
}

type RequiredError struct {
	Name string
	Tags map[string]string
}

func (e *RequiredError) Error() string {
	return fmt.Sprintf("config \"%s\" is required", e.Name)
}

// TemplateField is a structure used when rendering
// a template of the spec.
type TemplateField struct {
	Name   string
	Value  interface{}
	String string
	Type   reflect.Type
	Tags   map[string]string
}

func stringInSlice(list []string, a string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// trimsplit slices s into all substrings separated by sep and returns a
// slice of the substrings between the separator with all leading and trailing
// white space removed, as defined by Unicode.
func trimsplit(s, sep string) []string {
	trimmed := strings.Split(s, sep)
	for i := range trimmed {
		trimmed[i] = strings.TrimSpace(trimmed[i])
	}
	return trimmed
}
