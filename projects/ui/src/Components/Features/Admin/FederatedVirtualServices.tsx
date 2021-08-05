import React, { useEffect } from 'react';
import styled from '@emotion/styled/macro';
import {
  SoloTable,
  RenderStatus,
  TableActionCircle,
  TableActions,
  RenderCluster,
  RenderClustersList,
  RenderExpandableNamesList,
} from 'Components/Common/SoloTable';
import { SectionCard } from 'Components/Common/SectionCard';
import { ReactComponent as GlooIcon } from 'assets/Gloo.svg';
import { ReactComponent as DownloadIcon } from 'assets/download-icon.svg';
import { useNavigate } from 'react-router';
import { useListFederatedVirtualServices } from 'API/hooks';
import { VirtualService } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/gateway_resources_pb';
import { Loading } from 'Components/Common/Loading';
import { FederatedVirtualService } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/federated_gateway_resources_pb';
import { federatedGatewayResourceApi } from 'API/federated-gateway';
import { doDownload } from 'download-helper';
import { DataError } from 'Components/Common/DataError';

const GlooIconHolder = styled.div`
  display: flex;
  align-items: center;
  justify-items: center;

  svg {
    width: 25px;
  }
`;

type VirtualServiceTableFields = {
  key: string;
  name: string;
  namespace: string;
  clusters: string[];
  inNamespaces: string[];
  status: number;
  actions: FederatedVirtualService.AsObject;
};

type Props = {
  nameFilter?: string;
};
export const FederatedVirtualServices = (props: Props) => {
  const navigate = useNavigate();

  const [tableData, setTableData] = React.useState<VirtualServiceTableFields[]>(
    []
  );

  const {
    data: virtualServices,
    error: fedVSError,
  } = useListFederatedVirtualServices();

  useEffect(() => {
    if (virtualServices) {
      setTableData(
        virtualServices
          .filter(vs => vs.metadata?.name.includes(props.nameFilter ?? ''))
          .sort((gA, gB) => (gA.metadata?.name ?? '').localeCompare(gB.metadata?.name ?? '') || (gA.metadata?.namespace ?? '').localeCompare(gB.metadata?.namespace ?? ''))
          .map(vs => {
            let dataItem: VirtualServiceTableFields = {
              key:
                vs.metadata?.uid ??
                'An virtual service was provided with no UID',
              name: vs.metadata?.name ?? '',
              namespace: vs.metadata?.namespace ?? '',
              clusters: vs.spec?.placement?.clustersList ?? [],
              inNamespaces: vs.spec?.placement?.namespacesList ?? [],
              status: vs.status?.placementStatus?.state ?? 0,
              actions: vs,
            };

            return dataItem;
          })
      );
    } else {
      setTableData([]);
    }
  }, [virtualServices, props.nameFilter]);

  if (!!fedVSError) {
    return <DataError error={fedVSError} />;
  } else if (!virtualServices) {
    return <Loading message={`Retrieving virtual services...`} />;
  }

  const onDownloadVirtualService = (vs: VirtualService.AsObject) => {
    if (vs.metadata) {
      federatedGatewayResourceApi
        .getFederatedVirtualServiceYAML({
          name: vs.metadata.name!,
          namespace: vs.metadata.namespace!,
        })
        .then(vsYaml => {
          doDownload(
            vsYaml,
            vs.metadata?.namespace + '--' + vs.metadata?.name + '.yaml'
          );
        });
    }
  };

  let columns: any = [
    {
      title: 'Name',
      dataIndex: 'name',
      width: 200,
    },
    {
      title: 'Namespace',
      dataIndex: 'namespace',
    },
    {
      title: 'Placement Clusters',
      dataIndex: 'clusters',
      render: RenderClustersList,
    },
    {
      title: 'Placement Namespaces',
      dataIndex: 'inNamespaces',
      render: RenderExpandableNamesList,
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

  return (
    <SectionCard
      cardName={'Federated Virtual Services'}
      logoIcon={
        <GlooIconHolder>
          <GlooIcon />
        </GlooIconHolder>
      }
      noPadding={true}>
      {!!fedVSError ? (
        <DataError error={fedVSError} />
      ) : !virtualServices ? (
        <Loading message={'Retrieving federated Virtual Services...'} />
      ) : (
        <SoloTable
          columns={columns}
          dataSource={tableData}
          removePaging
          removeShadows
          curved={true}
        />
      )}
    </SectionCard>
  );
};
