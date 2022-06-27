import styled from '@emotion/styled/macro';
import { getUpstreamDetails } from 'API/gloo-resource';
import { RenderSimpleLink } from 'Components/Common/SoloLink';
import { RenderCluster, SoloTable } from 'Components/Common/SoloTable';
import { GetUpstreamDetailsResponse } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/gloo_resources_pb';
import { FailoverSchemeSpec } from 'proto/github.com/solo-io/solo-projects/projects/gloo-fed/api/fed/v1/failover_pb';
import React, { useEffect } from 'react';
import { di } from 'react-magnetic-di/macro';
import { colors } from 'Styles/colors';
import { WeightPercentageBlock } from 'Styles/StyledComponents/weight-block';

const ROW_HEIGHT = 40;

const TableHolder = styled.div`
  margin-top: 15px;
  font-size: 14px;

  table thead.ant-table-thead tr th {
    background: ${colors.marchGrey};
  }

  table tbody.ant-table-tbody tr.cluster {
    background: ${colors.januaryGrey};
  }
`;

type ConnectorLineProps = {
  first?: boolean;
};
const ConnectorLine = styled.div<ConnectorLineProps>`
  border-left: 1px dotted ${colors.juneGrey};
  border-bottom: 1px dotted ${colors.juneGrey};
  width: 8px;
  height: ${(props: ConnectorLineProps) =>
    props.first ? ROW_HEIGHT / 2 : ROW_HEIGHT}px;
  margin-top: -${(props: ConnectorLineProps) => (props.first ? ROW_HEIGHT / 2 : ROW_HEIGHT)}px;
`;

const ConnectorCircle = styled.div`
  width: 5px;
  height: 5px;
  border-radius: 100%;
  border: 1px solid ${colors.juneGrey};
  margin-right: 8px;
`;

const UpstreamContainer = styled.div`
  display: flex;
  margin-left: 8px;
  align-items: center;
`;

enum RowType {
  CLUSTER = 0,
  UPSTREAM = 1,
}

type RowData = {
  key: string;
  rowType: RowType;
  cluster?: string;
  rawWeight?: number;
  weightPercentage?: number;
  upstreamIndex?: number;
  upstream?: {
    name: string;
    namespace: string;
    cluster?: string;
    glooInstanceName?: string;
    glooInstanceNamespace?: string;
  };
};

const RenderWeight = (rowData: RowData) => {
  if (rowData.rowType !== RowType.CLUSTER) {
    return null;
  }
  if (rowData.rawWeight === undefined) {
    return <WeightPercentageBlock>EQUAL</WeightPercentageBlock>;
  }
  return (
    <WeightPercentageBlock percentage={rowData.weightPercentage}>
      {rowData.rawWeight}
    </WeightPercentageBlock>
  );
};

const RenderName = (rowData: RowData) => {
  if (rowData.rowType === RowType.CLUSTER) {
    return RenderCluster(`Cluster ${rowData.cluster ?? ''}`);
  }

  const name = rowData.upstream?.name;
  const namespace = rowData.upstream?.namespace;
  const cluster = rowData.upstream?.cluster;
  const glooInstanceName = rowData.upstream?.glooInstanceName;
  const glooInstanceNamespace = rowData.upstream?.glooInstanceNamespace;
  const link =
    name &&
    namespace &&
    cluster &&
    glooInstanceName &&
    glooInstanceNamespace ? (
      <RenderSimpleLink
        displayElement={name}
        link={`/gloo-instances/${glooInstanceNamespace}/${glooInstanceName}/upstreams/${cluster}/${namespace}/${name}`}
        inline
      />
    ) : (
      <span>{name}</span>
    );

  return (
    <UpstreamContainer>
      <ConnectorLine first={rowData.upstreamIndex === 0} />
      <ConnectorCircle />
      {link}
    </UpstreamContainer>
  );
};

const RenderNamespace = (rowData: RowData) => {
  if (rowData.rowType === RowType.CLUSTER) {
    return null;
  }

  return rowData.upstream?.namespace;
};

