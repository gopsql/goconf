package goconf

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"strconv"
	"strings"
)

func Marshal(c interface{}) ([]byte, error) {
	output := `package config

const (
`
	rt := reflect.TypeOf(c)
	rv := reflect.ValueOf(c)
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
		output += "\t" + rt.Field(i).Name + " = " + strconv.Quote(v) + "\n"
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
					if field.Kind() == reflect.String {
						field.SetString(value)
					} else if i, ok := field.Addr().Interface().(interface{ SetString(string) error }); ok {
						err := i.SetString(value)
						if err != nil {
							return err
						}
					}
				}
			}
		}
	}
	return nil
}
