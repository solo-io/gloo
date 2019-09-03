import styled from '@emotion/styled';
import { useCreateRoute, useUpdateRoute } from 'Api/useVirtualServiceClient';
import { Loading } from 'Components/Common/DisplayOnly/Loading';
import {
  SoloFormDropdown,
  SoloFormInput,
  SoloFormMetadataBasedDropdown,
  SoloFormMultipartStringCardsList,
  SoloFormMultiselect,
  SoloFormVirtualServiceTypeahead
} from 'Components/Common/Form/SoloFormField';
import {
  InputRow,
  SoloFormTemplate
} from 'Components/Common/Form/SoloFormTemplate';
import { SoloButton } from 'Components/Common/SoloButton';
import { Formik, FormikErrors } from 'formik';
import { VirtualService } from 'proto/github.com/solo-io/gloo/projects/gateway/api/v1/virtual_service_pb';
import { DestinationSpec as AWSDestinationSpec } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins/aws/aws_pb';
import { DestinationSpec as AzureDestinationSpec } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins/azure/azure_pb';
import { DestinationSpec as RestDestinationSpec } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins/rest/rest_pb';
import { DestinationSpec } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins_pb';
import {
  ConsulServiceDestination,
  Destination,
  HeaderMatcher,
  Matcher,
  QueryParameterMatcher,
  Route,
  RouteAction
} from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/proxy_pb';
import { Upstream } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/upstream_pb';
import { ResourceRef } from 'proto/github.com/solo-io/solo-kit/api/v1/ref_pb';
import {
  CreateRouteRequest,
  RouteInput,
  VirtualServiceDetails
} from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/virtualservice_pb';
import * as React from 'react';
import { useSelector } from 'react-redux';
import { RouteComponentProps, withRouter } from 'react-router';
import { AppState } from 'store';
import { colors, soloConstants } from 'Styles';
import { ButtonProgress } from 'Styles/CommonEmotions/button';
import { getRouteMatcher } from 'utils/helpers';
import * as yup from 'yup';
import { DestinationForm } from './DestinationForm';

enum PathSpecifierCase { // From gloo -> proxy_pb -> Matcher's namespace
  PATH_SPECIFIER_NOT_SET = 0,
  PREFIX = 1,
  EXACT = 2,
  REGEX = 3
}

export const PATH_SPECIFIERS = [
  {
    key: 'PREFIX',
    value: 'PREFIX',
    displayValue: 'Prefix'
  },
  {
    key: 'EXACT',
    value: 'EXACT',
    displayValue: 'Exact'
  },
  {
    key: 'REGEX',
    value: 'REGEX',
    displayValue: 'Regex'
  }
];

type MethodType =
  | 'POST'
  | 'PUT'
  | 'GET'
  | 'PATCH'
  | 'DELETE'
  | 'HEAD'
  | 'OPTIONS';
type RouteMethodsType = { [key in MethodType]: boolean };
export interface CreateRouteValuesType {
  virtualService: VirtualService.AsObject | undefined;
  upstream: Upstream.AsObject | undefined;
  destinationSpec: DestinationSpec.AsObject | undefined;
  path: string;
  matchType: 'PREFIX' | 'EXACT' | 'REGEX';
  headers: {
    name: string;
    value: string;
    regex: boolean;
  }[];
  queryParameters: {
    name: string;
    value: string;
    regex: boolean;
  }[];
  methods: RouteMethodsType;
}

