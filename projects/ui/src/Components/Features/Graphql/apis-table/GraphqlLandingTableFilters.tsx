import { GraphqlApi } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/graphql_pb';
import { isExecutableAPI, isStitchedAPI } from 'utils/graphql-helpers';

// -- TYPES -- //

export type FilterType = 'intersection' | 'union';
export interface ITypeFilter<T> {
  key: string;
  type: FilterType;
  displayValue: string;
  filterFn(data: T): boolean;
}
export interface IFilterGroup<T> {
  key: string;
  title: string;
  filters: ITypeFilter<T>[];
}

// -- FILTERS -- //
// Edit this to add new filters.
export const apiFilterGroups = [
  {
    key: 'api-type',
    title: 'API Type',
    filters: [
      {
        type: 'union',
        key: 'executable',
        displayValue: 'Executable',
        filterFn: d => isExecutableAPI(d),
      },
      {
        type: 'union',
        key: 'stitched',
        displayValue: 'Stitched',
        filterFn: d => isStitchedAPI(d),
      },
    ],
  },
] as IFilterGroup<GraphqlApi.AsObject>[];
