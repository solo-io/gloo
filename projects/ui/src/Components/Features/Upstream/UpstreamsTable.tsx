import React, { useEffect } from 'react';
import styled from '@emotion/styled/macro';
import {
  SoloTable,
  RenderStatus,
  TableActionCircle,
  TableActions,
  RenderCluster,
} from 'Components/Common/SoloTable';
import { SectionCard } from 'Components/Common/SectionCard';
import { ReactComponent as UpstreamIcon } from 'assets/upstream-icon.svg';
import { ReactComponent as DownloadIcon } from 'assets/download-icon.svg';
import { ReactComponent as FailoverIcon } from 'assets/GlooFed-Specific/failover-icon.svg';
import { colors } from 'Styles/colors';
import { useParams, useNavigate } from 'react-router';
import {
  useListClusterDetails,
  useListGlooInstances,
  useListUpstreams,
} from 'API/hooks';
import { Loading } from 'Components/Common/Loading';
import { objectMetasAreEqual } from 'API/helpers';
import { Upstream } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/gloo_resources_pb';
import { UpstreamStatus } from 'proto/github.com/solo-io/solo-apis/api/gloo/gloo/v1/upstream_pb';
import { SimpleLinkProps, RenderSimpleLink } from 'Components/Common/SoloLink';
import { glooResourceApi } from 'API/gloo-resource';
import { doDownload } from 'download-helper';
import { IconHolder } from 'Styles/StyledComponents/icons';
import { DataError } from 'Components/Common/DataError';
import { CheckboxFilterProps } from './UpstreamsLanding';
import { getUpstreamType } from 'utils/upstream-helpers';

const UpstreamIconHolder = styled.div`
  display: flex;
  align-items: center;
  justify-items: center;

  svg {
    width: 35px;
    max-width: none;
  }
`;

type TableHolderProps = { wholePage?: boolean };
const TableHolder = styled.div<TableHolderProps>`
  ${(props: TableHolderProps) =>
    props.wholePage
      ? ''
      : `
    table thead.ant-table-thead tr th {
      background: ${colors.marchGrey};
    }
  `};
`;

