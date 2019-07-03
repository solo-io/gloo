import * as React from 'react';
/** @jsx jsx */
import { jsx } from '@emotion/core';
import styled from '@emotion/styled/macro';
import { colors, soloConstants } from 'Styles';
import { CardCSS } from 'Styles/CommonEmotions/card';
import { hslToHSLA } from 'Styles/colors';

const Container = styled.div`
  ${CardCSS};
  position: relative;
  width: 235px;
  padding: 0;
  margin-right: 20px;
  height: fit-content;
`;

const MainSection = styled.div`
  padding: 12px;
  border-radius: ${soloConstants.radius}px ${soloConstants.radius}px 0 0;
`;

const CardTitle = styled.div`
  color: ${colors.novemberGrey};
  font-size: 16px;
  font-weight: 600;
`;

const CardSubtitle = styled.div`
  color: ${colors.novemberGrey};
  font-size: 12px;
`;

const Footer = styled.div`
  display: flex;
  justify-content: space-between;
  margin-top: 10px;
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
  margin: 0 12px;
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
`;

export interface CardType {
  cardTitle: string;
  cardSubtitle?: string;
  onRemoveCard?: (id: string) => any;
  id?: string;
  onExpand?: () => any;
  onClick?: () => any;
  details?: {
    title: string;
    value: string;
    valueDisplay?: React.ReactNode | Element;
  }[];
}

export const Card = (props: CardType) => {
  const [expanded, setExpanded] = React.useState<boolean>(false);

  const {
    cardTitle,
    cardSubtitle,
    onRemoveCard,
    onExpand,
    details,
    onClick
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

  return (
    <Container>
      <MainSection>
        <CardTitle>{cardTitle}</CardTitle>
        <CardSubtitle>{cardSubtitle}</CardSubtitle>
      </MainSection>
      <Footer onClick={handleFooterClick}>View Details</Footer>
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

          <Footer onClick={handleFooterClick}>Hide Details</Footer>
        </Expansion>
      )}
    </Container>
  );
};
