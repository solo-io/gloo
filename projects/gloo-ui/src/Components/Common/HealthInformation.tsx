import * as React from 'react';
/** @jsx jsx */
import { jsx } from '@emotion/core';
import styled from '@emotion/styled/macro';
import { colors, healthConstants } from 'Styles';
import { soloConstants } from 'Styles/constants';
import { ReactComponent as InfoPrompt } from 'assets/info-prompt.svg';
import { Status } from 'proto/github.com/solo-io/solo-kit/api/v1/status_pb';
import { hslToHSLA } from 'Styles/colors';

const Health = styled<'div', { health: number }>('div')`
  position: relative;
  color: ${props =>
    props.health === healthConstants.Good.value
      ? colors.forestGreen
      : props.health === healthConstants.Error.value
      ? colors.grapefruitOrange
      : colors.sunGold};
`;

const InfoPromptContainer = styled.div`
  position: absolute;
  right: -16px;
  top: 2px;
`;

const ExtraInfo = styled<'div', { health: number }>('div')`
  position: absolute;
  left: 15px;
  top: 15px;
  width: 125px;
  padding: 5px 8px;
  word-break: break-word;
  border: 1px solid
    ${props =>
      hslToHSLA(
        props.health === healthConstants.Good.value
          ? colors.forestGreen
          : props.health === healthConstants.Error.value
          ? colors.grapefruitOrange
          : colors.sunGold,
        0.9
      )};
  background: ${props =>
    hslToHSLA(
      props.health === healthConstants.Good.value
        ? colors.groveGreen
        : props.health === healthConstants.Error.value
        ? colors.tangerineOrange
        : colors.flashlightGold,
      0.9
    )};
  z-index: 2;
`;

interface Props {
  healthStatus: Status.AsObject | undefined;
}

export const HealthInformation = (props: Props) => {
  const [showExtraInfo, setShowExtraInfo] = React.useState(false);

  const usedState = props.healthStatus
    ? props.healthStatus.state
    : healthConstants.Pending.value;
  const status =
    // @ts-ignore
    healthConstants[
      Object.keys(healthConstants).find(key => {
        // @ts-ignore
        return healthConstants[key].value === usedState;
      })!
    ];

  return (
    <Health health={usedState}>
      {status.text}
      {!!props.healthStatus &&
        !!props.healthStatus.reason &&
        !!props.healthStatus.reason.length && (
          <InfoPromptContainer>
            <InfoPrompt
              onClick={() => setShowExtraInfo(s => !s)}
              style={{ cursor: 'pointer' }}
            />
            {showExtraInfo && (
              <ExtraInfo health={usedState}>
                {props.healthStatus.reason}
              </ExtraInfo>
            )}
          </InfoPromptContainer>
        )}
    </Health>
  );
};
