import styled from '@emotion/styled/macro';
import { TabList, TabPanel, TabPanels } from '@reach/tabs';
import { useGetGraphqlSchemaDetails, useListUpstreams } from 'API/hooks';
import { OptionType } from 'Components/Common/SoloDropdown';
import {
  SoloFormDropdown,
  SoloFormRadio,
  SoloFormRadioOption,
} from 'Components/Common/SoloFormComponents';
import { StyledModalTab, StyledModalTabs } from 'Components/Common/SoloModal';
import YamlDisplayer from 'Components/Common/YamlDisplayer';
import YamlEditor from 'Components/Common/YamlEditor';
import { Formik, FormikState, useFormikContext } from 'formik';
import React from 'react';
import { colors } from 'Styles/colors';
import {
  SoloButtonStyledComponent,
  SoloCancelButton,
  SoloNegativeButton,
} from 'Styles/StyledComponents/button';
import * as yup from 'yup';
import YAML from 'yaml';
import { graphqlApi } from 'API/graphql';
import { ValidateResolverYamlRequest } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/graphql_pb';
import {
  GrpcResolver,
  Resolution,
  RESTResolver,
} from 'proto/github.com/solo-io/solo-apis/api/gloo/graphql.gloo/v1alpha1/graphql_pb';
import { useParams } from 'react-router';
import { Value } from 'google-protobuf/google/protobuf/struct_pb';
import { StringValue } from 'google-protobuf/google/protobuf/wrappers_pb';
import ConfirmationModal from 'Components/Common/ConfirmationModal';

export const EditorContainer = styled.div<{ editMode: boolean }>`
  .ace_cursor {
    opacity: ${props => (props.editMode ? 1 : 0)};
  }
  cursor: ${props => (props.editMode ? 'text' : 'default')};
`;

export const IconButton = styled.button`
  display: inline-flex;
  cursor: pointer;
  border: none;
  outline: none !important;
  background: transparent;
  justify-content: center;
  align-items: center;
  color: ${props => colors.lakeBlue};
  cursor: pointer;

  &:disabled {
    opacity: 0.3;
    pointer-events: none;
    cursor: default;
  }
`;

export type ResolverWizardFormProps = {
  resolverType: 'REST' | 'gRPC';
  upstream: string;
  resolverConfig: string;
};

const validationSchema = yup.object().shape({
  resolverType: yup.string().required('You need to specify a resolver type.'),
  upstream: yup.string().required('You need to specify an upstream.'),
  resolverConfig: yup
    .string()
    .required('You need to specify a resolver configuration.'),
});

type ResolverTypeSectionProps = { isEdit: boolean };
export let apiTypeOptions = [
  {
    displayValue: 'REST',
    value: 'REST',
    subHeader:
      'Integrate with upstream REST APIs and customize HTTP request and response mappings.',
  },
  {
    displayValue: 'gRPC',
    value: 'gRPC',
    subHeader: 'Integrate with upstream gRPC APIs based on a proto definition.',
  },
] as SoloFormRadioOption[];

const ResolverTypeSection = ({ isEdit }: ResolverTypeSectionProps) => {
  const formik = useFormikContext<ResolverWizardFormProps>();

  return (
    <div className='w-full h-full p-6 pb-0'>
      <div
        className={'flex items-center mb-6 text-lg font-medium text-gray-800'}>
        {isEdit ? 'Edit' : 'Configure'} Resolver{' '}
      </div>
      <div className='grid grid-cols-2 gap-4 '>
        <SoloFormRadio<ResolverWizardFormProps>
          name='resolverType'
          isUpdate={Boolean(isEdit)}
          title='Resolver Type'
          options={apiTypeOptions}
          titleAbove
        />
      </div>
    </div>
  );
};

type UpstreamSectionProps = { isEdit: boolean };

