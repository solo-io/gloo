import { EnumTypeDefinitionNode } from 'graphql';
import React from 'react';
import {
  ISchemaDefinitionContent,
  ISchemaDefinitionRecursiveItem,
} from '../../SchemaDefinitionContent.type';
import Description from '../Description';
import DirectiveList from '../DirectiveList';
import { SchemaStyles } from '../SchemaStyles.style';

const EnumTypeDefinition: React.FC<
  ISchemaDefinitionContent<EnumTypeDefinitionNode>
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
      <DirectiveList {...newProps} />
      {node.values?.map((v, i) => (
        <SchemaStyles.Field key={v.name?.value ?? i}>
          <div>{v.name.value}</div>
          <Description {...newProps} isRoot={false} />
          <DirectiveList {...newProps} isRoot={false} />
        </SchemaStyles.Field>
      ))}
    </>
  );
};

export default EnumTypeDefinition;
