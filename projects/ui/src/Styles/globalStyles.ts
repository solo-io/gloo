import { css } from '@emotion/core';
import { colors } from './colors';

export const globalStyles = css`
  html,
  body {
    width: 100vw;
    height: 100vh;

    line-height: 19px;
  }

  body {
    font-family: 'Proxima Nova', 'Open Sans', 'Helvetica', 'Arial', 'sans-serif';
    margin: 0;
    padding: 0;
    min-height: 100vh;
    min-width: 100vw;
    background: ${colors.januaryGrey};

    * {
      box-sizing: border-box;
      font-family: 'Proxima Nova', 'Open Sans', 'Helvetica', 'Arial',
        'sans-serif';
    }

    a {
      color: #2196c9;

      &:hover,
      &:focus,
      &:active {
        color: #54b7e3;
      }
    }

    table {
      border-collapse: collapse;
    }

    input {
      font-size: 16px;
    }

    .ant-modal-content {
      border-radius: 10px;
      box-shadow: hsla(0, 0%, 0%, 0.1) 0 4px 9px;

      .ant-modal-title {
        font-size: 24px;
        line-height: 26px;
      }
    }

    .ant-popover {
      .ant-popover-content {
        min-width: 125px;

        .ant-popover-message-title {
          color: white;
        }
        .ant-popover-inner {
          background: ${colors.novemberGrey};
          border-radius: 2px;

          .ant-popover-inner-content {
            color: white;
          }
        }
        .ant-popover-arrow {
          border-color: ${colors.novemberGrey};
        }
      }
    }

    ul.ant-select-tree {
      width: 100%;
      /* margin-bottom: 15px; */
      line-height: 16px;
      padding: 0;

      .anticon-caret-down {
        visibility: hidden;
      }

      .treeIcon {
        display: inline-flex;
        justify-content: center;
        align-items: center;
        width: 18px;
        height: 18px;
        border-radius: 18px;
        background: ${colors.februaryGrey};

        svg {
          max-height: 14px;
          width: 14px;
        }
      }

      > li > ul {
        /*padding-left: 15px !important;*/
      }

      li {
        span.ant-select-tree-switcher {
          display: none;
        }

        .ant-select-tree-node-content-wrapper {
          width: calc(100%);
        }

        &.ant-select-tree-treenode-switcher-open > span {
          color: ${colors.juneGrey};
          cursor: default;
          background: white;
          padding: 0 5px !important;
          line-height: 18px;

          &:hover {
            background: white;
          }
        }
      }

      .ant-select-tree-node-content-wrapper {
        &.ant-select-tree-node-content-wrapper-open {
          + ul {
            padding: 0;
          }
        }

        &.ant-select-tree-node-content-wrapper-normal {
          padding-left: 20px !important;

          &:hover {
            background: ${colors.splashBlue};
          }

          &.ant-select-tree-node-content-wrapper-normal {
            padding-left: 20px !important;
            color: ${colors.novemberGrey};
            cursor: pointer;

            &:hover {
              background: ${colors.splashBlue};
            }

            &.ant-select-tree-node-selected {
              cursor: default;
              background: ${colors.februaryGrey};

              padding: 9px 15px 9px 11px;
              border: none;
              border-radius: 0;
              height: auto;
              outline: none;

              .ant-select-tree-node-selected__rendered {
                line-height: inherit;
                margin: 0;

                .ant-select-tree-node-selected-selected-value {
                  color: ${colors.septemberGrey};
                }
              }

              &:disabled {
                background: ${colors.aprilGrey};
              }
            }
          }
        }
      }

      .ant-select-dropdown {
        background: white;
        padding: 0 10px;
        border: 1px solid ${colors.aprilGrey};
        border-radius: 0 0 8px 8px;
        width: 0;

        .ant-select-item {
          display: flex;
          align-items: center;
          padding: 5px 0;
          color: ${colors.septemberGrey};
          cursor: pointer;

          svg {
            height: 20px;
            width: 20px;
            min-width: 20px;
            margin-right: 8px;
          }
        }
      }
    }

    .ace_editor span,
    .ace_editor textarea,
    .ace_line {
      font: 16px / normal 'Monaco', 'Menlo', 'Ubuntu Mono', 'Consolas',
        'source-code-pro', monospace;
    }
    .ace_layer {
      .ace_line,
      .ace_active-line {
        height: 22px !important;
      }
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
  }
`;
