import * as React from 'react';
/** @jsx jsx */
import { jsx } from '@emotion/core';

import styled from '@emotion/styled/macro';
import { colors } from 'Styles';
import { soloConstants } from 'Styles/constants';
import { SoloInput } from './SoloInput';
import { ReactComponent as GreenPlus } from 'assets/small-green-plus.svg';
const Container = styled.div`
  display: flex;
  flex-wrap: wrap;
  align-items: center;
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
  margin: 10px;
  white-space: nowrap;
`;
const NewStringPrompt = styled.div`
  display: flex;
  justify-content: space-between;
  width: 175px;
  align-items: center;
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
          <StringCard key={ind}>
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
          <div onClick={sendCreateNew}>
            <GreenPlus style={{ opacity: 0.5 }} />
          </div>
        </NewStringPrompt>
      )}
    </Container>
  );
};
