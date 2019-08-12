import {
  Redirect,
  Route,
  Switch,
  RouteComponentProps,
  RouteProps
} from 'react-router-dom';
/** @jsx jsx */
import { jsx } from '@emotion/core';
import styled from '@emotion/styled/macro';
import { VirtualServicesListing } from 'Components/Features/VirtualService/VirtualServicesListing';
import { UpstreamsListing } from 'Components/Features/Upstream/UpstreamsListing';
import { StatsLanding } from 'Components/Features/Stats/StatsLanding';
import { SettingsLanding } from 'Components/Features/Settings/SettingsLanding';
import { VirtualServiceDetails } from '../Features/VirtualService/Details/VirtualServiceDetails';
import { Overview } from 'Components/Features/Overview';
import { Admin } from 'Components/Features/Admin';

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
          path='/virtualservices/'
          exact
          render={(props: any) => <VirtualServicesListing {...props} />}
        />
        <Route
          path='/upstreams/'
          render={(props: any) => <UpstreamsListing {...props} />}
        />
        <Route
          path='/virtualservices/:virtualservicenamespace/:virtualservicename'
          exact
          render={(props: any) => <VirtualServiceDetails {...props} />}
        />
        <Route path='/admin' exact render={props => <Admin {...props} />} />
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
