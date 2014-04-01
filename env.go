package factorcfg

import (
	"os"
	"strings"
)

type EnvGetValue func(string) string

func envGetValue(key string) string {
	return os.Getenv(key)
}

type Env struct{}

func NewEnv() *Env {
	return &Env{}
}

func (e *Env) All() (map[string]interface{}, error) {
	env := os.Environ()
	m := make(map[string]interface{})
	for _, pair := range env {
		kv := strings.Split(pair, "=")
		if kv != nil && len(kv) >= 2 {
			m[kv[0]] = kv[1]
		}
	}

	if len(m) > 0 {
		return m, nil
	}

	return nil, nil
}

func (e *Env) Tag() string {
	return "env"
}

// func (e *Env) Apply(spec interface{}) error {
// 	return envApply(envGetValue, spec)
// }

// func envApply(getValue EnvGetValue, spec interface{}) error {
// 	s := reflect.ValueOf(spec)
// 	if s.Kind() != reflect.Ptr {
// 		return errors.New("spec must be a pointer to a struct")
// 	}

// 	s = s.Elem()
// 	if s.Kind() != reflect.Struct {
// 		return ErrInvalidSpecification
// 	}

// 	return envapply(getValue, s)
// }

// func envapply(getValue EnvGetValue, s reflect.Value) error {
// 	typeOfSpec := s.Type()
// 	for i := 0; i < s.NumField(); i++ {
// 		f := s.Field(i)
// 		t := typeOfSpec.Field(i)
// 		if f.CanSet() {
// 			fieldName := typeOfSpec.Field(i).Name

// 			if f.Kind() == reflect.Struct {
// 				envapply(getValue, f)
// 				continue
// 			}

// 			key := t.Tag.Get("env")

// 			value := getValue(key)
// 			if value == "" {
// 				continue
// 			}

// 			switch f.Kind() {
// 			case reflect.String:
// 				f.SetString(value)
// 			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
// 				intValue, err := strconv.ParseInt(value, 0, f.Type().Bits())
// 				if err != nil {
// 					return &ParseError{
// 						FieldName: fieldName,
// 						TypeName:  f.Kind().String(),
// 						Value:     value,
// 					}
// 				}
// 				f.SetInt(intValue)
// 			case reflect.Bool:
// 				boolValue, err := strconv.ParseBool(value)
// 				if err != nil {
// 					return &ParseError{
// 						FieldName: fieldName,
// 						TypeName:  f.Kind().String(),
// 						Value:     value,
// 					}
// 				}
// 				f.SetBool(boolValue)
// 			case reflect.Float32:
// 				floatValue, err := strconv.ParseFloat(value, f.Type().Bits())
// 				if err != nil {
// 					return &ParseError{
// 						FieldName: fieldName,
// 						TypeName:  f.Kind().String(),
// 						Value:     value,
// 					}
// 				}
// 				f.SetFloat(floatValue)
// 			case reflect.Slice:
// 				f.Set(reflect.ValueOf(trimsplit(value, ",")))
// 			}
// 		}
// 	}

// 	return nil
// }
