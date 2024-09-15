package hydra

import "fmt"

const debug = false

func p(a ...any) {
	if debug {
		fmt.Println(a)
	}
}
