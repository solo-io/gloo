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
import { SoloButtonStyledComponent } from 'Styles/StyledComponents/button';
import * as yup from 'yup';
import YAML from 'yaml';

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
  const formik = useFormikContext<ResolverWizardProps>();
  const [isValid, setIsValid] = React.useState(false);

  React.useEffect(() => {
    setTimeout(() => {
      formik.setFieldValue('resolverConfig', resolverConfig);
    }, 300);
  }, []);

  return (
    <div className='w-full h-full p-6 pb-0'>
      <div
        className={'flex items-center mb-6 text-lg font-medium text-gray-800'}
      >
        {isEdit ? 'Edit' : 'Configure'} Resolver{' '}
      </div>
      <div className=''>
        <div className='mb-2 '>
          <div>
            <YamlEditor
              name='resolverConfig'
              title='Resolver Configuration'
              defaultValue={resolverConfig}
              onChange={e => {
                formik.setFieldValue('resolverConfig', e);
              }}
              onInput={() => {
                setIsValid(false);
              }}
              value={formik.values.resolverConfig ?? ''}
            />
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
  }, [!!props.resolver?.name]);

  const getInitialResolverConfig = (resolver?: typeof props.resolver) => {
    if (resolver?.restResolver) {
      return YAML.stringify(resolver);
    }
    return '';
  };
  return (
    <div className='h-[650px]'>
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
                  <div className='w-full h-full p-6 pb-0'>
                    <div
                      className={
                        'flex items-center mb-6 text-lg font-medium text-gray-800'
                      }
                    >
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
                  </div>

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
