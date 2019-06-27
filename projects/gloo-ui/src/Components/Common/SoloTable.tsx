import * as React from 'react';
/** @jsx jsx */
import { jsx } from '@emotion/core';
import styled from '@emotion/styled/macro';

import { colors } from '../../Styles';
import Table, { ColumnProps } from 'antd/lib/table';

// To restyle table to match spec later
const TableContainer = styled.div``;

interface Props {
  columns: any[];
  dataSource: any[];
}

export const SoloTable = (props: Props) => {
  return (
    <TableContainer>
      <Table columns={props.columns} dataSource={props.dataSource} />
    </TableContainer>
  );
};
