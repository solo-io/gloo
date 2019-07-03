// package: gloo.solo.io
// file: github.com/solo-io/gloo/projects/gloo/api/v1/plugins.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_struct_pb from "google-protobuf/google/protobuf/struct_pb";
import * as gogoproto_gogo_pb from "../../../../../../../gogoproto/gogo_pb";
import * as github_com_solo_io_gloo_projects_gloo_api_v1_ssl_pb from "../../../../../../../github.com/solo-io/gloo/projects/gloo/api/v1/ssl_pb";
import * as github_com_solo_io_gloo_projects_gloo_api_v1_extensions_pb from "../../../../../../../github.com/solo-io/gloo/projects/gloo/api/v1/extensions_pb";
import * as github_com_solo_io_gloo_projects_gloo_api_v1_circuit_breaker_pb from "../../../../../../../github.com/solo-io/gloo/projects/gloo/api/v1/circuit_breaker_pb";
import * as github_com_solo_io_gloo_projects_gloo_api_v1_load_balancer_pb from "../../../../../../../github.com/solo-io/gloo/projects/gloo/api/v1/load_balancer_pb";
import * as github_com_solo_io_gloo_projects_gloo_api_v1_connection_pb from "../../../../../../../github.com/solo-io/gloo/projects/gloo/api/v1/connection_pb";
import * as github_com_solo_io_gloo_projects_gloo_api_v1_plugins_aws_aws_pb from "../../../../../../../github.com/solo-io/gloo/projects/gloo/api/v1/plugins/aws/aws_pb";
import * as github_com_solo_io_gloo_projects_gloo_api_v1_plugins_rest_rest_pb from "../../../../../../../github.com/solo-io/gloo/projects/gloo/api/v1/plugins/rest/rest_pb";
import * as github_com_solo_io_gloo_projects_gloo_api_v1_plugins_grpc_grpc_pb from "../../../../../../../github.com/solo-io/gloo/projects/gloo/api/v1/plugins/grpc/grpc_pb";
import * as github_com_solo_io_gloo_projects_gloo_api_v1_plugins_grpc_web_grpc_web_pb from "../../../../../../../github.com/solo-io/gloo/projects/gloo/api/v1/plugins/grpc_web/grpc_web_pb";
import * as github_com_solo_io_gloo_projects_gloo_api_v1_plugins_hcm_hcm_pb from "../../../../../../../github.com/solo-io/gloo/projects/gloo/api/v1/plugins/hcm/hcm_pb";
import * as github_com_solo_io_gloo_projects_gloo_api_v1_plugins_tcp_tcp_pb from "../../../../../../../github.com/solo-io/gloo/projects/gloo/api/v1/plugins/tcp/tcp_pb";
import * as github_com_solo_io_gloo_projects_gloo_api_v1_plugins_azure_azure_pb from "../../../../../../../github.com/solo-io/gloo/projects/gloo/api/v1/plugins/azure/azure_pb";
import * as github_com_solo_io_gloo_projects_gloo_api_v1_plugins_consul_consul_pb from "../../../../../../../github.com/solo-io/gloo/projects/gloo/api/v1/plugins/consul/consul_pb";
import * as github_com_solo_io_gloo_projects_gloo_api_v1_plugins_kubernetes_kubernetes_pb from "../../../../../../../github.com/solo-io/gloo/projects/gloo/api/v1/plugins/kubernetes/kubernetes_pb";
import * as github_com_solo_io_gloo_projects_gloo_api_v1_plugins_retries_retries_pb from "../../../../../../../github.com/solo-io/gloo/projects/gloo/api/v1/plugins/retries/retries_pb";
import * as github_com_solo_io_gloo_projects_gloo_api_v1_plugins_static_static_pb from "../../../../../../../github.com/solo-io/gloo/projects/gloo/api/v1/plugins/static/static_pb";
import * as github_com_solo_io_gloo_projects_gloo_api_v1_plugins_stats_stats_pb from "../../../../../../../github.com/solo-io/gloo/projects/gloo/api/v1/plugins/stats/stats_pb";
import * as github_com_solo_io_gloo_projects_gloo_api_v1_plugins_transformation_prefix_rewrite_pb from "../../../../../../../github.com/solo-io/gloo/projects/gloo/api/v1/plugins/transformation/prefix_rewrite_pb";
import * as github_com_solo_io_gloo_projects_gloo_api_v1_plugins_transformation_transformation_pb from "../../../../../../../github.com/solo-io/gloo/projects/gloo/api/v1/plugins/transformation/transformation_pb";
import * as github_com_solo_io_gloo_projects_gloo_api_v1_plugins_faultinjection_fault_pb from "../../../../../../../github.com/solo-io/gloo/projects/gloo/api/v1/plugins/faultinjection/fault_pb";
import * as google_protobuf_duration_pb from "google-protobuf/google/protobuf/duration_pb";

