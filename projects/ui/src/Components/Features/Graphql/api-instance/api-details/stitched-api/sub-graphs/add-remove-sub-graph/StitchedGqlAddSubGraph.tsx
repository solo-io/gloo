import { graphqlConfigApi } from 'API/graphql';
import {
  useGetConsoleOptions,
  useGetGraphqlApiDetails,
  useGetGraphqlApiYaml,
  useListGraphqlApis,
  usePageApiRef,
  usePageGlooInstance,
} from 'API/hooks';
import SoloAddButton from 'Components/Common/SoloAddButton';
import { OptionType, SoloDropdown } from 'Components/Common/SoloDropdown';
import { SoloModal } from 'Components/Common/SoloModal';
import React, { useMemo, useState } from 'react';
import { SoloButtonStyledComponent } from 'Styles/StyledComponents/button';

interface gqlOptionType extends OptionType {
  apiIndex: number;
}
const nameNamespaceKey = (api?: { name?: string; namespace?: string }) =>
  `${api?.name ?? ''} ${api?.namespace ?? ''}`;

const StitchedGqlAddSubGraph = () => {
  const apiRef = usePageApiRef();
  const { mutate: mutateYaml } = useGetGraphqlApiYaml(apiRef);
  const { data: graphqlApi, mutate: mutateDetails } =
    useGetGraphqlApiDetails(apiRef);
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [glooInstance] = usePageGlooInstance();
  const { data: graphqlApis } = useListGraphqlApis(glooInstance?.metadata);
  const { readonly } = useGetConsoleOptions();

  const addSubGraph = async () => {
    if (!graphqlApis || !selectedOption) {
      console.error('Unable to add sub graph.');
      return;
    }
    // Get the selected api to add, and create a new sub graph reference to it.
    const apiToAdd = graphqlApis[selectedOption.apiIndex];
    const newSubGraph = {
      name: apiToAdd.metadata?.name ?? '',
      namespace: apiToAdd.metadata?.namespace ?? '',
      // typeMergeMap: Array<[string, StitchedSchema.SubschemaConfig.TypeMergeConfig.AsObject]>,
      typeMergeMap: [],
    };
    // Update the api with a new spec that includes the sub graphs.
    const existingSubGraphs =
      graphqlApi?.spec?.stitchedSchema?.subschemasList ?? [];
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
    mutateYaml();
    mutateDetails();
  };

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

  if (readonly) return null;
  return (
    <div>
      <SoloAddButton
        disabled={!graphqlApis}
        onClick={() => setIsModalVisible(true)}>
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

            <div className='text-right mt-5'>
              <SoloButtonStyledComponent
                disabled={!selectedOption}
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
