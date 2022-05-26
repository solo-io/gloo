import { InputObjectTypeDefinitionNode } from 'graphql';
import React from 'react';
import {
  ISchemaDefinitionContent,
  ISchemaDefinitionRecursiveItem,
} from '../../SchemaDefinitionContent.type';
import Description from '../Description';
import DirectiveList from '../DirectiveList';
import FieldList from '../FieldList';

const InputObjectTypeDefinition: React.FC<
  ISchemaDefinitionContent<InputObjectTypeDefinitionNode>
> = props => {
  const { node } = props;
  const newProps: ISchemaDefinitionRecursiveItem = {
    ...props,
    objectType: node.name.value,
    isRoot: true,
    isEditable: false,
  };
  return (
    <>
      <Description {...newProps} />
      <DirectiveList {...newProps} />
      <FieldList {...newProps} />
    </>
  );
};

export default InputObjectTypeDefinition;
