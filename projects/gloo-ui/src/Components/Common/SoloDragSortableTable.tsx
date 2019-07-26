import * as React from 'react';
/** @jsx jsx */
import { jsx } from '@emotion/core';
import styled from '@emotion/styled/macro';
import { DndProvider, DragSource, DropTarget } from 'react-dnd';
import HTML5Backend from 'react-dnd-html5-backend';
import { TableContainer, TableProps } from './SoloTable';
import { soloConstants } from 'Styles';
import Table from 'antd/lib/table';
import { hslToHSLA, colors } from 'Styles/colors';

let draggingIndex = -1;

const BodyRow = (props: any) => {
  const {
    isOver,
    connectDragSource,
    connectDropTarget,
    moveRow,
    ...restProps
  } = props;

  const style = { ...restProps.style, cursor: 'move' };

  let { className } = restProps;
  if (isOver) {
    if (restProps.index > draggingIndex) {
      className += ' drop-over-downward';
    }
    if (restProps.index < draggingIndex) {
      className += ' drop-over-upward';
    }
  }
  return connectDragSource(
    connectDropTarget(<tr {...restProps} className={className} style={style} />)
  );
};

const rowSource = {
  beginDrag(props: any) {
    draggingIndex = props.index;
    return {
      index: props.index
    };
  }
};

const rowTarget = {
  drop(props: any, monitor: any) {
    const dragIndex = monitor.getItem().index;
    const hoverIndex = props.index;

    // Don't replace items with themselves
    if (dragIndex === hoverIndex) {
      return;
    }

    // Time to actually perform the action
    props.moveRow(dragIndex, hoverIndex);

    // Note: we're mutating the monitor item here!
    // Generally it's better to avoid mutations,
    // but it's good here for the sake of performance
    // to avoid expensive index searches.
    monitor.getItem().index = hoverIndex;
  }
};

const DragableBodyRow = DropTarget('row', rowTarget, (connect, monitor) => ({
  connectDropTarget: connect.dropTarget(),
  isOver: monitor.isOver()
}))(
  DragSource('row', rowSource, connect => ({
    connectDragSource: connect.dragSource()
  }))(BodyRow)
);

interface Props extends TableProps {
  moveRow: (dragIndex: number, hoverIndex: number) => any;
}

export const SoloDragSortableTable = (props: Props) => {
  const components = {
    body: {
      row: DragableBodyRow
    }
  };

  return (
    <TableContainer>
      <DndProvider backend={HTML5Backend}>
        <Table
          dataSource={props.dataSource}
          columns={props.columns}
          components={components}
          onRow={(record: any, index) => ({
            index,
            moveRow: props.moveRow,
            record,
            cols: props.columns
          })}
        />
      </DndProvider>
    </TableContainer>
  );
};
