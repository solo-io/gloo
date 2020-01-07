// package: gloo.solo.io
// file: gloo/projects/gloo/api/v1/upstream.proto

import * as jspb from "google-protobuf";
import * as gogoproto_gogo_pb from "../../../../../gogoproto/gogo_pb";
import * as extproto_ext_pb from "../../../../../extproto/ext_pb";
import * as solo_kit_api_v1_metadata_pb from "../../../../../solo-kit/api/v1/metadata_pb";
import * as gloo_projects_gloo_api_v1_ssl_pb from "../../../../../gloo/projects/gloo/api/v1/ssl_pb";
import * as gloo_projects_gloo_api_v1_circuit_breaker_pb from "../../../../../gloo/projects/gloo/api/v1/circuit_breaker_pb";
import * as gloo_projects_gloo_api_v1_load_balancer_pb from "../../../../../gloo/projects/gloo/api/v1/load_balancer_pb";
import * as gloo_projects_gloo_api_v1_connection_pb from "../../../../../gloo/projects/gloo/api/v1/connection_pb";
import * as gloo_projects_gloo_api_external_envoy_api_v2_core_health_check_pb from "../../../../../gloo/projects/gloo/api/external/envoy/api/v2/core/health_check_pb";
import * as solo_kit_api_v1_status_pb from "../../../../../solo-kit/api/v1/status_pb";
import * as gloo_projects_gloo_api_external_envoy_api_v2_cluster_outlier_detection_pb from "../../../../../gloo/projects/gloo/api/external/envoy/api/v2/cluster/outlier_detection_pb";
import * as solo_kit_api_v1_solo_kit_pb from "../../../../../solo-kit/api/v1/solo-kit_pb";
import * as gloo_projects_gloo_api_v1_options_static_static_pb from "../../../../../gloo/projects/gloo/api/v1/options/static/static_pb";
import * as gloo_projects_gloo_api_v1_options_pipe_pipe_pb from "../../../../../gloo/projects/gloo/api/v1/options/pipe/pipe_pb";
import * as gloo_projects_gloo_api_v1_options_kubernetes_kubernetes_pb from "../../../../../gloo/projects/gloo/api/v1/options/kubernetes/kubernetes_pb";
import * as gloo_projects_gloo_api_v1_options_aws_aws_pb from "../../../../../gloo/projects/gloo/api/v1/options/aws/aws_pb";
import * as gloo_projects_gloo_api_v1_options_azure_azure_pb from "../../../../../gloo/projects/gloo/api/v1/options/azure/azure_pb";
import * as gloo_projects_gloo_api_v1_options_consul_consul_pb from "../../../../../gloo/projects/gloo/api/v1/options/consul/consul_pb";
import * as gloo_projects_gloo_api_v1_options_aws_ec2_aws_ec2_pb from "../../../../../gloo/projects/gloo/api/v1/options/aws/ec2/aws_ec2_pb";
import * as gloo_projects_gloo_api_v1_options_pb from "../../../../../gloo/projects/gloo/api/v1/options_pb";

export class Upstream extends jspb.Message {
  hasStatus(): boolean;
  clearStatus(): void;
  getStatus(): solo_kit_api_v1_status_pb.Status | undefined;
  setStatus(value?: solo_kit_api_v1_status_pb.Status): void;

  hasMetadata(): boolean;
  clearMetadata(): void;
  getMetadata(): solo_kit_api_v1_metadata_pb.Metadata | undefined;
  setMetadata(value?: solo_kit_api_v1_metadata_pb.Metadata): void;

  hasDiscoveryMetadata(): boolean;
  clearDiscoveryMetadata(): void;
  getDiscoveryMetadata(): DiscoveryMetadata | undefined;
  setDiscoveryMetadata(value?: DiscoveryMetadata): void;

  hasSslConfig(): boolean;
  clearSslConfig(): void;
  getSslConfig(): gloo_projects_gloo_api_v1_ssl_pb.UpstreamSslConfig | undefined;
  setSslConfig(value?: gloo_projects_gloo_api_v1_ssl_pb.UpstreamSslConfig): void;

  hasCircuitBreakers(): boolean;
  clearCircuitBreakers(): void;
  getCircuitBreakers(): gloo_projects_gloo_api_v1_circuit_breaker_pb.CircuitBreakerConfig | undefined;
  setCircuitBreakers(value?: gloo_projects_gloo_api_v1_circuit_breaker_pb.CircuitBreakerConfig): void;

  hasLoadBalancerConfig(): boolean;
  clearLoadBalancerConfig(): void;
  getLoadBalancerConfig(): gloo_projects_gloo_api_v1_load_balancer_pb.LoadBalancerConfig | undefined;
  setLoadBalancerConfig(value?: gloo_projects_gloo_api_v1_load_balancer_pb.LoadBalancerConfig): void;

  hasConnectionConfig(): boolean;
  clearConnectionConfig(): void;
  getConnectionConfig(): gloo_projects_gloo_api_v1_connection_pb.ConnectionConfig | undefined;
  setConnectionConfig(value?: gloo_projects_gloo_api_v1_connection_pb.ConnectionConfig): void;

