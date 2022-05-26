import React from 'react';
import { ISchemaDefinitionRecursiveItem } from '../SchemaDefinitionContent.type';
import NameAndReturnType from './NameAndReturnType';
import { SchemaStyles } from './SchemaStyles.style';

const VariableDefinitions: React.FC<ISchemaDefinitionRecursiveItem> = props => {
  const { node } = props;
  if (!node.variableDefinitions.length) return null;
  return (
    <>
      {node.variableDefinitions.map((d: any, i: number) => (
        <SchemaStyles.Field key={d.name?.value ?? i}>
          <NameAndReturnType {...props} node={d} />
        </SchemaStyles.Field>
      ))}
    </>
  );
};

export default VariableDefinitions;
