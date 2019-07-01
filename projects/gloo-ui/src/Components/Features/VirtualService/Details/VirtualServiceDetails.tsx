import * as React from 'react';
import { SectionCard } from 'Components/Common/SectionCard';
import { ReactComponent as GlooIcon } from 'assets/GlooEE.svg';
import { Domains } from './Domains';
import { Routes } from './Routes';
import { Configuration } from './Configuration';
import styled from '@emotion/styled/macro';
import { colors, soloConstants } from 'Styles';
import { RouteComponentProps } from 'react-router';
import { Breadcrumb } from 'Components/Common/Breadcrumb';

const DetailsContent = styled.div`
  display: grid;
  grid-template-rows: auto 2fr 1fr;
  grid-column-gap: 30px;
`;

const DetailsSection = styled.div`
  width: 100%;
`;
export const DetailsSectionTitle = styled.div`
  font-size: 18px;
  font-weight: bold;
  color: ${colors.novemberGrey};
  margin-top: 10px;
  margin-bottom: 10px;
`;

interface Props extends RouteComponentProps<{ virtualservicename: string }> {}
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

export const VirtualServiceDetails = (props: Props) => {
  const { match } = props;
  return (
    <React.Fragment>
      <Breadcrumb />

      <SectionCard
        cardName={match.params ? match.params.virtualservicename : 'test'}
        logoIcon={<GlooIcon />}
        health={1}
        headerSecondaryInformation={headerInfo}
        healthMessage='Service Status'
        closeIcon>
        <DetailsContent>
          <DetailsSection>
            <Domains />
          </DetailsSection>
          <DetailsSection>
            <Routes />
          </DetailsSection>
          <DetailsSection>
            <Configuration />
          </DetailsSection>
        </DetailsContent>
      </SectionCard>
    </React.Fragment>
  );
};
