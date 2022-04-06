import { useGetConsoleOptions } from 'API/hooks';
import React from 'react';
import * as styles from './GraphqlLanding.style';
import { GraphqlLandingTableContainer } from './apis-table/GraphqlLandingTableContainer';
import { NewApiButton } from './new-api-modal/NewApiButton';

export const GraphqlLanding = () => {
  const { readonly } = useGetConsoleOptions();

  return (
    <styles.GraphqlLandingContainer>
      {!readonly && <NewApiButton />}
      <GraphqlLandingTableContainer />
    </styles.GraphqlLandingContainer>
  );
};
