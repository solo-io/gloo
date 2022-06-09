import React, { useEffect, useState } from 'react';
import styled from '@emotion/styled/macro';
import { Loading } from 'Components/Common/Loading';
import { useParams } from 'react-router';
import { HealthNotificationBox } from 'Components/Common/HealthNotificationBox';
import { SectionCard } from 'Components/Common/SectionCard';
import { ReactComponent as GlooIcon } from 'assets/Gloo.svg';
import { SubRouteTablesTable } from './RouteTable/SubRouteTablesTable';
import { useListVirtualServices } from 'API/hooks';
import { VirtualService } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/gateway_resources_pb';
import { StringCardsList } from 'Components/Common/StringCardsList';
import { AreaTitle } from 'Styles/StyledComponents/headings';
import { RateLimitSection } from './RouteTable/RateLimit';
import { ExternalAuthorizationSection } from './RouteTable/ExtAuth';
import {
  CardSubsectionWrapper,
  CardSubsectionContent,
} from 'Components/Common/Card';
import { DataError } from 'Components/Common/DataError';
import { Empty } from 'antd';
import { EmptyAsterisk } from 'Components/Common/EmptyAsterisk';

const GlooIconHolder = styled.div`
  display: flex;
  align-items: center;
  justify-items: center;

  svg {
    width: 25px;
  }
`;

const DomainsArea = styled.div`
  margin-bottom: 20px;
`;
const RoutesArea = styled.div`
  margin-bottom: 20px;
`;
const ConfigurationArea = styled.div`
  > div {
    display: grid;
    grid-template-columns: 1fr 1fr;
    grid-gap: 22px;
  }
`;

export const VirtualServiceDetails = () => {
  const { name, namespace, virtualservicename, virtualservicenamespace } =
    useParams();

  const { data: vsResponse, error: vsError } = useListVirtualServices({
    name: name!,
    namespace: namespace!,
  });
  const allVirtualServices = vsResponse?.virtualServicesList;

  const [virtualService, setVirtualService] =
    useState<VirtualService.AsObject>();

  useEffect(() => {
    if (!!allVirtualServices) {
      setVirtualService(
        allVirtualServices.find(
          vs =>
            vs.metadata?.name === virtualservicename &&
            vs.metadata?.namespace === virtualservicenamespace
        )
      );
    } else {
      setVirtualService(undefined);
    }
  }, [
    name,
    namespace,
    allVirtualServices,
    virtualservicename,
    virtualservicenamespace,
  ]);

  if (!!vsError) {
    return <DataError error={vsError} />;
  } else if (!allVirtualServices) {
    return <Loading message={'Retrieving virtual services...'} />;
  }

  return (
    <SectionCard
      cardName={virtualservicename!}
      logoIcon={
        <GlooIconHolder>
          <GlooIcon />
        </GlooIconHolder>
      }
      headerSecondaryInformation={[
        {
          title: 'Namespace',
          value: virtualservicenamespace,
        },
      ]}
      health={{
        state: virtualService?.status?.state ?? 0,
        reason: virtualService?.status?.reason,
      }}>
      {!!vsError ? (
        <DataError error={vsError} />
      ) : !virtualService ? (
        <Loading
          message={`Retrieving information for the virtual service ${virtualservicename}...`}
        />
      ) : (
        <>
          <HealthNotificationBox
            state={virtualService?.status?.state}
            reason={virtualService?.status?.reason}
          />
          <DomainsArea>
            <AreaTitle>Domains</AreaTitle>
            <CardSubsectionWrapper>
              <StringCardsList
                values={
                  virtualService.spec?.virtualHost?.domainsList.map(dom =>
                    dom === '*' ? (
                      <div style={{ padding: '0 10px' }}>
                        <EmptyAsterisk />
                      </div>
                    ) : (
                      dom
                    )
                  ) ?? []
                }
              />
            </CardSubsectionWrapper>
          </DomainsArea>

          <RoutesArea>
            <AreaTitle>Routes</AreaTitle>
            <SubRouteTablesTable onlyTable={true} />
          </RoutesArea>
          <ConfigurationArea>
            <AreaTitle>Configuration</AreaTitle>
            <CardSubsectionWrapper>
              <CardSubsectionContent>
                <ExternalAuthorizationSection
                  externalAuth={
                    virtualService.spec?.virtualHost?.options?.extauth
                  }
                />
              </CardSubsectionContent>
              <CardSubsectionContent>
                <RateLimitSection
                  rateLimits={
                    virtualService.spec?.virtualHost?.options?.ratelimitBasic
                  }
                />
              </CardSubsectionContent>
            </CardSubsectionWrapper>
          </ConfigurationArea>
        </>
      )}
    </SectionCard>
  );
};
