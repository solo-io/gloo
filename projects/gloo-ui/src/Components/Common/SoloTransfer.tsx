import styled from '@emotion/styled';
import * as React from 'react';
import { Transfer } from 'antd';
import { colors } from 'Styles';
import { ReactComponent as GreenPlus } from 'assets/small-green-plus.svg';
import { ReactComponent as RedX } from 'assets/small-red-x.svg';

const TransferBlock = styled.div`
  display: flex;
  justify-content: space-between;
  width: 100%;
`;

const ListHalf = styled.div`
  width: calc(50% - 14px);
  font-size: 16px;
`;

const ListTitle = styled.div`
  font-weight: 500;
  margin-bottom: 10px;
`;

const ListHolder = styled.div`
  color: ${colors.septemberGrey};
  padding: 9px;
  border 1px solid ${colors.aprilGrey};
  background: white;
`;

const Item = styled.div`
  display: flex;
  justify-content: space-between;

  svg {
    opacity: 0;
    pointer-events: none;
  }

  &:hover {
    svg {
      opacity: 1;
      pointer-events: all;
    }
  }
`;

type ListItemType = {
  value: string;
  displayValue?: string;
};
interface TransferProps {
  allOptionsListName: string;
  allOptions: ListItemType[];
  chosenOptionsListName: string;
  chosenOptions: ListItemType[];
  onChange: (newChosenOptions: ListItemType[]) => any;
}

export const SoloTransfer = (props: TransferProps) => {
  const {
    allOptionsListName,
    chosenOptionsListName,
    allOptions,
    chosenOptions,
    onChange
  } = props;

  const addItem = (addedItem: ListItemType) => {
    onChange([...chosenOptions, addedItem]);
  };

  const removeItem = (addedItem: ListItemType) => {
    onChange(chosenOptions.filter(lItem => lItem.value !== addedItem.value));
  };

  return (
    <TransferBlock>
      <ListHalf>
        <ListTitle>{allOptionsListName}</ListTitle>
        <ListHolder>
          {allOptions
            .filter(
              item =>
                !chosenOptions.find(
                  chosenItem => chosenItem.value === item.value
                )
            )
            .map(item => (
              <Item key={item.value}>
                {item.displayValue || item.value}{' '}
                <GreenPlus onClick={() => addItem(item)} />
              </Item>
            ))}
        </ListHolder>
      </ListHalf>

      <ListHalf>
        <ListTitle>{chosenOptionsListName}</ListTitle>
        <ListHolder>
          {chosenOptions.map(item => (
            <Item key={item.value}>
              {item.displayValue || item.value}{' '}
              <RedX onClick={() => removeItem(item)} />
            </Item>
          ))}
        </ListHolder>
      </ListHalf>
    </TransferBlock>
  );
};
