import * as React from 'react';
/** @jsx jsx */
import { jsx } from '@emotion/core';

import styled from '@emotion/styled/macro';
import { colors } from 'Styles';
import { soloConstants } from 'Styles/constants';
import { SoloInput, Label } from './SoloInput';
import { ReactComponent as GreenPlus } from 'assets/small-green-plus.svg';
import { ReactComponent as GreyX } from 'assets/small-grey-x.svg';
import {
  NewStringPrompt as SmallNewStringPrompt,
  StringCard as SmallStringCard,
  CardValue as SmallCardValue,
  PlusHolder,
  DeleteX as SmallDeleteX
} from './StringCardsList';

const Container = styled.div`
  display: flex;
  flex-wrap: wrap;
  align-items: center;
`;

export const StringCard = styled(SmallStringCard)`
  width: 200px;
  padding: 0;
  margin-left: 0;
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

export const CardName = styled(CardValue)`
  padding-left: 10px;
  background: transparent;
  border-top: none;
  border-bottom: none;
`;

const NewStringPrompt = styled(SmallNewStringPrompt)`
  width: 100%;
  margin: 0;
`;

export const DeleteX = styled(SmallDeleteX)`
  padding: 0 8px;
  margin-left: 0;
`;

export interface MultipartStringCardsProps {
  values: { name: string; value: string }[];
  valueDeleted: (indexDeleted: number) => any;
  createNew?: (newPair: { newName: string; newValue: string }) => any;
  valueIsValid?: (value: string) => boolean;
  nameIsValid?: (value: string) => boolean;
  createNewNamePromptText?: string;
  createNewValuePromptText?: string;
  title?: string;
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
    title
  } = props;

  const [newName, setNewName] = React.useState<string>('');
  const [newValue, setNewValue] = React.useState<string>('');

  const newNameChanged = (evt: React.ChangeEvent<HTMLInputElement>): void => {
    setNewName(evt.target.value);
  };
  const newValueChanged = (evt: React.ChangeEvent<HTMLInputElement>): void => {
    setNewValue(evt.target.value);
  };

  const sendCreateNew = () => {
    if (newValue.length > 0 && newName.length > 0) {
      createNew!({
        newName,
        newValue
      });
      setNewName('');
      setNewValue('');
    }
  };

  return (
    <div>
      {title && <Label>{title}</Label>}
      <Container>
        {values.map((value, ind) => {
          console.log(value);
          return (
            <StringCard
              key={value.name + ind}
              hasError={
                (!!valueIsValid ? !valueIsValid(value.value) : false) ||
                (!!nameIsValid ? nameIsValid(value.name) : false)
              }>
              <CardName>{value.name}</CardName>
              <CardValue>{value.value} </CardValue>
              <DeleteX onClick={() => valueDeleted(ind)}>
                <GreyX style={{ marginBottom: '-3px' }} />
              </DeleteX>
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
              <PlusHolder
                disabled={
                  !newValue.length ||
                  !newName.length ||
                  (!!nameIsValid ? !nameIsValid(newName) : false) ||
                  (!!valueIsValid ? !valueIsValid(newValue) : false)
                }
                onClick={sendCreateNew}>
                <GreenPlus style={{ width: '16px', height: '16px' }} />
              </PlusHolder>
            </NewStringPrompt>
          </div>
        )}
      </Container>
    </div>
  );
};
