import * as React from 'react';
/** @jsx jsx */
import { jsx } from '@emotion/core';

import styled from '@emotion/styled/macro';
import { colors } from 'Styles';
import { soloConstants } from 'Styles/constants';
import { SoloInput, Label } from './SoloInput';
import { ReactComponent as GreenPlusSVG } from 'assets/small-green-plus.svg';
import { ReactComponent as GreyX } from 'assets/small-grey-x.svg';
import {
  NewStringPrompt as SmallNewStringPrompt,
  StringCard as SmallStringCard,
  CardValue as SmallCardValue,
  PlusHolder,
  DeleteX as SmallDeleteX
} from './StringCardsList';
import { SoloCheckbox } from './SoloCheckbox';
import { CheckboxChangeEvent } from 'antd/lib/checkbox';
import { hslToHSLA } from 'Styles/colors';

const Container = styled.div`
  display: flex;
  flex-wrap: wrap;
  align-items: center;
`;

export const StringCard = styled(SmallStringCard)`
  width: 200px;
  padding: 0;
  margin-left: 0;
  flex-wrap: wrap;
`;

export const CardValue = styled(SmallCardValue)`
  max-width: none;
  min-width: 0;
  width: calc(50% - 16px);
  padding: 0 5px;
  padding-left: 8px;
  background: white;
  border-top: 1px solid ${colors.februaryGrey};
  border-bottom: 1px solid ${colors.februaryGrey};
`;

export const CardName = styled(SmallCardValue)`
  max-width: none;
  min-width: 0;
  width: calc(50% - 16px);
  padding: 0 5px;
  padding-left: 10px;
`;

const CardBool = styled.div`
  display: flex;
  background: white;
  width: 100%;
  padding: 0 8px;
  font-size: 12px;
  line-height: 18px;
  height: 18px;
`;
const CardBoolIndicator = styled<'div', { isTrue: boolean }>('div')`
  padding: 0 8px;
  background: ${props =>
    props.isTrue ? hslToHSLA(colors.forestGreen, 0.68) : 'transparent'};
  color: white;
  height: 18px;
  margin-right: 8px;
  border-radius: 0 0 8px 8px;
`;

const NewStringPrompt = styled(SmallNewStringPrompt)`
  width: 100%;
  margin: 0;
`;

export const DeleteX = styled(SmallDeleteX)`
  padding: 0 8px;
  margin-left: 0;
`;

const RegexCheckbox = styled.div`
  margin-left: 8px;

  .ant-checkbox-wrapper {
    margin-left: 3px;
  }
`;

const GreenPlus = styled(GreenPlusSVG)`
  width: 16px;
  height: 16px;
`;

export interface MultipartStringCardsProps {
  values?: { name: string; value?: string; boolValue?: boolean }[];
  valueDeleted?: (indexDeleted: number) => any;
  createNew?: (newPair: {
    newName: string;
    newValue: string;
    newBool?: boolean;
  }) => any;
  valueIsValid?: (value: string) => boolean;
  nameIsValid?: (value: string) => boolean;
  createNewNamePromptText?: string;
  createNewValuePromptText?: string;
  title?: string;
  nameSlotTitle?: string;
  valueSlotTitle?: string;
  valuesMayBeEmpty?: boolean;
  booleanFieldText?: string;
  boolSlotTitle?: string;
}

// This badly needs a better name
export const MultipartStringCardsList = (props: MultipartStringCardsProps) => {
  const {
    values,
    valueDeleted,
    createNew,
    valueIsValid,
    nameIsValid,
    createNewNamePromptText,
    createNewValuePromptText,
    title,
    nameSlotTitle,
    valueSlotTitle,
    valuesMayBeEmpty,
    booleanFieldText,
    boolSlotTitle
  } = props;

  const [newName, setNewName] = React.useState<string>('');
  const [newValue, setNewValue] = React.useState<string>('');
  const [newBool, setNewBool] = React.useState<boolean>(false);

  const newNameChanged = (evt: React.ChangeEvent<HTMLInputElement>): void => {
    setNewName(evt.target.value);
  };
  const newValueChanged = (evt: React.ChangeEvent<HTMLInputElement>): void => {
    setNewValue(evt.target.value);
  };
  const newBoolChanged = (evt: CheckboxChangeEvent): void => {
    setNewBool(evt.target.checked);
  };

  const sendCreateNew = () => {
    if (newValue.length > 0 && newName.length > 0) {
      // TODO: Use the slotTitles prop for this case as well
      createNew!({
        newName,
        newValue,
        newBool
      });
      setNewName('');
      setNewValue('');
      setNewBool(false);
    }
  };

  return (
    <div>
      {title && <Label>{title}</Label>}
      <Container>
        {values!.map((value, ind) => {
          // @ts-ignore
          const name = !!nameSlotTitle ? value[nameSlotTitle] : value.name;
          // @ts-ignore
          const val = !!valueSlotTitle ? value[valueSlotTitle] : value.value;

          const boolAccess = !!boolSlotTitle ? boolSlotTitle : 'boolValue';
          // @ts-ignore
          const hasBool = value[boolAccess] !== undefined;
          // @ts-ignore
          const hasBoolTrue = value[boolAccess];

          return (
            <StringCard key={name + ind}>
              <CardName
                hasError={
                  (!!valueIsValid ? !valueIsValid(val) : false) ||
                  (!!nameIsValid ? nameIsValid(name) : false)
                }>
                {name}
              </CardName>
              <CardValue>{val} </CardValue>
              <DeleteX
                onClick={() => valueDeleted!(ind)}
                hasError={
                  (!!valueIsValid ? !valueIsValid(val) : false) ||
                  (!!nameIsValid ? nameIsValid(name) : false)
                }>
                <GreyX style={{ marginBottom: '-3px' }} />
              </DeleteX>
              {hasBool && (
                <CardBool>
                  <CardBoolIndicator isTrue={hasBoolTrue}>
                    {hasBoolTrue && booleanFieldText}
                  </CardBoolIndicator>
                </CardBool>
              )}
            </StringCard>
          );
        })}
        {!!createNew && (
          <div>
            <NewStringPrompt>
              <div style={{ marginRight: '5px' }}>
                <SoloInput
                  value={newName}
                  placeholder={createNewNamePromptText}
                  onChange={newNameChanged}
                  error={
                    !!newName.length &&
                    (!!nameIsValid ? !nameIsValid(newName) : false)
                  }
                />
              </div>
              <SoloInput
                value={newValue}
                placeholder={createNewValuePromptText}
                onChange={newValueChanged}
                error={
                  !!newName.length &&
                  (!!valueIsValid ? !valueIsValid(newValue) : false)
                }
              />
              {!!booleanFieldText && !!booleanFieldText.length && (
                <RegexCheckbox>
                  <SoloCheckbox
                    checked={newBool}
                    title={`${booleanFieldText}?`}
                    onChange={newBoolChanged}
                  />
                </RegexCheckbox>
              )}
              <PlusHolder
                disabled={
                  (!newValue.length && !valuesMayBeEmpty) ||
                  !newName.length ||
                  (!!nameIsValid ? !nameIsValid(newName) : false) ||
                  (!!valueIsValid ? !valueIsValid(newValue) : false)
                }
                onClick={sendCreateNew}
                withRegex={!!booleanFieldText && !!booleanFieldText.length}>
                <GreenPlus />
              </PlusHolder>
            </NewStringPrompt>
          </div>
        )}
      </Container>
    </div>
  );
};
