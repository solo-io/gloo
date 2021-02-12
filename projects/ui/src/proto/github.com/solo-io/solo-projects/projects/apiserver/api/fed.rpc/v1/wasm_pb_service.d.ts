// package: fed.rpc.solo.io
// file: github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/wasm.proto

import * as github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_wasm_pb from "../../../../../../../../github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/wasm_pb";
import {grpc} from "@improbable-eng/grpc-web";

type WasmFilterApiListWasmFilters = {
  readonly methodName: string;
  readonly service: typeof WasmFilterApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_wasm_pb.ListWasmFiltersRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_wasm_pb.ListWasmFiltersResponse;
};

type WasmFilterApiDescribeWasmFilter = {
  readonly methodName: string;
  readonly service: typeof WasmFilterApi;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_wasm_pb.DescribeWasmFilterRequest;
  readonly responseType: typeof github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_wasm_pb.DescribeWasmFilterResponse;
};

export class WasmFilterApi {
  static readonly serviceName: string;
  static readonly ListWasmFilters: WasmFilterApiListWasmFilters;
  static readonly DescribeWasmFilter: WasmFilterApiDescribeWasmFilter;
}

export type ServiceError = { message: string, code: number; metadata: grpc.Metadata }
export type Status = { details: string, code: number; metadata: grpc.Metadata }

interface UnaryResponse {
  cancel(): void;
}
interface ResponseStream<T> {
  cancel(): void;
  on(type: 'data', handler: (message: T) => void): ResponseStream<T>;
  on(type: 'end', handler: (status?: Status) => void): ResponseStream<T>;
  on(type: 'status', handler: (status: Status) => void): ResponseStream<T>;
}
interface RequestStream<T> {
  write(message: T): RequestStream<T>;
  end(): void;
  cancel(): void;
  on(type: 'end', handler: (status?: Status) => void): RequestStream<T>;
  on(type: 'status', handler: (status: Status) => void): RequestStream<T>;
}
interface BidirectionalStream<ReqT, ResT> {
  write(message: ReqT): BidirectionalStream<ReqT, ResT>;
  end(): void;
  cancel(): void;
  on(type: 'data', handler: (message: ResT) => void): BidirectionalStream<ReqT, ResT>;
  on(type: 'end', handler: (status?: Status) => void): BidirectionalStream<ReqT, ResT>;
  on(type: 'status', handler: (status: Status) => void): BidirectionalStream<ReqT, ResT>;
}

export class WasmFilterApiClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  listWasmFilters(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_wasm_pb.ListWasmFiltersRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_wasm_pb.ListWasmFiltersResponse|null) => void
  ): UnaryResponse;
  listWasmFilters(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_wasm_pb.ListWasmFiltersRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_wasm_pb.ListWasmFiltersResponse|null) => void
  ): UnaryResponse;
  describeWasmFilter(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_wasm_pb.DescribeWasmFilterRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_wasm_pb.DescribeWasmFilterResponse|null) => void
  ): UnaryResponse;
  describeWasmFilter(
    requestMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_wasm_pb.DescribeWasmFilterRequest,
    callback: (error: ServiceError|null, responseMessage: github_com_solo_io_solo_projects_projects_apiserver_api_fed_rpc_v1_wasm_pb.DescribeWasmFilterResponse|null) => void
  ): UnaryResponse;
}

