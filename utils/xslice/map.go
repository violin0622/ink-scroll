package xslice

import (
	"iter"
)

// func MapFn[SliceA ~[]A, SliceB []B, A any, B any](
// 	sa SliceA,
// 	fn func(A) B,
// ) SliceB {
// 	sb := make([]B, 0, len(sa))
// 	for _, a := range sa {
// 		sb = append(sb, fn(a))
// 	}
// 	return sb
// }

func Map[S iter.Seq[A], A, B any](sa S, fn func(A) B) iter.Seq[B] {
	return func(yield func(B) bool) {
		for a := range sa {
			if !yield(fn(a)) {
				return
			}
		}
	}
}

func MapOr[S iter.Seq[A], A, B any](
	sa S,
	fn func(A) (B, error),
	d B,
) iter.Seq[B] {
	return func(yield func(B) bool) {
		for a := range sa {
			b, err := fn(a)
			if err != nil {
				b = d
			}
			if yield(b) {
				return
			}
		}
	}
}

func MapOrElse[S iter.Seq[A], A, B any](
	sa S,
	fn func(A) (B, error),
	df func() B,
) iter.Seq[B] {
	return func(yield func(B) bool) {
		for a := range sa {
			b, err := fn(a)
			if err != nil {
				b = df()
			}
			if yield(b) {
				return
			}
		}
	}
}

func MapIdx[S iter.Seq2[K, V], K, V, A, B any](
	s S, fn func(K, V) (A, B),
) iter.Seq2[A, B] {
	return func(yeild func(A, B) bool) {
		for k, v := range s {
			if !yeild(fn(k, v)) {
				return
			}
		}
	}
}

func MapIdxOr[S iter.Seq2[K, V], K, V, A, B any](
	sa S,
	fn func(K, V) (A, B, error),
	da A, db B,
) iter.Seq2[A, B] {
	return func(yeild func(A, B) bool) {
		for k, v := range sa {
			k, v, err := fn(k, v)
			if err != nil {
				k, v = da, db
			}
			if !yeild(k, v) {
				return
			}
		}
	}
}

func MapIdxOrElse[S iter.Seq2[K, V], K, V, A, B any](
	sa S,
	fn func(K, V) (A, B, error),
	df func(K, V) (A, B),
) iter.Seq2[A, B] {
	return func(yeild func(A, B) bool) {
		for k, v := range sa {
			a, b, err := fn(k, v)
			if err != nil {
				a, b = df(k, v)
			}
			if !yeild(a, b) {
				return
			}
		}
	}
}

func Collect[S iter.Seq2[int, E], E any](s S) []E {
	var arr []E
	for _, e := range s {
		arr = append(arr, e)
	}

	return arr
}

// func MapIdxFn[SliceA ~[]A, SliceB []B, A, B any](
// 	sa SliceA,
// 	fn func(int, A) B,
// ) SliceB {
// 	sb := make([]B, 0, len(sa))
// 	for i, a := range sa {
// 		sb = append(sb, fn(i, a))
// 	}
// 	return sb
// }

// func Parallel[A, B any, SliceA ~[]A, SliceB ~[]B](
// 	sa SliceA, sb SliceB, fn func(int),
// ) {
// 	for i := range min(len(sa), len(sb)) {
// 		fn(i)
// 	}
// }
