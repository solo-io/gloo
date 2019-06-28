import { Redirect, Route, Switch } from 'react-router-dom';
/** @jsx jsx */
import { jsx } from '@emotion/core';
import styled from '@emotion/styled/macro';
import { VirtualServicesListing } from 'Components/Features/VirtualService/VirtualServicesListing';
import { UpstreamsListing } from 'Components/Features/Upstream/UpstreamsListing';
import { StatsLanding } from 'Components/Features/Stats/StatsLanding';
import { SettingsLanding } from 'Components/Features/Settings/SettingsLanding';
import { VirtualServiceDetails } from '../Features/VirtualService/Details/VirtualServiceDetails';

const Container = styled.div`
  padding: 35px 0 20px;
  width: 1275px;
  max-width: 100vw;
  margin: 0 auto;
`;
export const Content = () => {
  return (
    <Container>
      <Switch>
        <Route
          path='/virtualservices/'
          exact
          render={(props: any) => <VirtualServicesListing {...props} />}
        />
        <Route
          path='/virtualservices/:virtualservicename/'
          exact
          render={(props: any) => <VirtualServicesListing {...props} />}
        />
        <Route
          path='/upstreams/'
          exact
          render={(props: any) => <UpstreamsListing {...props} />}
        />
        <Route
          path='/virtualservices/:virtualservicename/details'
          exact
          render={(props: any) => <VirtualServiceDetails {...props} />}
        />
        <Route
          path='/stats/'
          exact
          render={(props: any) => <StatsLanding {...props} />}
        />
        <Route
          path='/settings/'
          render={(props: any) => <SettingsLanding {...props} />}
        />
        <Redirect exact from='/' to='/virtualservices/' />
      </Switch>
    </Container>
  );
};
