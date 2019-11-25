// package: gloo.solo.io
// file: github.com/solo-io/gloo/projects/gloo/api/v1/upstream.proto

import * as jspb from "google-protobuf";
import * as gogoproto_gogo_pb from "../../../../../../../gogoproto/gogo_pb";
import * as github_com_solo_io_solo_kit_api_v1_metadata_pb from "../../../../../../../github.com/solo-io/solo-kit/api/v1/metadata_pb";
import * as github_com_solo_io_gloo_projects_gloo_api_v1_ssl_pb from "../../../../../../../github.com/solo-io/gloo/projects/gloo/api/v1/ssl_pb";
import * as github_com_solo_io_gloo_projects_gloo_api_v1_circuit_breaker_pb from "../../../../../../../github.com/solo-io/gloo/projects/gloo/api/v1/circuit_breaker_pb";
import * as github_com_solo_io_gloo_projects_gloo_api_v1_load_balancer_pb from "../../../../../../../github.com/solo-io/gloo/projects/gloo/api/v1/load_balancer_pb";
import * as github_com_solo_io_gloo_projects_gloo_api_v1_connection_pb from "../../../../../../../github.com/solo-io/gloo/projects/gloo/api/v1/connection_pb";
import * as github_com_solo_io_gloo_projects_gloo_api_external_envoy_api_v2_core_health_check_pb from "../../../../../../../github.com/solo-io/gloo/projects/gloo/api/external/envoy/api/v2/core/health_check_pb";
import * as github_com_solo_io_solo_kit_api_v1_status_pb from "../../../../../../../github.com/solo-io/solo-kit/api/v1/status_pb";
import * as github_com_solo_io_gloo_projects_gloo_api_external_envoy_api_v2_cluster_outlier_detection_pb from "../../../../../../../github.com/solo-io/gloo/projects/gloo/api/external/envoy/api/v2/cluster/outlier_detection_pb";
import * as github_com_solo_io_solo_kit_api_v1_solo_kit_pb from "../../../../../../../github.com/solo-io/solo-kit/api/v1/solo-kit_pb";
import * as github_com_solo_io_gloo_projects_gloo_api_v1_options_static_static_pb from "../../../../../../../github.com/solo-io/gloo/projects/gloo/api/v1/options/static/static_pb";
import * as github_com_solo_io_gloo_projects_gloo_api_v1_options_pipe_pipe_pb from "../../../../../../../github.com/solo-io/gloo/projects/gloo/api/v1/options/pipe/pipe_pb";
import * as github_com_solo_io_gloo_projects_gloo_api_v1_options_kubernetes_kubernetes_pb from "../../../../../../../github.com/solo-io/gloo/projects/gloo/api/v1/options/kubernetes/kubernetes_pb";
import * as github_com_solo_io_gloo_projects_gloo_api_v1_options_aws_aws_pb from "../../../../../../../github.com/solo-io/gloo/projects/gloo/api/v1/options/aws/aws_pb";
import * as github_com_solo_io_gloo_projects_gloo_api_v1_options_azure_azure_pb from "../../../../../../../github.com/solo-io/gloo/projects/gloo/api/v1/options/azure/azure_pb";
import * as github_com_solo_io_gloo_projects_gloo_api_v1_options_consul_consul_pb from "../../../../../../../github.com/solo-io/gloo/projects/gloo/api/v1/options/consul/consul_pb";
import * as github_com_solo_io_gloo_projects_gloo_api_v1_options_aws_ec2_aws_ec2_pb from "../../../../../../../github.com/solo-io/gloo/projects/gloo/api/v1/options/aws/ec2/aws_ec2_pb";
import * as github_com_solo_io_gloo_projects_gloo_api_v1_options_pb from "../../../../../../../github.com/solo-io/gloo/projects/gloo/api/v1/options_pb";

export class Upstream extends jspb.Message {
  hasStatus(): boolean;
  clearStatus(): void;
  getStatus(): github_com_solo_io_solo_kit_api_v1_status_pb.Status | undefined;
  setStatus(value?: github_com_solo_io_solo_kit_api_v1_status_pb.Status): void;

  hasMetadata(): boolean;
  clearMetadata(): void;
  getMetadata(): github_com_solo_io_solo_kit_api_v1_metadata_pb.Metadata | undefined;
  setMetadata(value?: github_com_solo_io_solo_kit_api_v1_metadata_pb.Metadata): void;

  hasDiscoveryMetadata(): boolean;
  clearDiscoveryMetadata(): void;
  getDiscoveryMetadata(): DiscoveryMetadata | undefined;
  setDiscoveryMetadata(value?: DiscoveryMetadata): void;

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

  clearHealthChecksList(): void;
  getHealthChecksList(): Array<github_com_solo_io_gloo_projects_gloo_api_external_envoy_api_v2_core_health_check_pb.HealthCheck>;
  setHealthChecksList(value: Array<github_com_solo_io_gloo_projects_gloo_api_external_envoy_api_v2_core_health_check_pb.HealthCheck>): void;
  addHealthChecks(value?: github_com_solo_io_gloo_projects_gloo_api_external_envoy_api_v2_core_health_check_pb.HealthCheck, index?: number): github_com_solo_io_gloo_projects_gloo_api_external_envoy_api_v2_core_health_check_pb.HealthCheck;

