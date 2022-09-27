module therealbroker

go 1.15

require (
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/gocql/gocql v1.2.0
	github.com/lib/pq v1.10.6
	github.com/onsi/ginkgo v1.16.5 // indirect
	github.com/onsi/gomega v1.18.1 // indirect
	github.com/prometheus/client_golang v1.12.2
	github.com/stretchr/testify v1.7.1
	go.opentelemetry.io/otel v1.9.0
	go.opentelemetry.io/otel/exporters/jaeger v1.9.0
	go.opentelemetry.io/otel/sdk v1.9.0
	golang.org/x/sys v0.0.0-20220804214406-8e32c043e418 // indirect
	google.golang.org/grpc v1.48.0
	google.golang.org/protobuf v1.28.0
)
