import * as React from 'react';
import { RouteComponentProps, withRouter } from 'react-router';

import styled from '@emotion/styled/macro';
import { colors } from 'Styles';

export const TallyContainer = styled<'div', { color: 'orange' | 'blue' }>(
  'div'
)`
  display: flex;
  padding: 8px 13px;
  line-height: 24px;
  border-radius: 8px;
  margin-bottom: 13px;

  ${props =>
    props.color === 'orange'
      ? `border: 1px solid ${colors.grapefruitOrange};
    background: ${colors.tangerineOrange};
    color: ${colors.pumpkinOrange};`
      : `border: 1px solid ${colors.lakeBlue};
    background: ${colors.splashBlue};
    color: ${colors.seaBlue};`}
`;

const TallyCount = styled.div`
  font-size: 24px;
  font-weight: 600;
  margin-right: 7px;
`;

const TallyDescription = styled.div`
  font-size: 16px;
  line-height: 16px;
`;

const MoreInfoLink = styled.div`
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
`;

interface Props extends RouteComponentProps {
  tallyCount: number | null;
  tallyDescription: string | null;
  moreInfoLink?: {
    prompt: string;
    link: string;
  };
  color: 'orange' | 'blue';
}

const TallyInformationDisplayC = (props: Props) => {
  const goToMoreInfo = (): void => {
    props.history.push(props.moreInfoLink!.link);
  };

  const countDisplay = (): string | number => {
    const count = props.tallyCount;

    if (count !== null) {
      if (count < 10) {
        return `0${count}`;
      } else {
        return count;
      }
    } else {
      return '?';
    }
  };

  return (
    <TallyContainer color={props.color}>
      <TallyCount>{countDisplay()}</TallyCount>
      <div>
        <TallyDescription>{props.tallyDescription}</TallyDescription>
        {!!props.moreInfoLink && (
          <MoreInfoLink onClick={goToMoreInfo}>
            {props.moreInfoLink.prompt}
          </MoreInfoLink>
        )}
      </div>
    </TallyContainer>
  );
};

export const TallyInformationDisplay = withRouter(TallyInformationDisplayC);
