import React, { useState } from 'react';
import styled from '@emotion/styled';
import Table from 'antd/lib/table';
import { colors, hslToHSLA } from 'Styles/colors';
import { SoloLink } from 'Components/Common/SoloLink';
import { HealthIndicator } from './HealthIndicator';
import { ReactComponent as ClusterIcon } from 'assets/cluster-instance-icon.svg';
import { Tooltip } from 'antd';
import { ColumnsType, ExpandableConfig } from 'antd/lib/table/interface';

const ListExtraCount = styled.div`
  margin-left: 4px;
  display: inline-block;
  font-size: 14px;
  line-height: 17px;
  border-radius: 8px;
  border: 1px solid ${colors.marchGrey};
  text-transform: lowercase;
  padding: 0 4px;
`;

const TableHealthCircleHolder = styled.div`
  display: flex;
  align-items: center;
  justify-content: left;

  > div {
    width: 12px;
    height: 12px;
    margin-left: 12px;
  }
`;

export const RenderStatus = (status: 0 | 1 | 2 | 3) => {
  return (
    <TableHealthCircleHolder>
      <HealthIndicator healthStatus={status} />
    </TableHealthCircleHolder>
  );
};

type ClusterCellHolderProps = {
  multiple?: boolean;
};
const ClusterCellHolder = styled.div<ClusterCellHolderProps>`
  display: flex;
  ${(props: ClusterCellHolderProps) =>
    props.multiple ? 'cursor: pointer;' : ''}

  svg {
    width: 18px;
    margin-right: 8px;

    * {
      fill: ${colors.seaBlue};
    }
  }
`;
export const RenderClustersList = (clusters: string[]) => {
  const [isOpen, setIsOpen] = useState(false);

  return (
    <ClusterCellHolder
      multiple={clusters.length > 1}
      onClick={() => setIsOpen(true)}>
      <ClusterIcon />
      {clusters.length === 0 ? null : clusters.length > 1 ? (
        <Tooltip
          title={
            <div>
              {clusters.map(cluster => (
                <div>{cluster}</div>
              ))}
            </div>
          }
          trigger='click'
          visible={isOpen}
          onVisibleChange={() => {
            setIsOpen(!isOpen);
          }}>
          <div>
            {clusters[0]} ({clusters.length})
          </div>
        </Tooltip>
      ) : (
        clusters[0]
      )}
    </ClusterCellHolder>
  );
};
export const RenderCluster = (clusterName: string) => {
  return (
    <ClusterCellHolder>
      <ClusterIcon />
      <SoloLink displayElement={clusterName} link={'/admin/clusters'} />
    </ClusterCellHolder>
  );
};

type ExpandableNamesHolderProps = {
  listLength: number;
};
const ExpandableNamesHolder = styled.div<ExpandableNamesHolderProps>`
  ${(props: ExpandableNamesHolderProps) =>
    props.listLength > 1 ? 'cursor: pointer;' : 'pointer-events: none;'}
`;
export const RenderExpandableNamesList = (names: string[]) => {
  const [isOpen, setIsOpen] = useState(false);

  return (
    <ExpandableNamesHolder
      listLength={names.length}
      onClick={() => setIsOpen(true)}>
      {names.length === 0 ? null : names.length > 1 ? (
        <Tooltip
          title={
            <div>
              {names.map(name => (
                <div>{names}</div>
              ))}
            </div>
          }
          trigger='click'
          visible={isOpen}
          onVisibleChange={() => {
            setIsOpen(!isOpen);
          }}>
          <div>
            {names[0]}{' '}
            {names.length > 1 && (
              <ListExtraCount>+ {names.length - 1}</ListExtraCount>
            )}
          </div>
        </Tooltip>
      ) : (
        names[0]
      )}
    </ExpandableNamesHolder>
  );
};

export const TableActions = styled.div`
  display: flex;
  padding-left: 8px;
`;

export const TableActionCircle = styled.div`
  width: 18px;
  height: 18px;
  line-height: 18px;
  display: flex;
  justify-content: center;
  align-items: center;
  font-weight: normal;
  color: ${colors.novemberGrey};
  border-radius: 18px;
  cursor: pointer;

  background: ${colors.marchGrey};

  &:hover,
  &:focus {
    background: ${colors.mayGrey};
  }

  &:active {
    background: ${colors.marchGrey};
  }

  svg {
    width: 10px;
    height: 10px;

    * {
      fill: ${colors.septemberGrey};
    }
  }
`;

