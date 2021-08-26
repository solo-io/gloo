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
import { ReactComponent as GlooIcon } from 'assets/Gloo.svg';
import { ReactComponent as DownloadIcon } from 'assets/download-icon.svg';
import { colors } from 'Styles/colors';
import { useParams, useNavigate } from 'react-router';
import {
  useIsGlooFedEnabled,
  useListClusterDetails,
  useListGlooInstances,
  useListVirtualServices,
} from 'API/hooks';
import { VirtualService } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/gateway_resources_pb';
import { Loading } from 'Components/Common/Loading';
import { objectMetasAreEqual } from 'API/helpers';
import { SimpleLinkProps, RenderSimpleLink } from 'Components/Common/SoloLink';
import { VirtualServiceStatus } from 'proto/github.com/solo-io/solo-apis/api/gloo/gateway/v1/virtual_service_pb';
import { gatewayResourceApi } from 'API/gateway-resources';
import { doDownload } from 'download-helper';
import { DataError } from 'Components/Common/DataError';
import { EmptyAsterisk } from 'Components/Common/EmptyAsterisk';

const GlooIconHolder = styled.div`
  display: flex;
  align-items: center;
  justify-items: center;

  svg {
    width: 25px;
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

type VirtualServiceTableFields = {
  key: string;
  name: SimpleLinkProps;
  domain: React.ReactNode;
  namespace: string;
  glooInstance?: { name: string; namespace: string };
  cluster?: string;
  routes: number;
  status: number;
  actions: VirtualService.AsObject;
};

type Props = {
  statusFilter?: VirtualServiceStatus.StateMap[keyof VirtualServiceStatus.StateMap];
  nameFilter?: string;
  glooInstanceFilter?: {
    name: string;
    namespace: string;
  };
};
export const VirtualServicesTable = (props: Props & TableHolderProps) => {
  const { name, namespace } = useParams();
  const navigate = useNavigate();

  const [tableData, setTableData] = React.useState<VirtualServiceTableFields[]>(
    []
  );

  const { data: glooInstances, error: glooError } = useListGlooInstances();
  const { data: clusterDetailsList, error: cError } = useListClusterDetails();
  const {
    data: glooFedCheckResponse,
    error: glooFedCheckError,
  } = useIsGlooFedEnabled();
  const isGlooFedEnabled = glooFedCheckResponse?.enabled;

  const multipleClustersOrInstances =
    (clusterDetailsList && clusterDetailsList.length > 1) ||
    (glooInstances && glooInstances.length > 1);

  const {
    data: virtualServices,
    error: virtualServicesError,
  } = useListVirtualServices(
    !!name && !!namespace
      ? {
          name,
          namespace,
        }
      : undefined
  );

  useEffect(() => {
    if (virtualServices) {
      setTableData(
        virtualServices
          .filter(
            vs =>
              vs.metadata?.name.includes(props.nameFilter ?? '') &&
              (props.statusFilter === undefined ||
                vs.status?.state === props.statusFilter) &&
              (!props.glooInstanceFilter ||
                objectMetasAreEqual(
                  {
                    name: vs.glooInstance?.name ?? '',
                    namespace: vs.glooInstance?.namespace ?? '',
                  },
                  props.glooInstanceFilter
                ))
          )
          .sort(
            (gA, gB) =>
              (gA.metadata?.name ?? '').localeCompare(
                gB.metadata?.name ?? ''
              ) ||
              (!props.wholePage
                ? 0
                : (gA.glooInstance?.name ?? '').localeCompare(
                    gB.glooInstance?.name ?? ''
                  ))
          )
          .map(vs => {
            let dataItem: VirtualServiceTableFields = {
              key:
                vs.metadata?.uid ??
                'An virtual service was provided with no UID',
              name: {
                displayElement: vs.metadata?.name ?? '',
                link: vs.metadata
                  ? isGlooFedEnabled
                    ? `/gloo-instances/${vs.glooInstance?.namespace}/${vs.glooInstance?.name}/virtual-services/${vs.metadata.clusterName}/${vs.metadata.namespace}/${vs.metadata.name}/`
                    : `/gloo-instances/${vs.glooInstance?.namespace}/${vs.glooInstance?.name}/virtual-services/${vs.metadata.namespace}/${vs.metadata.name}/`
                  : '',
              },
              namespace: vs.metadata?.namespace ?? '',
              domain: vs.spec?.virtualHost?.domainsList ? (
                vs.spec?.virtualHost?.domainsList.length === 1 &&
                vs.spec?.virtualHost?.domainsList[0] === '*' ? (
                  <EmptyAsterisk />
                ) : (
                  vs.spec?.virtualHost?.domainsList.join(', ')
                )
              ) : (
                ''
              ),
              routes: vs.spec?.virtualHost?.routesList.length ?? 0,
              status: vs.status?.state ?? 0,
              actions: vs,
            };

            if (props.wholePage) {
              dataItem['glooInstance'] = vs.glooInstance;
              dataItem['cluster'] = vs.metadata?.clusterName ?? '';
            }

            return dataItem;
          })
      );
    } else {
      setTableData([]);
    }
  }, [
    virtualServices,
    props.nameFilter,
    props.statusFilter,
    props.glooInstanceFilter,
  ]);

  if (!!virtualServicesError) {
    return <DataError error={virtualServicesError} />;
  } else if (!virtualServices) {
    return <Loading message={'Retrieving virtual services...'} />;
  }

  const onDownloadVirtualService = (vs: VirtualService.AsObject) => {
    if (vs.metadata) {
      gatewayResourceApi
        .getVirtualServiceYAML({
          name: vs.metadata.name!,
          namespace: vs.metadata.namespace!,
          clusterName: vs.metadata.clusterName,
        })
        .then(vsYaml => {
          doDownload(
            vsYaml,
            vs.metadata?.namespace + '--' + vs.metadata?.name + '.yaml'
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

  let columns: any = [
    {
      title: 'Name',
      dataIndex: 'name',
      width: 200,
      render: RenderSimpleLink,
    },
    {
      title: 'Domain',
      dataIndex: 'domain',
    },
    {
      title: 'Namespace',
      dataIndex: 'namespace',
    },
    {
      title: 'Routes',
      dataIndex: 'routes',
    },
    {
      title: 'Status',
      dataIndex: 'status',
      render: RenderStatus,
    },

    {
      title: 'Actions',
      dataIndex: 'actions',
      render: (vs: VirtualService.AsObject) => (
        <TableActions>
          <TableActionCircle onClick={() => onDownloadVirtualService(vs)}>
            <DownloadIcon />
          </TableActionCircle>
        </TableActions>
      ),
    },
  ];
  if (props.wholePage && multipleClustersOrInstances) {
    columns.splice(3, 0, {
      title: 'Cluster',
      dataIndex: 'cluster',
      render: RenderCluster,
    });
  }
  if (props.wholePage) {
    columns.splice(3, 0, {
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

export const VirtualServicesPageTable = (props: Props) => {
  return (
    <SectionCard
      cardName={'Virtual Services'}
      logoIcon={
        <GlooIconHolder>
          <GlooIcon />
        </GlooIconHolder>
      }
      noPadding={true}>
      <VirtualServicesTable {...props} wholePage={true} />
    </SectionCard>
  );
};
