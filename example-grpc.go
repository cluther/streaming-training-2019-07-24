package main

import (
	"context"
	"crypto/tls"
	"log"
	"math/rand"
	"os"
	"time"

	structpb "github.com/golang/protobuf/ptypes/struct"
	zenoss "github.com/zenoss/zenoss-protobufs/go/cloud/data_receiver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

const (
	address    = "api.zenoss.io:443"
	apiKey     = "<ZENOSS_API_KEY>"
	sourceType = "com.zenoss.example-grpc.go"
	app        = "example-grpc.go"
)

var (
	source string
)

func init() {
	rand.Seed(time.Now().UnixNano())

	hostname, err := os.Hostname()
	if err != nil {
		source = "localhost"
	} else {
		source = hostname
	}
}

func main() {
	client, err := getClient()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	ctx = metadata.AppendToOutgoingContext(ctx, "zenoss-api-key", apiKey)

	dimensions := map[string]string{
		"app":    app,
		"source": source,
	}

	metadata := &structpb.Struct{
		Fields: map[string]*structpb.Value{
			"source-type": valueFromString(sourceType),
			"source":      valueFromString(source),
			"name":        valueFromString(app),
		},
	}

	log.Printf("sending model for %s", app)
	modelStatus, err := client.PutModels(ctx, &zenoss.Models{
		DetailedResponse: true,
		Models: []*zenoss.Model{
			&zenoss.Model{
				Timestamp:      time.Now().UnixNano() / 1e6,
				Dimensions:     dimensions,
				MetadataFields: metadata,
			},
		},
	})

	if err != nil {
		log.Printf("error sending model: %v", err)
	} else {
		if modelStatus.GetFailed() > 0 {
			log.Printf("model not accepted: %v", modelStatus.GetMessage())
		}
	}

	log.Printf("sending random.number metric for %s", app)
	metricStatus, err := client.PutMetrics(ctx, &zenoss.Metrics{
		DetailedResponse: true,
		Metrics: []*zenoss.Metric{
			&zenoss.Metric{
				Timestamp:      time.Now().UnixNano() / 1e6,
				Dimensions:     dimensions,
				MetadataFields: metadata,
				Metric:         "random.number",
				Value:          float64(rand.Int31()),
			},
		},
	})

	if err != nil {
		log.Printf("error sending metric: %v", err)
	} else {
		if metricStatus.GetFailed() > 0 {
			log.Printf("metric not accepted: %v", metricStatus.GetMessage())
		}
	}
}

func getClient() (zenoss.DataReceiverServiceClient, error) {
	opt := grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{}))
	conn, err := grpc.Dial(address, opt)
	if err != nil {
		return nil, err
	}

	return zenoss.NewDataReceiverServiceClient(conn), nil
}

func valueFromString(s string) *structpb.Value {
	return &structpb.Value{
		Kind: &structpb.Value_StringValue{
			StringValue: s,
		},
	}
}
