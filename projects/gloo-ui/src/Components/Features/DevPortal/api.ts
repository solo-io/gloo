import { Empty } from 'google-protobuf/google/protobuf/empty_pb';
import { grpc } from '@improbable-eng/grpc-web';
import { PortalApi } from 'proto/dev-portal/api/grpc/admin/portal_pb_service';
import { host } from 'store';
import { Portal, PortalList } from 'proto/dev-portal/api/grpc/admin/portal_pb';
 
export const DevPortalApi = {
  listPortals
};

function listPortals(): Promise<Portal.AsObject[]> {
  return new Promise((resolve, reject) => {
    grpc.invoke(PortalApi.ListPortals, {
      request: new Empty(),
      host,
      metadata: new grpc.Metadata(),
      onHeaders: (headers: grpc.Metadata) => {
        // console.log('onheaders', headers);
      },
      onMessage: (message: PortalList) => {
        // console.log('message', message);
        if (message) {
          resolve(message.toObject().portalsList);
        }
      },
      onEnd: (
        status: grpc.Code,
        statusMessage: string,
        trailers: grpc.Metadata
      ) => {
        // console.log('onEnd', status, statusMessage, trailers);
        if (status !== grpc.Code.OK) {
          reject(statusMessage);
        }
      }
    });
  });
}
