import styled from '@emotion/styled/macro';
import {
  CheckboxFilterProps,
  SoloCheckbox,
} from 'Components/Common/SoloCheckbox';
import { SoloInput } from 'Components/Common/SoloInput';
import React, { useState } from 'react';
import { colors } from 'Styles/colors';
import { GraphqlPageTable } from './GraphqlTable';

const GraphqlLandingContainer = styled.div`
  display: grid;
  grid-template-columns: 200px 1fr;
  grid-gap: 28px;
`;

const GraphQLIconHolder = styled.div`
  display: flex;
  align-items: center;
  justify-items: center;

  svg {
    width: 35px;
    max-width: none;
  }
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
    margin-bottom: 8px;
  }
`;

export const GraphqlLanding = () => {
  const [nameFilter, setNameFilter] = useState('');

  const [typeFilters, setTypeFilters] = useState<CheckboxFilterProps[]>([
    { label: 'GraphQL', checked: false },
    { label: 'REST', checked: false },
    { label: 'gRPC', checked: false },
  ]);

  const changeNameFilter = (e: React.ChangeEvent<HTMLInputElement>) => {
    setNameFilter(e.target.value);
  };

  const changeTypeFilter = (ind: number, checked: boolean) => {
    const newArray = [...typeFilters];
    newArray[ind].checked = checked;
    setTypeFilters(newArray);
  };
  return (
    <GraphqlLandingContainer>
      <div>
        <SoloInput
          value={nameFilter}
          onChange={changeNameFilter}
          placeholder={'Filter by name...'}
        />

        <HorizontalDivider>
          <div>Status Filter</div>
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
        <div>Table</div>
        <GraphqlPageTable />
      </div>
    </GraphqlLandingContainer>
  );
};
