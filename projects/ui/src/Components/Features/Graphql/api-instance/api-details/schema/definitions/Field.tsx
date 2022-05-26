import { useGetConsoleOptions } from 'API/hooks';
import React from 'react';
import ResolverButton from '../ResolverButton';
import { ISchemaDefinitionRecursiveItem } from '../SchemaDefinitionContent.type';
import Description from './Description';
import DirectiveList from './DirectiveList';
import NameAndReturnType from './NameAndReturnType';
import { SchemaStyles } from './SchemaStyles.style';

const Field: React.FC<ISchemaDefinitionRecursiveItem> = props => {
  const { node, isEditable, objectType } = props;
  const { readonly } = useGetConsoleOptions();
  const canAddResolver = isEditable && !readonly;
  const resolveDirective = node?.directives?.find(
    (d: any) => d.name.value === 'resolve'
  );
  return (
    <SchemaStyles.Field>
      <div className='flex'>
        <div className='grow'>
          <NameAndReturnType {...props} />
          <Description {...props} />
        </div>
        <div>
          {canAddResolver &&
            (!node.directives?.length || !!resolveDirective) && (
              <ResolverButton
                objectType={objectType}
                field={node}
                resolveDirectiveExists={!!resolveDirective}
              />
            )}
        </div>
      </div>
      <DirectiveList {...props} />
    </SchemaStyles.Field>
  );
};

export default Field;
