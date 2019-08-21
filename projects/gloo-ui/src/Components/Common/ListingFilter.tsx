import * as React from 'react';
/** @jsx jsx */
import { jsx } from '@emotion/core';
import styled from '@emotion/styled/macro';
import { colors } from 'Styles';
import { SoloInput } from './SoloInput';
import { SoloCheckbox } from './SoloCheckbox';
import { SoloRadioGroup } from './SoloRadioGroup';

const FilterContainer = styled.div`
  display: flex;
`;
const Filters = styled.div`
  width: 190px;
  margin-right: 35px;
`;
const Content = styled.div`
  flex: 1;
`;

const StyledHeader = styled.div`
  overflow: hidden;
  text-align: center;
  padding: 15px;
  color: ${colors.septemberGrey};
  &:before,
  &:after {
    background-color: ${colors.aprilGrey};
    content: '';
    display: inline-block;
    height: 1px;
    position: relative;
    vertical-align: middle;
    width: 50%;
  }
  &:before {
    right: 0.5em;
    margin-left: -50%;
  }
  &:after {
    left: 0.5em;
    margin-right: -50%;
  }
`;

const FilterInput = styled.div`
  margin-bottom: 15px;
`;

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
  id: string;
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
  onChange?: (
    strings: StringFilterProps[],
    types: TypeFilterProps[],
    checkboxes: CheckboxFilterProps[],
    radios: RadioFilterProps[]
  ) => any;
  hideFilters?: boolean;
  showLabels?: boolean;
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

  React.useEffect(() => {
    if (filterProps.onChange) {
      filterProps.onChange(
        stringFilters,
        typesFilters,
        checkboxFilters,
        radioFilters
      );
    }
  }, [stringFilters, typesFilters, checkboxFilters, radioFilters]);

  return (
    <FilterContainer>
      {!filterProps.hideFilters && (
        <Filters>
          <FilterInput>
            {stringFilters.map((filter, ind) => {
              return (
                <SoloInput
                  key={filter.displayName}
                  value={filter.value!}
                  placeholder={filter.placeholder}
                  onChange={({ target }) => {
                    const newArray = [...stringFilters];
                    newArray[ind].value = target.value;

                    setStringFilters(newArray);
                  }}
                />
              );
            })}
          </FilterInput>
          <React.Fragment>
            {radioFilters.length > 0 && filterProps.showLabels && (
              <StyledHeader>Status Filter</StyledHeader>
            )}
          </React.Fragment>
          {radioFilters.map((filter, ind) => {
            return (
              <SoloRadioGroup
                key={ind}
                options={filter.options.map(option => {
                  return {
                    displayName: option.displayName,
                    id: option.id || option.displayName
                  };
                })}
                currentSelection={filter.choice}
                onChange={newValue => {
                  const newArray = [...radioFilters];
                  newArray[ind].choice = newValue;

                  setRadioFilters(newArray);
                }}
              />
            );
          })}
          {typesFilters.map((filter, ind) => {
            return (
              <SoloRadioGroup
                key={ind}
                options={filter.options.map(option => {
                  return {
                    displayName: option.displayName,
                    id: option.id || option.displayName
                  };
                })}
                currentSelection={filter.choice}
                withoutCheckboxes={true}
                forceAChoice={true}
                onChange={newValue => {
                  const newArray = [...typesFilters];
                  newArray[ind].choice = newValue;

                  setTypesFilters(newArray);
                }}
              />
            );
          })}
          {checkboxFilters.length > 0 && filterProps.showLabels && (
            <StyledHeader>Types Filter</StyledHeader>
          )}
          {checkboxFilters.map((filter, ind) => {
            return (
              <SoloCheckbox
                key={filter.displayName}
                title={filter.displayName}
                checked={filter.value!}
                withWrapper={true}
                onChange={evt => {
                  const newArray = [...checkboxFilters];
                  newArray[ind].value = evt.target.checked;

                  setCheckboxFilters(newArray);
                }}
              />
            );
          })}
        </Filters>
      )}
      <Content>
        {filterProps.filterFunction(
          stringFilters,
          typesFilters,
          checkboxFilters,
          radioFilters
        )}
      </Content>
    </FilterContainer>
  );
};
