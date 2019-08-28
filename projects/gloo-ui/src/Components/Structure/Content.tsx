import styled from '@emotion/styled';
import { AdminHub } from 'Components/Features/Admin/AdminHub';
import { AdminLanding } from 'Components/Features/Admin/AdminLanding';
import { Overview } from 'Components/Features/Overview';
import { SettingsLanding } from 'Components/Features/Settings/SettingsLanding';
import { StatsLanding } from 'Components/Features/Stats/StatsLanding';
import { UpstreamsListing } from 'Components/Features/Upstream/UpstreamsListing';
import { VirtualServicesListing } from 'Components/Features/VirtualService/VirtualServicesListing';
import React from 'react';
import { Redirect, Route, Switch } from 'react-router-dom';
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
        <Route path='/overview' render={props => <Overview {...props} />} />
        <Route
          path='/virtualservices/:virtualservicenamespace/:virtualservicename'
          exact
          render={(props: any) => <VirtualServiceDetails {...props} />}
        />
        <Route
          path='/virtualservices/'
          render={(props: any) => <VirtualServicesListing {...props} />}
        />
        <Route
          path='/upstreams/'
          render={(props: any) => <UpstreamsListing {...props} />}
        />
        <Route
          path='/admin'
          exact
          render={props => <AdminLanding {...props} />}
        />
        <Route
          path='/admin/:sublocation'
          render={props => <AdminHub {...props} />}
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
        <Redirect exact from='/' to='/overview/' />
      </Switch>
    </Container>
  );
};
