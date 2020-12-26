package sink

import "context"

type SinkServerConfiguration struct {
	ToChan chan *SunkMessage
}

type SinkServer interface {
	Serve(ctx context.Context) error
	Stop(ctx context.Context) error
}
