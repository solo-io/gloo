import React from 'react';
import { useCreateSecret, useGetUpstreamsList, useCreateRoute } from 'Api';
import { ResourceRef } from 'proto/github.com/solo-io/solo-kit/api/v1/ref_pb';
import { VirtualService } from 'proto/github.com/solo-io/gloo/projects/gateway/api/v1/virtual_service_pb';
import { Formik, Field } from 'formik';
import {
  SoloFormInput,
  TableFormWrapper,
  SoloFormDropdown,
  SoloFormMultiselect,
  SoloFormMetadataBasedDropdown,
  SoloFormMultipartStringCardsList,
  ErrorText
} from 'Components/Common/Form/SoloFormField';
import {
  CreateRouteValuesType,
  createRouteDefaultValues,
  PATH_SPECIFIERS
} from 'Components/Features/Route/CreateRouteModal';
import { ReactComponent as GreenPlus } from 'assets/small-green-plus.svg';
import { NamespacesContext } from 'GlooIApp';
import { Upstream } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/upstream_pb';
import { DestinationSpec } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins_pb';
import { DestinationSpec as AWSDestinationSpec } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins/aws/aws_pb';
import { DestinationSpec as AzureDestinationSpec } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins/azure/azure_pb';
import { DestinationSpec as RestDestinationSpec } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins/rest/rest_pb';
import { DestinationSpec as GrpcDestinationSpec } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins/grpc/grpc_pb';

import { ListUpstreamsRequest } from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/upstream_pb';
import {
  CreateRouteRequest,
  RouteInput
} from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/virtualservice_pb';
import {
  Route,
  Matcher,
  HeaderMatcher,
  RouteAction,
  Destination,
  QueryParameterMatcher
} from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/proxy_pb';

interface Props {
  virtualService?: VirtualService.AsObject;
  reloadVirtualService: (newVirtualService?: VirtualService.AsObject) => any;
}

