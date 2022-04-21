import { ColumnsType } from 'antd/lib/table';
import Tooltip from 'antd/lib/tooltip';
import { graphqlConfigApi } from 'API/graphql';
import { useGetConsoleOptions, useIsGlooFedEnabled } from 'API/hooks';
import { ReactComponent as DownloadIcon } from 'assets/download-icon.svg';
import { ReactComponent as GraphQLIcon } from 'assets/graphql-icon.svg';
import { ReactComponent as XIcon } from 'assets/x-icon.svg';
import { SectionCard } from 'Components/Common/SectionCard';
import { RenderSimpleLink } from 'Components/Common/SoloLink';
import {
  RenderStatus,
  SoloTable,
  TableActionCircle,
  TableActions,
} from 'Components/Common/SoloTable';
import { doDownload } from 'download-helper';
import { GraphQLApiStatus } from 'proto/github.com/solo-io/solo-apis/api/gloo/graphql.gloo/v1beta1/graphql_pb';
import { GlooInstance } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/glooinstance_pb';
import { GraphqlApi } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/graphql_pb';
import React, { useEffect, useState } from 'react';
import {
  isExecutableAPI,
  isStitchedAPI,
  makeGraphqlApiLink,
  makeGraphqlApiRef,
} from 'utils/graphql-helpers';
import { useDeleteApi } from 'utils/hooks';
import * as styles from '../GraphqlLanding.style';
import { NewApiButton } from '../new-api-modal/NewApiButton';

// *
// * This table component takes care of the transofrmation and
// * rendering (not the filtering), of graphqlApis passed to it.
// *
export const GraphqlLandingTable: React.FC<{
  title: string;
  glooInstance: GlooInstance.AsObject;
  graphqlApis: GraphqlApi.AsObject[];
}> = ({ title, glooInstance, graphqlApis }) => {
  const isGlooFedEnabled = useIsGlooFedEnabled().data?.enabled;
  const { readonly } = useGetConsoleOptions();
  const confirmDeleteApi = useDeleteApi();

  // --- TRANSFORM GRAPHQL -> TABLE DATA --- //
  const [tableData, setTableData] = useState<any[]>([]);
  useEffect(() => {
    // Map object properties to column names.
    const newTableData = graphqlApis.map(api => ({
      key: `${api.metadata?.uid}`,
      name: {
        displayElement: api.metadata?.name ?? '',
        link: makeGraphqlApiLink(
          api.metadata?.name,
          api.metadata?.namespace,
          api.metadata?.clusterName ?? '',
          api.glooInstance?.name ?? '',
          api.glooInstance?.namespace ?? '',
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
  }, [graphqlApis, isGlooFedEnabled]);

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
          <Tooltip title='Download schema'>
            <TableActionCircle
              data-testid='graphql-table-action-download'
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
          </Tooltip>
          <Tooltip title='Delete schema'>
            <TableActionCircle
              data-testid='graphql-table-action-delete'
              onClick={() => confirmDeleteApi(makeGraphqlApiRef(api))}>
              <XIcon />
            </TableActionCircle>
          </Tooltip>
        </TableActions>
      ),
    },
  ];

  return (
    <div className='pb-5'>
      <SectionCard
        key={title}
        cardName={title}
        logoIcon={
          <styles.GraphqlIconHolder>
            <GraphQLIcon />
          </styles.GraphqlIconHolder>
        }
        noPadding={true}>
        <styles.TableHolder wholePage={true}>
          {!readonly && <NewApiButton glooInstance={glooInstance} />}
          <SoloTable
            columns={columns}
            dataSource={tableData}
            removePaging
            flatTopped
            removeShadows
          />
        </styles.TableHolder>
      </SectionCard>
    </div>
  );
};
