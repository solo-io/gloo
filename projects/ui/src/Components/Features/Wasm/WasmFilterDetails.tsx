/* Should this really be in this folder? */
import React, { useState, useEffect } from 'react';
import styled from '@emotion/styled';
import { ReactComponent as FailoversIcon } from 'assets/filter-icon.svg';
import { ReactComponent as ViewIcon } from 'assets/view-icon.svg';
import { ReactComponent as WorkloadIcon } from 'assets/workload-icon.svg';
import { ReactComponent as DownloadIcon } from 'assets/download-icon.svg';
import { CardHeader } from 'Components/Common/Card';
import { HealthIndicator } from 'Components/Common/HealthIndicator';
import { DataError } from 'Components/Common/DataError';
import { Loading } from 'Components/Common/Loading';
import { SoloLink } from 'Components/Common/SoloLink';
import { useLocation, useParams } from 'react-router';
import { HealthNotificationBox } from 'Components/Common/HealthNotificationBox';
import { PolicyHeaderTitle } from './DetailModalHeader';
import YamlDisplayer from 'Components/Common/YamlDisplayer';
import { Tooltip } from 'antd';
import { WasmFilter } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/wasm_pb';
import { objectMetasAreEqual } from 'API/helpers';
import { colors } from 'Styles/colors';
import { useListGateways } from 'API/hooks';
import {
  RenderStatus,
  SoloTable,
  TableActionCircle,
  TableActions,
} from 'Components/Common/SoloTable';
import { ClusterObjectRef } from 'proto/github.com/solo-io/skv2/api/core/v1/core_pb';
import { ObjectMeta } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/common_pb';
import { gatewayResourceApi } from 'API/gateway-resources';
import { doDownload } from 'download-helper';

const FailoverDetailsContainer = styled.div`
  display: flex;
  flex-direction: column;
  min-height: 50vh;
  padding: 0 20px 20px;
`;

const DetailsHeader = styled(CardHeader)`
  position: relative;
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  background: white;
  line-height: normal;
  padding: 11px 110px 20px 0;
  border-bottom: 1px solid ${colors.februaryGrey};

  > div {
    display: flex;
    align-items: center;
    margin-bottom: 20px;
  }
`;

const HeaderImageHolder = styled.div`
  margin-right: 15px;
  height: 33px;
  width: 33px;
  min-width: 33px;
  border-radius: 100%;
  background: ${colors.seaBlue};
  display: flex;
  justify-content: center;
  align-items: center;

  img,
  svg {
    width: 18px;
    max-height: 18px;
    margin-top: 2px;

    * {
      fill: white;
    }
  }
`;

const HeaderTitleName = styled.div`
  font-size: 22px;
  color: ${colors.novemberGrey};
`;

const SecondaryInformation = styled.div`
  display: flex;
  flex-wrap: flex;
  align-items: center;
`;
type SecondaryInformationSectionProps = {
  noMaxWidth?: boolean;
};
const SecondaryInformationSection = styled.div<SecondaryInformationSectionProps>`
  display: flex;
  ${(props: SecondaryInformationSectionProps) =>
    props.noMaxWidth ? '' : 'max-width: 200px;'}
  font-size: 14px;
  line-height: 22px;
  height: 22px;
  padding: 0 12px;
  color: ${colors.novemberGrey};
  background: ${colors.februaryGrey};
  margin-left: 13px;
  border-radius: 16px;

  > span {
    overflow: hidden;
    white-space: nowrap;
    text-overflow: ellipsis;
  }
`;
const SecondaryInformationTitle = styled.div`
  font-weight: bold;
  margin-right: 4px;
  white-space: nowrap;
`;

const HealthContainer = styled.div`
  position: absolute;
  right: 35px;
  top: 18px;
  display: flex;
  justify-content: flex-end;
  align-items: center;
  flex: 1;
  text-align: right;
  font-size: 16px;
  font-weight: 600;
  color: ${colors.novemberGrey};
`;

const PolicyDetailsDescription = styled.div`
  border-radius: 8px;
  background: ${colors.januaryGrey};
  padding: 15px 13px;
  margin-bottom: 18px;
`;

const Title = styled.div`
  display: flex;
  font-size: 20px;
  line-height: 24px;
  font-weight: 500;
`;

const TableHolder = styled.div`
  margin-top: 24px;

  > div:nth-of-type(2) {
    margin-top: 18px;
    border: 1px solid ${colors.marchGrey};
    border-radius: 11px;
  }
`;

const DataTitle = styled(PolicyHeaderTitle)`
  margin-left: 0;
  margin-bottom: 18px;
`;

const RelativePlacer = styled.div`
  position: relative;
`;
const ImagePullViewer = styled.div`
  position: absolute;
  right: 0;
  top: 0;
  display: flex;
  align-items: center;
  color: ${colors.seaBlue};

  svg {
    margin-right: 8px;
  }
`;

const SubsetListItem = styled.div`
  display: inline-flex;
  font-size: 14px;
  line-height: 17px;
  border-radius: 16px;
  border: 1px solid ${colors.marchGrey};
  text-transform: lowercase;
  padding: 0 8px;
  margin-right: 8px;
`;

const SubsetItemTitle = styled.div`
  font-weight: 500;
  margin-right: 4px;
`;

