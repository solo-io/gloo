import React, { useState, useEffect } from 'react';
import styled from '@emotion/styled';
import { OverviewGlooInstancesBox, OverviewSmallBoxSummary } from '../Overview/OverviewBoxSummary';
import { ReactComponent as VirtualServiceIcon } from 'assets/virtualservice-icon.svg';
import { ReactComponent as UpstreamIcon } from 'assets/upstream-icon.svg';
import { ReactComponent as ClusterIcon } from 'assets/cluster-icon.svg';
import { ReactComponent as SuccessCircle } from 'assets/big-successful-checkmark.svg';
import { ReactComponent as WarningExclamation } from 'assets/big-warning-exclamation.svg';
import { GatewayStatus } from 'proto/github.com/solo-io/solo-apis/api/gloo/gateway/v1/gateway_pb';
import { ReactComponent as GatewayIcon } from 'assets/gateway.svg';
import { ReactComponent as ProxyIcon } from 'assets/proxy-icon.svg';
import { ReactComponent as GearIcon } from 'assets/gear-icon.svg';
import { ReactComponent as EnvoyIcon } from 'assets/envoy-logo.svg';
import { ReactComponent as WatchedNamespacesIcon } from 'assets/watched-namespace-icon.svg';
import { ReactComponent as SecretsIcon } from 'assets/cloud-key-icon.svg';
import { ReactComponent as HealthIcon } from 'assets/health-icon.svg';
import { Card, CardHeader, CardSubsectionContent, CardSubsectionWrapper } from 'Components/Common/Card';
import { colors } from 'Styles/colors';
import { SoloLink } from 'Components/Common/SoloLink';

const Container = styled.div``

const Main = styled.div`
    display: flex;
    flex-direction: row;
    flex-wrap: nowrap;
    gap: 10px;
    justify-content: space-between;
    align-items: center;
`;

const OverviewGlooInstancesBoxOverride = styled(OverviewGlooInstancesBox)`
    flex-direction: column;
    max-width: 300px;
`;

const GraphqlHeaderTitles = styled.div`
    display: flex;
    flex-direction: column;
`;

const GraphqlHeaderWrapper = styled.header`
    display: flex;
    flex-direction: row;
    justify-content: space-between;
    align-items: center;
`;

const BottomWrapper = styled.div`
    width: 400px;
`;

const StyledCardSubsectionWrapper = styled(CardSubsectionWrapper)`
    flex-basis: 100%;
`;

const LogoHolder = styled.div`
  svg {
    height: 28px;
  }
`;
const LogoRecolorHolder = styled(LogoHolder)`
  svg {
    * {
      fill: ${colors.seaBlue};
    }
  }
`;

const LogoRecolorAndResizeHolder = styled(LogoHolder)`
  svg {
    * {
      fill: ${colors.seaBlue};
      stroke-width: .2px;
    }
  }
`;

const Success = styled.div`
  border-radius: 100%;
  background: ${colors.forestGreen};
  height: 20px;
  width: 20px;
`;

const ItemWrapper = styled.div`
    display: flex;
    flex-direction: row;
    justify-content: space-between;
    width: 120px;
    padding-left: 20px;
    flex: 1;
    align-items: center;
    height: 40px;
`;

const StyledCard = styled(Card)`
  height: 405px;
`;

const BottomFooterWrapper = styled.footer`
  display: flex;
  flex-direction: row;
  align-items: center;
`;

const BottomContent = styled.div`
  padding-top: 20px;
  padding-bottom: 20px;
`;

const GearWrapper = styled.div`
  margin-left: 20px;
`;

/**
 * TODO:  This is a mock implementation for now.
 */
// GATEWAY
export const GraphqlOverview = () => {
    return (
        <Container>
            <CardSubsectionContent>
                <GraphqlHeaderWrapper>
                    <GraphqlHeaderTitles>
                        <h2>GraphQL Administration</h2>
                        <h3>Advanced Administration</h3>
                    </GraphqlHeaderTitles>
                    <div>
                        <HealthIcon />
                    </div>
                </GraphqlHeaderWrapper>
                <Main>
                    <StyledCardSubsectionWrapper>
                        <StyledCard>
                            <OverviewSmallBoxSummary
                                css={{
                                    justifyContent: 'space-between',
                                }}
                                title='Gateway Configuration'
                                logo={<ItemWrapper>
                                    <UpstreamIcon />
                                    <Success />
                                </ItemWrapper>}
                                description='Gateways are used to configure the protocols and ports for Envoy. Optionally, gateways can be associated with a specific set of virtual services.'
                                status={GatewayStatus.State.ACCEPTED}
                                count={2}
                                countDescription='Gateway Configurations are configured within Gloo Edge'
                                link='gateways'
                                descriptionTitle='View Gateways'
                            />
                        </StyledCard>
                    </StyledCardSubsectionWrapper>
                    <StyledCardSubsectionWrapper>
                        <StyledCard>
                            <OverviewSmallBoxSummary
                                title={'Envoy Configuration'}
                                logo={
                                    <ItemWrapper>
                                        <LogoRecolorHolder>
                                            <ProxyIcon />
                                        </LogoRecolorHolder>
                                        <Success />

                                    </ItemWrapper>
                                }
                                description={
                                    'Gloo generates proxy configs from upstreams, virtual services, and gateways, and then transforms them directly into Envoy config. If a proxy config is rejected, it means Envoy will not receive configuration updates.'
                                }
                                status={GatewayStatus.State.ACCEPTED}
                                count={1}
                                countDescription='Proxy Configurations are configured within Gloo Edge'
                                link='proxy/'
                                descriptionTitle='View Proxy'
                            />
                        </StyledCard>
                    </StyledCardSubsectionWrapper>
                    <StyledCardSubsectionWrapper>
                        <StyledCard>
                            <OverviewSmallBoxSummary
                                title={'Envoy Configuration'}
                                logo={
                                    <ItemWrapper>
                                        <LogoHolder>
                                            <EnvoyIcon />
                                        </LogoHolder>
                                        <Success />
                                    </ItemWrapper>
                                }
                                description={
                                    'This is the live config dump from Envoy. This is translated directly from the proxy config and should be updated any time the proxy configuration changes.'
                                }
                                status={GatewayStatus.State.ACCEPTED}
                                count={1}
                                countDescription='Proxy Configurations are configured within Gloo Edge'
                                link='envoy/'
                                descriptionTitle='View Envoy'
                            />
                        </StyledCard>
                    </StyledCardSubsectionWrapper>
                </Main>
                <div>
                    <BottomContent>
                        <h2>APIs</h2>
                    </BottomContent>
                    <BottomWrapper>
                        <CardSubsectionWrapper>
                            <Card>
                                <BottomFooterWrapper>
                                    <h2>
                                        GraphQL
                                    </h2>
                                    <GearWrapper>
                                        <LogoRecolorAndResizeHolder>
                                            <GearIcon />
                                        </LogoRecolorAndResizeHolder>
                                    </GearWrapper>
                                </BottomFooterWrapper>
                                <BottomContent>
                                    Graphql configuration including schema definitions, resolvers, GraphQL APIs, and environments.
                                </BottomContent>
                                <div>
                                    <SoloLink displayElement={`View Settings`} link='graphql' />
                                </div>
                            </Card>
                        </CardSubsectionWrapper>
                    </BottomWrapper>
                </div>
            </CardSubsectionContent>
        </Container>
    )
}
