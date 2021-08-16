import { Route } from 'proto/github.com/solo-io/solo-apis/api/gloo/gateway/v1/virtual_service_pb';

export function getRouteMethods(route: Route.AsObject) {
  if (route.matchersList[0] && route.matchersList[0].methodsList) {
    return route.matchersList[0].methodsList.join(', ');
  }
  return '*';
}

export function getRouteSingleUpstream(route: Route.AsObject) {
  if (route.routeAction) {
    if (route.routeAction.single && route.routeAction.single.upstream) {
      let functionName = '';
      if (!!route.routeAction.single.destinationSpec) {
        const spec = route.routeAction.single.destinationSpec;
        functionName = spec.aws?.logicalName ?? spec.azure?.functionName ?? spec.grpc?.pb_function ?? spec.rest?.functionName ?? '';
      }
      return `${route.routeAction.single.upstream.name}${functionName && ':' + functionName}`;
    }
  } else if (route.delegateAction) {
    return route.delegateAction.name;
  }
  return '';
}

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
  return { matcher, matchType };
}

export function getRouteHeaders(route: Route.AsObject) {
  if (route.matchersList[0] && route.matchersList[0].headersList) {
    return route.matchersList[0].headersList;
  }
  return [];
}

export function getRouteQueryParams(route: Route.AsObject) {
  if (route.matchersList[0] && route.matchersList[0].queryParametersList) {
    return route.matchersList[0].queryParametersList;
  }
  return [];
}
