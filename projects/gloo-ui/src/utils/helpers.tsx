import { Upstream } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/upstream_pb';
import { VirtualService } from 'proto/github.com/solo-io/gloo/projects/gateway/api/v1/virtual_service_pb';
import { Route } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/proxy_pb';
import { ResourceRef } from 'proto/github.com/solo-io/solo-kit/api/v1/ref_pb';

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

export function groupBy<T>(data: T[], getKey: (item: T) => string) {
  const map = new Map<string, T[]>();
  data.forEach(resource => {
    const key = getKey(resource);
    if (!map.get(key)) {
      map.set(key, [resource]);
    } else {
      map.get(key)!.push(resource);
    }
  });
  return map;
}

/* -------------------------------------------------------------------------- */
/*                                  UPSTREAMS                                 */
/* -------------------------------------------------------------------------- */

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

export const UPSTREAM_TYPES = [
  {
    key: 'AWS',
    value: 'AWS'
  },
  {
    key: 'Azure',
    value: 'Azure'
  },
  {
    key: 'Kubernetes',
    value: 'Kubernetes'
  },
  {
    key: 'Static',
    value: 'Static'
  }
];

export enum UPSTREAM_SPEC_TYPES {
  AZURE = 'Azure',
  KUBE = 'Kubernetes',
  AWS = 'AWS',
  STATIC = 'Static'
}

/* -------------------------------------------------------------------------- */
/*                              VIRTUAL SERVICES                              */
/* -------------------------------------------------------------------------- */

export function getVSDomains(virtualService: VirtualService.AsObject) {
  if (virtualService.virtualHost && virtualService.virtualHost.domainsList) {
    return virtualService.virtualHost.domainsList.join(', ');
  }
}

/* -------------------------------------------------------------------------- */
/*                                   ROUTES                                   */
/* -------------------------------------------------------------------------- */

export function getRouteMethods(route: Route.AsObject) {
  if (route.matcher && route.matcher.methodsList) {
    return route.matcher.methodsList.join(', ');
  }
  return '';
}

export function getRouteSingleUpstream(route: Route.AsObject) {
  if (route.routeAction) {
    if (route.routeAction.single && route.routeAction.single.upstream) {
      return route.routeAction.single.upstream;
    }
  }
  return {} as ResourceRef.AsObject;
}

// TODO: handle multi destination case
export function getRouteMultiUpstream(route: Route.AsObject) {}

export function getRouteMatcher(route: Route.AsObject) {
  let matcher = '';
  let matchType = 'PATH_SPECIFIER_NOT_SET';
  if (route.matcher) {
    if (route.matcher.prefix) {
      matcher = route.matcher.prefix;
      matchType = 'PREFIX';
    }
    if (route.matcher.exact) {
      matcher = route.matcher.exact;
      matchType = 'EXACT';
    }
    if (route.matcher.regex) {
      matcher = route.matcher.regex;
      matchType = 'REGEX';
    }
  }
  return {
    matcher,
    matchType
  };
}

export function getRouteHeaders(route: Route.AsObject) {
  if (route.matcher && route.matcher.headersList) {
    return route.matcher.headersList.map(header => header.name).join(', ');
  }
  return '';
}

export function getRouteQueryParams(route: Route.AsObject) {
  if (route.matcher && route.matcher.queryParametersList) {
    return route.matcher.queryParametersList
      .map(param => param.name)
      .join(', ');
  }
  return '';
}
