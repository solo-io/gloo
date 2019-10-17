import { host } from 'store';
import { grpc } from '@improbable-eng/grpc-web';
import { RouteTableApiClient } from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/routetable_pb_service';
import {
  GetRouteTableRequest,
  GetRouteTableResponse,
  ListRouteTablesRequest,
  ListRouteTablesResponse,
  CreateRouteTableRequest,
  CreateRouteTableResponse,
  UpdateRouteTableRequest,
  UpdateRouteTableResponse,
  UpdateRouteTableYamlRequest,
  DeleteRouteTableRequest,
  DeleteRouteTableResponse,
  RouteTableDetails
} from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/routetable_pb';
import { ResourceRef } from 'proto/github.com/solo-io/solo-kit/api/v1/ref_pb';

import {
  Destination,
  HeaderMatcher,
  Matcher,
  QueryParameterMatcher,
  RouteAction
} from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/proxy_pb';
import { Parameters } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins/transformation/parameters_pb';
import { StringValue } from 'google-protobuf/google/protobuf/wrappers_pb';
import {
  RoutePlugins,
  DestinationSpec
} from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins_pb';
import { DestinationSpec as AwsDestinationSpec } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins/aws/aws_pb';
import { DestinationSpec as AzureDestinationSpec } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins/azure/azure_pb';
import { DestinationSpec as GrpcDestinationSpec } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins/grpc/grpc_pb';
import { DestinationSpec as RestDestinationSpec } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins/rest/rest_pb';
import { Route } from 'proto/github.com/solo-io/gloo/projects/gateway/api/v1/virtual_service_pb';
import { EditedResourceYaml } from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/types_pb';
import { RouteTable } from 'proto/github.com/solo-io/gloo/projects/gateway/api/v1/route_table_pb';
import { Metadata } from 'proto/github.com/solo-io/solo-kit/api/v1/metadata_pb';

const client = new RouteTableApiClient(host, {
  transport: grpc.CrossBrowserHttpTransport({ withCredentials: false }),
  debug: true
});

