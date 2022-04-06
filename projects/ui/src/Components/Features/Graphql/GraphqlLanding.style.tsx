import styled from '@emotion/styled';
import { colors } from 'Styles/colors';

export const GraphqlLandingContainer = styled.div`
  position: relative;
  display: grid;
  grid-template-columns: 200px 1fr;
  grid-gap: 28px;
`;

export const HorizontalDivider = styled.div`
  position: relative;
  height: 1px;
  width: 100%;
  background: ${colors.marchGrey};
  margin: 35px 0;

  div {
    position: absolute;
    display: block;
    left: 0;
    right: 0;
    top: 50%;
    margin: -9px auto 0;
    width: 105px;
    text-align: center;
    color: ${colors.septemberGrey};
    background: ${colors.januaryGrey};
  }
`;

export const CheckboxWrapper = styled.div`
  > div {
    width: 190px;
    margin-bottom: 8px;
  }
`;

export const GraphqlIconHolder = styled.div`
  display: flex;
  align-items: center;
  justify-items: center;

  svg {
    width: 35px;
    max-width: none;
  }
`;

export type TableHolderProps = { wholePage?: boolean };
export const TableHolder = styled.div<TableHolderProps>`
  ${(props: TableHolderProps) =>
    props.wholePage
      ? ''
      : `
    table thead.ant-table-thead tr th {
      background: ${colors.marchGrey};
    }
  `};
`;
