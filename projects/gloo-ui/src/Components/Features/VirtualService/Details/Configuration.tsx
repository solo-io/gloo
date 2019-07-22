import * as React from 'react';
import styled from '@emotion/styled/macro';
import { colors } from 'Styles';
import { DetailsSectionTitle } from './VirtualServiceDetails';
import { ExtAuth } from './ExtAuth';
import { RateLimit } from './RateLimit';
import { IngressRateLimit } from 'proto/github.com/solo-io/solo-projects/projects/gloo/api/v1/plugins/ratelimit/ratelimit_pb';

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

interface Props {
  rates: IngressRateLimit.AsObject | undefined;
  rateLimitsChanged: (newRateLimits: IngressRateLimit.AsObject) => any;
}
export const Configuration = (props: Props) => {
  const { rates, rateLimitsChanged } = props;

  return (
    <React.Fragment>
      <DetailsSectionTitle>Configuration</DetailsSectionTitle>
      <ConfigContainer>
        <ConfigItem>
          <ExtAuth />
        </ConfigItem>
        <ConfigItem>
          <RateLimit rates={rates} rateLimitsChanged={rateLimitsChanged} />
        </ConfigItem>
      </ConfigContainer>
    </React.Fragment>
  );
};
