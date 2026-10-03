package nsq

import (
	"context"

	"github.com/nsqio/go-nsq"
	"google.golang.org/protobuf/proto"

	eventv1 "github.com/violin0622/ink-scroll/proto/gen/go/event/v1"
)

type sub struct {
	mq *mq
}

func (mq *mq) Consume(
	topic string,
	h func(context.Context, *eventv1.Event) error,
) (sub, error) {
	c, err := nsq.NewConsumer(topic, `tmp`, mq.cfg)
	if err != nil {
		return sub{}, err
	}

	c.AddHandler(nsq.HandlerFunc(func(m *nsq.Message) error {
		var e eventv1.Event
		if err := proto.Unmarshal(m.Body, &e); err != nil {
			return err
		}

		return h(context.TODO(), &e)
	}))

	return sub{}, nil
}
