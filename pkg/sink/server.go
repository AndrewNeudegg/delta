package sink

import "context"

type ServerConfiguration struct {
	ToChan chan *SunkMessage
}

type Server interface {
	Serve(ctx context.Context) error
	Stop(ctx context.Context) error
}
