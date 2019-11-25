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
import { ErrorBoundary } from 'Components/Features/Errors/ErrorBoundary';

const Container = styled.div`
  padding: 35px 0 20px;
  width: 1275px;
  max-width: 100vw;
  margin: 0 auto;
`;

export const Content = () => (
  <Container>
    <Switch>
      <Route path='/overview'>
        <ErrorBoundary
          fallback={<div>There was an error with the Overview</div>}>
          <Overview />
        </ErrorBoundary>
      </Route>
      <Route
        path='/virtualservices/:virtualservicenamespace/:virtualservicename'
        exact>
        <ErrorBoundary
          fallback={<div>There was an error with Virtual Service Details</div>}>
          <VirtualServiceDetails />
        </ErrorBoundary>
      </Route>
      <Route path='/routetables/:routetablenamespace/:routetablename' exact>
        <ErrorBoundary
          fallback={<div>There was an error with Route Table Details</div>}>
          <RouteTableDetails />
        </ErrorBoundary>
      </Route>
      <Route path='/virtualservices/'>
        <ErrorBoundary
          fallback={
            <div>There was an error with the Virtual Services Listing</div>
          }>
          <VirtualServicesListing />
        </ErrorBoundary>
      </Route>
      <Route path='/upstreams/'>
        <ErrorBoundary
          fallback={<div>There was an error with the Upstreams Listing</div>}>
          <UpstreamsListing />
        </ErrorBoundary>
      </Route>
      <Route path='/admin' exact>
        <ErrorBoundary
          fallback={<div>There was an error with the Admin section</div>}>
          <AdminLanding />
        </ErrorBoundary>
      </Route>
      <Route path='/admin/:sublocation'>
        <ErrorBoundary
          fallback={<div>There was an error with the Admin section</div>}>
          <AdminHub />
        </ErrorBoundary>
      </Route>
      <Route path='/stats/' exact>
        <ErrorBoundary
          fallback={<div>There was an error with the Stats section</div>}>
          <StatsLanding />
        </ErrorBoundary>
      </Route>
      <Route path='/settings/'>
        <ErrorBoundary
          fallback={<div>There was an error with the Settings section</div>}>
          <SettingsLanding />
        </ErrorBoundary>
      </Route>
      <Redirect exact from='/' to='/overview/' />
    </Switch>
  </Container>
);
