import * as React from 'react';
import styled from '@emotion/styled/macro';
import { colors } from 'Styles';
import { DetailsSectionTitle } from './VirtualServiceDetails';
import { ExtAuth } from './ExtAuth';
import { RateLimit } from './RateLimit';

const ConfigContainer = styled.div`
  display: grid;
  grid-template-columns: 1fr 1fr;
  background: ${colors.januaryGrey};
  height: 80%;
`;

const ConfigItem = styled.div`
  margin: 20px;
  padding: 10px;
  justify-items: center;
  background: white;
`;

export const Configuration = () => {
  return (
    <React.Fragment>
      <DetailsSectionTitle>Configuration</DetailsSectionTitle>
      <ConfigContainer>
        <ConfigItem>
          <ExtAuth />
        </ConfigItem>
        <ConfigItem>
          <RateLimit />
        </ConfigItem>
      </ConfigContainer>
    </React.Fragment>
  );
};
