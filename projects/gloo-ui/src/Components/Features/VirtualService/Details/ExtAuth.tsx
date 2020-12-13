import styled from '@emotion/styled';
import { Button, Spin } from 'antd';
import { Label } from 'Components/Common/SoloInput';
import { SoloModal } from 'Components/Common/SoloModal';
import { OAuth } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/enterprise/options/extauth/v1/extauth_pb';
import * as React from 'react';
import { colors } from 'Styles';
import { ExtAuthForm } from './ExtAuthForm';
import { ExtAuthPlugin } from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/virtualservice_pb';

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

const Unconfigured = styled.div`
  display: inline-block;
  margin-right: 13px;
  line-height: 32px;
`;

const UpdateButton = styled(Button)`
  padding: 0;
`;

const AuthInfo = styled.div`
  display: grid;
  grid-template-columns: 50% 50%;
`;

const InfoBlock = styled.div`
  display: flex;
  line-height: 24px;
  max-width: 95%;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
`;

const StrongLabel = styled(Label)`
  margin-bottom: 0;
  margin-right: 8px;
`;

interface Props {
  externalAuth?: ExtAuthPlugin.AsObject;
}

export const ExtAuth = (props: Props) => {
  const { externalAuth } = props;

  const [showExtAuthModal, setShowExtAuthModal] = React.useState(false);

  return (
    <div>
      <ConfigItemHeader>External Authorization </ConfigItemHeader>
      <div>
        {!!externalAuth &&
        !!externalAuth.value &&
        !!externalAuth.value.authConfigRefName ? (
          <AuthInfo>
            <InfoBlock>
              <StrongLabel>Auth ConfigRef Name:</StrongLabel>
              <InfoBlock>{externalAuth.value.authConfigRefName}</InfoBlock>
            </InfoBlock>
          </AuthInfo>
        ) : (
          <Unconfigured>
            External Authorization has not been configured.
          </Unconfigured>
        )}

        <UpdateButton
          type='link'
          onClick={() => setShowExtAuthModal(show => !show)}>
          Update Authorization.
        </UpdateButton>
        <SoloModal
          visible={showExtAuthModal}
          width={650}
          title='Update Authorization'
          onClose={() => setShowExtAuthModal(false)}>
          <>
            <Legend>
              Prior to creating an OAuth config, you must create a client
              secret.
              <br />
              Need help?{' '}
              <a
                href='https://docs.solo.io/gloo/latest/security/auth/'
                target='_blank'
                rel='noopener noreferrer'>
                View Authorization documentation.
              </a>
            </Legend>
            <ExtAuthForm externalAuth={externalAuth} />
          </>
        </SoloModal>
      </div>
    </div>
  );
};
