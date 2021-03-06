# gRPC/OpenCensus Demo

Prerequisites:

* Running Prometheus server
* Stackdriver Tracing and Monitoring account

In this demo, we are forking the gRPC hello world server and client.
We will have a Go gRPC client talking to a gRPC server written in Java.

OpenCensus is vendor agnostic, it can upload data to various
tracing and metric collection backends. For this demo, we will
use both Prometheus and Stackdriver Monitoring, and Stackdriver for tracing.
We could have also used Zipkin for tracing, or any backend with a supported exporter.

gRPC Java supports OpenCensus instrumentation libraries out
of the box. To enable this, you simply need to subscribe to the metrics that you
want to collect and pass your credentials to your chosen exporter.

Let's run the gRPC Java server:

```
$ cd server
$ ./startserver.sh
```

Now the server is listening ... Time to start sending some RPCs.

Let's take a look at the client first, see client/main.go source code.

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

You can see the exported distribution of latency with grpc method and grpc service labels.

![s](https://i.imgur.com/bkaP7an.png)

I will switch to the Prometheus UI. You can see that :9999/metrics target is
up and has been scraped. We should be retrieving some metrics.

![stats](https://i.imgur.com/JpLy8E0.png)


See `kubecon_demo_grpc_io_client_request_bytes_cum_bucket` for the current
request size distribution.

See the 50th percentile:

```
histogram_quantile(0.5, rate(kubecon_demo_grpc_io_client_request_bytes_cumulative_bucket[5m]))
```

And 70th:

```
histogram_quantile(0.7, rate(kubecon_demo_grpc_io_client_request_bytes_cumulative_bucket[5m]))
```

And 90th:

```
histogram_quantile(0.9, rate(kubecon_demo_grpc_io_client_request_bytes_cumulative_bucket[5m]))
```

We are also collecting the roundtrip latency, you can see the distribution as well.

![latency](https://i.imgur.com/LEqtb3d.png)

Metrics will be collected with the labels from the context
and grpc_service and grpc_method will be automatically added.

You can filter the metrics coming from the "SayHello" service by
filtering by the label:

```
kubecon_demo_grpc_io_client_roundtrip_latency_cumulative_bucket{grpc_service="SayHello"}
```

Talking about latency, we have a multi service system here.
Let's a take a look at the traces we collected to have more granular
understanding of the cause of latency.

Both gRPC client and server uploads the collected traces to the Stackdriver
Trace service. In the Google Cloud Console, we already can see some traces:

![traces](https://i.imgur.com/MB8dRki.png)

You see the both parts of the RPC call. Sent part is the outgoing
RPC call from the client that too 3ms in total. It took 1 ms for the
server to handle the request and respond, then it takes another 1ms
for the client to handle the response.

You can also generate reports and analysis from the collected traces and compare them.

![report](https://i.imgur.com/ip0LZ4G.png)

The 95th percentile is 4ms.

Now back to the gRPC Java server, we can take a look
at the exported metrics. In this demo, we are using Stackdriver Monitoring
to export the server metrics to. You can choose any vendor of your choice.

The following dashboard displays the mean and 99th percentile
round-trip latency collected from the server.

![stackdriver monitoring](https://i.imgur.com/n3CZzNl.png)

One more thing...

OpenCensus includes an optional dashboard that displays diagnostic performance information on the host upon which it is running.
This is very useful when diagnosing issues related to a particular host or when you need to inspect unfiltered and unsampled data.

![tracez](https://i.imgur.com/iQiprU5.png)

See http://localhost:9999/debug/tracez for traces.
