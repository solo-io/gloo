import * as React from 'react';
/** @jsx jsx */
import { jsx } from '@emotion/core';
import styled from '@emotion/styled/macro';
import { colors } from 'Styles';

export interface StringFilterProps {
  displayName: string;
  placeholder?: string;
  value?: string;
}

export interface CheckboxFilterProps {
  displayName: string;
  value?: boolean;
}

export interface RadioFilterProps {
  options: {
    id?: string;
    displayName: string;
  }[];
  choice?: string; //matched to id
}

export interface TypeFilterProps {
  options: {
    id?: string;
    displayName: string;
  }[];
  choice?: string; //matched to id
}

interface FilterProps {
  strings?: StringFilterProps[];
  types?: TypeFilterProps[];
  checkboxes?: CheckboxFilterProps[];
  radios?: RadioFilterProps[];
  filterFunction: (
    strings: StringFilterProps[],
    types: TypeFilterProps[],
    checkboxes: CheckboxFilterProps[],
    radios: RadioFilterProps[]
  ) => any;
}

export const ListingFilter = (filterProps: FilterProps) => {
  const [stringFilters, setStringFilters] = React.useState<StringFilterProps[]>(
    filterProps.strings
      ? filterProps.strings.map(stringFilter => {
          return {
            displayName: stringFilter.displayName,
            placeholder: stringFilter.placeholder,
            value: stringFilter.value || ''
          };
        })
      : []
  );
  const [typesFilters, setTypesFilters] = React.useState<TypeFilterProps[]>(
    filterProps.types
      ? filterProps.types.map(typeFilter => {
          return {
            ...typeFilter
          };
        })
      : []
  );
  const [checkboxFilters, setCheckboxFilters] = React.useState<
    CheckboxFilterProps[]
  >(
    filterProps.checkboxes
      ? filterProps.checkboxes.map(checkboxFilter => {
          return {
            ...checkboxFilter
          };
        })
      : []
  );
  const [radioFilters, setRadioFilters] = React.useState<RadioFilterProps[]>(
    filterProps.radios
      ? filterProps.radios.map(radioFilter => {
          return {
            ...radioFilter
          };
        })
      : []
  );

  return (
    <div>
      <div>Filters...</div>
      <div>
        List Container...
        {filterProps.filterFunction(
          stringFilters,
          typesFilters,
          checkboxFilters,
          radioFilters
        )}
      </div>
    </div>
  );
};
