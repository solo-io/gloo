import { DefinitionNode } from 'graphql';
import { SupportedDocumentNode } from 'utils/graphql-helpers';

export type ISchemaDefinitionContent<T extends DefinitionNode> = {
  isEditable: boolean;
  schema: SupportedDocumentNode;
  node: T;
  onReturnTypeClicked(t: string): void;
};

export type ISchemaDefinitionRecursiveItem = {
  isRoot: boolean;
  isEditable: boolean;
  schema: SupportedDocumentNode;
  node: any;
  objectType: string;
  canAddResolverThisLevel?: boolean;
  onReturnTypeClicked(t: string): void;
};
