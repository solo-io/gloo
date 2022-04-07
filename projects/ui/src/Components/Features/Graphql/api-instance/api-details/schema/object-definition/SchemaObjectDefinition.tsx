import {
  useGetConsoleOptions,
  useGetGraphqlApiDetails,
  useGetGraphqlApiYaml,
} from 'API/hooks';
import { ReactComponent as CodeIcon } from 'assets/code-icon.svg';
import { ReactComponent as RouteIcon } from 'assets/route-icon.svg';
import { SoloModal } from 'Components/Common/SoloModal';
import { FieldDefinitionNode, ObjectTypeDefinitionNode } from 'graphql';
import { ClusterObjectRef } from 'proto/github.com/solo-io/skv2/api/core/v1/core_pb';
import React, { useState } from 'react';
import { useVirtual } from 'react-virtual';
import { colors } from 'Styles/colors';
import {
  getFieldReturnType,
  hasResolutionForField,
  SupportedDocumentNode,
} from 'utils/graphql-helpers';
import * as styles from '../SchemaDefinitions.style';
import { ResolverWizard } from './resolver-wizard/ResolverWizard';

export const SchemaObjectDefinition: React.FC<{
  apiRef: ClusterObjectRef.AsObject;
  isEditable: boolean;
  schema: SupportedDocumentNode;
  objectTypeDefinition: ObjectTypeDefinitionNode;
  onReturnTypeClicked(t: string): void;
}> = ({
  apiRef,
  isEditable,
  schema,
  objectTypeDefinition,
  onReturnTypeClicked,
}) => {
  const { data: graphqlApi, mutate: mutateDetails } =
    useGetGraphqlApiDetails(apiRef);
  const { mutate: mutateYaml } = useGetGraphqlApiYaml(apiRef);
  const objectType = objectTypeDefinition.name.value;
  const fields = objectTypeDefinition.fields ?? [];
  const listRef = React.useRef<HTMLDivElement>(null);
  const { readonly } = useGetConsoleOptions();

  const rowVirtualizer = useVirtual({
    size: fields.length,
    parentRef: listRef,
    estimateSize: React.useCallback(() => 90, []),
    overscan: 1,
  });

  // --- RESOLVER CONFIG MODAL --- //
  const [selectedField, setSelectedFieldName] =
    useState<FieldDefinitionNode | null>(null);

  return (
    <div data-testid='resolver-item'>
      <SoloModal
        visible={selectedField !== null}
        width={750}
        onClose={() => setSelectedFieldName(null)}>
        <ResolverWizard
          apiRef={apiRef}
          field={selectedField}
          objectType={objectType}
          onClose={() => {
            setTimeout(() => {
              mutateDetails();
              mutateYaml();
            }, 300);
            setSelectedFieldName(null);
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
          {!readonly && isEditable && 'Resolver'}
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
            const field = fields[virtualRow.index];
            const fieldName = field.name.value ?? '';
            const resolutionExists = hasResolutionForField(graphqlApi, field);
            const returnType = getFieldReturnType(field);
            return (
              <div
                key={`${objectType}-${fieldName}`}
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
                      gridTemplateRows: field.description?.value
                        ? ' 1fr min-content'
                        : '1fr',
                      gridAutoRows: 'min-content',
                      columnGap: '5px',
                      rowGap: '5px',
                    }}>
                    <span className='flex items-center font-medium text-gray-900 '>
                      {fieldName}
                    </span>
                    <span
                      className='flex items-center text-sm text-gray-700 '
                      style={{ fontFamily: 'monospace' }}>
                      {returnType.parts.prefix}
                      {schema.definitions.find(
                        d => d.name.value === returnType.parts.base
                      ) ? (
                        <a
                          style={{ fontFamily: 'monospace' }}
                          onClick={() =>
                            onReturnTypeClicked(returnType.parts.base)
                          }>
                          {returnType.parts.base}
                        </a>
                      ) : (
                        <>{returnType.parts.base}</>
                      )}
                      {returnType.parts.suffix}
                    </span>
                    <span className={`flex items-center  justify-center`}>
                      {!readonly && isEditable && (
                        <span
                          data-testid={`resolver-${field.name.value}`}
                          className={`inline-flex items-center min-w-max p-1 px-2 ${
                            resolutionExists
                              ? 'focus:ring-blue-500gloo text-blue-700gloo bg-blue-200gloo  border-blue-600gloo hover:bg-blue-300gloo'
                              : 'focus:ring-gray-500 text-gray-700 bg-gray-300  border-gray-600 hover:bg-gray-200'
                          }   border rounded-full shadow-sm cursor-pointer  focus:outline-none focus:ring-2 focus:ring-offset-2 `}
                          onClick={() => setSelectedFieldName(field)}>
                          {resolutionExists && (
                            <RouteIcon
                              data-testid={`route-${field.name.value}`}
                              className='w-6 h-6 mr-1 fill-current text-blue-600gloo'
                            />
                          )}
                          {resolutionExists ? 'Resolver' : 'Define Resolver'}
                        </span>
                      )}
                    </span>

                    {field.description && (
                      <styles.OperationDescription>
                        {field.description?.value}
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
