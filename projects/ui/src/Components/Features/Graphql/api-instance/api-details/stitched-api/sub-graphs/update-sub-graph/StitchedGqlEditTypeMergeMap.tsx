import toast from 'react-hot-toast';
import { EditFilled } from '@ant-design/icons';
import { graphqlConfigApi } from 'API/graphql';
import {
  useGetConsoleOptions,
  useGetGraphqlApiDetails,
  usePageApiRef,
} from 'API/hooks';
import { SoloModal } from 'Components/Common/SoloModal';
import { TableActionCircle, TableActions } from 'Components/Common/SoloTable';
import { StitchedSchema } from 'proto/github.com/solo-io/solo-apis/api/gloo/graphql.gloo/v1beta1/graphql_pb';
import React, { useMemo, useState } from 'react';
import { SoloButtonStyledComponent } from 'Styles/StyledComponents/button';
import StitchedGqlTypeMergeMapConfig from '../type-merge-map/StitchedGqlTypeMergeMapConfig';
import { hotToastError } from 'utils/hooks';

const StitchedGqlEditTypeMergeMap: React.FC<{
  subGraphConfig: StitchedSchema.SubschemaConfig.AsObject;
  onAfterEdit(): void;
}> = ({ subGraphConfig, onAfterEdit }) => {
  const apiRef = usePageApiRef();
  const { data: graphqlApi } = useGetGraphqlApiDetails(apiRef);
  const [isModalVisible, setIsModalVisible] = useState(false);
  const { readonly } = useGetConsoleOptions();

  const existingSubGraph = useMemo(
    () =>
      graphqlApi?.spec?.stitchedSchema?.subschemasList?.find(
        s =>
          s.name === subGraphConfig.name &&
          s.namespace === subGraphConfig.namespace
      ),
    [graphqlApi]
  );

  // --- TYPE MERGE MAP --- //
  const [typeMergeMap, setTypeMergeMap] = useState<
    [string, StitchedSchema.SubschemaConfig.TypeMergeConfig.AsObject][]
  >([]);
  const saveTypeMergeMap = async () => {
    const subschemasList = graphqlApi?.spec?.stitchedSchema?.subschemasList;
    if (
      !isMergeMapValid ||
      !graphqlApi?.spec ||
      subschemasList === undefined ||
      existingSubGraph === undefined
    ) {
      console.error('Unable to update type merge map.');
      return;
    }
    //
    // Update the object reference in the list, and call the api.
    existingSubGraph.typeMergeMap = typeMergeMap;
    await toast.promise(
      graphqlConfigApi.updateGraphqlApi({
        graphqlApiRef: apiRef,
        spec: {
          ...graphqlApi.spec,
          stitchedSchema: { subschemasList },
        },
      }),
      {
        loading: 'Updating API...',
        success: 'API updated!',
        error: hotToastError,
      }
    );
    setIsModalVisible(false);
    onAfterEdit();
  };

  // --- VALIDATION --- //
  const [isMergeMapValid, setIsMergeMapValid] = useState(false);
  const canSubmit = isMergeMapValid && !!subGraphConfig;

  if (readonly) return null;
  return (
    <TableActions>
      <TableActionCircle onClick={() => setIsModalVisible(true)}>
        <EditFilled />
      </TableActionCircle>

      <SoloModal
        visible={isModalVisible}
        onClose={() => setIsModalVisible(false)}
        title={`Editing ${subGraphConfig.name}`}
        width={600}>
        <div className='p-5 pb-10 pt-0'>
          {!!subGraphConfig && (
            <StitchedGqlTypeMergeMapConfig
              onIsValidChange={isValid => setIsMergeMapValid(isValid)}
              initialTypeMergeMap={existingSubGraph?.typeMergeMap ?? []}
              onTypeMergeMapChange={m => setTypeMergeMap(m)}
              subGraphqlApiRef={{
                name: subGraphConfig.name ?? '',
                namespace: subGraphConfig.namespace ?? '',
                clusterName: apiRef.clusterName ?? '',
              }}
            />
          )}

          <div className='text-right mt-10'>
            <SoloButtonStyledComponent
              disabled={!canSubmit}
              onClick={saveTypeMergeMap}>
              Update Type Merge Map
            </SoloButtonStyledComponent>
          </div>
        </div>
      </SoloModal>
    </TableActions>
  );
};

export default StitchedGqlEditTypeMergeMap;
