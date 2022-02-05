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
		Hello  string
		Bool   bool
		Number int
		Uint   uint32
		Float  float64
	}{
		"Bar", "World", true, 123, 22, 1.23,
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
	//
	// 	Bool = true
	//
	// 	Number = 123
	//
	// 	Uint = 22
	//
	// 	Float = 1.23
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
	//   },
	//   {
	//     "Key": "Bool",
	//     "Value": "true",
	//     "Comment": ""
	//   },
	//   {
	//     "Key": "Number",
	//     "Value": "123",
	//     "Comment": ""
	//   },
	//   {
	//     "Key": "Uint",
	//     "Value": "22",
	//     "Comment": ""
	//   },
	//   {
	//     "Key": "Float",
	//     "Value": "1.23",
	//     "Comment": ""
	//   }
	// ]
}

func ExampleUnmarshal() {
	var c struct {
		Foo    string
		Hello  string
		Bool   bool
		Number int
		Uint   uint32
		Float  float64
	}
	fmt.Printf("before: %+v\n", c)
	err := Unmarshal([]byte(`
package config

const Foo = "BAR"

const (
	Hello  = "WORLD"
	Bool   = true
	Number = 123
	Uint   = 22
	Float  = 1.23
)
`), &c)
	if err != nil {
		panic(err)
	}
	fmt.Printf("after: %+v\n", c)
	// Output:
	// before: {Foo: Hello: Bool:false Number:0 Uint:0 Float:0}
	// after: {Foo:BAR Hello:WORLD Bool:true Number:123 Uint:22 Float:1.23}
}
