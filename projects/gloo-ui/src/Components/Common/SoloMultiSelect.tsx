/** @jsx jsx */
import { jsx } from '@emotion/core';
import { Select } from 'antd';
import { colors, soloConstants } from 'Styles';
import styled from '@emotion/styled/macro';
import { Label } from './SoloInput';

const SoloMultiselectBlock = styled(Select)`
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
      line-height: 33px;
      margin: 0;

      .ant-select-selection-selected-value {
        color: ${colors.septemberGrey};
      }

      .ant-select-selection__choice {
        background: ${colors.marchGrey};
        color: ${colors.novemberGrey};
        border-radius: ${soloConstants.smallRadius}px;
        line-height: 33px;
        height: 35px;

        .ant-select-selection__choice__content {
          margin-right: 8px;
        }

        .ant-select-selection__choice__remove {
          .anticon.anticon-close.ant-select-remove-icon {
            display: inline-flex;
            align-items: center;
            justify-content: center;
            height: 20px;
            width: 20px;
            border-radius: 20px;
            background: #9f9f9f;
            color: white;

            svg {
              height: 10px;
            }
          }
        }
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
export interface MultiselectProps {
  values: string[] | number[] | undefined[];
  options: OptionType[];
  onChange?: (newValues: string[] | number[]) => any;
  title?: string;
  placeholder?: string;
  defaultValues?: string[] | number[];
  onBlur?: (newValue: string | number) => any;
  disabled?: boolean;
}

export const SoloMultiSelect = (props: MultiselectProps) => {
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
      <SoloMultiselectBlock
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
      </SoloMultiselectBlock>
    </div>
  );
};