function getRouteTable(
  getRouteTableRequest: GetRouteTableRequest.AsObject
): Promise<GetRouteTableResponse> {
  return new Promise((resolve, reject) => {
    let { ref } = getRouteTableRequest!;
    let request = new GetRouteTableRequest();
    let routeTableRef = new ResourceRef();
    routeTableRef.setName(ref!.name);
    routeTableRef.setNamespace(ref!.namespace);
    request.setRef(routeTableRef);

    client.getRouteTable(request, (error, data) => {
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

function listRouteTables(): Promise<ListRouteTablesResponse.AsObject> {
  return new Promise((resolve, reject) => {
    let request = new ListRouteTablesRequest();
    client.listRouteTables(request, (error, data) => {
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

function createRouteTable(
  createRouteTableRequest: CreateRouteTableRequest.AsObject
): Promise<CreateRouteTableResponse.AsObject> {
  return new Promise((resolve, reject) => {
    let request = new CreateRouteTableRequest();
    let newRouteTable = new RouteTable();
    let { routeTable } = createRouteTableRequest!;
    if (routeTable !== undefined) {
      let newRoutes = routeTable.routesList.map(setInputRouteValues);
      newRouteTable.setRoutesList(newRoutes);
      if (routeTable.metadata !== undefined) {
        let newMetadata = new Metadata();
        newMetadata.setName(routeTable.metadata.name);
        newMetadata.setNamespace(routeTable.metadata.namespace);
        newRouteTable.setMetadata(newMetadata);
      }
    }

    request.setRouteTable(newRouteTable);
    client.createRouteTable(request, (error, data) => {
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

export function setInputRouteValues(route: Route.AsObject) {
  let updatedRoute = new Route();

  if (route !== undefined) {
    // matcher
    if (route.matchersList[0] !== undefined) {
      let {
        prefix,
        exact,
        regex,
        headersList,
        queryParametersList,
        methodsList
      } = route.matchersList[0]!;
      let newMatcher = new Matcher();
      if (prefix !== undefined && prefix !== '') {
        newMatcher.setPrefix(prefix);
      }
      if (exact !== undefined && exact !== '') {
        newMatcher.setExact(exact);
      }
      if (regex !== undefined && regex !== '') {
        newMatcher.setRegex(regex);
      }

      if (headersList !== undefined) {
        let newHeaderMatcherList = headersList.map(header => {
          let newHeaderMatcher = new HeaderMatcher();
          newHeaderMatcher.setName(header.value);
          newHeaderMatcher.setRegex(header.regex);
          newHeaderMatcher.setValue(header.value);
          return newHeaderMatcher;
        });
        newMatcher.setHeadersList(newHeaderMatcherList);
      }

      if (queryParametersList !== undefined) {
        let newQueryParamsList = queryParametersList.map(qp => {
          let newQueryParams = new QueryParameterMatcher();
          newQueryParams.setName(qp.name);
          newQueryParams.setRegex(qp.regex);
          newQueryParams.setValue(qp.value);
          return newQueryParams;
        });

        newMatcher.setQueryParametersList(newQueryParamsList);
      }
      newMatcher.setMethodsList(methodsList);
      updatedRoute.setMatchersList([newMatcher]);
    }
    // route action
    if (route.routeAction !== undefined) {
      let updatedRouteAction = new RouteAction();

      if (route.routeAction.single) {
        let updatedDestination = new Destination();
        if (route.routeAction.single.upstream !== undefined) {
          let updatedUpstreamRef = new ResourceRef();
          updatedUpstreamRef.setName(route.routeAction.single.upstream.name);
          updatedUpstreamRef.setNamespace(
            route.routeAction.single.upstream.namespace
          );
          updatedDestination.setUpstream(updatedUpstreamRef);
        }

        if (route.routeAction.single.destinationSpec !== undefined) {
          let updatedDestinationSpec = new DestinationSpec();

          /* ----------------------------------- AWS ---------------------------------- */

          if (route.routeAction.single.destinationSpec.aws !== undefined) {
            let {
              logicalName,
              invocationStyle,
              responseTransformation
            } = route.routeAction.single.destinationSpec.aws;
            let updatedAwsDestinationSpec = new AwsDestinationSpec();
            updatedAwsDestinationSpec.setLogicalName(logicalName);
            updatedAwsDestinationSpec.setInvocationStyle(invocationStyle);
            updatedAwsDestinationSpec.setResponseTransformation(
              responseTransformation
            );
            updatedDestinationSpec.setAws(updatedAwsDestinationSpec);
          }

          /* ---------------------------------- AZURE --------------------------------- */

          if (route.routeAction.single.destinationSpec.azure !== undefined) {
            let {
              functionName
            } = route.routeAction.single.destinationSpec.azure;
            let updatedAzureDestinationSpec = new AzureDestinationSpec();
            updatedAzureDestinationSpec.setFunctionName(functionName);
            updatedDestinationSpec.setAzure(updatedAzureDestinationSpec);
          }

          /* ---------------------------------- REST ---------------------------------- */

          if (route.routeAction.single.destinationSpec.rest !== undefined) {
            let {
              functionName,
              parameters,
              responseTransformation
            } = route.routeAction.single.destinationSpec.rest;
            let updatedRestDestinationSpec = new RestDestinationSpec();
            updatedRestDestinationSpec.setFunctionName(functionName);

            if (parameters !== undefined) {
              let updatedParams = new Parameters();

              if (parameters.path !== undefined) {
                let pathValue = new StringValue();
                pathValue.setValue(parameters.path.value);
                updatedParams.setPath(pathValue);
              }
              // TODO
              if (parameters.headersMap !== undefined) {
              }
              updatedRestDestinationSpec.setParameters(updatedParams);
            }
            // TODO
            // updatedRestDestinationSpec.setResponseTransformation()
            updatedDestinationSpec.setRest(updatedRestDestinationSpec);
          }

          /* ---------------------------------- GRPC ---------------------------------- */

          if (route.routeAction.single.destinationSpec.grpc !== undefined) {
            let {
              pb_function,
              pb_package,
              service,
              parameters
            } = route.routeAction.single.destinationSpec.grpc;
            let updatedGrpcDestinationSpec = new GrpcDestinationSpec();
            updatedGrpcDestinationSpec.setFunction(pb_function);
            updatedGrpcDestinationSpec.setPackage(pb_package);
            updatedGrpcDestinationSpec.setService(service);
            if (parameters !== undefined) {
              let grpcParams = new Parameters();
              if (parameters.path !== undefined) {
                let grpcPath = new StringValue();
                grpcPath.setValue(parameters.path!.value);
                grpcParams.setPath(grpcPath);
              }
              updatedGrpcDestinationSpec.setParameters(grpcParams);
            }
            updatedDestinationSpec.setGrpc(updatedGrpcDestinationSpec);
          }

          // TODO
          // if (route.routeAction.single.kube !== undefined) {}
          // TODO
          // if (route.routeAction.single.consul !== undefined) {}
          // TODO
          // if (route.routeAction.single.subset !== undefined) {}
          updatedDestination.setDestinationSpec(updatedDestinationSpec);
        }

        updatedRouteAction.setSingle(updatedDestination);
      }

      // TODO
      // if (route.routeAction.multi !== undefined) {}

      if (route.routeAction.upstreamGroup !== undefined) {
        let { name, namespace } = route.routeAction.upstreamGroup!;
        let upstreamGroupRef = new ResourceRef();
        upstreamGroupRef.setName(name);
        upstreamGroupRef.setNamespace(namespace);
        updatedRouteAction.setUpstreamGroup(upstreamGroupRef);
      }
      updatedRoute.setRouteAction(updatedRouteAction);
    }

    // updatedRoute.setRedirectAction();
    // updatedRoute.setDirectResponseAction();

    if (route.delegateAction !== undefined) {
      let routeTableRef = new ResourceRef();
      routeTableRef.setName(route.delegateAction.name);
      routeTableRef.setNamespace(route.delegateAction.namespace);

      updatedRoute.setDelegateAction(routeTableRef);
    }

    if (route.routePlugins !== undefined) {
      let updatedRoutePlugins = new RoutePlugins();
      let {
        transformations,
        faults,
        prefixRewrite,
        timeout,
        retries,
        extensions,
        tracing,
        shadowing,
        headerManipulation,
        hostRewrite,
        cors,
        lbHash,
        ratelimitBasic,
        ratelimit,
        waf,
        jwt,
        rbac,
        extauth
      } = route.routePlugins;
      if (prefixRewrite !== undefined) {
        updatedRoutePlugins.setPrefixRewrite(prefixRewrite);
      }
      updatedRoute.setRoutePlugins(updatedRoutePlugins);
    }
  }
  return updatedRoute;
}

function updateRouteTable(
  updateRouteTableRequest: UpdateRouteTableRequest.AsObject
): Promise<UpdateRouteTableResponse.AsObject> {
  return new Promise(async (resolve, reject) => {
    let request = new UpdateRouteTableRequest();
    let { routeTable } = updateRouteTableRequest!;

    if (routeTable !== undefined && routeTable.metadata !== undefined) {
      let newRoutes = routeTable.routesList.map(setInputRouteValues);
      let { name, namespace } = routeTable.metadata;
      let routeTableToUpdateResponse = await getRouteTable({
        ref: { name, namespace }
      });

      if (routeTableToUpdateResponse.getRouteTableDetails() !== undefined) {
        let routeTableDetailsToUpdate = routeTableToUpdateResponse.getRouteTableDetails()!;
        if (routeTableDetailsToUpdate.getRouteTable() !== undefined) {
          let routeTableToUpdate = routeTableDetailsToUpdate.getRouteTable()!;
          routeTableToUpdate.setRoutesList(newRoutes);
          routeTableDetailsToUpdate.setRouteTable(routeTableToUpdate);
          request.setRouteTable(routeTableToUpdate);
        }
      }
    }

    client.updateRouteTable(request, (error, data) => {
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

function updateRouteTableYaml(
  updateRouteTableYamlRequest: UpdateRouteTableYamlRequest.AsObject
): Promise<UpdateRouteTableResponse.AsObject> {
  return new Promise((resolve, reject) => {
    let { ref, editedYaml } = updateRouteTableYamlRequest!.editedYamlData!;
    let request = new UpdateRouteTableYamlRequest();
    let editedYamlData = new EditedResourceYaml();
    let routeTableRef = new ResourceRef();
    routeTableRef.setName(ref!.name);
    routeTableRef.setNamespace(ref!.namespace);
    editedYamlData.setRef(routeTableRef);
    editedYamlData.setEditedYaml(editedYaml);
    request.setEditedYamlData(editedYamlData);

    client.updateRouteTableYaml(request, (error, data) => {
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

function deleteRouteTable(
  deleteRouteTableRequest: DeleteRouteTableRequest.AsObject
): Promise<DeleteRouteTableResponse.AsObject> {
  return new Promise((resolve, reject) => {
    let { ref } = deleteRouteTableRequest!;
    let request = new DeleteRouteTableRequest();
    let routeTableRef = new ResourceRef();
    routeTableRef.setName(ref!.name);
    routeTableRef.setNamespace(ref!.namespace);
    request.setRef(routeTableRef);
    client.deleteRouteTable(request, (error, data) => {
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

export const routeTables = {
  listRouteTables,
  getRouteTable,
  deleteRouteTable,
  createRouteTable,
  updateRouteTable,
  updateRouteTableYaml
};
