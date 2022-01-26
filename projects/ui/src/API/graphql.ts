import { getObjectRefClassFromRefObj, host } from './helpers';
import { grpc } from '@improbable-eng/grpc-web';

import { ObjectRef } from 'proto/github.com/solo-io/skv2/api/core/v1/core_pb';
import petstoreSchema from '../Components/Features/Graphql/data/petstore.json';
import bookinfoSchema from '../Components/Features/Graphql/data/book-info.json';
import { bookInfoYaml } from '../Components/Features/Graphql/data/book-info-yaml';
import { petstoreYaml } from '../Components/Features/Graphql/data/petstore-yaml';
const graphqlApiClient = 'TODO';

export type GraphqlSchemaType = {
  metadata: typeof petstoreSchema.metadata | typeof bookinfoSchema.metadata;
  spec: typeof petstoreSchema.spec | typeof bookinfoSchema.spec;
  yamlConfig: string;
};

export type ResolutionType =
  | typeof petstoreSchema.spec.executableSchema.executor.local.resolutions
  | typeof bookinfoSchema.spec.executableSchema.executor.local.resolutions;

let petstoreResolutionArr = Object.entries(
  petstoreSchema.spec.executableSchema.executor.local.resolutions
);

let bookinfoResolutionArr = Object.entries(
  bookinfoSchema.spec.executableSchema.executor.local.resolutions
);

export type ResolutionMapType = [
  string,
  {
    restResolver: {
      request: {
        body?:
          | string
          | {
              [key: string]: string;
            };
        headers?: {
          ':method': string;
          ':path': string;
        };
      };
      response?: {
        [key: string]: string | { [key: string]: string };
      };
      upstreamRef: ObjectRef.AsObject;
    };
  }
][];

export const graphqlApi = {
  listGraphqlSchemas,
  getGraphqlSchema,
};

function listGraphqlSchemas(): Promise<GraphqlSchemaType[]> {
  return new Promise((resolve, reject) => {
    resolve([
      {
        metadata: petstoreSchema.metadata,
        spec: petstoreSchema.spec,
        yamlConfig: petstoreYaml,
      },
      {
        metadata: bookinfoSchema.metadata,
        spec: bookinfoSchema.spec,
        yamlConfig: bookInfoYaml,
      },
    ]);
  });
}

function getGraphqlSchema(
  graphqlSchemaRef: ObjectRef.AsObject
): Promise<GraphqlSchemaType> {
  return new Promise((resolve, reject) => {
    if (graphqlSchemaRef?.name === 'bookinfo-graphql') {
      resolve({
        metadata: bookinfoSchema.metadata,
        spec: bookinfoSchema.spec,
        yamlConfig: bookInfoYaml,
      });
    } else {
      resolve({
        metadata: petstoreSchema.metadata,
        spec: petstoreSchema.spec,
        yamlConfig: petstoreYaml,
      });
    }
  });
}
