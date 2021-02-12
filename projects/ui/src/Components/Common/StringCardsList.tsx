import styled from '@emotion/styled';
import { ReactComponent as GreenPlus } from 'assets/small-green-plus.svg';
import { ReactComponent as GreyX } from 'assets/small-grey-x.svg';
import * as React from 'react';
import { SoloInput, Label } from './SoloInput';
import { colors } from 'Styles/colors';

export const Container = styled.div`
  display: flex;
  flex-wrap: wrap;
  align-items: center;
`;

export const StringCard = styled.div`
  display: flex;
  justify-content: space-between;
  font-size: 16px;
  width: auto;
  margin-right: 10px;
  max-width: 500px;
  padding: 8px 10px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  cursor: default;
  border-radius: 8px;
  background: ${colors.marchGrey};
  color: ${colors.novemberGrey};
`;

export interface StringCardsListProps {
  values: React.ReactNode[];
}

// This badly needs a better name
export const StringCardsList = (props: StringCardsListProps) => {
  const { values } = props;

  return (
    <>
      <Container>
        {values.map((value, ind) => {
          return (
            <StringCard key={ind} title={value?.toLocaleString()}>
              {value}
            </StringCard>
          );
        })}
      </Container>
    </>
  );
};
