import { InterfaceTypeDefinitionNode } from 'graphql';
import React from 'react';
import {
  ISchemaDefinitionContent,
  ISchemaDefinitionRecursiveItem,
} from '../../SchemaDefinitionContent.type';
import Description from '../Description';
import DirectiveList from '../DirectiveList';
import FieldList from '../FieldList';

const InterfaceTypeDefinition: React.FC<
  ISchemaDefinitionContent<InterfaceTypeDefinitionNode>
> = props => {
  const { node } = props;
  const newProps: ISchemaDefinitionRecursiveItem = {
    ...props,
    isRoot: true,
    objectType: node.name.value,
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

export default InterfaceTypeDefinition;