export const TableNotificationCircle = styled.div`
  width: 18px;
  height: 18px;
  line-height: 18px;
  display: flex;
  justify-content: center;
  align-items: center;
  font-weight: normal;
  border-radius: 18px;
  cursor: pointer;
`;

type TableContainerProps = {
  removeShadows?: boolean;
  curved?: boolean;
  flatTopped?: boolean;
  withBorder?: boolean;
  rowHeight?: string;
};
export const TableContainer = styled.div<TableContainerProps>`
  width: 100%;

  ${(props: TableContainerProps) =>
    props.removeShadows
      ? props.withBorder
        ? `border: 1px solid ${colors.boxShadow};`
        : ''
      : `box-shadow: 0px 4px 9px ${colors.boxShadow};`};

  ${(props: TableContainerProps) =>
    props.curved ? `border-radius: 10px;` : ''};

  table {
    background: ${hslToHSLA(colors.marchGrey, 0.15)};
    width: 100%;

    .ant-table-thead {
      ${(props: TableContainerProps) =>
        !props.flatTopped
          ? `
        border-radius: 10px 10px 0 0`
          : ''};

      tr {
        background: ${colors.januaryGrey};

        th {
          background: transparent;
          cursor: default;
          font-weight: 600;
          padding: 12px 20px;
          color: ${colors.novemberGrey};
          text-align: left;
        }
      }

      ${(props: TableContainerProps) =>
        props.curved && !props.flatTopped
          ? `
              tr:first-of-type {
                > th:first-of-type {
                  border-top-left-radius: 10px;
                }
                > th:last-child {
                  border-top-right-radius: 10px;
                }
              }
            `
          : ''};
    }

    .ant-table-tbody {
      background: white;

      .ant-table-row {
        position: relative;

        :not(:last-child) {
          border-bottom: 1px solid ${colors.februaryGrey};
        }

        > td {
          border-color: ${colors.februaryGrey};
          cursor: default;
          padding: 0 20px;
          height: ${(props: TableContainerProps) => props.rowHeight ?? '60px'};
          color: ${colors.septemberGrey};
        }
        &:hover {
          > td {
            background: ${hslToHSLA(colors.marchGrey, 0.5)};
          }
        }
      }

      ${(props: TableContainerProps) =>
        props.curved
          ? `
            tr:last-child {
              > td:first-of-type {
                border-bottom-left-radius: 10px;
              }
              > td:last-child {
                border-bottom-right-radius: 10px;
              }
            }
          `
          : ''};
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

export interface TableProps {
  columns: ColumnsType<any>;
  dataSource: any[];
  formComponent?: React.FC;
  removeShadows?: boolean;
  removePaging?: boolean;
  curved?: boolean;
  flatTopped?: boolean;
  withBorder?: boolean;
  hideHeader?: boolean;
  rowHeight?: string;
  rowClassName?: (rowData: any, index: number) => string;
  expandable?: ExpandableConfig<any>;
}

// TODO: figure out if edit row should always be shown or always be last row
const EditableRow = ({ lastRowID, formComponent, isEmpty, ...props }: any) => {
  const isLastRow = lastRowID === props['data-row-key'];
  const FormComponent = formComponent;

  return (
    <>
      {isLastRow && !!formComponent ? (
        <tr>
          <FormComponent />
        </tr>
      ) : (
        <tr {...props} />
      )}
    </>
  );
};

export const SoloTable = (props: TableProps) => {
  const components = {
    body: {
      row: EditableRow,
    },
  };

  return (
    <TableContainer
      rowHeight={props.rowHeight}
      removeShadows={props.removeShadows}
      curved={props.curved}
      withBorder={props.withBorder}
      flatTopped={props.flatTopped}>
      <Table
        dataSource={props.dataSource}
        columns={props.columns}
        components={components}
        pagination={props.removePaging ? false : { position: ['bottomRight'] }}
        showHeader={!props.hideHeader}
        rowClassName={props.rowClassName}
        expandable={props.expandable}
      />
    </TableContainer>
  );
};
