import styled from '@emotion/styled';
import { ReactComponent as AWSLogo } from 'assets/aws-logo.svg';
import { ReactComponent as AzureLogo } from 'assets/azure-logo.svg';
import Gloo from 'assets/Gloo.svg';
// TODO: get svg format GRPC logo
// import { ReactComponent as GRPCLogo } from 'assets/grpc-logo.svg';
import GRPCLogo from 'assets/grpc-logo.png';
import { ReactComponent as KubeLogo } from 'assets/kube-logo.svg';
import { ReactComponent as RESTLogo } from 'assets/rest-logo.svg';
import { ReactComponent as StaticLogo } from 'assets/static-logo.svg';
import {
  VirtualService,
  Route
} from 'proto/github.com/solo-io/gloo/projects/gateway/api/v1/virtual_service_pb';
import { Upstream } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/upstream_pb';
import { ResourceRef } from 'proto/github.com/solo-io/solo-kit/api/v1/ref_pb';
import * as React from 'react';
import RT from 'assets/route-table-icon.png';
import { RouteTable } from 'proto/github.com/solo-io/gloo/projects/gateway/api/v1/route_table_pb';
type Resource =
  | VirtualService.AsObject
  | Upstream.AsObject
  | RouteTable.AsObject
  | number;

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
    case 'Route Table':
      return <img src={RT} style={{ width: '25px', paddingRight: '5px' }} />;
    default:
      return <img src={Gloo} style={{ width: '25px', paddingRight: '5px' }} />;
  }
}

export function getIconFromSpec(spec: Upstream.AsObject) {
  if (spec?.kube !== undefined)
    return (
      <KubeLogo
        style={{
          width: '25px',
          paddingRight: '5px'
        }}
      />
    );
  if (spec?.aws !== undefined)
    return (
      <AWSLogo
        style={{
          width: '25px',
          paddingRight: '5px'
        }}
      />
    );
  if (spec?.azure !== undefined)
    return (
      <AzureLogo
        style={{
          width: '25px',
          paddingRight: '5px'
        }}
      />
    );
  if (spec?.pb_static !== undefined)
    return (
      <StaticLogo
        style={{
          width: '25px',
          paddingRight: '5px'
        }}
      />
    );
  else {
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
  if (route.matchersList[0] && route.matchersList[0].methodsList) {
    if (route.matchersList[0].methodsList.length === 7) {
      // all options selected
      return '*';
    }

    return route.matchersList[0].methodsList.join(', ');
  }
  return '*';
}

export function getRouteSingleUpstream(route: Route.AsObject) {
  let functionName = '';
  if (route.routeAction) {
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
  } else if (route.delegateAction) {
    functionName = route.delegateAction.name;
    return functionName;
  }
  return functionName;
}

// TODO: handle multi destination case
export function getRouteMultiUpstream(route: Route.AsObject) {}

export function getRouteMatcher(route: Route.AsObject) {
  let matcher = '';
  let matchType = 'PATH_SPECIFIER_NOT_SET';
  if (route.matchersList[0]) {
    if (route.matchersList[0].prefix) {
      matcher = route.matchersList[0].prefix;
      matchType = 'PREFIX';
    }
    if (route.matchersList[0].exact) {
      matcher = route.matchersList[0].exact;
      matchType = 'EXACT';
    }
    if (route.matchersList[0].regex) {
      matcher = route.matchersList[0].regex;
      matchType = 'REGEX';
    }
  }
  return {
    matcher,
    matchType
  };
}

export function getRouteHeaders(route: Route.AsObject) {
  if (route.matchersList[0] && route.matchersList[0].headersList) {
    return route.matchersList[0].headersList.map(header => (
      <div>
        {header.name}: {header.value}
      </div>
    ));
  }
  return '';
}

export function getRouteQueryParams(route: Route.AsObject) {
  if (route.matchersList[0] && route.matchersList[0].queryParametersList) {
    return route.matchersList[0].queryParametersList.map(param => (
      <div>
        {param.name}: {param.value}
      </div>
    ));
  }
  return '';
}

export * from './upstreamHelpers';
export * from './virtualServiceHelpers';