const UpstreamSection = ({ isEdit }: UpstreamSectionProps) => {
  const formik = useFormikContext<ResolverWizardFormProps>();
  const { data: upstreams, error: upstreamsError } = useListUpstreams();

  return (
    <div className='w-full h-full p-6 pb-0'>
      <div
        className={'flex items-center mb-6 text-lg font-medium text-gray-800'}>
        {isEdit ? 'Edit' : 'Configure'} Resolver{' '}
      </div>
      <div className='grid gap-4 '>
        <div className='mb-2 '>
          <label className='text-base font-medium '>Upstream</label>
          <div className='mt-2'>
            <SoloFormDropdown
              name='upstream'
              value={formik.values.upstream}
              defaultValue={formik.values.upstream}
              options={upstreams
                ?.map(upstream => {
                  return {
                    key: upstream.metadata?.uid!,
                    value: `${upstream.metadata?.name!}::${
                      upstream.metadata?.namespace
                    }`,
                    displayValue: upstream.metadata?.name!,
                  };
                })
                .sort((upstream1, upstream2) =>
                  upstream1.displayValue === upstream2.displayValue
                    ? 0
                    : (upstream1?.displayValue ?? upstream1.value) >
                      (upstream2?.displayValue ?? upstream2.value)
                    ? 1
                    : -1
                )}
            />
          </div>
        </div>
      </div>
    </div>
  );
};

type ResolverConfigSectionProps = {
  isEdit: boolean;
  resolverConfig: string;
};

let demoConfig = `restResolver:
  request:
    headers:
      - :method: 
      - :path:
    queryParams:
    body:
  response:
    resultRoot:
    setters:`;
