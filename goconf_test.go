package goconf

import (
	"encoding/json"
	"fmt"
	"testing"
)

type bar struct{ Code string }

func (b *bar) SetString(input string) error { b.Code = input; return nil }
func (b bar) String() string                { return b.Code }

func Test(t *testing.T) {
	if ToConfigs(nil) != nil {
		t.Error("ToConfigs(nil) should return nil")
	}
	if _, err := Marshal(nil); err != ErrNotStruct {
		t.Error("Marshal(nil) should return ErrNotStruct")
	}
	if _, err := Marshal(1); err != ErrNotStruct {
		t.Error("Marshal(1) should return ErrNotStruct")
	}
	if _, err := Marshal("a"); err != ErrNotStruct {
		t.Error("Marshal(a) should return ErrNotStruct")
	}
	a := &struct {
		Name string
		Bar  *bar
		Old  string
	}{"foo", &bar{"bar"}, "old"}
	c, err := Marshal(a)
	if err != nil {
		panic(err)
	}
	var b struct {
		Bar  *bar
		Name string
		New  string
	}
	err = Unmarshal(c, &b)
	if err != nil {
		panic(err)
	}
	if a.Name != b.Name {
		t.Error("name should be equal")
	}
	if a.Bar.Code != b.Bar.Code {
		t.Error("code should be equal")
	}
}

func ExampleMarshal() {
	a := struct {
		Foo string `
Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor
incididunt ut labore et dolore magna aliqua.`
		Hello string
	}{
		"Bar", "World",
	}
	c, err := Marshal(a)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(c))
	b, err := json.MarshalIndent(ToConfigs(a), "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	// Output:
	// package config
	//
	// const (
	// 	// Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor
	// 	// incididunt ut labore et dolore magna aliqua.
	// 	Foo = "Bar"
	//
	// 	Hello = "World"
	// )
	//
	// [
	//   {
	//     "Key": "Foo",
	//     "Value": "Bar",
	//     "Comment": "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor\nincididunt ut labore et dolore magna aliqua."
	//   },
	//   {
	//     "Key": "Hello",
	//     "Value": "World",
	//     "Comment": ""
	//   }
	// ]
}

func ExampleUnmarshal() {
	var c struct {
		Foo   string
		Hello string
	}
	fmt.Printf("before: %+v\n", c)
	err := Unmarshal([]byte(`
package config

const Foo = "BAR"

const (
	Hello = "WORLD"
)
`), &c)
	if err != nil {
		panic(err)
	}
	fmt.Printf("after: %+v\n", c)
	// Output:
	// before: {Foo: Hello:}
	// after: {Foo:BAR Hello:WORLD}
}
