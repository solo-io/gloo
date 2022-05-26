import React from 'react';
import { ISchemaDefinitionRecursiveItem } from '../SchemaDefinitionContent.type';
import NameAndReturnType from './NameAndReturnType';
import { SchemaStyles } from './SchemaStyles.style';

const SelectionSet: React.FC<
  ISchemaDefinitionRecursiveItem & { outerPadding?: boolean }
> = props => {
  const { node, isRoot, outerPadding } = props;
  if (isRoot)
    return (
      <>
        {node.selectionSet?.selections.map((v: any, i: number) => (
          <SelectionSet
            {...props}
            canAddResolverThisLevel={false}
            node={v}
            isRoot={false}
            outerPadding={true}
            key={v.name?.value ?? i}
          />
        ))}
      </>
    );
  return (
    <SchemaStyles.SelectionSet outerPadding={!!outerPadding}>
      <NameAndReturnType {...props} />
      {!!node.selectionSet?.selections.length && (
        <>
          {' {'}
          <br />
          {node.selectionSet?.selections.map((v: any, i: number) => (
            <SelectionSet
              {...props}
              canAddResolverThisLevel={false}
              node={v}
              isRoot={false}
              outerPadding={false}
              key={v.name?.value ?? i}
            />
          ))}
          {'}'}
        </>
      )}
    </SchemaStyles.SelectionSet>
  );
};

export default SelectionSet;
