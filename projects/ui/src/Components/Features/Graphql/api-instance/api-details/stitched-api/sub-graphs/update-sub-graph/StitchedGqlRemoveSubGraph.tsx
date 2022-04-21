import { graphqlConfigApi } from 'API/graphql';
import {
  useGetConsoleOptions,
  useGetGraphqlApiDetails,
  useListGraphqlApis,
  usePageApiRef,
  usePageGlooInstance,
} from 'API/hooks';
import Tooltip from 'antd/lib/tooltip';
import { ReactComponent as XIcon } from 'assets/x-icon.svg';
import { TableActionCircle, TableActions } from 'Components/Common/SoloTable';
import { useConfirm } from 'Components/Context/ConfirmModalContext';
import { StitchedSchema } from 'proto/github.com/solo-io/solo-apis/api/gloo/graphql.gloo/v1beta1/graphql_pb';
import React from 'react';
import toast from 'react-hot-toast';
import { hotToastError } from 'utils/hooks';

const StitchedGqlRemoveSubGraph: React.FC<{
  subGraphConfig: StitchedSchema.SubschemaConfig.AsObject;
  onAfterRemove(): void;
}> = ({ subGraphConfig, onAfterRemove }) => {
  const apiRef = usePageApiRef();
  const { data: graphqlApi } = useGetGraphqlApiDetails(apiRef);
  const { glooInstance } = usePageGlooInstance();
  const { data: graphqlApis } = useListGraphqlApis(glooInstance?.metadata);
  const { readonly } = useGetConsoleOptions();
  const confirm = useConfirm();

  const removeSubGraph = async () => {
    if (!graphqlApis) throw new Error('Unable to add sub graph.');
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
    onAfterRemove();
  };

  if (readonly) return null;
  return (
    <Tooltip title='Remove sub graph'>
      <TableActions className={`${subGraphConfig.name}-actions`}>
        <TableActionCircle
          data-testid='remove-sub-graph'
          onClick={() =>
            confirm({
              confirmPrompt: 'remove this sub graph?',
              confirmButtonText: 'Remove',
              isNegative: true,
            }).then(() =>
              toast.promise(removeSubGraph(), {
                loading: 'Removing sub graph...',
                success: 'Sub graph removed!',
                error: hotToastError,
              })
            )
          }>
          <XIcon />
        </TableActionCircle>
      </TableActions>
    </Tooltip>
  );
};

export default StitchedGqlRemoveSubGraph;
