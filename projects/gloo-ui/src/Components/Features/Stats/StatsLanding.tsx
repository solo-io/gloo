import * as React from 'react';
/** @jsx jsx */
import { jsx } from '@emotion/core';

import styled from '@emotion/styled/macro';
import { withRouter, RouteComponentProps } from 'react-router';
import { colors } from 'Styles';

export interface RouteParams {
  //... eg, virtualservice?: string
}

function StatsLandingC({
  history,
  match,
  location
}: RouteComponentProps<RouteParams>) {
  return <div>This is the stats landing placeholder...</div>;
}

export const StatsLanding = withRouter(StatsLandingC);
