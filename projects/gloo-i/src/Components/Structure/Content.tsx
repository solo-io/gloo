import { Redirect, Route, Switch } from 'react-router-dom';
/** @jsx jsx */
import { jsx } from '@emotion/core';
import styled from '@emotion/styled/macro';

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
          render={(props: any) => <div>Hi!!</div>}
        />
        <Route
          path='/upstreams/'
          exact
          render={(props: any) => <div>Up!!</div>}
        />
        <Route
          path='/stats/'
          exact
          render={(props: any) => <div>Stat!!</div>}
        />
        <Route path='/settings/' render={() => <div>SETTINGS</div>} />

        <Redirect exact from='/' to='/virtualservices/' />
      </Switch>
    </Container>
  );
};
