# goconf

Use Go file as config file.

```go
a := struct {
	Foo string `
Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor
incididunt ut labore et dolore magna aliqua.`
	Hello string
}{
	"Bar", "World",
}
c, err := goconf.Marshal(a)
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
```

For custom data types, Marshal() uses its String() method, Unmarshal() uses its SetString() method.

Examples:
- [RSA private key](https://github.com/gopsql/jwt/blob/v1.0.0/jwt.go#L117-L150)
