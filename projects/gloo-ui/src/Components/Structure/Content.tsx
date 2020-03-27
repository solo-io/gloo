import styled from '@emotion/styled';
import { AdminHub } from 'Components/Features/Admin/AdminHub';
import { AdminLanding } from 'Components/Features/Admin/AdminLanding';
import { APIDetails } from 'Components/Features/DevPortal/apis/APIDetails';
import { DevPortal } from 'Components/Features/DevPortal/DevPortal';
import { DevPortalOverview } from 'Components/Features/DevPortal/DevPortalOverview';
import { PortalDetails } from 'Components/Features/DevPortal/portals/PortalDetails';
import { Overview } from 'Components/Features/Overview';
import { Settings } from 'Components/Features/Settings/SettingsDetails';
import { StatsLanding } from 'Components/Features/Stats/StatsLanding';
import { UpstreamDetails } from 'Components/Features/Upstream/UpstreamDetails';
import { UpstreamGroupDetails } from 'Components/Features/Upstream/UpstreamGroupDetails';
import { UpstreamsListing } from 'Components/Features/Upstream/UpstreamsListing';
import { RouteTableDetails } from 'Components/Features/VirtualService/RouteTableDetails';
import { VirtualServicesListing } from 'Components/Features/VirtualService/VirtualServicesListing';
import { PortalPageEditor } from 'Components/Features/DevPortal/portals/PortalPageEditor';
import React from 'react';
import { Redirect, Route, RouteProps, Switch } from 'react-router-dom';
import { configAPI } from 'store/config/api';
import useSWR from 'swr';
import { VirtualServiceDetails } from '../Features/VirtualService/Details/VirtualServiceDetails';

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

  if (isDeveloperPortalEnabled === undefined) {
    return null;
  }

  return (
    <Route {...props}>
      {isDeveloperPortalEnabled ? props.children : <Redirect to='/overview/' />}
    </Route>
  );
};

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
        <DevPortalRoute
          path='/dev-portal/portals/:portalname/page-editor/:pagename'
          exact>
          <PortalPageEditor />
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
