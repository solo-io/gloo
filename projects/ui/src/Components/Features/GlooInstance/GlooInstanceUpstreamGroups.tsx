import React, { useState } from 'react';
import styled from '@emotion/styled/macro';
import { SoloInput } from 'Components/Common/SoloInput';
import { SelectValue } from 'antd/lib/select';
import { UpstreamGroupsTable } from '../Upstream/UpstreamGroupsTable';
import { SoloDropdown } from 'Components/Common/SoloDropdown';
import { UpstreamStatus } from 'proto/github.com/solo-io/solo-apis/api/gloo/gloo/v1/upstream_pb';

const InputHolderRow = styled.div`
  display: grid;
  grid-gap: 15px;
  grid-template-columns: 1fr 1fr;
  margin-bottom: 25px;
`;

export const GlooInstanceUpstreamGroups = () => {
  const [nameFilter, setNameFilter] = useState('');
  const [statusFilter, setStatusFilter] =
    useState<UpstreamStatus.StateMap[keyof UpstreamStatus.StateMap]>();

  const changeNameFilter = (e: React.ChangeEvent<any>) => {
    setNameFilter(e.target.value);
  };
  const changeStatusFilter = (newValue: SelectValue) => {
    setStatusFilter(
      newValue as UpstreamStatus.StateMap[keyof UpstreamStatus.StateMap]
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
              key: UpstreamStatus.State.ACCEPTED.toString(),
              value: UpstreamStatus.State.ACCEPTED,
              displayValue: 'Accepted',
            },
            {
              key: UpstreamStatus.State.REJECTED.toString(),
              value: UpstreamStatus.State.REJECTED,
              displayValue: 'Rejected',
            },
            {
              key: UpstreamStatus.State.PENDING.toString(),
              value: UpstreamStatus.State.PENDING,
              displayValue: 'Pending',
            },

            {
              key: UpstreamStatus.State.WARNING.toString(),
              value: UpstreamStatus.State.WARNING,
              displayValue: 'Warning',
            },
          ]}
          placeholder={'Filter by status...'}
          onChange={changeStatusFilter}
        />
      </InputHolderRow>

      <UpstreamGroupsTable
        nameFilter={nameFilter}
        statusFilter={statusFilter}
        wholePage={false}
      />
    </>
  );
};
