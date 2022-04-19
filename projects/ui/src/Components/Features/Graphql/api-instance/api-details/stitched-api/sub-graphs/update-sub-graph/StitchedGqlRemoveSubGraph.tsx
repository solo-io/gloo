import { graphqlConfigApi } from 'API/graphql';
import {
  useGetConsoleOptions,
  useGetGraphqlApiDetails,
  useListGraphqlApis,
  usePageApiRef,
  usePageGlooInstance,
} from 'API/hooks';
import { ReactComponent as XIcon } from 'assets/x-icon.svg';
import ConfirmationModal from 'Components/Common/ConfirmationModal';
import { TableActionCircle, TableActions } from 'Components/Common/SoloTable';
import { StitchedSchema } from 'proto/github.com/solo-io/solo-apis/api/gloo/graphql.gloo/v1beta1/graphql_pb';
import React, { useState } from 'react';

const StitchedGqlRemoveSubGraph: React.FC<{
  subGraphConfig: StitchedSchema.SubschemaConfig.AsObject;
  onAfterRemove(): void;
}> = ({ subGraphConfig, onAfterRemove }) => {
  const apiRef = usePageApiRef();
  const { data: graphqlApi } = useGetGraphqlApiDetails(apiRef);
  const { glooInstance } = usePageGlooInstance();
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
    onAfterRemove();
  };

  if (readonly) return null;
  return (
    <TableActions className={`${subGraphConfig.name}-actions`}>
      <TableActionCircle
        data-testid='remove-sub-graph'
        onClick={() => setIsModalVisible(true)}>
        <XIcon />
      </TableActionCircle>
      <ConfirmationModal
        visible={isModalVisible}
        confirmPrompt='remove this sub graph?'
        confirmTestId='confirm-remove-sub-graph'
        confirmButtonText='Remove'
        goForIt={removeSubGraph}
        cancel={() => setIsModalVisible(false)}
        isNegative
      />
    </TableActions>
  );
};

export default StitchedGqlRemoveSubGraph;
