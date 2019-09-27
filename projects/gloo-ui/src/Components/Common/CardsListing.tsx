import styled from '@emotion/styled';
import * as React from 'react';
import { colors } from 'Styles';
import { Card, CardType } from './Card';

const ListContainer = styled.div`
  display: grid;
  grid-template-columns: repeat(auto-fill, 235px);
  grid-gap: 20px;
  margin-bottom: 20px;
`;

interface Props {
  cardsData: CardType[];
  title?: string;
  emptyContent?: React.ReactNode;
}

const CardsListingTitle = styled.div`
  font-size: 18px;
  color: ${colors.novemberGrey};
  margin-bottom: 15px;
  font-weight: 700;
`;
const Container = styled.div`
  display: flex;
  flex-direction: column;
`;

export const CardsListing = (props: Props) => {
  const { title, cardsData } = props;

  if (!cardsData.length) {
    return null;
  }

  return (
    <Container>
      {title && <CardsListingTitle>{title}</CardsListingTitle>}
      {!!cardsData.length || !!props.emptyContent ? (
        <ListContainer>
          {cardsData.map(cardInfo => {
            return (
              <Card
                key={cardInfo.cardTitle + (cardInfo.cardSubtitle || '')}
                {...cardInfo}
              />
            );
          })}
        </ListContainer>
      ) : (
        <>{props.emptyContent}</>
      )}
    </Container>
  );
};
