import React from 'react';
import styled from '@emotion/styled';
import { colors } from 'Styles/colors';
import {
  Card,
  CardSubsectionContent,
  CardSubsectionWrapper,
} from 'Components/Common/Card';
import { ReactComponent as HealthIcon } from 'assets/health-icon.svg';
import {
  AdminClustersBox,
  AdminFederatedResourcesBox,
} from './AdminBoxSummary';
import { AdminGlooInstancesTable } from './AdminGlooInstancesTable';
import { AppName } from 'Components/Common/AppName';

const Heading = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
`;

const HeadingTitle = styled.div`
  font-size: 22px;
  line-height: 26px;
  margin-bottom: 5px;
`;
const HeadingSubtitle = styled.div`
  font-size: 18px;
  line-height: 22px;
`;
const HeadingLogo = styled.div``;

const Section = styled(CardSubsectionWrapper)`
  margin-top: 20px;
`;
const TopSection = styled(Section)`
  display: grid;
  grid-gap: 22px;
  grid-template-columns: 1fr 1fr;
`;
const BottomSection = styled(Section)``;

export const AdminLanding = () => {
  return (
    <Card>
      <Heading>
        <div>
          <HeadingTitle>
            <AppName /> Administration
          </HeadingTitle>
          <HeadingSubtitle>
            Advanced Administration for your Gloo Edge Configuration
          </HeadingSubtitle>
        </div>
        <HeadingLogo>
          <HealthIcon />
        </HeadingLogo>
      </Heading>

      <TopSection>
        <AdminClustersBox />
        <AdminFederatedResourcesBox />
      </TopSection>
      <BottomSection>
        <AdminGlooInstancesTable />
      </BottomSection>
    </Card>
  );
};
