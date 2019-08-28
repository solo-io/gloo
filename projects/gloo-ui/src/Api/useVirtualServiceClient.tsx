import * as React from 'react';
import { requestReducer, Reducer, RequestAction } from './request-reducer';

import {
  GetVirtualServiceRequest,
  GetVirtualServiceResponse,
  ListVirtualServicesRequest,
  ListVirtualServicesResponse,
  CreateVirtualServiceRequest,
  CreateVirtualServiceResponse,
  UpdateVirtualServiceRequest,
  UpdateVirtualServiceResponse,
  DeleteVirtualServiceRequest,
  DeleteVirtualServiceResponse,
  CreateRouteRequest,
  CreateRouteResponse,
  UpdateRouteRequest,
  UpdateRouteResponse,
  DeleteRouteRequest,
  DeleteRouteResponse,
  SwapRoutesRequest,
  SwapRoutesResponse,
  ShiftRoutesRequest,
  ShiftRoutesResponse
} from '../proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/virtualservice_pb.d';
import { ServiceError } from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/virtualservice_pb_service';
import { client } from './v2/VirtualServiceClient';

// export class VirtualServiceApi {
//   static readonly serviceName: string;
//   static readonly GetVirtualService: VirtualServiceApiGetVirtualService;
//   static readonly ListVirtualServices: VirtualServiceApiListVirtualServices;
//   static readonly StreamVirtualServiceList: VirtualServiceApiStreamVirtualServiceList;
//   static readonly CreateVirtualService: VirtualServiceApiCreateVirtualService;
//   static readonly UpdateVirtualService: VirtualServiceApiUpdateVirtualService;
//   static readonly DeleteVirtualService: VirtualServiceApiDeleteVirtualService;
//   static readonly CreateRoute: VirtualServiceApiCreateRoute;
//   static readonly UpdateRoute: VirtualServiceApiUpdateRoute;
//   static readonly DeleteRoute: VirtualServiceApiDeleteRoute;
//   static readonly SwapRoutes: VirtualServiceApiSwapRoutes;
//   static readonly ShiftRoutes: VirtualServiceApiShiftRoutes;
// }

/* -------------------------------------------------------------------------- */
/*                             GET VIRTUAL SERVICE                            */
/* -------------------------------------------------------------------------- */

export const useGetVirtualService = (
  request: GetVirtualServiceRequest | null,
  initialData: GetVirtualServiceResponse.AsObject | null = null
) => {
  const [state, dispatch] = React.useReducer<
    Reducer<GetVirtualServiceResponse.AsObject | null>
  >(requestReducer, {
    isLoading: true,
    data: initialData
  });

  const mounted = React.useRef(true);

  const makeRequest = (request: GetVirtualServiceRequest | null) => {
    if (!request) {
      return;
    }
    dispatch({ type: RequestAction.START, payload: null });

    client.getVirtualService(
      request,
      (
        error: ServiceError | null,
        responseMessage: GetVirtualServiceResponse | null
      ) => {
        if (error) {
          console.error('Error:', error.message);
          console.error('Code:', error.code);
          console.error('Metadata:', error.metadata);
          if (!mounted.current) return;
          dispatch({ type: RequestAction.ERROR, payload: null, error });
        } else {
          const response = responseMessage;
          if (!mounted.current) return;
          dispatch({
            type: RequestAction.SUCCESS,
            payload: response!.toObject()
          });
        }
      }
    );
  };

  React.useEffect(() => {
    makeRequest(request || null);
    return () => {
      mounted.current = false;
    };
  }, []);

  return {
    data: state.data!,
    error: state.error,
    loading: state.isLoading,
    refetch: makeRequest
  };
};
/* -------------------------------------------------------------------------- */
/*                            LIST VIRTUAL SERVICES                           */
/* -------------------------------------------------------------------------- */

