// package: gloo.solo.io
// file: github.com/solo-io/gloo/projects/gloo/api/v1/settings.proto

import * as jspb from "google-protobuf";
import * as gogoproto_gogo_pb from "../../../../../../../gogoproto/gogo_pb";
import * as github_com_solo_io_solo_kit_api_v1_metadata_pb from "../../../../../../../github.com/solo-io/solo-kit/api/v1/metadata_pb";
import * as github_com_solo_io_solo_kit_api_v1_status_pb from "../../../../../../../github.com/solo-io/solo-kit/api/v1/status_pb";
import * as github_com_solo_io_gloo_projects_gloo_api_v1_extensions_pb from "../../../../../../../github.com/solo-io/gloo/projects/gloo/api/v1/extensions_pb";
import * as github_com_solo_io_gloo_projects_gloo_api_v1_circuit_breaker_pb from "../../../../../../../github.com/solo-io/gloo/projects/gloo/api/v1/circuit_breaker_pb";
import * as google_protobuf_duration_pb from "google-protobuf/google/protobuf/duration_pb";

export class Settings extends jspb.Message {
  getDiscoveryNamespace(): string;
  setDiscoveryNamespace(value: string): void;

  clearWatchNamespacesList(): void;
  getWatchNamespacesList(): Array<string>;
  setWatchNamespacesList(value: Array<string>): void;
  addWatchNamespaces(value: string, index?: number): string;

  hasKubernetesConfigSource(): boolean;
  clearKubernetesConfigSource(): void;
  getKubernetesConfigSource(): Settings.KubernetesCrds | undefined;
  setKubernetesConfigSource(value?: Settings.KubernetesCrds): void;

  hasDirectoryConfigSource(): boolean;
  clearDirectoryConfigSource(): void;
  getDirectoryConfigSource(): Settings.Directory | undefined;
  setDirectoryConfigSource(value?: Settings.Directory): void;

  hasKubernetesSecretSource(): boolean;
  clearKubernetesSecretSource(): void;
  getKubernetesSecretSource(): Settings.KubernetesSecrets | undefined;
  setKubernetesSecretSource(value?: Settings.KubernetesSecrets): void;

  hasVaultSecretSource(): boolean;
  clearVaultSecretSource(): void;
  getVaultSecretSource(): Settings.VaultSecrets | undefined;
  setVaultSecretSource(value?: Settings.VaultSecrets): void;

  hasDirectorySecretSource(): boolean;
  clearDirectorySecretSource(): void;
  getDirectorySecretSource(): Settings.Directory | undefined;
  setDirectorySecretSource(value?: Settings.Directory): void;

  hasKubernetesArtifactSource(): boolean;
  clearKubernetesArtifactSource(): void;
  getKubernetesArtifactSource(): Settings.KubernetesConfigmaps | undefined;
  setKubernetesArtifactSource(value?: Settings.KubernetesConfigmaps): void;

  hasDirectoryArtifactSource(): boolean;
  clearDirectoryArtifactSource(): void;
  getDirectoryArtifactSource(): Settings.Directory | undefined;
  setDirectoryArtifactSource(value?: Settings.Directory): void;

  getBindAddr(): string;
  setBindAddr(value: string): void;

  hasRefreshRate(): boolean;
  clearRefreshRate(): void;
  getRefreshRate(): google_protobuf_duration_pb.Duration | undefined;
  setRefreshRate(value?: google_protobuf_duration_pb.Duration): void;

  getDevMode(): boolean;
  setDevMode(value: boolean): void;

  getLinkerd(): boolean;
  setLinkerd(value: boolean): void;

  hasCircuitBreakers(): boolean;
  clearCircuitBreakers(): void;
  getCircuitBreakers(): github_com_solo_io_gloo_projects_gloo_api_v1_circuit_breaker_pb.CircuitBreakerConfig | undefined;
  setCircuitBreakers(value?: github_com_solo_io_gloo_projects_gloo_api_v1_circuit_breaker_pb.CircuitBreakerConfig): void;

  hasExtensions(): boolean;
  clearExtensions(): void;
  getExtensions(): github_com_solo_io_gloo_projects_gloo_api_v1_extensions_pb.Extensions | undefined;
  setExtensions(value?: github_com_solo_io_gloo_projects_gloo_api_v1_extensions_pb.Extensions): void;

  hasMetadata(): boolean;
  clearMetadata(): void;
  getMetadata(): github_com_solo_io_solo_kit_api_v1_metadata_pb.Metadata | undefined;
  setMetadata(value?: github_com_solo_io_solo_kit_api_v1_metadata_pb.Metadata): void;

  hasStatus(): boolean;
  clearStatus(): void;
  getStatus(): github_com_solo_io_solo_kit_api_v1_status_pb.Status | undefined;
  setStatus(value?: github_com_solo_io_solo_kit_api_v1_status_pb.Status): void;

  getConfigSourceCase(): Settings.ConfigSourceCase;
  getSecretSourceCase(): Settings.SecretSourceCase;
  getArtifactSourceCase(): Settings.ArtifactSourceCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Settings.AsObject;
  static toObject(includeInstance: boolean, msg: Settings): Settings.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Settings, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Settings;
  static deserializeBinaryFromReader(message: Settings, reader: jspb.BinaryReader): Settings;
}

export namespace Settings {
  export type AsObject = {
    discoveryNamespace: string,
    watchNamespacesList: Array<string>,
    kubernetesConfigSource?: Settings.KubernetesCrds.AsObject,
    directoryConfigSource?: Settings.Directory.AsObject,
    kubernetesSecretSource?: Settings.KubernetesSecrets.AsObject,
    vaultSecretSource?: Settings.VaultSecrets.AsObject,
    directorySecretSource?: Settings.Directory.AsObject,
    kubernetesArtifactSource?: Settings.KubernetesConfigmaps.AsObject,
    directoryArtifactSource?: Settings.Directory.AsObject,
    bindAddr: string,
    refreshRate?: google_protobuf_duration_pb.Duration.AsObject,
    devMode: boolean,
    linkerd: boolean,
    circuitBreakers?: github_com_solo_io_gloo_projects_gloo_api_v1_circuit_breaker_pb.CircuitBreakerConfig.AsObject,
    extensions?: github_com_solo_io_gloo_projects_gloo_api_v1_extensions_pb.Extensions.AsObject,
    metadata?: github_com_solo_io_solo_kit_api_v1_metadata_pb.Metadata.AsObject,
    status?: github_com_solo_io_solo_kit_api_v1_status_pb.Status.AsObject,
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

  export enum ConfigSourceCase {
    CONFIG_SOURCE_NOT_SET = 0,
    KUBERNETES_CONFIG_SOURCE = 4,
    DIRECTORY_CONFIG_SOURCE = 5,
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
  }
}