export class ListenerPlugins extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListenerPlugins.AsObject;
  static toObject(includeInstance: boolean, msg: ListenerPlugins): ListenerPlugins.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListenerPlugins, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListenerPlugins;
  static deserializeBinaryFromReader(message: ListenerPlugins, reader: jspb.BinaryReader): ListenerPlugins;
}

export namespace ListenerPlugins {
  export type AsObject = {
  }
}

export class HttpListenerPlugins extends jspb.Message {
  hasGrpcWeb(): boolean;
  clearGrpcWeb(): void;
  getGrpcWeb(): github_com_solo_io_gloo_projects_gloo_api_v1_plugins_grpc_web_grpc_web_pb.GrpcWeb | undefined;
  setGrpcWeb(value?: github_com_solo_io_gloo_projects_gloo_api_v1_plugins_grpc_web_grpc_web_pb.GrpcWeb): void;

  hasHttpConnectionManagerSettings(): boolean;
  clearHttpConnectionManagerSettings(): void;
  getHttpConnectionManagerSettings(): github_com_solo_io_gloo_projects_gloo_api_v1_plugins_hcm_hcm_pb.HttpConnectionManagerSettings | undefined;
  setHttpConnectionManagerSettings(value?: github_com_solo_io_gloo_projects_gloo_api_v1_plugins_hcm_hcm_pb.HttpConnectionManagerSettings): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): HttpListenerPlugins.AsObject;
  static toObject(includeInstance: boolean, msg: HttpListenerPlugins): HttpListenerPlugins.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: HttpListenerPlugins, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): HttpListenerPlugins;
  static deserializeBinaryFromReader(message: HttpListenerPlugins, reader: jspb.BinaryReader): HttpListenerPlugins;
}

export namespace HttpListenerPlugins {
  export type AsObject = {
    grpcWeb?: github_com_solo_io_gloo_projects_gloo_api_v1_plugins_grpc_web_grpc_web_pb.GrpcWeb.AsObject,
    httpConnectionManagerSettings?: github_com_solo_io_gloo_projects_gloo_api_v1_plugins_hcm_hcm_pb.HttpConnectionManagerSettings.AsObject,
  }
}

export class TcpListenerPlugins extends jspb.Message {
  hasTcpProxySettings(): boolean;
  clearTcpProxySettings(): void;
  getTcpProxySettings(): github_com_solo_io_gloo_projects_gloo_api_v1_plugins_tcp_tcp_pb.TcpProxySettings | undefined;
  setTcpProxySettings(value?: github_com_solo_io_gloo_projects_gloo_api_v1_plugins_tcp_tcp_pb.TcpProxySettings): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TcpListenerPlugins.AsObject;
  static toObject(includeInstance: boolean, msg: TcpListenerPlugins): TcpListenerPlugins.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: TcpListenerPlugins, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TcpListenerPlugins;
  static deserializeBinaryFromReader(message: TcpListenerPlugins, reader: jspb.BinaryReader): TcpListenerPlugins;
}

export namespace TcpListenerPlugins {
  export type AsObject = {
    tcpProxySettings?: github_com_solo_io_gloo_projects_gloo_api_v1_plugins_tcp_tcp_pb.TcpProxySettings.AsObject,
  }
}

export class VirtualHostPlugins extends jspb.Message {
  hasExtensions(): boolean;
  clearExtensions(): void;
  getExtensions(): github_com_solo_io_gloo_projects_gloo_api_v1_extensions_pb.Extensions | undefined;
  setExtensions(value?: github_com_solo_io_gloo_projects_gloo_api_v1_extensions_pb.Extensions): void;