export const NewRouteRowForm: React.FC<Props> = ({
  virtualService,
  reloadVirtualService
}) => {
  const {
    data: createdVirtualServiceData,
    refetch: makeRequest
  } = useCreateRoute(null);

  const [allUsableUpstreams, setAllUsableUpstreams] = React.useState<
    Upstream.AsObject[]
  >([]);

  const namespaces = React.useContext(NamespacesContext);
  let listUpstreamsRequest = React.useRef(new ListUpstreamsRequest());
  listUpstreamsRequest.current.setNamespacesList(namespaces);

  const {
    data: upstreamsData,
    error: upstreamsError,
    loading: upstreamsLoading
  } = useGetUpstreamsList(listUpstreamsRequest.current);
  React.useEffect(() => {
    setAllUsableUpstreams(
      !!upstreamsData
        ? upstreamsData.upstreamsList.filter(upstream => !!upstream.metadata)
        : []
    );
  }, [upstreamsData]);

  const initialValues = { ...createRouteDefaultValues, virtualService };

  const createRoute = (values: CreateRouteValuesType) => {
    let newRouteReq = new CreateRouteRequest();
    let reqRouteInput = new RouteInput();

    let virtualServiceResourceRef = new ResourceRef();
    virtualServiceResourceRef.setName(values.virtualService!.metadata!.name);
    virtualServiceResourceRef.setNamespace(
      values.virtualService!.metadata!.namespace
    );
    reqRouteInput.setVirtualServiceRef(virtualServiceResourceRef);

    //reqRouteInput.setIndex(vs.virtualHost!.routesList.length);

    /***
     *  ROUTE CREATION BEGINS
     * */
    let newRoute = new Route();
    let routeMatcher = new Matcher();
    switch (values.matchType) {
      case 'PREFIX':
        routeMatcher.setPrefix(values.path);
        break;
      case 'EXACT':
        routeMatcher.setExact(values.path);
        break;
      case 'REGEX':
        routeMatcher.setRegex(values.path);
        break;
    }

    let matcherHeaders: HeaderMatcher[] = values.headers.map(head => {
      const newMatcherHeader = new HeaderMatcher();
      newMatcherHeader.setName(head.name);
      newMatcherHeader.setValue(head.value);
      newMatcherHeader.setRegex(head.isRegex);

      return newMatcherHeader;
    });
    routeMatcher.setHeadersList(matcherHeaders);
    let matcherQueryParams: QueryParameterMatcher[] = values.queryParameters.map(
      queryParam => {
        const newMatcherQueryParam = new QueryParameterMatcher();
        newMatcherQueryParam.setName(queryParam.name);
        newMatcherQueryParam.setValue(queryParam.value);
        newMatcherQueryParam.setRegex(queryParam.isRegex);

        return newMatcherQueryParam;
      }
    );
    routeMatcher.setQueryParametersList(matcherQueryParams);
    routeMatcher.setMethodsList(
      //@ts-ignore
      Object.keys(values.methods).filter(key => values.methods[key])
    );
    newRoute.setMatcher(routeMatcher);

    /* Route->Destination Section */
    let newRouteAction = new RouteAction();
    let newDestination = new Destination();
    const upstreamSpec = values.upstream!.upstreamSpec!;
    let newDestinationResourceRef = new ResourceRef();
    newDestinationResourceRef.setName(values.upstream!.metadata!.name);
    newDestinationResourceRef.setNamespace(
      values.upstream!.metadata!.namespace
    );
    let newDestinationSpec = new DestinationSpec();

    if (!!upstreamSpec.aws) {
      newDestination.setUpstream(newDestinationResourceRef);
      let newAWSDestinationSpec = new AWSDestinationSpec();
      // TODO :: I have no idea what to set the values to
      //newAWSDestinationSpec.setInvocationStyle(0);
      newDestinationSpec.setAws(newAWSDestinationSpec);
    } else if (!!upstreamSpec.azure) {
      newDestination.setUpstream(newDestinationResourceRef);
      let newAzureDestinationSpec = new AzureDestinationSpec();
      // TODO :: I have no idea what to set the values to
      newDestinationSpec.setAzure(newAzureDestinationSpec);
    } /*else if (!!upstreamSpec.kube) {
      let newKubeServiceDestination = new KubernetesServiceDestination();
      newKubeServiceDestination.setRef(newDestinationResourceRef);
      // TODO :: I have no idea what to set the values to
      newDestination.setKube(newKubeServiceDestination);
      let newKubeDestinationSpec;
      // TODO:: How do we tell if it is rest or GRPC?
      //if() -> set DestinationSpec to grpc...
      newDestination.setDestinationSpec(newKubeDestinationSpec);
    } else if (!!upstreamSpec.consul) {
      let newConsulServiceDestination = new ConsulServiceDestination();
      // TODO :: I have no idea what to set the values to
      newDestination.setConsul(newConsulServiceDestination);
      let newConsulDestinationSpec;
      // TODO:: I have no idea what goes in this case
      newDestination.setDestinationSpec(newConsulDestinationSpec);
    }*/
    newDestination.setDestinationSpec(newDestinationSpec);
    newRouteAction.setSingle(newDestination);
    newRoute.setRouteAction(newRouteAction);

    // It looks like we don't see any of the other actions if
    // Route Action is taken??  But if they supplied
    // a path, shouldn't we do the redirect action?
    // Not clear on what the other actions would be based on?

    /*
    let newRedirectAction = new RedirectAction();
    //TODO:: Do we need to set anything else for this???
    if(values.matchType === PathSpecifierCase.PREFIX) {
      // TODO:: Is this correct??
      newRedirectAction.setPrefixRewrite("PREFIX");
    } else {
      newRedirectAction.setPathRedirect(values.path);
    }
    newRoute.setRedirectAction(newRedirectAction);*/

    reqRouteInput.setRoute(newRoute);
    /***
     *  ROUTE CREATION ENDS
     * */

    newRouteReq.setInput(reqRouteInput);
    makeRequest(newRouteReq);

    setTimeout(reloadVirtualService, 300);
  };

  return (
    <React.Fragment>
      <Formik<CreateRouteValuesType>
        initialValues={initialValues}
        onSubmit={createRoute}>
        {({ handleSubmit, values }) => (
          <React.Fragment>
            <TableFormWrapper>
              <Field
                name='path'
                placeholder='Path...'
                component={SoloFormInput}
              />
              <Field
                name='matchType'
                defaultValue={'PREFIX'}
                options={PATH_SPECIFIERS}
                component={SoloFormDropdown}
              />
              <Field
                name='methods'
                placeholder='Methods...'
                options={Object.keys(createRouteDefaultValues.methods).map(
                  key => {
                    return {
                      key: key,
                      value: key
                    };
                  }
                )}
                component={SoloFormMultiselect}
              />
              <Field
                name='upstream'
                placeholder='Upstream...'
                options={allUsableUpstreams}
                component={SoloFormMetadataBasedDropdown}
              />
              <Field
                name='headers'
                createNewNamePromptText={'Name...'}
                createNewValuePromptText={'Value...'}
                component={SoloFormMultipartStringCardsList}
              />
              <Field
                name='queryParameters'
                createNewNamePromptText={'Name...'}
                createNewValuePromptText={'Value...'}
                component={SoloFormMultipartStringCardsList}
              />
            </TableFormWrapper>
            <td>
              <GreenPlus
                style={{ cursor: 'pointer' }}
                onClick={() => handleSubmit()}
              />
              <ErrorText errorExists={false} />
            </td>
          </React.Fragment>
        )}
      </Formik>
    </React.Fragment>
  );
};
