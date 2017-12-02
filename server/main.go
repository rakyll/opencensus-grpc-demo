// Package server contains a demo gRPC Go server
// with OpenCensus enabled to collect metrics and
// report traces for the incoming RPCs.
package main

import (
	"context"
	"log"
	"math/rand"
	"net"
	"net/http"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"go.opencensus.io/plugin/grpc/grpcstats"
	"go.opencensus.io/stats"
	"go.opencensus.io/trace"

	"go.opencensus.io/exporter/stats/prometheus"
	sstackdriver "go.opencensus.io/exporter/stats/stackdriver"
	tstackdriver "go.opencensus.io/exporter/trace/stackdriver"
	ocgrpc "go.opencensus.io/plugin/grpc"

	pb "go.opencensus.io/examples/grpc/proto"
)

func main() {
	pe, se := exporters()
	stats.RegisterExporter(pe)
	trace.RegisterExporter(se)

	go func() {
		// Serve the prometheus metrics endpoint at localhost:9998.
		http.Handle("/metrics", pe)
		log.Fatal(http.ListenAndServe(":9998", nil))
	}()

	// Subscribe to collect some views.
	if err := grpcstats.RPCServerErrorCountView.Subscribe(); err != nil {
		log.Fatal(err)
	}
	if err := grpcstats.RPCServerRequestBytesView.Subscribe(); err != nil {
		log.Fatal(err)
	}

	stats.SetReportingPeriod(time.Second)
	trace.SetDefaultSampler(trace.AlwaysSample()) // for the demo

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Set up a new server with the OpenCensus
	// stats handler to enable stats and tracing.
	s := grpc.NewServer(grpc.StatsHandler(ocgrpc.NewServerStatsHandler()))
	pb.RegisterGreeterServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
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

// server is used to implement helloworld.GreeterServer.
type server struct{}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	ctx = trace.StartSpan(ctx, "auth")
	time.Sleep(time.Duration(40+rand.Intn(10)) * time.Millisecond)
	trace.EndSpan(ctx)

	ctx = trace.StartSpan(ctx, "render")
	time.Sleep(time.Duration(200+rand.Intn(50)) * time.Millisecond)
	trace.EndSpan(ctx)

	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}
