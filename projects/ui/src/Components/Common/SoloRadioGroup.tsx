import * as React from 'react';
import styled from '@emotion/styled';
import { Checkbox } from 'antd';
import { colors } from 'Styles/colors';
import { Label } from './SoloInput';

type CheckboxWrapperProps = {
  checked?: boolean;
  withoutCheckboxVisual?: boolean;
  forceAChoice?: boolean;
};
const CheckboxWrapper = styled.div<CheckboxWrapperProps>`
  display: flex;
  justify-content: space-between;
  border-radius: 10px;
  width: 190px;
  padding: 7px 7px 7px 16px;
  color: ${colors.septemberGrey};
  background: white;
  border: 1px solid ${colors.juneGrey};
  transition: background 0.3s, border 0.3s;
  margin-bottom: 15px;
  cursor: pointer;

  label {
    cursor: pointer;
    pointer-events: none;
  }

  &:last-child {
    margin-bottom: 0;
  }

  ${(props: CheckboxWrapperProps) =>
    !!props.checked
      ? `background: ${colors.dropBlue};
        border-color: ${colors.seaBlue};
        color: ${colors.seaBlue};

        ${props.forceAChoice ? 'cursor: default;' : ''}
        `
      : ``}

  ${(props: CheckboxWrapperProps) =>
    !!props.withoutCheckboxVisual
      ? `
        label { display: none; }

        .ant-checkbox {
          .ant-checkbox-inner {
            visibility: hidden; 
          }
        }`
      : `
      .ant-checkbox {
            
        input {
          display: none;
        }
        .ant-checkbox-inner {
          display: block;
          background: ${colors.januaryGrey};
          border: 1px solid ${colors.juneGrey};
          border-radius: 11px;
          width: 18px;
          height: 18px;
        }


        &.ant-checkbox-checked {
          .ant-checkbox-inner {
            background: ${colors.puddleBlue};
            border-color: ${colors.seaBlue};
            border-radius: 11px;

            &:after {
              display: block;
              transform: none;
              border: none;
              background: ${colors.seaBlue};
              width: 4px;
              height: 4px;
              border-radius: 4px;
              border-color: ${colors.seaBlue};
              border-width: 1px;
              left: 5px;
              top: 9px;

              transform: rotate(45deg) scale(1) translate(-37%, -66%);
            }
          }
        }
      }`}
`;

interface Props {
  options: {
    displayName: string;
    id: string | number;
  }[];
  currentSelection: string | number | undefined; //matches to id
  onChange: (idSelected: string | number | undefined) => any;
  withoutCheckboxes?: boolean;
  forceAChoice?: boolean;
  title?: string;
  className?: string;
}

export const SoloRadioGroup = (props: Props) => {
  const {
    options,
    currentSelection,
    onChange,
    withoutCheckboxes,
    forceAChoice,
    title,
    className,
  } = props;

  const attemptSelection = (selectedId: string | number) => {
    if (selectedId !== currentSelection) {
      onChange(selectedId);
    } else if (!forceAChoice) {
      onChange(undefined);
    }
  };

  return (
    <div>
      {title && <Label>{title}</Label>}

      <div className={className}>
        {options.map(option => {
          return (
            <CheckboxWrapper
              key={option.id}
              checked={option.id === currentSelection}
              onClick={() => attemptSelection(option.id)}
              withoutCheckboxVisual={withoutCheckboxes}
              forceAChoice={forceAChoice}>
              {option.displayName}
              {/** The checkbox below is only for the visual. The wrapper is the intended clickable */}
              <Checkbox
                checked={option.id === currentSelection}
                onChange={() => {}}
              />
            </CheckboxWrapper>
          );
        })}
      </div>
    </div>
  );
};
