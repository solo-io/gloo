#!/usr/local/bin/node

const graphql = require('graphql');
const { stitchSchemas } = require('@graphql-tools/stitch');
const { makeExecutableSchema } = require('@graphql-tools/schema');
const { msgToBase64String, base64StringToMsg } = require('./conversion');

const protoImportPath = process.env.GRAPHQL_PROTO_ROOT;
if (!protoImportPath) {
  console.error(
    'stitching tools script requires GRAPHQL_PROTO_ROOT environment variable'
  );
  process.exit(1);
}

const {
  GraphQLToolsStitchingInput,
  GraphQLToolsStitchingOutput,
} = require(protoImportPath +
  'github.com/solo-io/solo-projects/projects/gloo/api/enterprise/graphql/v1/stitching_info_pb');
const {
  MergedTypeConfig,
  FieldNodeMap,
  FieldNodes,
  FieldNode,
  Schemas,
} = require(protoImportPath +
  'github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/extensions/graphql/stitching_pb');

async function makeStitchedSchema(input) {
  const inMsg = base64StringToMsg(
    input,
    GraphQLToolsStitchingInput.deserializeBinary
  );
  let subschemas = [];
  let subschemaList = inMsg.getSubschemasList();
  for (let i = 0; i < subschemaList.length; i++) {
    subschemaConfig = {};
    let subschema = subschemaList[i];
    if (!subschema.getSchema()) {
      console.error(
        `subschema '${subschema.getName()}' is missing a schema definition`
      );
      process.exit(1);
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
    for (let fieldNodeType in si.fieldNodesByField[fieldNodeByFieldIdx]) {
      let newFieldNodeList = new FieldNodes();
      for (
        let i = 0;
        i < si.fieldNodesByField[fieldNodeByFieldIdx][fieldNodeType].length;
        i++
      ) {
        let newFieldNode = new FieldNode();
        newFieldNode.setName(
          si.fieldNodesByField[fieldNodeByFieldIdx][fieldNodeType][i].name.value
        );
        newFieldNodeList.addFieldNodes(newFieldNode);
      }
      newFieldNodeMap.getNodesMap().set(fieldNodeType, newFieldNodeList);
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

makeStitchedSchema(process.argv[2]).then(schema => {
  let newSi = convertStitchingInfo(schema.extensions.stitchingInfo);
  newSi.setStitchedSchema(graphql.printSchema(schema));
  let b64 = msgToBase64String(newSi);
  // This is the stdout output that the control plane reads to get the StitchingInfo message
  console.log(b64);
});
