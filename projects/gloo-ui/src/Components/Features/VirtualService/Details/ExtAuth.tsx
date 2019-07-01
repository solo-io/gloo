import * as React from 'react';
import styled from '@emotion/styled/macro';
import { colors } from 'Styles';
import { Button } from 'antd';
import { SoloModal } from 'Components/Common/SoloModal';
import { ExtAuthForm } from './ExtAuthForm';

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

export const ExtAuth = () => {
  const [showExtAuthModal, setShowExtAuthModal] = React.useState(false);
  return (
    <React.Fragment>
      <ConfigItemHeader>External Authorization </ConfigItemHeader>
      <div>
        External Authorization has not been configured.
        <Button type='link' onClick={e => setShowExtAuthModal(show => !show)}>
          Update Authorization.
        </Button>
        <SoloModal
          visible={showExtAuthModal}
          width={650}
          title='Update Authorization'
          onClose={() => setShowExtAuthModal(false)}>
          <React.Fragment>
            <Legend>
              Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam
              nonumy eirmod tempor. Need help? View Authorization documentation.
            </Legend>
            <ExtAuthForm />
          </React.Fragment>
        </SoloModal>
      </div>
    </React.Fragment>
  );
};
