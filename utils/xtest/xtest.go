// Package xtest provides a very simple short hand for
// naming and describe test cases for a long table driven test.
package xtest

import (
	"fmt"
	"strings"
)

type Base struct {
	id, desc string
}

func (b Base) String() string {
	if len(strings.Trim(b.desc, ` _/.`)) == 0 {
		return b.id
	}
	return fmt.Sprintf(`%s_%s`, b.id, b.desc)
}

func New(id string, desc ...string) Base {
	if len(id) == 0 {
		panic(`id cannot be empty string`)
	}

	return Base{id: id, desc: strings.Join(desc, ` `)}
}
