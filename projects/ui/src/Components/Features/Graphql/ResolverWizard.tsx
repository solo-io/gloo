import styled from '@emotion/styled/macro';
import { TabList, TabPanel, TabPanels } from '@reach/tabs';
import { useListUpstreams } from 'API/hooks';
import { OptionType } from 'Components/Common/SoloDropdown';
import { SoloFormDropdown } from 'Components/Common/SoloFormComponents';
import { StyledModalTab, StyledModalTabs } from 'Components/Common/SoloModal';
import YamlDisplayer from 'Components/Common/YamlDisplayer';
import YamlEditor from 'Components/Common/YamlEditor';
import { Formik, FormikState, useFormikContext } from 'formik';
import React from 'react';
import { colors } from 'Styles/colors';
import {
  SoloButtonStyledComponent,
  SoloCancelButton,
} from 'Styles/StyledComponents/button';
import * as yup from 'yup';
import YAML from 'yaml';
import { graphqlApi } from 'API/graphql';
import { ValidateResolverYamlRequest } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/graphql_pb';
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

export type ResolverWizardProps = {
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

const ResolverTypeSection = ({ isEdit }: ResolverTypeSectionProps) => {
  const formik = useFormikContext<ResolverWizardProps>();

  return (
    <div className='w-full h-full p-6 pb-0'>
      <div
        className={'flex items-center mb-6 text-lg font-medium text-gray-800'}
      >
        {isEdit ? 'Edit' : 'Configure'} Resolver{' '}
      </div>
      <div className='grid grid-cols-2 gap-4 '>
        <div className='grid grid-cols-2 col-span-2 gap-2 mb-2'>
          <label className='text-base font-medium '>Resolver Type</label>
          <div className='col-span-2 mt-3 -space-y-px bg-white rounded-md'>
            <div
              onClick={() => formik.setFieldValue('resolverType', 'REST')}
              className={`relative flex p-3 border ${
                formik.values.resolverType === 'REST'
                  ? ' border-blue-300gloo bg-blue-150gloo z-10 '
                  : 'border-gray-200'
              } rounded-tl-md rounded-tr-md`}
            >
              <div className='flex items-center h-5'>
                <input
                  type='radio'
                  readOnly
                  className='w-4 h-4 border-gray-300 cursor-pointer text-blue-600gloo focus:ring-blue-600gloo'
                  checked={formik.values.resolverType === 'REST'}
                />
              </div>
              <label className='flex flex-col ml-3 cursor-pointer'>
                <span
                  className={`block text-sm font-medium ${
                    formik.values.resolverType === 'REST'
                      ? ' text-blue-700gloo'
                      : 'text-gray-900'
                  } `}
                >
                  REST
                </span>
                {/* TODO: add copy explaining things */}
                <span className='block text-sm text-blue-700gloo'>
                  Integrate with upstream REST APIs and customize HTTP request
                  and response mappings.
                </span>
              </label>
            </div>

            <div
              className={`relative flex p-3 border ${
                formik.values.resolverType === 'gRPC'
                  ? ' border-blue-300gloo bg-blue-150gloo z-10 '
                  : 'border-gray-200'
              }rounded-bl-md rounded-br-md`}
              onClick={() => formik.setFieldValue('resolverType', 'gRPC')}
            >
              <div className='flex items-center h-5'>
                <input
                  type='radio'
                  readOnly
                  checked={formik.values.resolverType === 'gRPC'}
                  className='w-4 h-4 border-gray-300 cursor-pointer text-blue-600gloo focus:ring-blue-600gloo'
                />
              </div>
              <label
                htmlFor='settings-option-1'
                className='flex flex-col ml-3 cursor-pointer'
              >
                <span
                  className={`block text-sm font-medium ${
                    formik.values.resolverType === 'gRPC'
                      ? ' text-blue-700gloo'
                      : 'text-gray-900'
                  } `}
                >
                  gRPC
                </span>
                <span className='block text-sm text-blue-700gloo '>
                  Integrate with upstream gRPC APIs based on a proto definition.
                </span>
              </label>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

type UpstreamSectionProps = { isEdit: boolean };

const UpstreamSection = ({ isEdit }: UpstreamSectionProps) => {
  const formik = useFormikContext<ResolverWizardProps>();
  const { data: upstreams, error: upstreamsError } = useListUpstreams();

  return (
    <div className='w-full h-full p-6 pb-0'>
      <div
        className={'flex items-center mb-6 text-lg font-medium text-gray-800'}
      >
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
                    value: upstream.metadata?.name!,
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

const ResolverConfigSection = ({
  isEdit,
  resolverConfig,
}: ResolverConfigSectionProps) => {
  const { setFieldValue, values, dirty } =
    useFormikContext<ResolverWizardProps>();
  const [isValid, setIsValid] = React.useState(false);
  const [errorModal, setErrorModal] = React.useState(false);
  const [errorMessage, setErrorMessage] = React.useState('');
  React.useEffect(() => {
    setTimeout(() => {
      setFieldValue('resolverConfig', resolverConfig);
    }, 300);
    /* eslint-disable-next-line react-hooks/exhaustive-deps */
  }, []);

  const validateResolverSchema = (resolver: string) => {
    setIsValid(!isValid);
    try {
      graphqlApi
        .validateResolverYaml({
          yaml: resolver,
          resolverType:
            values.resolverType === 'REST'
              ? ValidateResolverYamlRequest.ResolverType.REST_RESOLVER
              : ValidateResolverYamlRequest.ResolverType.REST_RESOLVER,
        })
        .then(res => {
          console.log('res', res);
          setIsValid(true);
          setErrorMessage('');
        })
        .catch(err => {
          let [_, conversionError] = err.message?.split(
            'failed to convert options YAML to JSON: yaml:'
          ) as [string, string];
          let [__, yamlError] = err.message?.split(
            ' invalid options YAML:'
          ) as [string, string];
          if (conversionError) {
            setIsValid(false);
            setErrorMessage(`Error on ${conversionError}`);
          } else if (yamlError) {
            setIsValid(false);
            setErrorMessage(
              `Error: ${
                yamlError?.substring(yamlError.indexOf('):') + 2) ?? ''
              }`
            );
          }
        });
    } catch (error) {
      console.log(error);
    }
  };

  return (
    <div className='h-full p-6 pb-0 '>
      <div
        className={'flex items-center mb-2 text-lg font-medium text-gray-800'}
      >
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
                      } h-10 text-center`}
                    >
                      <div
                        style={{ backgroundColor: '#f2fef2' }}
                        className='p-2 text-green-400 border border-green-400 '
                      >
                        <div className='font-medium '>Valid</div>
                      </div>
                    </div>
                  ) : (
                    <div
                      className={`${
                        errorMessage.length > 0 ? 'opacity-100' : '  opacity-0'
                      } h-10`}
                    >
                      <div
                        style={{ backgroundColor: '#FEF2F2' }}
                        className='p-2 text-orange-400 border border-orange-400 '
                      >
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
                      width: '24vw',
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
                    fontSize={14}
                    showPrintMargin={false}
                    showGutter={true}
                    highlightActiveLine={true}
                    value={values.resolverConfig ?? ''}
                    readOnly={false}
                    setOptions={{
                      highlightGutterLine: true,
                      showGutter: true,
                      fontSize: 16,
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
                      }
                    >
                      Validate
                    </SoloButtonStyledComponent>

                    <SoloCancelButton
                      disabled={!dirty}
                      onClick={() => {
                        setFieldValue('resolverConfig', '');
                        setErrorMessage('');
                      }}
                    >
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

let res = {
  name: 'author',
  restResolver: {
    request: {
      headers: {
        ':method': 'GET',
        ':path': '/details/{$parent.id}',
      },
    },
    response: {
      resultRoot: 'author',
    },
    upstreamRef: {
      name: 'default-details-9080',
      namespace: 'gloo-system',
    },
  },
};
type ResolverWizardFormProps = {
  onClose: () => void;
  resolver?: typeof res;
};

export const ResolverWizard: React.FC<ResolverWizardFormProps> = props => {
  const [tabIndex, setTabIndex] = React.useState(0);
  const handleTabsChange = (index: number) => {
    setTabIndex(index);
  };
  const [isValid, setIsValid] = React.useState(false);

  const [isEdit, setIsEdit] = React.useState(Boolean(props.resolver));

  const submitResolverConfig = async (values: ResolverWizardProps) => {
    // TODO
  };

  const resolverTypeIsValid = (formik: FormikState<ResolverWizardProps>) => {
    return !formik.errors.resolverType;
  };

  const upstreamIsValid = (formik: FormikState<ResolverWizardProps>) => {
    return !formik.errors.upstream;
  };

  const resolverConfigIsValid = (formik: FormikState<ResolverWizardProps>) => {
    return !formik.errors.resolverConfig;
  };

  const formIsValid = (formik: FormikState<ResolverWizardProps>) =>
    resolverTypeIsValid(formik) &&
    upstreamIsValid(formik) &&
    resolverConfigIsValid(formik);

  React.useEffect(() => {
    setIsEdit(Boolean(props.resolver?.name));
    /* eslint-disable-next-line react-hooks/exhaustive-deps */
  }, [!!props.resolver?.name]);

  const getInitialResolverConfig = (resolver?: typeof props.resolver) => {
    if (resolver?.restResolver) {
      return YAML.stringify(resolver);
    }
    return '';
  };
  return (
    <div className='h-[700px]'>
      <Formik<ResolverWizardProps>
        initialValues={{
          resolverType: 'REST',
          upstream: props.resolver?.restResolver?.upstreamRef?.name! ?? '',
          resolverConfig: getInitialResolverConfig(props?.resolver),
        }}
        enableReinitialize
        validationSchema={validationSchema}
        onSubmit={submitResolverConfig}
      >
        {formik => (
          <>
            <StyledModalTabs
              style={{ backgroundColor: colors.oceanBlue }}
              className='grid h-full rounded-lg grid-cols-[150px_1fr]'
              index={tabIndex}
              onChange={handleTabsChange}
            >
              <TabList className='flex flex-col mt-6'>
                <StyledModalTab
                  isCompleted={!!formik.values.resolverType?.length}
                >
                  Resolver Type
                </StyledModalTab>

                <StyledModalTab isCompleted={!!formik.values.upstream?.length}>
                  Upstream
                </StyledModalTab>
                <StyledModalTab
                  isCompleted={!!formik.values.resolverConfig?.length}
                >
                  Resolver Config
                </StyledModalTab>
              </TabList>
              <TabPanels className='bg-white rounded-r-lg'>
                <TabPanel className='relative flex flex-col justify-between h-full pb-4 focus:outline-none'>
                  <ResolverTypeSection isEdit={false} />
                  <div className='flex items-center justify-between px-6 '>
                    <IconButton onClick={() => props.onClose()}>
                      Cancel
                    </IconButton>
                    <SoloButtonStyledComponent
                      onClick={() => setTabIndex(tabIndex + 1)}
                      disabled={!resolverTypeIsValid(formik)}
                    >
                      Next Step
                    </SoloButtonStyledComponent>
                  </div>
                </TabPanel>

                <TabPanel className='relative flex flex-col justify-between h-full pb-4 focus:outline-none'>
                  <UpstreamSection isEdit={false} />
                  <div className='flex items-center justify-between px-6 '>
                    <IconButton onClick={() => props.onClose()}>
                      Cancel
                    </IconButton>
                    <SoloButtonStyledComponent
                      onClick={() => setTabIndex(tabIndex + 1)}
                      disabled={!upstreamIsValid(formik)}
                    >
                      Next Step
                    </SoloButtonStyledComponent>
                  </div>
                </TabPanel>
                <TabPanel className='relative flex flex-col justify-between h-full pb-4 focus:outline-none'>
                  {/* <div className='w-full h-full p-6 pb-0'>
                    <div
                      className={
                        'flex items-center mb-6 text-lg font-medium text-gray-800'
                      }>
                      Resolver{' '}
                    </div>
                    <div className=''>
                      <div className='mb-2 '>
                        <div>
                          <YamlDisplayer
                            contentString={formik.values.resolverConfig}
                            copyable={false}
                          />
                        </div>
                      </div>
                    </div>
                  </div> */}
                  <ResolverConfigSection
                    isEdit
                    resolverConfig={formik.values.resolverConfig}
                  />

                  <div className='flex items-center justify-between px-6 '>
                    <IconButton onClick={() => props.onClose()}>
                      Cancel
                    </IconButton>
                    <SoloButtonStyledComponent
                      onClick={formik.handleSubmit as any}
                      disabled={!formik.isValid || !formIsValid(formik)}
                    >
                      Submit
                    </SoloButtonStyledComponent>
                  </div>
                </TabPanel>
              </TabPanels>
            </StyledModalTabs>
          </>
        )}
      </Formik>
    </div>
  );
};
