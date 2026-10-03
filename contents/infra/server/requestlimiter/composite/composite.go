package composite

import (
	"golang.org/x/time/rate"

	requestlimter "github.com/violin0622/ink-scroll/contents/infra/server/requestlimiter"
)

func New(rl ...requestlimter.RequestLimiter) composite {
	return composite{rls: rl}
}

type composite struct {
	rls []requestlimter.RequestLimiter
}

func (composite) Set(rate.Limit, int) {}
func (c composite) Request(ri requestlimter.RequestInfo) requestlimter.Result {
	for _, rl := range c.rls {
		if r := rl.Request(ri); r.Deny {
			return r
		}
	}

	return requestlimter.Result{}
}
