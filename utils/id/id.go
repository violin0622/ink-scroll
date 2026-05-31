// Package id defines some common ID operations.
package id

type ID[A comparable] interface {
	Equal(A) bool
}
