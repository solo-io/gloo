/* eslint-disable */
// package: gloo.solo.io
// file: github.com/solo-io/solo-apis/api/gloo/gloo/v1/settings.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_struct_pb from "google-protobuf/google/protobuf/struct_pb";
import * as extproto_ext_pb from "../../../../../../../extproto/ext_pb";
import * as github_com_solo_io_solo_kit_api_v1_solo_kit_pb from "../../../../../../../github.com/solo-io/solo-kit/api/v1/solo-kit_pb";
import * as github_com_solo_io_solo_kit_api_v1_ref_pb from "../../../../../../../github.com/solo-io/solo-kit/api/v1/ref_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_extensions_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/extensions_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_ratelimit_ratelimit_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/enterprise/options/ratelimit/ratelimit_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_caching_caching_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/enterprise/options/caching/caching_pb";
import * as github_com_solo_io_solo_apis_api_gloo_enterprise_gloo_v1_auth_config_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/enterprise.gloo/v1/auth_config_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_rbac_rbac_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/enterprise/options/rbac/rbac_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_circuit_breaker_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/circuit_breaker_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_ssl_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/ssl_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_extensions_aws_filter_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/extensions/aws/filter_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_consul_query_options_pb from "../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/options/consul/query_options_pb";
import * as google_protobuf_duration_pb from "google-protobuf/google/protobuf/duration_pb";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";

export class SettingsSpec extends jspb.Message {
  getDiscoveryNamespace(): string;
  setDiscoveryNamespace(value: string): void;

  clearWatchNamespacesList(): void;
  getWatchNamespacesList(): Array<string>;
  setWatchNamespacesList(value: Array<string>): void;
  addWatchNamespaces(value: string, index?: number): string;

  hasKubernetesConfigSource(): boolean;
  clearKubernetesConfigSource(): void;
  getKubernetesConfigSource(): SettingsSpec.KubernetesCrds | undefined;
  setKubernetesConfigSource(value?: SettingsSpec.KubernetesCrds): void;

  hasDirectoryConfigSource(): boolean;
  clearDirectoryConfigSource(): void;
  getDirectoryConfigSource(): SettingsSpec.Directory | undefined;
  setDirectoryConfigSource(value?: SettingsSpec.Directory): void;

  hasConsulKvSource(): boolean;
  clearConsulKvSource(): void;
  getConsulKvSource(): SettingsSpec.ConsulKv | undefined;
  setConsulKvSource(value?: SettingsSpec.ConsulKv): void;

  hasKubernetesSecretSource(): boolean;
  clearKubernetesSecretSource(): void;
  getKubernetesSecretSource(): SettingsSpec.KubernetesSecrets | undefined;
  setKubernetesSecretSource(value?: SettingsSpec.KubernetesSecrets): void;

  hasVaultSecretSource(): boolean;
  clearVaultSecretSource(): void;
  getVaultSecretSource(): SettingsSpec.VaultSecrets | undefined;
  setVaultSecretSource(value?: SettingsSpec.VaultSecrets): void;

  hasDirectorySecretSource(): boolean;
  clearDirectorySecretSource(): void;
  getDirectorySecretSource(): SettingsSpec.Directory | undefined;
  setDirectorySecretSource(value?: SettingsSpec.Directory): void;

  hasKubernetesArtifactSource(): boolean;
  clearKubernetesArtifactSource(): void;
  getKubernetesArtifactSource(): SettingsSpec.KubernetesConfigmaps | undefined;
  setKubernetesArtifactSource(value?: SettingsSpec.KubernetesConfigmaps): void;

  hasDirectoryArtifactSource(): boolean;
  clearDirectoryArtifactSource(): void;
  getDirectoryArtifactSource(): SettingsSpec.Directory | undefined;
  setDirectoryArtifactSource(value?: SettingsSpec.Directory): void;

  hasConsulKvArtifactSource(): boolean;
  clearConsulKvArtifactSource(): void;
  getConsulKvArtifactSource(): SettingsSpec.ConsulKv | undefined;
  setConsulKvArtifactSource(value?: SettingsSpec.ConsulKv): void;

  hasRefreshRate(): boolean;
  clearRefreshRate(): void;
  getRefreshRate(): google_protobuf_duration_pb.Duration | undefined;
  setRefreshRate(value?: google_protobuf_duration_pb.Duration): void;

  getDevMode(): boolean;
  setDevMode(value: boolean): void;

  getLinkerd(): boolean;
  setLinkerd(value: boolean): void;

  hasKnative(): boolean;
  clearKnative(): void;
  getKnative(): SettingsSpec.KnativeOptions | undefined;
  setKnative(value?: SettingsSpec.KnativeOptions): void;

  hasDiscovery(): boolean;
  clearDiscovery(): void;
  getDiscovery(): SettingsSpec.DiscoveryOptions | undefined;
  setDiscovery(value?: SettingsSpec.DiscoveryOptions): void;

  hasGloo(): boolean;
  clearGloo(): void;
  getGloo(): GlooOptions | undefined;
  setGloo(value?: GlooOptions): void;

  hasGateway(): boolean;
  clearGateway(): void;
  getGateway(): GatewayOptions | undefined;
  setGateway(value?: GatewayOptions): void;

  hasConsul(): boolean;
  clearConsul(): void;
  getConsul(): SettingsSpec.ConsulConfiguration | undefined;
  setConsul(value?: SettingsSpec.ConsulConfiguration): void;

  hasConsuldiscovery(): boolean;
  clearConsuldiscovery(): void;
  getConsuldiscovery(): SettingsSpec.ConsulUpstreamDiscoveryConfiguration | undefined;
  setConsuldiscovery(value?: SettingsSpec.ConsulUpstreamDiscoveryConfiguration): void;

  hasKubernetes(): boolean;
  clearKubernetes(): void;
  getKubernetes(): SettingsSpec.KubernetesConfiguration | undefined;
  setKubernetes(value?: SettingsSpec.KubernetesConfiguration): void;