  hasRetries(): boolean;
  clearRetries(): void;
  getRetries(): github_com_solo_io_gloo_projects_gloo_api_v1_plugins_retries_retries_pb.RetryPolicy | undefined;
  setRetries(value?: github_com_solo_io_gloo_projects_gloo_api_v1_plugins_retries_retries_pb.RetryPolicy): void;

  hasStats(): boolean;
  clearStats(): void;
  getStats(): github_com_solo_io_gloo_projects_gloo_api_v1_plugins_stats_stats_pb.Stats | undefined;
  setStats(value?: github_com_solo_io_gloo_projects_gloo_api_v1_plugins_stats_stats_pb.Stats): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VirtualHostPlugins.AsObject;
  static toObject(includeInstance: boolean, msg: VirtualHostPlugins): VirtualHostPlugins.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: VirtualHostPlugins, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VirtualHostPlugins;
  static deserializeBinaryFromReader(message: VirtualHostPlugins, reader: jspb.BinaryReader): VirtualHostPlugins;
}

export namespace VirtualHostPlugins {
  export type AsObject = {
    extensions?: github_com_solo_io_gloo_projects_gloo_api_v1_extensions_pb.Extensions.AsObject,
    retries?: github_com_solo_io_gloo_projects_gloo_api_v1_plugins_retries_retries_pb.RetryPolicy.AsObject,
    stats?: github_com_solo_io_gloo_projects_gloo_api_v1_plugins_stats_stats_pb.Stats.AsObject,
  }
}

export class RoutePlugins extends jspb.Message {
  hasTransformations(): boolean;
  clearTransformations(): void;
  getTransformations(): github_com_solo_io_gloo_projects_gloo_api_v1_plugins_transformation_transformation_pb.RouteTransformations | undefined;
  setTransformations(value?: github_com_solo_io_gloo_projects_gloo_api_v1_plugins_transformation_transformation_pb.RouteTransformations): void;

  hasFaults(): boolean;
  clearFaults(): void;
  getFaults(): github_com_solo_io_gloo_projects_gloo_api_v1_plugins_faultinjection_fault_pb.RouteFaults | undefined;
  setFaults(value?: github_com_solo_io_gloo_projects_gloo_api_v1_plugins_faultinjection_fault_pb.RouteFaults): void;

  hasPrefixRewrite(): boolean;
  clearPrefixRewrite(): void;
  getPrefixRewrite(): github_com_solo_io_gloo_projects_gloo_api_v1_plugins_transformation_prefix_rewrite_pb.PrefixRewrite | undefined;
  setPrefixRewrite(value?: github_com_solo_io_gloo_projects_gloo_api_v1_plugins_transformation_prefix_rewrite_pb.PrefixRewrite): void;

  hasTimeout(): boolean;
  clearTimeout(): void;
  getTimeout(): google_protobuf_duration_pb.Duration | undefined;
  setTimeout(value?: google_protobuf_duration_pb.Duration): void;

  hasRetries(): boolean;
  clearRetries(): void;
  getRetries(): github_com_solo_io_gloo_projects_gloo_api_v1_plugins_retries_retries_pb.RetryPolicy | undefined;
  setRetries(value?: github_com_solo_io_gloo_projects_gloo_api_v1_plugins_retries_retries_pb.RetryPolicy): void;

  hasExtensions(): boolean;
  clearExtensions(): void;
  getExtensions(): github_com_solo_io_gloo_projects_gloo_api_v1_extensions_pb.Extensions | undefined;
  setExtensions(value?: github_com_solo_io_gloo_projects_gloo_api_v1_extensions_pb.Extensions): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RoutePlugins.AsObject;
  static toObject(includeInstance: boolean, msg: RoutePlugins): RoutePlugins.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RoutePlugins, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RoutePlugins;
  static deserializeBinaryFromReader(message: RoutePlugins, reader: jspb.BinaryReader): RoutePlugins;
}

export namespace RoutePlugins {
  export type AsObject = {
    transformations?: github_com_solo_io_gloo_projects_gloo_api_v1_plugins_transformation_transformation_pb.RouteTransformations.AsObject,
    faults?: github_com_solo_io_gloo_projects_gloo_api_v1_plugins_faultinjection_fault_pb.RouteFaults.AsObject,
    prefixRewrite?: github_com_solo_io_gloo_projects_gloo_api_v1_plugins_transformation_prefix_rewrite_pb.PrefixRewrite.AsObject,
    timeout?: google_protobuf_duration_pb.Duration.AsObject,
    retries?: github_com_solo_io_gloo_projects_gloo_api_v1_plugins_retries_retries_pb.RetryPolicy.AsObject,
    extensions?: github_com_solo_io_gloo_projects_gloo_api_v1_extensions_pb.Extensions.AsObject,
  }
}

