import * as React from 'react';
import { SectionCard } from 'Components/Common/SectionCard';
import { ReactComponent as GlooIcon } from 'assets/GlooEE.svg';
import { Domains } from './Domains';
import { Routes } from './Routes';
import { Configuration } from './Configuration';
import styled from '@emotion/styled/macro';
import { colors, soloConstants, healthConstants } from 'Styles';
import { RouteComponentProps } from 'react-router';
import { Breadcrumb } from 'Components/Common/Breadcrumb';
import { ResourceRef } from 'proto/github.com/solo-io/solo-kit/api/v1/ref_pb';
import {
  GetVirtualServiceRequest,
  UpdateVirtualServiceRequest,
  VirtualServiceInput
} from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/virtualservice_pb';
import { useGetVirtualService, useUpdateVirtualService } from 'Api';
import {
  Route,
  Matcher,
  HeaderMatcher,
  QueryParameterMatcher,
  RouteAction,
  Destination
} from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/proxy_pb';
import { ErrorText } from './ExtAuthForm';
import { DestinationSpec } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins_pb';
import { DestinationSpec as AWSDestinationSpec } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins/aws/aws_pb';
import { DestinationSpec as AzureDestinationSpec } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins/azure/azure_pb';
import { DestinationSpec as RestDestinationSpec } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins/rest/rest_pb';
import { DestinationSpec as GrpcDestinationSpec } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins/grpc/grpc_pb';
import {
  IngressRateLimit,
  RateLimit
} from 'proto/github.com/solo-io/solo-projects/projects/gloo/api/v1/plugins/ratelimit/ratelimit_pb';
import {
  OAuth,
  CustomAuth
} from 'proto/github.com/solo-io/solo-projects/projects/gloo/api/v1/plugins/extauth/extauth_pb';
import { VirtualService } from 'proto/github.com/solo-io/gloo/projects/gateway/api/v1/virtual_service_pb';

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

