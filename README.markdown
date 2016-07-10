# cmus

`cmus` is a simple Go package for interacting with the [cmus][] console music player.

[cmus]: https://cmus.github.io/

## Usage

```go
package main

import (
	"fmt"
	"log"

	"github.com/stewart/cmus"
)

func main() {
	client := cmus.Client{}

	err := client.Connect()
	if err != nil {
		log.Fatal(err)
	}

	// get the status of the cmus instance
	status, err := client.Status()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(status)

	// pause or resume playback of music
	client.PlayPause()

	// seek back twenty seconds
	client.Seek("-0:20")

	// raise the volume
	client.Volume("+20")
}
```

Additional documentation of available methods can be found [on godoc.org][godocs].

[godocs]: https://godoc.org/github.com/stewart/cmus

## Todo

- [ ] test coverage
- [ ] parse `client.Status()` response into struct

## License

This project is licensed under the MIT License.

License can be found [here](LICENSE).
