import { TabList, TabPanel, TabPanels } from '@reach/tabs';
import { graphqlConfigApi } from 'API/graphql';
import { useGetConsoleOptions, useGetGraphqlApiDetails } from 'API/hooks';
import ConfirmationModal from 'Components/Common/ConfirmationModal';
import { StyledModalTab, StyledModalTabs } from 'Components/Common/SoloModal';
import { Formik, FormikState } from 'formik';
import { FieldDefinitionNode } from 'graphql';
import { ClusterObjectRef } from 'proto/github.com/solo-io/skv2/api/core/v1/core_pb';
import { Resolution } from 'proto/github.com/solo-io/solo-apis/api/gloo/graphql.gloo/v1beta1/graphql_pb';
import { ValidateSchemaDefinitionRequest } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/graphql_pb';
import React, { useEffect, useMemo, useState } from 'react';
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
import * as yup from 'yup';
import { createResolverItem, getResolverFromConfig } from './converters';
import { GrpcProtoCheck } from './grpcProtoCheck/GrpcProtoCheck';
import { ResolverConfigSection } from './ResolverConfigSection';
import { getType, ResolverTypeSection } from './ResolverTypeSection';
import * as styles from './ResolverWizard.styles';
import { UpstreamSection } from './UpstreamSection';

export type ResolverWizardFormProps = {
  resolverType: 'REST' | 'gRPC';
  upstream: string;
  resolverConfig: string;
  listOfResolvers: [string, Resolution.AsObject][];
  protoFile?: string;
};

// --- VALIDATION --- //
const validationSchema = yup.object().shape({
  resolverType: yup.string().required('You need to specify a resolver type.'),
  upstream: yup.string().required('You need to specify an upstream.'),
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
  const [isConfirmingDelete, setIsConfirmingDelete] = useState(false);
  useEffect(() => {
    // Reset when field is unset (the wizard is hidden).
    if (field === null) {
      setIsConfirmingDelete(false);
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
      await graphqlConfigApi.updateGraphqlApiResolver(apiRef, resolverItem);
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
      setIsConfirmingDelete(false);
    }
  };

  return (
    <div
      data-testid='resolver-wizard'
      className='relative min-h-[600px] max-h-[800px] h-[85vh]'>
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
                {/* TODO:  Toggle here on gRPC. */}
                <StyledModalTab
                  isSelected={
                    formik.values.resolverType === 'gRPC' && tabIndex === 1
                  }
                  className={`${
                    formik.values.resolverType === 'gRPC' ? 'visible' : 'hidden'
                  }`}
                  data-testid='resolver-gprc-proto-tab'
                  isCompleted={!!formik.values.protoFile?.length}>
                  gRPC Toggle
                </StyledModalTab>
                <StyledModalTab
                  isSelected={
                    (formik.values.resolverType === 'gRPC' && tabIndex === 2) ||
                    (formik.values.resolverType === 'REST' && tabIndex === 1)
                  }
                  data-testid='upstream-tab'
                  isCompleted={!!formik.values.upstream?.length}>
                  Upstream
                </StyledModalTab>
                <StyledModalTab
                  isSelected={
                    (formik.values.resolverType === 'gRPC' && tabIndex === 3) ||
                    (formik.values.resolverType === 'REST' && tabIndex === 2)
                  }
                  data-testid='resolver-config-tab'
                  isCompleted={!!formik.values.resolverConfig?.length}>
                  Resolver Config
                </StyledModalTab>
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
                {/* --- STEP 1: API TYPE --- */}
                <TabPanel className='relative flex flex-col justify-between h-full pb-4 focus:outline-none'>
                  <ResolverTypeSection />
                  {!readonly && !isNewResolution && (
                    <div className='ml-2'>
                      <SoloNegativeButton
                        data-testid='remove-configuration-btn'
                        onClick={() => setIsConfirmingDelete(true)}>
                        Remove Configuration
                      </SoloNegativeButton>
                    </div>
                  )}
                  <div className='flex items-center justify-between px-6 '>
                    <styles.IconButton onClick={onClose}>
                      Cancel
                    </styles.IconButton>
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
                {/* Step 2 or none: Get the gRPC ProtoFile  */}
                <TabPanel
                  className={`
                  ${
                    formik.values.resolverType !== 'gRPC' ? 'hidden' : 'visible'
                  }
                  relative flex-grow flex flex-col justify-between pb-4 focus:outline-none`}>
                  <GrpcProtoCheck
                    setWarningMessage={(message: string) => {
                      setWarningMessage(message);
                    }}
                    warningMessage={warningMessage}
                  />
                  <div className='flex items-center justify-between px-6 '>
                    <styles.IconButton onClick={onClose}>
                      Cancel
                    </styles.IconButton>
                    {!readonly && (
                      <SoloButtonStyledComponent
                        data-testid='resolver-wizard-submit'
                        onClick={() => setTabIndex(tabIndex + 1)}>
                        Next Step
                      </SoloButtonStyledComponent>
                    )}
                  </div>
                </TabPanel>
                {/* --- STEP 2 | 3: UPSTREAM --- */}
                <TabPanel className='relative flex-grow flex flex-col justify-between pb-4 focus:outline-none'>
                  <UpstreamSection existingUpstreamId={existingUpstreamId} />
                  <div className='flex items-center justify-between px-6 '>
                    <styles.IconButton onClick={onClose}>
                      Cancel
                    </styles.IconButton>
                    <SoloButtonStyledComponent
                      onClick={() => setTabIndex(tabIndex + 1)}
                      disabled={!upstreamIsValid(formik)}>
                      Next Step
                    </SoloButtonStyledComponent>
                  </div>
                </TabPanel>
                {/* --- STEP 3 | 4: CONFIG --- */}
                <TabPanel className='relative flex-grow flex flex-col justify-between pb-4 focus:outline-none'>
                  <ResolverConfigSection warningMessage={warningMessage} />
                  <div className='flex items-center justify-between px-6 '>
                    <styles.IconButton onClick={onClose}>
                      Cancel
                    </styles.IconButton>
                    {!readonly && (
                      <SoloButtonStyledComponent
                        data-testid='resolver-wizard-submit'
                        onClick={formik.handleSubmit as any}
                        disabled={!formik.isValid || !formIsValid(formik)}>
                        Submit
                      </SoloButtonStyledComponent>
                    )}
                  </div>
                </TabPanel>
              </TabPanels>
            </StyledModalTabs>
          </>
        )}
      </Formik>
      <ConfirmationModal
        visible={isConfirmingDelete}
        confirmPrompt='delete this Resolver'
        confirmButtonText='Delete'
        confirmTestId='confirm-delete-resolver'
        cancelTestId='cancel-delete-resolver'
        goForIt={deleteResolverConfig}
        cancel={() => setIsConfirmingDelete(false)}
        isNegative
      />
    </div>
  );
};
