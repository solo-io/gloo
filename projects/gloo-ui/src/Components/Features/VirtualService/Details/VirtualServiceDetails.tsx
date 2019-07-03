import * as React from 'react';
import { SectionCard } from 'Components/Common/SectionCard';
import { ReactComponent as GlooIcon } from 'assets/GlooEE.svg';
import { Domains } from './Domains';
import { Routes } from './Routes';
import { Configuration } from './Configuration';
import styled from '@emotion/styled/macro';
import { colors, soloConstants } from 'Styles';
import { RouteComponentProps } from 'react-router';
import { Breadcrumb } from 'Components/Common/Breadcrumb';
import { ResourceRef } from 'proto/github.com/solo-io/solo-kit/api/v1/ref_pb';
import { GetVirtualServiceRequest } from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/virtualservice_pb';
import { useGetVirtualService } from 'Api';
import { Route } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/proxy_pb';

const DetailsContent = styled.div`
  display: grid;
  grid-template-rows: auto 2fr 1fr;
  grid-column-gap: 30px;
`;

const DetailsSection = styled.div`
  width: 100%;
`;
export const DetailsSectionTitle = styled.div`
  font-size: 18px;
  font-weight: bold;
  color: ${colors.novemberGrey};
  margin-top: 10px;
  margin-bottom: 10px;
`;

interface Props
  extends RouteComponentProps<{
    virtualservicename: string;
    virtualservicenamespace: string;
  }> {}
const headerInfo = [
  {
    title: 'namespace',
    value: 'dio'
  },
  {
    title: 'namespace',
    value: 'default'
  }
];

export const VirtualServiceDetails = (props: Props) => {
  const { match } = props;
  console.log(match);
  const { virtualservicename, virtualservicenamespace } = match.params;
  let resourceRef = new ResourceRef();
  resourceRef.setName(virtualservicename);
  resourceRef.setNamespace(virtualservicenamespace);
  let vsRequest = new GetVirtualServiceRequest();
  vsRequest.setRef(resourceRef);
  const { data, loading, error } = useGetVirtualService(vsRequest);

  console.log(data);
  if (!data || loading) {
    return <div>Loading...</div>;
  }

  let routes: Route.AsObject[] = [];
  let domains: string[] = [];
  if (data.virtualService && data.virtualService.virtualHost) {
    routes = data.virtualService.virtualHost.routesList;
    domains = data.virtualService.virtualHost.domainsList;
  }

  return (
    <React.Fragment>
      <Breadcrumb />

      <SectionCard
        cardName={match.params ? match.params.virtualservicename : 'test'}
        logoIcon={<GlooIcon />}
        health={1}
        headerSecondaryInformation={headerInfo}
        healthMessage='Service Status'
        closeIcon>
        <DetailsContent>
          <DetailsSection>
            <Domains domains={domains} />
          </DetailsSection>
          <DetailsSection>
            <Routes routes={routes} />
          </DetailsSection>
          <DetailsSection>
            <Configuration />
          </DetailsSection>
        </DetailsContent>
      </SectionCard>
    </React.Fragment>
  );
};
