/* eslint-disable */
// package: gloo.solo.io
// file: github.com/solo-io/solo-apis/api/gloo/gloo/v1/upstream.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_struct_pb from "google-protobuf/google/protobuf/struct_pb";
import * as extproto_ext_pb from "../../../../../../../extproto/ext_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_ssl_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/ssl_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_circuit_breaker_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/circuit_breaker_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_load_balancer_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/load_balancer_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_connection_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/connection_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_api_v2_core_health_check_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/api/v2/core/health_check_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_api_v2_cluster_outlier_detection_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/api/v2/cluster/outlier_detection_pb";
import * as github_com_solo_io_solo_kit_api_v1_solo_kit_pb from "../../../../../../../github.com/solo-io/solo-kit/api/v1/solo-kit_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_static_static_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/options/static/static_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pipe_pipe_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/options/pipe/pipe_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_kubernetes_kubernetes_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/options/kubernetes/kubernetes_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_aws_aws_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/options/aws/aws_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_azure_azure_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/options/azure/azure_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_consul_consul_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/options/consul/consul_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_aws_ec2_aws_ec2_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/options/aws/ec2/aws_ec2_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_failover_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/failover_pb";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";

export class UpstreamSpec extends jspb.Message {
  hasDiscoveryMetadata(): boolean;
  clearDiscoveryMetadata(): void;
  getDiscoveryMetadata(): DiscoveryMetadata | undefined;
  setDiscoveryMetadata(value?: DiscoveryMetadata): void;

