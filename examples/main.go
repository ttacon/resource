package main

import (
	"fmt"

	"github.com/ttacon/resource"
)

var res = resource.Resource{
	Path:       "examples/file.txt",
	RelativeTo: "github.com/ttacon/resource",
	IsMain:     true,
}

func main() {
	data, err := res.Get()
	fmt.Println("err: ", err)
	fmt.Printf("data: %d\n", len(data))
}
