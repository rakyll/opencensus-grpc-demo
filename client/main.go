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

	"go.opencensus.io/exporter/stats/prometheus"
	sstackdriver "go.opencensus.io/exporter/stats/stackdriver"
	tstackdriver "go.opencensus.io/exporter/trace/stackdriver"

	pb "go.opencensus.io/examples/grpc/proto"
)

const address = "localhost:50051"

func main() {
	pe, se := exporters()
	stats.RegisterExporter(pe)
	trace.RegisterExporter(se)

	go func() {
		// Serve the proometheus metrics endpoint at localhost:9999.
		http.Handle("/metrics", pe)
		log.Fatal(http.ListenAndServe(":9999", nil))
	}()

	// Subscribe to collect client request count as a distribution
	// and the count of the errored RPCs.
	if err := grpcstats.RPCClientRoundTripLatencyView.Subscribe(); err != nil {
		log.Fatal(err)
	}
	if err := grpcstats.RPCClientErrorCountView.Subscribe(); err != nil {
		log.Fatal(err)
	}

	stats.SetReportingPeriod(time.Second)
	trace.SetDefaultSampler(trace.AlwaysSample()) // for the demo

	// Set up a connection to the server with the OpenCensus
	// stats handler to enable stats and tracing.
	conn, err := grpc.Dial(
		address,
		grpc.WithStatsHandler(ocgrpc.NewClientStatsHandler()),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	for {
		r, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: strings.Repeat("*", rand.Intn(65536))})
		if err == nil {
			log.Printf("Greeting: %s", r.Message)
		}
		time.Sleep(1 * time.Second)
	}
}

func exporters() (*prometheus.Exporter, *tstackdriver.Exporter) {
	pe, err := prometheus.NewExporter(prometheus.Options{
		Namespace: "kubecon_demo",
	})
	if err != nil {
		log.Fatal(err)
	}
	se, err := tstackdriver.NewExporter(tstackdriver.Options{
		ProjectID: "jbdtalks",
	})
	if err != nil {
		log.Fatal(err)
	}
	return pe, se
}

func stackdriverMonitoringExporter() *sstackdriver.Exporter {
	se, err := sstackdriver.NewExporter(sstackdriver.Options{
		ProjectID: "jbdtalks",
	})
	if err != nil {
		log.Fatal(err)
	}
	return se
}
