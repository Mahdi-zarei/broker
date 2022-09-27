package modules

import (
	"context"
	"errors"
	rpt "therealbroker/metrics"
	"therealbroker/package/broker"
)

type connection struct {
	pipe chan broker.Message
	ctx  context.Context
}

func NewConnection(size int, ctx context.Context) *connection {
	ch := make(chan broker.Message, size)
	return &connection{
		pipe: ch,
		ctx:  ctx,
	}
}

func (c connection) SendMessage(msg broker.Message) {
	c.pipe <- msg
}

func (c connection) GetChannel() chan broker.Message {
	return c.pipe
}

func (c connection) contextIsClosed() bool {
	return errors.Is(c.ctx.Err(), context.Canceled)
}

func (c connection) ShouldClose() bool {
	if c.contextIsClosed() {
		rpt.TotalSubs.Dec()
		close(c.pipe)
		return true
	}
	return false
}
