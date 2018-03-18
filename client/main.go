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

	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
	"go.opencensus.io/zpages"

	"go.opencensus.io/exporter/prometheus"
	"go.opencensus.io/exporter/stackdriver"

	pb "go.opencensus.io/examples/grpc/proto"
)

func main() {
	prometheusExporter, stackdriverExporter := exporters()
	view.RegisterExporter(prometheusExporter)
	trace.RegisterExporter(stackdriverExporter)

	// Subscribe to collect client request count as a distribution
	// and the count of the errored RPCs.
	views := []*view.View{
		ocgrpc.ClientRoundTripLatencyView,
		ocgrpc.ClientErrorCountView,
	}
	for _, v := range views {
		if err := v.Subscribe(); err != nil {
			log.Fatal(err)
		}
	}

	go func() {
		// Serve the prometheus metrics endpoint at localhost:9999.
		http.Handle("/metrics", prometheusExporter)
		http.Handle("/debug/", http.StripPrefix("/debug", zpages.Handler))
		log.Fatal(http.ListenAndServe(":9999", nil))
	}()

	view.SetReportingPeriod(1 * time.Second)

	// Set up a connection to the server with the OpenCensus
	// stats handler to enable stats and tracing.
	conn, err := grpc.Dial(
		"localhost:50051",
		grpc.WithStatsHandler(&ocgrpc.ClientHandler{}),
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
		resp, err := c.SayHello(ctx, &pb.HelloRequest{Name: strings.Repeat("*", rand.Intn(1<<16))})
		if err != nil {
			log.Printf("Failed to send request: %v", err)
		} else {
			log.Printf("Recieved %v", resp.Message)
		}
		time.Sleep(100 * time.Millisecond)
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
