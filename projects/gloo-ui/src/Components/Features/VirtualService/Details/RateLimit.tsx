import * as React from 'react';
import styled from '@emotion/styled/macro';
import { colors } from 'Styles';
import { Button } from 'antd';
import { SoloModal } from 'Components/Common/SoloModal';
import { ExtAuthForm } from './ExtAuthForm';
import { RateLimitForm } from './RateLimitForm';

const ConfigItemHeader = styled.div`
  display: flex;
  align-items: flex-start;
  font-weight: 600;
  font-size: 14px;
  color: ${colors.novemberGrey};
`;

const Legend = styled.div`
  background-color: ${colors.januaryGrey};
`;
export const RateLimit = () => {
  const [showRateLimitModal, setShowRateLimitModal] = React.useState(false);

  return (
    <div>
      <ConfigItemHeader>Rate Limits</ConfigItemHeader>
      <div>
        Rate Limits have not been configured.
        <Button type='link' onClick={e => setShowRateLimitModal(show => !show)}>
          Update Rate Limits.
        </Button>
        <SoloModal
          visible={showRateLimitModal}
          width={650}
          title='Update Rate Limit'
          onClose={() => setShowRateLimitModal(false)}>
          <React.Fragment>
            <Legend>
              Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam
              nonumy eirmod tempor. Need help? View Authorization documentation.
            </Legend>
            <RateLimitForm />
          </React.Fragment>
        </SoloModal>
      </div>
    </div>
  );
};
