import React from 'react';
import { ISchemaDefinitionRecursiveItem } from '../SchemaDefinitionContent.type';
import { SchemaStyles } from './SchemaStyles.style';

const Description: React.FC<ISchemaDefinitionRecursiveItem> = props => {
  const { node, isRoot } = props;
  if (!node?.description?.value) return null;

  const descriptionEl = (
    <SchemaStyles.Description>
      - {node.description.value}
    </SchemaStyles.Description>
  );
  if (isRoot) return <SchemaStyles.Field>{descriptionEl}</SchemaStyles.Field>;
  return descriptionEl;
};

export default Description;