export const createRouteDefaultValues: CreateRouteValuesType = {
  virtualService: new VirtualService().toObject(),
  upstream: new Upstream().toObject(),
  destinationSpec: undefined,
  path: '',
  matchType: 'PREFIX',
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
  virtualService: yup.object(),
  upstream: yup
    .object()
    .test(
      'Valid upstream',
      'Upstream must be set',
      upstream => !!upstream && !!upstream.metadata
    ),
  destinationSpec: yup.object().when('upstream', {
    is: upstream =>
      !!upstream && !!upstream.upstreamSpec && !!upstream.upstreamSpec.aws,
    then: yup.object().shape({
      aws: yup.object().shape({
        logicalName: yup
          .string()
          .required()
          .test(
            'len',
            'A lambda function must be selected',
            val => !!val && val.length > 0
          ),
        invocationStyle: yup.boolean(),
        responseTransformation: yup.boolean()
      })
    }),
    otherwise: yup.object()
  }),
  path: yup
    .string()
    .test('Valid Path', 'Paths begin with /', val => val && val[0] === '/'),
  matchType: yup.string(),
  headers: yup.array().of(
    yup.object().shape({
      name: yup.string(),
      value: yup.string(),
      regex: yup.boolean()
    })
  ),
  queryParameters: yup.array().of(
    yup.object().shape({
      name: yup.string(),
      value: yup.string(),
      regex: yup.boolean()
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

const FormContainer = styled.form`
  display: flex;
  flex-direction: column;
`;

const InnerSectionTitle = styled.div`
  color: ${colors.novemberGrey};
  font-size: 18px;
  line-height: 22px;
  margin: 13px 0;
`;

const InnerFormSectionContent = styled.div`
  background: white;
  border: 1px solid ${colors.marchGrey};
  border-radius: ${soloConstants.smallRadius}px;
  padding: 13px 8px;
  display: flex;
  flex-direction: column;

  > div {
  }
`;

export const HalfColumn = styled.div`
  width: calc(50% - 10px);
`;
const Footer = styled.div`
  display: flex;
  flex-direction: row;
  justify-content: flex-end;
  margin-top: 28px;
`;

interface Props extends RouteComponentProps {
  defaultVirtualService?: VirtualService.AsObject;
  defaultUpstream?: Upstream.AsObject;
  completeCreation: (newVirtualService?: VirtualService.AsObject) => any;
  existingRoute?: Route.AsObject;
  lockVirtualService?: boolean;
}

export const CreateRouteModalC = (props: Props) => {
  const namespacesList = useSelector(
    (state: AppState) => state.config.namespacesList
  );

  const virtualServicesList = useSelector(
    (state: AppState) => state.virtualServices.virtualServicesList
  );

  const upstreamsList = useSelector((state: AppState) =>
    state.upstreams.upstreamsList.map(u => u.upstream!)
  );

  const [
    allUsableVirtualServices,
    setAllUsableVirtualServices
  ] = React.useState<VirtualServiceDetails.AsObject[]>([]);
  const [allUsableUpstreams, setAllUsableUpstreams] = React.useState<
    Upstream.AsObject[]
  >([]);

  const {
    data: createdVirtualServiceData,
    refetch: makeRequest
  } = useCreateRoute(null);
  const {
    data: updatedVirtualServiceData,
    refetch: makeUpdateRequest
  } = useUpdateRoute(null);

  React.useEffect(() => {
    if (!!createdVirtualServiceData) {
      props.completeCreation(createdVirtualServiceData.virtualService);

      if (
        !!createdVirtualServiceData.virtualServiceDetails &&
        !!createdVirtualServiceData.virtualServiceDetails.virtualService
      ) {
        props.history.push({
          pathname: `/virtualservices/${
            createdVirtualServiceData.virtualServiceDetails.virtualService
              .metadata!.namespace
          }/${
            createdVirtualServiceData.virtualServiceDetails.virtualService
              .metadata!.name
          }`
        });
      }
    }
    if (!!updatedVirtualServiceData) {
      props.completeCreation(updatedVirtualServiceData.virtualService);
    }
  }, [createdVirtualServiceData, updatedVirtualServiceData]);

  React.useEffect(() => {
    setAllUsableVirtualServices(
      !!virtualServicesList
        ? virtualServicesList.filter(vs => !!vs.virtualService!.metadata)
        : []
    );
  }, [virtualServicesList.length]);

  React.useEffect(() => {
    setAllUsableUpstreams(
      !!upstreamsList
        ? upstreamsList.filter(upstream => !!upstream.metadata)
        : []
    );
  }, [upstreamsList.length]);

  if (!upstreamsList.length || !virtualServicesList.length) {
    return <Loading />;
  }
  // if (!!upstreamsError || !!virtualServicesError) {
  //   // @ts-ignore
  //   return <ErrorText>{upstreamsError || virtualServicesError}</ErrorText>;
  // }

  const { defaultUpstream, defaultVirtualService } = props;

  const createRoute = (values: CreateRouteValuesType) => {
    let newRouteReq = new CreateRouteRequest();
    let reqRouteInput = new RouteInput();

    let virtualServiceResourceRef = new ResourceRef();

    if (!!values.virtualService && values.virtualService.metadata) {
      virtualServiceResourceRef.setName(values.virtualService.metadata.name);
      virtualServiceResourceRef.setNamespace(
        values.virtualService.metadata.namespace
      );
      reqRouteInput.setVirtualServiceRef(virtualServiceResourceRef);
    }
    /* -------------------------- ROUTE CREATION BEGINS ------------------------- */
    console.log('values', values);
    let newRoute = new Route();
    let routeMatcher = new Matcher();
    switch (values.matchType) {
      case 'PREFIX':
        routeMatcher.setPrefix(values.path);
        break;
      case 'EXACT':
        routeMatcher.setExact(values.path);
        break;
      case 'REGEX':
        routeMatcher.setRegex(values.path);
        break;
    }

    let matcherHeaders: HeaderMatcher[] = values.headers.map(head => {
      const newMatcherHeader = new HeaderMatcher();
      newMatcherHeader.setName(head.name);
      newMatcherHeader.setValue(head.value);
      newMatcherHeader.setRegex(head.regex);

      return newMatcherHeader;
    });
    routeMatcher.setHeadersList(matcherHeaders);
    let matcherQueryParams: QueryParameterMatcher[] = values.queryParameters.map(
      queryParam => {
        const newMatcherQueryParam = new QueryParameterMatcher();
        newMatcherQueryParam.setName(queryParam.name);
        newMatcherQueryParam.setValue(queryParam.value);
        newMatcherQueryParam.setRegex(queryParam.regex);

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
    const upstreamSpec = values.upstream!.upstreamSpec!;
    let newDestinationResourceRef = new ResourceRef();
    newDestinationResourceRef.setName(values.upstream!.metadata!.name);
    newDestinationResourceRef.setNamespace(
      values.upstream!.metadata!.namespace
    );
    newDestination.setUpstream(newDestinationResourceRef);
    let newDestinationSpec = new DestinationSpec();
    /* ----------------------------- AWS DESTINATION ---------------------------- */
    if (
      !!upstreamSpec.aws &&
      !!values.destinationSpec &&
      values.destinationSpec.aws
    ) {
      const {
        logicalName,
        invocationStyle,
        responseTransformation
      } = values.destinationSpec.aws;
      let newAWSDestinationSpec = new AWSDestinationSpec();
      newAWSDestinationSpec.setLogicalName(logicalName);
      newAWSDestinationSpec.setInvocationStyle(invocationStyle);
      newAWSDestinationSpec.setResponseTransformation(responseTransformation);
      newDestinationSpec.setAws(newAWSDestinationSpec);
      newDestination.setDestinationSpec(newDestinationSpec);
      /* ---------------------------- AZURE DESTINATION --------------------------- */
    } else if (!!upstreamSpec.azure) {
      let newAzureDestinationSpec = new AzureDestinationSpec();
      newDestinationSpec.setAzure(newAzureDestinationSpec);
      newDestination.setDestinationSpec(newDestinationSpec);

      /* ---------------------------- KUBE DESTINATION ---------------------------- */
    } else if (!!upstreamSpec.kube) {
      if (
        values.destinationSpec &&
        values.destinationSpec.rest &&
        values.destinationSpec!.rest!.functionName
      ) {
        let kubeDestination = new RestDestinationSpec();
        kubeDestination.setFunctionName(
          values.destinationSpec!.rest!.functionName
        );
        newDestinationSpec.setRest(kubeDestination);
        newDestination.setDestinationSpec(newDestinationSpec);
      }

      /* --------------------------- CONSUL DESTINATION --------------------------- */
    } else if (!!upstreamSpec.consul) {
      let newConsulServiceDestination = new ConsulServiceDestination();
      // TODO :: I have no idea what to set the values to
      newDestination.setConsul(newConsulServiceDestination);
      let newConsulDestinationSpec;

      newDestination.setDestinationSpec(newConsulDestinationSpec);
    }

    newRouteAction.setSingle(newDestination);
    newRoute.setRouteAction(newRouteAction);

    reqRouteInput.setRoute(newRoute);

    /* --------------------------- ROUTE CREATION ENDS -------------------------- */

    newRouteReq.setInput(reqRouteInput);
    if (!!props.existingRoute) {
      makeUpdateRequest(newRouteReq);
    } else {
      makeRequest(newRouteReq);
    }
  };

  const isSubmittable = (
    errors: FormikErrors<CreateRouteValuesType>,
    isChanged: boolean
  ) => {
    return isChanged && !Object.keys(errors).length;
  };

  let existingRouteToInitialValues = null;

  if (!!props.existingRoute) {
    const { existingRoute } = props;

    let methodsList: RouteMethodsType = {
      POST: false,
      PUT: false,
      GET: false,
      PATCH: false,
      DELETE: false,
      HEAD: false,
      OPTIONS: false
    };
    existingRoute.matcher!.methodsList.forEach(methodName => {
      methodsList[methodName as MethodType] = true;
    });

    const existingRouteUpstream = allUsableUpstreams.find(
      upstream =>
        !!upstream.metadata &&
        !!existingRoute.routeAction &&
        !!existingRoute.routeAction.single &&
        !!existingRoute.routeAction.single.upstream &&
        upstream.metadata!.name ===
          existingRoute.routeAction!.single!.upstream!.name &&
        upstream.metadata!.namespace ===
          existingRoute.routeAction!.single!.upstream!.namespace
    );

    existingRouteToInitialValues = {
      virtualService: props.defaultVirtualService!,
      upstream: existingRouteUpstream,
      destinationSpec: existingRoute.routeAction!.single!.destinationSpec,
      headers: existingRoute.matcher!.headersList,
      methods: methodsList,
      path: getRouteMatcher(existingRoute).matcher,
      matchType: getRouteMatcher(existingRoute).matchType as any,
      queryParameters: existingRoute.matcher!.queryParametersList
    };
  }

  const initialValues: CreateRouteValuesType = !!existingRouteToInitialValues
    ? existingRouteToInitialValues
    : {
        ...createRouteDefaultValues,
        destinationSpec: defaultUpstream
          ? {
              aws: {
                logicalName: '',
                invocationStyle: 0,
                responseTransformation: false
              }
            }
          : undefined,
        virtualService: defaultVirtualService
          ? defaultVirtualService
          : createRouteDefaultValues.virtualService,
        upstream: defaultUpstream
          ? defaultUpstream
          : createRouteDefaultValues.upstream
      };

  return (
    <Formik
      initialValues={initialValues}
      enableReinitialize
      validationSchema={validationSchema}
      onSubmit={createRoute}>
      {({
        values,
        isSubmitting,
        handleSubmit,
        isValid,
        errors,
        dirty,
        setFieldValue
      }) => {
        return (
          <FormContainer data-testid='create-route-form'>
            <SoloFormTemplate>
              <InputRow>
                <React.Fragment>
                  <HalfColumn>
                    <SoloFormVirtualServiceTypeahead
                      name='virtualService'
                      title='Virtual Service'
                      value={values.virtualService}
                      placeholder='Virtual Service...'
                      defaultValue='Virtual Service'
                      options={allUsableVirtualServices.map(
                        vs => vs.virtualService!
                      )}
                      disabled={
                        !allUsableVirtualServices.length ||
                        (!!values.virtualService && props.lockVirtualService)
                      }
                    />
                  </HalfColumn>
                  {allUsableUpstreams.length && (
                    <HalfColumn>
                      <SoloFormMetadataBasedDropdown
                        testId='upstream-dropdown'
                        name='upstream'
                        title='Upstream'
                        value={values.upstream}
                        placeholder='Upstream...'
                        options={allUsableUpstreams}
                        onChange={newUpstream => {
                          if (newUpstream.upstreamSpec.aws) {
                            setFieldValue('destinationSpec', {
                              aws: {
                                logicalName: '',
                                invocationStyle: 0,
                                responseTransformation: false
                              }
                            });
                          }
                        }}
                      />
                    </HalfColumn>
                  )}
                </React.Fragment>
              </InputRow>
              {allUsableUpstreams.length && (
                <InputRow>
                  {!!values.upstream && (
                    <DestinationForm
                      name='destinationSpec'
                      upstreamSpec={values.upstream.upstreamSpec!}
                    />
                  )}
                </InputRow>
              )}
              <InputRow>
                <HalfColumn>
                  <SoloFormInput
                    name='path'
                    title='Path'
                    placeholder='Path...'
                  />
                </HalfColumn>
                <HalfColumn>
                  <SoloFormDropdown
                    name='matchType'
                    title='Match Type'
                    defaultValue={'PREFIX'}
                    options={PATH_SPECIFIERS}
                  />
                </HalfColumn>
              </InputRow>

              <InnerSectionTitle>
                <div>Match Options</div>
              </InnerSectionTitle>
              <InnerFormSectionContent>
                <InputRow>
                  <SoloFormMultipartStringCardsList
                    name='headers'
                    title='Headers'
                    values={values.headers}
                    valuesMayBeEmpty={true}
                    createNewNamePromptText={'Name...'}
                    createNewValuePromptText={'Value...'}
                    booleanFieldText={'Regex'}
                    boolSlotTitle={'regex'}
                  />
                </InputRow>
                <InputRow>
                  <SoloFormMultipartStringCardsList
                    name='queryParameters'
                    title='Query Parameters'
                    values={values.queryParameters}
                    valuesMayBeEmpty={true}
                    createNewNamePromptText={'Name...'}
                    createNewValuePromptText={'Value...'}
                    booleanFieldText={'Regex'}
                    boolSlotTitle={'regex'}
                  />
                </InputRow>
                <InputRow>
                  <SoloFormMultiselect
                    name='methods'
                    title='Methods'
                    placeholder='Methods...'
                    options={Object.keys(createRouteDefaultValues.methods).map(
                      key => {
                        return {
                          key: key,
                          value: key
                        };
                      }
                    )}
                  />
                </InputRow>
              </InnerFormSectionContent>
            </SoloFormTemplate>
            <Footer>
              <SoloButton
                onClick={handleSubmit}
                text={!!props.existingRoute ? 'Save Edits' : 'Create Route'}
                disabled={!isSubmittable(errors, dirty)}
                loading={isSubmitting}
                inProgressText={
                  !!props.existingRoute
                    ? 'Saving Routes...'
                    : 'Creating Route...'
                }>
                <ButtonProgress />
              </SoloButton>
            </Footer>
          </FormContainer>
        );
      }}
    </Formik>
  );
};

export const CreateRouteModal = withRouter(CreateRouteModalC);
