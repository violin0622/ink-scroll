package common

import "context"

type Aggregate[ID comparable] interface {
	ID() ID
	Kind() string
}

type Entity[ID comparable] interface {
	ID() ID
	Kind() string
}

type AggregateRepo[ID comparable, A Aggregate[ID]] interface {
	Get(context.Context, ID) (A, error)
	Save(context.Context, A) error
	Update(context.Context, ID, func(context.Context, *A) error) error
}
