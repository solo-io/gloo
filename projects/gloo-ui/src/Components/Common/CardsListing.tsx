import * as React from 'react';
/** @jsx jsx */
import { jsx } from '@emotion/core';

import styled from '@emotion/styled/macro';
import { CardType, Card } from './Card';

const ListContainer = styled.div`
  display: grid;
  grid-template-columns: repeat(auto-fill, 235px);
  grid-gap: 20px;
`;

interface Props {
  cardsData: CardType[];
}

export const CardsListing = (props: Props) => {
  return (
    <ListContainer>
      {props.cardsData.map(cardInfo => {
        return <Card key={cardInfo.id || cardInfo.cardTitle} {...cardInfo} />;
      })}
    </ListContainer>
  );
};
