// Package client contains a demo gRPC Go client
// with OpenCensus enabled to collect metrics and
// report traces for the outgoing RPCs.
package main

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"google.golang.org/grpc"

	ocgrpc "go.opencensus.io/plugin/grpc"
	"go.opencensus.io/plugin/grpc/grpcstats"
	"go.opencensus.io/stats"
	"go.opencensus.io/trace"

	"go.opencensus.io/exporter/prometheus"
	"go.opencensus.io/exporter/stackdriver"

	pb "go.opencensus.io/examples/grpc/proto"
)

func main() {
	prometheusExporter, stackdriverExporter := exporters()
	stats.RegisterExporter(prometheusExporter)
	trace.RegisterExporter(stackdriverExporter)

	// Subscribe to collect client request count as a distribution
	// and the count of the errored RPCs.
	views := []*stats.View{
		grpcstats.RPCClientRoundTripLatencyView,
		grpcstats.RPCClientErrorCountView,
		grpcstats.RPCClientRequestBytesView,
		grpcstats.RPCClientResponseBytesView,
	}
	for _, v := range views {
		if err := v.Subscribe(); err != nil {
			log.Fatal(err)
		}
	}

	go func() {
		// Serve the prometheus metrics endpoint at localhost:9999.
		http.Handle("/metrics", prometheusExporter)
		log.Fatal(http.ListenAndServe(":9999", nil))
	}()

	stats.SetReportingPeriod(1 * time.Second)

	// Set up a connection to the server with the OpenCensus
	// stats handler to enable stats and tracing.
	conn, err := grpc.Dial(
		"localhost:50051",
		grpc.WithStatsHandler(ocgrpc.NewClientStatsHandler()),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// For demoing purposes, always sample.
	trace.SetDefaultSampler(trace.AlwaysSample())

	ctx := context.Background()
	for {
		_, err := c.SayHello(ctx, &pb.HelloRequest{Name: strings.Repeat("*", rand.Intn(1<<16))})
		if err != nil {
			log.Printf("Failed to send request: %v", err)
		}
	}
}

func exporters() (*prometheus.Exporter, *stackdriver.Exporter) {
	pe, err := prometheus.NewExporter(prometheus.Options{
		Namespace: "kubecon_demo",
	})
	if err != nil {
		log.Fatal(err)
	}
	se, err := stackdriver.NewExporter(stackdriver.Options{
		ProjectID: "jbdtalks",
	})
	if err != nil {
		log.Fatal(err)
	}
	return pe, se
}
