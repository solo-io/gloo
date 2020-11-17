import React from 'react';
import { SectionCard } from 'Components/Common/SectionCard';
import { ReactComponent as WasmFilterIcon } from 'assets/wasm-filter-icon.svg';
import { SoloTable } from 'Components/Common/SoloTable';
import useSWR from 'swr';
import { gatewayAPI } from 'store/gateway/api';
import { Breadcrumb } from 'Components/Common/Breadcrumb';
import { WasmFilter } from 'proto/gloo/projects/gloo/api/v1/options/wasm/wasm_pb';
import { NavLink } from 'react-router-dom';
import css from '@emotion/css/macro';
import { colors } from 'Styles/colors';

export type WasmListingProps = {};

export const WasmListing: React.FC<WasmListingProps> = ({}) => {
  const { data: gatewaysList, error } = useSWR(
    'listGateways',
    gatewayAPI.listGateways
  );

  if (!gatewaysList && !error) {
    return <div>Loading...</div>;
  }

  console.log('gatewaysList', gatewaysList);

  function getWasmFilterData() {
    let wasmFiltersList = gatewaysList?.flatMap(gateway =>
      gateway?.gateway?.httpGateway?.options?.wasm?.filtersList?.map(filter => {
        return {
          ...filter,
          gateway: gateway.gateway?.metadata?.name
        };
      })
    )!;
    if (wasmFiltersList?.length > 0) {
      return wasmFiltersList.filter(filter => !!filter);
    }
    return [];
  }
  console.log('getWasmFilterData()', getWasmFilterData());
  return (
    <>
      <Breadcrumb />
      <SectionCard
        noPadding
        data-testid='upstreams-listing-section'
        cardName={'Wasm Filters'}
        logoIcon={<WasmFilterIcon />}>
        <div style={{ padding: '-20px' }}>
          <SoloTable
            pagination={{ hideOnSinglePage: true }}
            dataSource={getWasmFilterData()}
            columns={[
              {
                title: 'Filter Name',
                dataIndex: 'name'
              },
              {
                title: 'Gateway',
                dataIndex: 'gateway',
                render: (gateway: string) => (
                  <NavLink
                    css={css`
                      cursor: pointer;
                      color: ${colors.seaBlue};
                    `}
                    to={`/wasm/gateway/${gateway}`}>
                    {gateway}
                  </NavLink>
                )
              },
              {
                title: 'Root ID',
                dataIndex: 'rootId'
              },
              {
                title: 'Image',
                dataIndex: 'image'
              },
              {
                title: 'VM Type',
                dataIndex: 'vmType',
                render: (
                  vmType: WasmFilter.VmTypeMap[keyof WasmFilter.VmTypeMap]
                ) => (
                  <div className='flex items-center '>
                    {vmType === WasmFilter.VmType.V8 ? 'V8' : 'WAVM'}
                  </div>
                )
              }
            ]}
          />
        </div>
      </SectionCard>
    </>
  );
};
