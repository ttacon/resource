package main

import (
	"fmt"
	"github.com/ttacon/resource"
)

var (
	res = resource.Resource{
		Path:       "gopher.jpg",
		RelativeTo: "github.com/ttacon/resource/examples/image",
		IsMain:     true,
	}
)

func main() {
	data, _ := res.Get()
	fmt.Printf("file is %d bytes\n", len(data))
}
