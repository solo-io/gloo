import { useListGraphqlApis, usePageGlooInstance } from 'API/hooks';
import { DataError } from 'Components/Common/DataError';
import { Loading } from 'Components/Common/Loading';
import { SoloCheckbox } from 'Components/Common/SoloCheckbox';
import { SoloInput } from 'Components/Common/SoloInput';
import { GraphqlApi } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/graphql_pb';
import React, { useEffect, useState } from 'react';
import * as styles from '../GraphqlLanding.style';
import { GraphqlLandingTable } from './GraphqlLandingTable';
import { apiFilterGroups, ITypeFilter } from './GraphqlLandingTableFilters';

// *
// * This container component filters the list of graphql apis
// * and passes them to the GraphqlLandingTable.
// *
export const GraphqlLandingTableContainer = () => {
  const { glooInstance, glooInstances } = usePageGlooInstance();

  // --- FILTERS STATE --- //
  const [searchText, setSearchText] = useState('');
  const [filters, setFilters] = useState<{
    [key: string]: ITypeFilter<GraphqlApi.AsObject>;
  }>({});

  // --- FILTERED GRAPHQL APIS --- //
  const { data: graphqlApis, error: graphqlApisError } = useListGraphqlApis();
  const [filteredGraphqlApis, setFilteredGraphqlApis] = useState<
    GraphqlApi.AsObject[]
  >([]);
  useEffect(() => {
    if (!graphqlApis) {
      setFilteredGraphqlApis([]);
      return;
    }
    // Apply filters.
    const filterList = Object.keys(filters).map(k => filters[k]);
    const intersectionFilters = filterList.filter(
      f => f.type === 'intersection'
    );
    const unionFilters = filterList.filter(f => f.type === 'union');
    let newFilteredGqlApis = graphqlApis.filter(api => {
      if (
        // If there are intersection filters and any of them fail, remove this item.
        (unionFilters.length !== 0 &&
          intersectionFilters.find(f => !f.filterFn(api)) !== undefined) ||
        !api.metadata?.name.includes(searchText)
      )
        return false;
      // If there are no union filters or any of them succeed, keep this item.
      return (
        unionFilters.length === 0 ||
        unionFilters.find(f => f.filterFn(api)) !== undefined
      );
    });
    setFilteredGraphqlApis(newFilteredGqlApis);
  }, [graphqlApis, searchText, filters]);

  // --- ERROR, LOADING --- //
  if (!!graphqlApisError) return <DataError error={graphqlApisError} />;
  if (!graphqlApis) return <Loading message={'Retrieving GraphQL APIs...'} />;

  return (
    <>
      {/* --- FILTER CONTROLS --- */}
      <div>
        <SoloInput
          value={searchText}
          onChange={e => setSearchText(e.target.value)}
          placeholder={'Filter by name...'}
        />
        {apiFilterGroups.map(fg => (
          <div key={fg.key}>
            {fg.title && (
              <styles.HorizontalDivider>
                <div>{fg.title}</div>
              </styles.HorizontalDivider>
            )}
            <styles.CheckboxWrapper>
              {fg.filters.map((filter, ind) => {
                const nestedKey = `${fg.key}.${filter.key}`;
                return (
                  <SoloCheckbox
                    key={nestedKey}
                    title={filter.displayValue}
                    withWrapper={true}
                    onChange={evt => {
                      let newApiTypeFilters = { ...filters };
                      if (evt.target.checked)
                        newApiTypeFilters[nestedKey] = filter;
                      else delete newApiTypeFilters[nestedKey];
                      setFilters(newApiTypeFilters);
                    }}
                    checked={!!filters[nestedKey]}
                  />
                );
              })}
            </styles.CheckboxWrapper>
          </div>
        ))}
      </div>
      {/* --- TABLE --- */}
      {glooInstances &&
        (glooInstance ? (
          <GraphqlLandingTable
            title='GraphQL'
            glooInstance={glooInstance}
            graphqlApis={filteredGraphqlApis.filter(
              api =>
                api.glooInstance?.name === glooInstance.metadata?.name &&
                api.glooInstance?.namespace === glooInstance.metadata?.namespace
            )}
          />
        ) : (
          <div>
            {glooInstances.map(gi => (
              <GraphqlLandingTable
                key={`${gi.metadata?.name}:${gi.metadata?.namespace}`}
                title={`GraphQL - ${gi.metadata?.name} - ${gi.metadata?.namespace}`}
                glooInstance={gi}
                graphqlApis={filteredGraphqlApis.filter(
                  api =>
                    api.glooInstance?.name === gi.metadata?.name &&
                    api.glooInstance?.namespace === gi.metadata?.namespace
                )}
              />
            ))}
          </div>
        ))}
    </>
  );
};
