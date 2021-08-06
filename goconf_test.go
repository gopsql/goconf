package goconf

import (
	"encoding/json"
	"fmt"
)

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