  hasExtensions(): boolean;
  clearExtensions(): void;
  getExtensions(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_extensions_pb.Extensions | undefined;
  setExtensions(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_extensions_pb.Extensions): void;

  hasRatelimit(): boolean;
  clearRatelimit(): void;
  getRatelimit(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_ratelimit_ratelimit_pb.ServiceSettings | undefined;
  setRatelimit(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_ratelimit_ratelimit_pb.ServiceSettings): void;

  hasRatelimitServer(): boolean;
  clearRatelimitServer(): void;
  getRatelimitServer(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_ratelimit_ratelimit_pb.Settings | undefined;
  setRatelimitServer(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_ratelimit_ratelimit_pb.Settings): void;

  hasRbac(): boolean;
  clearRbac(): void;
  getRbac(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_rbac_rbac_pb.Settings | undefined;
  setRbac(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_rbac_rbac_pb.Settings): void;

  hasExtauth(): boolean;
  clearExtauth(): void;
  getExtauth(): github_com_solo_io_solo_apis_api_gloo_enterprise_gloo_v1_auth_config_pb.Settings | undefined;
  setExtauth(value?: github_com_solo_io_solo_apis_api_gloo_enterprise_gloo_v1_auth_config_pb.Settings): void;

  getNamedExtauthMap(): jspb.Map<string, github_com_solo_io_solo_apis_api_gloo_enterprise_gloo_v1_auth_config_pb.Settings>;
  clearNamedExtauthMap(): void;
  hasCachingServer(): boolean;
  clearCachingServer(): void;
  getCachingServer(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_caching_caching_pb.Settings | undefined;
  setCachingServer(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_caching_caching_pb.Settings): void;

  hasObservabilityoptions(): boolean;
  clearObservabilityoptions(): void;
  getObservabilityoptions(): SettingsSpec.ObservabilityOptions | undefined;
  setObservabilityoptions(value?: SettingsSpec.ObservabilityOptions): void;

  hasUpstreamoptions(): boolean;
  clearUpstreamoptions(): void;
  getUpstreamoptions(): UpstreamOptions | undefined;
  setUpstreamoptions(value?: UpstreamOptions): void;

  hasConsoleOptions(): boolean;
  clearConsoleOptions(): void;
  getConsoleOptions(): ConsoleOptions | undefined;
  setConsoleOptions(value?: ConsoleOptions): void;

  hasGraphqlOptions(): boolean;
  clearGraphqlOptions(): void;
  getGraphqlOptions(): GraphqlOptions | undefined;
  setGraphqlOptions(value?: GraphqlOptions): void;

  getConfigSourceCase(): SettingsSpec.ConfigSourceCase;
  getSecretSourceCase(): SettingsSpec.SecretSourceCase;
  getArtifactSourceCase(): SettingsSpec.ArtifactSourceCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SettingsSpec.AsObject;
  static toObject(includeInstance: boolean, msg: SettingsSpec): SettingsSpec.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: SettingsSpec, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SettingsSpec;
  static deserializeBinaryFromReader(message: SettingsSpec, reader: jspb.BinaryReader): SettingsSpec;
}

export namespace SettingsSpec {
  export type AsObject = {
    discoveryNamespace: string,
    watchNamespacesList: Array<string>,
    kubernetesConfigSource?: SettingsSpec.KubernetesCrds.AsObject,
    directoryConfigSource?: SettingsSpec.Directory.AsObject,
    consulKvSource?: SettingsSpec.ConsulKv.AsObject,
    kubernetesSecretSource?: SettingsSpec.KubernetesSecrets.AsObject,
    vaultSecretSource?: SettingsSpec.VaultSecrets.AsObject,
    directorySecretSource?: SettingsSpec.Directory.AsObject,
    kubernetesArtifactSource?: SettingsSpec.KubernetesConfigmaps.AsObject,
    directoryArtifactSource?: SettingsSpec.Directory.AsObject,
    consulKvArtifactSource?: SettingsSpec.ConsulKv.AsObject,
    refreshRate?: google_protobuf_duration_pb.Duration.AsObject,
    devMode: boolean,
    linkerd: boolean,
    knative?: SettingsSpec.KnativeOptions.AsObject,
    discovery?: SettingsSpec.DiscoveryOptions.AsObject,
    gloo?: GlooOptions.AsObject,
    gateway?: GatewayOptions.AsObject,
    consul?: SettingsSpec.ConsulConfiguration.AsObject,
    consuldiscovery?: SettingsSpec.ConsulUpstreamDiscoveryConfiguration.AsObject,
    kubernetes?: SettingsSpec.KubernetesConfiguration.AsObject,
    extensions?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_extensions_pb.Extensions.AsObject,
    ratelimit?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_ratelimit_ratelimit_pb.ServiceSettings.AsObject,
    ratelimitServer?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_ratelimit_ratelimit_pb.Settings.AsObject,
    rbac?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_rbac_rbac_pb.Settings.AsObject,
    extauth?: github_com_solo_io_solo_apis_api_gloo_enterprise_gloo_v1_auth_config_pb.Settings.AsObject,
    namedExtauthMap: Array<[string, github_com_solo_io_solo_apis_api_gloo_enterprise_gloo_v1_auth_config_pb.Settings.AsObject]>,
    cachingServer?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_options_caching_caching_pb.Settings.AsObject,
    observabilityoptions?: SettingsSpec.ObservabilityOptions.AsObject,
    upstreamoptions?: UpstreamOptions.AsObject,
    consoleOptions?: ConsoleOptions.AsObject,
    graphqlOptions?: GraphqlOptions.AsObject,
  }

  export class KubernetesCrds extends jspb.Message {
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): KubernetesCrds.AsObject;
    static toObject(includeInstance: boolean, msg: KubernetesCrds): KubernetesCrds.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: KubernetesCrds, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): KubernetesCrds;
    static deserializeBinaryFromReader(message: KubernetesCrds, reader: jspb.BinaryReader): KubernetesCrds;
  }

  export namespace KubernetesCrds {
    export type AsObject = {
    }
  }

  export class KubernetesSecrets extends jspb.Message {
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): KubernetesSecrets.AsObject;
    static toObject(includeInstance: boolean, msg: KubernetesSecrets): KubernetesSecrets.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: KubernetesSecrets, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): KubernetesSecrets;
    static deserializeBinaryFromReader(message: KubernetesSecrets, reader: jspb.BinaryReader): KubernetesSecrets;
  }

  export namespace KubernetesSecrets {
    export type AsObject = {
    }
  }

  export class VaultSecrets extends jspb.Message {
    getToken(): string;
    setToken(value: string): void;

    getAddress(): string;
    setAddress(value: string): void;

    getCaCert(): string;
    setCaCert(value: string): void;

    getCaPath(): string;
    setCaPath(value: string): void;

    getClientCert(): string;
    setClientCert(value: string): void;

    getClientKey(): string;
    setClientKey(value: string): void;

    getTlsServerName(): string;
    setTlsServerName(value: string): void;

    hasInsecure(): boolean;
    clearInsecure(): void;
    getInsecure(): google_protobuf_wrappers_pb.BoolValue | undefined;
    setInsecure(value?: google_protobuf_wrappers_pb.BoolValue): void;

    getRootKey(): string;
    setRootKey(value: string): void;

    getPathPrefix(): string;
    setPathPrefix(value: string): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): VaultSecrets.AsObject;
    static toObject(includeInstance: boolean, msg: VaultSecrets): VaultSecrets.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: VaultSecrets, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): VaultSecrets;
    static deserializeBinaryFromReader(message: VaultSecrets, reader: jspb.BinaryReader): VaultSecrets;
  }

  export namespace VaultSecrets {
    export type AsObject = {
      token: string,
      address: string,
      caCert: string,
      caPath: string,
      clientCert: string,
      clientKey: string,
      tlsServerName: string,
      insecure?: google_protobuf_wrappers_pb.BoolValue.AsObject,
      rootKey: string,
      pathPrefix: string,
    }
  }

  export class ConsulKv extends jspb.Message {
    getRootKey(): string;
    setRootKey(value: string): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): ConsulKv.AsObject;
    static toObject(includeInstance: boolean, msg: ConsulKv): ConsulKv.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: ConsulKv, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): ConsulKv;
    static deserializeBinaryFromReader(message: ConsulKv, reader: jspb.BinaryReader): ConsulKv;
  }

  export namespace ConsulKv {
    export type AsObject = {
      rootKey: string,
    }
  }

  export class KubernetesConfigmaps extends jspb.Message {
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): KubernetesConfigmaps.AsObject;
    static toObject(includeInstance: boolean, msg: KubernetesConfigmaps): KubernetesConfigmaps.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: KubernetesConfigmaps, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): KubernetesConfigmaps;
    static deserializeBinaryFromReader(message: KubernetesConfigmaps, reader: jspb.BinaryReader): KubernetesConfigmaps;
  }

  export namespace KubernetesConfigmaps {
    export type AsObject = {
    }
  }

  export class Directory extends jspb.Message {
    getDirectory(): string;
    setDirectory(value: string): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Directory.AsObject;
    static toObject(includeInstance: boolean, msg: Directory): Directory.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: Directory, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Directory;
    static deserializeBinaryFromReader(message: Directory, reader: jspb.BinaryReader): Directory;
  }

  export namespace Directory {
    export type AsObject = {
      directory: string,
    }
  }

  export class KnativeOptions extends jspb.Message {
    getClusterIngressProxyAddress(): string;
    setClusterIngressProxyAddress(value: string): void;

    getKnativeExternalProxyAddress(): string;
    setKnativeExternalProxyAddress(value: string): void;

    getKnativeInternalProxyAddress(): string;
    setKnativeInternalProxyAddress(value: string): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): KnativeOptions.AsObject;
    static toObject(includeInstance: boolean, msg: KnativeOptions): KnativeOptions.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: KnativeOptions, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): KnativeOptions;
    static deserializeBinaryFromReader(message: KnativeOptions, reader: jspb.BinaryReader): KnativeOptions;
  }

  export namespace KnativeOptions {
    export type AsObject = {
      clusterIngressProxyAddress: string,
      knativeExternalProxyAddress: string,
      knativeInternalProxyAddress: string,
    }
  }

  export class DiscoveryOptions extends jspb.Message {
    getFdsMode(): SettingsSpec.DiscoveryOptions.FdsModeMap[keyof SettingsSpec.DiscoveryOptions.FdsModeMap];
    setFdsMode(value: SettingsSpec.DiscoveryOptions.FdsModeMap[keyof SettingsSpec.DiscoveryOptions.FdsModeMap]): void;

    hasUdsOptions(): boolean;
    clearUdsOptions(): void;
    getUdsOptions(): SettingsSpec.DiscoveryOptions.UdsOptions | undefined;
    setUdsOptions(value?: SettingsSpec.DiscoveryOptions.UdsOptions): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): DiscoveryOptions.AsObject;
    static toObject(includeInstance: boolean, msg: DiscoveryOptions): DiscoveryOptions.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: DiscoveryOptions, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): DiscoveryOptions;
    static deserializeBinaryFromReader(message: DiscoveryOptions, reader: jspb.BinaryReader): DiscoveryOptions;
  }

  export namespace DiscoveryOptions {
    export type AsObject = {
      fdsMode: SettingsSpec.DiscoveryOptions.FdsModeMap[keyof SettingsSpec.DiscoveryOptions.FdsModeMap],
      udsOptions?: SettingsSpec.DiscoveryOptions.UdsOptions.AsObject,
    }

    export class UdsOptions extends jspb.Message {
      hasEnabled(): boolean;
      clearEnabled(): void;
      getEnabled(): google_protobuf_wrappers_pb.BoolValue | undefined;
      setEnabled(value?: google_protobuf_wrappers_pb.BoolValue): void;

      getWatchLabelsMap(): jspb.Map<string, string>;
      clearWatchLabelsMap(): void;
      serializeBinary(): Uint8Array;
      toObject(includeInstance?: boolean): UdsOptions.AsObject;
      static toObject(includeInstance: boolean, msg: UdsOptions): UdsOptions.AsObject;
      static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
      static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
      static serializeBinaryToWriter(message: UdsOptions, writer: jspb.BinaryWriter): void;
      static deserializeBinary(bytes: Uint8Array): UdsOptions;
      static deserializeBinaryFromReader(message: UdsOptions, reader: jspb.BinaryReader): UdsOptions;
    }

    export namespace UdsOptions {
      export type AsObject = {
        enabled?: google_protobuf_wrappers_pb.BoolValue.AsObject,
        watchLabelsMap: Array<[string, string]>,
      }
    }

    export interface FdsModeMap {
      BLACKLIST: 0;
      WHITELIST: 1;
      DISABLED: 2;
    }

    export const FdsMode: FdsModeMap;
  }

  export class ConsulConfiguration extends jspb.Message {
    getAddress(): string;
    setAddress(value: string): void;

    getDatacenter(): string;
    setDatacenter(value: string): void;

    getUsername(): string;
    setUsername(value: string): void;

    getPassword(): string;
    setPassword(value: string): void;

    getToken(): string;
    setToken(value: string): void;

    getCaFile(): string;
    setCaFile(value: string): void;

    getCaPath(): string;
    setCaPath(value: string): void;

    getCertFile(): string;
    setCertFile(value: string): void;

    getKeyFile(): string;
    setKeyFile(value: string): void;

    hasInsecureSkipVerify(): boolean;
    clearInsecureSkipVerify(): void;
    getInsecureSkipVerify(): google_protobuf_wrappers_pb.BoolValue | undefined;
    setInsecureSkipVerify(value?: google_protobuf_wrappers_pb.BoolValue): void;

    hasWaitTime(): boolean;
    clearWaitTime(): void;
    getWaitTime(): google_protobuf_duration_pb.Duration | undefined;
    setWaitTime(value?: google_protobuf_duration_pb.Duration): void;

    hasServiceDiscovery(): boolean;
    clearServiceDiscovery(): void;
    getServiceDiscovery(): SettingsSpec.ConsulConfiguration.ServiceDiscoveryOptions | undefined;
    setServiceDiscovery(value?: SettingsSpec.ConsulConfiguration.ServiceDiscoveryOptions): void;

    getHttpAddress(): string;
    setHttpAddress(value: string): void;

    getDnsAddress(): string;
    setDnsAddress(value: string): void;

    hasDnsPollingInterval(): boolean;
    clearDnsPollingInterval(): void;
    getDnsPollingInterval(): google_protobuf_duration_pb.Duration | undefined;
    setDnsPollingInterval(value?: google_protobuf_duration_pb.Duration): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): ConsulConfiguration.AsObject;
    static toObject(includeInstance: boolean, msg: ConsulConfiguration): ConsulConfiguration.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: ConsulConfiguration, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): ConsulConfiguration;
    static deserializeBinaryFromReader(message: ConsulConfiguration, reader: jspb.BinaryReader): ConsulConfiguration;
  }

  export namespace ConsulConfiguration {
    export type AsObject = {
      address: string,
      datacenter: string,
      username: string,
      password: string,
      token: string,
      caFile: string,
      caPath: string,
      certFile: string,
      keyFile: string,
      insecureSkipVerify?: google_protobuf_wrappers_pb.BoolValue.AsObject,
      waitTime?: google_protobuf_duration_pb.Duration.AsObject,
      serviceDiscovery?: SettingsSpec.ConsulConfiguration.ServiceDiscoveryOptions.AsObject,
      httpAddress: string,
      dnsAddress: string,
      dnsPollingInterval?: google_protobuf_duration_pb.Duration.AsObject,
    }

    export class ServiceDiscoveryOptions extends jspb.Message {
      clearDataCentersList(): void;
      getDataCentersList(): Array<string>;
      setDataCentersList(value: Array<string>): void;
      addDataCenters(value: string, index?: number): string;

      serializeBinary(): Uint8Array;
      toObject(includeInstance?: boolean): ServiceDiscoveryOptions.AsObject;
      static toObject(includeInstance: boolean, msg: ServiceDiscoveryOptions): ServiceDiscoveryOptions.AsObject;
      static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
      static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
      static serializeBinaryToWriter(message: ServiceDiscoveryOptions, writer: jspb.BinaryWriter): void;
      static deserializeBinary(bytes: Uint8Array): ServiceDiscoveryOptions;
      static deserializeBinaryFromReader(message: ServiceDiscoveryOptions, reader: jspb.BinaryReader): ServiceDiscoveryOptions;
    }

    export namespace ServiceDiscoveryOptions {
      export type AsObject = {
        dataCentersList: Array<string>,
      }
    }
  }

  export class ConsulUpstreamDiscoveryConfiguration extends jspb.Message {
    getUsetlstagging(): boolean;
    setUsetlstagging(value: boolean): void;

    getTlstagname(): string;
    setTlstagname(value: string): void;

    hasRootca(): boolean;
    clearRootca(): void;
    getRootca(): github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef | undefined;
    setRootca(value?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef): void;

    getSplittlsservices(): boolean;
    setSplittlsservices(value: boolean): void;

    getConsistencymode(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_consul_query_options_pb.ConsulConsistencyModesMap[keyof github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_consul_query_options_pb.ConsulConsistencyModesMap];
    setConsistencymode(value: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_consul_query_options_pb.ConsulConsistencyModesMap[keyof github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_consul_query_options_pb.ConsulConsistencyModesMap]): void;

    hasQueryOptions(): boolean;
    clearQueryOptions(): void;
    getQueryOptions(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_consul_query_options_pb.QueryOptions | undefined;
    setQueryOptions(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_consul_query_options_pb.QueryOptions): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): ConsulUpstreamDiscoveryConfiguration.AsObject;
    static toObject(includeInstance: boolean, msg: ConsulUpstreamDiscoveryConfiguration): ConsulUpstreamDiscoveryConfiguration.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: ConsulUpstreamDiscoveryConfiguration, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): ConsulUpstreamDiscoveryConfiguration;
    static deserializeBinaryFromReader(message: ConsulUpstreamDiscoveryConfiguration, reader: jspb.BinaryReader): ConsulUpstreamDiscoveryConfiguration;
  }

  export namespace ConsulUpstreamDiscoveryConfiguration {
    export type AsObject = {
      usetlstagging: boolean,
      tlstagname: string,
      rootca?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef.AsObject,
      splittlsservices: boolean,
      consistencymode: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_consul_query_options_pb.ConsulConsistencyModesMap[keyof github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_consul_query_options_pb.ConsulConsistencyModesMap],
      queryOptions?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_consul_query_options_pb.QueryOptions.AsObject,
    }
  }

  export class KubernetesConfiguration extends jspb.Message {
    hasRateLimits(): boolean;
    clearRateLimits(): void;
    getRateLimits(): SettingsSpec.KubernetesConfiguration.RateLimits | undefined;
    setRateLimits(value?: SettingsSpec.KubernetesConfiguration.RateLimits): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): KubernetesConfiguration.AsObject;
    static toObject(includeInstance: boolean, msg: KubernetesConfiguration): KubernetesConfiguration.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: KubernetesConfiguration, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): KubernetesConfiguration;
    static deserializeBinaryFromReader(message: KubernetesConfiguration, reader: jspb.BinaryReader): KubernetesConfiguration;
  }

  export namespace KubernetesConfiguration {
    export type AsObject = {
      rateLimits?: SettingsSpec.KubernetesConfiguration.RateLimits.AsObject,
    }

    export class RateLimits extends jspb.Message {
      getQps(): number;
      setQps(value: number): void;

      getBurst(): number;
      setBurst(value: number): void;

      serializeBinary(): Uint8Array;
      toObject(includeInstance?: boolean): RateLimits.AsObject;
      static toObject(includeInstance: boolean, msg: RateLimits): RateLimits.AsObject;
      static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
      static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
      static serializeBinaryToWriter(message: RateLimits, writer: jspb.BinaryWriter): void;
      static deserializeBinary(bytes: Uint8Array): RateLimits;
      static deserializeBinaryFromReader(message: RateLimits, reader: jspb.BinaryReader): RateLimits;
    }

    export namespace RateLimits {
      export type AsObject = {
        qps: number,
        burst: number,
      }
    }
  }

  export class ObservabilityOptions extends jspb.Message {
    hasGrafanaintegration(): boolean;
    clearGrafanaintegration(): void;
    getGrafanaintegration(): SettingsSpec.ObservabilityOptions.GrafanaIntegration | undefined;
    setGrafanaintegration(value?: SettingsSpec.ObservabilityOptions.GrafanaIntegration): void;

    getConfigstatusmetriclabelsMap(): jspb.Map<string, SettingsSpec.ObservabilityOptions.MetricLabels>;
    clearConfigstatusmetriclabelsMap(): void;
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): ObservabilityOptions.AsObject;
    static toObject(includeInstance: boolean, msg: ObservabilityOptions): ObservabilityOptions.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: ObservabilityOptions, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): ObservabilityOptions;
    static deserializeBinaryFromReader(message: ObservabilityOptions, reader: jspb.BinaryReader): ObservabilityOptions;
  }

  export namespace ObservabilityOptions {
    export type AsObject = {
      grafanaintegration?: SettingsSpec.ObservabilityOptions.GrafanaIntegration.AsObject,
      configstatusmetriclabelsMap: Array<[string, SettingsSpec.ObservabilityOptions.MetricLabels.AsObject]>,
    }

    export class GrafanaIntegration extends jspb.Message {
      hasDefaultDashboardFolderId(): boolean;
      clearDefaultDashboardFolderId(): void;
      getDefaultDashboardFolderId(): google_protobuf_wrappers_pb.UInt32Value | undefined;
      setDefaultDashboardFolderId(value?: google_protobuf_wrappers_pb.UInt32Value): void;

      serializeBinary(): Uint8Array;
      toObject(includeInstance?: boolean): GrafanaIntegration.AsObject;
      static toObject(includeInstance: boolean, msg: GrafanaIntegration): GrafanaIntegration.AsObject;
      static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
      static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
      static serializeBinaryToWriter(message: GrafanaIntegration, writer: jspb.BinaryWriter): void;
      static deserializeBinary(bytes: Uint8Array): GrafanaIntegration;
      static deserializeBinaryFromReader(message: GrafanaIntegration, reader: jspb.BinaryReader): GrafanaIntegration;
    }

    export namespace GrafanaIntegration {
      export type AsObject = {
        defaultDashboardFolderId?: google_protobuf_wrappers_pb.UInt32Value.AsObject,
      }
    }

    export class MetricLabels extends jspb.Message {
      getLabeltopathMap(): jspb.Map<string, string>;
      clearLabeltopathMap(): void;
      serializeBinary(): Uint8Array;
      toObject(includeInstance?: boolean): MetricLabels.AsObject;
      static toObject(includeInstance: boolean, msg: MetricLabels): MetricLabels.AsObject;
      static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
      static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
      static serializeBinaryToWriter(message: MetricLabels, writer: jspb.BinaryWriter): void;
      static deserializeBinary(bytes: Uint8Array): MetricLabels;
      static deserializeBinaryFromReader(message: MetricLabels, reader: jspb.BinaryReader): MetricLabels;
    }

    export namespace MetricLabels {
      export type AsObject = {
        labeltopathMap: Array<[string, string]>,
      }
    }
  }

  export enum ConfigSourceCase {
    CONFIG_SOURCE_NOT_SET = 0,
    KUBERNETES_CONFIG_SOURCE = 4,
    DIRECTORY_CONFIG_SOURCE = 5,
    CONSUL_KV_SOURCE = 21,
  }

  export enum SecretSourceCase {
    SECRET_SOURCE_NOT_SET = 0,
    KUBERNETES_SECRET_SOURCE = 6,
    VAULT_SECRET_SOURCE = 7,
    DIRECTORY_SECRET_SOURCE = 8,
  }

  export enum ArtifactSourceCase {
    ARTIFACT_SOURCE_NOT_SET = 0,
    KUBERNETES_ARTIFACT_SOURCE = 9,
    DIRECTORY_ARTIFACT_SOURCE = 10,
    CONSUL_KV_ARTIFACT_SOURCE = 23,
  }
}

