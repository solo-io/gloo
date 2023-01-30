#!/usr/local/bin/node
/**
 * The following methods are created by v8go. getInput() returnOutput(value)
 * v8go will use these to communicate the input and output from this JS file.
 * 
 * The Stitching job will stitch GraphQL Schemas together
 */
import 'core-js/actual';
const graphql = require('graphql');
const { stitchSchemas } = require('@graphql-tools/stitch');
const { makeExecutableSchema } = require('@graphql-tools/schema');
const { msgToBase64String, base64StringToMsg } = require('./conversion');

const {
  GraphQLToolsStitchingInput,
  GraphQLToolsStitchingOutput,
} = require("../../../../../ui/src/proto/github.com/solo-io/solo-projects/projects/gloo/api/enterprise/graphql/v1/stitching_info_pb");
const {
  MergedTypeConfig,
  FieldNodeMap,
  FieldNodes,
  FieldNode,
  Schemas,
} = require('../../../../../ui/src/proto/github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/extensions/graphql/stitching_pb')


function makeStitchedSchema(input) {
  const inMsg = base64StringToMsg(
    input,
    GraphQLToolsStitchingInput.deserializeBinary
  );
  let subschemas = [];
  let subschemaList = inMsg.getSubschemasList();
  for (let i = 0; i < subschemaList.length; i++) {
    let subschemaConfig = {};
    let subschema = subschemaList[i];
    if (!subschema.getSchema()) {
      throw new Error(
        `subschema '${subschema.getName()}' is missing a schema definition`
      );
    }
    subschemaConfig.schema = makeExecutableSchema({
      typeDefs: subschema.getSchema(),
    });
    subschemaConfig.name = subschema.getName();
    subschemaConfig.merge = {};
    subschema.getTypeMergeConfigMap().forEach((typeMergeConfig, typeName) => {
      subschemaConfig.merge[typeName] = {
        selectionSet: typeMergeConfig.getSelectionSet(),
        fieldName: typeMergeConfig.getFieldName(),
        // this doesn't do anything but is needed for type merging to be carried out by graphql-tools
        args: () => ({}),
      };
    });

    subschemas.push(subschemaConfig);
  }

  return stitchSchemas({
    subschemas: subschemas,
    //todo(sai) - make this configurable
    mergeTypes: true,
  });
}

function convertStitchingInfo(si) {
  let newSi = new GraphQLToolsStitchingOutput();
  let subschemaMap = new Map();
  for (let [key, value] of si.subschemaMap) {
    subschemaMap.set(value, key.name);
  }
  let fieldNodesByField = newSi.getFieldNodesByFieldMap();
  for (let fieldNodeByFieldIdx in si.fieldNodesByField) {
    let newFieldNodeMap = new FieldNodeMap();
    for (let fieldNodeTypeIdx in si.fieldNodesByField[fieldNodeByFieldIdx]) {
      let newFieldNodeList = new FieldNodes();
      for (
          let i = 0;
          i < si.fieldNodesByField[fieldNodeByFieldIdx][fieldNodeTypeIdx].length;
          i++
      ) {
        let newFieldNode = new FieldNode();
        newFieldNode.setName(
            si.fieldNodesByField[fieldNodeByFieldIdx][fieldNodeTypeIdx][i].name.value
        );
        newFieldNodeList.addFieldNodes(newFieldNode);
      }
      newFieldNodeMap.getNodesMap().set(fieldNodeTypeIdx, newFieldNodeList);
    }
    fieldNodesByField.set(fieldNodeByFieldIdx, newFieldNodeMap);
  }
  let fieldNodesByType = newSi.getFieldNodesByTypeMap();
  for (let k in si.fieldNodesByType) {
    let newFieldNodeList = new FieldNodes();
    for (let i = 0; i < si.fieldNodesByType[k].length; i++) {
      let newFieldNode = new FieldNode();
      newFieldNode.setName(si.fieldNodesByType[k][i].name.value);
      newFieldNodeList.addFieldNodes(newFieldNode);
    }
    fieldNodesByType.set(k, newFieldNodeList);
  }
  let mergedTypesMap = newSi.getMergedTypesMap();
  for (let type in si.mergedTypes) {
    let mt = si.mergedTypes[type];
    let a = si.mergedTypes[type].selectionSets;
    let newMtConfig = new MergedTypeConfig();
    newMtConfig.setTypeName(type);
    for (let [subschema, selectionSet] of mt.selectionSets) {
      newMtConfig
        .getSelectionSetsMap()
        .set(subschemaMap.get(subschema), graphql.print(selectionSet));
    }
    for (let [subschema, targetSubschemas] of mt.targetSubschemas) {
      let newTargetSubschemas = new Schemas();
      for (let targetSubschema of targetSubschemas) {
        newTargetSubschemas.addSchemas(subschemaMap.get(targetSubschema));
      }
      newMtConfig
        .getDeclarativeTargetSubschemasMap()
        .set(subschemaMap.get(subschema), newTargetSubschemas);
    }

    for (let [fieldName, subschema] of Object.entries(mt.uniqueFields)) {
      newMtConfig
        .getUniqueFieldsToSubschemaNameMap()
        .set(fieldName, subschemaMap.get(subschema));
    }
    mergedTypesMap.set(type, newMtConfig);
  }

  return newSi;
}

// getInput is defined in the go control plane with g8go
let schema = makeStitchedSchema(getInput());
let newSi = convertStitchingInfo(schema.extensions.stitchingInfo);
newSi.setStitchedSchema(graphql.printSchema(schema));
// returnOutput is defined in the go control plane and injected via v8go
returnOutput(msgToBase64String(newSi));
