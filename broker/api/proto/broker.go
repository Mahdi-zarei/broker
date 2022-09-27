package srv

import (
	"context"
	"go.opentelemetry.io/otel"
	"log"
	rpt "therealbroker/metrics"
	broker2 "therealbroker/package/broker"
	"time"
)

type Server struct {
	Src broker2.Broker
	UnimplementedBrokerServer
}

const name = "api"

func (s *Server) Publish(ctx context.Context, pr *PublishRequest) (*PublishResponse, error) {
	//log.Print("received request for publish")
	newCtx, span := otel.Tracer(name).Start(ctx, "Publish")
	defer span.End()

	timer := time.Now()

	msg := broker2.Message{
		Body:       string(pr.Body),
		Expiration: time.Duration(pr.ExpirationSeconds),
	}

	val, err := s.Src.Publish(newCtx, pr.Subject, msg)
	if err != nil {
		rpt.PublishSuccessCount.Inc()
		return nil, err
	}
	rpt.PublishSuccessCount.Inc()
	defer rpt.PublishTime.Observe(float64(time.Since(timer)))

	return &PublishResponse{
		Id: int32(val),
	}, nil
}

func (s *Server) Subscribe(sr *SubscribeRequest, stream Broker_SubscribeServer) error {
	//log.Printf("received subscribe request : %v\n", sr)

	ctx, cancel := context.WithCancel(context.Background())
	newCtx, span := otel.Tracer(name).Start(ctx, "Subscribe")

	timer := time.Now()
	ch, err := s.Src.Subscribe(newCtx, sr.Subject)
	rpt.SubscribeSuccessCount.Inc()
	rpt.SubscribeTime.Observe(float64(time.Since(timer)))

	if err != nil {
		cancel()
		return err
	}
	span.End()

	for {
		select {
		case x, ok := <-ch:
			if !ok {
				cancel()
				return nil
			}
			tmp := &MessageResponse{
				Body: []byte(x.Body),
			}
			err := stream.Send(tmp)
			if err != nil {
				log.Printf("stream of a \"%v\" subscriber got closed with error %v\n", sr.Subject, err)
				cancel()
				return err
			}
		}
	}
}

func (s *Server) Fetch(ctx context.Context, fr *FetchRequest) (*MessageResponse, error) {
	log.Printf("recieved fetch request for %v\n", fr)

	newCtx, span := otel.Tracer(name).Start(ctx, "Fetch")
	defer span.End()

	timer := time.Now()
	msg, err := s.Src.Fetch(newCtx, fr.Subject, int(fr.Id))
	rpt.FetchSuccessCount.Inc()
	defer rpt.FetchTime.Observe(float64(time.Since(timer)))
	if err != nil {
		return nil, err
	}

	tmp := &MessageResponse{
		Body: []byte(msg.Body),
	}
	return tmp, nil
}
