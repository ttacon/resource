resource
========

Retrieves source files relative to a go package, and also embeds them in the source code for next time.

The first time you fun your code using a resource, it will load it from disk and wil also
write a generated file for next time to speed things up (this is also useful for making a build for production).


ROADMAP
========

Soon will come the ability to either not write the generated file, and also
timestamps for development mode (so the data in the cache can me cleared for
changing assets).

Examples
========

given the following code, 
```
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
```

If gopher.jpg is 333K, and we this code (hitting disk), the first time we see:
```
time image
file is 341503 bytes

real	0m0.046s
user	0m0.020s
sys		0m0.014s
```

Building and running a second time, with the now generated file as part of the
package, we see:
```
time image
file is 341503 bytes

real 0m0.007s
user 0m0.001s
sys	 0m0.003s
```