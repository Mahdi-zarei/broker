package broker

import (
	"context"
	"errors"
	"go.opentelemetry.io/otel"
	"log"
	"sync"
	"therealbroker/Internal/database"
	"therealbroker/Internal/modules"
	redisMem "therealbroker/Internal/redis"
	rpt "therealbroker/metrics"
	"therealbroker/package/broker"
	"time"
)

const name = "internal"

type Module struct {
	subs      map[string]modules.Subscribers
	currentid int
	syncer sync.Mutex
	closed bool
}

func NewModule() *Module {
	return &Module{
		subs:      map[string]modules.Subscribers{},
		closed:    false,
		currentid: -1,
	}
}

func (m *Module) Close() error {
	if m.closed {
		return broker.ErrUnavailable
	}
	m.syncer.Lock()
	m.closed = true
	m.subs = nil
	m.syncer.Unlock()
	return nil
}

func (m *Module) assertSubscribersExists(subject string) {
	if _, ok := m.subs[subject]; !ok {
		m.subs[subject] = *modules.NewSubscribersList()
	}
}

func (m *Module) dispatchMessages(ctx context.Context, msg broker.Message, subject string) {
	m.syncer.Lock()
	defer m.syncer.Unlock()
	subList, ok := m.subs[subject]
	if !ok {
		return
	}
	for ch, _ := range subList {
		if ch.ShouldClose() {
			delete(subList, ch)
			continue
		}
		ch.SendMessage(msg)
	}
}

func (m *Module) Publish(ctx context.Context, subject string, msg broker.Message) (int, error) {
	if m.closed {
		return 0, broker.ErrUnavailable
	}
	newCtx, span := otel.Tracer(name).Start(ctx, "innerPublish")
	defer span.End()

	id, cr := redisMem.GetNewId(), time.Now()

	m.dispatchMessages(ctx, msg, subject)

	data := modules.NewRowObject(id, msg.Body, cr, int(msg.Expiration), subject)
	database.StoreToDatabase(newCtx, data)

	return id, nil
}

func (m *Module) Subscribe(ctx context.Context, subject string) (<-chan broker.Message, error) {
	if m.closed {
		return nil, broker.ErrUnavailable
	}
	if errors.Is(ctx.Err(), context.Canceled) {
		return nil, nil
	}
	log.Printf("new subscription of subject %v\n", subject)

	_, span := otel.Tracer(name).Start(ctx, "innerSubscribe")
	defer span.End()

	m.syncer.Lock()
	defer m.syncer.Unlock()

	m.assertSubscribersExists(subject)

	clt := *modules.NewConnection(1000, ctx)
	m.subs[subject].AddNewClient(clt)

	rpt.TotalSubs.Inc()

	return clt.GetChannel(), nil
}

func (m *Module) Fetch(ctx context.Context, subject string, id int) (broker.Message, error) {
	if m.closed {
		return broker.Message{}, broker.ErrUnavailable
	}
	log.Printf("fetching on subject %v\n", subject)

	data := database.GetMessage(id)
	if data.Id == -1 || data.Subject != subject {
		return broker.Message{}, broker.ErrInvalidID
	}
	if t := time.Since(data.Creation); t > time.Duration(1e9*data.Expiration) {
		return broker.Message{}, broker.ErrExpiredID
	}

	return broker.Message{
		Body:       data.Body,
		Expiration: time.Duration(data.Expiration),
	}, nil
}
