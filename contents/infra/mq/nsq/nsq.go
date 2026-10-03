// Package nsq is adaptor package of mq to nsq.
package nsq

import (
	"github.com/go-logr/logr"
	"github.com/nsqio/go-nsq"

	"github.com/violin0622/ink-scroll/utils/xerr"
)

const nsqAddr = `localhost:4150`

type mq struct {
	log *logr.Logger
	cfg *nsq.Config
	p   *nsq.Producer
}

func NewProducer(
	log *logr.Logger,
) (m *mq, err error) {
	xerr.Expect(&err, `new nsq Producer`)

	c := nsq.NewConfig()
	p, err := nsq.NewProducer(nsqAddr, c)
	if err != nil {
		return nil, err
	}

	return &mq{log: log, p: p}, nil
}