  hasOutlierDetection(): boolean;
  clearOutlierDetection(): void;
  getOutlierDetection(): github_com_solo_io_gloo_projects_gloo_api_external_envoy_api_v2_cluster_outlier_detection_pb.OutlierDetection | undefined;
  setOutlierDetection(value?: github_com_solo_io_gloo_projects_gloo_api_external_envoy_api_v2_cluster_outlier_detection_pb.OutlierDetection): void;

  getUseHttp2(): boolean;
  setUseHttp2(value: boolean): void;

  hasKube(): boolean;
  clearKube(): void;
  getKube(): github_com_solo_io_gloo_projects_gloo_api_v1_options_kubernetes_kubernetes_pb.UpstreamSpec | undefined;
  setKube(value?: github_com_solo_io_gloo_projects_gloo_api_v1_options_kubernetes_kubernetes_pb.UpstreamSpec): void;

  hasStatic(): boolean;
  clearStatic(): void;
  getStatic(): github_com_solo_io_gloo_projects_gloo_api_v1_options_static_static_pb.UpstreamSpec | undefined;
  setStatic(value?: github_com_solo_io_gloo_projects_gloo_api_v1_options_static_static_pb.UpstreamSpec): void;

  hasPipe(): boolean;
  clearPipe(): void;
  getPipe(): github_com_solo_io_gloo_projects_gloo_api_v1_options_pipe_pipe_pb.UpstreamSpec | undefined;
  setPipe(value?: github_com_solo_io_gloo_projects_gloo_api_v1_options_pipe_pipe_pb.UpstreamSpec): void;

  hasAws(): boolean;
  clearAws(): void;
  getAws(): github_com_solo_io_gloo_projects_gloo_api_v1_options_aws_aws_pb.UpstreamSpec | undefined;
  setAws(value?: github_com_solo_io_gloo_projects_gloo_api_v1_options_aws_aws_pb.UpstreamSpec): void;

  hasAzure(): boolean;
  clearAzure(): void;
  getAzure(): github_com_solo_io_gloo_projects_gloo_api_v1_options_azure_azure_pb.UpstreamSpec | undefined;
  setAzure(value?: github_com_solo_io_gloo_projects_gloo_api_v1_options_azure_azure_pb.UpstreamSpec): void;

  hasConsul(): boolean;
  clearConsul(): void;
  getConsul(): github_com_solo_io_gloo_projects_gloo_api_v1_options_consul_consul_pb.UpstreamSpec | undefined;
  setConsul(value?: github_com_solo_io_gloo_projects_gloo_api_v1_options_consul_consul_pb.UpstreamSpec): void;

  hasAwsEc2(): boolean;
  clearAwsEc2(): void;
  getAwsEc2(): github_com_solo_io_gloo_projects_gloo_api_v1_options_aws_ec2_aws_ec2_pb.UpstreamSpec | undefined;
  setAwsEc2(value?: github_com_solo_io_gloo_projects_gloo_api_v1_options_aws_ec2_aws_ec2_pb.UpstreamSpec): void;

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
    status?: github_com_solo_io_solo_kit_api_v1_status_pb.Status.AsObject,
    metadata?: github_com_solo_io_solo_kit_api_v1_metadata_pb.Metadata.AsObject,
    discoveryMetadata?: DiscoveryMetadata.AsObject,
    sslConfig?: github_com_solo_io_gloo_projects_gloo_api_v1_ssl_pb.UpstreamSslConfig.AsObject,
    circuitBreakers?: github_com_solo_io_gloo_projects_gloo_api_v1_circuit_breaker_pb.CircuitBreakerConfig.AsObject,
    loadBalancerConfig?: github_com_solo_io_gloo_projects_gloo_api_v1_load_balancer_pb.LoadBalancerConfig.AsObject,
    connectionConfig?: github_com_solo_io_gloo_projects_gloo_api_v1_connection_pb.ConnectionConfig.AsObject,
    healthChecksList: Array<github_com_solo_io_gloo_projects_gloo_api_external_envoy_api_v2_core_health_check_pb.HealthCheck.AsObject>,
    outlierDetection?: github_com_solo_io_gloo_projects_gloo_api_external_envoy_api_v2_cluster_outlier_detection_pb.OutlierDetection.AsObject,
    useHttp2: boolean,
    kube?: github_com_solo_io_gloo_projects_gloo_api_v1_options_kubernetes_kubernetes_pb.UpstreamSpec.AsObject,
    pb_static?: github_com_solo_io_gloo_projects_gloo_api_v1_options_static_static_pb.UpstreamSpec.AsObject,
    pipe?: github_com_solo_io_gloo_projects_gloo_api_v1_options_pipe_pipe_pb.UpstreamSpec.AsObject,
    aws?: github_com_solo_io_gloo_projects_gloo_api_v1_options_aws_aws_pb.UpstreamSpec.AsObject,
    azure?: github_com_solo_io_gloo_projects_gloo_api_v1_options_azure_azure_pb.UpstreamSpec.AsObject,
    consul?: github_com_solo_io_gloo_projects_gloo_api_v1_options_consul_consul_pb.UpstreamSpec.AsObject,
    awsEc2?: github_com_solo_io_gloo_projects_gloo_api_v1_options_aws_ec2_aws_ec2_pb.UpstreamSpec.AsObject,
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

