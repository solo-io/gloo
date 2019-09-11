import styled from '@emotion/styled';
import { Button } from 'antd';
import { Label } from 'Components/Common/SoloInput';
import { SoloModal } from 'Components/Common/SoloModal';
import * as React from 'react';
import { colors } from 'Styles';
import { RateLimitForm, timeOptions } from './RateLimitForm';
import { RateLimitPlugin } from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/virtualservice_pb';

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
  rateLimits: RateLimitPlugin.AsObject | undefined;
}
export const RateLimit = (props: Props) => {
  const { rateLimits } = props;

  const [showRateLimitModal, setShowRateLimitModal] = React.useState(false);

  return (
    <div>
      <ConfigItemHeader>Rate Limits</ConfigItemHeader>
      <div>
        {!!rateLimits && !!rateLimits.value ? (
          <RatesInfo>
            <InfoBlock>
              <StrongLabel>Authorized Limits:</StrongLabel>
              {rateLimits.value.authorizedLimits ? (
                <InfoBlock>
                  {rateLimits.value.authorizedLimits.requestsPerUnit}
                  <PerText>per</PerText>
                  {
                    timeOptions.find(
                      opt =>
                        opt.value === rateLimits.value!.authorizedLimits!.unit
                    )!.displayValue
                  }
                </InfoBlock>
              ) : (
                'None'
              )}
            </InfoBlock>
            <InfoBlock>
              <StrongLabel>Anonymous Limits:</StrongLabel>
              {!!rateLimits.value.anonymousLimits ? (
                <InfoBlock>
                  {rateLimits.value.anonymousLimits.requestsPerUnit}
                  <PerText>per</PerText>
                  {
                    timeOptions.find(
                      opt =>
                        opt.value === rateLimits.value!.anonymousLimits!.unit
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
          onClick={() => setShowRateLimitModal(show => !show)}>
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
                href='https://gloo.solo.io/gloo_routing/virtual_services/rate_limiting/simple/'
                target='_blank'
                rel='noopener noreferrer'>
                View Rate Limit documentation.
              </a>
            </Legend>
            <RateLimitForm rateLimits={rateLimits} />
          </React.Fragment>
        </SoloModal>
      </div>
    </div>
  );
};
