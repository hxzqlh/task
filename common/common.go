package common

import (
	"context"
	"io"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	MongodbUri  = "mongodb://admin:123456@127.0.0.1:27017"
	EtcdAddr    = "127.0.0.1:2379"
	NatsAddr    = "nats://127.0.0.1:4222"
	TaskApiAddr = ":8888"
	JaegerAddr  = "127.0.0.1:6831"
)

const (
	TaskServiceName       = "go.micro.service.task"
	TaskClientName        = "go.micro.client.task"
	TaskTopicName         = "go.micro.topic.task"
	TaskApiName           = "go.micro.api.task"
	ArchiementServiceName = "go.micro.service.achievement"
)

func ConnectMongo(uri string, timeout time.Duration) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	// Ping the primary
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	return client, nil
}

func NewJaegerTracer(serviceName string, addr string) (opentracing.Tracer, io.Closer, error) {
	//cfg := jaegercfg.Configuration{
	//	ServiceName: serviceName,
	//	Sampler: &jaegercfg.SamplerConfig{
	//		Type:  jaeger.SamplerTypeConst,
	//		Param: 1,
	//	},
	//	Reporter: &jaegercfg.ReporterConfig{
	//		LogSpans:            true,
	//		BufferFlushInterval: 1 * time.Second,
	//	},
	//}
	//
	//sender, err := jaeger.NewUDPTransport(addr, 0)
	//if err != nil {
	//	return nil, nil, err
	//}
	//
	//reporter := jaeger.NewRemoteReporter(sender)
	//tracer, closer, err := cfg.NewTracer(jaegercfg.Reporter(reporter))
	//
	//return tracer, closer, err
	cfg := jaegercfg.Configuration{
		ServiceName: serviceName,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
		},
	}

	tracer, closer, err := cfg.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
	return tracer, closer, err
}
