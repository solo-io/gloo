import styled from '@emotion/styled/macro';
import { useGetGraphqlSchemaDetails } from 'API/hooks';
import { ReactComponent as RouteIcon } from 'assets/route-icon.svg';
import { SoloModal } from 'Components/Common/SoloModal';
import { ReactComponent as CodeIcon } from 'assets/code-icon.svg';
import { ReactComponent as GraphQLIcon } from 'assets/graphql-icon.svg';
import React from 'react';
import { useParams } from 'react-router';
import { useVirtual } from 'react-virtual';
import { colors } from 'Styles/colors';
import tw from 'twin.macro';
import gql from 'graphql-tag';
import {
  buildASTSchema,
  EnumTypeDefinitionNode,
  EnumValueDefinitionNode,
  FieldDefinitionNode,
  GraphQLSchema,
  NamedTypeNode,
  ObjectTypeDefinitionNode,
  printSchema,
  //@ts-ignore
} from 'graphql';
import { ResolverWizard } from './ResolverWizard';
import { GraphqlIconHolder } from './GraphqlTable';
import { ClusterObjectRef } from 'proto/github.com/solo-io/skv2/api/core/v1/core_pb';
import { mapSchema, getDirective, MapperKind } from '@graphql-tools/utils';
import { EnumResolver } from './components/EnumResolver';
import { makeExecutableSchema } from '@graphql-tools/schema';

type ArrowToggleProps = { active?: boolean };
export const ArrowToggle = styled('div')<ArrowToggleProps>`
  position: absolute;
  left: 1rem;

  &:before,
  &:after {
    position: absolute;
    content: '';
    display: block;
    width: 8px;
    height: 1px;
    background: ${colors.septemberGrey};
    transition: transform 0.3s;
  }

  &:before {
    right: 5px;
    border-top-left-radius: 10px;
    border-bottom-left-radius: 10px;
    transform: rotate(${props => (props.active ? '-' : '')}45deg);
  }

  &:after {
    right: 1px;
    transform: rotate(${props => (props.active ? '' : '-')}45deg);
  }
`;

export const OperationDescription = styled('div')`
  ${tw`w-full overflow-y-scroll text-sm text-gray-600 whitespace-normal`};
  grid-column: span 3 / span 3;
  /* Hide scrollbar for Chrome, Safari and Opera */
  &::-webkit-scrollbar {
    display: none !important;
  }

  /* Hide scrollbar for IE, Edge and Firefox */
  & {
    -ms-overflow-style: none !important; /* IE and Edge */
    scrollbar-width: none !important; /* Firefox */
  }
`;
type EndpointCircleProps = {
  isFirstSubrow?: boolean;
};

export const EndpointCircle = styled.div<EndpointCircleProps>`
  ${tw`relative w-3 h-3 pl-3 mx-3 border border-gray-400 rounded-full`};

  &:before {
    content: '';
    position: absolute;
    border-left: 1px dotted ${colors.aprilGrey};
    border-bottom: 1px dotted ${colors.aprilGrey};
    left: -0.75rem;
    width: 0.6rem;
    ${props =>
      props.isFirstSubrow
        ? `top: -.9rem; height: 1.2rem;`
        : `top: -.9rem; height: 1.2rem;`};
  }
`;

