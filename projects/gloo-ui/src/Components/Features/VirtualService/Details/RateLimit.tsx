import * as React from 'react';
import styled from '@emotion/styled/macro';
import { colors } from 'Styles';
import { Button } from 'antd';
import { SoloModal } from 'Components/Common/SoloModal';
import { RateLimitForm, timeOptions } from './RateLimitForm';
import { IngressRateLimit } from 'proto/github.com/solo-io/solo-projects/projects/gloo/api/v1/plugins/ratelimit/ratelimit_pb';
import { Label } from 'Components/Common/SoloInput';

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

const UpdateButton = styled(Button)`
  padding: 0;
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

interface Props {
  rates: IngressRateLimit.AsObject | undefined;
  rateLimitsChanged: (newRateLimits: IngressRateLimit.AsObject) => any;
}
export const RateLimit = (props: Props) => {
  const { rates, rateLimitsChanged } = props;

  const [showRateLimitModal, setShowRateLimitModal] = React.useState(false);

  return (
    <div>
      <ConfigItemHeader>Rate Limits</ConfigItemHeader>
      <div>
        {!!rates ? (
          <RatesInfo>
            <InfoBlock>
              <StrongLabel>Authorized Limits:</StrongLabel>
              {!!rates.authorizedLimits ? (
                <InfoBlock>
                  {rates.authorizedLimits.requestsPerUnit}
                  <PerText>per</PerText>
                  {
                    timeOptions.find(
                      opt => opt.value === rates.authorizedLimits!.unit
                    )!.displayValue
                  }
                </InfoBlock>
              ) : (
                'None'
              )}
            </InfoBlock>
            <InfoBlock>
              <StrongLabel>Anonymous Limits:</StrongLabel>
              {!!rates.anonymousLimits ? (
                <InfoBlock>
                  {rates.anonymousLimits.requestsPerUnit}
                  <PerText>per</PerText>
                  {
                    timeOptions.find(
                      opt => opt.value === rates.anonymousLimits!.unit
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
        <UpdateButton
          type='link'
          onClick={e => setShowRateLimitModal(show => !show)}>
          Update Rate Limits.
        </UpdateButton>
        <SoloModal
          visible={showRateLimitModal}
          width={650}
          title='Update Rate Limit'
          onClose={() => setShowRateLimitModal(false)}>
          <React.Fragment>
            <Legend>
              You can set different limits for both authorized and anonymous
              users.
              <br />
              Need help?{' '}
              <a
                href='https://gloo.solo.io/enterprise/rate_limiting/ratelimit/'
                target='_blank'>
                View Rate Limit documentation.
              </a>
            </Legend>
            <RateLimitForm
              rates={rates}
              rateLimitsChanged={rateLimitsChanged}
            />
          </React.Fragment>
        </SoloModal>
      </div>
    </div>
  );
};
