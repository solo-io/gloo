import styled from '@emotion/styled/macro';
import { TabList, TabPanel, TabPanels } from '@reach/tabs';
import {
  useGetGraphqlApiDetails,
  useGetGraphqlApiYaml,
  useGetConsoleOptions,
} from 'API/hooks';
import { StyledModalTab, StyledModalTabs } from 'Components/Common/SoloModal';
import { Formik, FormikState } from 'formik';
import React from 'react';
import { colors } from 'Styles/colors';
import {
  SoloButtonStyledComponent,
  SoloNegativeButton,
} from 'Styles/StyledComponents/button';
import * as yup from 'yup';
import YAML from 'yaml';
import { graphqlConfigApi } from 'API/graphql';
import { Resolution } from 'proto/github.com/solo-io/solo-apis/api/gloo/graphql.gloo/v1alpha1/graphql_pb';
import { useParams } from 'react-router';
import ConfirmationModal from 'Components/Common/ConfirmationModal';
import { getType, ResolverTypeSection } from './ResolverTypeSection';
import { UpstreamSection } from './UpstreamSection';
import { ResolverConfigSection } from './ResolverConfigSection';
import protobuf from 'protobufjs';
import { toProto3JSON } from 'proto3-json-serializer';
import GQLJsonDescriptor from 'Components/Features/Graphql/data/graphql.json';
import { ValidateSchemaDefinitionRequest } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/graphql_pb';
import { createResolverItem, getResolverFromConfig } from './converters';

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
  listOfResolvers: [string, Resolution.AsObject][];
};

export const getUpstream = (resolver: Resolution.AsObject): string => {
  return `
  ${
    resolver?.restResolver?.upstreamRef?.name!
      ? `${resolver?.restResolver?.upstreamRef?.name!}::${resolver?.restResolver
          ?.upstreamRef?.namespace!}`
      : resolver?.grpcResolver?.upstreamRef?.name!
      ? `${resolver?.grpcResolver?.upstreamRef?.name!}::${resolver?.grpcResolver
          ?.upstreamRef?.namespace!}`
      : ''
  }`.trim();
};

export const getUpstreamFromMap = (
  resolutionsMap: Array<[string, Resolution.AsObject]>,
  resolverName: string
) => {
  const resolutionsMapItem = resolutionsMap?.find(
    ([rN]) => rN === resolverName
  )?.[1];
  if (resolutionsMapItem) {
    return getUpstream(resolutionsMapItem);
  }
  return '';
};

const validationSchema = yup.object().shape({
  resolverType: yup.string().required('You need to specify a resolver type.'),
  upstream: yup.string().required('You need to specify an upstream.'),
  resolverConfig: yup
    .string()
    .required('You need to specify a resolver configuration.'),
});

type ResolverWizardProps = {
  onClose: () => void;
  resolver?: Resolution.AsObject;
  resolverName?: string;
  hasDirective?: boolean;
  fieldWithDirective?: string;
  fieldWithoutDirective?: string;
  altFieldWithDirective?: string;
};

