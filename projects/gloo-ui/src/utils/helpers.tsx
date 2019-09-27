import styled from '@emotion/styled';
import { ReactComponent as AWSLogo } from 'assets/aws-logo.svg';
import { ReactComponent as AzureLogo } from 'assets/azure-logo.svg';
import { ReactComponent as Gloo } from 'assets/Gloo.svg';
// TODO: get svg format GRPC logo
// import { ReactComponent as GRPCLogo } from 'assets/grpc-logo.svg';
import GRPCLogo from 'assets/grpc-logo.png';
import { ReactComponent as KubeLogo } from 'assets/kube-logo.svg';
import { ReactComponent as RESTLogo } from 'assets/rest-logo.svg';
import { ReactComponent as StaticLogo } from 'assets/static-logo.svg';
import { VirtualService } from 'proto/github.com/solo-io/gloo/projects/gateway/api/v1/virtual_service_pb';
import { Route } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/proxy_pb';
import { Upstream } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/upstream_pb';
import { ResourceRef } from 'proto/github.com/solo-io/solo-kit/api/v1/ref_pb';
import * as React from 'react';

type Resource = VirtualService.AsObject | Upstream.AsObject | number;

const StyledGRPCLogo = styled.img`
  width: 20px;
  max-height: 25px;
`;
/* -------------------------------------------------------------------------- */
/*                                   GENERAL                                  */
/* -------------------------------------------------------------------------- */

export function getResourceStatus(resource: Resource) {
  const status =
    typeof resource === 'number'
      ? resource
      : !resource.status
      ? 'Pending'
      : resource.status!.state;
  switch (status) {
    case 0:
      return 'Pending';
    case 1:
      return 'Accepted';
    case 2:
      return 'Rejected';
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

export function getIcon(type: string) {
  switch (type) {
    case 'Kubernetes':
      return <KubeLogo />;
    case 'Aws':
      return <AWSLogo />;
    case 'Azure':
      return <AzureLogo />;
    case 'GRPC':
      return <StyledGRPCLogo src={GRPCLogo} />;
    case 'Static':
      return <StaticLogo />;
    case 'REST':
      return <RESTLogo />;
    default:
      return <Gloo />;
  }
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
    if (route.matcher.methodsList.length === 7) {
      // all options selected
      return '*';
    }

    return route.matcher.methodsList.join(', ');
  }
  return '*';
}

export function getRouteSingleUpstream(route: Route.AsObject) {
  if (route.routeAction) {
    let functionName = '';
    if (route.routeAction.single && route.routeAction.single.upstream) {
      if (!!route.routeAction.single.destinationSpec) {
        if (route.routeAction.single.destinationSpec.aws) {
          functionName =
            route.routeAction.single.destinationSpec.aws.logicalName;
        }
        if (route.routeAction.single.destinationSpec.azure) {
          functionName =
            route.routeAction.single.destinationSpec.azure.functionName;
        }
        if (route.routeAction.single.destinationSpec.grpc) {
          functionName =
            route.routeAction.single.destinationSpec.grpc.pb_function;
        }
        if (route.routeAction.single.destinationSpec.rest) {
          functionName =
            route.routeAction.single.destinationSpec.rest.functionName;
        }
      }
      return `${route.routeAction.single.upstream.name}${
        !functionName ? '' : ':' + functionName
      }`;
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
    return route.matcher.headersList.map(header => (
      <div>
        {header.name}: {header.value}
      </div>
    ));
  }
  return '';
}

export function getRouteQueryParams(route: Route.AsObject) {
  if (route.matcher && route.matcher.queryParametersList) {
    return route.matcher.queryParametersList.map(param => (
      <div>
        {param.name}: {param.value}
      </div>
    ));
  }
  return '';
}

export * from './upstreamHelpers';
export * from './virtualServiceHelpers';
