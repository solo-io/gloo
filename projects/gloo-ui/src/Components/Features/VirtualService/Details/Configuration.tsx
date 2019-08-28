import styled from '@emotion/styled';
import { OAuth } from 'proto/github.com/solo-io/solo-projects/projects/gloo/api/v1/plugins/extauth/extauth_pb';
import { IngressRateLimit } from 'proto/github.com/solo-io/solo-projects/projects/gloo/api/v1/plugins/ratelimit/ratelimit_pb';
import * as React from 'react';
import { colors, soloConstants } from 'Styles';
import { ExtAuth } from './ExtAuth';
import { RateLimit } from './RateLimit';
import { DetailsSectionTitle } from './VirtualServiceDetails';

const ConfigContainer = styled.div`
  display: grid;
  grid-template-columns: 1fr 1fr;
  background: ${colors.januaryGrey};
  height: 80%;
  border-radius: ${soloConstants.smallRadius}px;
`;

const ConfigItem = styled.div`
  margin: 20px;
  padding: 10px;
  justify-items: center;
  background: white;
`;

interface Props {
  rates: IngressRateLimit.AsObject | undefined;
  rateLimitsChanged: (newRateLimits: IngressRateLimit.AsObject) => any;
  externalAuth: OAuth.AsObject | undefined;
  externalAuthChanged: (newExternalAuth: OAuth.AsObject) => any;
}
export const Configuration = (props: Props) => {
  const { rates, rateLimitsChanged, externalAuth, externalAuthChanged } = props;

  return (
    <React.Fragment>
      <DetailsSectionTitle>Configuration</DetailsSectionTitle>
      <ConfigContainer>
        <ConfigItem>
          <ExtAuth
            externalAuth={externalAuth}
            externalAuthChanged={externalAuthChanged}
          />
        </ConfigItem>
        <ConfigItem>
          <RateLimit rates={rates} rateLimitsChanged={rateLimitsChanged} />
        </ConfigItem>
      </ConfigContainer>
    </React.Fragment>
  );
};
