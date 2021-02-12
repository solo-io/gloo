import React from 'react';
import styled from '@emotion/styled/macro';
import { Label } from 'Components/Common/SoloInput';
import { colors } from 'Styles/colors';
import { ExtAuthExtension } from 'proto/github.com/solo-io/solo-apis/api/gloo/enterprise.gloo/v1/auth_config_pb';

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
  externalAuth?: ExtAuthExtension.AsObject;
}

export const ExternalAuthorizationSection = ({ externalAuth }: Props) => {
  return (
    <div>
      <ConfigItemHeader>External Authorization </ConfigItemHeader>
      <div>
        {externalAuth?.configRef?.name ? (
          <AuthInfo>
            <InfoBlock>
              <StrongLabel>Auth ConfigRef Name:</StrongLabel>
              <InfoBlock>{externalAuth?.configRef?.name}</InfoBlock>
            </InfoBlock>
          </AuthInfo>
        ) : (
          <Unconfigured>
            External Authorization has not been configured.
          </Unconfigured>
        )}
      </div>
    </div>
  );
};
