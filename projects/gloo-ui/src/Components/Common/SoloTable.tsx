import * as React from 'react';
/** @jsx jsx */
import { jsx } from '@emotion/core';
import styled from '@emotion/styled/macro';

import { colors, soloConstants } from '../../Styles';
import Table, { ColumnProps } from 'antd/lib/table';
import { hslToHSLA } from 'Styles/colors';

// To restyle table to match spec later
const TableContainer = styled.div`
  box-shadow: 0px 4px 9px ${colors.boxShadow};

  .ant-table-wrapper {
    background: ${hslToHSLA(colors.marchGrey, 0.15)};

    .ant-table-thead {
      border-radius: ${soloConstants.radius}px ${soloConstants.radius}px 0 0;

      tr {
        background: ${colors.marchGrey};

        .ant-table-column-title {
          cursor: default;
          font-weight: 600;
          color: ${colors.novemberGrey};
        }
      }
    }

    .ant-table-tbody {
      background: white;

      .ant-table-row {
        position: relative;

        > td {
          border-color: ${colors.februaryGrey};
        }
        &:hover {
          > td {
            background: ${hslToHSLA(colors.marchGrey, 0.5)};
          }
        }
      }
    }

    .ant-table-pagination {
      &.ant-pagination {
        margin: 0;
      }

      &[unselectable='unselectable'] {
        opacity: 0.25;
      }

      .ant-pagination-prev,
      .ant-pagination-next {
      }

      a,
      a.ant-pagination-item-link,
      .ant-pagination-item-active {
        background: none;
        border: none;
        color: ${colors.septemberGrey};
        line-height: 30px;
        height: 30px;
      }
    }

    .ant-empty-description {
      color: white;

      &::after {
        content: 'No Matches';
        position: absolute;
        left: 0;
        right: 0;
        text-align: center;
        color: ${colors.juneGrey};
      }
    }
  }
`;

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
