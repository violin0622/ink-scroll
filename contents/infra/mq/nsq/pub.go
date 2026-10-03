package nsq

import (
	"context"

	"github.com/samber/lo"
	"google.golang.org/protobuf/proto"

	novelv1 "github.com/violin0622/ink-scroll/proto/gen/go/contents/novel/v1"
	"github.com/violin0622/ink-scroll/utils/xerr"
)

type pub struct {
	topic string
	mq    *mq
}

func (mq *mq) Producer(topic string) pub {
	return pub{topic: topic, mq: mq}
}

// Pub implements [outbox.EventBus].
func (p pub) Pub(ctx context.Context, nes []*novelv1.Event) (err error) {
	xerr.Expect(&err, `publish novel events`)

	msgs, err := lo.MapErr(
		nes,
		func(ne *novelv1.Event, _ int) ([]byte, error) { return proto.Marshal(ne) },
	)

	return p.mq.p.MultiPublish(p.topic, msgs)
}
