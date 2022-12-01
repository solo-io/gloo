import { ObjectTypeDefinitionNode } from 'graphql';
import React from 'react';
import FieldTypeValue from '../../FieldTypeValue';
import {
  ISchemaDefinitionContent,
  ISchemaDefinitionRecursiveItem,
} from '../../SchemaDefinitionContent.type';
import Description from '../Description';
import DirectiveList from '../DirectiveList';
import FieldList from '../FieldList';
import { SchemaStyles } from '../SchemaStyles.style';

const ObjectTypeDefinition: React.FC<
  ISchemaDefinitionContent<ObjectTypeDefinitionNode>
> = props => {
  const { node } = props;
  const newProps: ISchemaDefinitionRecursiveItem = {
    ...props,
    isRoot: true,
    objectType: node.name.value,
  };
  return (
    <>
      {!!node.interfaces?.length && (
        <SchemaStyles.Field>
          <SchemaStyles.Description className='inline-block'>
            Implements:&nbsp;
          </SchemaStyles.Description>
          {node.interfaces.map((f, i) => (
            <div className='inline-block' key={f.name.value}>
              {i > 0 && ', '}
              <FieldTypeValue {...newProps} field={f} />
            </div>
          ))}
        </SchemaStyles.Field>
      )}
      <Description {...newProps} />
      <DirectiveList {...newProps} />
      <FieldList {...newProps} />
    </>
  );
};

export default ObjectTypeDefinition;
