import styled from '@emotion/styled/macro';
import { TabList, TabPanel, TabPanels } from '@reach/tabs';
import { graphqlConfigApi } from 'API/graphql';
import {
  useGetConsoleOptions,
  useGetGraphqlApiDetails,
  useGetGraphqlApiYaml,
} from 'API/hooks';
import ConfirmationModal from 'Components/Common/ConfirmationModal';
import { StyledModalTab, StyledModalTabs } from 'Components/Common/SoloModal';
import { Formik, FormikState } from 'formik';
import { Kind, ObjectTypeDefinitionNode } from 'graphql';
import { Resolution } from 'proto/github.com/solo-io/solo-apis/api/gloo/graphql.gloo/v1alpha1/graphql_pb';
import { ValidateSchemaDefinitionRequest } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/graphql_pb';
import React, { useEffect, useMemo, useState } from 'react';
import { useParams } from 'react-router';
import { colors } from 'Styles/colors';
import {
  SoloButtonStyledComponent,
  SoloNegativeButton,
} from 'Styles/StyledComponents/button';
import { supportedDefinitionTypes } from 'utils/graphql-helpers';
import * as yup from 'yup';
import { getFieldTypeParts } from '../ExeGqlObjectDefinition';
import { createResolverItem, getResolverFromConfig } from './converters';
import { ResolverConfigSection } from './ResolverConfigSection';
import { getType, ResolverTypeSection } from './ResolverTypeSection';
import { UpstreamSection } from './UpstreamSection';

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

const validationSchema = yup.object().shape({
  resolverType: yup.string().required('You need to specify a resolver type.'),
  upstream: yup.string().required('You need to specify an upstream.'),
  resolverConfig: yup
    .string()
    .required('You need to specify a resolver configuration.'),
});

export const ResolverWizard: React.FC<{
  onClose: () => void;
  resolverName: string;
  objectType: string;
  schemaDefinitions: supportedDefinitionTypes[];
}> = props => {
  const { resolverName, objectType, schemaDefinitions } = props;
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

  const [resolver, setResolver] = useState<Resolution.AsObject>();
  const isNewResolver = useMemo(() => !!resolver, [resolver]);
  const [fieldReturnType, setFieldReturnType] = useState('');
  useEffect(() => {
    if (!resolverName || !objectType) return;
    //
    // Get the current resolver from the schema.
    let [_currentResolverName, currentResolver] =
      graphqlApi?.spec?.executableSchema?.executor?.local?.resolutionsMap.find(
        ([rName, _resolver]) => rName.includes(resolverName)
      ) ?? ['', undefined];
    //
    // Find the definition and field for the selected resolver.
    const definition = schemaDefinitions.find(
      d => d.kind === Kind.OBJECT_TYPE_DEFINITION && d.name.value === objectType
    ) as ObjectTypeDefinitionNode | undefined;
    if (definition === undefined) return;
    const field = definition.fields?.find(f => f.name.value === resolverName);
    if (field === undefined) return;
    //
    // Find the base field type (this could be a nested list).
    const [typePrefix, baseType, typeSuffix] = getFieldTypeParts(field);
    let newFieldReturnType = typePrefix + baseType + typeSuffix;
    //
    // Set state.
    setFieldReturnType(newFieldReturnType);
    setResolver(currentResolver);
  }, [schemaDefinitions, resolverName]);

  const resolutionsMap =
    graphqlApi?.spec?.executableSchema?.executor?.local?.resolutionsMap ?? [];

  const listOfResolvers = resolutionsMap.filter(
    ([rName]: [rName: string, rObject: Resolution.AsObject]) => {
      return props.resolverName !== rName;
    }
  );

  const getUpstreamFromMap = () => {
    // --------------------------------- //
    //
    // TODO: refactor this with the logic in the graphql.ts update function.
    //
    // Find the definition and field for the resolver to update.
    const definition = schemaDefinitions.find(
      (d: any) =>
        d.kind === Kind.OBJECT_TYPE_DEFINITION && d.name.value === objectType
    ) as ObjectTypeDefinitionNode | undefined;
    if (definition === undefined) return '';
    const resolverField = definition.fields?.find(
      f => f.name.value === resolverName
    );
    if (resolverField === undefined) return '';
    //
    // Try to get the '@resolve(...)' directive.
    // This is how we can check if it existed previously.
    const resolveDirective = resolverField.directives?.find(
      d => d.kind === Kind.DIRECTIVE && d.name.value === 'resolve'
    );
    if (!resolveDirective) return '';
    //
    // Get the resolver directives 'name' argument.
    // '@resolve(name: "...")'
    const resolverDirectiveArg = resolveDirective.arguments?.find(
      a => a.name.value === 'name'
    );
    if (
      !resolverDirectiveArg ||
      resolverDirectiveArg.value.kind !== Kind.STRING
    )
      return '';
    const resolverDirectiveName = resolverDirectiveArg.value.value;
    //
    // --------------------------------- //
    const resolutionsMapItem = resolutionsMap?.find(
      ([rN]) => rN === resolverDirectiveName
    )?.[1];
    if (!resolutionsMapItem) return '';
    return getUpstream(resolutionsMapItem);
  };

  const [tabIndex, setTabIndex] = React.useState(0);
  const [warningMessage, setWarningMessage] = React.useState('');
  const handleTabsChange = (index: number) => {
    setTabIndex(index);
  };
  const isEdit = Boolean(resolver);
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
      isNewResolver,
      fieldReturnType,
      objectType,
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
          });
      })
      .catch(err => {
        if (typeof err === 'object') {
          setWarningMessage(err.message);
        } else {
          setWarningMessage(err);
        }
        props.onClose();
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
        resolverName,
        objectType,
        isNewResolver,
        fieldReturnType,
      },
      true
    ).then(() => {
      setTimeout(() => {
        mutate();
        mutateSchemaYaml();
      }, 300);
      props.onClose();
    })
    .catch((err) => {
      if (typeof err === 'object') {
        setWarningMessage(err.message);
      } else {
        setWarningMessage(err);
      }
      props.onClose();
    });
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
          resolverType: getType(resolver),
          upstream: getUpstreamFromMap(),
          resolverConfig: getResolverFromConfig(resolver),
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
                    existingUpstream={resolver ? getUpstream(resolver) : ''}
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
