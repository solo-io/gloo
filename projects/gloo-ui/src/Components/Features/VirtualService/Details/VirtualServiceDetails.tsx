import styled from '@emotion/styled';
import { Spin } from 'antd';
import { ReactComponent as GlooIcon } from 'assets/GlooEE.svg';
import { Breadcrumb } from 'Components/Common/Breadcrumb';
import { ConfigDisplayer } from 'Components/Common/DisplayOnly/ConfigDisplayer';
import { FileDownloadLink } from 'Components/Common/FileDownloadLink';
import { SectionCard } from 'Components/Common/SectionCard';
import * as React from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { useHistory, useParams } from 'react-router';
import { AppState } from 'store';
import { updateVirtualServiceYaml } from 'store/virtualServices/actions';
import { colors, healthConstants, soloConstants } from 'Styles';
import { Domains } from './Domains';
import { ExtAuth } from './ExtAuth';
import { RateLimit } from './RateLimit';
import { Routes } from './Routes';
import { RouteParent } from '../RouteTableDetails';

export const ConfigContainer = styled.div`
  display: grid;
  grid-template-columns: 1fr 1fr;
  background: ${colors.januaryGrey};
  height: 80%;
  border-radius: ${soloConstants.smallRadius}px;
`;

export const ConfigItem = styled.div`
  margin: 20px;
  padding: 10px;
  justify-items: center;
  background: white;
`;

type DetailsContentProps = { configurationShowing?: boolean };
export const DetailsContent = styled.div`
  position: relative;
  display: grid;
  grid-template-rows: ${(props: DetailsContentProps) =>
      props.configurationShowing ? 'auto' : ''} auto 1fr 1fr;
  grid-template-columns: 100%;
  grid-column-gap: 30px;
`;

export const YamlLink = styled.div`
  position: absolute;
  top: 10px;
  right: 0;
  display: flex;
`;

export const ConfigurationToggle = styled.div`
  cursor: pointer;
  color: ${colors.seaBlue};
  font-size: 14px;
  margin-right: 8px;
`;

export const DetailsSection = styled.div`
  width: 100%;
`;

export const DetailsSectionTitle = styled.div`
  font-size: 18px;
  font-weight: bold;
  color: ${colors.novemberGrey};
  margin-top: 10px;
  margin-bottom: 10px;
`;

export const VirtualServiceDetails = () => {
  let history = useHistory();
  let { virtualservicename, virtualservicenamespace } = useParams();

  const [showConfiguration, setShowConfiguration] = React.useState(false);

  const virtualServicesList = useSelector(
    (state: AppState) => state.virtualServices.virtualServicesList
  );

  const yamlError = useSelector(
    (state: AppState) => state.virtualServices.yamlParseError
  );

  const dispatch = useDispatch();

  let virtualServiceDetails = virtualServicesList.find(
    vsD => vsD?.virtualService?.metadata?.name === virtualservicename
  )!;

  if (!virtualServiceDetails?.virtualService?.virtualHost) {
    return (
      <>
        <Breadcrumb />
        <Spin size='large' />
      </>
    );
  }

  let { virtualService, raw, plugins } = virtualServiceDetails!;

  let {
    routesList,
    domainsList
  } = virtualServiceDetails?.virtualService?.virtualHost;

  let rateLimits;
  let externalAuth;

  if (plugins) {
    rateLimits = plugins.rateLimit;
    externalAuth = plugins.extAuth;
  }

  const saveYamlChange = (newYaml: string) => {
    dispatch(
      updateVirtualServiceYaml({
        editedYamlData: {
          editedYaml: newYaml,
          ref: {
            name: virtualService!.metadata!.name,
            namespace: virtualService!.metadata!.namespace
          }
        }
      })
    );
  };

  const headerInfo = [
    {
      title: 'namespace',
      value: virtualservicenamespace!
    }
  ];
  return (
    <>
      <Breadcrumb />

      <SectionCard
        cardName={
          virtualService!.displayName.length
            ? virtualService!.displayName
            : virtualservicename
            ? virtualservicename
            : 'Error'
        }
        logoIcon={<GlooIcon />}
        health={
          virtualService!.status
            ? virtualService!.status!.state
            : healthConstants.Pending.value
        }
        headerSecondaryInformation={headerInfo}
        healthMessage={
          virtualService!.status && virtualService!.status!.reason.length
            ? virtualService!.status!.reason
            : 'Service Status'
        }
        onClose={() => history.push(`/virtualservices/`)}>
        <DetailsContent configurationShowing={showConfiguration}>
          {!!raw && (
            <YamlLink>
              <ConfigurationToggle
                onClick={() => setShowConfiguration(s => !s)}>
                {showConfiguration ? 'Hide' : 'View'} YAML Configuration
              </ConfigurationToggle>
              <FileDownloadLink
                fileContent={raw.content}
                fileName={raw.fileName}
              />
            </YamlLink>
          )}
          {showConfiguration && (
            <DetailsSection>
              <DetailsSectionTitle>YAML Configuration</DetailsSectionTitle>
              <ConfigDisplayer
                content={raw ? raw.content : ''}
                asEditor
                yamlError={yamlError}
                saveEdits={saveYamlChange}
              />
            </DetailsSection>
          )}
          <DetailsSection>
            <Domains
              domains={domainsList}
              vsRef={{
                name: virtualservicename!,
                namespace: virtualservicenamespace!
              }}
            />
          </DetailsSection>
          <DetailsSection>
            <Routes
              routes={routesList}
              routeParent={virtualService as RouteParent}
            />
          </DetailsSection>
          <DetailsSection>
            <>
              <DetailsSectionTitle>Configuration</DetailsSectionTitle>
              <ConfigContainer>
                <ConfigItem>
                  <ExtAuth externalAuth={externalAuth} />
                </ConfigItem>
                <ConfigItem>
                  <RateLimit rateLimits={rateLimits} />
                </ConfigItem>
              </ConfigContainer>
            </>
          </DetailsSection>
        </DetailsContent>
      </SectionCard>
    </>
  );
};
