package io.grpc.examples.helloworld;

import com.google.common.collect.ImmutableSet;
import io.opencensus.contrib.grpc.metrics.RpcViewConstants;
import io.opencensus.stats.Stats;
import io.opencensus.stats.View;
import java.util.Set;

final class GrpcViews {
  private static final Set<View> RPC_VIEW_SET =
      ImmutableSet.of(
          RpcViewConstants.RPC_CLIENT_ERROR_COUNT_VIEW,
          RpcViewConstants.RPC_CLIENT_ROUNDTRIP_LATENCY_VIEW,
          RpcViewConstants.RPC_CLIENT_REQUEST_BYTES_VIEW,
          RpcViewConstants.RPC_CLIENT_RESPONSE_BYTES_VIEW,
          RpcViewConstants.RPC_CLIENT_REQUEST_COUNT_VIEW,
          RpcViewConstants.RPC_CLIENT_RESPONSE_COUNT_VIEW,
          RpcViewConstants.RPC_CLIENT_UNCOMPRESSED_REQUEST_BYTES_VIEW,
          RpcViewConstants.RPC_CLIENT_UNCOMPRESSED_RESPONSE_BYTES_VIEW,
          RpcViewConstants.RPC_CLIENT_SERVER_ELAPSED_TIME_VIEW,
          RpcViewConstants.RPC_SERVER_ERROR_COUNT_VIEW,
          RpcViewConstants.RPC_SERVER_SERVER_LATENCY_VIEW,
          RpcViewConstants.RPC_SERVER_SERVER_ELAPSED_TIME_VIEW,
          RpcViewConstants.RPC_SERVER_REQUEST_BYTES_VIEW,
          RpcViewConstants.RPC_SERVER_RESPONSE_BYTES_VIEW,
          RpcViewConstants.RPC_SERVER_REQUEST_COUNT_VIEW,
          RpcViewConstants.RPC_SERVER_RESPONSE_COUNT_VIEW,
          RpcViewConstants.RPC_SERVER_UNCOMPRESSED_REQUEST_BYTES_VIEW,
          RpcViewConstants.RPC_SERVER_UNCOMPRESSED_RESPONSE_BYTES_VIEW);

  static void registerViews() {
    for (View view : RPC_VIEW_SET) {
      Stats.getViewManager().registerView(view);
    }
  }
}
