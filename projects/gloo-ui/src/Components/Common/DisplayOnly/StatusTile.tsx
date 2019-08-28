import styled from '@emotion/styled';
import * as React from 'react';
import { RouteComponentProps, withRouter } from 'react-router';
import { colors, soloConstants } from 'Styles';
import { HealthIndicator } from '../HealthIndicator';

const StatusTileContainer = styled.div`
  height: 100%;
  border-radius: 8px;
  background: ${colors.januaryGrey};
  padding: 17px;
`;

type StatusTileInformationProps = { horizontal?: boolean };
const StatusTileInformation = styled.div`
  position: relative;
  border-radius: 8px;
  background: white;
  padding: 15px 18px 18px 15px;
  height: 100%;

  ${(props: StatusTileInformationProps) =>
    props.horizontal ? `display: flex;` : ''}
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
  /* margin-left: ${soloConstants.largeBuffer}px; */
  margin: 10px;

  svg {
    width: 120px;
    height: 120px;
    border-radius: 50%;
    border: 1px solid ${colors.marchGrey};

  }
`;

type DescriptionProps = { minHeight?: string };
const Description = styled.div`
  color: ${colors.novemberGrey};
  font-size: 16px;
  line-height: 19px;
  margin-bottom: 23px;
  min-height: ${(props: DescriptionProps) =>
    props.minHeight !== undefined ? props.minHeight : '0'};
`;

const Content = styled.div`
  padding-bottom: 5px;
`;

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
    top: 55%;
    margin-top: -${soloConstants.largeBuffer}px;
    border-right: ${soloConstants.smallBuffer}px solid ${colors.januaryGrey};
    border-top: 13px solid transparent;
    border-bottom: 13px solid transparent;
  }
`;

const Link = styled.div`
  position: absolute;
  bottom: 10px;
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
  descriptionMinHeight?: string;
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
            <Description minHeight={props.descriptionMinHeight}>
              {props.description}
            </Description>
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
