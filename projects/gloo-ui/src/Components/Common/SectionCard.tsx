import * as React from 'react';
/** @jsx jsx */
import { jsx } from '@emotion/core';
import styled from '@emotion/styled/macro';

import { colors, soloConstants } from 'Styles';
import { CardCSS } from 'Styles/CommonEmotions/card';
import { HealthIndicator } from './HealthIndicator';

const CardBlock = styled.div`
  ${CardCSS};
  margin-bottom: 30px;
  padding: 0;
  @media (max-width: 1380px) {
    margin-bottom: 45px;
  }
`;

const Header = styled.div`
  display: flex;
  align-items: center;
  width: 100%;
  background: ${colors.marchGrey};
  padding: 13px;
  border-radius: ${soloConstants.radius}px ${soloConstants.radius}px 0 0;
`;

const HeaderImageHolder = styled.div`
  margin-right: 15px;
  height: 33px;
  width: 33px;
  border-radius: 100%;
  background: white;
  display: flex;
  justify-content: center;
  align-items: center;

  img,
  svg {
    width: 20px;
    max-height: 25px;
  }
`;

const HeaderTitleSection = styled.div`
  max-width: calc(100% - 300px);
`;
const HeaderTitleName = styled.div`
  width: 100%;
  font-size: 22px;
  color: ${colors.novemberGrey};
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
  text-transform: capitalize;
`;

const SecondaryInformation = styled.div`
  display: flex;
  align-items: center;
`;
const SecondaryInformationSection = styled.div`
  font-size: 14px;
  line-height: 22px;
  height: 22px;
  padding: 0 12px;
  color: ${colors.novemberGrey};
  background: white;
  margin-left: 13px;
  border-radius: ${soloConstants.largeRadius}px;
`;
const SecondaryInformationTitle = styled.span`
  font-weight: bold;
`;

const HealthContainer = styled.div`
  flex: 1;
  text-align: right;
  font-size: 16px;
  color: ${colors.novemberGrey};
`;

const BodyContainer = styled.div`
  padding: 20px;
`;

const CardAddonsContainer = styled.div`
  padding: ${soloConstants.buffer}px ${soloConstants.smallBuffer}px 10px;
  display: grid;
  grid-template-columns: 1fr 1fr;
  grid-column-gap: 30px;
`;

const CardAddons = styled.div`
  width: 100%;
`;
const CardAddonsTitle = styled.div`
  font-size: 18px;
  font-weight: bold;
  color: ${colors.novemberGrey};
  margin-bottom: 10px;
`;

const MoreTease = styled.div`
  margin: 10px 0 0;
  font-size: 14px;
  color: ${colors.septemberGrey};
`;

interface Props {
  cardName: string;
  logoIcon?: React.ReactNode;
  headerSecondaryInformation?: {
    title?: string;
    value: string;
  }[];
  health?: number;
  healthMessage?: string;
  closeIcon?: boolean;
}

export const SectionCard: React.FunctionComponent<Props> = props => {
  const {
    logoIcon,
    cardName,
    children,
    headerSecondaryInformation,
    health,
    healthMessage,
    closeIcon
  } = props;

  return (
    <CardBlock>
      <Header>
        {logoIcon && <HeaderImageHolder>{logoIcon}</HeaderImageHolder>}
        <HeaderTitleSection>
          <HeaderTitleName>{cardName}</HeaderTitleName>
        </HeaderTitleSection>

        {!!headerSecondaryInformation && (
          <SecondaryInformation>
            {headerSecondaryInformation.map(info => {
              return (
                <SecondaryInformationSection key={info.value}>
                  {!!info.title && (
                    <SecondaryInformationTitle>
                      {info.title}:{' '}
                    </SecondaryInformationTitle>
                  )}
                  {info.value}
                </SecondaryInformationSection>
              );
            })}
          </SecondaryInformation>
        )}
        {health && (
          <HealthContainer>
            {healthMessage || ''}
            <HealthIndicator healthStatus={health} />
          </HealthContainer>
        )}
        {closeIcon && <div style={{ padding: '10px' }}>X</div>}
      </Header>
      <BodyContainer>{children}</BodyContainer>
    </CardBlock>
  );
};
