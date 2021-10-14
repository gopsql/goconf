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

func ToConfigs(c interface{}) (configs []Config) {
	rt := reflect.TypeOf(c)
	rv := reflect.ValueOf(c)
	for i := 0; i < rt.NumField(); i++ {
		v := fmt.Sprint(rv.Field(i).Interface())
		tag := strings.TrimSpace(string(rt.Field(i).Tag))
		configs = append(configs, Config{
			Key:     rt.Field(i).Name,
			Value:   v,
			Comment: tag,
		})
	}
	return
}

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
	for i := 0; i < rt.NumField(); i++ {
		if i > 0 {
			output += "\n"
		}
		if tag := string(rt.Field(i).Tag); tag != "" {
			lines := strings.Split(strings.TrimSpace(tag), "\n")
			for _, line := range lines {
				output += "\t// " + line + "\n"
			}
		}
		v := fmt.Sprint(rv.Field(i).Interface())
		output += "\t" + rt.Field(i).Name + " = " + quoteString(v) + "\n"
	}
	output += `)
`
	return []byte(output), nil
}

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
				bl, ok := v.(*ast.BasicLit)
				if !ok {
					continue
				}
				value, err := strconv.Unquote(bl.Value)
				if err != nil {
					return err
				}
				for _, name := range vs.Names {
					field := rv.Elem().FieldByName(name.Name)
					if !field.IsValid() {
						continue
					}
					if field.Kind() == reflect.Ptr && field.IsNil() {
						field.Set(reflect.New(field.Type().Elem()))
					}
					if field.Kind() == reflect.String {
						field.SetString(value)
					} else if i, ok := field.Interface().(CanSetString); ok {
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
	return nil
}

func quoteString(in string) string {
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
