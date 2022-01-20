import React from 'react';
import { useParams } from 'react-router';
import styled from '@emotion/styled/macro';
import { ReactComponent as UpstreamGroupIcon } from 'assets/upstream-group-icon.svg';
import { UpstreamStatus } from 'proto/github.com/solo-io/solo-apis/api/gloo/gloo/v1/upstream_pb';
import { glooResourceApi } from 'API/gloo-resource';
import { useGetUpstreamGroupDetails } from 'API/hooks';
import { SectionCard } from 'Components/Common/SectionCard';
import AreaHeader from 'Components/Common/AreaHeader';
import { HealthNotificationBox } from 'Components/Common/HealthNotificationBox';
import { IconHolder } from 'Styles/StyledComponents/icons';
import { Loading } from 'Components/Common/Loading';
import { UpstreamGroupWeightsTable } from './UpstreamGroupWeightsTable';
import { DataError } from 'Components/Common/DataError';

const ConfigArea = styled.div`
  margin-bottom: 20px;
`;

export const UpstreamGroupDetails = () => {
  const {
    name = '',
    namespace = '',
    upstreamGroupName = '',
    upstreamGroupNamespace = '',
    upstreamGroupClusterName = '',
  } = useParams();

  const { data: upstreamGroup, error: ugError } = useGetUpstreamGroupDetails(
    { name, namespace },
    {
      name: upstreamGroupName,
      namespace: upstreamGroupNamespace,
      clusterName: upstreamGroupClusterName,
    }
  );

  if (!!ugError) {
    return <DataError error={ugError} />;
  } else if (!upstreamGroup) {
    return (
      <Loading message={`Retrieving information for ${upstreamGroupName}...`} />
    );
  }

  const loadYaml = async () => {
    try {
      const yaml = await glooResourceApi.getUpstreamGroupYAML({
        name: upstreamGroupName,
        namespace: upstreamGroupNamespace,
        clusterName: upstreamGroupClusterName,
      });
      return yaml;
    } catch (error) {
      console.error(error);
    }
    return '';
  };

  return (
    <SectionCard
      cardName={upstreamGroupName}
      headerSecondaryInformation={[
        { title: 'namespace', value: upstreamGroupNamespace },
      ]}
      health={{
        title: 'Group Status',
        state: upstreamGroup?.status?.state ?? UpstreamStatus.State.PENDING,
        reason: upstreamGroup?.status?.reason,
      }}
      logoIcon={
        <IconHolder>
          <UpstreamGroupIcon />
        </IconHolder>
      }
    >
      {!!ugError ? (
        <DataError error={ugError} />
      ) : !upstreamGroup ? (
        <Loading
          message={`Retrieving upstream group: ${upstreamGroupName}...`}
        />
      ) : (
        <>
          <HealthNotificationBox
            state={upstreamGroup?.status?.state}
            reason={upstreamGroup?.status?.reason}
          />
          <ConfigArea>
            <AreaHeader
              title='Upstreams'
              contentTitle={`${upstreamGroupNamespace}--${upstreamGroupName}.yaml`}
              onLoadContent={loadYaml}
            />
            <UpstreamGroupWeightsTable
              destinations={upstreamGroup.spec?.destinationsList}
            />
          </ConfigArea>
        </>
      )}
    </SectionCard>
  );
};
