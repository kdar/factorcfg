package factorcfg

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"text/template"
)

func Render(spec interface{}, tmpl *template.Template) ([]byte, error) {
	s := reflect.ValueOf(spec)
	if s.Kind() != reflect.Ptr {
		return nil, errors.New("spec must be a pointer to a struct")
	}

	s = s.Elem()
	if s.Kind() != reflect.Struct {
		return nil, ErrInvalidSpecification
	}

	fields := getTemplateFields(nil, s)

	buf := &bytes.Buffer{}
	err := tmpl.Execute(buf, fields)
	return buf.Bytes(), err
}

func getTemplateFields(keys []string, s reflect.Value) []TemplateField {
	var fields []TemplateField

	typeOfSpec := s.Type()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		t := typeOfSpec.Field(i)
		if f.CanSet() {
			fieldName := typeOfSpec.Field(i).Name

			if f.Kind() == reflect.Struct {
				keys = append(keys, fieldName)
				fields = append(fields, getTemplateFields(keys, f)...)
				continue
			}

			fields = append(fields, TemplateField{
				Name:   fieldName,
				Value:  f.Interface(),
				String: fmt.Sprintf("%#v", f.Interface()),
				Type:   t.Type,
				Tags:   getTags(t.Tag),
			})
		}
	}

	return fields
}

func getTags(tag reflect.StructTag) (m map[string]string) {
	for tag != "" {
		if m == nil {
			m = make(map[string]string)
		}

		// skip leading space
		i := 0
		for i < len(tag) && tag[i] == ' ' {
			i++
		}
		tag = tag[i:]
		if tag == "" {
			break
		}

		// scan to colon.
		// a space or a quote is a syntax error
		i = 0
		for i < len(tag) && tag[i] != ' ' && tag[i] != ':' && tag[i] != '"' {
			i++
		}
		if i+1 >= len(tag) || tag[i] != ':' || tag[i+1] != '"' {
			break
		}
		name := string(tag[:i])
		tag = tag[i+1:]

		// scan quoted string to find value
		i = 1
		for i < len(tag) && tag[i] != '"' {
			if tag[i] == '\\' {
				i++
			}
			i++
		}
		if i >= len(tag) {
			break
		}
		qvalue := string(tag[:i+1])
		tag = tag[i+1:]

		value, _ := strconv.Unquote(qvalue)
		m[name] = value
	}

	return
}
