import React, { useState } from 'react';
import styled from '@emotion/styled/macro';
import { UpstreamsPageTable } from './UpstreamsTable';
import { UpstreamGroupsPageTable } from './UpstreamGroupsTable';
import { colors } from 'Styles/colors';
import { UpstreamStatus } from 'proto/github.com/solo-io/solo-apis/api/gloo/gloo/v1/upstream_pb';
import { SoloInput } from 'Components/Common/SoloInput';
import { SoloRadioGroup } from 'Components/Common/SoloRadioGroup';
import {
  CheckboxFilterProps,
  SoloCheckbox,
} from 'Components/Common/SoloCheckbox';
import {
  TYPE_AWS,
  TYPE_AWS_EC2,
  TYPE_AZURE,
  TYPE_CONSUL,
  TYPE_KUBE,
  TYPE_STATIC,
} from 'utils/upstream-helpers';

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

const CheckboxWrapper = styled.div`
  > div {
    width: 190px;
    margin-bottom: 15px;
  }
`;

const UpstreamGroupToggleWrapper = styled.div`
  margin-top: 15px;
`;

const CHECKBOX_DEFAULT_FILTERS: CheckboxFilterProps[] = [
  TYPE_AWS,
  TYPE_AZURE,
  TYPE_CONSUL,
  TYPE_KUBE,
  TYPE_AWS_EC2,
  TYPE_STATIC,
].map(s => ({ label: s, checked: false }));

export const UpstreamsLanding = () => {
  const [nameFilter, setNameFilter] = useState('');
  const [showUG, setShowUG] = useState(false);
  const [statusFilter, setStatusFilter] =
    useState<UpstreamStatus.StateMap[keyof UpstreamStatus.StateMap]>();
  const [typeFilters, setTypeFilters] = useState<CheckboxFilterProps[]>(
    CHECKBOX_DEFAULT_FILTERS
  );

  const changeNameFilter = (e: React.ChangeEvent<HTMLInputElement>) => {
    setNameFilter(e.target.value);
  };

  const changeStatus = (idSelected: string | number | undefined) => {
    setStatusFilter(
      idSelected as UpstreamStatus.StateMap[keyof UpstreamStatus.StateMap]
    );
  };

  const changeTypeFilter = (ind: number, checked: boolean) => {
    const newArray = [...typeFilters];
    newArray[ind].checked = checked;
    setTypeFilters(newArray);
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

        <HorizontalDivider>
          <div>Types Filter</div>
        </HorizontalDivider>
        <CheckboxWrapper>
          {typeFilters.map((filter, ind) => {
            return (
              <SoloCheckbox
                key={filter.label}
                title={filter.label}
                checked={filter.checked}
                withWrapper={true}
                onChange={evt => changeTypeFilter(ind, evt.target.checked)}
              />
            );
          })}
        </CheckboxWrapper>
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
          typeFilters={typeFilters}
        />
      </div>
    </UpstreamLandingContainer>
  );
};
