import * as React from 'react';
import { Navigate, useLocation } from 'react-router-dom';

type Props = {
  children: React.ReactElement;
  reroute?: boolean;
};

export const EnableGraphqlFeature: React.FC<Props> = props => {
  let query = new URLSearchParams(useLocation().search);
  // @ts-ignore
  const graphqlIntegrationEnabled =
    query.get('graphql') === 'enabled' ||
    query.get('graphql') === 'on' ||
    process.env.REACT_APP_GRAPHQL_INTEGRATION === 'true';
  if (!graphqlIntegrationEnabled) {
    return props.reroute ? <Navigate to='/' replace /> : null;
  }
  return props.children;
};
