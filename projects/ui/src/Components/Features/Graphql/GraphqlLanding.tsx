import { ReactComponent as GreenPlus } from 'assets/small-green-plus.svg';
import {
  CheckboxFilterProps,
  SoloCheckbox,
} from 'Components/Common/SoloCheckbox';
import { SoloInput } from 'Components/Common/SoloInput';
import React, { useState } from 'react';
import * as styles from './GraphqlLanding.style';
import { GraphqlPageTable } from './GraphqlTable';
import { NewApiModal } from './new-api-modal/NewApiModal';

export enum APIType {
  REST = 'REST',
  GRPC = 'gRPC',
  GRAPHQL = 'GraphQL',
}

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
    <styles.GraphqlLandingContainer className='relative'>
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
        <styles.HorizontalDivider>
          <div>API Type</div>
        </styles.HorizontalDivider>
        <styles.CheckboxWrapper>
          {typeFilters.map((filter, ind) => (
            <SoloCheckbox
              disabled={true}
              key={filter.label}
              title={filter.label}
              checked={filter.checked}
              withWrapper={true}
              onChange={evt => changeTypeFilter(filter, evt.target.checked)}
            />
          ))}
        </styles.CheckboxWrapper>
      </div>
      <div>
        <GraphqlPageTable typeFilters={typeFilters} nameFilter={nameFilter} />
      </div>
      <NewApiModal
        show={showGraphqlModal}
        onClose={() => setShowGraphqlModal(false)}
      />
    </styles.GraphqlLandingContainer>
  );
};