type Props = {
  group: FailoverSchemeSpec.FailoverEndpoints.AsObject;
  isWeighted: boolean;
};

const COLUMNS: any = [
  {
    title: 'Weight',
    width: '200px',
    render: RenderWeight,
  },
  {
    title: 'Target',
    render: RenderName,
  },
  {
    title: 'Namespace',
    width: '250px',
    render: RenderNamespace,
  },
];

const UpstreamFailoverGroupTable = ({ group, isWeighted }: Props) => {
  di(getUpstreamDetails);
  const [tableData, setTableData] = React.useState<RowData[]>([]);

  useEffect(() => {
    (async () => {
      //
      // This makes a bunch of parallel requests to get the details for each referenced upstream.
      //
      // - First level is this group's priorityGroupList index
      // - Second level is that priorityGroupList's upstreamsList index
      //
      // So allUpstreamDetailsRequests[1][2] is the details request
      // for `group.priorityGroupList[1].upstreamList[2];`
      //
      const allUpstreamDetailsPromises =
        [] as Promise<GetUpstreamDetailsResponse.AsObject>[][];
      group.priorityGroupList.forEach(priorityGroup => {
        const pgUpstreamsDetailsPromises =
          [] as Promise<GetUpstreamDetailsResponse.AsObject>[];
        priorityGroup.upstreamsList.forEach(upstream => {
          pgUpstreamsDetailsPromises.push(
            getUpstreamDetails({
              clusterName: priorityGroup.cluster,
              name: upstream.name,
              namespace: upstream.namespace,
            })
          );
        });
        allUpstreamDetailsPromises.push(pgUpstreamsDetailsPromises);
      });
      //
      // Wait for the responses, (stored in the same way).
      const allUpstreamDetailsResponses =
        [] as GetUpstreamDetailsResponse.AsObject[][];
      for (let pgIdx = 0; pgIdx < group.priorityGroupList.length; pgIdx++) {
        allUpstreamDetailsResponses.push(
          await Promise.all(allUpstreamDetailsPromises[pgIdx])
        );
      }
      //
      // Prepare the table data.
      const tableData: RowData[] = [];
      const totalWeight = group.priorityGroupList
        ?.map(pgl => pgl.localityWeight?.value ?? 0)
        .reduce((sum, w) => sum + w, 0);
      group.priorityGroupList.forEach((priorityGroup, pgIdx) => {
        tableData.push({
          key: `cluster-${pgIdx}`,
          rowType: RowType.CLUSTER,
          cluster: priorityGroup.cluster,
          rawWeight: isWeighted
            ? priorityGroup.localityWeight?.value ?? 0
            : undefined,
          weightPercentage: isWeighted
            ? ((priorityGroup.localityWeight?.value ?? 0) * 100) / totalWeight
            : undefined,
        });
        priorityGroup.upstreamsList.forEach((upstream, uIdx) => {
          const upstreamDetails =
            allUpstreamDetailsResponses[pgIdx][uIdx].upstream;
          tableData.push({
            key: `cluster-${pgIdx}-upstream-${uIdx}`,
            rowType: RowType.UPSTREAM,
            upstreamIndex: uIdx,
            upstream: {
              name: upstream.name,
              namespace: upstream.namespace,
              cluster: priorityGroup.cluster,
              glooInstanceName: upstreamDetails?.glooInstance?.name ?? '',
              glooInstanceNamespace:
                upstreamDetails?.glooInstance?.namespace ?? '',
            },
          });
        });
      });
      setTableData(tableData);
    })();
  }, [group, isWeighted]);

  return (
    <TableHolder>
      <SoloTable
        rowHeight={`${ROW_HEIGHT}px`}
        columns={COLUMNS}
        dataSource={tableData}
        rowClassName={(rowData: RowData) =>
          rowData.rowType === RowType.CLUSTER ? 'cluster' : ''
        }
        removePaging
        removeShadows
        curved
        withBorder
      />
    </TableHolder>
  );
};

export default UpstreamFailoverGroupTable;