const NameColumn = styled.div`
  display: flex;
  align-items: center;

  svg {
    margin-right: 8px;

    * {
      fill: ${colors.seaBlue};
    }
  }
`;

const Clickable = styled.div`
  display: flex;
  align-items: center;
  cursor: pointer;
`;

type GatewaysTableFields = {
  key: string;
  glooInstance: string;
  cluster: string;
  status: number;
  actions?: ClusterObjectRef.AsObject;
};

export const WasmFilterDetails = ({
  wasmFilter,
}: {
  wasmFilter: WasmFilter.AsObject;
}) => {
  const [imagePullVisible, setImagePullVisibile] = useState(false);
  const [gatewaysTableData, setGatewaysTableData] = useState<
    GatewaysTableFields[]
  >([]);

  useEffect(() => {
    const newGatewayTableData: GatewaysTableFields[] = wasmFilter.locationsList
      .map(location => {
        return {
          key:
            (location.gatewayRef?.name ?? '') +
            location.gatewayRef?.namespace +
            location.gatewayRef?.clusterName,
          glooInstance: location.glooInstanceRef?.name ?? 'Not found',
          cluster: location.gatewayRef?.clusterName ?? 'Not found',
          status: location.gatewayStatus?.state ?? 0,
          actions: location.gatewayRef,
        };
      })
      .sort((rowA, rowB) => rowA.key.localeCompare(rowB.key));

    setGatewaysTableData(newGatewayTableData);
  }, [wasmFilter.locationsList]);

  const onDownloadGateway = (gatewayRef: ClusterObjectRef.AsObject) => {
    gatewayResourceApi.getGatewayYAML(gatewayRef).then(gatewayYaml => {
      doDownload(
        gatewayYaml,
        gatewayRef.namespace +
          '-' +
          gatewayRef.name +
          '--' +
          gatewayRef.clusterName +
          '.yaml'
      );
    });
  };

  let gatewayColumns: any = [
    {
      title: 'Gloo Instance',
      dataIndex: 'glooInstance',
    },
    {
      title: 'Cluster',
      dataIndex: 'cluster',
    },
    {
      title: 'Status',
      dataIndex: 'status',
      render: RenderStatus,
    },
    {
      title: 'Actions',
      dataIndex: 'actions',
      render: (gateway: ClusterObjectRef.AsObject) => (
        <TableActions>
          {!!gateway && (
            <TableActionCircle onClick={() => onDownloadGateway(gateway)}>
              <DownloadIcon />
            </TableActionCircle>
          )}
        </TableActions>
      ),
    },
  ];

  return (
    <FailoverDetailsContainer>
      <DetailsHeader>
        <HeaderImageHolder>
          <FailoversIcon />
        </HeaderImageHolder>
        <HeaderTitleName title={wasmFilter.name}>
          {wasmFilter.name}
        </HeaderTitleName>
        <div>
          {!!wasmFilter.rootId.length && (
            <SecondaryInformation>
              <SecondaryInformationSection>
                <SecondaryInformationTitle>
                  <div>Root ID:</div>
                </SecondaryInformationTitle>
                <span>{wasmFilter.rootId}</span>
              </SecondaryInformationSection>
            </SecondaryInformation>
          )}
        </div>
      </DetailsHeader>

      <PolicyDetailsDescription>
        Wasm Filters allow Gloo Mesh administrators extend the functionality of
        the sidecar proxy attached to any Istio Workload.
      </PolicyDetailsDescription>

      <RelativePlacer>
        <DataTitle>Source</DataTitle>
        <PolicyDetailsDescription
          style={{ border: `1px solid ${colors.marchGrey}` }}>
          {wasmFilter.source ?? 'No source not found'}
        </PolicyDetailsDescription>
      </RelativePlacer>

      <RelativePlacer>
        <DataTitle>Config </DataTitle>
        {!wasmFilter.config ? (
          <div>No filter configuration provided</div>
        ) : (
          <YamlDisplayer contentString={wasmFilter.config} copyable={true} />
        )}
      </RelativePlacer>

      <TableHolder>
        <Title>
          Gateways
          {wasmFilter.locationsList[0]?.gatewayRef?.name && (
            <SecondaryInformation>
              <SecondaryInformationSection noMaxWidth={true}>
                <SecondaryInformationTitle>
                  <div>Gateways Name</div>
                </SecondaryInformationTitle>
                {wasmFilter.locationsList[0].gatewayRef.name}
              </SecondaryInformationSection>
            </SecondaryInformation>
          )}
          {wasmFilter.locationsList[0]?.gatewayRef?.namespace && (
            <SecondaryInformation>
              <SecondaryInformationSection noMaxWidth={true}>
                <SecondaryInformationTitle>
                  <div>Gateways Namespace</div>
                </SecondaryInformationTitle>
                {wasmFilter.locationsList[0].gatewayRef.namespace}
              </SecondaryInformationSection>
            </SecondaryInformation>
          )}
        </Title>
        <SoloTable
          columns={gatewayColumns}
          dataSource={gatewaysTableData}
          removePaging
          removeShadows
          curved={true}
        />
      </TableHolder>
    </FailoverDetailsContainer>
  );
};
