import { Kind } from 'graphql';
import React from 'react';
import { ISchemaDefinitionRecursiveItem } from '../SchemaDefinitionContent.type';
import Field from './Field';

const FieldList: React.FC<ISchemaDefinitionRecursiveItem> = props => {
  const { node } = props;
  return (
    <>
      {node.fields?.map((f: any, i: number) => (
        <Field
          {...props}
          canAddResolverThisLevel={node.kind === Kind.OBJECT_TYPE_DEFINITION}
          node={f}
          isRoot={false}
          key={f.name?.value ?? i}
        />
      ))}
    </>
  );
};

export default FieldList;
