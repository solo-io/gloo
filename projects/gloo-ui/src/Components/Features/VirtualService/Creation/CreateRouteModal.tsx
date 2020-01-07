import styled from '@emotion/styled';
import {
  RouteDestinationDropdown,
  SoloFormDropdown,
  SoloFormInput,
  SoloFormMultipartStringCardsList,
  SoloFormMultiselect,
  SoloRouteParentDropdown
} from 'Components/Common/Form/SoloFormField';
import {
  InputRow,
  SoloFormTemplate
} from 'Components/Common/Form/SoloFormTemplate';
import { SoloButton } from 'Components/Common/SoloButton';
import { Formik, FormikErrors } from 'formik';
import { uniqBy } from 'lodash';
import { RouteTable } from 'proto/gloo/projects/gateway/api/v1/route_table_pb';
import {
  Route,
  VirtualService
} from 'proto/gloo/projects/gateway/api/v1/virtual_service_pb';
import { Upstream } from 'proto/gloo/projects/gloo/api/v1/upstream_pb';
import { Metadata } from 'proto/solo-kit/api/v1/metadata_pb';
import { VirtualServiceDetails } from 'proto/solo-projects/projects/grpcserver/api/v1/virtualservice_pb';
import * as React from 'react';
import { shallowEqual, useDispatch, useSelector } from 'react-redux';
import { useHistory } from 'react-router';
import { AppState } from 'store';
import { createRoute } from 'store/virtualServices/actions';
import { colors, soloConstants } from 'Styles';
import { ButtonProgress } from 'Styles/CommonEmotions/button';
import { getRouteMatcher } from 'utils/helpers';
import * as yup from 'yup';
import { RouteParent } from '../RouteTableDetails';
import { DestinationForm } from './DestinationForm';
import { updateRouteTable } from 'store/routeTables/actions';
import { DestinationSpec } from 'proto/gloo/projects/gloo/api/v1/options_pb';
import {
  HeaderMatcher,
  QueryParameterMatcher
} from 'proto/gloo/projects/gloo/api/v1/core/matchers/matchers_pb';

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
  background: ${colors.februaryGrey};
  border: 1px solid ${colors.marchGrey};
  border-radius: ${soloConstants.smallRadius}px;
  padding: 13px 8px;
  display: flex;
  flex-direction: column;
`;

export const HalfColumn = styled.div`
  width: calc(50% - 10px);
`;

const Thirds = styled.div`
  width: calc(33.33% - 10px);
`;

const Footer = styled.div`
  display: flex;
  flex-direction: row;
  justify-content: flex-end;
  margin-top: 28px;
