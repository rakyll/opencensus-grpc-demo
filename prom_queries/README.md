# Prometheus Queries


## Client Request Bytes rate over 5 minutes

```
rate(kubecon_demo_grpc_io_client_request_bytes_cumulative_bucket[5m])
```

## Client Request Bytes 50th percentile over 30 minutes


```
histogram_quantile(0.5, rate(kubecon_demo_grpc_io_client_request_bytes_cumulative_bucket[30m]))
```

```
histogram_quantile(0.7, rate(kubecon_demo_grpc_io_client_roundtrip_latency_distribution_cumulative_bucket[30m]))
```