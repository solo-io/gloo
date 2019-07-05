import * as React from 'react';
/** @jsx jsx */
import { jsx } from '@emotion/core';
import { AutoComplete } from 'antd';
import { colors, soloConstants } from 'Styles';
import styled from '@emotion/styled/macro';
import { Label } from './SoloInput';
import { SelectValue } from 'antd/lib/select';
import { DataSourceItemType } from 'antd/lib/auto-complete';

const { Option } = AutoComplete;

const SoloAutocompleteBlock = styled(AutoComplete)`
  width: 100%;
  /* margin-bottom: 15px; */
  line-height: 40px;

  &.ant-select {
    .ant-select-search--inline {
      float: none;
    }
    .ant-select-selection--single {
      height: 36px;
    }
  }

  .ant-select-selection {
    width: 100%;
    height: auto;

    .ant-select-selection__rendered {
      line-height: inherit;
      margin: 0;

      .ant-select-selection-selected-value {
        color: ${colors.septemberGrey};
      }
    }

    &:disabled {
      background: ${colors.aprilGrey};
    }

    input.ant-input {
      height: auto;
      line-height: 16px;
      padding: 9px 15px 9px 11px;
      border: 1px solid ${colors.aprilGrey} !important;
      border-radius: ${soloConstants.smallRadius}px;
      outline: none;
      color: ${colors.septemberGrey};

      &:focus {
        outline: none;
      }
    }

    .ant-select-arrow {
      display: block;
    }
  }
`;

interface DropdownProps {
  presetOptions: string[];
  onChange: (newValue: string) => any;
  title?: string;
  placeholder?: string;
  defaultValue?: string | number;
  onBlur?: (newValue: string | number) => any;
  disabled?: boolean;
}

export const SoloTypeahead = (props: DropdownProps) => {
  const [typeInText, setTypeInText] = React.useState<string>('');

  const { title, disabled, placeholder, presetOptions, onChange } = props;

  const handleChange = (value: SelectValue): void => {
    onChange(value as string);
  };
  const getOptions = (): DataSourceItemType[] => {
    return presetOptions
      .filter(
        opt =>
          opt.toLowerCase().includes(typeInText.toLowerCase()) &&
          opt.toLowerCase() !== typeInText.toLowerCase()
      )
      .concat(typeInText.length ? [typeInText] : [])
      .map(renderOption);
  };

  const renderOption = (value: string) => {
    return (
      // @ts-ignore
      <Option key={value} value={value}>
        {value}
      </Option>
    );
  };

  return (
    <React.Fragment>
      {title && <Label>{title}</Label>}
      <SoloAutocompleteBlock
        disabled={disabled}
        onChange={handleChange}
        defaultValue={
          props.presetOptions.length
            ? props.presetOptions[0]
            : 'New namespace...'
        }
        onSearch={setTypeInText}
        dataSource={getOptions()}
        placeholder={placeholder}
        optionLabelProp='value'
      />
    </React.Fragment>
  );
};
