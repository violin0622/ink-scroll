// Package nop provides a no op requestlimiter.
package nop

import (
	"math"
	"sync/atomic"

	requestlimter "github.com/violin0622/ink-scroll/contents/infra/server/requestlimiter"
)

type nop struct {
	total atomic.Uint64
}

func New() nop { return nop{} }

func (n *nop) Request(requestlimter.RequestInfo) (r requestlimter.Result) {
	n.total.Add(1)
	return r
}

func (n *nop) Stats() requestlimter.Stats {
	return requestlimter.Stats{
		Total:  n.total.Load(),
		Denied: 0,
		Limit:  math.NaN(),
		Tokens: math.NaN(),
		Burst:  math.MaxInt,
	}
}
