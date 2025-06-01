package assert

import (
	"fmt"
	"strings"
)

func Eq[T comparable](expected T, got T, message ...string) {
	if expected == got {
		return
	}

	panicMsg := fmt.Sprintf("Assertion failed: expected %v, got %v", expected, got)
	if len(message) > 0 {
		panicMsg = strings.Join(message, " ")
	}
	panic(panicMsg)
}
