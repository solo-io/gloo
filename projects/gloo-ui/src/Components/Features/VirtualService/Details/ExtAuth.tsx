import * as React from 'react';
import styled from '@emotion/styled/macro';
import { colors } from 'Styles';
import { Button } from 'antd';
import { SoloModal } from 'Components/Common/SoloModal';
import { ExtAuthForm } from './ExtAuthForm';
import { OAuth } from 'proto/github.com/solo-io/solo-projects/projects/gloo/api/v1/plugins/extauth/extauth_pb';
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

interface Props {
  externalAuth: OAuth.AsObject | undefined;
  externalAuthChanged: (newExternalAuth: OAuth.AsObject) => any;
}
export const ExtAuth = (props: Props) => {
  const { externalAuth, externalAuthChanged } = props;

  const [showExtAuthModal, setShowExtAuthModal] = React.useState(false);

  return (
    <div>
      <ConfigItemHeader>External Authorization </ConfigItemHeader>
      <div>
        {!!externalAuth ? (
          <AuthInfo>
            <InfoBlock>
              <StrongLabel>Client ID:</StrongLabel>
              <InfoBlock>{externalAuth.clientId}</InfoBlock>
            </InfoBlock>
            <InfoBlock>
              <StrongLabel>Callback Path:</StrongLabel>
              <InfoBlock>{externalAuth.callbackPath}</InfoBlock>
            </InfoBlock>
            <InfoBlock>
              <StrongLabel>Issuer URL:</StrongLabel>
              <InfoBlock>{externalAuth.issuerUrl}</InfoBlock>
            </InfoBlock>
            <InfoBlock>
              <StrongLabel>App URL:</StrongLabel>
              <InfoBlock>{externalAuth.appUrl}</InfoBlock>
            </InfoBlock>
            <InfoBlock>
              <StrongLabel>Secret Ref Name:</StrongLabel>
              {!!externalAuth.clientSecretRef ? (
                <InfoBlock>{externalAuth.clientSecretRef.name}</InfoBlock>
              ) : (
                'None'
              )}
            </InfoBlock>
            <InfoBlock>
              <StrongLabel>Secret Ref Namespace:</StrongLabel>
              {!!externalAuth.clientSecretRef ? (
                <InfoBlock>{externalAuth.clientSecretRef.namespace}</InfoBlock>
              ) : (
                'None'
              )}
            </InfoBlock>
          </AuthInfo>
        ) : (
          <Unconfigured>
            External Authorization has not been configured.
          </Unconfigured>
        )}

        <UpdateButton
          type='link'
          onClick={e => setShowExtAuthModal(show => !show)}>
          Update Authorization.
        </UpdateButton>
        <SoloModal
          visible={showExtAuthModal}
          width={650}
          title='Update Authorization'
          onClose={() => setShowExtAuthModal(false)}>
          <React.Fragment>
            <Legend>
              Prior to creating an OAuth config, you must create a client
              secret.
              <br />
              Need help?{' '}
              <a
                href='https://gloo.solo.io/enterprise/authentication/auth/'
                target='_blank'>
                View Authorization documentation.
              </a>
            </Legend>

            <ExtAuthForm
              externalAuth={externalAuth}
              externalAuthChanged={externalAuthChanged}
            />
          </React.Fragment>
        </SoloModal>
      </div>
    </div>
  );
};
