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
import { Struct } from 'google-protobuf/google/protobuf/struct_pb';
import { VirtualService } from 'proto/github.com/solo-io/gloo/projects/gateway/api/v1/virtual_service_pb';

const DetailsContent = styled.div`
  display: grid;
  grid-template-rows: auto 1fr 1fr;
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
    if (!!data) {
      setVirtualService(data.virtualService);
    }
  }, [loading]);

  React.useEffect(() => {
    if (!!updateData) {
      setVirtualService(updateData.virtualService);
    }
  }, [updateLoading]);

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

  let configsMap: Map<string, Struct.AsObject> | undefined = undefined;
  let rateLimits: IngressRateLimit.AsObject | undefined = undefined;
  let externalAuth: OAuth.AsObject | undefined = undefined;
  if (
    !!virtualService.virtualHost &&
    !!virtualService.virtualHost!.virtualHostPlugins &&
    !!virtualService.virtualHost!.virtualHostPlugins!.extensions
  ) {
    configsMap = new Map(
      virtualService.virtualHost!.virtualHostPlugins!.extensions!.configsMap
    );
  }
  if (!!configsMap && !!configsMap.get('rate-limit')) {
    const fieldsMap = new Map(configsMap.get('rate-limit')!.fieldsMap);

    let anonLimit = undefined;
    if (!!fieldsMap.get('anonymous_limits')) {
      const structValues = new Map(
        fieldsMap.get('anonymous_limits')!.structValue!.fieldsMap
      );

      anonLimit = {
        // @ts-ignore
        unit: RateLimit.Unit[structValues.get('unit')!.stringValue],
        requestsPerUnit: structValues.get('requests_per_unit')!.numberValue
      };
    }
    let authLimit = undefined;
    if (!!fieldsMap.get('authorized_limits')) {
      const structValues = new Map(
        fieldsMap.get('authorized_limits')!.structValue!.fieldsMap
      );

      authLimit = {
        // @ts-ignore
        unit: RateLimit.Unit[structValues.get('unit')!.stringValue],
        requestsPerUnit: structValues.get('requests_per_unit')!.numberValue
      };
    }

    rateLimits = {
      anonymousLimits: anonLimit,
      authorizedLimits: authLimit
    };
  }
  if (!!configsMap && !!configsMap.get('extauth')) {
    let fieldsMap = new Map(configsMap.get('extauth')!.fieldsMap);
    if (!!fieldsMap.get('oauth')) {
      fieldsMap = new Map(fieldsMap.get('oauth')!.structValue!.fieldsMap);
    }

    const appUrl = fieldsMap.get('app_url')!.stringValue;
    const clientId = fieldsMap.get('client_id')!.stringValue;
    const issuerUrl = fieldsMap.get('issuer_url')!.stringValue;
    const callbackPath = fieldsMap.get('callback_path')!.stringValue;
    let clientSecretRef = undefined;
    if (
      !!fieldsMap.get('client_secret_ref') &&
      !!fieldsMap.get('client_secret_ref')!.stringValue.length
    ) {
      const structValues = new Map(
        fieldsMap.get('client_secret_ref')!.structValue!.fieldsMap
      );

      clientSecretRef = {
        name: structValues.get('name')!.stringValue,
        namespace: structValues.get('namespace')!.stringValue
      };
    }

    externalAuth = {
      clientId,
      clientSecretRef,
      issuerUrl,
      appUrl,
      callbackPath
    };
  }

  const updateVirtualService = (newInfo: {
    newDomainsList?: string[];
    newRoutesList?: Route.AsObject[];
    newRateLimits?: IngressRateLimit.AsObject;
    newOAuth?: OAuth.AsObject;
  }) => {
    console.log(newInfo);
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

    /** RATE LIMITS */
    let newRateLimits = new IngressRateLimit();
    const usedRateLimits = !!newInfo.newRateLimits
      ? newInfo.newRateLimits
      : !!configsMap && !!configsMap.get('rate-limit')
      ? configsMap.get('rate-limit')
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
        newRateLimits.setAnonymousLimits(anonLimit);
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
        newRateLimits.setAuthorizedLimits(authLimit);
      }
    }
    virtualServiceInput.setRateLimitConfig(newRateLimits);

    /** AUTHORIZATIONS */
    /*if (!!configsMap && !!configsMap.get('basic-auth')) {
      const existingBasicAuth = configsMap.get('basic-auth');
      let basicAuth = new VirtualServiceInput.BasicAuthInput();
      // @ts-ignore
      basicAuth.setSpecCsv(existingBasicAuth.specCsv);
      // @ts-ignore
      basicAuth.setRealm(existingBasicAuth.realm);
      virtualServiceInput.setBasicAuth(basicAuth);
    }*/
    if (!!configsMap && !!configsMap.get('extauth')) {
      if (newInfo.newOAuth) {
        const usedOAuth = newInfo.newOAuth || configsMap!.get('extauth');
        let oAuth = new OAuth();
        // @ts-ignore
        oAuth.setClientId(usedOAuth.clientId);
        // @ts-ignore
        oAuth.setCallbackPath(usedOAuth.callbackPath);
        // @ts-ignore
        oAuth.setIssuerUrl(usedOAuth.issuerUrl);
        // @ts-ignore
        oAuth.setAppUrl(usedOAuth.appUrl);
        // @ts-ignore
        if (!!usedOAuth!.clientSecretRef) {
          let clientSecretRef = new ResourceRef();
          // @ts-ignore
          clientSecretRef.setName(usedOAuth.clientSecretRef!.name);
          clientSecretRef.setNamespace(
            // @ts-ignore
            usedOAuth.clientSecretRef!.namespace
          );
          oAuth.setClientSecretRef(clientSecretRef);
        }
        virtualServiceInput.setOauth(oAuth);
      } else {
      }
    }
    if (!!configsMap && !!configsMap.get('custom-auth')) {
      let customAuth = new CustomAuth();
      virtualServiceInput.setCustomAuth(customAuth);
    }

    let updateRequest = new UpdateVirtualServiceRequest();
    updateRequest.setInput(virtualServiceInput);
    makeUpdateRequest(updateRequest);
  };

  const domainsChanged = (newDomainsList: string[]) => {
    updateVirtualService({ newDomainsList });
  };
  const ratesChanged = (newRateLimits: IngressRateLimit.AsObject) => {
    updateVirtualService({ newRateLimits });
  };
  const externalAuthChanged = (newOAuth: OAuth.AsObject) => {
    updateVirtualService({ newOAuth });
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
        cardName={
          virtualService.displayName.length
            ? virtualService.displayName
            : match.params
            ? match.params.virtualservicename
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
            <Configuration
              externalAuth={externalAuth}
              externalAuthChanged={externalAuthChanged}
              rates={rateLimits}
              rateLimitsChanged={ratesChanged}
            />
          </DetailsSection>
        </DetailsContent>
      </SectionCard>
    </React.Fragment>
  );
};
