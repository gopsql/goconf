package goconf

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"strconv"
	"strings"
)

type Config struct {
	Key     string
	Value   string
	Comment string
}

type CanSetString interface {
	SetString(string) error
}

var (
	ErrNotStruct = errors.New("target must be a struct or a pointer to a struct")
)

// ToConfigs parses a struct and returns list of its keys, values and comments
// (struct tag values).
func ToConfigs(c interface{}) (configs []Config) {
	if c == nil {
		return
	}
	rt := reflect.TypeOf(c)
	if rt.Kind() != reflect.Struct {
		return
	}
	rv := reflect.ValueOf(c)
	for i := 0; i < rt.NumField(); i++ {
		ft := rt.Field(i)
		if ft.PkgPath != "" { // ignore unexported fields
			continue
		}
		v := fmt.Sprint(rv.Field(i).Interface())
		tag := strings.TrimSpace(string(ft.Tag))
		configs = append(configs, Config{
			Key:     ft.Name,
			Value:   v,
			Comment: tag,
		})
	}
	return
}

// Marshal collects exported keys, values and comments (tag values) of a struct
// and puts them in a Go file with constant strings.
func Marshal(c interface{}) ([]byte, error) {
	output := `package config

const (
`
	rt := reflect.TypeOf(c)
	rv := reflect.ValueOf(c)
	if !rv.IsValid() {
		return nil, ErrNotStruct
	}
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
		rv = rv.Elem()
	}
	if rt.Kind() != reflect.Struct {
		return nil, ErrNotStruct
	}
	count := 0
	for i := 0; i < rt.NumField(); i++ {
		ft := rt.Field(i)
		if ft.PkgPath != "" { // ignore unexported fields
			continue
		}
		if count > 0 {
			output += "\n"
		}
		if tag := string(rt.Field(i).Tag); tag != "" {
			lines := strings.Split(strings.TrimSpace(tag), "\n")
			for _, line := range lines {
				output += "\t// " + line + "\n"
			}
		}
		v := fmt.Sprint(rv.Field(i).Interface())
		output += "\t" + ft.Name + " = " + quoteString(ft, v) + "\n"
		count += 1
	}
	output += `)
`
	return []byte(output), nil
}

// Unmarshal gets all constants in a Go file and assign them to the struct
// provided.
func Unmarshal(input []byte, v interface{}) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", string(input), 0)
	if err != nil {
		return err
	}
	rv := reflect.ValueOf(v)
	for _, decl := range f.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		if gd.Tok != token.CONST {
			continue
		}
		for _, sp := range gd.Specs {
			vs, ok := sp.(*ast.ValueSpec)
			if !ok {
				continue
			}
			for _, v := range vs.Values {
				var value string
				var err error
				if bl, ok := v.(*ast.BasicLit); ok {
					if bl.Kind == token.STRING {
						value, err = strconv.Unquote(bl.Value)
					} else {
						value = bl.Value
					}
				} else if i, ok := v.(*ast.Ident); ok {
					value = i.Name
				} else {
					continue
				}
				if err != nil {
					return err
				}
				for _, name := range vs.Names {
					field := rv.Elem().FieldByName(name.Name)
					if !field.IsValid() {
						continue
					}
					kind := field.Kind()
					if kind == reflect.Ptr && field.IsNil() {
						field.Set(reflect.New(field.Type().Elem()))
					}
					switch field.Kind() {
					case reflect.String:
						field.SetString(value)
					case reflect.Bool:
						field.SetBool(value == "true")
					case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
						i, err := strconv.ParseInt(value, 10, 64)
						if err != nil {
							return err
						}
						field.SetInt(i)
					case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
						i, err := strconv.ParseUint(value, 10, 64)
						if err != nil {
							return err
						}
						field.SetUint(i)
					case reflect.Float32, reflect.Float64:
						f, err := strconv.ParseFloat(value, 64)
						if err != nil {
							return err
						}
						field.SetFloat(f)
					default:
						if i, ok := field.Interface().(CanSetString); ok {
							if err := i.SetString(value); err != nil {
								return err
							}
						} else if i, ok := field.Addr().Interface().(CanSetString); ok {
							if err := i.SetString(value); err != nil {
								return err
							}
						}
					}
				}
			}
		}
	}
	return nil
}

func quoteString(field reflect.StructField, in string) string {
	switch field.Type.Kind() {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return in
	}
	if useBackquote(in) {
		return "`" + in + "`"
	}
	return strconv.Quote(in)
}

func useBackquote(in string) bool {
	if strings.Count(in, "\n") == 0 {
		return false
	}
	lines := strings.Split(in, "\n")
	for _, line := range lines {
		if !strconv.CanBackquote(line) {
			return false
		}
	}
	return true
}