export const useListVirtualServices = (
  request: ListVirtualServicesRequest | null,
  initialData: ListVirtualServicesResponse.AsObject | null = null
) => {
  const [state, dispatch] = React.useReducer<
    Reducer<ListVirtualServicesResponse.AsObject | null>
  >(requestReducer, {
    isLoading: true,
    data: initialData
  });

  const mounted = React.useRef(true);

  const makeRequest = (request: ListVirtualServicesRequest | null) => {
    if (!request) {
      return;
    }
    dispatch({ type: RequestAction.START, payload: null });

    client.listVirtualServices(
      request,
      (
        error: ServiceError | null,
        responseMessage: ListVirtualServicesResponse | null
      ) => {
        if (error) {
          console.error('Error:', error.message);
          console.error('Code:', error.code);
          console.error('Metadata:', error.metadata);
          if (!mounted.current) return;
          dispatch({ type: RequestAction.ERROR, payload: null, error });
        } else {
          const response = responseMessage;
          if (!mounted.current) return;
          dispatch({
            type: RequestAction.SUCCESS,
            payload: response!.toObject()
          });
        }
      }
    );
  };

  React.useEffect(() => {
    makeRequest(request || null);
    return () => {
      mounted.current = false;
    };
  }, []);

  return {
    data: state.data!,
    error: state.error,
    loading: state.isLoading,
    refetch: makeRequest
  };
};
/* -------------------------------------------------------------------------- */
/*                           CREATE VIRTUAL SERVICE                           */
/* -------------------------------------------------------------------------- */

export const useCreateVirtualService = (
  request: CreateVirtualServiceRequest | null,
  initialData: CreateVirtualServiceResponse.AsObject | null = null
) => {
  const [state, dispatch] = React.useReducer<
    Reducer<CreateVirtualServiceResponse.AsObject | null>
  >(requestReducer, {
    isLoading: true,
    data: initialData
  });

  const mounted = React.useRef(true);

  const makeRequest = (request: CreateVirtualServiceRequest | null) => {
    if (!request) {
      return;
    }
    dispatch({ type: RequestAction.START, payload: null });

    client.createVirtualService(
      request,
      (
        error: ServiceError | null,
        responseMessage: CreateVirtualServiceResponse | null
      ) => {
        if (error) {
          console.error('Error:', error.message);
          console.error('Code:', error.code);
          console.error('Metadata:', error.metadata);
          if (!mounted.current) return;
          dispatch({ type: RequestAction.ERROR, payload: null, error });
        } else {
          const response = responseMessage;
          if (!mounted.current) return;
          dispatch({
            type: RequestAction.SUCCESS,
            payload: response!.toObject()
          });
        }
      }
    );
  };

  React.useEffect(() => {
    makeRequest(request || null);
    return () => {
      mounted.current = false;
    };
  }, []);

  return {
    data: state.data!,
    error: state.error,
    loading: state.isLoading,
    refetch: makeRequest
  };
};
/* -------------------------------------------------------------------------- */
/*                           UPDATE VIRTUAL SERVICE                           */
/* -------------------------------------------------------------------------- */

export const useUpdateVirtualService = (
  request: UpdateVirtualServiceRequest | null,
  initialData: UpdateVirtualServiceResponse.AsObject | null = null
) => {
  const [state, dispatch] = React.useReducer<
    Reducer<UpdateVirtualServiceResponse.AsObject | null>
  >(requestReducer, {
    isLoading: true,
    data: initialData
  });

  const mounted = React.useRef(true);

  const makeRequest = (request: UpdateVirtualServiceRequest | null) => {
    if (!request) {
      dispatch({ type: RequestAction.INITIALREFETCH, payload: null });
      return;
    }
    dispatch({ type: RequestAction.START, payload: null });

    client.updateVirtualService(
      request,
      (
        error: ServiceError | null,
        responseMessage: UpdateVirtualServiceResponse | null
      ) => {
        if (error) {
          console.error('Error:', error.message);
          console.error('Code:', error.code);
          console.error('Metadata:', error.metadata);
          if (!mounted.current) return;
          dispatch({ type: RequestAction.ERROR, payload: null, error });
        } else {
          const response = responseMessage;
          if (!mounted.current) return;
          dispatch({
            type: RequestAction.SUCCESS,
            payload: response!.toObject()
          });
        }
      }
    );
  };

  React.useEffect(() => {
    makeRequest(request || null);
    return () => {
      mounted.current = false;
    };
  }, []);

  return {
    data: state.data!,
    error: state.error,
    loading: state.isLoading,
    refetch: makeRequest
  };
};
/* -------------------------------------------------------------------------- */
/*                           DELETE VIRTUAL SERVICE                           */
/* -------------------------------------------------------------------------- */

