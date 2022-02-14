import {
  getClusterRefClassFromClusterRefObj,
  getObjectRefClassFromRefObj,
  host,
} from './helpers';
import { grpc } from '@improbable-eng/grpc-web';

import {
  ClusterObjectRef,
  ObjectRef,
} from 'proto/github.com/solo-io/skv2/api/core/v1/core_pb';
import petstoreSchema from '../Components/Features/Graphql/data/petstore.json';
import bookinfoSchema from '../Components/Features/Graphql/data/book-info.json';
import { bookInfoYaml } from '../Components/Features/Graphql/data/book-info-yaml';
import { petstoreYaml } from '../Components/Features/Graphql/data/petstore-yaml';
import { GraphqlApiClient } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/graphql_pb_service';
import {
  GetGraphqlSchemaRequest,
  GraphqlSchema,
  ListGraphqlSchemasRequest,
  GetGraphqlSchemaYamlRequest,
} from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/graphql_pb';
const graphqlApiClient = new GraphqlApiClient(host, {
  transport: grpc.CrossBrowserHttpTransport({ withCredentials: false }),
  debug: true,
});

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
  getGraphqlSchemaYaml,
};

function listGraphqlSchemas(
  glooInstanceRef?: ObjectRef.AsObject
): Promise<GraphqlSchema.AsObject[]> {
  let request = new ListGraphqlSchemasRequest();
  if (glooInstanceRef) {
    request.setGlooInstanceRef(getObjectRefClassFromRefObj(glooInstanceRef));
  }
  return new Promise((resolve, reject) => {
    graphqlApiClient.listGraphqlSchemas(request, (error, data) => {
      if (error !== null) {
        console.error('Error:', error.message);
        console.error('Code:', error.code);
        console.error('Metadata:', error.metadata);
        reject(error);
      } else {
        resolve(data!.toObject().graphqlSchemasList);
      }
    });
  });
}

function getGraphqlSchema(
  graphqlSchemaRef: ClusterObjectRef.AsObject
): Promise<GraphqlSchema.AsObject> {
  let request = new GetGraphqlSchemaRequest();
  request.setGraphqlSchemaRef(
    getClusterRefClassFromClusterRefObj(graphqlSchemaRef)
  );

  return new Promise((resolve, reject) => {
    graphqlApiClient.getGraphqlSchema(request, (error, data) => {
      if (error !== null) {
        console.error('Error:', error.message);
        console.error('Code:', error.code);
        console.error('Metadata:', error.metadata);
        reject(error);
      } else {
        resolve(data!.toObject()!.graphqlSchema!);
      }
    });
  });
}

function getGraphqlSchemaYaml(
  graphqlSchemaRef: ClusterObjectRef.AsObject
): Promise<string> {
  let request = new GetGraphqlSchemaYamlRequest();
  request.setGraphqlSchemaRef(
    getClusterRefClassFromClusterRefObj(graphqlSchemaRef)
  );

  return new Promise((resolve, reject) => {
    graphqlApiClient.getGraphqlSchemaYaml(request, (error, data) => {
      if (error !== null) {
        console.error('Error:', error.message);
        console.error('Code:', error.code);
        console.error('Metadata:', error.metadata);
        reject(error);
      } else {
        resolve(data!.toObject().yamlData?.yaml ?? 'None');
      }
    });
  });
}
