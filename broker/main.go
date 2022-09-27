package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"therealbroker/Internal/broker"
	"therealbroker/Internal/database"
	redisMem "therealbroker/Internal/redis"
	srv "therealbroker/api/proto"
	rpt "therealbroker/metrics"
)

const (
	service     = "broker"
	environment = "production"
	id          = 1
)

func tracerProvider(url string) (*tracesdk.TracerProvider, error) {
	// Create the Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, err
	}
	tp := tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// Record information about this application in a Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(service),
			attribute.String("environment", environment),
			attribute.Int64("ID", id),
		)),
	)
	return tp, nil
}

func setupPrometheus() {
	prometheus.MustRegister(rpt.SubscribeFailedCount, rpt.SubscribeSuccessCount, rpt.FetchFailedCount, rpt.FetchSuccessCount,
		rpt.PublishFailedCount, rpt.PublishSuccessCount, rpt.PublishTime, rpt.FetchTime, rpt.SubscribeTime, rpt.TotalSubs)

	http.Handle("/metrics", promhttp.Handler())
	go func() {
		err := http.ListenAndServe(":8000", nil)
		if err != nil {
			log.Fatal("Failed to start prometheus server with error ", err)
		}
	}()
}

func setupGRPC(module *broker.Module) {
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatal("could not start GPRC server with error ", err)
		return
	}
	s := grpc.NewServer()
	srv.RegisterBrokerServer(s, &srv.Server{
		Src:                       module,
		UnimplementedBrokerServer: srv.UnimplementedBrokerServer{},
	})

	err = s.Serve(lis)
	if err != nil {
		log.Fatal("Error in GPRC server ", err)
	}
}

func setupJeager() {
	tp, err := tracerProvider("http://jgrint:14268/api/traces")
	if err != nil {
		log.Fatal(err)
	}
	otel.SetTracerProvider(tp)
}

func main() {
	//setupJeager()

	setupPrometheus()

	module := broker.NewModule()

	go redisMem.InitRedisClient()

	go database.InitPostgresDB()

	// this function blocks as long as the GRPC server is up
	setupGRPC(module)
}
