import styled from '@emotion/styled/macro';
import { glooResourceApi } from 'API/gloo-resource';
import { useGetUpstreamDetails, useIsGlooFedEnabled } from 'API/hooks';
import { ReactComponent as UpstreamIcon } from 'assets/upstream-icon.svg';
import AreaHeader from 'Components/Common/AreaHeader';
import { DataError } from 'Components/Common/DataError';
import { HealthNotificationBox } from 'Components/Common/HealthNotificationBox';
import { Loading } from 'Components/Common/Loading';
import { SectionCard } from 'Components/Common/SectionCard';
import { UpstreamStatus } from 'proto/github.com/solo-io/solo-apis/api/gloo/gloo/v1/upstream_pb';
import React from 'react';
import { di } from 'react-magnetic-di/macro';
import { useParams } from 'react-router';
import { IconHolder } from 'Styles/StyledComponents/icons';
import { getUpstreamType } from 'utils/upstream-helpers';
import UpstreamConfiguration from './UpstreamConfiguration';
import UpstreamFailoverGroups from './UpstreamFailoverGroups';

const ConfigArea = styled.div`
  margin-bottom: 20px;
`;

export const UpstreamDetails = () => {
  di(useParams, useIsGlooFedEnabled, useGetUpstreamDetails);
  const {
    name = '',
    upstreamName = '',
    upstreamNamespace = '',
    upstreamClusterName = '',
  } = useParams();

  const { data: upResponse, error: upstreamError } = useGetUpstreamDetails({
    name: upstreamName,
    namespace: upstreamNamespace,
    clusterName: upstreamClusterName,
  });
  const upstream = upResponse?.upstream;

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

        {isGlooFedEnabled && <UpstreamFailoverGroups upstream={upstream} />}
      </>
    </SectionCard>
  );
};
