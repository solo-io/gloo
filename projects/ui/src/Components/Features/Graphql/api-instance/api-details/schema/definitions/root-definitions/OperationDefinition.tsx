import { OperationDefinitionNode } from 'graphql';
import React from 'react';
import {
  ISchemaDefinitionContent,
  ISchemaDefinitionRecursiveItem,
} from '../../SchemaDefinitionContent.type';
import Description from '../Description';
import DirectiveList from '../DirectiveList';
import SelectionSet from '../SelectionSet';
import VariableDefinitions from '../VariableDefinitions';

const OperationDefinition: React.FC<
  ISchemaDefinitionContent<OperationDefinitionNode>
> = props => {
  const { node } = props;
  const newProps: ISchemaDefinitionRecursiveItem = {
    ...props,
    isRoot: true,
    objectType: node.name?.value ?? '',
  };
  return (
    <>
      <Description {...newProps} />
      <VariableDefinitions {...newProps} />
      <DirectiveList {...newProps} />
      <SelectionSet {...newProps} />
    </>
  );
};

export default OperationDefinition;
