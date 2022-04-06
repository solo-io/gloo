import { TableActions, TableActionCircle } from 'Components/Common/SoloTable';
import { ReactComponent as XIcon } from 'assets/x-icon.svg';
import { StitchedSchema } from 'proto/github.com/solo-io/solo-apis/api/gloo/graphql.gloo/v1beta1/graphql_pb';
import React, { useState } from 'react';
import { graphqlConfigApi } from 'API/graphql';
import {
  usePageApiRef,
  useGetGraphqlApiDetails,
  useListGraphqlApis,
  useGetGraphqlApiYaml,
  usePageGlooInstance,
  useGetConsoleOptions,
} from 'API/hooks';
import ConfirmationModal from 'Components/Common/ConfirmationModal';

const StitchedGqlRemoveSubGraph: React.FC<{
  subGraphConfig: StitchedSchema.SubschemaConfig.AsObject;
}> = ({ subGraphConfig }) => {
  const apiRef = usePageApiRef();
  const { mutate: mutateYaml } = useGetGraphqlApiYaml(apiRef);
  const { data: graphqlApi, mutate: mutateDetails } =
    useGetGraphqlApiDetails(apiRef);
  const [glooInstance] = usePageGlooInstance();
  const { data: graphqlApis } = useListGraphqlApis(glooInstance?.metadata);
  const { readonly } = useGetConsoleOptions();

  const [isModalVisible, setIsModalVisible] = useState(false);

  const removeSubGraph = async () => {
    if (!graphqlApis) {
      console.error('Unable to add sub graph.');
      return;
    }
    const existingSubGraphs =
      graphqlApi?.spec?.stitchedSchema?.subschemasList ?? [];
    const newSubschemasList = existingSubGraphs.filter(
      g =>
        !(
          g.name === subGraphConfig.name &&
          g.namespace === subGraphConfig.namespace
        )
    );
    // Updates the api with a new spec that filters out the selected sub-graph.
    await graphqlConfigApi.updateGraphqlApi({
      graphqlApiRef: apiRef,
      spec: {
        stitchedSchema: {
          subschemasList: newSubschemasList,
        },
        allowedQueryHashesList: [],
      },
    });
    setIsModalVisible(false);
    mutateYaml();
    mutateDetails();
  };

  if (readonly) return null;
  return (
    <TableActions>
      <TableActionCircle onClick={() => setIsModalVisible(true)}>
        <XIcon />
      </TableActionCircle>
      <ConfirmationModal
        visible={isModalVisible}
        confirmPrompt='remove this sub graph?'
        confirmButtonText='Remove'
        goForIt={removeSubGraph}
        cancel={() => setIsModalVisible(false)}
        isNegative
      />
    </TableActions>
  );
};

export default StitchedGqlRemoveSubGraph;
