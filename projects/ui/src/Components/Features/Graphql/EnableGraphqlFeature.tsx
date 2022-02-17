import * as React from 'react';
import { Navigate } from 'react-router-dom';
import { useIsGraphqlEnabled } from 'API/hooks';

type Props = {
  children: React.ReactElement;
  reroute?: boolean;
};

export const EnableGraphqlFeature: React.FC<Props> = props => {
  const { data: graphqlIntegrationEnabled, error: graphqlCheckError } =
    useIsGraphqlEnabled();
  if (graphqlCheckError || !graphqlIntegrationEnabled) {
    return props.reroute ? <Navigate to='/' replace /> : null;
  }
  return props.children;
};
