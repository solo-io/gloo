import { useGetConsoleOptions } from 'API/hooks';
import React from 'react';
import { ISchemaDefinitionRecursiveItem } from '../SchemaDefinitionContent.type';
import NameAndReturnType from './NameAndReturnType';
import { SchemaStyles } from './SchemaStyles.style';

const DirectiveList: React.FC<ISchemaDefinitionRecursiveItem> = props => {
  const { node, isEditable } = props;
  const { readonly } = useGetConsoleOptions();
  const canAddResolver = isEditable && !readonly;
  if (!node.directives?.length) return null;
  return (
    <>
      {node.directives?.map((d: any, i: number) => {
        //
        // This skips showing the resolve directive only if
        // the add/edit resolver button will be shown.
        if (d.name.value === 'resolve' && canAddResolver) return null;
        //
        // Other than that case, it shows all the directives.
        return (
          <SchemaStyles.Field key={d.name?.value ?? i} spaceY={0.5}>
            <NameAndReturnType {...props} node={d} />
          </SchemaStyles.Field>
        );
      })}
    </>
  );
};

export default DirectiveList;
