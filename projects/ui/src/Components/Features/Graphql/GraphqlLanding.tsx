import React from 'react';
import { GraphqlLandingTableContainer } from './apis-table/GraphqlLandingTableContainer';
import * as styles from './GraphqlLanding.style';

export const GraphqlLanding = () => {
  return (
    <styles.GraphqlLandingContainer>
      <GraphqlLandingTableContainer />
    </styles.GraphqlLandingContainer>
  );
};