export class DestinationSpec extends jspb.Message {
  hasAws(): boolean;
  clearAws(): void;
  getAws(): github_com_solo_io_gloo_projects_gloo_api_v1_plugins_aws_aws_pb.DestinationSpec | undefined;
  setAws(value?: github_com_solo_io_gloo_projects_gloo_api_v1_plugins_aws_aws_pb.DestinationSpec): void;

  hasAzure(): boolean;
  clearAzure(): void;
  getAzure(): github_com_solo_io_gloo_projects_gloo_api_v1_plugins_azure_azure_pb.DestinationSpec | undefined;
  setAzure(value?: github_com_solo_io_gloo_projects_gloo_api_v1_plugins_azure_azure_pb.DestinationSpec): void;

  hasRest(): boolean;
  clearRest(): void;
  getRest(): github_com_solo_io_gloo_projects_gloo_api_v1_plugins_rest_rest_pb.DestinationSpec | undefined;
  setRest(value?: github_com_solo_io_gloo_projects_gloo_api_v1_plugins_rest_rest_pb.DestinationSpec): void;

  hasGrpc(): boolean;
  clearGrpc(): void;
  getGrpc(): github_com_solo_io_gloo_projects_gloo_api_v1_plugins_grpc_grpc_pb.DestinationSpec | undefined;
  setGrpc(value?: github_com_solo_io_gloo_projects_gloo_api_v1_plugins_grpc_grpc_pb.DestinationSpec): void;

  getDestinationTypeCase(): DestinationSpec.DestinationTypeCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DestinationSpec.AsObject;
  static toObject(includeInstance: boolean, msg: DestinationSpec): DestinationSpec.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DestinationSpec, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DestinationSpec;
  static deserializeBinaryFromReader(message: DestinationSpec, reader: jspb.BinaryReader): DestinationSpec;
}

export namespace DestinationSpec {
  export type AsObject = {
    aws?: github_com_solo_io_gloo_projects_gloo_api_v1_plugins_aws_aws_pb.DestinationSpec.AsObject,
    azure?: github_com_solo_io_gloo_projects_gloo_api_v1_plugins_azure_azure_pb.DestinationSpec.AsObject,
    rest?: github_com_solo_io_gloo_projects_gloo_api_v1_plugins_rest_rest_pb.DestinationSpec.AsObject,
    grpc?: github_com_solo_io_gloo_projects_gloo_api_v1_plugins_grpc_grpc_pb.DestinationSpec.AsObject,
  }

  export enum DestinationTypeCase {
    DESTINATION_TYPE_NOT_SET = 0,
    AWS = 1,
    AZURE = 2,
    REST = 3,
    GRPC = 4,
  }
}

export class UpstreamSpec extends jspb.Message {
  hasSslConfig(): boolean;
  clearSslConfig(): void;
  getSslConfig(): github_com_solo_io_gloo_projects_gloo_api_v1_ssl_pb.UpstreamSslConfig | undefined;
  setSslConfig(value?: github_com_solo_io_gloo_projects_gloo_api_v1_ssl_pb.UpstreamSslConfig): void;

  hasCircuitBreakers(): boolean;
  clearCircuitBreakers(): void;
  getCircuitBreakers(): github_com_solo_io_gloo_projects_gloo_api_v1_circuit_breaker_pb.CircuitBreakerConfig | undefined;
  setCircuitBreakers(value?: github_com_solo_io_gloo_projects_gloo_api_v1_circuit_breaker_pb.CircuitBreakerConfig): void;

  hasLoadBalancerConfig(): boolean;
  clearLoadBalancerConfig(): void;
  getLoadBalancerConfig(): github_com_solo_io_gloo_projects_gloo_api_v1_load_balancer_pb.LoadBalancerConfig | undefined;
  setLoadBalancerConfig(value?: github_com_solo_io_gloo_projects_gloo_api_v1_load_balancer_pb.LoadBalancerConfig): void;

