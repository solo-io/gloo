import React, { useState } from 'react';
import styled from '@emotion/styled/macro';
import { SoloInput } from 'Components/Common/SoloInput';
import { SelectValue } from 'antd/lib/select';
import { SoloDropdown } from 'Components/Common/SoloDropdown';
import { RouteTableStatus } from 'proto/github.com/solo-io/solo-apis/api/gloo/gateway/v1/route_table_pb';
import { RouteTablesPageCardContents } from '../VirtualService/RouteTablesTable';

const InputHolderRow = styled.div`
  display: grid;
  grid-gap: 15px;
  grid-template-columns: 1fr 1fr;
  margin-bottom: 25px;
`;

export const GlooInstanceRouteTables = () => {
  const [nameFilter, setNameFilter] = useState('');
  const [statusFilter, setStatusFilter] = useState<
  RouteTableStatus.StateMap[keyof RouteTableStatus.StateMap]
  >();

  const changeNameFilter = (e: React.ChangeEvent<any>) => {
    setNameFilter(e.target.value);
  };
  const changeStatusFilter = (newValue: SelectValue) => {
    setStatusFilter(
      newValue as RouteTableStatus.StateMap[keyof RouteTableStatus.StateMap]
    );
  };

  return (
    <>
      <InputHolderRow>
        <SoloInput
          placeholder={'Filter by name...'}
          onChange={changeNameFilter}
          value={nameFilter}
        />
        <SoloDropdown
          value={statusFilter}
          options={[
            {
              key: RouteTableStatus.State.ACCEPTED.toString(),
              value: RouteTableStatus.State.ACCEPTED,
              displayValue: 'Accepted',
            },
            {
              key: RouteTableStatus.State.REJECTED.toString(),
              value: RouteTableStatus.State.REJECTED,
              displayValue: 'Rejected',
            },
            {
              key: RouteTableStatus.State.PENDING.toString(),
              value: RouteTableStatus.State.PENDING,
              displayValue: 'Pending',
            },

            {
              key: RouteTableStatus.State.WARNING.toString(),
              value: RouteTableStatus.State.WARNING,
              displayValue: 'Warning',
            },
          ]}
          placeholder={'Filter by status...'}
          onChange={changeStatusFilter}
        />
      </InputHolderRow>

      <RouteTablesPageCardContents
        nameFilter={nameFilter}
        statusFilter={statusFilter}
        wholePage={false}
      />
    </>
  );
};
