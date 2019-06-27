import * as React from 'react';
/** @jsx jsx */
import { css, jsx } from '@emotion/core';
import { colors, soloConstants } from '../../Styles';
import styled from '@emotion/styled/macro';
import { Checkbox } from 'antd';

const InputStyling = css`
  border-radius: 10px;
  width: 190px;
  padding: 10px 16px;
  margin-bottom: 20px;
  background: white;
  border: 1px solid ${colors.juneGrey};
`;

const CheckboxStyling = css`
  .ant-checkbox {
    .ant-checkbox-inner {
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
          left: 6px;
          top: 6px;

          transform: rotate(45deg) scale(1) translate(-37%, -66%);
        }
      }
    }
  }
`;

const CheckboxWrapper = styled<
  'div',
  { checked?: boolean; withoutCheckboxVisual?: boolean }
>('div')`
  ${InputStyling}
  display: flex;
  justify-content: space-between;
  padding: 7px 7px 7px 16px;
  color: ${colors.septemberGrey};
  transition: background ${soloConstants.transitionTime},
    border ${soloConstants.transitionTime};

  ${props =>
    !!props.checked
      ? `background: ${colors.dropBlue};
        border-color: ${colors.seaBlue};
        cursor: default;`
      : `cursor: pointer;`}

  ${props =>
    !!props.withoutCheckboxVisual
      ? `.ant-checkbox {
          .ant-checkbox-inner {
            visibility: hidden; 
          }
        }`
      : CheckboxStyling}
`;

interface Props {
  options: {
    displayName: string;
    id: string;
  }[];
  currentSelection: string | undefined; //matches to id
  onChange: (idSelected: string | undefined) => any;
  withoutCheckboxes?: boolean;
  forceAChoice?: boolean;
}

export const SoloRadioGroup = (props: Props) => {
  const {
    options,
    currentSelection,
    onChange,
    withoutCheckboxes,
    forceAChoice
  } = props;

  const attemptSelection = (selectedId: string) => {
    if (selectedId !== currentSelection) {
      onChange(selectedId);
    } else if (!forceAChoice) {
      onChange(undefined);
    }
  };

  return (
    <div>
      {options.map(option => {
        return (
          <CheckboxWrapper
            checked={option.id === currentSelection}
            onClick={() => attemptSelection(option.id)}
            withoutCheckboxVisual={withoutCheckboxes}>
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
  );
};
