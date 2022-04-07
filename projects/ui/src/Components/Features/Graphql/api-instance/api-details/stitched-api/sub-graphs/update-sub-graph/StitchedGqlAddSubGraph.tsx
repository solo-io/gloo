import { graphqlConfigApi } from 'API/graphql';
import {
  useGetConsoleOptions,
  useGetGraphqlApiDetails,
  useListGraphqlApis,
  usePageApiRef,
  usePageGlooInstance,
} from 'API/hooks';
import SoloAddButton from 'Components/Common/SoloAddButton';
import { OptionType, SoloDropdown } from 'Components/Common/SoloDropdown';
import { SoloModal } from 'Components/Common/SoloModal';
import { StitchedSchema } from 'proto/github.com/solo-io/solo-apis/api/gloo/graphql.gloo/v1beta1/graphql_pb';
import React, { useMemo, useState } from 'react';
import { SoloButtonStyledComponent } from 'Styles/StyledComponents/button';
import StitchedGqlTypeMergeMapConfig from '../type-merge-map/StitchedGqlTypeMergeMapConfig';

interface gqlOptionType extends OptionType {
  apiIndex: number;
}
const nameNamespaceKey = (api?: { name?: string; namespace?: string }) =>
  `${api?.name ?? ''} ${api?.namespace ?? ''}`;

const StitchedGqlAddSubGraph: React.FC<{ onAfterAdd(): void }> = ({
  onAfterAdd,
}) => {
  const apiRef = usePageApiRef();
  const { data: graphqlApi } = useGetGraphqlApiDetails(apiRef);
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [glooInstance] = usePageGlooInstance();
  const { data: graphqlApis } = useListGraphqlApis(glooInstance?.metadata);
  const { readonly } = useGetConsoleOptions();

  // --- TYPE MERGE MAP --- //
  const [typeMergeMap, setTypeMergeMap] = useState<
    [string, StitchedSchema.SubschemaConfig.TypeMergeConfig.AsObject][]
  >([]);
  const addSubGraph = async () => {
    if (
      !selectedOption ||
      !apiToAdd ||
      (isShowingTypeMergeMap && !isTypeMergeMapValid) ||
      !graphqlApi
    ) {
      console.error('Unable to update type merge map.');
      return;
    }
    //
    // Get the selected api to add, and create a new sub graph reference to it.
    const newSubGraph = {
      name: apiToAdd.metadata?.name ?? '',
      namespace: apiToAdd.metadata?.namespace ?? '',
      typeMergeMap,
    };
    //
    // Update the api with a new spec that includes the sub graphs.
    const existingSubGraphs =
      graphqlApi.spec?.stitchedSchema?.subschemasList ?? [];
    await graphqlConfigApi.updateGraphqlApi({
      graphqlApiRef: apiRef,
      spec: {
        stitchedSchema: {
          subschemasList: [...existingSubGraphs, newSubGraph],
        },
        allowedQueryHashesList: [],
      },
    });
    setIsModalVisible(false);
    setSelectedOption(undefined);
    onAfterAdd();
  };
  // TODO: The StitchedGqlTypeMergeMapConfig UI is parsing the type merge map and keeping it in state - this should be cleaner.
  // Need useMemo here to prevent max update depth error.
  const initialTypeMergeMapForConfig = useMemo(() => [], []);

  // -- SUB GRAPH SELECTION -- //
  const [selectedOption, setSelectedOption] = useState<gqlOptionType>();
  const options = useMemo(() => {
    if (!graphqlApis || !graphqlApi) return [];
    const thisApi = nameNamespaceKey(graphqlApi.metadata);
    return graphqlApis
      .map((a, idx) => ({
        key: nameNamespaceKey(a.metadata),
        apiIndex: idx,
        value: nameNamespaceKey(a.metadata),
        displayValue: a.metadata?.name ?? '',
      }))
      .filter(
        option =>
          // Filters this api out of the dropdown.
          option?.value !== thisApi &&
          // Filters any already-added sub graphs out of the dropdown.
          !graphqlApi?.spec?.stitchedSchema?.subschemasList.find(
            s => option?.value === nameNamespaceKey(s)
          )
      );
  }, [graphqlApis, graphqlApi]);
  let apiToAdd = useMemo(() => {
    if (
      !graphqlApis ||
      !selectedOption ||
      selectedOption.apiIndex < 0 ||
      selectedOption.apiIndex >= graphqlApis.length
    )
      return undefined;
    return graphqlApis[selectedOption.apiIndex];
  }, [graphqlApis, selectedOption]);
  const subGraphConfig = useMemo(
    () =>
      ({
        name: apiToAdd?.metadata?.name ?? '',
        namespace: apiToAdd?.metadata?.namespace ?? '',
        typeMergeMap: [],
      } as StitchedSchema.SubschemaConfig.AsObject),
    [apiToAdd]
  );
  // Show the type merge map if:
  //   - There is metadata for it.
  //   - There is > 2 apis (one api is the one we're adding,
  //     and the other one is this stitched one)
  // TODO: After we check for conflicts, we can populate the initial type merge map and check if that exists here instead.
  const isShowingTypeMergeMap =
    apiToAdd?.metadata && graphqlApis && graphqlApis.length > 2;

  // --- VALIDATION --- //
  const [isTypeMergeMapValid, setIsTypeMergeMapValid] = useState(false);
  const canSubmit =
    (!isShowingTypeMergeMap || isTypeMergeMapValid) && !!selectedOption;

  if (readonly) return null;
  return (
    <div>
      <SoloAddButton onClick={() => setIsModalVisible(true)}>
        Add Sub Graph
      </SoloAddButton>

      {!!graphqlApis && (
        <SoloModal
          visible={isModalVisible}
          onClose={() => setIsModalVisible(false)}
          title='Add Sub Graph'
          width={600}>
          <div className='p-5 pb-10 pt-3'>
            <SoloDropdown
              title='Select a GraphQL API'
              options={options}
              value={selectedOption?.value ?? ''}
              onChange={value =>
                setSelectedOption(options.find(o => o.value === value))
              }
              searchable={true}
            />

            {isShowingTypeMergeMap && (
              <StitchedGqlTypeMergeMapConfig
                onIsValidChange={isValid => setIsTypeMergeMapValid(isValid)}
                initialTypeMergeMap={initialTypeMergeMapForConfig}
                onTypeMergeMapChange={m => setTypeMergeMap(m)}
                apiRef={apiRef}
                subGraphConfig={subGraphConfig}
              />
            )}

            <div className='text-right mt-10'>
              <SoloButtonStyledComponent
                disabled={!canSubmit}
                onClick={addSubGraph}>
                Add Sub Graph
              </SoloButtonStyledComponent>
            </div>
          </div>
        </SoloModal>
      )}
    </div>
  );
};

export default StitchedGqlAddSubGraph;
