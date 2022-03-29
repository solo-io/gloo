import { BootstrapApiClient } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/bootstrap_pb_service';
import { host } from './helpers';
import { grpc } from '@improbable-eng/grpc-web';
import {
  GetConsoleOptionsRequest,
  GetConsoleOptionsResponse,
  GlooFedCheckRequest,
  GlooFedCheckResponse,
  GraphqlCheckRequest,
} from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/bootstrap_pb';
import { ServiceError } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/gateway_resources_pb_service';

const bootstrapApiClient = new BootstrapApiClient(host, {
  transport: grpc.CrossBrowserHttpTransport({ withCredentials: false }),
  debug: true,
});

export const bootstrapApi = {
  isGlooFedEnabled,
  isGraphqlEnabled,
  getConsoleOptions,
};

function isGlooFedEnabled(): Promise<GlooFedCheckResponse.AsObject> {
  const request = new GlooFedCheckRequest();
  return new Promise((resolve, reject) => {
    bootstrapApiClient.isGlooFedEnabled(request, (error, data) => {
      if (error !== null) {
        console.error('Error:', error.message);
        console.error('Code:', error.code);
        console.error('Metadata:', error.metadata);
        reject(error);
      } else {
        resolve(data!.toObject());
      }
    });
  });
}

function isGraphqlEnabled(): Promise<boolean> {
  const request = new GraphqlCheckRequest();
  return new Promise((resolve, reject) => {
    bootstrapApiClient.isGraphqlEnabled(request, (error, data) => {
      if (error !== null) {
        console.error('Error:', error.message);
        console.error('Code:', error.code);
        console.error('Metadata:', error.metadata);
        reject(error);
      } else {
        resolve(data!.toObject().enabled);
      }
    });
  });
}

function getConsoleOptions(): Promise<GetConsoleOptionsResponse.AsObject | null> {
  const message = new GetConsoleOptionsRequest();
  return new Promise((resolve, reject) => {
    return bootstrapApiClient.getConsoleOptions(
      message,
      (err, responseMessage) => {
        if (err) {
          reject(err.message);
        } else {
          if (responseMessage) {
            resolve(responseMessage!.toObject());
          } else {
            resolve(null);
          }
        }
      }
    );
  });
}
