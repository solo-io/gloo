import * as React from 'react';
/** @jsx jsx */
import { jsx } from '@emotion/core';

import styled from '@emotion/styled/macro';

import { ReactComponent as ColorlessList } from 'assets/listView/list-view-colorless.svg';
import { ReactComponent as ColorlessTile } from 'assets/tileView/tile-view-colorless.svg';
import { colors } from 'Styles';

const Container = styled.div``;

const TileIcon = styled<
  React.FunctionComponent,
  { selected?: boolean; onClick: () => any }
>(ColorlessTile)`
  margin-right: 10px;
  cursor: pointer;

  .colorlessTileView {
    fill: ${colors.juneGrey};
  }
  &:hover {
    .colorlessTileView {
      fill: ${colors.novemberGrey};
    }
  }

  ${props =>
    props.selected
      ? `
        pointer-events: none;

        .colorlessTileView {
          fill: ${colors.lakeBlue};
        }
      `
      : ``};
`;

const ListIcon = styled<
  React.FunctionComponent,
  { selected?: boolean; onClick: () => any }
>(ColorlessList)`
  cursor: pointer;

  .colorlessListView {
    fill: ${colors.juneGrey};
  }
  &:hover {
    .colorlessTileView {
      fill: ${colors.novemberGrey};
    }
  }

  ${props =>
    props.selected
      ? `
        pointer-events: none;

        .colorlessListView {
          fill: ${colors.lakeBlue};
        }
      `
      : ``};
`;

interface Props {
  listIsSelected: boolean;
  onToggle: () => any;
  disabled?: boolean;
  tileselected?: boolean; // If we ever need them both false or true for some reason
}

export const CatalogTableToggle = (props: Props) => {
  const { listIsSelected, onToggle, disabled, tileselected } = props;

  const doToggle = () => {
    if (!disabled) {
      onToggle();
    }
  };

  return (
    <Container>
      <TileIcon
        selected={!listIsSelected || !!tileselected}
        onClick={doToggle}
      />
      <ListIcon selected={listIsSelected} onClick={doToggle} />
    </Container>
  );
};
