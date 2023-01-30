#!/usr/local/bin/node
/**
 * The following methods are created by v8go. getInput() returnOutput(value)
 * v8go will use these to communicate the input and output from this JS file.
 * 
 * The schema diff will report on the differences of the schema's of an original and new schema.
 */
import 'core-js/actual';
const { diff, DiffRule, CriticalityLevel } = require('@graphql-inspector/core');
const { makeExecutableSchema } = require('@graphql-tools/schema');
const { msgToBase64String, base64StringToMsg } = require('./conversion');
const {
  GraphQLInspectorDiffInput,
  GraphQLInspectorDiffOutput,
} = require('../../../../../ui/src/proto/github.com/solo-io/solo-projects/projects/gloo/api/enterprise/graphql/v1/diff_pb');
const { GraphqlOptions } = require('../../../../../ui/src/proto/github.com/solo-io/solo-apis/api/gloo/gloo/v1/settings_pb');


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


// getInput is defined in the control plane by v8go
var input = getInput();
const inMsg = base64StringToMsg(
  input,
  GraphQLInspectorDiffInput.deserializeBinary
);
const oldSchemaStr = directiveDefinitions + inMsg.getOldSchema();
const newSchemaStr = directiveDefinitions + inMsg.getNewSchema();
const ruleEnums = inMsg.getRulesList();

// convert to diff function input
const oldSchema = makeExecutableSchema({
  typeDefs: oldSchemaStr,
});
const newSchema = makeExecutableSchema({
  typeDefs: newSchemaStr,
});
const rules = ruleEnums?.length ? ruleEnums.map(convertRuleFromMsg) : [];

diff(oldSchema, newSchema, rules)
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
    // returnOutput is defined in the control plane by v8go
    returnOutput(msgToBase64String(output));
  })
  .catch(err => {
    throw new Error(err)
  });

// message conversion functions
function convertRuleFromMsg(ruleMsg) {
  switch (ruleMsg) {
    case GraphqlOptions.SchemaChangeValidationOptions.ProcessingRule
      .RULE_DANGEROUS_TO_BREAKING:
      return DiffRule.dangerousBreaking;
    case GraphqlOptions.SchemaChangeValidationOptions.ProcessingRule
      .RULE_DEPRECATED_FIELD_REMOVAL_DANGEROUS:
      return DiffRule.suppressRemovalOfDeprecatedField;
    case GraphqlOptions.SchemaChangeValidationOptions.ProcessingRule
      .RULE_IGNORE_DESCRIPTION_CHANGES:
      return DiffRule.ignoreDescriptionChanges;
    // TODO support RULE_IGNORE_UNREACHABLE -> safeUnreachable when https://github.com/kamilkisiela/graphql-inspector/issues/2063 is fixed
    // tracking issue here https://github.com/solo-io/solo-projects/issues/3853
    default:
      throw new Error('unexpected rule: ' + ruleMsg)
  }
}

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
