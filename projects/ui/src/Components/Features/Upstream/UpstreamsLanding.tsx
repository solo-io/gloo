import React, { useState } from 'react';
import styled from '@emotion/styled/macro';
import { UpstreamsPageTable } from './UpstreamsTable';
import { UpstreamGroupsPageTable } from './UpstreamGroupsTable';
import { colors } from 'Styles/colors';
import { UpstreamStatus } from 'proto/github.com/solo-io/solo-apis/api/gloo/gloo/v1/upstream_pb';
import { SoloInput } from 'Components/Common/SoloInput';
import { SoloRadioGroup } from 'Components/Common/SoloRadioGroup';
import {
  SoloCheckbox,
} from 'Components/Common/SoloCheckbox';

const UpstreamLandingContainer = styled.div`
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
    width: 105px;
    text-align: center;
    color: ${colors.septemberGrey};
    background: ${colors.januaryGrey};
  }
`;

const UpstreamGroupToggleWrapper = styled.div`
  margin-top: 15px;
`;

export const UpstreamsLanding = () => {
  const [nameFilter, setNameFilter] = useState('');
  const [showUG, setShowUG] = useState(false);
  const [statusFilter, setStatusFilter] =
    useState<UpstreamStatus.StateMap[keyof UpstreamStatus.StateMap]>();

  const changeNameFilter = (e: React.ChangeEvent<HTMLInputElement>) => {
    setNameFilter(e.target.value);
  };

  const changeStatus = (idSelected: string | number | undefined) => {
    setStatusFilter(
      idSelected as UpstreamStatus.StateMap[keyof UpstreamStatus.StateMap]
    );
  };

  return (
    <UpstreamLandingContainer>
      <div>
        <SoloInput
          value={nameFilter}
          onChange={changeNameFilter}
          placeholder={'Filter by name...'}
        />

        <UpstreamGroupToggleWrapper>
          <SoloCheckbox
            title={'Show Upstream Groups'}
            checked={showUG}
            withWrapper={true}
            onChange={evt => setShowUG(evt.target.checked)}
          />
        </UpstreamGroupToggleWrapper>

        <HorizontalDivider>
          <div>Status Filter</div>
        </HorizontalDivider>
        <SoloRadioGroup
          options={[
            {
              displayName: 'Accepted',
              id: UpstreamStatus.State.ACCEPTED,
            },
            {
              displayName: 'Rejected',
              id: UpstreamStatus.State.REJECTED,
            },
            {
              displayName: 'Pending',
              id: UpstreamStatus.State.PENDING,
            },

            {
              displayName: 'Warning',
              id: UpstreamStatus.State.WARNING,
            },
          ]}
          currentSelection={statusFilter}
          onChange={changeStatus}
        />
      </div>
      <div>
        {showUG && (
          <UpstreamGroupsPageTable
            nameFilter={nameFilter}
            statusFilter={statusFilter}
          />
        )}
        <UpstreamsPageTable
          nameFilter={nameFilter}
          statusFilter={statusFilter}
        />
      </div>
    </UpstreamLandingContainer>
  );
};
