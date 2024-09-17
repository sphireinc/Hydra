package hydra

import "testing"

func TestDebug(t *testing.T) {
	debug = true
	p("Test")
	debug = false
}
