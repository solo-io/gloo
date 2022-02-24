import React from 'react';
import { useParams } from 'react-router';
import styled from '@emotion/styled/macro';
import { ReactComponent as UpstreamIcon } from 'assets/upstream-icon.svg';
import { UpstreamStatus } from 'proto/github.com/solo-io/solo-apis/api/gloo/gloo/v1/upstream_pb';
import { glooResourceApi } from 'API/gloo-resource';
import { useGetUpstreamDetails, useIsGlooFedEnabled } from 'API/hooks';
import { getUpstreamType } from 'utils/upstream-helpers';
import { SectionCard } from 'Components/Common/SectionCard';
import AreaHeader from 'Components/Common/AreaHeader';
import { HealthNotificationBox } from 'Components/Common/HealthNotificationBox';
import { IconHolder } from 'Styles/StyledComponents/icons';
import UpstreamConfiguration from './UpstreamConfiguration';
import UpstreamFailoverGroups from './UpstreamFailoverGroups';
import { DataError } from 'Components/Common/DataError';
import { Loading } from 'Components/Common/Loading';

const ConfigArea = styled.div`
  margin-bottom: 20px;
`;

export const UpstreamDetails = () => {
  const {
    name = '',
    namespace = '',
    upstreamName = '',
    upstreamNamespace = '',
    upstreamClusterName = '',
  } = useParams();

  const { data: upstream, error: upstreamError } = useGetUpstreamDetails(
    { name, namespace },
    {
      name: upstreamName,
      namespace: upstreamNamespace,
      clusterName: upstreamClusterName,
    }
  );

  const { data: glooFedCheckResponse, error: glooFedCheckError } =
    useIsGlooFedEnabled();
  const isGlooFedEnabled = glooFedCheckResponse?.enabled;

  if (!!upstreamError) {
    return <DataError error={upstreamError} />;
  } else if (!upstream) {
    return <Loading message={`Retrieving upstream for ${name}...`} />;
  }

  const loadYaml = async () => {
    if (!upstreamName || !upstreamNamespace) {
      return '';
    }

    try {
      const yaml = await glooResourceApi.getUpstreamYAML({
        name: upstreamName,
        namespace: upstreamNamespace,
        clusterName: upstreamClusterName,
      });
      return yaml;
    } catch (error) {
      console.error(error);
    }
    return '';
  };

  return (
    <SectionCard
      cardName={upstreamName}
      headerSecondaryInformation={[
        { title: 'name', value: upstreamName },
        { title: 'namespace', value: upstreamNamespace },
        { title: 'type', value: getUpstreamType(upstream) },
      ]}
      health={{
        title: 'Upstream Status',
        state: upstream?.status?.state ?? UpstreamStatus.State.PENDING,
        reason: upstream?.status?.reason,
      }}
      logoIcon={
        <IconHolder>
          <UpstreamIcon />
        </IconHolder>
      }>
      <>
        <HealthNotificationBox
          state={upstream?.status?.state}
          reason={upstream?.status?.reason}
        />
        <ConfigArea>
          <AreaHeader
            title='Configuration'
            contentTitle={`${upstreamNamespace}--${upstreamName}.yaml`}
            onLoadContent={loadYaml}
          />
          <UpstreamConfiguration upstream={upstream} />
        </ConfigArea>

        {isGlooFedEnabled && (
          <UpstreamFailoverGroups
            upstreamName={upstreamName}
            upstreamNamespace={upstreamNamespace}
            upstreamClusterName={upstreamClusterName}
          />
        )}
      </>
    </SectionCard>
  );
};
