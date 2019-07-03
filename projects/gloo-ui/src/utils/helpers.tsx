import { Upstream } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/upstream_pb';
import { VirtualService } from 'proto/github.com/solo-io/gloo/projects/gateway/api/v1/virtual_service_pb';

type Resource = VirtualService.AsObject | Upstream.AsObject;

export function getResourceStatus(resource: Resource) {
  switch (resource.status!.state) {
    case 0:
      return 'PENDING';
    case 1:
      return 'ACCEPTED';
    case 2:
      return 'REJECTED';
    default:
      return '';
  }
}

export function getUpstreamType(upstream: Upstream.AsObject) {
  let upstreamType = '';
  if (!!upstream.upstreamSpec!.aws) {
    upstreamType = 'AWS';
  }
  if (!!upstream.upstreamSpec!.azure) {
    upstreamType = 'Azure';
  }

  if (!!upstream.upstreamSpec!.consul) {
    upstreamType = 'Consul';
  }

  if (!!upstream.upstreamSpec!.kube) {
    upstreamType = 'Kubernetes';
  }
  return upstreamType;
}

export function getVSDomains(virtualService: VirtualService.AsObject) {
  if (virtualService.virtualHost && virtualService.virtualHost.domainsList) {
    return virtualService.virtualHost.domainsList.join(', ');
  }
}
