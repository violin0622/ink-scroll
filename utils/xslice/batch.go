package xslice

import "iter"

func Batch[S ~[]A, A any](sa S, b int) iter.Seq[S] {
	var cur int
	return func(yield func(S) bool) {
		if !yield(sa[cur : cur+b]) {
			return
		}
	}
}
