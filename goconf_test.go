package goconf

import (
	"fmt"
)

func ExampleMarshal() {
	c, err := Marshal(struct {
		Foo string `
Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor
incididunt ut labore et dolore magna aliqua.`
		Hello string
	}{
		"Bar", "World",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(string(c))
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