export const useDeleteVirtualService = (
  request: DeleteVirtualServiceRequest | null,
  initialData: DeleteVirtualServiceResponse.AsObject | null = null
) => {
  const [state, dispatch] = React.useReducer<
    Reducer<DeleteVirtualServiceResponse.AsObject | null>
  >(requestReducer, {
    isLoading: true,
    data: initialData
  });

  const mounted = React.useRef(true);

  const makeRequest = (request: DeleteVirtualServiceRequest | null) => {
    if (!request) {
      return;
    }
    dispatch({ type: RequestAction.START, payload: null });

    client.deleteVirtualService(
      request,
      (
        error: ServiceError | null,
        responseMessage: DeleteVirtualServiceResponse | null
      ) => {
        if (error) {
          console.error('Error:', error.message);
          console.error('Code:', error.code);
          console.error('Metadata:', error.metadata);
          if (!mounted.current) return;
          dispatch({ type: RequestAction.ERROR, payload: null, error });
        } else {
          const response = responseMessage;
          if (!mounted.current) return;
          dispatch({
            type: RequestAction.SUCCESS,
            payload: response!.toObject()
          });
        }
      }
    );
  };

  React.useEffect(() => {
    makeRequest(request || null);
    return () => {
      mounted.current = false;
    };
  }, []);

  return {
    data: state.data!,
    error: state.error,
    loading: state.isLoading,
    refetch: makeRequest
  };
};
/* --------------------------------- ROUTES --------------------------------- */

/* -------------------------------------------------------------------------- */
/*                                CREATE ROUTE                                */
/* -------------------------------------------------------------------------- */

export const useCreateRoute = (
  request: CreateRouteRequest | null,
  initialData: CreateRouteResponse.AsObject | null = null
) => {
  const [state, dispatch] = React.useReducer<
    Reducer<CreateRouteResponse.AsObject | null>
  >(requestReducer, {
    isLoading: true,
    data: initialData
  });

  const mounted = React.useRef(true);

  const makeRequest = (request: CreateRouteRequest | null) => {
    if (!request) {
      return;
    }
    dispatch({ type: RequestAction.START, payload: null });

    client.createRoute(
      request,
      (
        error: ServiceError | null,
        responseMessage: CreateRouteResponse | null
      ) => {
        if (error) {
          console.error('Error:', error.message);
          console.error('Code:', error.code);
          console.error('Metadata:', error.metadata);
          if (!mounted.current) return;
          dispatch({ type: RequestAction.ERROR, payload: null, error });
        } else {
          const response = responseMessage;
          if (!mounted.current) return;
          dispatch({
            type: RequestAction.SUCCESS,
            payload: response!.toObject()
          });
        }
      }
    );
  };

  React.useEffect(() => {
    makeRequest(request || null);
    return () => {
      mounted.current = false;
    };
  }, []);

  return {
    data: state.data!,
    error: state.error,
    loading: state.isLoading,
    refetch: makeRequest
  };
};

/* -------------------------------------------------------------------------- */
/*                                UPDATE ROUTE                                */
/* -------------------------------------------------------------------------- */

export const useUpdateRoute = (
  request: UpdateRouteRequest | null,
  initialData: UpdateRouteResponse.AsObject | null = null
) => {
  const [state, dispatch] = React.useReducer<
    Reducer<UpdateRouteResponse.AsObject | null>
  >(requestReducer, {
    isLoading: true,
    data: initialData
  });

  const mounted = React.useRef(true);

  const makeRequest = (request: UpdateRouteRequest | null) => {
    if (!request) {
      return;
    }
    dispatch({ type: RequestAction.START, payload: null });

    client.updateRoute(
      request,
      (
        error: ServiceError | null,
        responseMessage: UpdateRouteResponse | null
      ) => {
        if (error) {
          console.error('Error:', error.message);
          console.error('Code:', error.code);
          console.error('Metadata:', error.metadata);
          if (!mounted.current) return;
          dispatch({ type: RequestAction.ERROR, payload: null, error });
        } else {
          const response = responseMessage;
          if (!mounted.current) return;
          dispatch({
            type: RequestAction.SUCCESS,
            payload: response!.toObject()
          });
        }
      }
    );
  };

  React.useEffect(() => {
    makeRequest(request || null);
    return () => {
      mounted.current = false;
    };
  }, []);

  return {
    data: state.data!,
    error: state.error,
    loading: state.isLoading,
    refetch: makeRequest
  };
};
/* -------------------------------------------------------------------------- */
/*                                DELETE ROUTE                                */
/* -------------------------------------------------------------------------- */

