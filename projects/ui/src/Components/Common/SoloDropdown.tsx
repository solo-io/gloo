import React from 'react';
import { Select } from 'antd';
import styled from '@emotion/styled';
import { Label } from './SoloInput';
import { colors } from 'Styles/colors';
import { valueType } from 'antd/lib/statistic/utils';
import { SelectValue } from 'antd/lib/select';

// COMMON CHOICES
export const StatusDropdownChoices: OptionType[] = [
  {
    key: 'none',
    value: '',
    displayValue: 'Neither',
  },
  {
    key: 'accepted',
    value: 'Accepted',
  },
  {
    key: 'rejected',
    value: 'Rejected',
  },
];

export const SoloDropdownBlock = styled(Select)`
  position: relative;
  width: inherit;
  line-height: 16px;

  .ant-select-selector {
    padding: 7px 15px 7px 11px;
    border-radius: 8px;
    border: 1px solid ${colors.aprilGrey};
    line-height: 16px;
    outline: none;
    height: auto;
    cursor: pointer;

    .ant-select-selection-search,
    .ant-select-selection-search input {
      width: 0;
      padding: 0;
      border: none;
    }
    .ant-select-selection-placeholder {
      color: ${colors.juneGrey};
    }

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
  }

  .ant-select-arrow {
    position: absolute;
    right: 10px;
    top: 10px;
  }
`;

export interface OptionType {
  key: string;
  disabled?: boolean;
  value: string | number;
  displayValue?: React.ReactNode;
  icon?: JSX.Element;
}
export interface DropdownProps {
  value: string | number | undefined;
  options: OptionType[];
  onChange?: (newValue: SelectValue) => any;
  title?: string;
  placeholder?: string;
  defaultValue?: string | number;
  onBlur?: (evt: React.FocusEvent) => any;
  disabled?: boolean;
  testId?: string;
  error?: any;
}

export const SoloDropdown = (props: DropdownProps) => {
  const {
    title,
    disabled,
    defaultValue,
    value,
    placeholder,
    options,
    onChange,
    onBlur,
    testId,
    error,
    ...rest
  } = props;

  const getDefaultValue = (): string | number => {
    if (typeof defaultValue === undefined) {
      return '';
    }

    return defaultValue!;
  };

  return (
    <div style={{ width: '100%' }}>
      {title && <Label>{title}</Label>}
      <SoloDropdownBlock
        data-testid={testId}
        dropdownClassName={testId}
        value={value}
        dropdownMatchSelectWidth={false}
        defaultValue={getDefaultValue()}
        onChange={onChange}
        onBlur={onBlur}
        disabled={disabled}
        placeholder={placeholder}
        {...rest}>
        {options.map((opt: OptionType) => (
          <Select.Option
            key={opt.key}
            value={opt.value}
            disabled={opt.disabled}>
            {opt.icon} {opt.displayValue || opt.value}
          </Select.Option>
        ))}
      </SoloDropdownBlock>
    </div>
  );
};