export const ResolverWizard: React.FC<ResolverWizardProps> = props => {
  let {
    hasDirective,
    fieldWithDirective,
    fieldWithoutDirective,
    altFieldWithDirective,
  } = props;
  const {
    graphqlApiName = '',
    graphqlApiNamespace = '',
    graphqlApiClusterName = '',
  } = useParams();

  const { data: graphqlApi, mutate } = useGetGraphqlApiDetails({
    name: graphqlApiName,
    namespace: graphqlApiNamespace,
    clusterName: graphqlApiClusterName,
  });

  const { readonly } = useGetConsoleOptions();

  const { mutate: mutateSchemaYaml } = useGetGraphqlApiYaml({
    name: graphqlApiName,
    namespace: graphqlApiNamespace,
    clusterName: graphqlApiClusterName,
  });

  const resolutionsMap =
    graphqlApi?.spec?.executableSchema?.executor?.local?.resolutionsMap ?? [];

  const listOfResolvers = resolutionsMap.filter(
    ([rName]: [rName: string, rObject: Resolution.AsObject]) => {
      return props.resolverName !== rName;
    }
  );

  const [tabIndex, setTabIndex] = React.useState(0);
  const [warningMessage, setWarningMessage] = React.useState('');
  const handleTabsChange = (index: number) => {
    setTabIndex(index);
  };
  const isEdit = Boolean(props.resolver);
  const [attemptUpdateSchema, setAttemptUpdateSchema] = React.useState(false);

  const submitResolverConfig = async (values: ResolverWizardFormProps) => {
    let { resolverConfig, resolverType, upstream } = values;
    /*
     `parsedResolverConfig` can be formatted in different ways:
     - `restResolver.[request | response | spanName | ...]`....
     - `grpcResolver.[request | response | spanName | ...]`...
     - `[request | response | spanName | ...]`...
    */

    const apiRef = {
      name: graphqlApiName,
      namespace: graphqlApiNamespace,
      clusterName: graphqlApiClusterName,
    };
    const extras = {
      hasDirective,
      fieldWithDirective,
      fieldWithoutDirective,
      altFieldWithDirective,
    };
    let resolverItem: any;
    try {
      resolverItem = createResolverItem(
        resolverConfig,
        resolverType,
        props.resolverName ?? '',
        upstream,
        extras
      );
    } catch (err: any) {
      setWarningMessage(err.message);
      return;
    }

    setWarningMessage('');
    let validationObject =
      new ValidateSchemaDefinitionRequest().toObject() as any;
    const spec = (
      await graphqlConfigApi.getGraphqlApiWithResolver(apiRef, resolverItem)
    ).toObject();
    validationObject = {
      ...validationObject,
      spec,
      apiRef,
      resolverItem,
    };
    await graphqlConfigApi
      .validateSchema(validationObject)
      .then(_res => {
        return graphqlConfigApi
          .updateGraphqlApiResolver(apiRef, resolverItem)
          .then(_res => {
            mutate();
            mutateSchemaYaml();
            props.onClose();
          })
          .catch(err => {
            setWarningMessage(err.message);
            console.error({ err });
          });
      })
      .catch(err => {
        setWarningMessage(err.message);
      });
  };
  const removeResolverConfig = async () => {
    await graphqlConfigApi.updateGraphqlApiResolver(
      {
        name: graphqlApiName,
        namespace: graphqlApiNamespace,
        clusterName: graphqlApiClusterName,
      },
      {
        resolverName: props.resolverName!,
        hasDirective,
        fieldWithDirective,
        fieldWithoutDirective,
        altFieldWithDirective,
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

  return (
    <div data-testid='resolver-wizard' className=' h-[800px]'>
      <Formik<ResolverWizardFormProps>
        initialValues={{
          resolverType: getType(props.resolver),
          upstream: getUpstreamFromMap(
            resolutionsMap,
            props.resolverName ?? ''
          ),
          resolverConfig: getResolverFromConfig(props.resolver),
          listOfResolvers,
        }}
        enableReinitialize
        validateOnMount={true}
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
                  {!readonly && (
                    <div className='ml-2'>
                      <SoloNegativeButton
                        onClick={() => {
                          setAttemptUpdateSchema(true);
                        }}>
                        Remove Configuration
                      </SoloNegativeButton>
                    </div>
                  )}
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
                  <UpstreamSection
                    isEdit={isEdit}
                    existingUpstream={
                      props.resolver ? getUpstream(props.resolver) : ''
                    }
                  />
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
                      warningMessage={warningMessage}
                      isEdit={isEdit}
                    />
                  )}
                  <div className='flex items-center justify-between px-6 '>
                    <IconButton onClick={() => props.onClose()}>
                      Cancel
                    </IconButton>
                    {!readonly && (
                      <SoloButtonStyledComponent
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
