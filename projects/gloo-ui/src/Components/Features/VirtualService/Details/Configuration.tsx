import * as React from 'react';
import styled from '@emotion/styled/macro';
import { colors } from 'Styles';
import { DetailsSectionTitle } from './VirtualServiceDetails';
import { Button } from 'antd';
import { SoloModal } from 'Components/Common/SoloModal';
import { ExtAuthForm } from './ExtAuthForm';

const ConfigContainer = styled.div`
  display: grid;
  grid-template-columns: 1fr 1fr;
  background: ${colors.januaryGrey};
  height: 80%;
`;

const ConfigItem = styled.div`
  margin: 20px;
  padding: 10px;
  justify-items: center;
  background: white;
`;

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

export const Configuration = () => {
  const [showExtAuthModal, setShowExtAuthModal] = React.useState(false);
  const [showRateLimitModal, setShowRateLimitModal] = React.useState(false);
  return (
    <React.Fragment>
      <DetailsSectionTitle>Configuration</DetailsSectionTitle>
      <ConfigContainer>
        <ConfigItem>
          <ConfigItemHeader>External Authorization </ConfigItemHeader>
          <div>
            External Authorization has not been configured.
            <Button
              type='link'
              onClick={e => setShowExtAuthModal(show => !show)}>
              Update Authorization.
            </Button>
            <SoloModal
              visible={showExtAuthModal}
              width={650}
              title='Update Authorization'
              onClose={() => setShowExtAuthModal(false)}>
              <React.Fragment>
                <Legend>
                  Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed
                  diam nonumy eirmod tempor. Need help? View Authorization
                  documentation.
                </Legend>
                <ExtAuthForm />
              </React.Fragment>
            </SoloModal>
          </div>
        </ConfigItem>
        <ConfigItem>
          <ConfigItemHeader>Rate Limits</ConfigItemHeader>
          <div>
            Rate Limits have not been configured.
            <Button
              type='link'
              onClick={e => setShowRateLimitModal(show => !show)}>
              Update Rate Limits.
            </Button>
            <SoloModal
              visible={showRateLimitModal}
              width={650}
              title='Update Authorization'
              onClose={() => setShowRateLimitModal(false)}>
              Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam
              nonumy eirmod tempor. Need help? View Authorization documentation.
            </SoloModal>
          </div>
        </ConfigItem>
      </ConfigContainer>
    </React.Fragment>
  );
};
