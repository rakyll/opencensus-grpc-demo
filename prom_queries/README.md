# Prometheus Queries

## Server Request Bytes 70th percentile over 30 minutes

```
 histogram_quantile(0.7, rate(kubecon_demo_grpc_io_server_request_bytes_distribution_cumulative_bucket[30m]))
```
