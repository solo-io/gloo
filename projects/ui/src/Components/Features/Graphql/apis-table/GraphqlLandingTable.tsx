import { graphqlConfigApi } from 'API/graphql';
import { ReactComponent as DownloadIcon } from 'assets/download-icon.svg';
import {
  useIsGlooFedEnabled,
  useListGraphqlApis,
  usePageGlooInstance,
} from 'API/hooks';
import { ReactComponent as GraphQLIcon } from 'assets/graphql-icon.svg';
import { ReactComponent as XIcon } from 'assets/x-icon.svg';
import ConfirmationModal from 'Components/Common/ConfirmationModal';
import ErrorModal from 'Components/Common/ErrorModal';
import { SectionCard } from 'Components/Common/SectionCard';
import { RenderSimpleLink } from 'Components/Common/SoloLink';
import {
  RenderStatus,
  SoloTable,
  TableActionCircle,
  TableActions,
} from 'Components/Common/SoloTable';
import { doDownload } from 'download-helper';
import { GraphqlApi } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/graphql_pb';
import React, { useEffect, useState } from 'react';
import {
  isExecutableAPI,
  isStitchedAPI,
  makeGraphqlApiLink,
  makeGraphqlApiRef,
} from 'utils/graphql-helpers';
import { useDeleteAPI } from 'utils/hooks';
import * as styles from '../GraphqlLanding.style';
import { ColumnsType } from 'antd/lib/table';
import { GraphQLApiStatus } from 'proto/github.com/solo-io/solo-apis/api/gloo/graphql.gloo/v1beta1/graphql_pb';

// *
// * This table component takes care of the transofrmation and
// * rendering (not the filtering), of graphqlApis passed to it.
// *
export const GraphqlLandingTable: React.FC<{
  graphqlApis: GraphqlApi.AsObject[];
}> = ({ graphqlApis }) => {
  const isGlooFedEnabled = useIsGlooFedEnabled().data?.enabled;
  const [glooInstance] = usePageGlooInstance();
  const { mutate } = useListGraphqlApis();
  const {
    isDeleting,
    triggerDelete,
    cancelDelete,
    closeErrorModal,
    errorModalIsOpen,
    errorDeleteModalProps,
    deleteFn,
  } = useDeleteAPI({
    revalidate: mutate,
    optimistic: true,
  });

  // --- TRANSFORM GRAPHQL -> TABLE DATA --- //
  const [tableData, setTableData] = useState<any[]>([]);
  useEffect(() => {
    if (!glooInstance) return;
    // Map object properties to column names.
    const newTableData = graphqlApis.map(api => ({
      key: `${api.metadata?.uid}`,
      name: {
        displayElement: api.metadata?.name ?? '',
        link: makeGraphqlApiLink(
          api.metadata?.name,
          api.metadata?.namespace,
          glooInstance?.metadata?.clusterName ?? '',
          glooInstance?.metadata?.name ?? '',
          glooInstance?.metadata?.namespace ?? '',
          isGlooFedEnabled
        ),
      },
      namespace: api.metadata?.namespace,
      cluster: api.metadata?.clusterName ?? '',
      // TODO: Update this when pending/accepted states return correctly.
      status:
        api.status?.state === undefined ||
        api.status?.state === GraphQLApiStatus.State.PENDING
          ? GraphQLApiStatus.State.ACCEPTED
          : api.status?.state,
      apiType: isExecutableAPI(api) ? 'Executable' : 'Stitched',
      resolvers: isExecutableAPI(api)
        ? (api as any).executable.numResolvers
        : '',
      subGraphs: isStitchedAPI(api) ? (api as any).stitched.numSubschemas : '',
      api,
    }));
    setTableData(newTableData);
  }, [graphqlApis, isGlooFedEnabled, glooInstance]);

  // --- COLUMNS --- //
  let columns: ColumnsType = [
    {
      title: 'Name',
      dataIndex: 'name',
      width: 200,
      render: RenderSimpleLink,
    },
    {
      title: 'Namespace',
      dataIndex: 'namespace',
    },
    {
      title: 'Type',
      dataIndex: 'apiType',
    },
    {
      title: 'Resolvers',
      dataIndex: 'resolvers',
      align: 'center',
    },
    {
      title: 'Sub Graphs',
      dataIndex: 'subGraphs',
      align: 'center',
    },
    {
      title: 'Status',
      dataIndex: 'status',
      render: RenderStatus,
    },
    {
      title: 'Actions',
      dataIndex: 'api',
      align: 'center',
      render: (api: GraphqlApi.AsObject) => (
        <TableActions className='space-x-3 justify-center'>
          <TableActionCircle
            onClick={() => {
              if (!api.metadata) return;
              graphqlConfigApi
                .getGraphqlApiYaml(makeGraphqlApiRef(api))
                .then(gqlApiYaml => {
                  doDownload(
                    gqlApiYaml,
                    api.metadata?.namespace +
                      '--' +
                      api.metadata?.name +
                      '.yaml'
                  );
                });
            }}>
            <DownloadIcon />
          </TableActionCircle>
          <TableActionCircle
            onClick={() => triggerDelete(makeGraphqlApiRef(api))}>
            <XIcon />
          </TableActionCircle>
        </TableActions>
      ),
    },
  ];

  return (
    <SectionCard
      key='graphql'
      cardName='GraphQL'
      logoIcon={
        <styles.GraphqlIconHolder>
          <GraphQLIcon />
        </styles.GraphqlIconHolder>
      }
      noPadding={true}>
      <styles.TableHolder wholePage={true}>
        <SoloTable
          columns={columns}
          dataSource={tableData}
          removePaging
          flatTopped
          removeShadows
        />
        <ConfirmationModal
          visible={isDeleting}
          confirmPrompt='delete this API'
          confirmButtonText='Delete'
          goForIt={deleteFn}
          cancel={cancelDelete}
          isNegative
        />
        <ErrorModal
          {...errorDeleteModalProps}
          cancel={closeErrorModal}
          visible={errorModalIsOpen}
        />
      </styles.TableHolder>
    </SectionCard>
  );
};
