import { Metadata } from 'proto/github.com/solo-io/solo-kit/api/v1/metadata_pb';

export function createVirtualServiceId(
  virtualServiceMetadata: Metadata.AsObject
): string {
  return `${virtualServiceMetadata!.name}-.-${
    virtualServiceMetadata!.namespace
  }`;
}

export function parseVirtualServiceId(
  virtualServiceId: string
): {
  name: string;
  namespace: string;
} {
  const idData = virtualServiceId.split('-.-');

  return {
    name: idData[0],
    namespace: idData[1]
  };
}