  clearHealthChecksList(): void;
  getHealthChecksList(): Array<gloo_projects_gloo_api_external_envoy_api_v2_core_health_check_pb.HealthCheck>;
  setHealthChecksList(value: Array<gloo_projects_gloo_api_external_envoy_api_v2_core_health_check_pb.HealthCheck>): void;
  addHealthChecks(value?: gloo_projects_gloo_api_external_envoy_api_v2_core_health_check_pb.HealthCheck, index?: number): gloo_projects_gloo_api_external_envoy_api_v2_core_health_check_pb.HealthCheck;

  hasOutlierDetection(): boolean;
  clearOutlierDetection(): void;
  getOutlierDetection(): gloo_projects_gloo_api_external_envoy_api_v2_cluster_outlier_detection_pb.OutlierDetection | undefined;
  setOutlierDetection(value?: gloo_projects_gloo_api_external_envoy_api_v2_cluster_outlier_detection_pb.OutlierDetection): void;

  getUseHttp2(): boolean;
  setUseHttp2(value: boolean): void;

  hasKube(): boolean;
  clearKube(): void;
  getKube(): gloo_projects_gloo_api_v1_options_kubernetes_kubernetes_pb.UpstreamSpec | undefined;
  setKube(value?: gloo_projects_gloo_api_v1_options_kubernetes_kubernetes_pb.UpstreamSpec): void;

  hasStatic(): boolean;
  clearStatic(): void;
  getStatic(): gloo_projects_gloo_api_v1_options_static_static_pb.UpstreamSpec | undefined;
  setStatic(value?: gloo_projects_gloo_api_v1_options_static_static_pb.UpstreamSpec): void;

  hasPipe(): boolean;
  clearPipe(): void;
  getPipe(): gloo_projects_gloo_api_v1_options_pipe_pipe_pb.UpstreamSpec | undefined;
  setPipe(value?: gloo_projects_gloo_api_v1_options_pipe_pipe_pb.UpstreamSpec): void;

  hasAws(): boolean;
  clearAws(): void;
  getAws(): gloo_projects_gloo_api_v1_options_aws_aws_pb.UpstreamSpec | undefined;
  setAws(value?: gloo_projects_gloo_api_v1_options_aws_aws_pb.UpstreamSpec): void;

  hasAzure(): boolean;
  clearAzure(): void;
  getAzure(): gloo_projects_gloo_api_v1_options_azure_azure_pb.UpstreamSpec | undefined;
  setAzure(value?: gloo_projects_gloo_api_v1_options_azure_azure_pb.UpstreamSpec): void;

  hasConsul(): boolean;
  clearConsul(): void;
  getConsul(): gloo_projects_gloo_api_v1_options_consul_consul_pb.UpstreamSpec | undefined;
  setConsul(value?: gloo_projects_gloo_api_v1_options_consul_consul_pb.UpstreamSpec): void;

  hasAwsEc2(): boolean;
  clearAwsEc2(): void;
  getAwsEc2(): gloo_projects_gloo_api_v1_options_aws_ec2_aws_ec2_pb.UpstreamSpec | undefined;
  setAwsEc2(value?: gloo_projects_gloo_api_v1_options_aws_ec2_aws_ec2_pb.UpstreamSpec): void;

  getUpstreamTypeCase(): Upstream.UpstreamTypeCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Upstream.AsObject;
  static toObject(includeInstance: boolean, msg: Upstream): Upstream.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Upstream, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Upstream;
  static deserializeBinaryFromReader(message: Upstream, reader: jspb.BinaryReader): Upstream;
}

export namespace Upstream {
  export type AsObject = {
    status?: solo_kit_api_v1_status_pb.Status.AsObject,
    metadata?: solo_kit_api_v1_metadata_pb.Metadata.AsObject,
    discoveryMetadata?: DiscoveryMetadata.AsObject,
    sslConfig?: gloo_projects_gloo_api_v1_ssl_pb.UpstreamSslConfig.AsObject,
    circuitBreakers?: gloo_projects_gloo_api_v1_circuit_breaker_pb.CircuitBreakerConfig.AsObject,
    loadBalancerConfig?: gloo_projects_gloo_api_v1_load_balancer_pb.LoadBalancerConfig.AsObject,
    connectionConfig?: gloo_projects_gloo_api_v1_connection_pb.ConnectionConfig.AsObject,
    healthChecksList: Array<gloo_projects_gloo_api_external_envoy_api_v2_core_health_check_pb.HealthCheck.AsObject>,
    outlierDetection?: gloo_projects_gloo_api_external_envoy_api_v2_cluster_outlier_detection_pb.OutlierDetection.AsObject,
    useHttp2: boolean,
    kube?: gloo_projects_gloo_api_v1_options_kubernetes_kubernetes_pb.UpstreamSpec.AsObject,
    pb_static?: gloo_projects_gloo_api_v1_options_static_static_pb.UpstreamSpec.AsObject,
    pipe?: gloo_projects_gloo_api_v1_options_pipe_pipe_pb.UpstreamSpec.AsObject,
    aws?: gloo_projects_gloo_api_v1_options_aws_aws_pb.UpstreamSpec.AsObject,
    azure?: gloo_projects_gloo_api_v1_options_azure_azure_pb.UpstreamSpec.AsObject,
    consul?: gloo_projects_gloo_api_v1_options_consul_consul_pb.UpstreamSpec.AsObject,
    awsEc2?: gloo_projects_gloo_api_v1_options_aws_ec2_aws_ec2_pb.UpstreamSpec.AsObject,
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
  }
}

