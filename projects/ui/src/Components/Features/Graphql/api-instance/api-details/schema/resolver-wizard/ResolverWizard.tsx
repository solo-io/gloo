import { TabList, TabPanel, TabPanels } from '@reach/tabs';
import { graphqlConfigApi } from 'API/graphql';
import { useGetConsoleOptions, useGetGraphqlApiDetails } from 'API/hooks';
import { StyledModalTab, StyledModalTabs } from 'Components/Common/SoloModal';
import { useConfirm } from 'Components/Context/ConfirmModalContext';
import { Formik, FormikState } from 'formik';
import { FieldDefinitionNode } from 'graphql';
import { ClusterObjectRef } from 'proto/github.com/solo-io/skv2/api/core/v1/core_pb';
import { Resolution } from 'proto/github.com/solo-io/solo-apis/api/gloo/graphql.gloo/v1beta1/graphql_pb';
import { ValidateSchemaDefinitionRequest } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/graphql_pb';
import React, { useEffect, useMemo } from 'react';
import toast from 'react-hot-toast';
import { colors } from 'Styles/colors';
import {
  SoloButtonStyledComponent,
  SoloNegativeButton,
} from 'Styles/StyledComponents/button';
import {
  getFieldReturnType,
  getResolution,
  getResolveDirectiveName,
  getUpstreamId,
  getUpstreamRef,
} from 'utils/graphql-helpers';
import { hotToastError } from 'utils/hooks';
import * as yup from 'yup';
import { createResolverItem, getResolverFromConfig } from './converters';
import { GrpcProtoCheck } from './grpcProtoCheck/GrpcProtoCheck';
import { ResolverConfigSection } from './ResolverConfigSection';
import { getType, ResolverTypeSection } from './ResolverTypeSection';
import * as ResolverWizardStyles from './ResolverWizard.styles';
import { UpstreamSection } from './UpstreamSection';

export type ResolverType = 'gRPC' | 'REST' | 'Mock';

export type ResolverWizardFormProps = {
  resolverType: ResolverType;
  upstream: string;
  resolverConfig: string;
  listOfResolvers: [string, Resolution.AsObject][];
  protoFile?: string;
};

// --- VALIDATION --- //
const validationSchema = yup.object().shape({
  resolverType: yup.string().required('You need to specify a resolver type.'),
  upstream: yup
    .string()
    .test(
      'is-upstream-valid',
      'You need to specify an upstream.',
      (item, context) => {
        // Upstream is undefined for Mock resolvers.
        const values = context.parent as ResolverWizardFormProps;
        if (values.resolverType === 'Mock') return true;
        return !!item;
      }
    ),
  resolverConfig: yup
    .string()
    .required('You need to specify a resolver configuration.'),
});
const resolverTypeIsValid = (formik: FormikState<ResolverWizardFormProps>) =>
  !formik.errors.resolverType;
const upstreamIsValid = (formik: FormikState<ResolverWizardFormProps>) =>
  !formik.errors.upstream;
const resolverConfigIsValid = (formik: FormikState<ResolverWizardFormProps>) =>
  !formik.errors.resolverConfig;
const protoFileIsValid = (formik: FormikState<ResolverWizardFormProps>) =>
  !formik.errors.protoFile;
const formIsValid = (formik: FormikState<ResolverWizardFormProps>) =>
  resolverTypeIsValid(formik) &&
  upstreamIsValid(formik) &&
  resolverConfigIsValid(formik) &&
  protoFileIsValid(formik);