const ResolverConfigSection = ({
  isEdit,
  resolverConfig,
}: ResolverConfigSectionProps) => {
  const { setFieldValue, values, dirty, errors } =
    useFormikContext<ResolverWizardFormProps>();
  const [isValid, setIsValid] = React.useState(false);
  const [errorModal, setErrorModal] = React.useState(false);
  const [errorMessage, setErrorMessage] = React.useState('');

  React.useEffect(() => {
    if (values.resolverConfig?.length > 3) {
      let config = YAML.parse(values.resolverConfig) as any;

      let configToDisplay: Partial<
        GrpcResolver.AsObject | RESTResolver.AsObject
      > = {};

      if (values.resolverType === 'REST') {
        //@ts-ignore
        configToDisplay = {
          //@ts-ignore
          ...((config.restResolver?.request ?? config?.request) && {
            request: {
              ...((config?.restResolver?.request ?? config?.request)?.body
                ?.stringValue && {
                body: (config?.restResolver?.request ?? config?.request)?.body
                  ?.stringValue,
              }),
              //@ts-ignore
              ...((config.restResolver?.request ?? config?.request)?.headersMap
                ?.length > 0 && {
                headers: Object.fromEntries(
                  //@ts-ignore
                  (config.restResolver?.request ?? config?.request)?.headersMap
                ),
              }),
              ...((config.restResolver?.request ?? config?.request)
                ?.queryParamsMap?.length > 0 && {
                queryParams: Object.fromEntries(
                  (config.restResolver?.request ?? config?.request)
                    ?.queryParamsMap
                ),
              }),
            },
          }),
          //@ts-ignore
          ...((config.restResolver?.response ?? config?.response) && {
            response: {
              ...((config.restResolver?.response ?? config?.response)
                ?.resultRoot && {
                resultRoot: '',
              }),
              ...((config.restResolver?.response ?? config?.response)
                ?.settersMap?.length > 0 && {
                setters: Object.fromEntries(
                  (config.restResolver?.response ?? config?.response)
                    ?.settersMap
                ),
              }),
            },
          }),
          ...((config.restResolver?.spanName ?? config?.spanName) && {
            spanName: config.restResolver?.spanName ?? config?.spanName,
          }),
        };
      } else {
        //@ts-ignore
        configToDisplay = {
          //@ts-ignore
          ...((config?.grpcResolver?.requestTransform ??
            config?.requestTransform) && {
            requestTransform: {
              ...((
                config?.grpcResolver?.requestTransform ??
                config?.requestTransform
              )?.methodName && {
                methodName: (
                  config?.grpcResolver?.requestTransform ??
                  config?.requestTransform
                )?.methodName,
              }),
              ...((
                config?.grpcResolver?.requestTransform ??
                config?.requestTransform
              )?.requestMetadataMap?.length > 0 && {
                requestMetadata: Object.fromEntries(
                  (
                    config?.grpcResolver?.requestTransform ??
                    config?.requestTransform
                  )?.requestMetadataMap
                ),
              }),
              ...((
                config?.grpcResolver?.requestTransform ??
                config?.requestTransform
              )?.serviceName && {
                serviceName: (
                  config?.grpcResolver?.requestTransform ??
                  config?.requestTransform
                )?.serviceName,
              }),
              ...((
                config?.grpcResolver?.requestTransform ??
                config?.requestTransform
              )?.outgoingMessageJson && {
                outgoingMessageJson: (
                  config?.grpcResolver?.requestTransform ??
                  config?.requestTransform
                )?.outgoingMessageJson.stringValue,
              }),
            },
          }),
          ...((config.grpcResolver?.spanName ?? config?.spanName) && {
            spanName: config.grpcResolver?.spanName ?? config?.spanName,
          }),
        };
      }

      setFieldValue('resolverConfig', YAML.stringify(configToDisplay));
      // setFieldValue('resolverConfig', 'request:');
    } else {
      setFieldValue('resolverConfig', demoConfig);
    }

    let res = `request:
    headers:
      - :method: GET
      - :path: /api/v1/products`;
  }, [
    setFieldValue,
    values.upstream,
    values.resolverConfig,
    values.resolverType,
  ]);

  const validateResolverSchema = async (resolver: string) => {
    setIsValid(!isValid);
    try {
      let res = await graphqlApi.validateResolverYaml({
        yaml: resolver,
        resolverType:
          values.resolverType === 'REST'
            ? ValidateResolverYamlRequest.ResolverType.REST_RESOLVER
            : ValidateResolverYamlRequest.ResolverType.REST_RESOLVER,
      });
      setIsValid(true);
      setErrorMessage('');
    } catch (err: any) {
      let [_, conversionError] = err?.message?.split(
        'failed to convert options YAML to JSON: yaml:'
      ) as [string, string];
      let [__, yamlError] = err?.message?.split(' invalid options YAML:') as [
        string,
        string
      ];
      if (conversionError) {
        setIsValid(false);
        setErrorMessage(`Error on ${conversionError}`);
      } else if (yamlError) {
        setIsValid(false);
        setErrorMessage(
          `Error: ${yamlError?.substring(yamlError.indexOf('):') + 2) ?? ''}`
        );
      }
    }
  };

  return (
    <div className='h-full p-6 pb-0 '>
      <div
        className={'flex items-center mb-2 text-lg font-medium text-gray-800'}>
        {isEdit ? 'Edit' : 'Configure'} Resolver{' '}
      </div>
      <div className=''>
        <div className='mb-2 '>
          <div>
            <EditorContainer editMode={true}>
              <div className=''>
                <div className='' style={{ height: 'min-content' }}>
                  {isValid ? (
                    <div
                      className={`${
                        isValid ? 'opacity-100' : 'opacity-0'
                      } h-10 text-center`}>
                      <div
                        style={{ backgroundColor: '#f2fef2' }}
                        className='p-2 text-green-400 border border-green-400 '>
                        <div className='font-medium '>Valid</div>
                      </div>
                    </div>
                  ) : (
                    <div
                      className={`${
                        errorMessage.length > 0 ? 'opacity-100' : '  opacity-0'
                      } h-10`}>
                      <div
                        style={{ backgroundColor: '#FEF2F2' }}
                        className='p-2 text-orange-400 border border-orange-400 '>
                        <div className='font-medium '>
                          {errorMessage?.split(',')[0]}
                        </div>
                        <ul className='pl-2 list-disc'>
                          {errorMessage?.split(',')[1]}
                        </ul>
                      </div>
                    </div>
                  )}
                </div>

                <div className='mt-2'></div>
              </div>
              <div className='flex flex-col w-full '>
                <>
                  <YamlEditor
                    mode='yaml'
                    theme='chrome'
                    name='resolverConfiguration'
                    style={{
                      width: '100%',
                      maxHeight: '36vh',
                      cursor: 'text',
                    }}
                    onChange={e => {
                      setFieldValue('resolverConfig', e);
                    }}
                    focus={true}
                    onInput={() => {
                      setIsValid(false);
                    }}
                    fontSize={16}
                    showPrintMargin={false}
                    showGutter={true}
                    highlightActiveLine={true}
                    value={values.resolverConfig}
                    readOnly={false}
                    setOptions={{
                      highlightGutterLine: true,
                      showGutter: true,
                      fontFamily: 'monospace',
                      enableBasicAutocompletion: true,
                      enableLiveAutocompletion: true,
                      showLineNumbers: true,
                      tabSize: 2,
                    }}
                  />
                  <div className='flex gap-3 mt-2'>
                    <SoloButtonStyledComponent
                      data-testid='save-route-options-changes-button '
                      disabled={!dirty || isValid}
                      onClick={() =>
                        validateResolverSchema(values.resolverConfig ?? '')
                      }>
                      Validate
                    </SoloButtonStyledComponent>

                    <SoloCancelButton
                      disabled={!dirty}
                      onClick={() => {
                        setFieldValue('resolverConfig', demoConfig);
                        setErrorMessage('');
                      }}>
                      Reset
                    </SoloCancelButton>
                  </div>
                </>
              </div>
            </EditorContainer>
          </div>
        </div>
      </div>
    </div>
  );
};

