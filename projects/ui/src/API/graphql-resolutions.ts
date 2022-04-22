import jspb from 'google-protobuf';
import { ASTNode, FieldDefinitionNode, Kind, print, visit } from 'graphql';
import {
  ExecutableSchema,
  Resolution,
} from 'proto/github.com/solo-io/solo-apis/api/gloo/graphql.gloo/v1beta1/graphql_pb';
import {
  getResolveDirectiveName,
  parseSchemaString,
  SupportedDocumentNode,
} from 'utils/graphql-helpers';
import { ResolverItem } from './graphql';

//
// --- ADD RESOLVE DIRECTIVE --- //
//
/**
 *
 * @param parsedSchema
 * @param objectType
 * @param newResolveDirectiveName
 * @param field
 * @param currentExSchema
 * @returns The updated schema object.
 */
const addResolveDirectiveToField = (
  parsedSchema: SupportedDocumentNode,
  objectType: string,
  newResolveDirectiveName: string,
  field: FieldDefinitionNode,
  currentExSchema: ExecutableSchema
) => {
  //
  // Traverse the parsed schema.
  // visit() does a depth first search, and we can return new nodes
  // in each enter() function to replace them. A new AST is returned.
  var newSchema = visit(parsedSchema, {
    enter(node, key, parent, path, ancestors) {
      // @return
      //   undefined: no action
      //   false: skip visiting this node
      //   visitor.BREAK: stop visiting altogether
      //   null: delete this node
      //   any value: replace this node with the returned value
      //
      // At this object type definition.
      if (
        node.kind === Kind.OBJECT_TYPE_DEFINITION &&
        node.name.value === objectType
      )
        return visit(node, {
          enter(node, key, parent, path, ancestors) {
            //
            // At this field.
            if (
              node.kind === Kind.FIELD_DEFINITION &&
              node.name.value === field.name.value
            )
              // Replace the field, adding in the resolve directive.
              return {
                ...field,
                directives: [
                  ...(field.directives ?? []),
                  {
                    kind: Kind.DIRECTIVE,
                    name: {
                      kind: Kind.NAME,
                      value: 'resolve',
                    },
                    arguments: [
                      {
                        kind: Kind.ARGUMENT,
                        name: {
                          kind: Kind.NAME,
                          value: 'name',
                        },
                        value: {
                          kind: Kind.STRING,
                          value: newResolveDirectiveName,
                        },
                      },
                    ],
                  },
                ],
              } as FieldDefinitionNode;
          },
        });
    },
  });
  //
  // Serialize the newSchema that we just made, and set that as the schema definition.
  const newSchemaString = print(newSchema);
  currentExSchema.setSchemaDefinition(newSchemaString);
  return newSchema;
};

//
// --- REMOVE RESOLVE DIRECTIVE --- //
//
/**
 *
 * @param parsedSchema
 * @param objectType
 * @param resolutionName
 * @param field
 * @param currentExSchema
 * @returns The updated schema object.
 */
const removeResolveDirectiveFromField = (
  parsedSchema: SupportedDocumentNode,
  objectType: string,
  resolutionName: string,
  field: FieldDefinitionNode,
  currentExSchema: ExecutableSchema
) => {
  //
  // Traverse the parsed schema.
  // visit() does a depth first search, and we can return new nodes
  // in each enter() function to replace them. A new AST is returned.
  var newSchema = visit(parsedSchema, {
    enter(node, key, parent, path, ancestors) {
      // @return
      //   undefined: no action
      //   false: skip visiting this node
      //   visitor.BREAK: stop visiting altogether
      //   null: delete this node
      //   any value: replace this node with the returned value
      //
      // At this object type definition.
      if (
        node.kind === Kind.OBJECT_TYPE_DEFINITION &&
        node.name.value === objectType
      )
        return visit(node, {
          enter(node, key, parent, path, ancestors) {
            //
            // At this field.
            if (
              node.kind === Kind.FIELD_DEFINITION &&
              node.name.value === field.name.value
            )
              return visit(node, {
                enter(node, key, parent, path, ancestors) {
                  //
                  // Return null to delete the resolve directive.
                  if (
                    node.kind === Kind.DIRECTIVE &&
                    node.name.value === 'resolve' &&
                    node.arguments?.find(
                      a =>
                        a.name.value === 'name' &&
                        a.value.kind === Kind.STRING &&
                        a.value.value === resolutionName
                    ) !== undefined
                  )
                    return null;
                },
              });
          },
        });
    },
  });
  //
  // Serialize the newSchema that we just made, and set that as the schema definition.
  const newSchemaString = print(newSchema as ASTNode);
  currentExSchema.setSchemaDefinition(newSchemaString);
};

