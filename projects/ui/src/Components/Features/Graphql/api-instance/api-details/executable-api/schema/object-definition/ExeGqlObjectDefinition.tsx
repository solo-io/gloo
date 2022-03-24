import { useGetGraphqlApiDetails } from 'API/hooks';
import { ReactComponent as CodeIcon } from 'assets/code-icon.svg';
import { ReactComponent as RouteIcon } from 'assets/route-icon.svg';
import { SoloModal } from 'Components/Common/SoloModal';
import {
  EnumTypeDefinitionNode,
  FieldDefinitionNode,
  Kind,
  NamedTypeNode,
  ObjectTypeDefinitionNode,
} from 'graphql';
import { ClusterObjectRef } from 'proto/github.com/solo-io/skv2/api/core/v1/core_pb';
import React from 'react';
import { useVirtual } from 'react-virtual';
import { colors } from 'Styles/colors';
import * as styles from '../ExecutableGraphqlSchemaDefinitions.style';
import { ResolverWizard } from './resolver-wizard/ResolverWizard';

export const ExeGqlObjectDefinition: React.FC<{
  apiRef: ClusterObjectRef.AsObject;
  resolverType: string;
  onReturnTypeClicked(t: string): void;
  schemaDefinitions: (ObjectTypeDefinitionNode | EnumTypeDefinitionNode)[];
  fields: readonly FieldDefinitionNode[];
}> = ({
  apiRef,
  resolverType,
  onReturnTypeClicked,
  schemaDefinitions,
  fields,
}) => {
  const listRef = React.useRef<HTMLDivElement>(null);
  const resolverKey = `${apiRef.namespace}-${apiRef.name}-${resolverType}`;
  const { data: graphqlApi, mutate } = useGetGraphqlApiDetails(apiRef);
  const rowVirtualizer = useVirtual({
    size: fields?.length ?? 0,
    parentRef: listRef,
    estimateSize: React.useCallback(() => 90, []),
    overscan: 1,
  });

  // --- RESOLVER CONFIG MODAL --- //
  const [currentResolver, setCurrentResolver] = React.useState<any>();
  const [currentResolverName, setCurrentResolverName] = React.useState('');
  const [hasDirective, setHasDirective] = React.useState(false);
  const [fieldWithDirective, setFieldWithDirective] = React.useState('');
  const [fieldWithoutDirective, setFieldWithoutDirective] = React.useState('');
  const [isModalOpen, setIsModalOpen] = React.useState(false);
  function handleResolverConfigModal(
    resolverName: string,
    resolverType: string
  ) {
    let [_currentResolverName, currentResolver] =
      graphqlApi?.spec?.executableSchema?.executor?.local?.resolutionsMap.find(
        ([rName, _resolver]) => rName.includes(resolverName)
      ) ?? ['', ''];
    setCurrentResolver(currentResolver);
    setCurrentResolverName(resolverName);
    //
    // Get the definition and field for the selected resolver from the schema.
    const definition = schemaDefinitions.find(
      d =>
        d.kind === Kind.OBJECT_TYPE_DEFINITION && d.name.value === resolverType
    ) as ObjectTypeDefinitionNode | undefined;
    if (definition === undefined) return;
    const field = definition.fields?.find(f => f.name.value === resolverName);
    if (field === undefined) return;
    //
    // Find the base field type (this could be a nested list).
    let baseField = field.type;
    let typePrefix = '';
    let typeSuffix = '';
    while (baseField?.kind === Kind.LIST_TYPE) {
      typePrefix += '[';
      typeSuffix += ']';
      baseField = baseField.type;
    }
    if (baseField?.kind !== Kind.NAMED_TYPE) return;
    let fieldType = typePrefix + baseField?.name?.value + typeSuffix;
    //
    // Build the directive strings for graphql.
    let fieldWithoutDirective = `${resolverName}: ${fieldType}]`;
    let fieldWithDirective = `${resolverName}: ${fieldType} @resolve(name: "${resolverName}")`;
    //
    // Update state.
    setHasDirective(
      !!graphqlApi?.spec?.executableSchema?.schemaDefinition.includes(
        fieldWithDirective
      )
    );
    setFieldWithDirective(fieldWithDirective);
    setFieldWithoutDirective(fieldWithoutDirective);
    setIsModalOpen(true);
  }

  return (
    <div data-testid='resolver-item' key={resolverKey}>
      <SoloModal
        visible={isModalOpen}
        width={750}
        onClose={() => setIsModalOpen(false)}>
        <ResolverWizard
          resolver={currentResolver}
          hasDirective={hasDirective}
          fieldWithDirective={fieldWithDirective}
          fieldWithoutDirective={fieldWithoutDirective}
          resolverName={currentResolverName}
          onClose={() => {
            setIsModalOpen(false);
            mutate();
          }}
        />
      </SoloModal>

      <div
        className='relative flex flex-col w-full py-3 border'
        style={{
          backgroundColor: colors.lightJanuaryGrey,
          display: 'grid',
          flexWrap: 'wrap',
          // Duplicating the content's gridTemplateColumns so it is centered.
          gridTemplateColumns: '1fr 1fr  minmax(120px, 200px) 105px',
          gridTemplateRows: '1fr',
          gridAutoRows: 'min-content',
          columnGap: '15px',
        }}>
        <span className='flex items-center justify-start ml-6 font-medium text-gray-900 '>
          Field Name
        </span>
        <span className='flex items-center justify-start ml-8 font-medium text-gray-900 '>
          Type
        </span>
        <span className='flex items-center justify-center ml-8 font-medium text-gray-900 '>
          Resolver
        </span>
        <span />
      </div>

      <div
        ref={listRef}
        style={{
          height: `${fields?.length * 90 < 400 ? fields!.length * 90 : 400}px`,
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
              !!graphqlApi?.spec?.executableSchema?.executor?.local?.resolutionsMap?.find(
                ([rN, r]) => rN.includes(fields[virtualRow.index].name?.value)
              );
            const getReturnType = () => {
              let prefix = '';
              let suffix = '';
              let cur = op.type;
              while (cur.kind === Kind.LIST_TYPE) {
                cur = cur.type;
                suffix += '[]';
              }
              if (cur.kind === Kind.NAMED_TYPE)
                return [prefix, cur.name.value, suffix] as [
                  string,
                  string,
                  string
                ];
              else return ['', '', ''] as [string, string, string];
            };
            const [returnTypePrefix, baseReturnType, returnTypeSuffix] =
              getReturnType();
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
                    <span
                      className='flex items-center text-sm text-gray-700 '
                      style={{ fontFamily: 'monospace' }}>
                      {returnTypePrefix}
                      {schemaDefinitions.find(
                        d => d.name.value === baseReturnType
                      ) ? (
                        <a
                          style={{ fontFamily: 'monospace' }}
                          onClick={() => onReturnTypeClicked(baseReturnType)}>
                          {baseReturnType}
                        </a>
                      ) : (
                        <>{baseReturnType}</>
                      )}
                      {returnTypeSuffix}
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
                      <styles.OperationDescription>
                        {op.description?.value}
                      </styles.OperationDescription>
                    )}
                  </div>
                </div>
              </div>
            );
          })}
        </div>
      </div>
    </div>
  );
};
