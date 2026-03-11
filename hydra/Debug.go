package hydra

import "fmt"

var debug = false

func p(a ...any) {
	if debug {
		fmt.Println(a...)
	}
}
