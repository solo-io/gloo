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
import { ReactComponent as GatewayIcon } from 'assets/gateway.svg';
import { ReactComponent as DownloadIcon } from 'assets/download-icon.svg';
import { useNavigate } from 'react-router';
import { useListFederatedGateways } from 'API/hooks';
import { Loading } from 'Components/Common/Loading';
import { federatedGlooResourceApi } from 'API/federated-gloo';
import { doDownload } from 'download-helper';
import { IconHolder } from 'Styles/StyledComponents/icons';
import { colors } from 'Styles/colors';
import { FederatedGateway } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/federated_gateway_resources_pb';
import { federatedGatewayResourceApi } from 'API/federated-gateway';
import { DataError } from 'Components/Common/DataError';

const GatewayIconHolder = styled.div`
  display: flex;
  align-items: center;
  justify-items: center;

  svg {
    width: 32px;
    margin-top: 2px;

    * {
      fill: ${colors.seaBlue};
    }
  }
`;

type GatewayTableFields = {
  key: string;
  name: string;
  namespace: string;
  clusters: string[];
  inNamespaces: string[];
  status: number;
  actions: FederatedGateway.AsObject;
};

// Named uniquely from the others  to avoid conflict with GRPC Class
export const FederatedGateways = () => {
  const [tableData, setTableData] = React.useState<GatewayTableFields[]>([]);

  const { data: gateways, error: fedGError } = useListFederatedGateways();

  useEffect(() => {
    if (gateways) {
      setTableData(
        gateways
        .sort((gA, gB) => (gA.metadata?.name ?? '').localeCompare(gB.metadata?.name ?? '') || (gA.metadata?.namespace ?? '').localeCompare(gB.metadata?.namespace ?? ''))
        .map(gateway => {
          return {
            key: gateway.metadata?.uid ?? 'An gateway was provided with no UID',
            name: gateway.metadata?.name ?? '',
            namespace: gateway.metadata?.namespace ?? '',
            clusters: gateway.spec?.placement?.clustersList ?? [],
            inNamespaces: gateway.spec?.placement?.namespacesList ?? [],
            status: gateway.status?.placementStatus?.state ?? 0,
            actions: gateway,
          };
        })
      );
    } else {
      setTableData([]);
    }
  }, [gateways]);

  if (!!fedGError) {
    return <DataError error={fedGError} />;
  } else if (!gateways) {
    return <Loading message={`Retrieving gateways...`} />;
  }

  const onDownloadGateway = (gateway: FederatedGateway.AsObject) => {
    if (gateway.metadata) {
      federatedGatewayResourceApi
        .getFederatedGatewayYAML({
          name: gateway.metadata.name!,
          namespace: gateway.metadata.namespace!,
        })
        .then(gatewayYaml => {
          doDownload(
            gatewayYaml,
            gateway.metadata?.namespace +
              '--' +
              gateway.metadata?.name +
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
      render: (gateway: FederatedGateway.AsObject) => (
        <TableActions>
          <TableActionCircle onClick={() => onDownloadGateway(gateway)}>
            <DownloadIcon />
          </TableActionCircle>
        </TableActions>
      ),
    },
  ];

  return (
    <SectionCard
      cardName={'Federated Gateways'}
      logoIcon={
        <GatewayIconHolder>
          <GatewayIcon />
        </GatewayIconHolder>
      }
      noPadding={true}>
      {!!fedGError ? (
        <DataError error={fedGError} />
      ) : !gateways ? (
        <Loading message={'Retrieving federated Gateways...'} />
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
