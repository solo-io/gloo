import * as React from 'react';
import { RouteComponentProps, withRouter } from 'react-router';

import styled from '@emotion/styled/macro';
import { colors } from 'Styles';
import { ReactComponent as SuccessCheckmark } from 'assets/big-successful-checkmark.svg';

const StateContainer = styled.div`
  display: flex;
  line-height: 19px;
  border-radius: 8px;
  margin-bottom: 12px;
  font-size: 16px;
`;

const SuccessIcon = styled.div`
  margin-right: 10px;

  > svg {
    height: 38px;
    width: 38px;
  }
`;

const Congratulations = styled.div`
  font-size: 18px;
`;

interface Props {
  typeOfItem: string;
}

export const GoodStateCongratulations = (props: Props) => {
  return (
    <StateContainer>
      <SuccessIcon>
        <SuccessCheckmark />
      </SuccessIcon>
      <div>
        <Congratulations>Congratulations</Congratulations>
        All of your {props.typeOfItem} are configured without any issues.
      </div>
    </StateContainer>
  );
};
