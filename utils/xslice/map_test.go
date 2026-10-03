package xslice_test

import (
	"fmt"
	"slices"
	"strconv"

	"github.com/violin0622/ink-scroll/utils/xslice"
)

func ExampleMap() {
	a := []int{1, 2, 3}

	b := xslice.Map(
		slices.Values(a),
		func(a int) int { return a * 3 })

	c := xslice.Map(
		slices.Values(a),
		func(a int) string { return strconv.Itoa(a) })

	fmt.Println(slices.Collect(b))
	fmt.Println(slices.Collect(c))

	//Output:
	//[3 6 9]
	//[1 2 3]
}

func ExampleMapIdx() {
	a := []int{1, 2, 3}

	b := xslice.MapIdx(
		slices.All(a),
		func(i, a int) (int, int) { return i, a * 3 })

	for i, n := range b {
		fmt.Println(i, n)
	}

	// Output:
	// 0 3
	// 1 6
	// 2 9
}
