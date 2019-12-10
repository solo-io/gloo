import styled from '@emotion/styled';
import { ReactComponent as GreenPlus } from 'assets/small-green-plus.svg';
import { ReactComponent as GreyX } from 'assets/small-grey-x.svg';
import * as React from 'react';
import { colors } from 'Styles';
import { soloConstants } from 'Styles/constants';
import { SoloInput, Label } from './SoloInput';
import { SoloTypeahead } from './SoloTypeahead';

export const Container = styled.div`
  display: flex;
  flex-wrap: wrap;
  align-items: center;
`;

const Title = styled.div`
  font-size: 22px;
  font-weight: bold;
  color: ${colors.novemberGrey};
  line-height: normal;
  padding: ${soloConstants.largeBuffer}px ${soloConstants.smallBuffer}px 13px;
`;
type StringCardProps = { limitWidth?: boolean };
export const StringCard = styled.div`
  display: flex;
  justify-content: space-between;
  line-height: 33px;
  font-size: 16px;
  width: ${(props: StringCardProps) => (props.limitWidth ? '175px' : 'auto')};
  margin: 10px;
  white-space: nowrap;
`;
type HasErrorProps = { hasError?: boolean };
export const CardValue = styled.div`
  min-width: 100px;
  max-width: 500px;
  padding-left: 10px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  cursor: default;
  border-radius: ${soloConstants.smallRadius}px 0 0
    ${soloConstants.smallRadius}px;

  ${(props: HasErrorProps) =>
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
export const DeleteX = styled.div`
  padding-right: 10px;
  padding-left: 5px;
  cursor: pointer;
  border-radius: 0 ${soloConstants.smallRadius}px ${soloConstants.smallRadius}px
    0;

  ${(props: HasErrorProps) =>
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

type PlusHolderProps = { disabled: boolean; withRegex?: boolean };

export const PlusHolder = styled.div`
  ${(props: PlusHolderProps) =>
    props.disabled
      ? `opacity: .5;
    pointer-events: none;`
      : ''}

  position: absolute;
  right: ${(props: PlusHolderProps) => (props.withRegex ? '-23px' : '7px')};
  top: 10px;
  cursor: pointer;
  z-index: 5;
`;

export interface StringCardsListProps {
  values: string[];
  label?: string;
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
    label,
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
    if (!!newValue.length && (!!valueIsValid ? valueIsValid(newValue) : true)) {
      createNew!(newValue);
      setNewValue('');
    }
  };

  return (
    <>
      {label && <Label>{label}</Label>}
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
                defaultValue={createNewPromptText}
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
    </>
  );
};
