import { graphqlConfigApi } from 'API/graphql';
import { useListGraphqlApis } from 'API/hooks';
import { useConfirm } from 'Components/Context/ConfirmModalContext';
import { ClusterObjectRef } from 'proto/github.com/solo-io/skv2/api/core/v1/core_pb';
import toast from 'react-hot-toast';
import { useNavigate } from 'react-router';

export const hotToastError = (e: any) => {
  toast.error(e?.message ?? e);
  return 'Error';
};

export function useDeleteApi() {
  const { mutate: mutateApiList } = useListGraphqlApis();
  const navigate = useNavigate();
  const confirm = useConfirm();
  return (apiRef: ClusterObjectRef.AsObject) =>
    confirm({
      confirmPrompt: `delete the API, ${apiRef.name}`,
      confirmButtonText: 'Delete',
      isNegative: true,
    }).then(() =>
      toast
        .promise(graphqlConfigApi.deleteGraphqlApi(apiRef), {
          loading: 'Deleting API...',
          success: () => {
            mutateApiList();
            return 'API Deleted!';
          },
          error: hotToastError,
        })
        .finally(() => navigate('/apis/'))
    );
}
