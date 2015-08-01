# malog

Fetch the latest updates on [Metal Archives](http://www.metal-archives.com/).

Example:

```
package main

import (
	"fmt"
	"time"

	"github.com/mvrilo/malog"
)

func main() {
	println("[malog] Started")
	res, err := malog.Fetch()
	for {
		select {
		case e := <-err:
			fmt.Println("[malog] ", e)
		case r := <-res:
			fmt.Printf("[malog] %s %s: %s %s\n", r.Title, r.Type, r.Name, r.URL)
		case <-time.Tick(1 * time.Minute):
			println("[malog] Fetch")
			res, err = malog.Fetch()
		}
	}
}
```

### AUTHOR

Murilo Santana <<mvrilo@gmail.com>>

### LICENSE

MIT
