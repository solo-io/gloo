import { UnionTypeDefinitionNode } from 'graphql';
import React from 'react';
import FieldTypeValue from '../../FieldTypeValue';
import {
  ISchemaDefinitionContent,
  ISchemaDefinitionRecursiveItem,
} from '../../SchemaDefinitionContent.type';
import Description from '../Description';
import DirectiveList from '../DirectiveList';
import { SchemaStyles } from '../SchemaStyles.style';

const UnionTypeDefinition: React.FC<
  ISchemaDefinitionContent<UnionTypeDefinitionNode>
> = props => {
  const { node } = props;
  const newProps: ISchemaDefinitionRecursiveItem = {
    ...props,
    isRoot: true,
    objectType: node.name.value,
  };
  return (
    <>
      <Description {...newProps} />
      <DirectiveList {...newProps} />
      {node.types?.map((t, i) => (
        <SchemaStyles.Field key={t.name?.value ?? i}>
          <FieldTypeValue {...newProps} field={t} />
        </SchemaStyles.Field>
      ))}
    </>
  );
};

export default UnionTypeDefinition;
