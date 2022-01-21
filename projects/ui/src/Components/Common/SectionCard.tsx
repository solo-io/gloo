import styled from '@emotion/styled';
import React from 'react';
import { colors } from 'Styles/colors';
import { HealthIndicator } from './HealthIndicator';
import { Card, CardHeader } from './Card';
import { StatusType } from 'utils/health-status';

const CardBlock = styled(Card)`
  margin-bottom: 30px;
  padding: 0;

  @media (max-width: 1380px) {
    margin-bottom: 45px;
  }
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
    width: 30px;
    max-height: 30px;
  }
`;

const HeaderTitleSection = styled.div`
  max-width: calc(100% - 300px);
  min-height: 58px;
`;
const HeaderTitleName = styled.div`
  width: 100%;
  font-size: 22px;
  color: ${colors.novemberGrey};
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
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
  border-radius: 16px;
`;
const SecondaryInformationTitle = styled.span`
  font-weight: bold;
`;

const HealthContainer = styled.div`
  display: flex;
  justify-content: flex-end;
  align-items: center;
  flex: 1;
  text-align: right;
  font-size: 16px;
  font-weight: 600;
  color: ${colors.novemberGrey};
`;
const CloseIcon = styled.div`
  font-size: 21px;
  line-height: 17px;
  margin-left: 23px;
  margin-top: 2px;
  font-weight: 100;
  color: ${colors.juneGrey};
  cursor: pointer;
`;

type BodyContainerProps = { noPadding: boolean };
const BodyContainer = styled.div<BodyContainerProps>`
  padding: ${(props: BodyContainerProps) => (props.noPadding ? '' : '20px;')};
`;

interface Props {
  cardName: string;
  logoIcon?: React.ReactNode;
  headerSecondaryInformation?: {
    title?: string;
    value: React.ReactNode;
  }[];
  health?: {
    state: number;
    type?: StatusType;
    title?: string;
    reason?: string;
  };
  onClose?: () => void;
  secondaryComponent?: React.ReactNode;
  noPadding?: boolean;
}

export const SectionCard: React.FunctionComponent<Props> = props => {
  const {
    logoIcon,
    cardName,
    children,
    headerSecondaryInformation,
    health,
    onClose,
    secondaryComponent,
    noPadding,
  } = props;

  return (
    <CardBlock>
      <CardHeader>
        {logoIcon && <HeaderImageHolder>{logoIcon}</HeaderImageHolder>}
        <HeaderTitleSection>
          <HeaderTitleName>{cardName}</HeaderTitleName>
        </HeaderTitleSection>
        {secondaryComponent}
        {!!headerSecondaryInformation && (
          <SecondaryInformation>
            {headerSecondaryInformation.map(info => {
              return (
                <SecondaryInformationSection key={info.title}>
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
        <HealthContainer>
          {!!health?.title && (
            <div style={{ marginRight: '10px' }}>{health.title}</div>
          )}
          {!!health && (
            <HealthIndicator
              healthStatus={health.state}
              statusType={health.type}
              issueText={health.reason}
            />
          )}
        </HealthContainer>
        {onClose && <CloseIcon onClick={onClose}>X</CloseIcon>}
      </CardHeader>
      <BodyContainer noPadding={!!noPadding}>{children}</BodyContainer>
    </CardBlock>
  );
};
