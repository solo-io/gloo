import React, { useEffect } from 'react';
import styled from '@emotion/styled/macro';
import {
  SoloTable,
  RenderStatus,
  TableActionCircle,
  RenderCluster,
  RenderExpandableNamesList,
} from 'Components/Common/SoloTable';
import { SectionCard } from 'Components/Common/SectionCard';
import { ReactComponent as FilterIcon } from 'assets/filter-icon.svg';
import { ReactComponent as DownloadIcon } from 'assets/download-icon.svg';
import { colors } from 'Styles/colors';
import { useParams, useNavigate } from 'react-router';
import {
  useGetWasmFilter,
  useListGateways,
  useListWasmFilters,
} from 'API/hooks';
import { Loading } from 'Components/Common/Loading';
import { objectMetasAreEqual } from 'API/helpers';
import { SimpleLinkProps, RenderSimpleLink } from 'Components/Common/SoloLink';
import { VirtualServiceStatus } from 'proto/github.com/solo-io/solo-apis/api/gloo/gateway/v1/virtual_service_pb';
import { gatewayResourceApi } from 'API/gateway-resources';
import { doDownload } from 'download-helper';
import { DataError } from 'Components/Common/DataError';
import {
  ClusterObjectRef,
  ObjectRef,
} from 'proto/github.com/solo-io/skv2/api/core/v1/core_pb';
import { ObjectMeta } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/common_pb';
import { WasmFilter } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/wasm_pb';
import { SoloModal } from 'Components/Common/SoloModal';
import { WasmFilterDetails } from './WasmFilterDetails';

const FilterIconHolder = styled.div`
  display: flex;
  align-items: center;
  justify-items: center;

  svg {
    width: 21px;
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

type WasmFilterTableFields = {
  key: string;
  name: SimpleLinkProps;
  gateways: string[];
  glooInstance: string[];
};

type Props = {
  nameFilter: string;
  instanceFilters: ObjectRef.AsObject[];
  gatewayFilters: ClusterObjectRef.AsObject[];
};
export const WasmPageTable = ({
  nameFilter,
  instanceFilters,
  gatewayFilters,
}: Props) => {
  const navigate = useNavigate();
  const { filterName } = useParams();

  const [tableData, setTableData] = React.useState<WasmFilterTableFields[]>([]);
  const [
    filterOfInterest,
    setFilterOfInterest,
  ] = React.useState<WasmFilter.AsObject>();

  const { data: wasmFilters, error: wasmFiltersError } = useListWasmFilters();

  useEffect(() => {
    if (wasmFilters) {
      let newTableData: WasmFilterTableFields[] = [];

      wasmFilters
        .filter(
          wasmFilter =>
            (!gatewayFilters.length ||
              gatewayFilters.some(gFilter => {
                return wasmFilter.locationsList.some(location => {
                  return objectMetasAreEqual(location.gatewayRef, gFilter);
                });
              })) &&
            (!instanceFilters.length ||
              instanceFilters.some(iFilter => {
                return wasmFilter.locationsList.some(location => {
                  return objectMetasAreEqual(location.glooInstanceRef, iFilter);
                });
              })) &&
            wasmFilter.name.includes(nameFilter)
        )
        .sort((gA, gB) => gA.name.localeCompare(gB.name) || (gA.locationsList[0]?.glooInstanceRef?.name ?? '').localeCompare(gB.locationsList[0]?.glooInstanceRef?.name ?? ''))
        .forEach(filter => {
          newTableData.push({
            key: filter.name,
            name: {
              displayElement: (
                <div style={{ display: 'flex' }}>
                  <FilterIconHolder>
                    <FilterIcon />
                  </FilterIconHolder>
                  <span style={{ marginLeft: '8px' }}>{filter.name}</span>
                </div>
              ),
              link: `${filter.name}/`,
            },
            gateways: filter.locationsList.map(
              loc => loc.gatewayRef?.name ?? ''
            ),
            glooInstance: filter.locationsList.map(
              loc => loc.glooInstanceRef?.name ?? ''
            ),
          });
        });

      setTableData(newTableData);
    } else {
      setTableData([]);
    }
  }, [wasmFilters, nameFilter, gatewayFilters, instanceFilters]);

  useEffect(() => {
    if (!!filterName) {
      const foundFilter = !!wasmFilters
        ? wasmFilters.find(filter => filter.name === filterName)
        : undefined;

      setFilterOfInterest(foundFilter);
    } else {
      setFilterOfInterest(undefined);
    }
  }, [filterName, wasmFilters]);

  if (!!filterName && !filterOfInterest) {
    return (
      <DataError error={{ message: `No filter named ${filterName} found.` }} />
    );
  }

  if (!!wasmFiltersError) {
    return <DataError error={wasmFiltersError} />;
  } else if (!wasmFilters) {
    return <Loading message={'Retrieving Wasm Filters...'} />;
  }

  let columns: any = [
    {
      title: 'Name',
      dataIndex: 'name',
      width: 200,
      render: RenderSimpleLink,
    },
    {
      title: 'Gateways',
      dataIndex: 'gateways',
      render: RenderExpandableNamesList,
    },
    {
      title: 'Gloo Instance',
      dataIndex: 'glooInstance',
      render: RenderExpandableNamesList,
    },
  ];

  return (
    <>
      <SectionCard
        cardName={'Wasm Filters'}
        logoIcon={
          <FilterIconHolder>
            <FilterIcon />
          </FilterIconHolder>
        }
        noPadding={true}>
        <TableHolder wholePage={true}>
          <SoloTable
            columns={columns}
            dataSource={tableData}
            removePaging
            removeShadows
            curved={true}
          />
        </TableHolder>
      </SectionCard>

      {!!filterOfInterest && (
        <SoloModal
          visible={true}
          width={1040}
          onClose={() => navigate('/wasm-filters/')}>
          <WasmFilterDetails wasmFilter={filterOfInterest} />
        </SoloModal>
      )}
    </>
  );
};
