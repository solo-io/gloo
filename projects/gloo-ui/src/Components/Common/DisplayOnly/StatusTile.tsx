import * as React from 'react';
import { RouteComponentProps, withRouter } from 'react-router';

import styled from '@emotion/styled/macro';
import { colors, soloConstants } from 'Styles';
import { HealthIndicator } from '../HealthIndicator';

const StatusTileContainer = styled.div`
  height: 100%;
  border-radius: 8px;
  background: ${colors.januaryGrey};
  padding: 17px;
`;

const StatusTileInformation = styled<'div', { horizontal?: boolean }>('div')`
  position: relative;
  border-radius: 8px;
  background: white;
  padding: 15px 18px 39px 15px;
  height: 100%;

  ${props => (props.horizontal ? `display: flex;` : '')}
`;

const Title = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-size: 20px;
  line-height: 24px;
  color: ${colors.novemberGrey};
  margin-bottom: 10px;

  svg {
    height: 24px;
    margin-left: 8px;
  }
`;

const HorizontalTitle = styled.div`
  display: flex;
  align-items: center;
  height: 140px;
  margin-left: ${soloConstants.largeBuffer}px;
  margin-right: 30px;

  svg {
    width: 85px;
  }
`;

const Description = styled.div`
  color: ${colors.novemberGrey};
  font-size: 16px;
  line-height: 19px;
  margin-bottom: 23px;
  min-height: 95px;
`;

const Content = styled.div``;

const HorizontalContent = styled.div`
  position: relative;
  flex: 1;
  border-radius: 8px;
  margin-left: 23px;
  background: ${colors.januaryGrey};
  padding: ${soloConstants.buffer}px ${soloConstants.smallBuffer}px;

  &:before {
    content: '';
    position: absolute;
    left: -${soloConstants.smallBuffer}px;
    top: 50%;
    margin-top: -${soloConstants.largeBuffer}px;
    border-right: ${soloConstants.smallBuffer}px solid ${colors.januaryGrey};
    border-top: 13px solid transparent;
    border-bottom: 13px solid transparent;
  }
`;

const Link = styled.div`
  position: absolute;
  bottom: 15px;
  cursor: pointer;
  color: ${colors.seaBlue};
  font-size: 14px;
`;

interface Props extends RouteComponentProps {
  titleText?: string;
  titleIcon?: Element | React.ReactElement;
  description?: string;
  children?: React.ReactChild;
  exploreMoreLink?: {
    prompt: string;
    link: string;
  };
  horizontal?: boolean;
  healthStatus?: number;
}

const StatusTileC = (props: Props) => {
  const goToLink = (): void => {
    props.history.push(props.exploreMoreLink!.link);
  };

  return (
    <StatusTileContainer>
      <StatusTileInformation horizontal={props.horizontal}>
        {!props.horizontal ? (
          <React.Fragment>
            <Title>
              <div>
                {props.titleText}
                {props.titleIcon}
              </div>

              {props.healthStatus !== undefined && (
                <HealthIndicator healthStatus={props.healthStatus} />
              )}
            </Title>
            <Description>{props.description}</Description>
            <Content>{props.children}</Content>
            {!!props.exploreMoreLink && (
              <Link onClick={goToLink}>{props.exploreMoreLink.prompt}</Link>
            )}
          </React.Fragment>
        ) : (
          <React.Fragment>
            <HorizontalTitle>
              {props.titleText}
              {props.titleIcon}
            </HorizontalTitle>
            <HorizontalContent>
              {props.children}
              {!!props.exploreMoreLink && (
                <Link onClick={goToLink}>{props.exploreMoreLink.prompt}</Link>
              )}
            </HorizontalContent>
          </React.Fragment>
        )}
      </StatusTileInformation>
    </StatusTileContainer>
  );
};

export const StatusTile = withRouter(StatusTileC);
