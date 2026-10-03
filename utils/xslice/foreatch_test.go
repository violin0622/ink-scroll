package xslice_test

import (
	"fmt"
	"slices"
	"testing"

	"github.com/violin0622/ink-scroll/utils/xslice"
)

func ExampleForeach() {
	xslice.Foreach(
		slices.Values([]string{`a`, `b`, `c`}),
		func(s string) { fmt.Print(s) })

	// Output: abc
}

func TestForeach(t *testing.T) {
}
