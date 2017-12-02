# gRPC/OpenCensus Demo

Prerequisites:

* Running Prometheus server
* Stackdriver Tracing and Monitoring account

In this demo, we are forking the gRPC hello world server and client.
We will have a Go gRPC client talking to a gRPC server written in Java.

OpenCensus is vendor agnostic, it can upload data to various
tracing and metric collection backends. In this demo, we will
use Prometheus and Stackdriver Monitoring and Stackdriver Tracing.

gRPC Java is supporting OpenCensus instrumentation libraries out
of the box. The user only need to subscribe to the metrics they
want to collect and setup the exporter.

Let's run the gRPC Java server:

```
$ cd java
$ ./startserver.sh
```

Now the server is listening... Time to start sending some RPCs.

Let's take a look at the client first, client/main.go.

You can see we set a Prometheus exporter for stats and a Stackdriver
Tracing exporter for the traces.

We need to register the Prometheus metrics handler so it can
scrap and read the metrics we collected.

I am subscribing to some views; so the library collects,
aggregates and exports them.

Let's run the gRPC Go client:

```
$ go run client/main.go
```

Now, the client is tracing all the outgoing requests and collecting
stats from the outgoing RPCs. Let's see the Prometheus [metrics](http://localhost:9999/metrics)
endpoint to see the exported metrics.

You can see the exported distribution of latency with grpc method and
grpc service labels.

![s](https://i.imgur.com/bkaP7an.png)

I will switch to the Prometheus UI. You can see that :9999/metrics target is
up and has been scraped. We should be retrieving some metrics.

![stats](https://i.imgur.com/JpLy8E0.png)


See `kubecon_demo_grpc_io_client_request_bytes_cum_bucket` for the current
request size distribution.

See the rate of each bucket in the past 5 minutes:

```
rate(kubecon_demo_grpc_io_client_request_bytes_cumulative_bucket[5m])
```

See the 90th percentile:

```
histogram_quantile(0.9, rate(kubecon_demo_grpc_io_client_request_bytes_cumulative_bucket[5m]))
```

To be continued...
