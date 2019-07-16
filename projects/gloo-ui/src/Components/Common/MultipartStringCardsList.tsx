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

const StringCard = styled(SmallStringCard)`
  width: 200px;
  padding: 0;
  margin-left: 0;
`;

const CardValue = styled(SmallCardValue)`
  max-width: none;
  width: 50%;
  padding: 0 5px;
  padding-left: 8px;
  background: white;
  border-top: 1px solid ${colors.februaryGrey};
  border-bottom: 1px solid ${colors.februaryGrey};
`;
const CardName = styled(CardValue)`
  padding-left: 10px;
  background: transparent;
`;

const NewStringPrompt = styled(SmallNewStringPrompt)`
  width: 100%;
  margin: 0;
`;

const DeleteX = styled(SmallDeleteX)`
  padding: 0 8px;
`;

interface Props {
  values: { name: string; value: string }[];
  valueDeleted: (indexDeleted: number) => any;
  createNew?: (newPair: { newName: string; newValue: string }) => any;
  createNewNamePromptText?: string;
  createNewValuePromptText?: string;
  title?: string;
}

// This badly needs a better name
export const MultipartStringCardsList = (props: Props) => {
  const {
    values,
    valueDeleted,
    createNew,
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
          return (
            <StringCard key={value.name + ind}>
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
                />
              </div>
              <SoloInput
                value={newValue}
                placeholder={createNewValuePromptText}
                onChange={newValueChanged}
              />
              <PlusHolder
                disabled={!newValue.length || !newName.length}
                onClick={sendCreateNew}>
                <GreenPlus style={{ marginBottom: '-3px' }} />
              </PlusHolder>
            </NewStringPrompt>
          </div>
        )}
      </Container>
    </div>
  );
};
