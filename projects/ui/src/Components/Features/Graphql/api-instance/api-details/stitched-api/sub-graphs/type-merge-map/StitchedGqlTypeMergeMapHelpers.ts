import { StitchedSchema } from 'proto/github.com/solo-io/solo-apis/api/gloo/graphql.gloo/v1beta1/graphql_pb';
import { arrayMapToObject, objectToArrayMap } from 'utils/graphql-helpers';
import YAML from 'yaml';

export type ParsedTypeMergeMap = [
  string,
  StitchedSchema.SubschemaConfig.TypeMergeConfig.AsObject
][];
export type TypeMergeMapStringFormat = {
  typeName: string;
  typeMergeConfig: string;
}[];

/**
 * @param typeMergeMap
 * @returns The Type Merge Map but with serialized configs.
 */
export const typeMergeMapToStringFormat = (
  typeMergeMap: ParsedTypeMergeMap
) => {
  const newTMMap = [] as TypeMergeMapStringFormat;
  typeMergeMap.forEach(mapping => {
    const typeName = mapping[0];
    const parsedMergeConfig = mapping[1];
    // TODO: Fix argsmap > args naming.
    // Convert the args array to an object.
    if (
      parsedMergeConfig.argsMap !== undefined &&
      parsedMergeConfig.argsMap.length > 0
    ) {
      parsedMergeConfig.argsMap = arrayMapToObject<any>(
        parsedMergeConfig.argsMap
      );
    }
    // Stringify the config to show in the text editor.
    YAML.scalarOptions.null.nullStr = '';
    const typeMergeConfig = YAML.stringify(parsedMergeConfig);
    newTMMap.push({ typeName, typeMergeConfig });
  });
  return newTMMap;
};

/**
 * @param typeMergeMapStringFormat The Type Merge Map with serialized configs.
 * @returns The parsed Type Merge Map object.
 */
export const typeMergeMapFromStringFormat = (
  typeMergeMapStringFormat: TypeMergeMapStringFormat
) => {
  let parsedMap = [] as ParsedTypeMergeMap;
  for (let i = 0; i < typeMergeMapStringFormat.length; i++) {
    const { typeName, typeMergeConfig } = typeMergeMapStringFormat[i];
    let parsedMergeConfig: any;
    try {
      parsedMergeConfig = YAML.parse(typeMergeConfig);
      // TODO: Fix argsmap > args naming.
      if (parsedMergeConfig.argsMap)
        parsedMergeConfig.argsMap = objectToArrayMap(parsedMergeConfig.argsMap);
      parsedMap.push([typeName, parsedMergeConfig]);
    } catch (err) {
      throw new Error(`${typeName}: ${(err as any).message}`);
    }
  }
  return parsedMap;
};

/**
 * @param parsedMap
 * @returns True if it is valid, otherwise it will throw an error with the validation message.
 */
export const validateTypeMergeMap = (parsedMap: ParsedTypeMergeMap) => {
  parsedMap.forEach(m => {
    const parsedMergeConfig = m[1];
    const configKeys = Object.keys(parsedMergeConfig);
    // TODO: Fix argsmap > args naming.
    if (
      configKeys.length !== 3 ||
      !configKeys.includes('argsMap') ||
      !configKeys.includes('queryName') ||
      !configKeys.includes('selectionSet')
    )
      throw new Error(
        `${m[0]}): Must include values for 'argsMap', 'queryName', and 'selectionSet' only.`
      );
    if (parsedMergeConfig.argsMap === null) parsedMergeConfig.argsMap = [];
    else if (!parsedMergeConfig.argsMap.indexOf)
      throw new Error(`${m[0]}: Must include a valid 'argsMap'.`);
    if (typeof parsedMergeConfig.queryName !== 'string')
      throw new Error(`${m[0]}: Must include a valid 'queryName'.`);
    if (typeof parsedMergeConfig.selectionSet !== 'string')
      throw new Error(`${m[0]}: Must include a valid 'selectionSet'.`);
  });
  return true;
};
