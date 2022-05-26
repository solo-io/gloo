import { DirectiveDefinitionNode } from 'graphql';
import React from 'react';
import {
  ISchemaDefinitionContent,
  ISchemaDefinitionRecursiveItem,
} from '../../SchemaDefinitionContent.type';
import Description from '../Description';
import NameAndReturnType from '../NameAndReturnType';
import { SchemaStyles } from '../SchemaStyles.style';

const DirectiveDefinition: React.FC<
  ISchemaDefinitionContent<DirectiveDefinitionNode>
> = props => {
  const { node } = props;
  const newProps: ISchemaDefinitionRecursiveItem = {
    ...props,
    objectType: node.name.value,
    isRoot: true,
  };
  return (
    <>
      <Description {...newProps} />
      <SchemaStyles.Field>
        <NameAndReturnType {...newProps} />
      </SchemaStyles.Field>
    </>
  );
};

export default DirectiveDefinition;
