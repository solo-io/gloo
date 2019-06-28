import * as React from 'react';
/** @jsx jsx */
import { jsx } from '@emotion/core';

import styled from '@emotion/styled/macro';
import { colors } from 'Styles';
import { soloConstants } from 'Styles/constants';
import { SoloInput } from './SoloInput';

const Container = styled.div`
  display: flex;
  flex-wrap: wrap;
`;

const StringCard = styled.div`
  display: flex;
  justify-content: space-between;
  border-radius: ${soloConstants.smallRadius}px;
  background: ${colors.marchGrey};
  color: ${colors.novemberGrey};
  padding: 10px;
  font-size: 16px;
  width: 175px;
  margin-right: 10px;
  margin-bottom: 23px;
`;
const NewStringPrompt = styled.div`
  display: flex;
  justify-content: space-between;
  width: 175px;
`;

interface Props {
  values: string[];
  valueDeleted: (indexDeleted: number) => any;
  createNew?: (newValue: string) => any;
  createNewPromptText?: string;
}

// This badly needs a better name
export const StringCardsList = (props: Props) => {
  const { values, valueDeleted, createNew, createNewPromptText } = props;

  const [newValue, setNewValue] = React.useState<string>('');

  const newValueChanged = (evt: React.ChangeEvent<HTMLInputElement>): void => {
    setNewValue(evt.target.value);
  };

  const sendCreateNew = () => {
    if (newValue.length > 0) {
      createNew!(newValue);
      setNewValue('');
    }
  };

  return (
    <Container>
      {values.map((value, ind) => {
        return (
          <StringCard>
            {value} <div onClick={() => valueDeleted(ind)}>X</div>
          </StringCard>
        );
      })}
      {!!createNew && (
        <NewStringPrompt>
          <SoloInput
            value={newValue}
            placeholder={createNewPromptText}
            onChange={newValueChanged}
          />{' '}
          <div onClick={sendCreateNew}>+</div>
        </NewStringPrompt>
      )}
    </Container>
  );
};
