import * as React from 'react';
/** @jsx jsx */
import { jsx } from '@emotion/core';
import styled from '@emotion/styled/macro';
import {
  colors,
  soloConstants,
  healthConstants,
  TableActionCircle
} from 'Styles';
import { CardCSS } from 'Styles/CommonEmotions/card';
import { hslToHSLA } from 'Styles/colors';
import { HealthIndicator } from './HealthIndicator';
import { Popconfirm } from 'antd';

const Container = styled.div`
  ${CardCSS};
  position: relative;
  width: 235px;
  padding: 0;
  margin-right: 20px;
  height: fit-content;
`;

const MainSection = styled.div`
  padding: 12px 6px 12px 12px;
  border-radius: ${soloConstants.radius}px ${soloConstants.radius}px 0 0;
`;

const CardTitle = styled.div`
  display: flex;
  justify-content: space-between;
  color: ${colors.novemberGrey};
  font-size: 16px;
  line-height: 20px;
  font-weight: 600;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
`;

const CardTitleText = styled.div`
  max-width: 175px;
  display: block;
  word-break: break-all;
  white-space: normal;
`;

const CardSubtitle = styled.div`
  color: ${colors.novemberGrey};
  font-size: 12px;
  min-height: 18px;
  max-height: 18px;
  line-height: 18px;
  overflow: hidden;
  text-overflow: ellipsis;
  word-break: break-all;
`;

const Footer = styled.div`
  position: relative;
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: ${hslToHSLA(colors.marchGrey, 0.15)};
  color: ${colors.seaBlue};
  font-size: 14px;
  line-height: 30px;
  height: 30px;
  padding: 0 6px 0 12px;
  border-radius: 0 0 ${soloConstants.radius}px ${soloConstants.radius}px;
  cursor: pointer;
`;

const Expansion = styled.div`
  position: absolute;
  top: calc(100% - 30px);
  left: 0;
  right: 0;
  background: white;
  box-shadow: 0px 5px 6px ${colors.darkerBoxShadow};
  border-radius: 0 0 ${soloConstants.radius}px ${soloConstants.radius}px;
  z-index: 2;
`;

const ExpandedDetails = styled.div`
  margin: 0 12px 18px;
  border-top: 1px solid ${colors.aprilGrey};
`;

const Detail = styled.div`
  display: flex;
  margin-top: 10px;
  font-size: 12px;
`;

const DetailTitle = styled.div`
  color: ${colors.novemberGrey};
  font-weight: 600;
  width: 70px;
`;

const DetailContent = styled.div`
  color: ${colors.septemberGrey};
  text-overflow: clip;
  white-space: normal;
  word-break: break-all;
`;

const ActionCircle = styled(TableActionCircle)`
  display: inline-block;
  margin: 1px;
  color: ${colors.septemberGrey};
`;

const ArrowToggle = styled<'div', { active?: boolean }>('div')`
  position: absolute;
  right: 8px;
  top: ${props => (props.active ? '14' : '15')}px;

  &:before,
  &:after {
    position: absolute;
    content: '';
    display: block;
    width: 8px;
    height: 1px;
    background: white;
    transition: transform 0.5s;
  }

  &:before {
    right: 5px;
    border-top-left-radius: 10px;
    border-bottom-left-radius: 10px;
    transform: rotate(${props => (props.active ? '-' : '')}45deg);
  }

  &:after {
    right: 1px;
    transform: rotate(${props => (props.active ? '' : '-')}45deg);
  }
`;

export interface CardType {
  cardTitle: string;
  cardSubtitle?: string;
  onRemoveCard?: () => any;
  id?: string;
  onExpand?: () => any;
  onClick?: () => any;
  details?: {
    title: string;
    value: string;
    valueDisplay?: React.ReactNode | Element;
  }[];
  healthStatus?: number;
  onCreate?: () => any;
  extraInfoComponent?: React.FC;
}

export const Card = (props: CardType) => {
  const [expanded, setExpanded] = React.useState<boolean>(false);

  const {
    cardTitle,
    cardSubtitle,
    onRemoveCard,
    onExpand,
    details,
    onClick,
    healthStatus,
    onCreate,
    extraInfoComponent
  } = props;

  const handleFooterClick = () => {
    if (onClick) {
      onClick();
    } else {
      if (!!onExpand && !expanded) {
        onExpand();
      }

      setExpanded(exp => !exp);
    }
  };

  let ExtraInformation = null;
  if (!!extraInfoComponent) {
    ExtraInformation = extraInfoComponent;
  }

  return (
    <Container>
      <MainSection>
        <CardTitle>
          <CardTitleText>{cardTitle}</CardTitleText>
          <div>
            {!!onCreate && <ActionCircle onClick={onCreate}>+</ActionCircle>}
            <Popconfirm
              title='Are you sure you want to delete this virtual service？'
              okText='Yes'
              cancelText='No'>
              {!!onRemoveCard && (
                <ActionCircle onClick={onRemoveCard}>x</ActionCircle>
              )}
            </Popconfirm>
          </div>
        </CardTitle>
        <CardSubtitle>
          {cardSubtitle && cardSubtitle.length ? cardSubtitle : '   '}
        </CardSubtitle>
      </MainSection>
      <Footer onClick={handleFooterClick}>
        <span>View Details</span>
        <HealthIndicator
          healthStatus={healthStatus || healthConstants.Pending.value}
        />
        {!!onExpand && <ArrowToggle />}
      </Footer>
      {expanded && (
        <Expansion>
          <ExpandedDetails>
            {details &&
              details.map(detail => {
                return (
                  <Detail key={detail.title}>
                    <DetailTitle>{detail.title}:</DetailTitle>
                    <DetailContent>
                      {!!detail.valueDisplay
                        ? detail.valueDisplay
                        : detail.value}
                    </DetailContent>
                  </Detail>
                );
              })}
          </ExpandedDetails>
          {!!ExtraInformation && <ExtraInformation />}
          <Footer onClick={handleFooterClick}>
            <span>Hide Details</span>
            <HealthIndicator
              healthStatus={healthStatus || healthConstants.Pending.value}
            />
            <ArrowToggle active={true} />
          </Footer>
        </Expansion>
      )}
    </Container>
  );
};
