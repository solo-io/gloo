import styled from '@emotion/styled';
import { AdminHub } from 'Components/Features/Admin/AdminHub';
import { AdminLanding } from 'Components/Features/Admin/AdminLanding';
import { APIDetails } from 'Components/Features/DevPortal/apis/APIDetails';
import { APIListing } from 'Components/Features/DevPortal/apis/APIListing';
import { DevPortalOverview } from 'Components/Features/DevPortal/DevPortalOverview';
import { PortalDetails } from 'Components/Features/DevPortal/portals/PortalDetails';
import { PortalsListing } from 'Components/Features/DevPortal/portals/PortalsListing';
import { SwaggerExplorer } from 'Components/Features/DevPortal/SwaggerExplorer';
import { ErrorBoundary } from 'Components/Features/Errors/ErrorBoundary';
import { Overview } from 'Components/Features/Overview';
import { Settings } from 'Components/Features/Settings/SettingsDetails';
import { SettingsOverview } from 'Components/Features/Settings/SettingsOverview';
import { StatsLanding } from 'Components/Features/Stats/StatsLanding';
import { UpstreamDetails } from 'Components/Features/Upstream/UpstreamDetails';
import { UpstreamGroupDetails } from 'Components/Features/Upstream/UpstreamGroupDetails';
import { UpstreamsListing } from 'Components/Features/Upstream/UpstreamsListing';
import { RouteTableDetails } from 'Components/Features/VirtualService/RouteTableDetails';
import { VirtualServicesListing } from 'Components/Features/VirtualService/VirtualServicesListing';
import React from 'react';
import { Redirect, Route, Switch, RouteProps } from 'react-router-dom';
import { VirtualServiceDetails } from '../Features/VirtualService/Details/VirtualServiceDetails';
import { UserGroups } from 'Components/Features/DevPortal/users/UserGroups';
import { DevPortal } from 'Components/Features/DevPortal/DevPortal';
import useSWR, { cache } from 'swr';
import { configAPI } from 'store/config/api';
import { WatchedNamespacesPage } from 'Components/Features/Settings/WatchedNamespacesPage';
import { SecretsPage } from 'Components/Features/Settings/SecretsPage';

const Container = styled.div`
  padding: 35px 0 20px;
  width: 1275px;
  max-width: 100vw;
  margin: 0 auto;
`;

const DevPortalRoute: React.FC<RouteProps> = props => {
  const {
    data: isDeveloperPortalEnabled,
    error: isDeveloperPortalEnabledError
  } = useSWR('isDeveloperPortalEnabled', configAPI.isDeveloperPortalEnabled, {
    refreshInterval: 0
  });
  return (
    <Route {...props}>
      {isDeveloperPortalEnabled ? props.children : <Redirect to='/overview/' />}
    </Route>
  );
};

export const Content = () => {
  const {
    data: isDeveloperPortalEnabled,
    error: isDeveloperPortalEnabledError
  } = useSWR('isDeveloperPortalEnabled', configAPI.isDeveloperPortalEnabled, {
    refreshInterval: 0
  });

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
        <Route path='/upstreams/:upstreamnamespace/:upstreamname' exact>
          <UpstreamDetails />
        </Route>
        <Route
          path='/upstreams/upstreamgroups/:upstreamgroupnamespace/:upstreamgroupname'
          exact>
          <UpstreamGroupDetails />
        </Route>
        <Route path='/upstreams/'>
          <UpstreamsListing />
        </Route>
        <Route path='/admin' exact>
          <AdminLanding />
        </Route>
        <DevPortalRoute path='/dev-portal' exact>
          <DevPortalOverview />
        </DevPortalRoute>
        <DevPortalRoute path='/dev-portal/portals' exact>
          <DevPortal />
        </DevPortalRoute>
        <DevPortalRoute path='/dev-portal/apis' exact>
          <DevPortal />
        </DevPortalRoute>
        <DevPortalRoute path='/dev-portal/users' exact>
          <DevPortal />
        </DevPortalRoute>
        <DevPortalRoute path='/dev-portal/api-key-scopes' exact>
          <DevPortal />{' '}
        </DevPortalRoute>
        <DevPortalRoute path='/dev-portal/api-keys' exact>
          <DevPortal />
        </DevPortalRoute>
        <DevPortalRoute path='/dev-portal/portals/:portalname' exact>
          <PortalDetails />
        </DevPortalRoute>
        <DevPortalRoute path='/dev-portal/apis/:apiname'>
          <APIDetails />
        </DevPortalRoute>
        <Route path='/admin/:sublocation'>
          <AdminHub />
        </Route>
        <Route path='/stats/' exact>
          <StatsLanding />
        </Route>

        <Route path='/settings/:settingsnamespace/:settingsname'>
          <Settings />
        </Route>
        <Redirect exact from='/routetables/' to='/virtualservices/' />
        <Redirect exact from='/' to='/overview/' />
      </Switch>
    </Container>
  );
};
