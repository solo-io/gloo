import React from 'react';
import styled from '@emotion/styled/macro';
import { FailoverScheme } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/failover_scheme_pb';
import { FailoverSchemeStatus } from 'proto/github.com/solo-io/solo-projects/projects/gloo-fed/api/fed/v1/failover_pb';
import { failoverSchemeApi } from 'API/failover-scheme';
import { AreaTitle } from 'Styles/StyledComponents/headings';
import AreaHeader from 'Components/Common/AreaHeader';
import { CardSubsectionWrapper } from 'Components/Common/Card';
import { StatusType } from 'utils/health-status';
import UpstreamFailoverGroup from './UpstreamFailoverGroup';

const NoFailoverWrapper = styled(CardSubsectionWrapper)`
  text-align: center;
`;

const FailoverGroupsContainer = styled.div`
  > div:not(:first-of-type) {
    margin-top: 20px;
  }
`;

type Props = {
  failoverScheme?: FailoverScheme.AsObject;
};

const UpstreamFailoverGroups = ({ failoverScheme }: Props) => {
  if (!failoverScheme) {
    return (
      <>
        <AreaTitle>Failover Groups</AreaTitle>
        <NoFailoverWrapper>
          No failover groups have been configured.
        </NoFailoverWrapper>
      </>
    );
  }

  const name = failoverScheme?.metadata?.name;
  const namespace = failoverScheme?.metadata?.namespace;

  const loadYaml = async () => {
    if (!name || !namespace) {
      return '';
    }

    try {
      const yaml = await failoverSchemeApi.getFailoverSchemeYAML({
        name,
        namespace,
      });
      return yaml;
    } catch (error) {
      console.error(error);
    }
    return '';
  };

  return (
    <div>
      <AreaHeader
        title='Failover Groups'
        contentTitle={`${namespace}--${name}.yaml`}
        onLoadContent={loadYaml}
        health={{
          title: 'Failover Status',
          type: StatusType.FAILOVER,
          state:
            failoverScheme.status?.state ?? FailoverSchemeStatus.State.PENDING,
          reason: failoverScheme.status?.message,
        }}
      />
      <FailoverGroupsContainer>
        {failoverScheme.spec?.failoverGroupsList?.map((group, idx) => (
          <UpstreamFailoverGroup key={idx} priority={idx + 1} group={group} />
        ))}
      </FailoverGroupsContainer>
    </div>
  );
};

export default UpstreamFailoverGroups;
