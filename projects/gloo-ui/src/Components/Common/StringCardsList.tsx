import * as React from 'react';
/** @jsx jsx */
import { jsx } from '@emotion/core';

import styled from '@emotion/styled/macro';
import { colors } from 'Styles';
import { soloConstants } from 'Styles/constants';
import { SoloInput } from './SoloInput';
import { ReactComponent as GreenPlus } from 'assets/small-green-plus.svg';
import { ReactComponent as GreyX } from 'assets/small-grey-x.svg';

export const Container = styled.div`
  display: flex;
  flex-wrap: wrap;
  align-items: center;
`;

export const StringCard = styled.div`
  display: flex;
  justify-content: space-between;
  border-radius: ${soloConstants.smallRadius}px;
  background: ${colors.marchGrey};
  color: ${colors.novemberGrey};
  padding: 0 10px;
  line-height: 33px;
  font-size: 16px;
  width: 175px;
  margin: 10px;
  white-space: nowrap;
`;
export const CardValue = styled.div`
  max-width: 160px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
`;
export const DeleteX = styled.div`
  cursor: pointer;
`;

export const NewStringPrompt = styled.div`
  display: flex;
  justify-content: space-between;
  width: 175px;
  align-items: center;
  margin: 0 10px;
`;
export const PlusHolder = styled<'div', { disabled: boolean }>('div')`
  ${props =>
    props.disabled
      ? `opacity: .5;
    pointer-events: none;`
      : ''}

  margin-left: 5px;
  cursor: pointer;
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
            <CardValue>{value}</CardValue>
            <DeleteX onClick={() => valueDeleted(ind)}>
              <GreyX style={{ marginBottom: '-3px' }} />
            </DeleteX>
          </StringCard>
        );
      })}
      {!!createNew && (
        <NewStringPrompt>
          <SoloInput
            value={newValue}
            placeholder={createNewPromptText}
            onChange={newValueChanged}
          />
          <PlusHolder disabled={!newValue.length} onClick={sendCreateNew}>
            <GreenPlus style={{ marginBottom: '-3px' }} />
          </PlusHolder>
        </NewStringPrompt>
      )}
    </Container>
  );
};
