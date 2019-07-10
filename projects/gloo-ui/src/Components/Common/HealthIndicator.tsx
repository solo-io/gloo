import * as React from 'react';
/** @jsx jsx */
import { jsx } from '@emotion/core';
import styled from '@emotion/styled/macro';
import { colors, healthConstants } from 'Styles';
import { soloConstants } from 'Styles/constants';

const HealthIndicatorCircle = styled<'div', { health: number }>('div')`
  display: inline-block;
  height: 18px;
  width: 18px;
  border-radius: 18px;
  margin-left: 10px;
  ${props =>
    props.health === healthConstants.Good.value
      ? `background: ${colors.forestGreen};`
      : props.health === healthConstants.Error.value
      ? `background: ${colors.grapefruitOrange};`
      : `background: ${colors.sunGold};`}
`;

interface Props {
  healthStatus: number;
}

export const HealthIndicator = (props: Props) => {
  return <HealthIndicatorCircle health={props.healthStatus} />;
};
