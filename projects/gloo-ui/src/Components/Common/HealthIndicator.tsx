import styled from '@emotion/styled';
import * as React from 'react';
import { colors, healthConstants } from 'Styles';

type HealthIndicatorCircleProps = { health: number };
const HealthIndicatorCircle = styled.div`
  display: inline-block;
  height: 18px;
  width: 18px;
  border-radius: 18px;
  margin-left: 10px;
  ${(props: HealthIndicatorCircleProps) =>
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
