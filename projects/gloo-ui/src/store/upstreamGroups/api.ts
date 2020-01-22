import { UpstreamGroupApiClient } from 'proto/solo-projects/projects/grpcserver/api/v1/upstreamgroup_pb_service';
import { grpc } from '@improbable-eng/grpc-web';
import { host } from 'store';
import {
  GetUpstreamGroupRequest,
  UpstreamGroupDetails,
  ListUpstreamGroupsRequest,
  ListUpstreamGroupsResponse,
  CreateUpstreamGroupRequest,
  CreateUpstreamGroupResponse,
  UpdateUpstreamGroupRequest,
  UpdateUpstreamGroupResponse,
  DeleteUpstreamGroupRequest,
  DeleteUpstreamGroupResponse
} from 'proto/solo-projects/projects/grpcserver/api/v1/upstreamgroup_pb';
import { ResourceRef } from 'proto/solo-kit/api/v1/ref_pb';
import {
  UpstreamGroup,
  WeightedDestination,
  Destination
} from 'proto/gloo/projects/gloo/api/v1/proxy_pb';
import { Metadata } from 'proto/solo-kit/api/v1/metadata_pb';

export const client = new UpstreamGroupApiClient(host, {
  transport: grpc.CrossBrowserHttpTransport({ withCredentials: false }),
  debug: true
});

// get
function getUpstreamGroup(
  getUpstreamGroupRequest: GetUpstreamGroupRequest.AsObject
): Promise<UpstreamGroupDetails> {
  return new Promise((resolve, reject) => {
    let request = new GetUpstreamGroupRequest();
    let ref = new ResourceRef();
    ref.setName(getUpstreamGroupRequest.ref!.name);
    ref.setNamespace(getUpstreamGroupRequest.ref!.namespace);
    request.setRef(ref);

    client.getUpstreamGroup(request, (error, data) => {
      if (error !== null) {
        console.error('Error:', error.message);
        console.error('Code:', error.code);
        console.error('Metadata:', error.metadata);
        reject(error);
      } else {
        resolve(data!.getUpstreamGroupDetails());
      }
    });
  });
}

// list
function listUpstreamGroups(): Promise<ListUpstreamGroupsResponse.AsObject> {
  return new Promise((resolve, reject) => {
    let request = new ListUpstreamGroupsRequest();

    client.listUpstreamGroups(request, (error, data) => {
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

function setWeightedDestinationValues(
  weightedDest: WeightedDestination.AsObject,
  weightedDestToUpdate = new WeightedDestination()
) {
  let { destination, weight, options } = weightedDest!;
  if (destination !== undefined) {
    let newDest = new Destination();
    let { upstream, kube, consul, destinationSpec, subset } = destination!;

    if (upstream !== undefined) {
      let upstreamRef = new ResourceRef();
      upstreamRef.setName(upstream.name);
      upstreamRef.setNamespace(upstream.namespace);
      newDest.setUpstream(upstreamRef);
    }

    // TODO
    if (kube !== undefined) {
      // let newKubeDest = new KubernetesServiceDestination()
    }

    // TODO
    if (consul !== undefined) {
    }

    // TODO
    if (destinationSpec !== undefined) {
    }
    weightedDestToUpdate.setDestination(newDest);
  }

  if (weight !== undefined) {
    weightedDestToUpdate.setWeight(weight);
  }

  // TODO
  if (options !== undefined) {
  }
  return weightedDestToUpdate;
}

function setUpstreamGroupValues(
  upstreamGroup: UpstreamGroup.AsObject,
  upstreamGroupToUpdate = new UpstreamGroup()
): UpstreamGroup {
  let { destinationsList, metadata, status } = upstreamGroup!;
  if (metadata !== undefined) {
    let { name, namespace } = metadata;
    let newMetadata = new Metadata();
    newMetadata.setName(name);
    newMetadata.setNamespace(namespace);
    upstreamGroupToUpdate.setMetadata(newMetadata);
  }
  if (destinationsList !== undefined) {
    let newDestinationsList = destinationsList.map(weightedDest => {
      let newDest = setWeightedDestinationValues(weightedDest);
      return newDest;
    });
    upstreamGroupToUpdate.setDestinationsList(newDestinationsList);
  }

  return upstreamGroupToUpdate;
}

// create
function createUpstreamGroup(
  createUpstreamGroupRequest: CreateUpstreamGroupRequest.AsObject
): Promise<CreateUpstreamGroupResponse.AsObject> {
  return new Promise((resolve, reject) => {
    let request = new CreateUpstreamGroupRequest();
    let { upstreamGroup } = createUpstreamGroupRequest!;

    if (upstreamGroup !== undefined) {
      let inputUpstreamGroup = setUpstreamGroupValues(upstreamGroup);
      request.setUpstreamGroup(inputUpstreamGroup);
    }

    client.createUpstreamGroup(request, (error, data) => {
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

// update
function updateUpstreamGroup(
  updateUpstreamGroupRequest: UpdateUpstreamGroupRequest.AsObject
): Promise<UpdateUpstreamGroupResponse.AsObject> {
  return new Promise(async (resolve, reject) => {
    let request = new UpdateUpstreamGroupRequest();
    let { upstreamGroup } = updateUpstreamGroupRequest;
    if (upstreamGroup !== undefined && upstreamGroup.metadata !== undefined) {
      let { name, namespace } = upstreamGroup.metadata;
      let upstreamGroupToUpdate = await getUpstreamGroup({
        ref: { name, namespace }
      });
      let updatedUpstreamGroup = setUpstreamGroupValues(
        upstreamGroup,
        upstreamGroupToUpdate.getUpstreamGroup()
      );
      request.setUpstreamGroup(updatedUpstreamGroup);
    }
    client.updateUpstreamGroup(request, (error, data) => {
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

// delete
function deleteUpstreamGroup(
  deleteUpstreamGroupRequest: DeleteUpstreamGroupRequest.AsObject
): Promise<DeleteUpstreamGroupResponse> {
  return new Promise((resolve, reject) => {
    let { ref } = deleteUpstreamGroupRequest!;
    let request = new DeleteUpstreamGroupRequest();
    let upstreamGroupRef = new ResourceRef();
    upstreamGroupRef.setName(ref!.name);
    upstreamGroupRef.setNamespace(ref!.namespace);
    request.setRef(upstreamGroupRef);

    client.deleteUpstreamGroup(request, (error, data) => {
      if (error !== null) {
        console.error('Error:', error.message);
        console.error('Code:', error.code);
        console.error('Metadata:', error.metadata);
        reject(error);
      } else {
        resolve(data!);
      }
    });
  });
}
export const upstreamGroupAPI = {
  listUpstreamGroups,
  createUpstreamGroup,
  updateUpstreamGroup,
  deleteUpstreamGroup
};
