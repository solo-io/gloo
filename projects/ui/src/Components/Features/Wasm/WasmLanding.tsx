import React, { useEffect, useState } from 'react';
import styled from '@emotion/styled/macro';
import { WasmPageTable } from './WasmTable';
import { SoloInput } from 'Components/Common/SoloInput';
import { SoloRadioGroup } from 'Components/Common/SoloRadioGroup';
import { VirtualServiceStatus } from 'proto/github.com/solo-io/solo-apis/api/gloo/gateway/v1/virtual_service_pb';
import { colors } from 'Styles/colors';
import { useListGateways, useListWasmFilters } from 'API/hooks';
import {
  ClusterObjectRef,
  ObjectRef,
} from 'proto/github.com/solo-io/skv2/api/core/v1/core_pb';
import { objectMetasAreEqual } from 'API/helpers';
import { ObjectMeta } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/common_pb';
import { useParams } from 'react-router';
import { DataError } from 'Components/Common/DataError';
import { Loading } from 'Components/Common/Loading';
import { SoloCheckbox } from 'Components/Common/SoloCheckbox';

const VirtualServiceLandingContainer = styled.div`
  display: grid;
  grid-template-columns: 200px 1fr;
  grid-gap: 28px;
`;

const HorizontalDivider = styled.div`
  position: relative;
  height: 1px;
  width: 100%;
  background: ${colors.marchGrey};
  margin: 35px 0;

  div {
    position: absolute;
    display: block;
    left: 0;
    right: 0;
    top: 50%;
    margin: -9px auto 0;
    width: 150px;
    text-align: center;
    color: ${colors.septemberGrey};
    background: ${colors.januaryGrey};
  }
`;

const FilterCheckboxHolder = styled.div`
  margin-bottom: 8px;

  &:last-child {
    margin-bottom: 0;
  }
`;

export const WasmLanding = () => {
  const [nameFilter, setNameFilter] = useState('');
  const [gatewayFiltersList, setGatewayFiltersList] = useState<
    { key: ClusterObjectRef.AsObject; filtering: boolean }[]
  >([]);
  const [instanceFiltersList, setInstanceFiltersList] = useState<
    { key: ObjectRef.AsObject; filtering: boolean }[]
  >([]);

  const { data: wasmFilters, error: wasmError } = useListWasmFilters();

  useEffect(() => {
    if (wasmFilters?.length) {
      let newGatewayFilters: {
        key: ClusterObjectRef.AsObject;
        filtering: boolean;
      }[] = [];
      let newInstanceFilters: {
        key: ObjectRef.AsObject;
        filtering: boolean;
      }[] = [];

      wasmFilters.forEach(filter => {
        filter.locationsList.forEach(location => {
          if (
            !newGatewayFilters.some(newGateway =>
              objectMetasAreEqual(newGateway.key, location.gatewayRef)
            )
          ) {
            newGatewayFilters.push({
              key: location.gatewayRef!,
              filtering: false,
            });
          }

          if (
            !newInstanceFilters.some(newInstance =>
              objectMetasAreEqual(newInstance.key, location.glooInstanceRef)
            )
          ) {
            newInstanceFilters.push({
              key: location.glooInstanceRef!,
              filtering: false,
            });
          }
        });
      });

      setGatewayFiltersList(newGatewayFilters);
      setInstanceFiltersList(newInstanceFilters);
    } else {
      setGatewayFiltersList([]);
      setInstanceFiltersList([]);
    }
  }, [wasmFilters]);

  if (!!wasmError) {
    return <DataError error={wasmError} />;
  } else if (!wasmFilters) {
    return <Loading message={'Retrieving Wasm Filters...'} />;
  }

  const changeNameFilter = (e: React.ChangeEvent<HTMLInputElement>) => {
    setNameFilter(e.target.value ?? '');
  };

  const toggleGatewayFilter = (
    newKey: ClusterObjectRef.AsObject,
    newVal: boolean
  ) => {
    const newGatewayFilters = gatewayFiltersList.map(gFilter => {
      if (objectMetasAreEqual(gFilter.key, newKey)) {
        return {
          key: newKey,
          filtering: newVal,
        };
      } else {
        return gFilter;
      }
    });

    setGatewayFiltersList(newGatewayFilters);
  };

  const toggleInstanceFilter = (
    newKey: ObjectRef.AsObject,
    newVal: boolean
  ) => {
    const newInstaceFilters = instanceFiltersList.map(iFilter => {
      if (objectMetasAreEqual(iFilter.key, newKey)) {
        return {
          key: newKey,
          filtering: newVal,
        };
      } else {
        return iFilter;
      }
    });

    setInstanceFiltersList(newInstaceFilters);
  };

  return (
    <VirtualServiceLandingContainer>
      <div>
        <SoloInput
          value={nameFilter}
          onChange={changeNameFilter}
          placeholder={'Filter by name...'}
        />

        {!!gatewayFiltersList.length && (
          <>
            <HorizontalDivider>
              <div>Gateway Filter</div>
            </HorizontalDivider>
            {gatewayFiltersList.map(gFilter => (
              <FilterCheckboxHolder
                key={
                  gFilter.key.name ??
                  '' + gFilter.key.namespace + gFilter.key.clusterName
                }>
                <SoloCheckbox
                  withWrapper={true}
                  checked={gFilter.filtering}
                  onChange={e =>
                    toggleGatewayFilter(gFilter.key, e.target.checked)
                  }
                  title={gFilter.key.name}
                />
              </FilterCheckboxHolder>
            ))}
          </>
        )}

        {!!instanceFiltersList.length && (
          <>
            <HorizontalDivider>
              <div>Gloo Instance Filter</div>
            </HorizontalDivider>
            {instanceFiltersList.map(iFilter => (
              <FilterCheckboxHolder
                key={iFilter.key.name ?? '' + iFilter.key.namespace}>
                <SoloCheckbox
                  withWrapper={true}
                  checked={iFilter.filtering}
                  onChange={e =>
                    toggleInstanceFilter(iFilter.key, e.target.checked)
                  }
                  title={iFilter.key.name}
                />
              </FilterCheckboxHolder>
            ))}
          </>
        )}
      </div>
      <WasmPageTable
        nameFilter={nameFilter}
        gatewayFilters={gatewayFiltersList
          .filter(gFilter => gFilter.filtering)
          .map(gFilter => gFilter.key)}
        instanceFilters={instanceFiltersList
          .filter(iFilter => iFilter.filtering)
          .map(iFilter => iFilter.key)}
      />
    </VirtualServiceLandingContainer>
  );
};
