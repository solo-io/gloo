import { DefinitionNode, Kind } from 'graphql';
import React from 'react';
import { Spacer } from 'Styles/StyledComponents/spacer';
import DirectiveDefinition from './definitions/root-definitions/DirectiveDefinition';
import EnumTypeDefinition from './definitions/root-definitions/EnumTypeDefinition';
import InputObjectTypeDefinition from './definitions/root-definitions/InputObjectTypeDefinition';
import InterfaceTypeDefinition from './definitions/root-definitions/InterfaceTypeDefinition';
import ObjectTypeDefinition from './definitions/root-definitions/ObjectTypeDefinition';
import OperationDefinition from './definitions/root-definitions/OperationDefinition';
import UnionTypeDefinition from './definitions/root-definitions/UnionTypeDefinition';
import { ISchemaDefinitionContent } from './SchemaDefinitionContent.type';

const SchemaDefinitionContent: React.FC<
  ISchemaDefinitionContent<DefinitionNode>
> = props => {
  const { node } = props;
  return (
    <Spacer px={3} py={1}>
      {node.kind === Kind.ENUM_TYPE_DEFINITION && (
        <EnumTypeDefinition {...props} node={node} />
      )}
      {node.kind === Kind.OBJECT_TYPE_DEFINITION && (
        <ObjectTypeDefinition {...props} node={node} />
      )}
      {node.kind === Kind.OPERATION_DEFINITION && (
        <OperationDefinition {...props} node={node} />
      )}
      {node.kind === Kind.INPUT_OBJECT_TYPE_DEFINITION && (
        <InputObjectTypeDefinition {...props} node={node} />
      )}
      {node.kind === Kind.UNION_TYPE_DEFINITION && (
        <UnionTypeDefinition {...props} node={node} />
      )}
      {node.kind === Kind.INTERFACE_TYPE_DEFINITION && (
        <InterfaceTypeDefinition {...props} node={node} />
      )}
      {node.kind === Kind.DIRECTIVE_DEFINITION && (
        <DirectiveDefinition {...props} node={node} />
      )}
    </Spacer>
  );
};

export default SchemaDefinitionContent;
