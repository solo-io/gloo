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
import { RouteTableDetails } from 'Components/Features/VirtualService/RouteTableDetails';

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
        <Route path='/overview'>
          <Overview />
        </Route>
        <Route
          path='/virtualservices/:virtualservicenamespace/:virtualservicename'
          exact>
          <VirtualServiceDetails />
        </Route>
        <Route path='/routetables/:routetablenamespace/:routetablename' exact>
          <RouteTableDetails />
        </Route>
        <Route path='/virtualservices/'>
          <VirtualServicesListing />
        </Route>
        <Route path='/upstreams/'>
          <UpstreamsListing />
        </Route>
        <Route path='/admin' exact>
          <AdminLanding />
        </Route>
        <Route path='/admin/:sublocation'>
          <AdminHub />
        </Route>
        <Route path='/stats/' exact>
          <StatsLanding />
        </Route>
        <Route path='/settings/'>
          <SettingsLanding />
        </Route>
        <Redirect exact from='/' to='/overview/' />
      </Switch>
    </Container>
  );
};
