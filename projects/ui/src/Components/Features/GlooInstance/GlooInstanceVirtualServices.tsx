import React, { useState } from 'react';
import styled from '@emotion/styled/macro';
import { SoloInput } from 'Components/Common/SoloInput';
import { VirtualServicesTable } from '../VirtualService/VirtualServicesTable';
import { SelectValue } from 'antd/lib/select';
import { SoloDropdown } from 'Components/Common/SoloDropdown';
import { VirtualServiceStatus } from 'proto/github.com/solo-io/solo-apis/api/gloo/gateway/v1/virtual_service_pb';

const InputHolderRow = styled.div`
  display: grid;
  grid-gap: 15px;
  grid-template-columns: 1fr 1fr;
  margin-bottom: 25px;
`;

export const GlooInstanceVirtualServices = () => {
  const [nameFilter, setNameFilter] = useState('');
  const [statusFilter, setStatusFilter] =
    useState<
      VirtualServiceStatus.StateMap[keyof VirtualServiceStatus.StateMap]
    >();

  const changeNameFilter = (e: React.ChangeEvent<any>) => {
    setNameFilter(e.target.value);
  };
  const changeStatusFilter = (newValue: SelectValue) => {
    setStatusFilter(
      newValue as VirtualServiceStatus.StateMap[keyof VirtualServiceStatus.StateMap]
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
              key: VirtualServiceStatus.State.ACCEPTED.toString(),
              value: VirtualServiceStatus.State.ACCEPTED,
              displayValue: 'Accepted',
            },
            {
              key: VirtualServiceStatus.State.REJECTED.toString(),
              value: VirtualServiceStatus.State.REJECTED,
              displayValue: 'Rejected',
            },
            {
              key: VirtualServiceStatus.State.PENDING.toString(),
              value: VirtualServiceStatus.State.PENDING,
              displayValue: 'Pending',
            },

            {
              key: VirtualServiceStatus.State.WARNING.toString(),
              value: VirtualServiceStatus.State.WARNING,
              displayValue: 'Warning',
            },
          ]}
          placeholder={'Filter by status...'}
          onChange={changeStatusFilter}
        />
      </InputHolderRow>

      <VirtualServicesTable
        nameFilter={nameFilter}
        statusFilter={statusFilter}
        wholePage={false}
      />
    </>
  );
};