export const ResolverItem: React.FC<{
  resolverType: string;
  fields: FieldDefinitionNode[] | EnumTypeDefinitionNode[];
  handleResolverConfigModal: (
    resolverName: string,
    resolverType: string
  ) => void;
}> = props => {
  const { resolverType, fields, handleResolverConfigModal } = props;
  const [isOpen, setIsOpen] = React.useState(
    resolverType === 'Query' || resolverType === 'Mutation'
  );
  const listRef = React.useRef<HTMLDivElement>(null);
  const {
    graphqlSchemaName = '',
    graphqlSchemaNamespace = '',
    graphqlSchemaClusterName = '',
  } = useParams();
  const resolverKey = `${graphqlSchemaNamespace}-${graphqlSchemaName}-${resolverType}`;
  const { data: graphqlSchema, error: graphqlSchemaError } =
    useGetGraphqlSchemaDetails({
      name: graphqlSchemaName,
      namespace: graphqlSchemaNamespace,
      clusterName: graphqlSchemaClusterName,
    });
  const rowVirtualizer = useVirtual({
    size: fields?.length ?? 0,
    parentRef: listRef,
    estimateSize: React.useCallback(() => 90, []),
    overscan: 1,
  });

  return (
    <div key={resolverKey}>
      <div className='relative flex flex-col w-full bg-gray-200 border h-28'>
        <div className='flex items-center justify-between gap-5 pt-4 my-2 ml-4 h-14 '>
          <div className='flex items-center mr-3'>
            <GraphqlIconHolder>
              <GraphQLIcon className='w-4 h-4 fill-current' />
            </GraphqlIconHolder>
            <span className='flex items-center font-medium text-gray-900 whitespace-nowrap'>
              {resolverType}
            </span>
          </div>
        </div>
        <div className='flex items-center justify-between w-full px-6 py-4 text-sm font-medium text-gray-900 whitespace-nowrap'>
          <div
            className='relative flex-wrap justify-between w-full h-full text-sm '
            style={{
              display: 'grid',
              flexWrap: 'wrap',
              gridTemplateColumns: '1fr 1fr  1fr',
              gridTemplateRows: '1fr',
              gridAutoRows: 'min-content',
              columnGap: '15px',
            }}>
            <span className='flex items-center justify-start ml-6 font-medium text-gray-900 '>
              Field Name
            </span>
            <span className='flex items-center justify-start ml-8 font-medium text-gray-900 '>
              Return Type
            </span>

            <span className='flex items-center justify-center ml-8 font-medium text-gray-900 '>
              Resolver
            </span>
          </div>
        </div>
        <div
          className='absolute top-0 right-0 flex items-center w-10 h-10 p-4 mr-2 cursor-pointer '
          onClick={() => setIsOpen(!isOpen)}>
          <ArrowToggle active={isOpen} className='self-center m-4 ' />
        </div>
      </div>
      {isOpen && (
        <div
          ref={listRef}
          style={{
            height: `${
              fields?.length * 90 < 400 ? fields!.length * 90 : 400
            }px`,
            width: `100%`,
            overflow: 'auto',
          }}>
          <div
            style={{
              height: `${rowVirtualizer.totalSize}px`,
              width: '100%',
              position: 'relative',
            }}>
            {rowVirtualizer.virtualItems.map(virtualRow => {
              const op = fields[virtualRow.index] as FieldDefinitionNode;
              let hasResolver =
                !!graphqlSchema?.spec?.executableSchema?.executor?.local?.resolutionsMap?.find(
                  ([rN, r]) => rN.includes(fields[virtualRow.index].name?.value)
                );
              return (
                <div
                  key={`${resolverType}-${op.name?.value}`}
                  className={`flex h-20 p-2 pl-0 border `}
                  style={{
                    position: 'absolute',
                    top: 0,
                    left: 0,
                    width: '100%',
                    height: `${virtualRow.size}px`,
                    transform: `translateY(${virtualRow.start}px)`,
                  }}>
                  <div className='flex items-center px-4 text-sm font-medium text-gray-900 whitespace-nowrap'>
                    <CodeIcon className='w-4 h-4 ml-2 mr-3 fill-current text-blue-600gloo' />
                  </div>
                  <div className='relative flex items-center w-full text-sm text-gray-500 whitespace-nowrap'>
                    <div
                      className='relative flex-wrap justify-between w-full h-full text-sm '
                      style={{
                        display: 'grid',
                        flexWrap: 'wrap',
                        gridTemplateColumns:
                          '1fr 1fr  minmax(120px, 200px) 105px',
                        gridTemplateRows: op.description?.value
                          ? ' 1fr min-content'
                          : '1fr',
                        gridAutoRows: 'min-content',
                        columnGap: '5px',
                        rowGap: '5px',
                      }}>
                      <span className='flex items-center font-medium text-gray-900 '>
                        {fields[virtualRow.index].name?.value ?? ''}
                      </span>
                      <span className='flex items-center text-sm text-gray-700 '>
                        {(op.type as NamedTypeNode).name?.value}
                      </span>
                      <span className={`flex items-center  justify-center`}>
                        <span
                          className={`inline-flex items-center min-w-max p-1 px-2 ${
                            hasResolver
                              ? 'focus:ring-blue-500gloo text-blue-700gloo bg-blue-200gloo  border-blue-600gloo hover:bg-blue-300gloo'
                              : 'focus:ring-gray-500 text-gray-700 bg-gray-300  border-gray-600 hover:bg-gray-200'
                          }   border rounded-full shadow-sm cursor-pointer  focus:outline-none focus:ring-2 focus:ring-offset-2 `}
                          onClick={() => {
                            handleResolverConfigModal(
                              fields[virtualRow.index].name?.value ?? '',
                              resolverType
                            );
                          }}>
                          {hasResolver && (
                            <RouteIcon className='w-6 h-6 mr-1 fill-current text-blue-600gloo' />
                          )}
                          {hasResolver ? 'Resolver' : 'Define Resolver'}
                        </span>
                      </span>

                      {op.description && (
                        <OperationDescription className='w-full overflow-y-scroll text-sm text-gray-600 whitespace-normal'>
                          {op.description?.value}
                        </OperationDescription>
                      )}
                    </div>
                  </div>
                </div>
              );
            })}
          </div>
        </div>
      )}
    </div>
  );
};

