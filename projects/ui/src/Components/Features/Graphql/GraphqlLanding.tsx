import styled from '@emotion/styled/macro';
import { ReactComponent as GreenPlus } from 'assets/small-green-plus.svg';
import {
  CheckboxFilterProps,
  SoloCheckbox,
} from 'Components/Common/SoloCheckbox';
import { SoloInput } from 'Components/Common/SoloInput';
import React, { useState } from 'react';
import { colors } from 'Styles/colors';
import { GraphqlPageTable } from './GraphqlTable';
import { NewApiModal } from './NewApiModal';

export enum APIType {
  REST = 'REST',
  GRPC = 'gRPC',
  GRAPHQL = 'GraphQL',
}
const GraphqlLandingContainer = styled.div`
  display: grid;
  grid-template-columns: 200px 1fr;
  grid-gap: 28px;
`;

const SecondaryComponent = styled.div`
  position: absolute;
  right: 20px;
`;

const Button = styled.button`
  color: ${colors.oceanBlue};
  &:hover {
    cursor: pointer;
    color: ${colors.seaBlue};
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

const API_TYPES: CheckboxFilterProps[] = [
  {
    checked: true,
    label: APIType.GRAPHQL,
  },
  {
    checked: false,
    label: APIType.REST,
  },
  {
    checked: false,
    label: APIType.GRPC,
  },
];

export const GraphqlLanding = () => {
  const [nameFilter, setNameFilter] = useState('');
  const [showGraphqlModal, setShowGraphqlModal] = React.useState(false);

  const openModal = () => setShowGraphqlModal(true);

  const [typeFilters, setTypeFilters] =
    useState<CheckboxFilterProps[]>(API_TYPES);
  const changeNameFilter = (e: React.ChangeEvent<HTMLInputElement>) => {
    setNameFilter(e.target.value);
  };

  const toggleGraphqlModal = () => {
    setShowGraphqlModal(!showGraphqlModal);
  };

  const changeTypeFilter = (filter: CheckboxFilterProps, checked: boolean) => {
    setTypeFilters(
      typeFilters.map(f => {
        if (f.label === filter.label) {
          return {
            ...f,
            checked,
          };
        } else return f;
      })
    );
  };
  return (
    <>
      <GraphqlLandingContainer className='relative'>
        <span
          onClick={openModal}
          className='absolute right-0 flex items-center text-green-400 cursor-pointer -top-8 hover:text-green-300'>
          <GreenPlus className='w-6 mr-1 fill-current' />
          <span className='text-gray-700'> Create API</span>
        </span>
        <div>
          <SoloInput
            value={nameFilter}
            onChange={changeNameFilter}
            placeholder={'Filter by name...'}
          />
          <HorizontalDivider>
            <div>API Type</div>
          </HorizontalDivider>
          <CheckboxWrapper>
            {typeFilters.map((filter, ind) => {
              return (
                <SoloCheckbox
                  disabled={true}
                  key={filter.label}
                  title={filter.label}
                  checked={filter.checked}
                  withWrapper={true}
                  onChange={evt => changeTypeFilter(filter, evt.target.checked)}
                />
              );
            })}
          </CheckboxWrapper>
        </div>
        <div>
          <GraphqlPageTable typeFilters={typeFilters} />
        </div>
      </GraphqlLandingContainer>

      <NewApiModal
        showNewModal={showGraphqlModal}
        toggleNewModal={toggleGraphqlModal}
      />
    </>
  );
};
