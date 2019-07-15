/** @jsx jsx */
import { jsx } from '@emotion/core';
import { Select } from 'antd';
import { colors, soloConstants } from 'Styles';
import styled from '@emotion/styled/macro';
import { Label } from './SoloInput';

const SoloDropdownBlock = styled(Select)`
  width: 100%;
  /* margin-bottom: 15px; */
  line-height: 16px;

  .ant-select-selection {
    width: 100%;
    padding: 9px 15px 9px 11px;
    border: 1px solid ${colors.aprilGrey};
    border-radius: ${soloConstants.smallRadius}px;
    height: auto;
    outline: none;

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
`;

export interface OptionType {
  key: string;
  disabled?: boolean;
  value: string | number;
}
export interface DropdownProps {
  values: string[] | number[] | undefined[];
  options: OptionType[];
  onChange?: (newValues: string[] | number[]) => any;
  title?: string;
  placeholder?: string;
  defaultValues?: string[] | number[];
  onBlur?: (newValue: string | number) => any;
  disabled?: boolean;
}

export const SoloMultiSelect = (props: DropdownProps) => {
  const {
    title,
    // disabled,
    defaultValues,
    values,
    placeholder,
    options,
    onChange,
    onBlur
  } = props;

  return (
    <div style={{ width: '100%' }}>
      {title && <Label>{title}</Label>}
      <SoloDropdownBlock
        mode='multiple'
        value={values}
        defaultValue={defaultValues || [] /**
      //@ts-ignore */}
        onChange={onChange /**
        //@ts-ignore */}
        onBlur={onBlur}
        placeholder={placeholder}>
        {options.map((opt: OptionType) => (
          <Select.Option
            key={opt.key}
            value={opt.value}
            disabled={opt.disabled}>
            {opt.value}
          </Select.Option>
        ))}
      </SoloDropdownBlock>
    </div>
  );
};