function defineResolveDirective() {
  let directiveName = 'resolve';
  return {
    mockedDirectiveTypeDefs: `directive @${directiveName}(name: String) on FIELD_DEFINITION  `,
    mockedDirectiveTransformer: (schema: GraphQLSchema) =>
      mapSchema(schema, {
        [MapperKind.OBJECT_FIELD]: fieldConfig => {
          const mockedDirective = getDirective(
            schema,
            fieldConfig,
            directiveName
          )?.[0];
          if (mockedDirective) {
            fieldConfig.deprecationReason = mockedDirective['name'];
            return fieldConfig;
          }
        },
        [MapperKind.ENUM_VALUE]: enumValueConfig => {
          const mockedDirective = getDirective(
            schema,
            enumValueConfig,
            directiveName
          )?.[0];
          if (mockedDirective) {
            enumValueConfig.deprecationReason = mockedDirective['name'];
            return enumValueConfig;
          }
        },
      }),
  };
}
type ResolversTableType = {
  schemaRef: ClusterObjectRef.AsObject;
};
const ResolversTable: React.FC<ResolversTableType> = props => {
  const { schemaRef } = props;
  const {
    data: graphqlSchema,
    error: graphqlSchemaError,
    mutate,
  } = useGetGraphqlSchemaDetails({
    name: schemaRef.name,
    namespace: schemaRef.namespace,
    clusterName: schemaRef.clusterName,
  });

  const [currentResolver, setCurrentResolver] = React.useState<any>();
  const [currentResolverName, setCurrentResolverName] = React.useState('');
  const [currentResolverType, setCurrentResolverType] = React.useState('');
  const [hasDirective, setHasDirective] = React.useState(false);
  const [fieldWithDirective, setFieldWithDirective] = React.useState('');
  const [fieldWithoutDirective, setFieldWithoutDirective] = React.useState('');
  const [modalOpen, setModalOpen] = React.useState(false);
  const openModal = () => setModalOpen(true);
  const closeModal = () => setModalOpen(false);

  const listRef = React.useRef<HTMLDivElement>(null);

  const [fieldTypesMap, setFieldTypesMap] = React.useState<
    [string, FieldDefinitionNode[]][]
  >([]);
  const [enumTypesMap, setEnumTypesMap] = React.useState<
    [string, EnumValueDefinitionNode[]][]
  >([]);

  React.useEffect(() => {
    if (graphqlSchema) {
      let query = gql`
        ${graphqlSchema.spec?.executableSchema?.schemaDefinition}
      `;
      if (query) {
        let objectTypeDefs = query.definitions.filter(
          (def: any) => def.kind === 'ObjectTypeDefinition'
        ) as ObjectTypeDefinitionNode[];
        let enumTypeDefs = query.definitions.filter(
          (def: any) => def.kind === 'EnumTypeDefinition'
        ) as EnumTypeDefinitionNode[];

        const enumFieldDefinitions = enumTypeDefs?.map(ot => [
          `Enum ${ot.name.value}`,
          (ot).values?.filter(
            f => f?.kind === 'EnumValueDefinition'
          ) as EnumValueDefinitionNode[],
        ]) as [string, EnumValueDefinitionNode[]][];

        let fieldDefinitions = objectTypeDefs.map(ot => [
          ot.name.value,
          (ot).fields?.filter(
            f => f?.kind === 'FieldDefinition'
          ) as FieldDefinitionNode[],
        ]) as [string, FieldDefinitionNode[]][];
        setFieldTypesMap(fieldDefinitions);
        setEnumTypesMap(enumFieldDefinitions);
      }
    }
  }, [graphqlSchema]);

  function handleResolverConfigModal(
    resolverName: string,
    resolverType: string
  ) {
    let [currentResolverName, currentResolver] =
      graphqlSchema?.spec?.executableSchema?.executor?.local?.resolutionsMap.find(
        ([rName, resolver]) => rName.includes(resolverName)
      ) ?? ['', ''];
    setCurrentResolver(currentResolver);
    setCurrentResolverName(resolverName);
    setCurrentResolverType(resolverType);

    let isListType = Object.fromEntries(fieldTypesMap)[resolverType]?.some(
      f => f.name.value === resolverName && f.type.kind === 'ListType'
    );

    let fieldType = isListType
      ? fieldTypesMap
          .find(([t, f]) => t === resolverType)?.[1]
          //@ts-ignore
          ?.find(f => f.name.value === resolverName)?.type?.type?.name?.value
      : fieldTypesMap
          .find(([t, f]) => t === resolverType)?.[1]
          //@ts-ignore
          ?.find(f => f.name.value === resolverName)?.type?.name?.value;

    let fieldWithDirective = '';
    let fieldWithoutDirective = '';
    if (isListType) {
      fieldWithoutDirective = `${resolverName}: [${fieldType}]`;
      fieldWithDirective = `${resolverName}: [${fieldType}] @resolve(name: "${resolverName}")`;
    } else {
      fieldWithoutDirective = `${resolverName}: ${fieldType}`;

      fieldWithDirective = `${resolverName}: ${fieldType} @resolve(name: "${resolverName}")`;
    }

    setHasDirective(
      !!graphqlSchema?.spec?.executableSchema?.schemaDefinition.includes(
        fieldWithDirective
      )
    );
    setFieldWithDirective(fieldWithDirective);
    setFieldWithoutDirective(fieldWithoutDirective);

    openModal();
  }

  return (
    <>
      <div className='flex flex-col w-full '>
        <div
          className='relative space-y-6 overflow-x-hidden overflow-y-scroll '
          ref={listRef}>
          {fieldTypesMap
            ?.sort(([aTypeName], [bTypeName]) => {
              // Ordering: Query, mutation, Everything else.
              if (aTypeName === 'Query') {
                return -1;
              } else if (bTypeName === 'Query') {
                return 1;
              }
              if (aTypeName === 'Mutation') {
                return -1;
              } else if (bTypeName === 'Mutation') {
                return 1;
              }
              return 0;
            })
            .map(([typeName, fields]) => {
              return (
                <ResolverItem
                  key={`${schemaRef.namespace}-${schemaRef.name}-${typeName}`}
                  resolverType={typeName}
                  fields={fields}
                  handleResolverConfigModal={handleResolverConfigModal}
                />
              );
            })}
            {
              enumTypesMap.map(([typeName, fields]) => {
                return (
                  <EnumResolver key={typeName} fields={fields} resolverType={typeName} />
                )
              }
            )}

          <SoloModal visible={modalOpen} width={750} onClose={closeModal}>
            <ResolverWizard
              resolver={currentResolver}
              hasDirective={hasDirective}
              fieldWithDirective={fieldWithDirective}
              fieldWithoutDirective={fieldWithoutDirective}
              resolverName={currentResolverName}
              onClose={() => {
                closeModal();
                mutate();
              }}
            />
          </SoloModal>
        </div>
      </div>
    </>
  );
};

export { ResolversTable as default };
