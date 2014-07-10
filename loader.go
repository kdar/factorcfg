package factorcfg

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type Loader struct {
	handlers   *handlerList
	all        map[string]map[string]interface{}
	savedError error
}

// Create a Loader
func NewLoader() *Loader {
	return &Loader{
		handlers: &handlerList{make([]Handler, 0)},
	}
}

// Use a custom Handler for loading configuration
func (l *Loader) Use(a ...Handler) *Loader {
	for _, v := range a {
		l.handlers.Add(v)
	}
	return l
}

// Load configuration with registered Handlers
func (l *Loader) Load(spec interface{}) error {
	l.all = make(map[string]map[string]interface{})
	defer func() {
		l.all = nil
	}()
	for _, a := range l.handlers.list {
		m, err := a.All()
		if err != nil {
			return err
		}
		if _, ok := l.all[a.Tag()]; ok {
			for k, v := range m {
				l.all[a.Tag()][k] = v
			}
		} else {
			l.all[a.Tag()] = m
		}
	}

	l.apply(spec)
	return l.savedError
}

func (l *Loader) saveError(err error) {
	if l.savedError == nil {
		l.savedError = err
	}
}

func (l *Loader) apply(spec interface{}) {
	s := reflect.ValueOf(spec)
	if s.Kind() != reflect.Ptr {
		l.saveError(errors.New("spec must be a pointer to a struct"))
		return
	}

	s = s.Elem()
	if s.Kind() != reflect.Struct {
		l.saveError(ErrInvalidSpecification)
		return
	}

	l.applyr(s)
}

func (l *Loader) applyr(s reflect.Value) {
	typeOfSpec := s.Type()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		if !f.CanSet() {
			continue
		}

		t := typeOfSpec.Field(i)
		opts := strings.Split(t.Tag.Get("opts"), ",")
		fieldName := typeOfSpec.Field(i).Name

		if f.Kind() == reflect.Struct {
			l.applyr(f)
			continue
		}

		var value interface{}
		tags := getTags(t.Tag)
		for tag, key := range tags {
			if m, ok := l.all[tag]; ok {
				value = m[key]
			}
		}

		if value == nil {
			if stringInSlice(opts, "required") {
				l.saveError(&RequiredError{
					Name: fieldName,
					Tags: tags,
				})
			}
			continue
		}

		switch f.Kind() {
		case reflect.String:
			switch vt := value.(type) {
			case string:
				f.SetString(vt)
			default:
				f.SetString(fmt.Sprintf("%v", value))
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			switch vt := value.(type) {
			case string:
				intValue, err := strconv.ParseInt(vt, 0, f.Type().Bits())
				if err != nil {
					l.saveError(&ParseError{
						FieldName: fieldName,
						TypeName:  f.Kind().String(),
						Value:     value,
					})
					continue
				}
				f.SetInt(intValue)
			default:
				f.Set(reflect.ValueOf(value))
			}

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			switch vt := value.(type) {
			case string:
				intValue, err := strconv.ParseUint(vt, 0, f.Type().Bits())
				if err != nil {
					l.saveError(&ParseError{
						FieldName: fieldName,
						TypeName:  f.Kind().String(),
						Value:     value,
					})
					continue
				}

				f.SetUint(intValue)
			default:
				f.Set(reflect.ValueOf(value))
			}

		case reflect.Bool:
			switch vt := value.(type) {
			case string:
				boolValue, err := strconv.ParseBool(vt)
				if err != nil {
					l.saveError(&ParseError{
						FieldName: fieldName,
						TypeName:  f.Kind().String(),
						Value:     value,
					})
					continue
				}
				f.SetBool(boolValue)
			default:
				f.Set(reflect.ValueOf(value))
			}
		case reflect.Float32:
			switch vt := value.(type) {
			case string:
				floatValue, err := strconv.ParseFloat(vt, f.Type().Bits())
				if err != nil {
					l.saveError(&ParseError{
						FieldName: fieldName,
						TypeName:  f.Kind().String(),
						Value:     value,
					})
					continue
				}
				f.SetFloat(floatValue)
			default:
				f.Set(reflect.ValueOf(value))
			}
		case reflect.Slice:
			switch vt := value.(type) {
			case string:
				f.Set(reflect.ValueOf(trimsplit(vt, ",")))
			default:
				f.Set(reflect.ValueOf(value))
			}
		}
	}
}
