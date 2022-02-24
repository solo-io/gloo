import React, { useEffect, useState } from 'react';
import styled from '@emotion/styled/macro';
import { TabPanels, Tabs } from '@reach/tabs';
import { Loading } from 'Components/Common/Loading';
import { useParams } from 'react-router';
import { GlooInstance } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/glooinstance_pb';
import { SectionCard } from 'Components/Common/SectionCard';
import { ReactComponent as GlooIcon } from 'assets/Gloo.svg';
import { ReactComponent as MeshIcon } from 'assets/mesh-icon.svg';
import { ReactComponent as ClusterIcon } from 'assets/cluster-icon.svg';
import { ReactComponent as NamespaceIcon } from 'assets/namespace-icon.svg';
import { ReactComponent as VersionsIcon } from 'assets/versions-icon.svg';
import { ReactComponent as GearIcon } from 'assets/gear-icon.svg';
import {
  FolderTab,
  FolderTabList,
  StyledTabPanel,
  FolderTabContent,
} from 'Components/Common/Tabs';
import { GlooInstanceUpstreams } from './GlooInstanceUpstreams';
import { GlooInstanceVirtualServices } from './GlooInstanceVirtualServices';
import { GlooInstanceUpstreamGroups } from './GlooInstanceUpstreamGroups';
import { CardWhiteSubsection } from 'Components/Common/Card';
import { colors } from 'Styles/colors';
import { SoloLink } from 'Components/Common/SoloLink';
import { useListGlooInstances } from 'API/hooks';
import { GlooInstanceIssues } from './GlooInstanceIssues';
import { DataError } from 'Components/Common/DataError';
import { GlooInstanceRouteTables } from './GlooInstanceRouteTables';
import { formatTimestamp } from 'utils';

const GlooIconHolder = styled.div`
  display: flex;
  align-items: center;
  justify-items: center;

  svg {
    width: 25px;
  }
`;

const QuickStats = styled(CardWhiteSubsection)`
  display: flex;
  align-items: center;
  margin-bottom: 28px;

  > div {
    display: flex;
    align-items: center;
    margin-right: 45px;
  }

  svg {
    height: 30px;
    margin-right: 10px;
  }
`;
const QuickStatTitle = styled.div`
  font-weight: 600;
  margin-right: 4px;
`;

const Divider = styled.div`
  width: 1px;
  height: 42px;
  background: ${colors.marchGrey};
`;

const OrangeIconHolder = styled.div`
  svg * {
    fill: ${colors.grapefruitOrange};
  }
`;
const BlueIconHolder = styled.div`
  svg * {
    fill: ${colors.seaBlue};
  }
`;

const TabsContainer = styled.div`
  position: relative;
`;

const AdminLinkHolder = styled.div`
  position: absolute;
  right: 0px;
  top: 14px;

  > div {
    svg {
      width: 24px;
      margin-left: 8px;

      * {
        fill: ${colors.seaBlue};
        stroke: ${colors.seaBlue};
      }
    }
  }
`;

const AdminLink = styled.div`
  display: flex;
  align-items: center;
`;

/*const tabsMap: { [key: string]: number } = {
  "virtual-services": 0,
  upstreams: 1,
  "upstream-groups": 2,
};*/

export const GlooInstancesDetails = () => {
  const { name, namespace } = useParams();

  const { data: glooInstances, error: instancesError } = useListGlooInstances();

  const [glooInstance, setGlooInstance] = useState<GlooInstance.AsObject>();
  const [tabIndex, setTabIndex] = React.useState(0); // 'virtual-services'

  useEffect(() => {
    if (!!glooInstances) {
      setGlooInstance(
        glooInstances.find(
          instance =>
            instance.metadata?.name === name &&
            instance.metadata?.namespace === namespace
        )
      );
    } else {
      setGlooInstance(undefined);
    }
  }, [name, namespace, glooInstances]);

  if (!!instancesError) {
    return <DataError error={instancesError} />;
  } else if (!glooInstances) {
    return (
      <Loading message={`Retrieving information on instance: ${name}...`} />
    );
  }

  const handleTabsChange = (index: number) => {
    setTabIndex(index);
  };

  const secondaryHeaderInfo = [
    {
      title: 'Namespace',
      value: namespace,
    },
  ];
  if (glooInstance?.metadata?.creationTimestamp) {
    secondaryHeaderInfo.unshift({
      title: 'Last Updated',
      value: formatTimestamp(glooInstance?.metadata?.creationTimestamp),
    });
  }

  return (
    <SectionCard
      cardName={name!}
      logoIcon={
        <GlooIconHolder>
          <GlooIcon />
        </GlooIconHolder>
      }
      headerSecondaryInformation={secondaryHeaderInfo}>
      {!!instancesError ? (
        <DataError error={instancesError} />
      ) : !glooInstance ? (
        <Loading
          message={`Retrieving information for the instance ${name}...`}
        />
      ) : (
        <div>
          <GlooInstanceIssues glooInstance={glooInstance} />
          <QuickStats>
            <div>
              <OrangeIconHolder>
                <MeshIcon />
              </OrangeIconHolder>
              <QuickStatTitle>Region: </QuickStatTitle>{' '}
              {glooInstance.spec?.region}
            </div>
            <div>
              <QuickStatTitle>Zones:</QuickStatTitle>{' '}
              {glooInstance.spec?.proxiesList
                .map(prox => prox.zonesList.join(', '))
                .join(', ')}
            </div>

            <Divider />

            <div>
              <BlueIconHolder>
                <ClusterIcon />
              </BlueIconHolder>
              <QuickStatTitle>Cluster:</QuickStatTitle>{' '}
              {glooInstance.spec?.cluster}
            </div>
            <div>
              <BlueIconHolder>
                <NamespaceIcon />
              </BlueIconHolder>
              <QuickStatTitle>Namespace:</QuickStatTitle>{' '}
              {glooInstance.metadata?.namespace}
            </div>
            <div>
              <BlueIconHolder>
                <VersionsIcon />
              </BlueIconHolder>
              <QuickStatTitle>Version:</QuickStatTitle>{' '}
              {glooInstance.spec?.controlPlane?.version}
            </div>
          </QuickStats>

          <TabsContainer>
            <AdminLinkHolder>
              <SoloLink
                link={`gloo-admin`}
                displayElement={
                  <AdminLink>
                    <div>Gloo Admin Settings</div> <GearIcon />
                  </AdminLink>
                }
              />
            </AdminLinkHolder>
            <Tabs index={tabIndex} onChange={handleTabsChange}>
              <FolderTabList>
                <FolderTab>Virtual Services</FolderTab>
                <FolderTab>Route Tables</FolderTab>
                <FolderTab>Upstreams</FolderTab>
                <FolderTab>Upstream Groups</FolderTab>
              </FolderTabList>

              <TabPanels>
                <StyledTabPanel>
                  <FolderTabContent>
                    <GlooInstanceVirtualServices />
                  </FolderTabContent>
                </StyledTabPanel>
                <StyledTabPanel>
                  <FolderTabContent>
                    <GlooInstanceRouteTables />
                  </FolderTabContent>
                </StyledTabPanel>
                <StyledTabPanel>
                  <FolderTabContent>
                    <GlooInstanceUpstreams />
                  </FolderTabContent>
                </StyledTabPanel>
                <StyledTabPanel>
                  <FolderTabContent>
                    <GlooInstanceUpstreamGroups />
                  </FolderTabContent>
                </StyledTabPanel>
              </TabPanels>
            </Tabs>
          </TabsContainer>
        </div>
      )}
    </SectionCard>
  );
};