export const VirtualServiceDetails = (props: Props) => {
  const { match, history } = props;
  const { virtualservicename, virtualservicenamespace } = match.params;

  const [virtualService, setVirtualService] = React.useState<
    VirtualService.AsObject | undefined
  >(undefined);

  let resourceRef = new ResourceRef();
  resourceRef.setName(virtualservicename);
  resourceRef.setNamespace(virtualservicenamespace);
  let vsRequest = new GetVirtualServiceRequest();
  vsRequest.setRef(resourceRef);
  const { data, loading, error, refetch: getVSRefetch } = useGetVirtualService(
    vsRequest
  );

  const {
    data: updateData,
    loading: updateLoading,
    refetch: makeUpdateRequest
  } = useUpdateVirtualService(null);

  React.useEffect(() => {
    console.log(data);
    if (!!data) {
      setVirtualService(data.virtualService);
    }
  }, [loading]);

  React.useEffect(() => {
    console.log(updateData);
    if (!!updateData) {
      setVirtualService(updateData.virtualService);
    }
  }, [updateLoading]);

  console.log({ data, loading, updateLoading });
  if (!virtualService && (loading || updateLoading)) {
    return (
      <React.Fragment>
        <Breadcrumb />

        <SectionCard
          cardName={'Loading...'}
          logoIcon={<GlooIcon />}
          health={healthConstants.Pending.value}
          healthMessage={'Loading...'}
          onClose={() => history.push(`/virtualservices/`)}
        />
      </React.Fragment>
    );
  }
  if (!!error || !virtualService) {
    return <ErrorText>{error}</ErrorText>;
  }

  const reloadVirtualService = (
    newVirtualService?: VirtualService.AsObject
  ) => {
    if (newVirtualService) {
      setVirtualService(newVirtualService);
    } else {
      getVSRefetch(vsRequest);
    }
  };

  let routes: Route.AsObject[] = [];
  let domains: string[] = [];
  if (!!virtualService!.virtualHost) {
    routes = virtualService!.virtualHost!.routesList;
    domains = virtualService!.virtualHost!.domainsList;
  }

  const updateVirtualService = (newInfo: {
    newDomainsList?: string[];
    newRoutesList?: Route.AsObject[];
    newRateLimits?: IngressRateLimit.AsObject;
    newOAuth?: OAuth.AsObject;
  }) => {
    let virtualServiceInput = new VirtualServiceInput();
    let vsRef = new ResourceRef();
    vsRef.setName(virtualService!.metadata!.name);
    vsRef.setNamespace(virtualService!.metadata!.namespace);
    virtualServiceInput.setRef(vsRef);
    virtualServiceInput.setDisplayName(virtualService!.displayName);
    virtualServiceInput.setDomainsList(
      !!newInfo.newDomainsList
        ? newInfo.newDomainsList
        : virtualService!.virtualHost!.domainsList
    );
    const routesList: Route[] = (!!newInfo.newRoutesList
      ? newInfo.newRoutesList
      : virtualService!.virtualHost!.routesList
    ).map((rt: Route.AsObject) => {
      let newRoute = new Route();

      let routeMatcher = new Matcher();
      if (!!rt.matcher!.prefix) {
        routeMatcher.setPrefix(rt.matcher!.prefix);
      } else if (!!rt.matcher!.exact) {
        routeMatcher.setExact(rt.matcher!.exact);
      } else if (!!rt.matcher!.regex) {
        routeMatcher.setRegex(rt.matcher!.regex);
      }

      let matcherHeaders: HeaderMatcher[] = rt.matcher!.headersList.map(
        head => {
          const newMatcherHeader = new HeaderMatcher();
          newMatcherHeader.setName(head.name);
          newMatcherHeader.setValue(head.value);
          newMatcherHeader.setRegex(head.regex);

          return newMatcherHeader;
        }
      );
      routeMatcher.setHeadersList(matcherHeaders);
      let matcherQueryParams: QueryParameterMatcher[] = rt.matcher!.queryParametersList.map(
        queryParam => {
          const newMatcherQueryParam = new QueryParameterMatcher();
          newMatcherQueryParam.setName(queryParam.name);
          newMatcherQueryParam.setValue(queryParam.value);
          newMatcherQueryParam.setRegex(queryParam.regex);

          return newMatcherQueryParam;
        }
      );
      routeMatcher.setQueryParametersList(matcherQueryParams);
      routeMatcher.setMethodsList(rt.matcher!.methodsList);
      newRoute.setMatcher(routeMatcher);

      let newRouteAction = new RouteAction();
      let newDestination = new Destination();

      if (!!rt.routeAction!.single) {
        const singleDestination = rt.routeAction!.single!;

        if (!!singleDestination.upstream) {
          let newDestinationResourceRef = new ResourceRef();
          newDestinationResourceRef.setName(singleDestination.upstream!.name);
          newDestinationResourceRef.setNamespace(
            singleDestination.upstream!.namespace
          );
          newDestination.setUpstream(newDestinationResourceRef);
        }

        if (!!singleDestination.destinationSpec!) {
          let newDestinationSpec = new DestinationSpec();
          if (!!singleDestination.destinationSpec!.aws) {
            const currentAWS = singleDestination.destinationSpec!.aws!;
            let newAWSDestinationSpec = new AWSDestinationSpec();
            newAWSDestinationSpec.setInvocationStyle(
              currentAWS.invocationStyle
            );
            newAWSDestinationSpec.setLogicalName(currentAWS.logicalName);
            newAWSDestinationSpec.setResponseTransformation(
              currentAWS.responseTransformation
            );
            newDestinationSpec.setAws(newAWSDestinationSpec);
          } else if (!!singleDestination.destinationSpec!.azure) {
            const currentAzure = singleDestination.destinationSpec!.azure!;
            let newAzureDestinationSpec = new AzureDestinationSpec();
            newAzureDestinationSpec.setFunctionName(currentAzure.functionName);
            newDestinationSpec.setAzure(newAzureDestinationSpec);
          }
          newDestination.setDestinationSpec(newDestinationSpec);
        }
        newRouteAction.setSingle(newDestination);
      }
      newRoute.setRouteAction(newRouteAction);

      return newRoute;
    });
    virtualServiceInput.setRoutesList(routesList);

    if (!!virtualService!.sslConfig && !!virtualService!.sslConfig!.secretRef) {
      let secretRef = new ResourceRef();
      secretRef.setName(virtualService!.sslConfig!.secretRef!.name);
      secretRef.setNamespace(virtualService!.sslConfig!.secretRef!.namespace);
      virtualServiceInput.setSecretRef(secretRef);
    }

    if (
      !!virtualService!.virtualHost &&
      !!virtualService!.virtualHost!.virtualHostPlugins &&
      !!virtualService!.virtualHost!.virtualHostPlugins!.extensions &&
      !!virtualService!.virtualHost!.virtualHostPlugins!.extensions!.configsMap
    ) {
      const configsMap = virtualService!.virtualHost!.virtualHostPlugins!
        .extensions!.configsMap;

      /** RATE LIMITS */
      let rateLimits = new IngressRateLimit();
      const currentRateLimitIndex = configsMap.findIndex(
        config => config[0] === 'rate-limit'
      );
      const usedRateLimits = newInfo.newRateLimits
        ? newInfo.newRateLimits
        : currentRateLimitIndex !== -1
        ? configsMap[currentRateLimitIndex][1]
        : undefined;
      if (!!usedRateLimits) {
        //@ts-ignore
        if (!!usedRateLimits.anonymousLimits) {
          const anonLimit = new RateLimit();
          //@ts-ignore
          anonLimit.setUnit(usedRateLimits.anonymousLimits!.unit);
          anonLimit.setRequestsPerUnit(
            //@ts-ignore
            usedRateLimits.anonymousLimits!.requestsPerUnit
          );
          rateLimits.setAnonymousLimits(anonLimit);
        }
        //@ts-ignore
        if (!!usedRateLimits.authorizedLimits) {
          const authLimit = new RateLimit();
          //@ts-ignore
          authLimit.setUnit(usedRateLimits.authorizedLimits!.unit);
          authLimit.setRequestsPerUnit(
            //@ts-ignore
            usedRateLimits.authorizedLimits!.requestsPerUnit
          );
          rateLimits.setAuthorizedLimits(authLimit);
        }
      }
      virtualServiceInput.setRateLimitConfig(rateLimits);

      /** AUTHORIZATIONS */
      const basicAuthIndex = configsMap.findIndex(
        config => config[0] === 'basic-auth'
      );
      if (basicAuthIndex !== -1) {
        const existingBasicAuth = configsMap[basicAuthIndex][1];
        let basicAuth = new VirtualServiceInput.BasicAuthInput();
        // @ts-ignore
        basicAuth.setSpecCsv(existingBasicAuth.specCsv);
        // @ts-ignore
        basicAuth.setRealm(existingBasicAuth.realm);
        virtualServiceInput.setBasicAuth(basicAuth);
      }
      const oAuthIndex = configsMap.findIndex(config => config[0] === 'oauth');
      if (oAuthIndex !== -1) {
        const existingOAuth = configsMap[oAuthIndex][1];
        const usedOAuth = newInfo.newOAuth || existingOAuth;
        let oAuth = new OAuth();
        // @ts-ignore
        oAuth.setClientId(usedOAuth.clientId);
        // @ts-ignore
        oAuth.setCallbackPath(usedOAuth.callbackPath);
        // @ts-ignore
        oAuth.setIssuerUrl(usedOAuth.issuerUrl);
        // @ts-ignore
        oAuth.setAppUrl(usedOAuth.appUrl);
        let clientSecretRef = new ResourceRef();
        // @ts-ignore
        clientSecretRef.setName(usedOAuth.clientSecretRef!.name);
        clientSecretRef.setNamespace(
          // @ts-ignore
          usedOAuth.clientSecretRef!.namespace
        );
        oAuth.setClientSecretRef(clientSecretRef);
        virtualServiceInput.setOauth(oAuth);
      }
      if (!!configsMap.find(config => config[0] === 'custom-auth')) {
        let customAuth = new CustomAuth();
        virtualServiceInput.setCustomAuth(customAuth);
      }
    }

    let updateRequest = new UpdateVirtualServiceRequest();
    updateRequest.setInput(virtualServiceInput);
    makeUpdateRequest(updateRequest);
  };

  const domainsChanged = (newDomainsList: string[]) => {
    updateVirtualService({ newDomainsList });
  };
  const routesChanged = (newRoutesList: Route.AsObject[]) => {
    updateVirtualService({ newRoutesList });
  };

  const headerInfo = [
    {
      title: 'namespace',
      value: virtualservicenamespace
    }
  ];

  return (
    <React.Fragment>
      <Breadcrumb />

      <SectionCard
        cardName={match.params ? match.params.virtualservicename : 'test'}
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
        <DetailsContent>
          <DetailsSection>
            <Domains domains={domains} domainsChanged={domainsChanged} />
          </DetailsSection>
          <DetailsSection>
            <Routes
              routes={routes}
              virtualService={virtualService!}
              routesChanged={routesChanged}
              reloadVirtualService={reloadVirtualService}
            />
          </DetailsSection>
          <DetailsSection>
            <Configuration />
          </DetailsSection>
        </DetailsContent>
      </SectionCard>
    </React.Fragment>
  );
};