type Props = {
  statusFilter?: UpstreamStatus.StateMap[keyof UpstreamStatus.StateMap];
  nameFilter?: string;
  typeFilters?: CheckboxFilterProps[];
  glooInstanceFilter?: {
    name: string;
    namespace: string;
  };
};
export const UpstreamsTable = (props: Props & TableHolderProps) => {
  const { name, namespace } = useParams();
  const navigate = useNavigate();

  const [tableData, setTableData] = React.useState<
    {
      key: string;
      name: SimpleLinkProps;
      namespace: string;
      glooInstance?: { name: string; namespace: string };
      cluster: string;
      failover: boolean;
      status: number;
      actions: Upstream.AsObject;
    }[]
  >([]);

  const { data: upstreams, error: upstreamsError } = useListUpstreams(
    !!name && !!namespace
      ? {
          name,
          namespace,
        }
      : undefined
  );

  const { data: glooInstances, error: glooError } = useListGlooInstances();
  const { data: clusterDetailsList, error: cError } = useListClusterDetails();

  const multipleClustersOrInstances =
    (clusterDetailsList && clusterDetailsList.length > 1) ||
    (glooInstances && glooInstances.length > 1);

  useEffect(() => {
    if (upstreams) {
      let typeCheckboxesNotSet = props.typeFilters?.every(c => !c.checked);
      setTableData(
        upstreams
          .filter(
            upstream =>
              upstream.metadata?.name.includes(props.nameFilter ?? '') &&
              (props.statusFilter === undefined ||
                upstream.status?.state === props.statusFilter) &&
                (typeCheckboxesNotSet || props.typeFilters?.find(c => c.label === getUpstreamType(upstream))?.checked) &&
              (!props.glooInstanceFilter ||
                objectMetasAreEqual(
                  {
                    name: upstream.glooInstance?.name ?? '',
                    namespace: upstream.glooInstance?.namespace ?? '',
                  },
                  props.glooInstanceFilter
                ))
          )
          .sort((gA, gB) => (gA.metadata?.name ?? '').localeCompare(gB.metadata?.name ?? '') || (!props.wholePage ? 0 : (gA.glooInstance?.name ?? '').localeCompare(gB.glooInstance?.name ?? '')))
          .map(upstream => {
            const glooInstNamespace = upstream.glooInstance?.namespace;
            const glooInstName = upstream.glooInstance?.name;
            const upstreamCluster = upstream.metadata?.clusterName ?? '';
            const upstreamNamespace = upstream.metadata?.namespace ?? '';
            const upstreamName = upstream.metadata?.name ?? '';
            const link =
              glooInstNamespace &&
              glooInstName &&
              upstreamCluster &&
              upstreamNamespace &&
              upstreamName
                ? `/gloo-instances/${glooInstNamespace}/${glooInstName}/upstreams/${upstreamCluster}/${upstreamNamespace}/${upstreamName}`
                : '';
            return {
              key:
                upstream.metadata?.uid ??
                'An upstream was provided with no UID',
              name: {
                displayElement: upstreamName,
                link,
              },
              namespace: upstreamNamespace,
              glooInstance: upstream.glooInstance,
              cluster: upstreamCluster,
              failover: !!upstream.spec?.failover,
              status: upstream.status?.state ?? UpstreamStatus.State.PENDING,
              actions: upstream,
            };
          })
      );
    } else {
      setTableData([]);
    }
  }, [
    upstreams,
    props.nameFilter,
    props.statusFilter,
    props.typeFilters,
    props.glooInstanceFilter,
  ]);

  if (!!upstreamsError) {
    return <DataError error={upstreamsError} />;
  } else if (!upstreams) {
    return <Loading message={'Retrieving upstreams...'} />;
  }

  const onDownloadUpstream = (upstream: Upstream.AsObject) => {
    if (upstream.metadata) {
      glooResourceApi
        .getUpstreamYAML({
          name: upstream.metadata.name,
          namespace: upstream.metadata.namespace,
          clusterName: upstream.metadata.clusterName,
        })
        .then(upstreamYaml => {
          doDownload(
            upstreamYaml,
            upstream.metadata?.namespace +
              '--' +
              upstream.metadata?.name +
              '.yaml'
          );
        });
    }
  };

  const renderGlooInstanceList = (glooInstance: {
    name: string;
    namespace: string;
  }) => {
    return (
      <div
        onClick={() =>
          navigate(
            `/gloo-instances/${glooInstance.namespace}/${glooInstance.name}/`
          )
        }>
        {glooInstance.name}
      </div>
    );
  };

  const renderFailover = (failoverExists: boolean) => {
    return failoverExists ? (
      <IconHolder>
        <FailoverIcon />
      </IconHolder>
    ) : (
      <React.Fragment />
    );
  };

  let columns: any = [
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
      title: 'Failover',
      dataIndex: 'failover',
      render: renderFailover,
    },
    {
      title: 'Status',
      dataIndex: 'status',
      render: RenderStatus,
    },

    {
      title: 'Actions',
      dataIndex: 'actions',
      render: (upstream: Upstream.AsObject) => (
        <TableActions>
          <TableActionCircle onClick={() => onDownloadUpstream(upstream)}>
            <DownloadIcon />
          </TableActionCircle>
        </TableActions>
      ),
    },
  ];

  if (props.wholePage && multipleClustersOrInstances) {
    columns.splice(2, 0, {
      title: 'Cluster',
      dataIndex: 'cluster',
      render: RenderCluster,
    });
  }
  if (props.wholePage) {
    columns.splice(2, 0, {
      title: 'Gloo Instance',
      dataIndex: 'glooInstance',
      render: renderGlooInstanceList,
    });
  }

  return (
    <TableHolder wholePage={props.wholePage}>
      <SoloTable
        columns={columns}
        dataSource={tableData}
        removePaging
        removeShadows
        curved={props.wholePage}
      />
    </TableHolder>
  );
};

export const UpstreamsPageTable = (props: Props) => {
  return (
    <SectionCard
      cardName={'Upstreams'}
      logoIcon={
        <UpstreamIconHolder>
          <UpstreamIcon />
        </UpstreamIconHolder>
      }
      noPadding={true}>
      <UpstreamsTable {...props} wholePage={true} />
    </SectionCard>
  );
};