  hasSslConfig(): boolean;
  clearSslConfig(): void;
  getSslConfig(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_ssl_pb.UpstreamSslConfig | undefined;
  setSslConfig(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_ssl_pb.UpstreamSslConfig): void;

  hasCircuitBreakers(): boolean;
  clearCircuitBreakers(): void;
  getCircuitBreakers(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_circuit_breaker_pb.CircuitBreakerConfig | undefined;
  setCircuitBreakers(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_circuit_breaker_pb.CircuitBreakerConfig): void;

  hasLoadBalancerConfig(): boolean;
  clearLoadBalancerConfig(): void;
  getLoadBalancerConfig(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_load_balancer_pb.LoadBalancerConfig | undefined;
  setLoadBalancerConfig(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_load_balancer_pb.LoadBalancerConfig): void;

  hasConnectionConfig(): boolean;
  clearConnectionConfig(): void;
  getConnectionConfig(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_connection_pb.ConnectionConfig | undefined;
  setConnectionConfig(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_connection_pb.ConnectionConfig): void;

  clearHealthChecksList(): void;
  getHealthChecksList(): Array<github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_api_v2_core_health_check_pb.HealthCheck>;
  setHealthChecksList(value: Array<github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_api_v2_core_health_check_pb.HealthCheck>): void;
  addHealthChecks(value?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_api_v2_core_health_check_pb.HealthCheck, index?: number): github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_api_v2_core_health_check_pb.HealthCheck;

  hasOutlierDetection(): boolean;
  clearOutlierDetection(): void;
  getOutlierDetection(): github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_api_v2_cluster_outlier_detection_pb.OutlierDetection | undefined;
  setOutlierDetection(value?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_api_v2_cluster_outlier_detection_pb.OutlierDetection): void;

  hasUseHttp2(): boolean;
  clearUseHttp2(): void;
  getUseHttp2(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setUseHttp2(value?: google_protobuf_wrappers_pb.BoolValue): void;

  hasKube(): boolean;
  clearKube(): void;
  getKube(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_kubernetes_kubernetes_pb.UpstreamSpec | undefined;
  setKube(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_kubernetes_kubernetes_pb.UpstreamSpec): void;

  hasStatic(): boolean;
  clearStatic(): void;
  getStatic(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_static_static_pb.UpstreamSpec | undefined;
  setStatic(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_static_static_pb.UpstreamSpec): void;

  hasPipe(): boolean;
  clearPipe(): void;
  getPipe(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pipe_pipe_pb.UpstreamSpec | undefined;
  setPipe(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pipe_pipe_pb.UpstreamSpec): void;

  hasAws(): boolean;
  clearAws(): void;
  getAws(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_aws_aws_pb.UpstreamSpec | undefined;
  setAws(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_aws_aws_pb.UpstreamSpec): void;

  hasAzure(): boolean;
  clearAzure(): void;
  getAzure(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_azure_azure_pb.UpstreamSpec | undefined;
  setAzure(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_azure_azure_pb.UpstreamSpec): void;

  hasConsul(): boolean;
  clearConsul(): void;
  getConsul(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_consul_consul_pb.UpstreamSpec | undefined;
  setConsul(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_consul_consul_pb.UpstreamSpec): void;

  hasAwsEc2(): boolean;
  clearAwsEc2(): void;
  getAwsEc2(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_aws_ec2_aws_ec2_pb.UpstreamSpec | undefined;
  setAwsEc2(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_aws_ec2_aws_ec2_pb.UpstreamSpec): void;

  hasFailover(): boolean;
  clearFailover(): void;
  getFailover(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_failover_pb.Failover | undefined;
  setFailover(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_failover_pb.Failover): void;

  hasInitialStreamWindowSize(): boolean;
  clearInitialStreamWindowSize(): void;
  getInitialStreamWindowSize(): google_protobuf_wrappers_pb.UInt32Value | undefined;
  setInitialStreamWindowSize(value?: google_protobuf_wrappers_pb.UInt32Value): void;

  hasInitialConnectionWindowSize(): boolean;
  clearInitialConnectionWindowSize(): void;
  getInitialConnectionWindowSize(): google_protobuf_wrappers_pb.UInt32Value | undefined;
  setInitialConnectionWindowSize(value?: google_protobuf_wrappers_pb.UInt32Value): void;

  hasMaxConcurrentStreams(): boolean;
  clearMaxConcurrentStreams(): void;
  getMaxConcurrentStreams(): google_protobuf_wrappers_pb.UInt32Value | undefined;
  setMaxConcurrentStreams(value?: google_protobuf_wrappers_pb.UInt32Value): void;

  hasHttpProxyHostname(): boolean;
  clearHttpProxyHostname(): void;
  getHttpProxyHostname(): google_protobuf_wrappers_pb.StringValue | undefined;
  setHttpProxyHostname(value?: google_protobuf_wrappers_pb.StringValue): void;

  hasIgnoreHealthOnHostRemoval(): boolean;
  clearIgnoreHealthOnHostRemoval(): void;
  getIgnoreHealthOnHostRemoval(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setIgnoreHealthOnHostRemoval(value?: google_protobuf_wrappers_pb.BoolValue): void;

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
    discoveryMetadata?: DiscoveryMetadata.AsObject,
    sslConfig?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_ssl_pb.UpstreamSslConfig.AsObject,
    circuitBreakers?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_circuit_breaker_pb.CircuitBreakerConfig.AsObject,
    loadBalancerConfig?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_load_balancer_pb.LoadBalancerConfig.AsObject,
    connectionConfig?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_connection_pb.ConnectionConfig.AsObject,
    healthChecksList: Array<github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_api_v2_core_health_check_pb.HealthCheck.AsObject>,
    outlierDetection?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_api_v2_cluster_outlier_detection_pb.OutlierDetection.AsObject,
    useHttp2?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    kube?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_kubernetes_kubernetes_pb.UpstreamSpec.AsObject,
    pb_static?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_static_static_pb.UpstreamSpec.AsObject,
    pipe?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_pipe_pipe_pb.UpstreamSpec.AsObject,
    aws?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_aws_aws_pb.UpstreamSpec.AsObject,
    azure?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_azure_azure_pb.UpstreamSpec.AsObject,
    consul?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_consul_consul_pb.UpstreamSpec.AsObject,
    awsEc2?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_aws_ec2_aws_ec2_pb.UpstreamSpec.AsObject,
    failover?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_failover_pb.Failover.AsObject,
    initialStreamWindowSize?: google_protobuf_wrappers_pb.UInt32Value.AsObject,
    initialConnectionWindowSize?: google_protobuf_wrappers_pb.UInt32Value.AsObject,
    maxConcurrentStreams?: google_protobuf_wrappers_pb.UInt32Value.AsObject,
    httpProxyHostname?: google_protobuf_wrappers_pb.StringValue.AsObject,
    ignoreHealthOnHostRemoval?: google_protobuf_wrappers_pb.BoolValue.AsObject,
  }

  export enum UpstreamTypeCase {
    UPSTREAM_TYPE_NOT_SET = 0,
    KUBE = 11,
    STATIC = 12,
    PIPE = 13,
    AWS = 14,
    AZURE = 15,
    CONSUL = 16,
    AWS_EC2 = 17,
  }
}

export class DiscoveryMetadata extends jspb.Message {
  getLabelsMap(): jspb.Map<string, string>;
  clearLabelsMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DiscoveryMetadata.AsObject;
  static toObject(includeInstance: boolean, msg: DiscoveryMetadata): DiscoveryMetadata.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DiscoveryMetadata, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DiscoveryMetadata;
  static deserializeBinaryFromReader(message: DiscoveryMetadata, reader: jspb.BinaryReader): DiscoveryMetadata;
}

export namespace DiscoveryMetadata {
  export type AsObject = {
    labelsMap: Array<[string, string]>,
  }
}

export class UpstreamStatus extends jspb.Message {
  getState(): UpstreamStatus.StateMap[keyof UpstreamStatus.StateMap];
  setState(value: UpstreamStatus.StateMap[keyof UpstreamStatus.StateMap]): void;

  getReason(): string;
  setReason(value: string): void;

  getReportedBy(): string;
  setReportedBy(value: string): void;

  getSubresourceStatusesMap(): jspb.Map<string, UpstreamStatus>;
  clearSubresourceStatusesMap(): void;
  hasDetails(): boolean;
  clearDetails(): void;
  getDetails(): google_protobuf_struct_pb.Struct | undefined;
  setDetails(value?: google_protobuf_struct_pb.Struct): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpstreamStatus.AsObject;
  static toObject(includeInstance: boolean, msg: UpstreamStatus): UpstreamStatus.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpstreamStatus, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpstreamStatus;
  static deserializeBinaryFromReader(message: UpstreamStatus, reader: jspb.BinaryReader): UpstreamStatus;
}

export namespace UpstreamStatus {
  export type AsObject = {
    state: UpstreamStatus.StateMap[keyof UpstreamStatus.StateMap],
    reason: string,
    reportedBy: string,
    subresourceStatusesMap: Array<[string, UpstreamStatus.AsObject]>,
    details?: google_protobuf_struct_pb.Struct.AsObject,
  }

  export interface StateMap {
    PENDING: 0;
    ACCEPTED: 1;
    REJECTED: 2;
    WARNING: 3;
  }

  export const State: StateMap;
}

export class UpstreamNamespacedStatuses extends jspb.Message {
  getStatusesMap(): jspb.Map<string, UpstreamStatus>;
  clearStatusesMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpstreamNamespacedStatuses.AsObject;
  static toObject(includeInstance: boolean, msg: UpstreamNamespacedStatuses): UpstreamNamespacedStatuses.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpstreamNamespacedStatuses, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpstreamNamespacedStatuses;
  static deserializeBinaryFromReader(message: UpstreamNamespacedStatuses, reader: jspb.BinaryReader): UpstreamNamespacedStatuses;
}

export namespace UpstreamNamespacedStatuses {
  export type AsObject = {
    statusesMap: Array<[string, UpstreamStatus.AsObject]>,
  }
}