type ResolverWizardProps = {
  onClose: () => void;
  resolver?: Resolution.AsObject;
  resolverName?: string;
};

export const ResolverWizard: React.FC<ResolverWizardProps> = props => {
  const {
    graphqlSchemaName = '',
    graphqlSchemaNamespace = '',
    graphqlSchemaClusterName = '',
  } = useParams();

  const { data: graphqlSchema, mutate } = useGetGraphqlSchemaDetails({
    name: graphqlSchemaName,
    namespace: graphqlSchemaNamespace,
    clusterName: graphqlSchemaClusterName,
  });
  const [tabIndex, setTabIndex] = React.useState(0);
  const handleTabsChange = (index: number) => {
    setTabIndex(index);
  };
  const [isValid, setIsValid] = React.useState(false);
  const [isEdit, setIsEdit] = React.useState(Boolean(props.resolver));
  const [attemptUpdateSchema, setAttemptUpdateSchema] = React.useState(false);

  const submitResolverConfig = async (values: ResolverWizardFormProps) => {
    let { resolverConfig, resolverType, upstream } = values;

    /*
     `parsedResolverConfig` can be formatted in different ways: 
     - `restResolver.[request | response | spanName | ...]`....
     - `grpcResolver.[request | response | spanName | ...]`...
     - `[request | response | spanName | ...]`...
    */
    let parsedResolverConfig = YAML.parse(resolverConfig);
    let headersMap: [string, string][] = [];
    let queryParamsMap: [string, string][] = [];
    let settersMap: [string, string][] = [];
    let requestMetadataMap: [string, string][] = [];
    let serviceName = '';
    let methodName = '';
    let outgoingMessageJson;

    let resultRoot = parsedResolverConfig?.response
      ? parsedResolverConfig?.response?.resultRoot
      : parsedResolverConfig?.grpcResolver?.resultRoot;
    let spanName = parsedResolverConfig?.spanName
      ? parsedResolverConfig?.grpcResolver.spanName
      : parsedResolverConfig?.restResolver?.spanName ??
        parsedResolverConfig?.spanName;
    let grpcRequest = parsedResolverConfig?.requestTransform;
    let request = parsedResolverConfig?.request
      ? parsedResolverConfig?.request
      : parsedResolverConfig?.restResolver?.request;
    let response = parsedResolverConfig?.response
      ? parsedResolverConfig?.response
      : parsedResolverConfig?.restResolver?.response;

    if (resolverType === 'REST') {
      if (parsedResolverConfig?.restResolver) {
        headersMap =
          parsedResolverConfig?.restResolver?.request?.headers?.flatMap(
            (item: any) => Object.entries(item)
          );
        queryParamsMap =
          parsedResolverConfig?.restResolver?.request?.queryParams?.flatMap(
            (item: any) => Object.entries(item)
          );
        settersMap =
          parsedResolverConfig?.restResolver?.response?.settersMap?.flatMap(
            (item: any) => Object.entries(item)
          );
        resultRoot = parsedResolverConfig?.restResolver?.response?.resultRoot;

        spanName = parsedResolverConfig?.restResolver?.spanName;
      } else {
        headersMap = parsedResolverConfig?.request?.headers?.flatMap(
          (item: any) => Object.entries(item)
        );
        queryParamsMap = parsedResolverConfig?.request?.queryParams?.flatMap(
          (item: any) => Object.entries(item)
        );
        settersMap = parsedResolverConfig?.response?.settersMap?.flatMap(
          (item: any) => Object.entries(item)
        );
        resultRoot = parsedResolverConfig?.response?.resultRoot;
        spanName = parsedResolverConfig?.spanName;
      }
    } else {
      if (parsedResolverConfig?.grpcResolver) {
        requestMetadataMap =
          parsedResolverConfig?.grpcResolver?.requestTransform?.requestMetadataMap?.flatMap(
            (item: any) => Object.entries(item)
          );
        serviceName =
          parsedResolverConfig?.grpcResolver?.requestTransform?.serviceName;
        methodName =
          parsedResolverConfig?.grpcResolver?.requestTransform?.methodName;
        spanName = parsedResolverConfig?.grpcResolver?.spanName;
      } else {
        requestMetadataMap =
          parsedResolverConfig?.requestTransform?.requestMetadataMap?.flatMap(
            (item: any) => Object.entries(item)
          );
        serviceName = parsedResolverConfig?.requestTransform?.serviceName;
        methodName = parsedResolverConfig?.requestTransform?.methodName;
        spanName = parsedResolverConfig?.spanName;
      }
    }

    let [upstreamName, upstreamNamespace] = upstream.split('::');

    let res = `
    request:
      headers:
        - :method: GET
        - :path: /api/v1/products`;

    await graphqlApi.updateGraphqlSchemaResolver(
      {
        name: graphqlSchemaName,
        namespace: graphqlSchemaNamespace,
        clusterName: graphqlSchemaClusterName,
      },
      {
        upstreamRef: {
          name: upstreamName!,
          namespace: upstreamNamespace!,
        },
        //@ts-ignore
        ...(request && {
          request: {
            headersMap,
            queryParamsMap,
          },
        }),
        resolverName: props.resolverName!,
        //@ts-ignore
        ...(grpcRequest && {
          grpcRequest: {
            methodName,
            requestMetadataMap,
            serviceName,
            outgoingMessageJson: undefined,
          },
        }),
        resolverType,
        //@ts-ignore
        ...(response && { response: { resultRoot, settersMap } }),
        spanName,
      }
    );

    mutate();
    props.onClose();
  };
  const removeResolverConfig = async () => {
    await graphqlApi.updateGraphqlSchemaResolver(
      {
        name: graphqlSchemaName,
        namespace: graphqlSchemaNamespace,
        clusterName: graphqlSchemaClusterName,
      },
      {
        resolverName: props.resolverName!,
      },
      true
    );
    setTimeout(() => {
      mutate();
    }, 300);
    props.onClose();
  };
  const resolverTypeIsValid = (
    formik: FormikState<ResolverWizardFormProps>
  ) => {
    return !formik.errors.resolverType;
  };

  const upstreamIsValid = (formik: FormikState<ResolverWizardFormProps>) => {
    return !formik.errors.upstream;
  };

  const resolverConfigIsValid = (
    formik: FormikState<ResolverWizardFormProps>
  ) => {
    return !formik.errors.resolverConfig;
  };

  const formIsValid = (formik: FormikState<ResolverWizardFormProps>) =>
    resolverTypeIsValid(formik) &&
    upstreamIsValid(formik) &&
    resolverConfigIsValid(formik);

  React.useEffect(() => {
    setIsEdit(Boolean(props.resolver));
  }, [props.resolver]);

  const getInitialResolverConfig = (resolver?: typeof props.resolver) => {
    if (resolver?.restResolver) {
      return YAML.stringify(resolver);
    }
    return '';
  };

  return (
    <div className='h-[700px]'>
      <Formik<ResolverWizardFormProps>
        initialValues={{
          resolverType: 'REST',
          upstream: props.resolver?.restResolver?.upstreamRef?.name!
            ? `${props.resolver?.restResolver?.upstreamRef?.name!}::${props
                .resolver?.restResolver?.upstreamRef?.namespace!}`
            : props.resolver?.grpcResolver?.upstreamRef?.name!
            ? `${props.resolver?.grpcResolver?.upstreamRef?.name!}::${props
                .resolver?.grpcResolver?.upstreamRef?.namespace!}`
            : '',
          resolverConfig: getInitialResolverConfig(props?.resolver),
        }}
        enableReinitialize
        validationSchema={validationSchema}
        onSubmit={submitResolverConfig}>
        {formik => (
          <>
            <StyledModalTabs
              style={{ backgroundColor: colors.oceanBlue }}
              className='grid h-full rounded-lg grid-cols-[150px_1fr]'
              index={tabIndex}
              onChange={handleTabsChange}>
              <TabList className='flex flex-col mt-6'>
                <StyledModalTab
                  isCompleted={!!formik.values.resolverType?.length}>
                  Resolver Type
                </StyledModalTab>

                <StyledModalTab isCompleted={!!formik.values.upstream?.length}>
                  Upstream
                </StyledModalTab>
                <StyledModalTab
                  isCompleted={!!formik.values.resolverConfig?.length}>
                  Resolver Config
                </StyledModalTab>
              </TabList>
              <TabPanels className='bg-white rounded-r-lg'>
                <TabPanel className='relative flex flex-col justify-between h-full pb-4 focus:outline-none'>
                  <ResolverTypeSection isEdit={isEdit} />
                  <div className='ml-2'>
                    <SoloNegativeButton
                      onClick={() => {
                        setAttemptUpdateSchema(true);
                      }}>
                      Remove Configuration
                    </SoloNegativeButton>
                  </div>
                  <div className='flex items-center justify-between px-6 '>
                    <IconButton onClick={() => props.onClose()}>
                      Cancel
                    </IconButton>
                    <SoloButtonStyledComponent
                      onClick={() => setTabIndex(tabIndex + 1)}
                      disabled={!resolverTypeIsValid(formik)}>
                      Next Step
                    </SoloButtonStyledComponent>
                  </div>
                </TabPanel>

                <TabPanel className='relative flex flex-col justify-between h-full pb-4 focus:outline-none'>
                  <UpstreamSection isEdit={isEdit} />
                  <div className='flex items-center justify-between px-6 '>
                    <IconButton onClick={() => props.onClose()}>
                      Cancel
                    </IconButton>
                    <SoloButtonStyledComponent
                      onClick={() => setTabIndex(tabIndex + 1)}
                      disabled={!upstreamIsValid(formik)}>
                      Next Step
                    </SoloButtonStyledComponent>
                  </div>
                </TabPanel>
                <TabPanel className='relative flex flex-col justify-between h-full pb-4 focus:outline-none'>
                  {tabIndex === 2 && (
                    <ResolverConfigSection
                      isEdit={isEdit}
                      resolverConfig={formik.values.resolverConfig}
                    />
                  )}

                  <div className='flex items-center justify-between px-6 '>
                    <IconButton onClick={() => props.onClose()}>
                      Cancel
                    </IconButton>
                    <SoloButtonStyledComponent
                      onClick={formik.handleSubmit as any}
                      disabled={!formik.isValid || !formIsValid(formik)}>
                      Submit
                    </SoloButtonStyledComponent>
                  </div>
                </TabPanel>
              </TabPanels>
            </StyledModalTabs>
          </>
        )}
      </Formik>
      <ConfirmationModal
        visible={attemptUpdateSchema}
        confirmPrompt='delete this Resolver'
        confirmButtonText='Delete'
        goForIt={removeResolverConfig}
        cancel={() => {
          setAttemptUpdateSchema(false);
        }}
        isNegative
      />
    </div>
  );
};
