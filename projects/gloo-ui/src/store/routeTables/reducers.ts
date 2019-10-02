import { RouteTableDetails } from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/routetable_pb';
import { RouteTableActionTypes, RouteTableAction } from './types';
import { SoloWarning } from 'Components/Common/SoloWarningContent';

export interface RouteTableState {
  routeTablesList: RouteTableDetails.AsObject[];
  yamlParseError: boolean;
}

const initialState: RouteTableState = {
  routeTablesList: [],
  yamlParseError: false
};

export function routeTablesReducer(
  state = initialState,
  action: RouteTableActionTypes
): RouteTableState {
  switch (action.type) {
    case RouteTableAction.LIST_ROUTE_TABLES:
      return {
        ...state,
        routeTablesList: action.payload
      };

    case RouteTableAction.UPDATE_ROUTE_TABLE_YAML_ERROR:
      SoloWarning(
        'There was an error updating the route table.',
        action.payload
      );

      return {
        ...state,
        yamlParseError: true
      };

    default:
      return { ...state, yamlParseError: false };
  }
}