export class UpstreamOptions extends jspb.Message {
  hasSslParameters(): boolean;
  clearSslParameters(): void;
  getSslParameters(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_ssl_pb.SslParameters | undefined;
  setSslParameters(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_ssl_pb.SslParameters): void;

  getGlobalAnnotationsMap(): jspb.Map<string, string>;
  clearGlobalAnnotationsMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpstreamOptions.AsObject;
  static toObject(includeInstance: boolean, msg: UpstreamOptions): UpstreamOptions.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpstreamOptions, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpstreamOptions;
  static deserializeBinaryFromReader(message: UpstreamOptions, reader: jspb.BinaryReader): UpstreamOptions;
}

export namespace UpstreamOptions {
  export type AsObject = {
    sslParameters?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_ssl_pb.SslParameters.AsObject,
    globalAnnotationsMap: Array<[string, string]>,
  }
}

export class GlooOptions extends jspb.Message {
  getXdsBindAddr(): string;
  setXdsBindAddr(value: string): void;

  getValidationBindAddr(): string;
  setValidationBindAddr(value: string): void;

  hasCircuitBreakers(): boolean;
  clearCircuitBreakers(): void;
  getCircuitBreakers(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_circuit_breaker_pb.CircuitBreakerConfig | undefined;
  setCircuitBreakers(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_circuit_breaker_pb.CircuitBreakerConfig): void;

  hasEndpointsWarmingTimeout(): boolean;
  clearEndpointsWarmingTimeout(): void;
  getEndpointsWarmingTimeout(): google_protobuf_duration_pb.Duration | undefined;
  setEndpointsWarmingTimeout(value?: google_protobuf_duration_pb.Duration): void;

  hasAwsOptions(): boolean;
  clearAwsOptions(): void;
  getAwsOptions(): GlooOptions.AWSOptions | undefined;
  setAwsOptions(value?: GlooOptions.AWSOptions): void;

  hasInvalidConfigPolicy(): boolean;
  clearInvalidConfigPolicy(): void;
  getInvalidConfigPolicy(): GlooOptions.InvalidConfigPolicy | undefined;
  setInvalidConfigPolicy(value?: GlooOptions.InvalidConfigPolicy): void;

  getDisableKubernetesDestinations(): boolean;
  setDisableKubernetesDestinations(value: boolean): void;

  hasDisableGrpcWeb(): boolean;
  clearDisableGrpcWeb(): void;
  getDisableGrpcWeb(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setDisableGrpcWeb(value?: google_protobuf_wrappers_pb.BoolValue): void;

  hasDisableProxyGarbageCollection(): boolean;
  clearDisableProxyGarbageCollection(): void;
  getDisableProxyGarbageCollection(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setDisableProxyGarbageCollection(value?: google_protobuf_wrappers_pb.BoolValue): void;

  hasRegexMaxProgramSize(): boolean;
  clearRegexMaxProgramSize(): void;
  getRegexMaxProgramSize(): google_protobuf_wrappers_pb.UInt32Value | undefined;
  setRegexMaxProgramSize(value?: google_protobuf_wrappers_pb.UInt32Value): void;

  getRestXdsBindAddr(): string;
  setRestXdsBindAddr(value: string): void;

  hasEnableRestEds(): boolean;
  clearEnableRestEds(): void;
  getEnableRestEds(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setEnableRestEds(value?: google_protobuf_wrappers_pb.BoolValue): void;

  hasFailoverUpstreamDnsPollingInterval(): boolean;
  clearFailoverUpstreamDnsPollingInterval(): void;
  getFailoverUpstreamDnsPollingInterval(): google_protobuf_duration_pb.Duration | undefined;
  setFailoverUpstreamDnsPollingInterval(value?: google_protobuf_duration_pb.Duration): void;

  hasRemoveUnusedFilters(): boolean;
  clearRemoveUnusedFilters(): void;
  getRemoveUnusedFilters(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setRemoveUnusedFilters(value?: google_protobuf_wrappers_pb.BoolValue): void;

  getProxyDebugBindAddr(): string;
  setProxyDebugBindAddr(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GlooOptions.AsObject;
  static toObject(includeInstance: boolean, msg: GlooOptions): GlooOptions.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GlooOptions, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GlooOptions;
  static deserializeBinaryFromReader(message: GlooOptions, reader: jspb.BinaryReader): GlooOptions;
}

export namespace GlooOptions {
  export type AsObject = {
    xdsBindAddr: string,
    validationBindAddr: string,
    circuitBreakers?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_circuit_breaker_pb.CircuitBreakerConfig.AsObject,
    endpointsWarmingTimeout?: google_protobuf_duration_pb.Duration.AsObject,
    awsOptions?: GlooOptions.AWSOptions.AsObject,
    invalidConfigPolicy?: GlooOptions.InvalidConfigPolicy.AsObject,
    disableKubernetesDestinations: boolean,
    disableGrpcWeb?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    disableProxyGarbageCollection?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    regexMaxProgramSize?: google_protobuf_wrappers_pb.UInt32Value.AsObject,
    restXdsBindAddr: string,
    enableRestEds?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    failoverUpstreamDnsPollingInterval?: google_protobuf_duration_pb.Duration.AsObject,
    removeUnusedFilters?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    proxyDebugBindAddr: string,
  }

  export class AWSOptions extends jspb.Message {
    hasEnableCredentialsDiscovey(): boolean;
    clearEnableCredentialsDiscovey(): void;
    getEnableCredentialsDiscovey(): boolean;
    setEnableCredentialsDiscovey(value: boolean): void;

    hasServiceAccountCredentials(): boolean;
    clearServiceAccountCredentials(): void;
    getServiceAccountCredentials(): github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_extensions_aws_filter_pb.AWSLambdaConfig.ServiceAccountCredentials | undefined;
    setServiceAccountCredentials(value?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_extensions_aws_filter_pb.AWSLambdaConfig.ServiceAccountCredentials): void;

    hasPropagateOriginalRouting(): boolean;
    clearPropagateOriginalRouting(): void;
    getPropagateOriginalRouting(): google_protobuf_wrappers_pb.BoolValue | undefined;
    setPropagateOriginalRouting(value?: google_protobuf_wrappers_pb.BoolValue): void;

    hasCredentialRefreshDelay(): boolean;
    clearCredentialRefreshDelay(): void;
    getCredentialRefreshDelay(): google_protobuf_duration_pb.Duration | undefined;
    setCredentialRefreshDelay(value?: google_protobuf_duration_pb.Duration): void;

    getCredentialsFetcherCase(): AWSOptions.CredentialsFetcherCase;
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): AWSOptions.AsObject;
    static toObject(includeInstance: boolean, msg: AWSOptions): AWSOptions.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: AWSOptions, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): AWSOptions;
    static deserializeBinaryFromReader(message: AWSOptions, reader: jspb.BinaryReader): AWSOptions;
  }

  export namespace AWSOptions {
    export type AsObject = {
      enableCredentialsDiscovey: boolean,
      serviceAccountCredentials?: github_com_solo_io_solo_apis_api_gloo_gloo_external_envoy_extensions_aws_filter_pb.AWSLambdaConfig.ServiceAccountCredentials.AsObject,
      propagateOriginalRouting?: google_protobuf_wrappers_pb.BoolValue.AsObject,
      credentialRefreshDelay?: google_protobuf_duration_pb.Duration.AsObject,
    }

    export enum CredentialsFetcherCase {
      CREDENTIALS_FETCHER_NOT_SET = 0,
      ENABLE_CREDENTIALS_DISCOVEY = 1,
      SERVICE_ACCOUNT_CREDENTIALS = 2,
    }
  }

  export class InvalidConfigPolicy extends jspb.Message {
    getReplaceInvalidRoutes(): boolean;
    setReplaceInvalidRoutes(value: boolean): void;

    getInvalidRouteResponseCode(): number;
    setInvalidRouteResponseCode(value: number): void;

    getInvalidRouteResponseBody(): string;
    setInvalidRouteResponseBody(value: string): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): InvalidConfigPolicy.AsObject;
    static toObject(includeInstance: boolean, msg: InvalidConfigPolicy): InvalidConfigPolicy.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: InvalidConfigPolicy, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): InvalidConfigPolicy;
    static deserializeBinaryFromReader(message: InvalidConfigPolicy, reader: jspb.BinaryReader): InvalidConfigPolicy;
  }

  export namespace InvalidConfigPolicy {
    export type AsObject = {
      replaceInvalidRoutes: boolean,
      invalidRouteResponseCode: number,
      invalidRouteResponseBody: string,
    }
  }
}

export class VirtualServiceOptions extends jspb.Message {
  hasOneWayTls(): boolean;
  clearOneWayTls(): void;
  getOneWayTls(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setOneWayTls(value?: google_protobuf_wrappers_pb.BoolValue): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VirtualServiceOptions.AsObject;
  static toObject(includeInstance: boolean, msg: VirtualServiceOptions): VirtualServiceOptions.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: VirtualServiceOptions, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VirtualServiceOptions;
  static deserializeBinaryFromReader(message: VirtualServiceOptions, reader: jspb.BinaryReader): VirtualServiceOptions;
}

export namespace VirtualServiceOptions {
  export type AsObject = {
    oneWayTls?: google_protobuf_wrappers_pb.BoolValue.AsObject,
  }
}

export class GatewayOptions extends jspb.Message {
  getValidationServerAddr(): string;
  setValidationServerAddr(value: string): void;

  hasValidation(): boolean;
  clearValidation(): void;
  getValidation(): GatewayOptions.ValidationOptions | undefined;
  setValidation(value?: GatewayOptions.ValidationOptions): void;

  getReadGatewaysFromAllNamespaces(): boolean;
  setReadGatewaysFromAllNamespaces(value: boolean): void;

  getAlwaysSortRouteTableRoutes(): boolean;
  setAlwaysSortRouteTableRoutes(value: boolean): void;

  getCompressedProxySpec(): boolean;
  setCompressedProxySpec(value: boolean): void;

  hasVirtualServiceOptions(): boolean;
  clearVirtualServiceOptions(): void;
  getVirtualServiceOptions(): VirtualServiceOptions | undefined;
  setVirtualServiceOptions(value?: VirtualServiceOptions): void;

  hasPersistProxySpec(): boolean;
  clearPersistProxySpec(): void;
  getPersistProxySpec(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setPersistProxySpec(value?: google_protobuf_wrappers_pb.BoolValue): void;

  hasEnableGatewayController(): boolean;
  clearEnableGatewayController(): void;
  getEnableGatewayController(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setEnableGatewayController(value?: google_protobuf_wrappers_pb.BoolValue): void;

  hasIsolateVirtualHostsBySslConfig(): boolean;
  clearIsolateVirtualHostsBySslConfig(): void;
  getIsolateVirtualHostsBySslConfig(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setIsolateVirtualHostsBySslConfig(value?: google_protobuf_wrappers_pb.BoolValue): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GatewayOptions.AsObject;
  static toObject(includeInstance: boolean, msg: GatewayOptions): GatewayOptions.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GatewayOptions, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GatewayOptions;
  static deserializeBinaryFromReader(message: GatewayOptions, reader: jspb.BinaryReader): GatewayOptions;
}

export namespace GatewayOptions {
  export type AsObject = {
    validationServerAddr: string,
    validation?: GatewayOptions.ValidationOptions.AsObject,
    readGatewaysFromAllNamespaces: boolean,
    alwaysSortRouteTableRoutes: boolean,
    compressedProxySpec: boolean,
    virtualServiceOptions?: VirtualServiceOptions.AsObject,
    persistProxySpec?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    enableGatewayController?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    isolateVirtualHostsBySslConfig?: google_protobuf_wrappers_pb.BoolValue.AsObject,
  }

  export class ValidationOptions extends jspb.Message {
    getProxyValidationServerAddr(): string;
    setProxyValidationServerAddr(value: string): void;

    getValidationWebhookTlsCert(): string;
    setValidationWebhookTlsCert(value: string): void;

    getValidationWebhookTlsKey(): string;
    setValidationWebhookTlsKey(value: string): void;

    getIgnoreGlooValidationFailure(): boolean;
    setIgnoreGlooValidationFailure(value: boolean): void;

    hasAlwaysAccept(): boolean;
    clearAlwaysAccept(): void;
    getAlwaysAccept(): google_protobuf_wrappers_pb.BoolValue | undefined;
    setAlwaysAccept(value?: google_protobuf_wrappers_pb.BoolValue): void;

    hasAllowWarnings(): boolean;
    clearAllowWarnings(): void;
    getAllowWarnings(): google_protobuf_wrappers_pb.BoolValue | undefined;
    setAllowWarnings(value?: google_protobuf_wrappers_pb.BoolValue): void;

    hasWarnRouteShortCircuiting(): boolean;
    clearWarnRouteShortCircuiting(): void;
    getWarnRouteShortCircuiting(): google_protobuf_wrappers_pb.BoolValue | undefined;
    setWarnRouteShortCircuiting(value?: google_protobuf_wrappers_pb.BoolValue): void;

    hasDisableTransformationValidation(): boolean;
    clearDisableTransformationValidation(): void;
    getDisableTransformationValidation(): google_protobuf_wrappers_pb.BoolValue | undefined;
    setDisableTransformationValidation(value?: google_protobuf_wrappers_pb.BoolValue): void;

    hasValidationServerGrpcMaxSizeBytes(): boolean;
    clearValidationServerGrpcMaxSizeBytes(): void;
    getValidationServerGrpcMaxSizeBytes(): google_protobuf_wrappers_pb.Int32Value | undefined;
    setValidationServerGrpcMaxSizeBytes(value?: google_protobuf_wrappers_pb.Int32Value): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): ValidationOptions.AsObject;
    static toObject(includeInstance: boolean, msg: ValidationOptions): ValidationOptions.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: ValidationOptions, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): ValidationOptions;
    static deserializeBinaryFromReader(message: ValidationOptions, reader: jspb.BinaryReader): ValidationOptions;
  }

  export namespace ValidationOptions {
    export type AsObject = {
      proxyValidationServerAddr: string,
      validationWebhookTlsCert: string,
      validationWebhookTlsKey: string,
      ignoreGlooValidationFailure: boolean,
      alwaysAccept?: google_protobuf_wrappers_pb.BoolValue.AsObject,
      allowWarnings?: google_protobuf_wrappers_pb.BoolValue.AsObject,
      warnRouteShortCircuiting?: google_protobuf_wrappers_pb.BoolValue.AsObject,
      disableTransformationValidation?: google_protobuf_wrappers_pb.BoolValue.AsObject,
      validationServerGrpcMaxSizeBytes?: google_protobuf_wrappers_pb.Int32Value.AsObject,
    }
  }
}

export class ConsoleOptions extends jspb.Message {
  hasReadOnly(): boolean;
  clearReadOnly(): void;
  getReadOnly(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setReadOnly(value?: google_protobuf_wrappers_pb.BoolValue): void;

  hasApiExplorerEnabled(): boolean;
  clearApiExplorerEnabled(): void;
  getApiExplorerEnabled(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setApiExplorerEnabled(value?: google_protobuf_wrappers_pb.BoolValue): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ConsoleOptions.AsObject;
  static toObject(includeInstance: boolean, msg: ConsoleOptions): ConsoleOptions.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ConsoleOptions, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ConsoleOptions;
  static deserializeBinaryFromReader(message: ConsoleOptions, reader: jspb.BinaryReader): ConsoleOptions;
}

export namespace ConsoleOptions {
  export type AsObject = {
    readOnly?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    apiExplorerEnabled?: google_protobuf_wrappers_pb.BoolValue.AsObject,
  }
}

export class GraphqlOptions extends jspb.Message {
  hasSchemaChangeValidationOptions(): boolean;
  clearSchemaChangeValidationOptions(): void;
  getSchemaChangeValidationOptions(): GraphqlOptions.SchemaChangeValidationOptions | undefined;
  setSchemaChangeValidationOptions(value?: GraphqlOptions.SchemaChangeValidationOptions): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GraphqlOptions.AsObject;
  static toObject(includeInstance: boolean, msg: GraphqlOptions): GraphqlOptions.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GraphqlOptions, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GraphqlOptions;
  static deserializeBinaryFromReader(message: GraphqlOptions, reader: jspb.BinaryReader): GraphqlOptions;
}

export namespace GraphqlOptions {
  export type AsObject = {
    schemaChangeValidationOptions?: GraphqlOptions.SchemaChangeValidationOptions.AsObject,
  }

  export class SchemaChangeValidationOptions extends jspb.Message {
    hasRejectBreakingChanges(): boolean;
    clearRejectBreakingChanges(): void;
    getRejectBreakingChanges(): google_protobuf_wrappers_pb.BoolValue | undefined;
    setRejectBreakingChanges(value?: google_protobuf_wrappers_pb.BoolValue): void;

    clearProcessingRulesList(): void;
    getProcessingRulesList(): Array<GraphqlOptions.SchemaChangeValidationOptions.ProcessingRuleMap[keyof GraphqlOptions.SchemaChangeValidationOptions.ProcessingRuleMap]>;
    setProcessingRulesList(value: Array<GraphqlOptions.SchemaChangeValidationOptions.ProcessingRuleMap[keyof GraphqlOptions.SchemaChangeValidationOptions.ProcessingRuleMap]>): void;
    addProcessingRules(value: GraphqlOptions.SchemaChangeValidationOptions.ProcessingRuleMap[keyof GraphqlOptions.SchemaChangeValidationOptions.ProcessingRuleMap], index?: number): GraphqlOptions.SchemaChangeValidationOptions.ProcessingRuleMap[keyof GraphqlOptions.SchemaChangeValidationOptions.ProcessingRuleMap];

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): SchemaChangeValidationOptions.AsObject;
    static toObject(includeInstance: boolean, msg: SchemaChangeValidationOptions): SchemaChangeValidationOptions.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: SchemaChangeValidationOptions, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): SchemaChangeValidationOptions;
    static deserializeBinaryFromReader(message: SchemaChangeValidationOptions, reader: jspb.BinaryReader): SchemaChangeValidationOptions;
  }

  export namespace SchemaChangeValidationOptions {
    export type AsObject = {
      rejectBreakingChanges?: google_protobuf_wrappers_pb.BoolValue.AsObject,
      processingRulesList: Array<GraphqlOptions.SchemaChangeValidationOptions.ProcessingRuleMap[keyof GraphqlOptions.SchemaChangeValidationOptions.ProcessingRuleMap]>,
    }

    export interface ProcessingRuleMap {
      RULE_UNSPECIFIED: 0;
      RULE_DANGEROUS_TO_BREAKING: 1;
      RULE_DEPRECATED_FIELD_REMOVAL_DANGEROUS: 2;
      RULE_IGNORE_DESCRIPTION_CHANGES: 3;
      RULE_IGNORE_UNREACHABLE: 4;
    }

    export const ProcessingRule: ProcessingRuleMap;
  }
}

export class SettingsStatus extends jspb.Message {
  getState(): SettingsStatus.StateMap[keyof SettingsStatus.StateMap];
  setState(value: SettingsStatus.StateMap[keyof SettingsStatus.StateMap]): void;

  getReason(): string;
  setReason(value: string): void;

  getReportedBy(): string;
  setReportedBy(value: string): void;

  getSubresourceStatusesMap(): jspb.Map<string, SettingsStatus>;
  clearSubresourceStatusesMap(): void;
  hasDetails(): boolean;
  clearDetails(): void;
  getDetails(): google_protobuf_struct_pb.Struct | undefined;
  setDetails(value?: google_protobuf_struct_pb.Struct): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SettingsStatus.AsObject;
  static toObject(includeInstance: boolean, msg: SettingsStatus): SettingsStatus.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: SettingsStatus, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SettingsStatus;
  static deserializeBinaryFromReader(message: SettingsStatus, reader: jspb.BinaryReader): SettingsStatus;
}

export namespace SettingsStatus {
  export type AsObject = {
    state: SettingsStatus.StateMap[keyof SettingsStatus.StateMap],
    reason: string,
    reportedBy: string,
    subresourceStatusesMap: Array<[string, SettingsStatus.AsObject]>,
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

export class SettingsNamespacedStatuses extends jspb.Message {
  getStatusesMap(): jspb.Map<string, SettingsStatus>;
  clearStatusesMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SettingsNamespacedStatuses.AsObject;
  static toObject(includeInstance: boolean, msg: SettingsNamespacedStatuses): SettingsNamespacedStatuses.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: SettingsNamespacedStatuses, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SettingsNamespacedStatuses;
  static deserializeBinaryFromReader(message: SettingsNamespacedStatuses, reader: jspb.BinaryReader): SettingsNamespacedStatuses;
}

export namespace SettingsNamespacedStatuses {
  export type AsObject = {
    statusesMap: Array<[string, SettingsStatus.AsObject]>,
  }
}
