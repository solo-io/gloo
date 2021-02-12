import React from 'react';
import styled from '@emotion/styled/macro';
import { colors } from 'Styles/colors';
import { Label } from 'Components/Common/SoloInput';
import { IngressRateLimit } from 'proto/github.com/solo-io/solo-apis/api/gloo/gloo/v1/enterprise/options/ratelimit/ratelimit_pb';
import { RateLimit } from 'proto/github.com/solo-io/solo-apis/api/rate-limiter/v1alpha1/ratelimit_pb';

const ConfigItemHeader = styled.div`
  display: flex;
  align-items: flex-start;
  font-weight: 600;
  font-size: 14px;
  color: ${colors.novemberGrey};
`;

const Legend = styled.div`
  background-color: ${colors.januaryGrey};
  padding: 13px 13px 13px 10px;
  margin-bottom: 23px;
`;

const Unconfigured = styled.div`
  display: inline-block;
  margin-right: 13px;
  line-height: 32px;
`;

const RatesInfo = styled.div`
  display: grid;
  grid-template-columns: 1fr 1fr;
`;

const InfoBlock = styled.div`
  display: flex;
  line-height: 24px;
`;

const StrongLabel = styled(Label)`
  margin-bottom: 0;
  margin-right: 8px;
`;

const PerText = styled.div`
  line-height: 24px;
  margin: 0 4px;
  color: ${colors.juneGrey};
`;

const timeOptions = [
  {
    key: RateLimit.Unit.DAY,
    value: RateLimit.Unit.DAY,
    displayValue: 'Day',
  },
  {
    key: RateLimit.Unit.HOUR,
    value: RateLimit.Unit.HOUR,
    displayValue: 'Hour',
  },
  {
    key: RateLimit.Unit.MINUTE,
    value: RateLimit.Unit.MINUTE,
    displayValue: 'Minute',
  },
  {
    key: RateLimit.Unit.SECOND,
    value: RateLimit.Unit.SECOND,
    displayValue: 'Second',
  },
];

interface Props {
  rateLimits?: IngressRateLimit.AsObject;
}
export const RateLimitSection = ({ rateLimits }: Props) => {
  return (
    <div>
      <ConfigItemHeader>Rate Limits</ConfigItemHeader>
      <div>
        {!!rateLimits?.anonymousLimits || !!rateLimits?.authorizedLimits ? (
          <RatesInfo>
            <InfoBlock>
              <StrongLabel>Authorized Limits:</StrongLabel>
              {!!rateLimits.authorizedLimits ? (
                <InfoBlock>
                  {rateLimits.authorizedLimits.requestsPerUnit}
                  <PerText>per</PerText>
                  {
                    timeOptions.find(
                      opt => opt.value === rateLimits.authorizedLimits!.unit
                    )!.displayValue
                  }
                </InfoBlock>
              ) : (
                'None'
              )}
            </InfoBlock>
            <InfoBlock>
              <StrongLabel>Anonymous Limits:</StrongLabel>
              {!!rateLimits.anonymousLimits ? (
                <InfoBlock>
                  {rateLimits.anonymousLimits.requestsPerUnit}
                  <PerText>per</PerText>
                  {
                    timeOptions.find(
                      opt => opt.value === rateLimits.anonymousLimits!.unit
                    )!.displayValue
                  }
                </InfoBlock>
              ) : (
                'None'
              )}
            </InfoBlock>
          </RatesInfo>
        ) : (
          <Unconfigured>Rate Limits have not been configured.</Unconfigured>
        )}
      </div>
    </div>
  );
};