//
// --- UDPATE SCHEMA AND RESOLUTION MAP --- //
//
/**
 * Updates an executable schema and resolution map:
 * - First, it updates the field in the schema, adding or removing the resolver directive on the field.
 * - Second, it updates the resolutions map, adding, updating, or removing the resolution for the field's resolver directive.
 * @param resolverItem
 * @param newResolution
 * @param currentExSchema
 * @param currResolMap
 * @param shouldDelete
 * @returns
 */
export const updateSchemaAndResolutionMap = (
  resolverItem: ResolverItem,
  newResolution: Resolution,
  currentExSchema: ExecutableSchema,
  currResolMap: jspb.Map<string, Resolution>,
  shouldDelete?: boolean
) => {
  const invalidUpdate = (msg?: string) => {
    const message =
      'Error while updating schema and resolution map' +
      (msg ? ': ' + msg : '.');
    console.error(message);
    throw new Error(message);
  };
  const { objectType, field, isNewResolution } = resolverItem;
  if (!objectType || !field)
    return invalidUpdate('Object type and field name must be supplied');
  const fieldName = resolverItem.field.name.value;
  //
  // Get the parsed values
  let currentSchemaDef = currentExSchema.getSchemaDefinition();
  const parsedSchema = parseSchemaString(currentSchemaDef);
  //
  // Perform the action to the schema and resolution map.
  if (!shouldDelete) {
    if (isNewResolution) {
      //
      // --- ADD --- //
      //
      // Generate a resolve directive name (this just has to be unique).
      const newResolveDirectiveName = `${objectType}|${fieldName}`;
      // Update the schema.
      addResolveDirectiveToField(
        parsedSchema,
        objectType,
        newResolveDirectiveName,
        field,
        currentExSchema
      );
      // Update the resolution map.
      currResolMap.set(newResolveDirectiveName, newResolution);
    } else {
      const existingResolveDirectiveName = getResolveDirectiveName(field);
      if (!existingResolveDirectiveName)
        return invalidUpdate('The resolution to update does not have a name.');
      //
      // --- UPDATE --- //
      //
      // Update the resolution map.
      currResolMap.set(existingResolveDirectiveName, newResolution);
    }
  } else {
    if (isNewResolution)
      return invalidUpdate('The resolution to delete does not exist.');
    //
    // --- DELETE --- //
    //
    // Get the resolve directive name to delete.
    const existingResolveDirectiveName = getResolveDirectiveName(field);
    if (!existingResolveDirectiveName)
      return invalidUpdate('The resolution to delete does not have a name.');
    //
    // Update the schema with the removed resolve directive.
    removeResolveDirectiveFromField(
      parsedSchema,
      objectType,
      existingResolveDirectiveName,
      field,
      currentExSchema
    );
    //
    // And delete the resolution from the resolution map if it exists (which it should).
    if (currResolMap.has(existingResolveDirectiveName))
      currResolMap.del(existingResolveDirectiveName);
    else
      console.warn(
        `An @resolve directive was found with the name "${existingResolveDirectiveName}", but this value ` +
          `was not found in the resolutions map. This @resolve directive has been removed from ` +
          `the schema for the field, "${fieldName}".`
      );
  }
};
