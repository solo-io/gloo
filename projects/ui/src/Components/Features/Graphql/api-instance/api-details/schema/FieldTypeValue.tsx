import { FieldDefinitionNode, TypeNode } from 'graphql';
import React from 'react';
import {
  getFieldReturnType,
  SupportedDocumentNode,
} from 'utils/graphql-helpers';

const FieldTypeValue: React.FC<{
  schema: SupportedDocumentNode;
  field: FieldDefinitionNode | TypeNode;
  onReturnTypeClicked(returnType: string): void;
}> = ({ schema, field, onReturnTypeClicked }) => {
  const returnType = getFieldReturnType(field);
  return (
    <span
      className='flex items-center text-sm text-gray-700 '
      style={{ fontFamily: 'monospace' }}>
      {returnType.parts.prefix}
      {schema.definitions.find(d => d.name?.value === returnType.parts.base) ? (
        <a
          style={{ fontFamily: 'monospace' }}
          onClick={() => onReturnTypeClicked(returnType.parts.base)}>
          {returnType.parts.base}
        </a>
      ) : (
        <>{returnType.parts.base}</>
      )}
      {returnType.parts.suffix}
    </span>
  );
};

export default FieldTypeValue;
