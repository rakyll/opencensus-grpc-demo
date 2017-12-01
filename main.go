// Package main contains a demo gRPC Go client
// with OpenCensus enabled to collect metrics and
// report traces for the outgoing RPCs.
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"google.golang.org/grpc"

	"go.opencensus.io/exporter/stats/prometheus"
	"go.opencensus.io/exporter/trace/stackdriver"
	ocgrpc "go.opencensus.io/plugin/grpc"
	"go.opencensus.io/plugin/grpc/grpcstats"
	"go.opencensus.io/stats"
	"go.opencensus.io/trace"

	pb "go.opencensus.io/examples/grpc/proto"
)

const address = "localhost:50051"

func main() {
	pe, se := exporters()
	stats.RegisterExporter(pe)
	trace.RegisterExporter(se)

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

	// Contact the server and print out its response.
	name := "kubecon"
	if len(os.Args) > 1 {
		name = os.Args[1]
	}

	go func() {
		for {
			r, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: name})
			if err == nil {
				log.Printf("Greeting: %s", r.Message)
			}
			time.Sleep(5 * time.Second)
		}
	}()

	http.Handle("/metrics", pe)
	log.Fatal(http.ListenAndServe(":9999", nil))
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
