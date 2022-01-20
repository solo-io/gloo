import { FailoverSchemeApiClient } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/failover_scheme_pb_service';
import {
  host,
  getObjectRefClassFromRefObj,
  getClusterRefClassFromClusterRefObj,
} from './helpers';
import { grpc } from '@improbable-eng/grpc-web';
import {
  GetFailoverSchemeRequest,
  GetFailoverSchemeYamlRequest,
  FailoverScheme,
} from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/failover_scheme_pb';
import {
  ObjectRef,
  ClusterObjectRef,
} from 'proto/github.com/solo-io/skv2/api/core/v1/core_pb';

const failoverSchemeApiClient = new FailoverSchemeApiClient(host, {
  transport: grpc.CrossBrowserHttpTransport({ withCredentials: false }),
  debug: true,
});

export const failoverSchemeApi = {
  getFailoverScheme,
  getFailoverSchemeYAML,
};

function getFailoverScheme(
  upstreamRef: ClusterObjectRef.AsObject
): Promise<FailoverScheme.AsObject> {
  let request = new GetFailoverSchemeRequest();
  request.setUpstreamRef(getClusterRefClassFromClusterRefObj(upstreamRef));

  return new Promise((resolve, reject) => {
    failoverSchemeApiClient.getFailoverScheme(request, (error, data) => {
      if (error !== null) {
        console.error('Error:', error.message);
        console.error('Code:', error.code);
        console.error('Metadata:', error.metadata);
        reject(error);
      } else {
        resolve(data!.toObject().failoverScheme!);
      }
    });
  });
}

function getFailoverSchemeYAML(objectRef: ObjectRef.AsObject): Promise<string> {
  let request = new GetFailoverSchemeYamlRequest();
  request.setFailoverSchemeRef(getObjectRefClassFromRefObj(objectRef));

  return new Promise((resolve, reject) => {
    failoverSchemeApiClient.getFailoverSchemeYaml(request, (error, data) => {
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
