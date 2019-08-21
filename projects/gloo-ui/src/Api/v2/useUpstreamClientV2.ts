import { upstreams, UpstreamSpecificValues } from './UpstreamClient';
import { UPSTREAM_SPEC_TYPES } from 'utils/upstreamHelpers';
import { useSendRequest } from './requestReducerV2';

// export function useGetUpstreamsListV2(listUpstreamRequest: { namespacesList: string[] }) {
//   return useSendRequest(listUpstreamRequest, upstreams.getUpstreamsList);
// }

// export function useGetUpstreamV2(getUpstreamRequest: {
//   name: string;
//   namespace: string;
// }) {
//   return useSendRequest(getUpstreamRequest, upstreams.getUpstream);
// }

// export function useCreateUpstreamV2(variables: {
//   name: string;
//   namespace: string;
//   type: UPSTREAM_SPEC_TYPES;
//   values: UpstreamSpecificValues;
// }) {
//   return useSendRequest(variables, upstreams.createUpstream);
// }

// export function useDeleteUpstreamV2(variables: {
//   name: string;
//   namespace: string;
// }) {
//   return useSendRequest(variables, upstreams.deleteUpstream);
// }
