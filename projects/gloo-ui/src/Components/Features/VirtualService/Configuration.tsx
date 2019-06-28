import * as React from 'react';
import styled from '@emotion/styled/macro';
import { colors } from 'Styles';

const ConfigContainer = styled.div`
  display: grid;
  grid-template-columns: 1fr 1fr;
  background: ${colors.januaryGrey};
  height: 80%;
`;

const ConfigItem = styled.div`
  display: flex;
  margin: 20px;
  align-items: flex-start;
  justify-items: center;
  background: white;
`;

const ConfigItemHeader = styled.div`
  font-weight: 600;
  font-size: 14px;
  color: ${colors.novemberGrey};
`;
export const Configuration = () => {
  return (
    <ConfigContainer>
      <ConfigItem>
        <ConfigItemHeader>External Auth </ConfigItemHeader>
      </ConfigItem>
      <ConfigItem>
        <ConfigItemHeader>Rate Limits</ConfigItemHeader>
      </ConfigItem>
    </ConfigContainer>
  );
};
