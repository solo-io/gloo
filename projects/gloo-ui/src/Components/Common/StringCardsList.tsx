import * as React from 'react';
/** @jsx jsx */
import { jsx } from '@emotion/core';

import styled from '@emotion/styled/macro';
import { colors } from 'Styles';
import { soloConstants } from 'Styles/constants';
import { SoloInput } from './SoloInput';
import { ReactComponent as GreenPlus } from 'assets/small-green-plus.svg';
import { ReactComponent as GreyX } from 'assets/small-grey-x.svg';
import { SoloTypeahead } from './SoloTypeahead';

export const Container = styled.div`
  display: flex;
  flex-wrap: wrap;
  align-items: center;
`;

export const StringCard = styled<'div', { limitWidth?: boolean }>('div')`
  display: flex;
  justify-content: space-between;
  line-height: 33px;
  font-size: 16px;
  width: ${props => (props.limitWidth ? '175px' : 'auto')};
  margin: 10px;
  white-space: nowrap;
`;
export const CardValue = styled<'div', { hasError?: boolean }>('div')`
  min-width: 100px;
  max-width: 500px;
  padding-left: 10px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  cursor: default;
  border-radius: ${soloConstants.smallRadius}px 0 0
    ${soloConstants.smallRadius}px;

  ${props =>
    props.hasError
      ? `background: ${colors.tangerineOrange};
      color: ${colors.pumpkinOrange};
      
      .greyX-c {
        fill: ${colors.pumpkinOrange};
      }`
      : `
      background: ${colors.marchGrey};
      color: ${colors.novemberGrey};`}
`;
export const DeleteX = styled<'div', { hasError?: boolean }>('div')`
  padding-right: 10px;
  padding-left: 5px;
  cursor: pointer;
  border-radius: 0 ${soloConstants.smallRadius}px ${soloConstants.smallRadius}px
    0;

  ${props =>
    props.hasError
      ? `background: ${colors.tangerineOrange};
    color: ${colors.pumpkinOrange};
    
    .greyX-c {
      fill: ${colors.pumpkinOrange};
    }`
      : `
    background: ${colors.marchGrey};
    color: ${colors.novemberGrey};`}
`;

export const NewStringPrompt = styled.div`
  position: relative;
  display: flex;
  justify-content: space-between;
  width: 175px;
  align-items: center;
  margin: 10px;
`;
export const PlusHolder = styled<
  'div',
  { disabled: boolean; withRegex?: boolean }
>('div')`
  ${props =>
    props.disabled
      ? `opacity: .5;
    pointer-events: none;`
      : ''}

  position: absolute;
  right: ${props => (props.withRegex ? '-23px' : '7px')};
  top: 10px;
  cursor: pointer;
  z-index: 5;
`;

export interface StringCardsListProps {
  values: string[];
  valueDeleted: (indexDeleted: number) => any;
  createNew?: (newValue: string) => any;
  valueIsValid?: (value: string) => boolean;
  createNewPromptText?: string;
  asTypeahead?: boolean;
  presetOptions?: string[];
}

// This badly needs a better name
export const StringCardsList = (props: StringCardsListProps) => {
  const {
    values,
    valueDeleted,
    createNew,
    createNewPromptText,
    valueIsValid,
    asTypeahead,
    presetOptions
  } = props;

  const [newValue, setNewValue] = React.useState<string>('');

  const newValueChanged = (evt: React.ChangeEvent<HTMLInputElement>): void => {
    setNewValue(evt.target.value);
  };

  const sendCreateNew = () => {
    createNew!(newValue);
    setNewValue('');
  };

  return (
    <Container>
      {values.map((value, ind) => {
        return (
          <StringCard key={ind}>
            <CardValue
              title={value}
              hasError={!!valueIsValid ? !valueIsValid(value) : false}>
              {value}
            </CardValue>
            <DeleteX
              onClick={() => valueDeleted(ind)}
              hasError={!!valueIsValid ? !valueIsValid(value) : false}>
              <GreyX style={{ marginBottom: '-3px' }} />
            </DeleteX>
          </StringCard>
        );
      })}
      {!!createNew && (
        <NewStringPrompt>
          {asTypeahead ? (
            <SoloTypeahead
              placeholder={createNewPromptText}
              onChange={value => setNewValue(value)}
              presetOptions={presetOptions!.map(pO => {
                return { value: pO };
              })}
              onKeyPress={(e: React.KeyboardEvent) =>
                e.key === 'Enter' ? sendCreateNew() : {}
              }
              hideArrow
            />
          ) : (
            <SoloInput
              value={newValue}
              placeholder={createNewPromptText}
              onChange={newValueChanged}
              onKeyPress={(e: React.KeyboardEvent) =>
                e.key === 'Enter' ? sendCreateNew() : {}
              }
              error={
                !!newValue.length &&
                (!!valueIsValid ? !valueIsValid(newValue) : false)
              }
            />
          )}
          <PlusHolder
            disabled={
              !newValue.length ||
              (!!valueIsValid ? !valueIsValid(newValue) : false)
            }
            onClick={sendCreateNew}>
            <GreenPlus
              style={{ width: '16px', height: '16px', cursor: 'pointer' }}
            />
          </PlusHolder>
        </NewStringPrompt>
      )}
    </Container>
  );
};
