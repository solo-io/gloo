import styled from '@emotion/styled';
import * as React from 'react';
import { Transfer } from 'antd';
import { colors } from 'Styles';
import { ReactComponent as GreenPlus } from 'assets/small-green-plus.svg';
import { ReactComponent as RedX } from 'assets/small-red-x.svg';
import { ReactComponent as NoSelectedList } from 'assets/no-selected-list.svg';
import { css } from '@emotion/core';

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
  border-radius: 8px;
  border: 1px solid ${colors.aprilGrey};
  height: 100%;
  background: white;
  max-height: 200px;
  overflow-y: scroll;
`;

const Item = styled.div`
  display: flex;
  align-items: center;
  cursor: pointer;
  justify-content: space-between;

  svg {
    opacity: 0;
    pointer-events: none;
  }

  &:hover {
    background-color: ${colors.februaryGrey};
    svg {
      opacity: 1;
      pointer-events: all;
    }
  }
`;

export type ListItemType = {
  name: string;
  namespace: string;
  displayValue?: string;
};

interface TransferProps {
  allOptionsListName: string;
  allOptions: ListItemType[];
  chosenOptionsListName: string;
  chosenOptions: ListItemType[];
  onChange: (newChosenOptions: ListItemType[]) => void;
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

  const removeItem = (itemToRemove: ListItemType) => {
    onChange(
      chosenOptions.filter(
        item =>
          item.name !== itemToRemove.name &&
          item.namespace !== itemToRemove.namespace
      )
    );
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
                  chosenItem =>
                    chosenItem.name === item.name &&
                    chosenItem.namespace === item.namespace
                )
            )
            .map(item => (
              <Item key={`${item.name}-${item.namespace}`}>
                {item.displayValue || `${item.name}-${item.namespace}`}
                <span className='text-green-400 cursor-pointer hover:text-green-300'>
                  <GreenPlus
                    className='w-4 h-4 fill-current'
                    onClick={() => addItem(item)}
                  />
                </span>
              </Item>
            ))}
        </ListHolder>
      </ListHalf>

      <ListHalf>
        <ListTitle>{chosenOptionsListName}</ListTitle>
        <ListHolder>
          {chosenOptions.length === 0 && (
            <div className='flex flex-col items-center justify-center w-full h-full bg-gray-100 rounded-lg'>
              <NoSelectedList className='w-12 h-12' />
              <div className='mt-2 text-gray-500'>Nothing Selected</div>
            </div>
          )}
          {chosenOptions.map(item => (
            <Item key={`${item.name}-${item.namespace}`}>
              {item.displayValue || `${item.name}-${item.namespace}`}
              <RedX onClick={() => removeItem(item)} />
            </Item>
          ))}
        </ListHolder>
      </ListHalf>
    </TransferBlock>
  );
};
