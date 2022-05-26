import { Kind } from 'graphql';
import React from 'react';
import { Spacer } from 'Styles/StyledComponents/spacer';
import { getFieldReturnType } from 'utils/graphql-helpers';
import FieldTypeValue from '../FieldTypeValue';
import { ISchemaDefinitionRecursiveItem } from '../SchemaDefinitionContent.type';

const NameAndReturnType: React.FC<ISchemaDefinitionRecursiveItem> = props => {
  const { node } = props;
  const returnType = getFieldReturnType(node);
  return (
    <>
      {node.name?.value &&
        (node.kind === Kind.DIRECTIVE || node.kind === Kind.DIRECTIVE_DEFINITION
          ? '@' + node.name.value
          : node.name.value)}
      {!!node.variable?.name?.value && '$' + node.variable?.name?.value}
      {!!node.arguments?.length && (
        <>
          (<br />
          {node.arguments.map((a: any) => (
            <Spacer pl={3} key={a.name.value}>
              {a.name.value}:{' '}
              <span className='inline-block'>
                <FieldTypeValue
                  schema={props.schema}
                  field={a}
                  onReturnTypeClicked={props.onReturnTypeClicked}
                />
                {a.value?.value && ' ' + a.value.value}
              </span>
            </Spacer>
          ))}
          )
        </>
      )}
      {!!returnType?.fullType && (
        <>
          :{' '}
          <span className='inline-block'>
            <FieldTypeValue
              schema={props.schema}
              field={node}
              onReturnTypeClicked={props.onReturnTypeClicked}
            />
          </span>
        </>
      )}
      {!!node.locations?.length &&
        ' on ' + node.locations.map((l: any) => l.value).join(', ')}
    </>
  );
};

export default NameAndReturnType;
