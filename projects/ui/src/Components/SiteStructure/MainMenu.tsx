import React, { useMemo } from 'react';
import styled from '@emotion/styled';
import { ReactComponent as GlooFedIcon } from 'assets/GlooFed-Specific/gloo-edge-logo-white-text.svg';
import {
  useIsGraphqlEnabled,
  useListClusterDetails,
  useListGlooInstances,
} from 'API/hooks';
import SoloNavbar, { ISoloNavLink } from 'Components/Common/SoloNavbar';
import { SoloLink } from 'Components/Common/SoloLink';
import { useLocation } from 'react-router';
import { ReactComponent as GearIcon } from 'assets/gear-icon.svg';
import { ReactComponent as AdminGearHover } from 'assets/admin-settings-hover.svg';

const GearIconHolder = styled.div`
  display: flex;
  align-items: center;

  svg {
    width: 28px;
  }
`;

export const MainMenu = () => {
  const routerLocation = useLocation();
  const { data: graphqlIntegrationEnabled, error: graphqlCheckError } =
    useIsGraphqlEnabled();
  const isGraphQLEnabled = !graphqlCheckError && graphqlIntegrationEnabled;

  const { data: glooInstances, error: glooError } = useListGlooInstances();
  const { data: clusterDetailsList, error: cError } = useListClusterDetails();

  const multipleClustersOrInstances =
    (clusterDetailsList && clusterDetailsList.length > 1) ||
    (glooInstances && glooInstances.length > 1);

  const navLinks = useMemo(() => {
    const links = [] as ISoloNavLink[];
    links.push({ name: 'Overview', href: '/', exact: true });
    if (multipleClustersOrInstances)
      links.push({ name: 'Gloo Instances', href: '/gloo-instances/' });
    links.push(
      { name: 'Virtual Services', href: '/virtual-services/' },
      { name: 'Upstreams', href: '/upstreams/' },
      { name: 'Wasm', href: '/wasm-filters/' }
    );
    const apiRoute =
      glooInstances?.length === 1
        ? `/gloo-instances/${glooInstances[0].metadata?.namespace}/${glooInstances[0].metadata?.name}/apis/`
        : '/apis/';
    if (isGraphQLEnabled) links.push({ name: 'APIs', href: apiRoute });
    return links;
  }, [multipleClustersOrInstances, isGraphQLEnabled]);

  return (
    <SoloNavbar
      BrandComponent={GlooFedIcon}
      navLinks={navLinks}
      SettingsComponent={() => (
        <SoloLink
          link={
            clusterDetailsList?.length === 1 && glooInstances?.length === 1
              ? `gloo-instances/${
                  clusterDetailsList[0]!.glooInstancesList[0].metadata
                    ?.namespace
                }/${
                  clusterDetailsList[0]!.glooInstancesList[0].metadata?.name
                }/gloo-admin/`
              : '/admin/'
          }
          displayElement={
            <GearIconHolder>
              {routerLocation.pathname.includes(
                multipleClustersOrInstances ? '/admin/' : '/gloo-admin/'
              ) ? (
                <AdminGearHover />
              ) : (
                <GearIcon />
              )}
            </GearIconHolder>
          }
        />
      )}
    />
  );
};