  hasConnectionConfig(): boolean;
  clearConnectionConfig(): void;
  getConnectionConfig(): github_com_solo_io_gloo_projects_gloo_api_v1_connection_pb.ConnectionConfig | undefined;
  setConnectionConfig(value?: github_com_solo_io_gloo_projects_gloo_api_v1_connection_pb.ConnectionConfig): void;

  getUseHttp2(): boolean;
  setUseHttp2(value: boolean): void;

  hasKube(): boolean;
  clearKube(): void;
  getKube(): github_com_solo_io_gloo_projects_gloo_api_v1_plugins_kubernetes_kubernetes_pb.UpstreamSpec | undefined;
  setKube(value?: github_com_solo_io_gloo_projects_gloo_api_v1_plugins_kubernetes_kubernetes_pb.UpstreamSpec): void;

  hasStatic(): boolean;
  clearStatic(): void;
  getStatic(): github_com_solo_io_gloo_projects_gloo_api_v1_plugins_static_static_pb.UpstreamSpec | undefined;
  setStatic(value?: github_com_solo_io_gloo_projects_gloo_api_v1_plugins_static_static_pb.UpstreamSpec): void;

  hasAws(): boolean;
  clearAws(): void;
  getAws(): github_com_solo_io_gloo_projects_gloo_api_v1_plugins_aws_aws_pb.UpstreamSpec | undefined;
  setAws(value?: github_com_solo_io_gloo_projects_gloo_api_v1_plugins_aws_aws_pb.UpstreamSpec): void;

  hasAzure(): boolean;
  clearAzure(): void;
  getAzure(): github_com_solo_io_gloo_projects_gloo_api_v1_plugins_azure_azure_pb.UpstreamSpec | undefined;
  setAzure(value?: github_com_solo_io_gloo_projects_gloo_api_v1_plugins_azure_azure_pb.UpstreamSpec): void;

  hasConsul(): boolean;
  clearConsul(): void;
  getConsul(): github_com_solo_io_gloo_projects_gloo_api_v1_plugins_consul_consul_pb.UpstreamSpec | undefined;
  setConsul(value?: github_com_solo_io_gloo_projects_gloo_api_v1_plugins_consul_consul_pb.UpstreamSpec): void;

  getUpstreamTypeCase(): UpstreamSpec.UpstreamTypeCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpstreamSpec.AsObject;
  static toObject(includeInstance: boolean, msg: UpstreamSpec): UpstreamSpec.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpstreamSpec, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpstreamSpec;
  static deserializeBinaryFromReader(message: UpstreamSpec, reader: jspb.BinaryReader): UpstreamSpec;
}

export namespace UpstreamSpec {
  export type AsObject = {
    sslConfig?: github_com_solo_io_gloo_projects_gloo_api_v1_ssl_pb.UpstreamSslConfig.AsObject,
    circuitBreakers?: github_com_solo_io_gloo_projects_gloo_api_v1_circuit_breaker_pb.CircuitBreakerConfig.AsObject,
    loadBalancerConfig?: github_com_solo_io_gloo_projects_gloo_api_v1_load_balancer_pb.LoadBalancerConfig.AsObject,
    connectionConfig?: github_com_solo_io_gloo_projects_gloo_api_v1_connection_pb.ConnectionConfig.AsObject,
    useHttp2: boolean,
    kube?: github_com_solo_io_gloo_projects_gloo_api_v1_plugins_kubernetes_kubernetes_pb.UpstreamSpec.AsObject,
    pb_static?: github_com_solo_io_gloo_projects_gloo_api_v1_plugins_static_static_pb.UpstreamSpec.AsObject,
    aws?: github_com_solo_io_gloo_projects_gloo_api_v1_plugins_aws_aws_pb.UpstreamSpec.AsObject,
    azure?: github_com_solo_io_gloo_projects_gloo_api_v1_plugins_azure_azure_pb.UpstreamSpec.AsObject,
    consul?: github_com_solo_io_gloo_projects_gloo_api_v1_plugins_consul_consul_pb.UpstreamSpec.AsObject,
  }

  export enum UpstreamTypeCase {
    UPSTREAM_TYPE_NOT_SET = 0,
    KUBE = 1,
    STATIC = 4,
    AWS = 2,
    AZURE = 3,
    CONSUL = 5,
  }
}