export const useDeleteRoute = (
  request: DeleteRouteRequest | null,
  initialData: DeleteRouteResponse.AsObject | null = null
) => {
  const [state, dispatch] = React.useReducer<
    Reducer<DeleteRouteResponse.AsObject | null>
  >(requestReducer, {
    isLoading: true,
    data: initialData
  });

  const mounted = React.useRef(true);

  const makeRequest = (request: DeleteRouteRequest | null) => {
    if (!request) {
      return;
    }
    dispatch({ type: RequestAction.START, payload: null });

    client.deleteRoute(
      request,
      (
        error: ServiceError | null,
        responseMessage: DeleteRouteResponse | null
      ) => {
        if (error) {
          console.error('Error:', error.message);
          console.error('Code:', error.code);
          console.error('Metadata:', error.metadata);
          if (!mounted.current) return;
          dispatch({ type: RequestAction.ERROR, payload: null, error });
        } else {
          const response = responseMessage;
          if (!mounted.current) return;
          dispatch({
            type: RequestAction.SUCCESS,
            payload: response!.toObject()
          });
        }
      }
    );
  };

  React.useEffect(() => {
    makeRequest(request || null);
    return () => {
      mounted.current = false;
    };
  }, []);

  return {
    data: state.data!,
    error: state.error,
    loading: state.isLoading,
    refetch: makeRequest
  };
};
/* -------------------------------------------------------------------------- */
/*                                 SWAP ROUTES                                */
/* -------------------------------------------------------------------------- */

export const useSwapRoutes = (
  request: SwapRoutesRequest | null,
  initialData: SwapRoutesResponse.AsObject | null = null
) => {
  const [state, dispatch] = React.useReducer<
    Reducer<SwapRoutesResponse.AsObject | null>
  >(requestReducer, {
    isLoading: true,
    data: initialData
  });

  const mounted = React.useRef(true);

  const makeRequest = (request: SwapRoutesRequest | null) => {
    if (!request) {
      return;
    }
    dispatch({ type: RequestAction.START, payload: null });

    client.swapRoutes(
      request,
      (
        error: ServiceError | null,
        responseMessage: SwapRoutesResponse | null
      ) => {
        if (error) {
          console.error('Error:', error.message);
          console.error('Code:', error.code);
          console.error('Metadata:', error.metadata);
          if (!mounted.current) return;
          dispatch({ type: RequestAction.ERROR, payload: null, error });
        } else {
          const response = responseMessage;
          if (!mounted.current) return;
          dispatch({
            type: RequestAction.SUCCESS,
            payload: response!.toObject()
          });
        }
      }
    );
  };

  React.useEffect(() => {
    makeRequest(request || null);
    return () => {
      mounted.current = false;
    };
  }, []);

  return {
    data: state.data!,
    error: state.error,
    loading: state.isLoading,
    refetch: makeRequest
  };
};
/* -------------------------------------------------------------------------- */
/*                                SHIFT ROUTES                                */
/* -------------------------------------------------------------------------- */

export const useShiftRoutes = (
  request: ShiftRoutesRequest | null,
  initialData: ShiftRoutesResponse.AsObject | null = null
) => {
  const [state, dispatch] = React.useReducer<
    Reducer<ShiftRoutesResponse.AsObject | null>
  >(requestReducer, {
    isLoading: true,
    data: initialData
  });

  const mounted = React.useRef(true);

  const makeRequest = (request: ShiftRoutesRequest | null) => {
    if (!request) {
      return;
    }
    dispatch({ type: RequestAction.START, payload: null });

    client.shiftRoutes(
      request,
      (
        error: ServiceError | null,
        responseMessage: ShiftRoutesResponse | null
      ) => {
        if (error) {
          console.error('Error:', error.message);
          console.error('Code:', error.code);
          console.error('Metadata:', error.metadata);
          if (!mounted.current) return;
          dispatch({ type: RequestAction.ERROR, payload: null, error });
        } else {
          const response = responseMessage;
          if (!mounted.current) return;
          dispatch({
            type: RequestAction.SUCCESS,
            payload: response!.toObject()
          });
        }
      }
    );
  };

  React.useEffect(() => {
    makeRequest(request || null);
    return () => {
      mounted.current = false;
    };
  }, []);

  return {
    data: state.data!,
    error: state.error,
    loading: state.isLoading,
    refetch: makeRequest
  };
};
