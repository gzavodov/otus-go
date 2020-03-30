module github.com/gzavodov/otus-go/calendar

go 1.13

require (
	contrib.go.opencensus.io/exporter/prometheus v0.1.0
	github.com/golang/protobuf v1.3.2
	github.com/gzavodov/otus-go/banner-rotation v0.0.0-20200322190258-04645b9de684
	github.com/jackc/pgx/v4 v4.5.0
	github.com/prometheus/client_golang v1.5.0
	github.com/slok/go-http-metrics v0.6.1
	github.com/streadway/amqp v0.0.0-20200108173154-1c71cc93ed71
	go.uber.org/zap v1.14.0
	google.golang.org/grpc v1.25.1
)
