import styled from '@emotion/styled';
import { colors } from 'Styles/colors';

export const TableActions = styled.div`
  display: grid;
  grid-template-columns: 18px 18px 18px;
  grid-gap: 5px;
`;

export const TableActionCircle = styled.div`
  width: 18px;
  height: 18px;
  line-height: 18px;
  text-align: center;
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
`;

export const TableHealthCircleHolder = styled.div`
  display: inline;

  > div {
    width: 10px;
    height: 10px;
    margin-left: 0;
    margin-right: 5px;
  }
`;
