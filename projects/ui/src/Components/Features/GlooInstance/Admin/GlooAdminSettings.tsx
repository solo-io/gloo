import React, { useEffect } from 'react';
import { useParams } from 'react-router';
import { colors } from 'Styles/colors';
import styled from '@emotion/styled';
import { SectionCard } from 'Components/Common/SectionCard';
import { ReactComponent as GearIcon } from 'assets/gear-icon.svg';
import { ReactComponent as DocumentsIcon } from 'assets/document.svg';
import { Loading } from 'Components/Common/Loading';
import { useListSettings } from 'API/hooks';
import { glooResourceApi } from 'API/gloo-resource';
import { doDownload } from 'download-helper';
import YamlDisplayer from 'Components/Common/YamlDisplayer';
import { Settings } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/gloo_resources_pb';
import { IconHolder } from 'Styles/StyledComponents/icons';
import { DataError } from 'Components/Common/DataError';

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

export const GlooAdminSettings = () => {
  const { name = '', namespace = '' } = useParams();

  const { data: settings, error: sError } = useListSettings({
    name,
    namespace,
  });

  const [yamlsOpen, setYamlsOpen] = React.useState<{
    [key: string]: boolean; // setting uid
  }>({});
  useEffect(() => {
    const yamlsListBySettingUid: { [key: string]: boolean } = {};
    if (settings?.length) {
      settings
        .filter(setting => !!setting.metadata)
        .forEach(setting => {
          yamlsListBySettingUid[setting.metadata!.uid] = false;
        });
    }

    setYamlsOpen(yamlsListBySettingUid);

    // expand yaml for first one by default
    if (settings?.length) {
      toggleView(settings[0]);
    }
    /* eslint-disable-next-line react-hooks/exhaustive-deps */
  }, [settings]);

  const [swaggerContentByUid, setSwaggerContentByUid] = React.useState<{
    [key: string]: string;
  }>({});

  if (!!sError) {
    return <DataError error={sError} />;
  } else if (!settings) {
    return <Loading message={`Retrieving settings for ${name}...`} />;
  }

  const toggleView = (setting: Settings.AsObject) => {
    let viewables = { ...yamlsOpen };
    viewables[setting.metadata!.uid] = !viewables[setting.metadata!.uid];
    setYamlsOpen(viewables);

    if (setting.metadata && !swaggerContentByUid[setting.metadata!.uid]) {
      glooResourceApi
        .getSettingYAML({
          name: setting.metadata.name,
          namespace: setting.metadata.namespace,
          clusterName: setting.metadata.clusterName,
        })
        .then(settingYaml => {
          let swaggers = { ...swaggerContentByUid };
          swaggers[setting.metadata!.uid] = settingYaml;
          setSwaggerContentByUid(swaggers);
        });
    }
  };

  const onDownloadSetting = (setting: Settings.AsObject) => {
    if (setting.metadata && !swaggerContentByUid[setting.metadata.uid]) {
      glooResourceApi
        .getSettingYAML({
          name: setting.metadata.name,
          namespace: setting.metadata.namespace,
          clusterName: setting.metadata.clusterName,
        })
        .then(settingYaml => {
          doDownload(settingYaml, setting.metadata?.name + '.yaml');

          let swaggers = { ...swaggerContentByUid };
          swaggers[setting.metadata!.uid] = settingYaml;
          setSwaggerContentByUid(swaggers);
        });
    } else {
      doDownload(
        swaggerContentByUid[setting.metadata!.uid],
        setting.metadata?.name + '.yaml'
      );
    }
  };

  return (
    <div>
      {settings?.map(setting => {
        let secondaryHeaderInfo = [
          {
            title: 'Namespace',
            value: setting.metadata?.namespace,
          },
        ];
        if (setting.spec?.rbac?.requireRbac !== undefined) {
          secondaryHeaderInfo.push({
            title: 'RBAC',
            value: setting.spec.rbac.requireRbac ? 'Require' : 'False',
          });
        }

        return (
          <SectionCard
            key={
              setting.metadata?.uid ??
              'Setting supplied with no unique identifier'
            }
            logoIcon={
              <IconHolder
                width={20}
                applyColor={{ strokeNotFill: true, color: colors.seaBlue }}
              >
                <GearIcon />
              </IconHolder>
            }
            cardName={setting.metadata!.name}
            headerSecondaryInformation={secondaryHeaderInfo}
          >
            <TitleRow>
              <div>Settings (Read Only)</div>
              <Actionables>
                <div onClick={() => toggleView(setting)}>
                  {yamlsOpen[setting.metadata!.uid] ? 'Hide' : 'View'} Raw
                  Config
                </div>
                <div onClick={() => onDownloadSetting(setting)}>
                  <DocumentsIcon /> {setting.metadata?.name}.yaml
                </div>
              </Actionables>
            </TitleRow>
            {yamlsOpen[setting.metadata!.uid] &&
              (!!swaggerContentByUid[setting.metadata!.uid] ? (
                <YamlDisplayer
                  contentString={swaggerContentByUid[setting.metadata!.uid]}
                />
              ) : (
                <Loading message={'Retrieving configuration...'} />
              ))}
          </SectionCard>
        );
      })}
    </div>
  );
};
