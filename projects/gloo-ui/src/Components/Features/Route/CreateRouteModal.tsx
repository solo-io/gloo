import * as React from 'react';
/** @jsx jsx */
import { jsx } from '@emotion/core';
import styled from '@emotion/styled/macro';
import { colors, soloConstants } from 'Styles';
import {
  SoloFormTemplate,
  InputRow
} from 'Components/Common/Form/SoloFormTemplate';
import {
  SoloFormInput,
  SoloFormTypeahead,
  SoloFormDropdown
} from 'Components/Common/Form/SoloFormField';
import { Field, Formik } from 'formik';
import * as yup from 'yup';

import { VirtualService } from 'proto/github.com/solo-io/gloo/projects/gateway/api/v1/virtual_service_pb';
import { Upstream } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/upstream_pb';
import {
  CreateRouteRequest,
  RouteInput,
  ListVirtualServicesRequest
} from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/virtualservice_pb';
import {
  useCreateRoute,
  useGetUpstreamsList,
  useListVirtualServices
} from 'Api';
import { ResourceRef } from 'proto/github.com/solo-io/solo-kit/api/v1/ref_pb';
import {
  Route,
  Matcher,
  HeaderMatcher,
  QueryParameterMatcher,
  RouteAction,
  Destination,
  KubernetesServiceDestination,
  ConsulServiceDestination,
  RedirectAction
} from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/proxy_pb';
import { DestinationSpec } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins_pb';
import { DestinationSpec as AWSDestinationSpec } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins/aws/aws_pb';
import { DestinationSpec as AzureDestinationSpec } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins/azure/azure_pb';
import { DestinationSpec as RestDestinationSpec } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins/rest/rest_pb';
import { DestinationSpec as GrpcDestinationSpec } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins/grpc/grpc_pb';
import { SoloDropdown } from 'Components/Common/SoloDropdown';
import { ListUpstreamsRequest } from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/upstream_pb';
import { Loading } from 'Components/Common/Loading';
import { ErrorText } from '../VirtualService/Details/ExtAuthForm';
import { NamespacesContext } from 'GlooIApp';
import {
  createUpstreamId,
  parseUpstreamId,
  createVirtualServiceId,
  parseVirtualServiceId
} from 'utils/helpers';

enum PathSpecifierCase { // From gloo -> proxy_pb -> Matcher's namespace
  PATH_SPECIFIER_NOT_SET = 0,
  PREFIX = 1,
  EXACT = 2,
  REGEX = 3
}

interface CreateRouteValuesType {
  virtualServiceId: string | undefined;
  upstreamId: string | undefined;
  path: string;
  matchType: PathSpecifierCase;
  headers: {
    name: string;
    value: string;
    isRegex: boolean;
  }[];
  queryParameters: {
    name: string;
    value: string;
    isRegex: boolean;
  }[];
  methods: {
    POST: boolean;
    PUT: boolean;
    GET: boolean;
    PATCH: boolean;
    DELETE: boolean;
    HEAD: boolean;
    OPTIONS: boolean;
  };
}

const defaultValues: CreateRouteValuesType = {
  virtualServiceId: '',
  upstreamId: '',
  path: '',
  matchType: PathSpecifierCase.PREFIX,
  headers: [],
  queryParameters: [],
  methods: {
    POST: false,
    PUT: false,
    GET: false,
    PATCH: false,
    DELETE: false,
    HEAD: false,
    OPTIONS: false
  }
};

const validationSchema = yup.object().shape({
  region: yup.string(),
  virtualServiceId: yup
    .string()
    .test('len', 'A virtual service must be chosen', val => val.length > 0),
  upstreamId: yup
    .string()
    .test('len', 'A upstream must be chosen', val => val.length > 0),
  path: yup.string(),
  matchType: yup.number(),
  headers: yup.array().of(
    yup.object().shape({
      name: yup.string(),
      value: yup.string(),
      isRegex: yup.boolean()
    })
  ),
  queryParameters: yup.array().of(
    yup.object().shape({
      name: yup.string(),
      value: yup.string(),
      isRegex: yup.boolean()
    })
  ),
  methods: yup.object().shape({
    POST: yup.boolean(),
    PUT: yup.boolean(),
    GET: yup.boolean(),
    PATCH: yup.boolean(),
    DELETE: yup.boolean(),
    HEAD: yup.boolean(),
    OPTIONS: yup.boolean()
  })
});

const FormContainer = styled.div`
  display: flex;
  flex-direction: column;
`;

interface Props {
  defaultVirtualService?: VirtualService.AsObject;
  defaultUpstream?: Upstream.AsObject;
}