//
// --- COMPONENT --- //
//
export const ResolverWizard: React.FC<{
  apiRef: ClusterObjectRef.AsObject;
  field: FieldDefinitionNode | null;
  objectType: string;
  onClose: () => void;
}> = ({ apiRef, field, objectType, onClose }) => {
  const { data: graphqlApi } = useGetGraphqlApiDetails(apiRef);
  const { readonly } = useGetConsoleOptions();
  const confirm = useConfirm();

  // --- STATE (FIELD, RESOLVER) --- //
  const fieldName = field?.name.value ?? '';
  const fieldReturnType = useMemo(
    () => getFieldReturnType(field).fullType,
    [field]
  );
  const resolution = useMemo(
    () => getResolution(graphqlApi, getResolveDirectiveName(field)),
    [field]
  );
  const isNewResolution = useMemo(() => !resolution, [resolution]);
  const resolutionsMap =
    graphqlApi?.spec?.executableSchema?.executor?.local?.resolutionsMap ?? [];
  const listOfResolvers = resolutionsMap.filter(
    ([rName]: [rName: string, rObject: Resolution.AsObject]) => {
      return fieldName !== rName;
    }
  );
  const existingUpstreamId = useMemo(
    () => getUpstreamId(getUpstreamRef(resolution)),
    [resolution]
  );
  // --- STATE (TAB, WARNING, CONFIRM) --- //
  const [tabIndex, setTabIndex] = React.useState(0);
  const [warningMessage, setWarningMessage] = React.useState('');
  useEffect(() => {
    // Reset when field is unset (the wizard is hidden).
    if (field === null) {
      setWarningMessage('');
    }
  }, [field]);

  //
  // --- ADD + UPDATE --- //
  //
  const submitResolverConfig = async (values: ResolverWizardFormProps) => {
    try {
      if (!field) throw new Error('Field does not exist.');
      const { resolverConfig, resolverType, upstream } = values;
      //
      // Create ResolverItem (the parameter for validation + updates).
      const resolverItem = createResolverItem(
        resolverConfig,
        resolverType,
        field,
        upstream,
        {
          isNewResolution,
          fieldReturnType,
          objectType,
        }
      );
      const spec = (
        await graphqlConfigApi.getGraphqlApiWithResolver(apiRef, resolverItem)
      ).toObject();
      //
      // Validate schema.
      await graphqlConfigApi.validateSchema({
        ...new ValidateSchemaDefinitionRequest().toObject(),
        spec,
        apiRef,
        resolverItem,
      });
      //
      // Make the update and close the modal.
      await toast.promise(
        graphqlConfigApi.updateGraphqlApiResolver(apiRef, resolverItem),
        {
          loading: 'Updating API..',
          success: 'API updated!',
          error: hotToastError,
        }
      );
      onClose();
    } catch (err: any) {
      setWarningMessage(err?.message ?? err);
    }
  };

  //
  // --- DELETE --- //
  //
  const deleteResolverConfig = async () => {
    try {
      if (!field) throw new Error('Field does not exist.');
      await graphqlConfigApi.updateGraphqlApiResolver(
        apiRef,
        {
          field,
          objectType,
          isNewResolution,
          fieldReturnType,
        },
        true
      );
      onClose();
    } catch (err: any) {
      setWarningMessage(err.message ?? err);
    }
  };

  return (
    <div
      data-testid='resolver-wizard'
      className='relative min-h-[700px] max-h-[800px] h-[75vh]'>
      <Formik<ResolverWizardFormProps>
        initialValues={{
          resolverType: getType(resolution),
          upstream: existingUpstreamId,
          resolverConfig: getResolverFromConfig(resolution),
          listOfResolvers,
          protoFile: '',
        }}
        enableReinitialize
        validateOnMount={true}
        validationSchema={validationSchema}
        onSubmit={submitResolverConfig}>
        {formik => (
          <>
            <StyledModalTabs
              style={{ backgroundColor: colors.oceanBlue }}
              className='grid rounded-lg grid-cols-[150px_1fr] absolute top-0 left-0 w-full h-full'
              index={tabIndex}
              onChange={setTabIndex}>
              <TabList className='flex flex-col mt-6'>
                {/* --- SIDEBAR --- */}
                <StyledModalTab
                  isSelected={tabIndex === 0}
                  data-testid='resolver-type-tab'
                  isCompleted={!!formik.values.resolverType?.length}>
                  Resolver Type
                </StyledModalTab>
                {
                  // ==============
                  // SIDEBAR (gRPC)
                  // ==============
                  formik.values.resolverType === 'gRPC' && (
                    <>
                      <StyledModalTab
                        isSelected={tabIndex === 1}
                        data-testid='resolver-gprc-proto-tab'
                        isCompleted={!!formik.values.protoFile?.length}>
                        gRPC Toggle
                      </StyledModalTab>
                      <StyledModalTab
                        isSelected={tabIndex === 2}
                        data-testid='upstream-tab'
                        isCompleted={!!formik.values.upstream?.length}>
                        Upstream
                      </StyledModalTab>
                      <StyledModalTab
                        isSelected={tabIndex === 3}
                        data-testid='resolver-config-tab'
                        isCompleted={!!formik.values.resolverConfig?.length}>
                        Resolver Config
                      </StyledModalTab>
                    </>
                  )
                }
                {
                  // ==============
                  // SIDEBAR (REST)
                  // ==============
                  formik.values.resolverType === 'REST' && (
                    <>
                      <StyledModalTab
                        isSelected={tabIndex === 1}
                        data-testid='upstream-tab'
                        isCompleted={!!formik.values.upstream?.length}>
                        Upstream
                      </StyledModalTab>
                      <StyledModalTab
                        isSelected={tabIndex === 2}
                        data-testid='resolver-config-tab'
                        isCompleted={!!formik.values.resolverConfig?.length}>
                        Resolver Config
                      </StyledModalTab>
                    </>
                  )
                }
                {
                  // ==============
                  // SIDEBAR (Mock)
                  // ==============
                  formik.values.resolverType === 'Mock' && (
                    <>
                      <StyledModalTab
                        isSelected={tabIndex === 1}
                        data-testid='resolver-config-tab'
                        isCompleted={!!formik.values.resolverConfig?.length}>
                        Resolver Config
                      </StyledModalTab>
                    </>
                  )
                }
              </TabList>

              <TabPanels className='bg-white rounded-r-lg flex flex-col h-full'>
                <div
                  className={
                    'flex items-center mb-6 text-lg font-medium text-gray-800 px-6 pt-6'
                  }>
                  {isNewResolution ? 'Configuring' : 'Editing'}&nbsp;Resolver
                  for:&nbsp;
                  <b>{fieldName}</b>
                </div>
                {/* --- STEP 1: RESOLVED API TYPE --- */}
                <TabPanel className='relative flex flex-col justify-between h-full pb-4 focus:outline-none'>
                  <ResolverTypeSection
                    setWarningMessage={message => setWarningMessage(message)}
                  />
                  {!readonly && !isNewResolution && (
                    <div className='ml-2'>
                      <SoloNegativeButton
                        data-testid='remove-configuration-btn'
                        onClick={() =>
                          confirm({
                            confirmButtonText: 'Delete',
                            confirmPrompt: 'delete this Resolver',
                            isNegative: true,
                          }).then(() =>
                            toast.promise(deleteResolverConfig(), {
                              loading: 'Deleting Resolver Config',
                              success: 'Resolver deleted!',
                              error: hotToastError,
                            })
                          )
                        }>
                        Remove Configuration
                      </SoloNegativeButton>
                    </div>
                  )}
                  <div className='flex items-center justify-between px-6 '>
                    <ResolverWizardStyles.IconButton onClick={onClose}>
                      Cancel
                    </ResolverWizardStyles.IconButton>
                    <SoloButtonStyledComponent
                      onClick={() => {
                        if (formik.values.resolverType === 'gRPC') {
                          setTabIndex(tabIndex + 1);
                        } else {
                          setTabIndex(tabIndex + 2);
                        }
                      }}
                      disabled={!resolverTypeIsValid(formik)}>
                      Next Step
                    </SoloButtonStyledComponent>
                  </div>
                </TabPanel>

                {
                  // ==============
                  // CONTENT (gRPC)
                  // ==============
                  formik.values.resolverType === 'gRPC' && (
                    <>
                      {/* STEP 2: Get the gRPC ProtoFile  */}
                      <TabPanel
                        className={`relative flex-grow flex flex-col justify-between pb-4 focus:outline-none`}>
                        <GrpcProtoCheck
                          setWarningMessage={(message: string) => {
                            setWarningMessage(message);
                          }}
                          warningMessage={warningMessage}
                        />
                        <div className='flex items-center justify-between px-6 '>
                          <ResolverWizardStyles.IconButton onClick={onClose}>
                            Cancel
                          </ResolverWizardStyles.IconButton>
                          <SoloButtonStyledComponent
                            onClick={() => setTabIndex(tabIndex + 1)}>
                            Next Step
                          </SoloButtonStyledComponent>
                        </div>
                      </TabPanel>
                      {/* --- STEP 3: UPSTREAM --- */}
                      <TabPanel className='relative flex-grow flex flex-col justify-between pb-4 focus:outline-none'>
                        <UpstreamSection
                          setWarningMessage={message =>
                            setWarningMessage(message)
                          }
                          onCancel={onClose}
                          nextButtonDisabled={!upstreamIsValid(formik)}
                          onNextClicked={() => setTabIndex(tabIndex + 1)}
                          existingUpstreamId={existingUpstreamId}
                        />
                      </TabPanel>
                      {/* --- STEP 4: CONFIG --- */}
                      <TabPanel className='relative flex-grow flex flex-col justify-between pb-4 focus:outline-none'>
                        <ResolverConfigSection
                          onCancel={onClose}
                          submitDisabled={
                            !formik.isValid || !formIsValid(formik)
                          }
                          warningMessage={warningMessage}
                        />
                      </TabPanel>
                    </>
                  )
                }

                {
                  // ==============
                  // CONTENT (REST)
                  // ==============
                  formik.values.resolverType === 'REST' && (
                    <>
                      {/* --- STEP 2: UPSTREAM --- */}
                      <TabPanel className='relative flex-grow flex flex-col justify-between pb-4 focus:outline-none'>
                        <UpstreamSection
                          setWarningMessage={message =>
                            setWarningMessage(message)
                          }
                          onCancel={onClose}
                          nextButtonDisabled={!upstreamIsValid(formik)}
                          onNextClicked={() => setTabIndex(tabIndex + 1)}
                          existingUpstreamId={existingUpstreamId}
                        />
                      </TabPanel>
                      {/* --- STEP 3: CONFIG --- */}
                      <TabPanel className='relative flex-grow flex flex-col justify-between pb-4 focus:outline-none'>
                        <ResolverConfigSection
                          onCancel={onClose}
                          submitDisabled={
                            !formik.isValid || !formIsValid(formik)
                          }
                          warningMessage={warningMessage}
                        />
                      </TabPanel>
                    </>
                  )
                }

                {
                  // ==============
                  // CONTENT (Mock)
                  // ==============
                  formik.values.resolverType === 'Mock' && (
                    <>
                      {/* --- STEP 2: CONFIG --- */}
                      <TabPanel className='relative flex-grow flex flex-col justify-between pb-4 focus:outline-none'>
                        <ResolverConfigSection
                          onCancel={onClose}
                          submitDisabled={
                            !formik.isValid || !formIsValid(formik)
                          }
                          warningMessage={warningMessage}
                        />
                      </TabPanel>
                    </>
                  )
                }
              </TabPanels>
            </StyledModalTabs>
          </>
        )}
      </Formik>
    </div>
  );
};
