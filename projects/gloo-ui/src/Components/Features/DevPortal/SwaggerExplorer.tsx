// import SwaggerUI from 'swagger-ui-react';
// import 'swagger-ui-react/swagger-ui.css';
// import { css } from '@emotion/core';

// export const SwaggerExplorer = () => {
//   return (
//     <div className='bg-white rounded-lg'>
//       <SwaggerUI url='https://petstore.swagger.io/v2/swagger.json' />;
//     </div>
//   );
// };
import * as React from 'react';
import styled from '@emotion/styled';
import SwaggerUI from 'swagger-ui-react';
import 'swagger-ui-react/swagger-ui.css';
import { colors, hslToHSLA } from 'Styles/colors';

type SwaggerHolderContainerProps = {
  visibleTag: string;
};
const SwaggerHolderContainer = styled.div`
  .information-container.wrapper {
    display: none;
  }
  .swagger-ui {
    position: relative;
    padding-top: 45px;
    @media (max-width: 800px) {
      padding-top: 90px;
    }
    .scheme-container {
      position: absolute;
      top: 0;
      right: 0;
      max-width: 300px;
      padding: 0;
      background: none;
      box-shadow: none;
      @media (max-width: 800px) {
        right: auto;
        left: 30px;
      }
      .schemes-title {
        display: none;
      }
      .schemes > label select {
        line-height: 22px;
      }
      .btn {
        padding: 5px 12px;
        margin-left: 5px;
        border-radius: 8px;
        border: none;
        line-height: 35px;
        min-width: 150px;
        padding: 0;
        display: flex;
        justify-content: center;
        align-items: center;
        background: ${colors.juneGrey};
        color: white;
        &.authorize {
          margin-right: 0;
          background: ${colors.forestGreen};
          color: white;
          svg {
            height: 16px;
            margin-left: 8px;
            fill: ${colors.groveGreen};
          }
        }
        span {
          padding: 0;
          text-align: center;
        }
      }
      label select {
        outline: none;
        cursor: pointer;
        border: 1px solid ${colors.juneGrey};
        color: ${colors.juneGrey};
        background: white;
        box-shadow: none;
      }
    }
    .wrapper {
      padding: 0;
      max-width: none;
      margin-bottom: 25px;
      > .block {
        background: white;
        border-radius: 10px;
        box-shadow: 0 0 15px ${colors.boxShadow};
        padding: 0 25px 25px;
      }
      + .wrapper > .block {
        padding: 18px 15px;
      }
    }
    /* HEADER + CONTROLING VISIBILITY */
    h4.opblock-tag {
      display: none;
      pointer-events: none;
      line-height: 26px;
      padding: 0;
      color: ${colors.novemberGrey};
      a.nostyle {
        font-size: 22px;
      }
      > small {
        font-size: 16px;
        color: ${colors.septemberGrey};
      }
      > div > small {
        color: ${colors.novemberGrey};
        font-size: 14px;
        a {
          color: ${colors.seaBlue};
          font-size: 16px;
        }
      }
      button {
        display: none;
      }
      & + div {
        display: none;
      }
    }
    h4#operations-tag-${(props: SwaggerHolderContainerProps) =>
        props.visibleTag} {
      display: flex;
      border: none;
      & + div {
        display: block;
      }
    }
    /* REQUEST OPTIONS */
    .opblock {
      background: white;
      border: 1px solid ${colors.juneGrey};
      box-shadow: none;
      .opblock-summary-method {
        padding: 6px 0;
      }
      /* Get */
      &.opblock-get {
        &.is-open {
          border: 1px solid ${hslToHSLA(colors.seaBlue, 0.8)};
        }
        .opblock-summary-method {
          color: ${colors.seaBlue};
          background: ${colors.dropBlue};
          border: 1px solid ${colors.seaBlue};
        }
      }
      /* Post */
      &.opblock-post {
        &.is-open {
          border: 1px solid ${hslToHSLA(colors.forestGreen, 0.8)};
        }
        .opblock-summary-method {
          color: ${colors.forestGreen};
          background: ${colors.groveGreen};
          border: 1px solid ${colors.forestGreen};
        }
      }
      /* Delete */
      &.opblock-delete {
        &.is-open {
          border: 1px solid ${hslToHSLA(colors.grapefruitOrange, 0.8)};
        }
        .opblock-summary-method {
          color: ${colors.grapefruitOrange};
          background: #ffeade;
          border: 1px solid ${colors.grapefruitOrange};
        }
      }
      /* Put */
      &.opblock-put {
        &.is-open {
          border: 1px solid ${hslToHSLA(colors.sunGold, 0.8)};
        }
        .opblock-summary-method {
          color: ${colors.sunGold};
          background: ${colors.flashlightGold};
          border: 1px solid ${colors.sunGold};
        }
      }
      &.opblock-deprecated {
        &.is-open {
          border: 1px solid ${colors.juneGrey};
        }
        .opblock-summary-method {
          color: ${colors.juneGrey};
          background: ${colors.januaryGrey};
          border: 1px solid ${colors.juneGrey};
        }
      }
    }
    /* MODAL */
    .modal-ux {
      .auth-btn-wrapper {
        justify-content: flex-end;
      }
      .modal-ux-header {
        border: none;
        padding: 23px 0 0;
        h3 {
          font-size: 22px;
          font-weight: 500;
        }
        svg {
          fill: ${colors.juneGrey};
        }
      }
      .modal-ux-content {
        padding: 20px 0;
        h6 {
          margin: 0;
        }
        .wrapper {
          margin-bottom: 0;
        }
        .auth-container {
          border: none;
        }
        .scope-def {
          border-radius: 7px;
          background: ${colors.januaryGrey};
        }
        .scopes .checkbox input[type='checkbox'] {
          & + label {
            cursor: pointer;
            .item {
              box-shadow: none;
              border: 1px solid ${colors.septemberGrey};
              background-color: ${colors.januaryGrey};
              border-radius: 5px;
              height: 21px;
              width: 21px;
            }
          }
          &:checked + label .item {
            background-color: ${colors.puddleBlue};
            border-color: ${colors.seaBlue};
          }
        }
      }
    }
    /* MODELS */
    section.models {
      padding: 0;
      margin: 0;
      border: none;
      > h4 {
        font-size: 18px;
        padding: 10px 0;
        svg {
          fill: ${colors.juneGrey};
          width: 12px;
        }
        &:hover {
          background: transparent;
        }
      }
      .model-container {
        margin: 15px 0;
        border-radius: 7px;
        border: 1px solid ${colors.marchGrey};
        background: ${colors.januaryGrey};
        line-height: 50px;
        .model-box {
          padding: 0;
          width: 100%;
        }
        > .model-box {
          padding: 0 16px;
          display: flex;
          justify-content: space-between;
          > .model-box > .model > span {
            display: grid;
            width: 100%;
            grid-template-columns: 1fr 1fr;
            > span:nth-child(2) {
              text-align: right;
            }
            .brace-open,
            .inner-object,
            .brace-close {
              grid-column: 1 / span 2;
            }
          }
        }
      }
    }
  }
`;
const OperationTagsHolder = styled.div`
  position: relative;
  margin-top: 25px;
`;
const OperationTagsList = styled.div`
  position: absolute;
  bottom: -45px;
  left: 28px;
  height: 45px;
  display: flex;
  z-index: 2;
  @media (max-width: 800px) {
    bottom: -90px;
  }
`;
type OperationTagProps = {
  isActive: boolean;
};
const OperationTag = styled.div`
  border: 1px solid ${colors.marchGrey};
  border-bottom: none;
  padding: 15px 20px;
  border-radius: 10px 10px 0 0;
  margin-right: 5px;
  font-weight: 500;
  cursor: pointer;
  ${(props: OperationTagProps) =>
    props.isActive
      ? `
    background: white;
    color: ${colors.seaBlue};
    `
      : `
    background: ${colors.februaryGrey};
    color: ${colors.septemberGrey};
    `}
`;
interface SwaggerHolderProps {
  swaggerJSON: object;
}
export const SwaggerExplorer = ({ swaggerJSON }: SwaggerHolderProps) => {
  let tagsList: string[];
  // @ts-ignore
  tagsList = swaggerJSON['tags'].map(t => t.name);
  console.log(tagsList);
  const [activeTag, setActiveTag] = React.useState('');
  const clickActiveTag = (newActiveTag?: string) => {
    console.log(newActiveTag, activeTag);
    if (!!newActiveTag && newActiveTag !== activeTag) {
      setActiveTag(newActiveTag);
    } else {
      // On switch, make sure things are 'open'
      const tagHeaader = document.getElementById(`operations-tag-${activeTag}`);
      console.log(
        !!tagHeaader,
        tagHeaader?.getAttribute('data-is-open') === 'true'
      );
      if (!!tagHeaader && tagHeaader.getAttribute('data-is-open') === 'false') {
        if (typeof tagHeaader.click === 'function') {
          tagHeaader.click();
        } else if (typeof tagHeaader.onclick === 'function') {
          // @ts-ignore
          tagHeaader.onclick();
        }
      }
    }
  };
  React.useEffect(() => {
    clickActiveTag(tagsList[0]);
  }, []);
  React.useEffect(() => {
    clickActiveTag();
  }, [activeTag]);
  return (
    <SwaggerHolderContainer visibleTag={activeTag}>
      <OperationTagsHolder>
        <OperationTagsList>
          {tagsList.map(tag => (
            <OperationTag
              isActive={tag === activeTag}
              onClick={() => setActiveTag(tag)}>
              {tag}
            </OperationTag>
          ))}
        </OperationTagsList>
      </OperationTagsHolder>
      <SwaggerUI spec={swaggerJSON} />
    </SwaggerHolderContainer>
  );
};
