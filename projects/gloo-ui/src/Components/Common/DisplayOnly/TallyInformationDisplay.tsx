import styled from '@emotion/styled';
import * as React from 'react';
import { RouteComponentProps, withRouter } from 'react-router';
import { colors } from 'Styles';
import { Divider } from 'antd';

type TallyContainerProps = { color: 'orange' | 'blue' };
export const TallyContainer = styled.div`
  display: flex;
  align-items: center;
  padding: 8px 13px;
  line-height: 24px;
  border-radius: 8px;
  margin-bottom: 13px;

  ${(props: TallyContainerProps) =>
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
  line-height: 1;
`;

const TallyDetail = styled<'div', { showLink: boolean }>('div')`
  display: flex;
  width: 100%;
  align-items: center;
  justify-content: flex-start;
  align-items: center;
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
      <TallyDetail showLink={!!props.moreInfoLink}>
        <TallyDescription>{props.tallyDescription}</TallyDescription>
        {!!props.moreInfoLink && (
          <React.Fragment>
            <Divider
              style={{
                height: '20px',
                display: 'flex',
                width: '2px',
                alignSelf: 'center',
                backgroundColor: colors.grapefruitOrange
              }}
              type='vertical'
            />
            <MoreInfoLink onClick={goToMoreInfo}>
              {props.moreInfoLink.prompt}
            </MoreInfoLink>
          </React.Fragment>
        )}
      </TallyDetail>
    </TallyContainer>
  );
};

export const TallyInformationDisplay = withRouter(TallyInformationDisplayC);