`;

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

let httpMethods = ['POST', 'PUT', 'GET', 'PATCH', 'DELETE', 'HEAD', 'OPTIONS'];

export interface CreateRouteValuesType {
  routeParent: { name: string; namespace: string };
  routeParentKind: 'virtualService' | 'routeTable' | '';
  destinationType: 'Upstream' | 'Route Table';
  routeDestination: Upstream.AsObject | RouteTable.AsObject | undefined;
  destinationSpec: DestinationSpec.AsObject | undefined;
  path: string;
  matchType: 'PREFIX' | 'EXACT' | 'REGEX';
  headers: HeaderMatcher.AsObject[];
  queryParameters: QueryParameterMatcher.AsObject[];
  methods: string[];
}

const validationSchema = yup.object().shape({
  region: yup.string(),
  virtualService: yup.object(),
  upstream: yup.object(),
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
  methods: yup.array().of(yup.string())
});

interface Props {
  defaultRouteParent?: RouteTable.AsObject | VirtualService.AsObject;
  defaultUpstream?: Upstream.AsObject;
  completeCreation: () => any;
  existingRoute?: Route.AsObject;
  createRouteFn?: (values: CreateRouteValuesType) => void;
}

function arePropsEqual(
  oldProps: Readonly<Props>,
  newProps: Readonly<Props>
): boolean {
  return true;
}

export const CreateRouteModal = React.memo((props: Props) => {
  const createRouteDefaultValues: CreateRouteValuesType = {
    routeParent: { name: '', namespace: '' },
    routeParentKind:
      !!props.defaultRouteParent && 'virtualHost' in props.defaultRouteParent
        ? 'virtualService'
        : 'routeTable',
    destinationSpec: undefined,
    routeDestination: props.defaultUpstream ? props.defaultUpstream : undefined,
    destinationType: 'Upstream',
    path: '',
    matchType: 'PREFIX',
    headers: [],
    queryParameters: [],
    methods: []
  };

  let history = useHistory();
  const dispatch = useDispatch();
  const virtualServicesList = useSelector(
    (state: AppState) => state.virtualServices.virtualServicesList,
    shallowEqual
  );
  const upstreamsList = useSelector(
    (state: AppState) => state.upstreams.upstreamsList.map(u => u.upstream!),
    shallowEqual
  );

  const routeTablesList = useSelector(
    (state: AppState) => state.routeTables?.routeTablesList
  );

  const upstreamsListMD = useSelector(
    (state: AppState) =>
      state.upstreams.upstreamsList.map(_ => _.upstream!.metadata!),
    shallowEqual
  );

  const routeTablesListMD = useSelector(
    (state: AppState) =>
      state.routeTables.routeTablesList.map(_ => _.routeTable!.metadata!),
    shallowEqual
  );

  let RToptionsList = uniqBy(routeTablesListMD, obj => obj.name);
  let USoptionsList = uniqBy(upstreamsListMD, obj => obj.name);

  const [fallbackVS, setFallbackVS] = React.useState<
    VirtualServiceDetails.AsObject
  >();

  const { defaultUpstream, defaultRouteParent } = props;

  const handleCreateRoute = (values: CreateRouteValuesType) => {
    let virtualService = virtualServicesList.find(
      vsd =>
        vsd?.virtualService?.metadata?.name === values.routeParent?.name &&
        vsd?.virtualService?.metadata?.namespace ===
          values.routeParent?.namespace
    );
    let routeTable = routeTablesList.find(
      rt =>
        rt?.routeTable?.metadata?.name === values.routeParent?.name &&
        rt?.routeTable?.metadata?.namespace === values.routeParent?.namespace
    );

    if (routeTable !== undefined && values.routeParentKind === 'routeTable') {
      let newRoutesList = routeTable?.routeTable?.routesList || [];

      let destination;
      if (values.destinationType === 'Route Table') {
        destination = {
          delegateAction: {
            name: values.routeDestination!.metadata!.name,
            namespace: values.routeDestination!.metadata!.namespace
          }
        };
      } else if (values.destinationType === 'Upstream') {
        let destinationSpec;
        if (values.destinationSpec !== undefined) {
          destinationSpec = values.destinationSpec;
        }
        destination = {
          routeAction: {
            single: {
              upstream: {
                name: values.routeDestination!.metadata!.name,
                namespace: values.routeDestination!.metadata!.namespace
              },
              destinationSpec
            }
          }
        };
      }
      dispatch(
        updateRouteTable({
          routeTable: {
            ...routeTable,
            metadata: {
              ...routeTable.routeTable?.metadata!,
              name: values.routeParent?.name,
              namespace: values?.routeParent?.namespace
            },
            routesList: [
              {
                matchersList: [
                  {
                    prefix: values.matchType === 'PREFIX' ? values.path : '',
                    exact: values.matchType === 'EXACT' ? values.path : '',
                    regex: values.matchType === 'REGEX' ? values.path : '',
                    methodsList: values.methods,
                    headersList: values.headers,
                    queryParametersList: values.queryParameters
                  }
                ],
                ...destination
              },
              ...newRoutesList
            ]
          }
        })
      );
      history.push({
        pathname: `/routetables/${values.routeParent?.namespace}/${values.routeParent?.name}`
      });
      props.completeCreation();
    } else if (
      virtualService !== undefined &&
      values.routeParentKind === 'virtualService'
    ) {
      // vs -> rt
      if (values.destinationType === 'Route Table') {
        dispatch(
          createRoute({
            input: {
              ...(!!values.routeParent && {
                virtualServiceRef: {
                  name: values.routeParent?.name,
                  namespace: values.routeParent?.namespace
                }
              }),
              index: 0,
              route: {
                matchersList: [
                  {
                    prefix: values.matchType === 'PREFIX' ? values.path : '',
                    exact: values.matchType === 'EXACT' ? values.path : '',
                    regex: values.matchType === 'REGEX' ? values.path : '',
                    methodsList: values.methods,
                    headersList: values.headers,
                    queryParametersList: values.queryParameters
                  }
                ],
                delegateAction: {
                  name: values.routeDestination!.metadata!.name,
                  namespace: values.routeDestination!.metadata!.namespace
                }
              }
            }
          })
        );
        props.completeCreation();
      } else if (values.destinationType === 'Upstream') {
        let destinationSpec;
        if (values.destinationSpec !== undefined) {
          destinationSpec = values.destinationSpec;
        }
        // vs -> us
        dispatch(
          createRoute({
            input: {
              ...(!!values.routeParent && {
                virtualServiceRef: {
                  name: values.routeParent?.name,
                  namespace: values.routeParent?.namespace
                }
              }),
              index: 0,
              route: {
                matchersList: [
                  {
                    prefix: values.matchType === 'PREFIX' ? values.path : '',
                    exact: values.matchType === 'EXACT' ? values.path : '',
                    regex: values.matchType === 'REGEX' ? values.path : '',
                    methodsList: values.methods,
                    headersList: values.headers,
                    queryParametersList: values.queryParameters
                  }
                ],
                routeAction: {
                  single: {
                    upstream: {
                      name: values.routeDestination!.metadata!.name,
                      namespace: values.routeDestination!.metadata!.namespace
                    },
                    destinationSpec
                  }
                }
              }
            }
          })
        );
        props.completeCreation();

        if (virtualService !== undefined) {
          history.push({
            pathname: `/virtualservices/${values.routeParent?.namespace}/${values.routeParent?.name}`
          });
        } else {
          if (!!fallbackVS && !!fallbackVS.virtualService) {
            history.push({
              pathname: `/virtualservices/${
                fallbackVS.virtualService.metadata!.namespace
              }/${fallbackVS.virtualService.metadata!.name}`
            });
          } else {
            history.push({
              pathname: `/virtualservices/gloo-system/default`
            });
          }
        }
      }
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

    const existingRouteUpstream = upstreamsList.find(
      upstream =>
        upstream?.metadata?.name ===
          existingRoute?.routeAction?.single?.upstream?.name &&
        upstream?.metadata?.namespace ===
          existingRoute?.routeAction?.single?.upstream?.namespace
    );

    existingRouteToInitialValues = {
      routeParent: props.defaultRouteParent ? props.defaultRouteParent! : '',
      routeDestination: existingRouteUpstream,
      destinationSpec: existingRoute.routeAction!.single!.destinationSpec,
      destinationType: 'Upstream',
      headers: existingRoute.matchersList[0]?.headersList,
      methods: existingRoute.matchersList[0]?.methodsList,
      path: getRouteMatcher(existingRoute).matcher,
      matchType: getRouteMatcher(existingRoute).matchType as any,
      queryParameters: existingRoute.matchersList[0]?.queryParametersList
    };
  }

  const initialValues: CreateRouteValuesType = {
    ...createRouteDefaultValues,
    routeParent: {
      name: defaultRouteParent?.metadata?.name || '',
      namespace: defaultRouteParent?.metadata?.namespace || ''
    },

    destinationSpec:
      defaultUpstream && defaultUpstream?.aws !== undefined
        ? {
            aws: {
              logicalName: '',
              invocationStyle: 0,
              responseTransformation: false
            }
          }
        : undefined,
    destinationType: createRouteDefaultValues.destinationType
    // routeDestination:
    //   createRouteDefaultValues.routeDestination || defaultUpstream
    // routeTable: createRouteDefaultValues.routeTable
  };

  return (
    <Formik
      initialValues={initialValues}
      enableReinitialize
      validationSchema={validationSchema}
      onSubmit={handleCreateRoute}>
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
                <Thirds>
                  <div>
                    <SoloRouteParentDropdown
                      title='Route Container'
                      name='routeParent'
                      defaultValue={defaultRouteParent?.metadata?.name!}
                    />
                  </div>
                </Thirds>

                <Thirds>
                  <SoloFormDropdown
                    name='destinationType'
                    testId='destination-dropdown'
                    title='Destination Type'
                    defaultValue={'Upstream'}
                    options={['Upstream', 'Route Table'].map(dest => {
                      return { key: dest, value: dest };
                    })}
                  />
                </Thirds>
                {upstreamsList.length > 0 && (
                  <Thirds>
                    <RouteDestinationDropdown
                      testId='upstream-dropdown'
                      name='routeDestination'
                      title='Destination'
                      value={values.routeDestination!}
                      optionsList={
                        values.destinationType === 'Route Table'
                          ? RToptionsList
                          : USoptionsList
                      }
                      destinationType={values.destinationType}
                      onChange={newUpstream => {
                        if (newUpstream?.aws !== undefined) {
                          setFieldValue('destinationSpec', {
                            aws: {
                              logicalName: '',
                              invocationStyle: 0,
                              responseTransformation: false
                            }
                          });
                        } else {
                          setFieldValue('destinationSpec', undefined);
                        }
                      }}
                    />
                  </Thirds>
                )}
              </InputRow>
              <InputRow>
                {values.routeDestination !== undefined &&
                  values.destinationType === 'Upstream' &&
                  'aws' in values.routeDestination && (
                    <DestinationForm
                      name='destinationSpec'
                      upstream={values.routeDestination}
                    />
                  )}
              </InputRow>
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
                    disabled={!!props.createRouteFn}
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
                    options={httpMethods.map(key => {
                      return {
                        key: key,
                        value: key
                      };
                    })}
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
}, arePropsEqual);
