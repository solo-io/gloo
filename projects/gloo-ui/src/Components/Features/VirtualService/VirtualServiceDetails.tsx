import * as React from 'react';
import { SectionCard } from 'Components/Common/SectionCard';
import { ReactComponent as GlooIcon } from 'assets/GlooEE.svg';
import { Domains } from './Domains';
import { Routes } from './Routes';
import { Configuration } from './Configuration';
import styled from '@emotion/styled/macro';
import { colors, soloConstants } from 'Styles';

const CardAddonsContainer = styled.div`
  /* padding: ${soloConstants.buffer}px ${soloConstants.smallBuffer}px 10px;
  display: flex;
  flex-direction: column; */
  display: grid;
  grid-template-rows: auto 2fr 1fr;
  grid-column-gap: 30px;
`;

const CardAddons = styled.div`
  width: 100%;
`;
const CardAddonsTitle = styled.div`
  font-size: 18px;
  font-weight: bold;
  color: ${colors.novemberGrey};
  margin-bottom: 10px;
`;

const MoreTease = styled.div`
  margin: 10px 0 0;
  font-size: 14px;
  color: ${colors.septemberGrey};
`;

const headerInfo = [
  {
    title: 'namespace',
    value: 'dio'
  },
  {
    title: 'namespace',
    value: 'default'
  }
];

export const VirtualServiceDetails = () => {
  return (
    <React.Fragment>
      <SectionCard
        cardName='Dio-test'
        logoIcon={<GlooIcon />}
        health={1}
        headerSecondaryInformation={headerInfo}
        healthMessage='Service Status'
        closeIcon>
        <CardAddonsContainer>
          <CardAddons>
            <CardAddonsTitle>Domains</CardAddonsTitle>
            <Domains />
          </CardAddons>
          <CardAddons>
            <CardAddonsTitle>Routes</CardAddonsTitle>
            <Routes />
          </CardAddons>
          <CardAddons>
            <CardAddonsTitle>Configuration</CardAddonsTitle>
            <Configuration />
          </CardAddons>
        </CardAddonsContainer>
      </SectionCard>
    </React.Fragment>
  );
};
