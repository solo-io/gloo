import React, { useEffect } from 'react';
import styled from '@emotion/styled/macro';
import {
  SoloTable,
  RenderStatus,
  TableActionCircle,
  TableActions,
  RenderClustersList,
  RenderExpandableNamesList,
} from 'Components/Common/SoloTable';
import { SectionCard } from 'Components/Common/SectionCard';
import { ReactComponent as RouteIcon } from 'assets/route-icon.svg';
import { ReactComponent as DownloadIcon } from 'assets/download-icon.svg';
import { useNavigate } from 'react-router';
import { useListFederatedRouteTables } from 'API/hooks';
import { Loading } from 'Components/Common/Loading';
import { doDownload } from 'download-helper';
import { IconHolder } from 'Styles/StyledComponents/icons';
import { colors } from 'Styles/colors';
import { FederatedRouteTable } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/federated_gateway_resources_pb';
import { federatedGatewayResourceApi } from 'API/federated-gateway';
import { DataError } from 'Components/Common/DataError';

const RTIconHolder = styled.div`
  display: flex;
  align-items: center;
  justify-items: center;

  svg {
    width: 26px;
    margin-left: -6px;

    * {
      fill: ${colors.seaBlue};
    }
  }
`;

type RouteTableTableFields = {
  key: string;
  name: string;
  namespace: string;
  clusters: string[];
  inNamespaces: string[];
  status: number;
  actions: FederatedRouteTable.AsObject;
};

export const FederatedRouteTables = () => {
  const [tableData, setTableData] = React.useState<RouteTableTableFields[]>([]);

  const {
    data: routeTables,
    error: fedRTError,
  } = useListFederatedRouteTables();

  useEffect(() => {
    if (routeTables) {
      setTableData(
        routeTables
        .sort((gA, gB) => (gA.metadata?.name ?? '').localeCompare(gB.metadata?.name ?? '') || (gA.metadata?.namespace ?? '').localeCompare(gB.metadata?.namespace ?? ''))
        .map(routeTable => {
          return {
            key:
              routeTable.metadata?.uid ??
              'A route table was provided with no UID',
            name: routeTable.metadata?.name ?? '',
            namespace: routeTable.metadata?.namespace ?? '',
            clusters: routeTable.spec?.placement?.clustersList ?? [],
            inNamespaces: routeTable.spec?.placement?.namespacesList ?? [],
            status: routeTable.status?.placementStatus?.state ?? 0,
            actions: routeTable,
          };
        })
      );
    } else {
      setTableData([]);
    }
  }, [routeTables]);

  if (!!fedRTError) {
    return <DataError error={fedRTError} />;
  } else if (!routeTables) {
    return <Loading message={`Retrieving route tables...`} />;
  }

  const onDownloadRouteTable = (routeTable: FederatedRouteTable.AsObject) => {
    if (routeTable.metadata) {
      federatedGatewayResourceApi
        .getFederatedRouteTableYAML({
          name: routeTable.metadata.name!,
          namespace: routeTable.metadata.namespace!,
        })
        .then(routeTableYaml => {
          doDownload(
            routeTableYaml,
            routeTable.metadata?.namespace +
              '--' +
              routeTable.metadata?.name +
              '.yaml'
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
      render: (routeTable: FederatedRouteTable.AsObject) => (
        <TableActions>
          <TableActionCircle onClick={() => onDownloadRouteTable(routeTable)}>
            <DownloadIcon />
          </TableActionCircle>
        </TableActions>
      ),
    },
  ];

  return (
    <SectionCard
      cardName={'Federated Route Tables'}
      logoIcon={
        <RTIconHolder>
          <RouteIcon />
        </RTIconHolder>
      }
      noPadding={true}>
      {!!fedRTError ? (
        <DataError error={fedRTError} />
      ) : !routeTables ? (
        <Loading message={'Retrieving federated Route Tables...'} />
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
