import { css } from '@emotion/core';
import { colors } from 'Styles/colors';

export const resolversTableStyles = css`
  .ant-collapse-content > .ant-collapse-content-box {
    padding: 0px;
  }

  .ant-collapse .ant-collapse-header {
    background-color: ${colors.januaryGrey};
    &:hover {
      background-color: ${colors.februaryGrey};
    }
    &:active {
      background-color: ${colors.marchGrey};
    }
    user-select: none;
  }
`;