export const CreateRouteModal = (props: Props) => {
  const [
    allUsableVirtualServices,
    setAllUsableVirtualServices
  ] = React.useState<VirtualService.AsObject[]>([]);
  const [allUsableUpstreams, setAllUsableUpstreams] = React.useState<
    Upstream.AsObject[]
  >([]);

  const { refetch: makeRequest } = useCreateRoute(null);
  let listVirtualServicesRequest = React.useRef(
    new ListVirtualServicesRequest()
  );
  let listUpstreamsRequest = React.useRef(new ListUpstreamsRequest());
  const namespaces = React.useContext(NamespacesContext);
  listVirtualServicesRequest.current.setNamespacesList(namespaces);
  listUpstreamsRequest.current.setNamespacesList(namespaces);

  const {
    data: upstreamsData,
    error: upstreamsError,
    loading: upstreamsLoading
  } = useGetUpstreamsList(listUpstreamsRequest.current);
  const {
    data: virtualServicesData,
    error: virtualServicesError,
    loading: virtualServicesLoading
  } = useListVirtualServices(listVirtualServicesRequest.current);

  React.useEffect(() => {
    setAllUsableVirtualServices(
      !!virtualServicesData
        ? virtualServicesData.virtualServicesList.filter(vs => !!vs.metadata)
        : []
    );
  }, [virtualServicesData]);
  React.useEffect(() => {
    setAllUsableUpstreams(
      !!upstreamsData
        ? upstreamsData.upstreamsList.filter(upstream => !!upstream.metadata)
        : []
    );
  }, [upstreamsData]);

  if (upstreamsLoading || virtualServicesLoading) {
    return <Loading />;
  }
  if (!!upstreamsError || !!virtualServicesError) {
    // @ts-ignore
    return <ErrorText>{upstreamsError || virtualServicesError}</ErrorText>;
  }

  const { defaultUpstream, defaultVirtualService } = props;

  const createRoute = (values: CreateRouteValuesType) => {
    let newRouteReq = new CreateRouteRequest();
    let reqRouteInput = new RouteInput();

    const { name: vsName, namespace: vsNamespace } = parseVirtualServiceId(
      values.virtualServiceId!
    );
    let virtualServiceResourceRef = new ResourceRef();
    virtualServiceResourceRef.setName(vsName);
    virtualServiceResourceRef.setNamespace(vsNamespace);
    reqRouteInput.setVirtualServiceRef(virtualServiceResourceRef);

    //reqRouteInput.setIndex(vs.virtualHost!.routesList.length);

    const {
      name: upstreamName,
      namespace: upstreamNamespace
    } = parseUpstreamId(values.upstreamId!);
    const usedUpstream = allUsableUpstreams.find(
      upstream =>
        upstream.metadata!.name === upstreamName &&
        upstream.metadata!.namespace === upstreamNamespace
    )!;
    /***
     *  ROUTE CREATION BEGINS
     * */
    let newRoute = new Route();
    let routeMatcher = new Matcher();
    switch (values.matchType) {
      case PathSpecifierCase.EXACT:
        {
          routeMatcher.setPrefix('EXACT');
        }
        break;
      case PathSpecifierCase.REGEX:
        {
          routeMatcher.setPrefix('REGEX');
        }
        break;
      case PathSpecifierCase.PREFIX:
      default:
        {
          routeMatcher.setPrefix('PREFIX');
        }
        break;
    }
    let matcherHeaders: HeaderMatcher[] = values.headers.map(head => {
      const newMatcherHeader = new HeaderMatcher();
      newMatcherHeader.setName(head.name);
      newMatcherHeader.setValue(head.value);
      newMatcherHeader.setRegex(head.isRegex);

      return newMatcherHeader;
    });
    routeMatcher.setHeadersList(matcherHeaders);
    let matcherQueryParams: QueryParameterMatcher[] = values.queryParameters.map(
      queryParam => {
        const newMatcherQueryParam = new QueryParameterMatcher();
        newMatcherQueryParam.setName(queryParam.name);
        newMatcherQueryParam.setValue(queryParam.value);
        newMatcherQueryParam.setRegex(queryParam.isRegex);

        return newMatcherQueryParam;
      }
    );
    routeMatcher.setQueryParametersList(matcherQueryParams);
    routeMatcher.setMethodsList(
      //@ts-ignore
      Object.keys(values.methods).filter(key => values.methods[key])
    );
    newRoute.setMatcher(routeMatcher);

    /* Route->Destination Section */
    let newRouteAction = new RouteAction();
    let newDestination = new Destination();
    const upstreamSpec = usedUpstream.upstreamSpec!;
    let newDestinationResourceRef = new ResourceRef();
    newDestinationResourceRef.setName(usedUpstream.metadata!.name);
    newDestinationResourceRef.setNamespace(usedUpstream!.metadata!.namespace);
    let newDestinationSpec = new DestinationSpec();

    if (!!upstreamSpec.aws) {
      newDestination.setUpstream(newDestinationResourceRef);
      let newAWSDestinationSpec = new AWSDestinationSpec();
      // TODO :: I have no idea what to set the values to
      //newAWSDestinationSpec.setInvocationStyle(0);
      newDestinationSpec.setAws(newAWSDestinationSpec);
    } else if (!!upstreamSpec.azure) {
      newDestination.setUpstream(newDestinationResourceRef);
      let newAzureDestinationSpec = new AzureDestinationSpec();
      // TODO :: I have no idea what to set the values to
      newDestinationSpec.setAzure(newAzureDestinationSpec);
    } /*else if (!!upstreamSpec.kube) {
      let newKubeServiceDestination = new KubernetesServiceDestination();
      newKubeServiceDestination.setRef(newDestinationResourceRef);
      // TODO :: I have no idea what to set the values to
      newDestination.setKube(newKubeServiceDestination);
      let newKubeDestinationSpec;
      // TODO:: How do we tell if it is rest or GRPC?
      //if() -> set DestinationSpec to grpc...
      newDestination.setDestinationSpec(newKubeDestinationSpec);
    } else if (!!upstreamSpec.consul) {
      let newConsulServiceDestination = new ConsulServiceDestination();
      // TODO :: I have no idea what to set the values to
      newDestination.setConsul(newConsulServiceDestination);
      let newConsulDestinationSpec;
      // TODO:: I have no idea what goes in this case
      newDestination.setDestinationSpec(newConsulDestinationSpec);
    }*/
    newDestination.setDestinationSpec(newDestinationSpec);
    newRouteAction.setSingle(newDestination);
    newRoute.setRouteAction(newRouteAction);

    // It looks like we don't see any of the other actions if
    // Route Action is taken??  But if they supplied
    // a path, shouldn't we do the redirect action?
    // Not clear on what the other actions would be based on?

    /*
    let newRedirectAction = new RedirectAction();
    //TODO:: Do we need to set anything else for this???
    if(values.matchType === PathSpecifierCase.PREFIX) {
      // TODO:: Is this correct??
      newRedirectAction.setPrefixRewrite("PREFIX");
    } else {
      newRedirectAction.setPathRedirect(values.path);
    }
    newRoute.setRedirectAction(newRedirectAction);*/

    reqRouteInput.setRoute(newRoute);
    /***
     *  ROUTE CREATION ENDS
     * */

    newRouteReq.setInput(reqRouteInput);
    makeRequest(newRouteReq);
  };

  const initialValues: CreateRouteValuesType = {
    ...defaultValues,
    virtualServiceId: defaultVirtualService
      ? createVirtualServiceId(defaultVirtualService.metadata!)
      : defaultValues.virtualServiceId,
    upstreamId: defaultUpstream
      ? createUpstreamId(defaultUpstream.metadata!)
      : defaultValues.upstreamId
  };

  return (
    <Formik
      initialValues={initialValues}
      validationSchema={validationSchema}
      onSubmit={createRoute}>
      {({ values, isSubmitting, handleSubmit }) => {
        console.log(values);
        console.log(allUsableVirtualServices);

        return (
          <FormContainer>
            <SoloFormTemplate formHeader='Create Route'>
              {allUsableVirtualServices.length && (
                <InputRow>
                  <Field
                    name='virtualServiceId'
                    title='Virtual Service'
                    placeholder='Virtual Service...'
                    options={allUsableVirtualServices.map(vs => {
                      return {
                        key: createVirtualServiceId(
                          vs.metadata!
                        ) /*`${vs.metadata!.name}-${vs.metadata!.namespace}`*/,
                        value: vs.displayName.length
                          ? vs.displayName
                          : vs.metadata!.name
                      };
                    })}
                    component={SoloFormDropdown}
                  />
                </InputRow>
              )}
              <InputRow>
                <Field
                  name='path'
                  title='Path'
                  placeholder='Path...'
                  component={SoloFormInput}
                />
              </InputRow>
              <InputRow>
                <Field
                  name='upstreamId'
                  title='Upstream'
                  placeholder='Upstream...'
                  options={allUsableUpstreams.map(upstream => {
                    return {
                      key: createUpstreamId(
                        upstream.metadata!
                      ) /*`${upstream.metadata!.name}-${
                      upstream.metadata!.namespace
                    }`*/,
                      value: upstream.metadata!.name
                    };
                  })}
                  component={SoloFormDropdown}
                />
              </InputRow>
            </SoloFormTemplate>
          </FormContainer>
        );
      }}
    </Formik>
  );
};
