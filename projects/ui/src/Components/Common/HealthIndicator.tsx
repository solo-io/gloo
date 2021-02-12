import React from 'react';
import styled from '@emotion/styled/macro';
import { StatusType, getHealthColor } from 'utils/health-status';

type HealthIndicatorCircleProps = {
  backgroundColor: string;
  borderColor?: string;
  small?: boolean;
};
const HealthIndicatorCircle = styled.div<HealthIndicatorCircleProps>`
  display: inline-block;

  ${(props: HealthIndicatorCircleProps) =>
    props.small
      ? `height: 11px;
  width: 11px;`
      : `height: 18px;
  width: 18px;`}

  border-radius: 18px;
  ${(props: HealthIndicatorCircleProps) =>
    `background: ${props.backgroundColor}`}
  ${(props: HealthIndicatorCircleProps) =>
    props.borderColor && `border: 1px solid ${props.borderColor};`}
`;

interface Props {
  healthStatus: number;
  statusType?: StatusType;
  issueText?: string;
  small?: boolean;
}

export const HealthIndicator = (props: Props) => {
  const statusType = props.statusType ?? StatusType.DEFAULT;
  const { backgroundColor, borderColor } = getHealthColor(
    props.healthStatus,
    statusType
  );

  return (
    <HealthIndicatorCircle
      backgroundColor={backgroundColor}
      borderColor={borderColor}
      title={props.issueText}
      small={props.small}
    />
  );
};
