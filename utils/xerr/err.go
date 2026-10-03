// Package xerr provides a simple function to wrap an error in fmt.Errorf.
// Use like this:
//
//	func foo(stuff int) (err error){
//	  xerr.Expect(&err, `Do something with %d`, stuff)
//	}
package xerr

import "fmt"

func Expect(err *error, f string, a ...any) {
	if err == nil || *err == nil {
		return
	}

	*err = fmt.Errorf(`%s: %w`, fmt.Sprintf(f, a...), *err)
}
