import * as React from 'react';
import { SoloTable } from 'Components/Common/SoloTable';
import { DetailsSectionTitle } from './VirtualServiceDetails';
import { Route } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/proxy_pb';
import {
  getRouteMethods,
  getRouteSingleUpstream,
  getRouteMatcher,
  getRouteHeaders,
  getRouteQueryParams
} from 'utils/helpers';

const routeColumns = [
  {
    title: 'Matcher',
    dataIndex: 'matcher',
    width: 200
  },
  {
    title: 'Path Match Type',
    dataIndex: 'pathMatch'
  },
  {
    title: 'Methods',
    dataIndex: 'method'
  },
  {
    title: 'Destination',
    dataIndex: 'upstreamName'
  },
  {
    title: 'Headers',
    dataIndex: 'header'
  },
  {
    title: 'Query Parameters',
    dataIndex: 'queryParams'
  },
  {
    title: 'Actions',
    dataIndex: 'actions',
    render: (text: any) => <div>ACTION!</div>
  }
];

interface Props {
  routes: Route.AsObject[];
}

export const Routes: React.FC<Props> = props => {
  const getRouteData = (routes: Route.AsObject[]) => {
    return routes.map(route => {
      const upstreamName = getRouteSingleUpstream(route).name || '';
      const { matcher, matchType } = getRouteMatcher(route);
      console.log(route);
      return {
        key: route.matcher!.prefix,
        matcher: matcher,
        pathMatch: matchType,
        method: getRouteMethods(route),
        upstreamName: upstreamName,
        header: getRouteHeaders(route),
        queryParams: getRouteQueryParams(route),
        actions: ''
      };
    });
  };

  return (
    <React.Fragment>
      <DetailsSectionTitle>Routes</DetailsSectionTitle>
      <SoloTable
        columns={routeColumns}
        dataSource={getRouteData(props.routes)}
      />
    </React.Fragment>
  );
};
