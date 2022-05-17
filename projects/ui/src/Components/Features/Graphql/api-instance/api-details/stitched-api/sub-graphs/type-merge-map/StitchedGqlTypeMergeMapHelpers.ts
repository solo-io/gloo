import { StitchedSchema } from 'proto/github.com/solo-io/solo-apis/api/gloo/graphql.gloo/v1beta1/graphql_pb';
import YAML from 'yaml';
import {
  gqlJsonRoot,
  postUnmarshallProtoValues,
  preMarshallProtoValues,
} from '../../../schema/object-definition/resolver-wizard/converters';

export type ParsedTypeMergeMap = [
  string,
  StitchedSchema.SubschemaConfig.TypeMergeConfig.AsObject
][];
export type TypeMergeMapStringFormat = {
  typeName: string;
  typeMergeConfig: string;
}[];

/**
 * The generated type merge config JSON descriptor.
 */
const pbTypeObj = gqlJsonRoot.lookupType(
  'graphql.gloo.solo.io.StitchedSchema.SubschemaConfig.TypeMergeConfig'
);

/**
 * @param typeMergeMap
 * @returns The Type Merge Map but with serialized configs.
 */
export const typeMergeMapToStringFormat = (
  typeMergeMap: ParsedTypeMergeMap
) => {
  const postUnmarshalledMap = [] as TypeMergeMapStringFormat;
  YAML.scalarOptions.null.nullStr = '';
  typeMergeMap.forEach(([typeName, typeMergeConfigAsObject]) => {
    const parsedTmm = postUnmarshallProtoValues(
      typeMergeConfigAsObject,
      pbTypeObj.toJSON()
    );
    const typeMergeConfig = YAML.stringify(parsedTmm, { simpleKeys: true });
    postUnmarshalledMap.push({ typeName, typeMergeConfig });
  });
  return postUnmarshalledMap;
};

/**
 * @param typeMergeMapStringFormat The Type Merge Map with serialized configs.
 * @returns The parsed Type Merge Map object.
 */
export const typeMergeMapFromStringFormat = (
  typeMergeMapStringFormat: TypeMergeMapStringFormat
) => {
  const preMarshalledMap = [] as any;
  YAML.scalarOptions.null.nullStr = '';
  typeMergeMapStringFormat.forEach(({ typeName, typeMergeConfig }) => {
    let jsonConfig = YAML.parse(typeMergeConfig, { simpleKeys: true });
    const preMarshalledConfig = preMarshallProtoValues(
      jsonConfig,
      pbTypeObj.toJSON()
    );
    preMarshalledMap.push([typeName, preMarshalledConfig]);
  });
  return preMarshalledMap;
};
