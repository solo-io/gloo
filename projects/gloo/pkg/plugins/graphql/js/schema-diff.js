#!/usr/local/bin/node

const { diff, CriticalityLevel } = require('@graphql-inspector/core');
const { makeExecutableSchema } = require('@graphql-tools/schema');
const { msgToBase64String, base64StringToMsg } = require('./conversion');

// In order for graphql-inspector to be able to validate schemas with custom directives,
// we need to include the directive definitions in each schema.
const directiveDefinitions = `
directive @resolve(name: String!) on FIELD_DEFINITION

enum CacheControl {
  unset
  private
  public
}

directive @cacheControl(
  maxAge: Int
  inheritMaxAge: Boolean
  scope: CacheControl = unset
) on FIELD_DEFINITION | OBJECT

`;

const protoImportPath = process.env.GRAPHQL_PROTO_ROOT;
if (!protoImportPath) {
  console.error(
    'schema diff script requires GRAPHQL_PROTO_ROOT environment variable'
  );
  process.exit(1);
}

const {
  GraphQLInspectorDiffInput,
  GraphQLInspectorDiffOutput,
} = require(protoImportPath +
  'github.com/solo-io/solo-projects/projects/gloo/api/enterprise/graphql/v1/diff_pb');

// the input to this script should be deserializable into a GraphQLInspectorDiffInput message
const input = process.argv[2];
const inMsg = base64StringToMsg(
  input,
  GraphQLInspectorDiffInput.deserializeBinary
);
const oldSchemaStr = directiveDefinitions + inMsg.getOldSchema();
const newSchemaStr = directiveDefinitions + inMsg.getNewSchema();

const oldSchema = makeExecutableSchema({
  typeDefs: oldSchemaStr,
});
const newSchema = makeExecutableSchema({
  typeDefs: newSchemaStr,
});

diff(oldSchema, newSchema)
  .then(changes => {
    // sort by criticality level (breaking -> non-dangerous -> breaking), then by path, then by change type
    changes.sort((change1, change2) => {
      const criticalityDiff =
        convertCriticalityLevelToMsg(change2.criticality?.level) -
        convertCriticalityLevelToMsg(change1.criticality?.level);
      if (criticalityDiff != 0) {
        return criticalityDiff;
      }

      const pathCompare = (change1.path ?? '').localeCompare(
        change2.path ?? ''
      );
      if (pathCompare != 0) {
        return pathCompare;
      }

      return (change1.type ?? '').localeCompare(change2.type ?? '');
    });

    const output = new GraphQLInspectorDiffOutput();
    const changeMsgs = changes.map(convertChangeToMsg);
    output.setChangesList(changeMsgs);

    // This is the stdout output that the control plane reads to get the GraphQLInspectorDiffOutput message
    console.log(msgToBase64String(output));
  })
  .catch(err => console.error(err));

// message conversion functions
function convertChangeToMsg(change) {
  const changeMsg = new GraphQLInspectorDiffOutput.Change();
  changeMsg.setMessage(change.message);
  changeMsg.setPath(change.path);
  changeMsg.setChangeType(change.type);
  changeMsg.setCriticality(convertCriticalityToMsg(change.criticality));
  return changeMsg;
}

function convertCriticalityToMsg(criticality) {
  const criticalityMsg = new GraphQLInspectorDiffOutput.Criticality();
  criticalityMsg.setLevel(convertCriticalityLevelToMsg(criticality?.level));
  criticalityMsg.setReason(criticality?.reason);
  return criticalityMsg;
}

function convertCriticalityLevelToMsg(criticalityLevel) {
  switch (criticalityLevel) {
    case CriticalityLevel.Breaking:
      return GraphQLInspectorDiffOutput.CriticalityLevel.BREAKING;
    case CriticalityLevel.Dangerous:
      return GraphQLInspectorDiffOutput.CriticalityLevel.DANGEROUS;
    case CriticalityLevel.NonBreaking:
    default:
      return GraphQLInspectorDiffOutput.CriticalityLevel.NON_BREAKING;
  }
}
