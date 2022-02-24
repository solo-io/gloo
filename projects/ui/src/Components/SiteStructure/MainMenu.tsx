import React from 'react';
import { colors } from 'Styles/colors';
import styled from '@emotion/styled';
import { NestedLink } from 'Components/Common/NestedLink';
import { EnableGraphqlFeature } from 'Components/Features/Graphql/EnableGraphqlFeature';
import { ReactComponent as GlooFedIcon } from 'assets/GlooFed-Specific/gloo-edge-logo-white-text.svg';
import { ReactComponent as GearIcon } from 'assets/gear-icon.svg';
import { ReactComponent as AdminGearHover } from 'assets/admin-settings-hover.svg';
import { SoloLink } from 'Components/Common/SoloLink';
import {
  useIsGlooFedEnabled,
  useListClusterDetails,
  useListGlooInstances,
} from 'API/hooks';
import { useLocation } from 'react-router';
import { Link } from 'react-router-dom';

const GlooIconHolder = styled(Link)`
  height: 35px;
  margin-right: 5px;

  svg {
    height: 35px;
  }
`;

const Container = styled.div`
  min-width: 1070px;
  height: 55px;
  line-height: 36px;
  background: ${colors.seaBlue};
`;
const InnerContainer = styled.div`
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  height: 100%;
  width: 1275px;
  max-width: 100vw;
  margin: 0 auto;
`;

const AppTitleBox = styled.div`
  display: flex;
  text-decoration: none;
  font-weight: 500;
  color: ${colors.puddleBlue};
  font-size: 18px;
  margin-right: 50px;
  margin-bottom: 8px;
  padding-right: 50px;
  border-right: 1px solid ${colors.lakeBlue};
  cursor: default;
`;

const Flexor = styled.div`
  display: flex;
  align-items: center;
`;

const GearIconHolder = styled.div`
  display: flex;
  align-items: center;

  svg {
    width: 28px;
  }
`;

export const MainMenu = () => {
  const routerLocation = useLocation();
  const { data: glooFedCheckResponse, error: glooFedCheckError } =
    useIsGlooFedEnabled();

  const isGlooFedEnabled = glooFedCheckResponse?.enabled;
  const { data: glooInstances, error: glooError } = useListGlooInstances();
  const { data: clusterDetailsList, error: cError } = useListClusterDetails();

  const multipleClustersOrInstances =
    (clusterDetailsList && clusterDetailsList.length > 1) ||
    (glooInstances && glooInstances.length > 1);

  return (
    <Container>
      <InnerContainer>
        <Flexor>
          <AppTitleBox>
            <GlooIconHolder to='/'>
              <GlooFedIcon />
            </GlooIconHolder>
          </AppTitleBox>
          <NestedLink to='/' exact={true}>
            Overview
          </NestedLink>
          {multipleClustersOrInstances && (
            <NestedLink exact to='/gloo-instances/'>
              Gloo Instances
            </NestedLink>
          )}
          <NestedLink to='/virtual-services/' exact>
            Virtual Services
          </NestedLink>
          <NestedLink to='/upstreams/' exact>
            Upstreams
          </NestedLink>
          <NestedLink to='/wasm-filters/' exact>
            Wasm
          </NestedLink>
          <EnableGraphqlFeature>
            <NestedLink
              to={
                clusterDetailsList?.length === 1 && glooInstances?.length === 1
                  ? `/gloo-instances/${
                      clusterDetailsList[0]!.glooInstancesList[0].metadata
                        ?.namespace
                    }/${
                      clusterDetailsList[0]!.glooInstancesList[0].metadata?.name
                    }/apis/`
                  : '/apis/'
              }
              exact>
              APIs
            </NestedLink>
          </EnableGraphqlFeature>
        </Flexor>
        <Flexor style={{ height: '100%' }}>
          {/*<HelpHolder>
            <NewsletterLink
                      href="https://share.hsforms.com/1gyYyUxhWTOaoDF_LCNrXkw31z0a"
                      target="_blank"
                    >
                      Get Updates
                    </NewsletterLink>
                  </div>
                  <VersionDisplay>Version: unknown</VersionDisplay>
                </div>
              }
            >
              <HelpBubble />
            </Popover>
            </HelpHolder>*/}
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
        </Flexor>
      </InnerContainer>
    </Container>
  );
};
