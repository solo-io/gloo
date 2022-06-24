import React from 'react';
import { useParams } from 'react-router';
import { colors } from 'Styles/colors';
import styled from '@emotion/styled';
import { SectionCard } from 'Components/Common/SectionCard';
import { ReactComponent as EnvoyLogo } from 'assets/envoy-logo.svg';
import { ReactComponent as DocumentsIcon } from 'assets/document.svg';
import { Loading } from 'Components/Common/Loading';
import { useGetConfigDumps } from 'API/hooks';
import { doDownload } from 'download-helper';
import { di } from 'react-magnetic-di/macro';
import YamlDisplayer from 'Components/Common/YamlDisplayer';
import { IconHolder } from 'Styles/StyledComponents/icons';
import { DataError } from 'Components/Common/DataError';
import { UpstreamStatus } from 'proto/github.com/solo-io/solo-apis/api/gloo/gloo/v1/upstream_pb';
import { ConfigDump } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/glooinstance_pb';

const TitleRow = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
  font-size: 18px;
`;

const Actionables = styled.div`
  display: flex;
  align-items: center;
  color: ${colors.seaBlue};

  > div {
    display: flex;
    align-items: center;
    margin-left: 20px;
    cursor: pointer;
  }

  svg {
    margin-right: 8px;
  }
`;

export const GlooAdminEnvoy = () => {
  di(useParams, useGetConfigDumps);
  const { name, namespace } = useParams();

  const { data: configsList, error: configsListError } = useGetConfigDumps({
    name: name!,
    namespace: namespace!,
  });

  const onDownloadConfig = (configDump: ConfigDump.AsObject) => {
    doDownload(configDump.raw, configDump.name + '.json');
  };

  if (!!configsListError) {
    return <DataError error={configsListError} />;
  } else if (!configsList) {
    return (
      <Loading message={`Retrievng configuration dumps for: ${name}...`} />
    );
  }

  return (
    <div>
      {configsList.map(configDump => {
        return (
          <SectionCard
            key={configDump.name}
            logoIcon={
              <IconHolder width={20}>
                <EnvoyLogo />
              </IconHolder>
            }
            cardName={configDump.name}
            health={{
              state:
                configDump.error.length > 0
                  ? UpstreamStatus.State.REJECTED
                  : UpstreamStatus.State.ACCEPTED,
              title: 'Config Status',
            }}>
            {configDump.error.length > 0 ? (
              <DataError error={{ message: configDump.error }} />
            ) : (
              <>
                <TitleRow>
                  <div>Raw Config (Read Only)</div>
                  <Actionables>
                    <div onClick={() => onDownloadConfig(configDump)}>
                      <DocumentsIcon /> {configDump.name}.json
                    </div>
                  </Actionables>
                </TitleRow>
                <YamlDisplayer contentString={configDump.raw} />
              </>
            )}
          </SectionCard>
        );
      })}
    </div>
  );
};
