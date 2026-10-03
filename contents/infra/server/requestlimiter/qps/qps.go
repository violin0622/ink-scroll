package qps

import (
	"golang.org/x/time/rate"
	epb "google.golang.org/genproto/googleapis/rpc/errdetails"

	requestlimter "github.com/violin0622/ink-scroll/contents/infra/server/requestlimiter"
)

var _ requestlimter.RequestLimiter = (*qps)(nil)

type qps struct {
	rl  *rate.Limiter
	api string
}

func New(api string) qps {
	return qps{
		rl:  rate.NewLimiter(0, 0),
		api: api,
	}
}

func (l *qps) Set(r rate.Limit, b int) {
	l.rl.SetBurst(b)
	l.rl.SetLimit(r)
}

// Request implements RequestLimiter.
func (l *qps) Request(requestlimter.RequestInfo) (r requestlimter.Result) {
	if l.rl.Allow() {
		return r
	}

	r.Deny = true
	r.Details = append(r.Details, &epb.QuotaFailure{
		Violations: []*epb.QuotaFailure_Violation{{
			Subject:     `TotalQPS`,
			Description: `Server has overall QPS upper limit.`,
			ApiService:  l.api,
			QuotaValue:  int64(l.rl.Limit()),
		}},
	})

	return r
}
