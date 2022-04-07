import { css } from '@emotion/core';
import styled from '@emotion/styled';
import tw from 'twin.macro';

export const globalStyles = css`
  .ant-collapse-content > .ant-collapse-content-box {
    padding: 0px;
  }
`;

export const OperationDescription = styled('div')`
  ${tw`w-full overflow-y-scroll text-sm text-gray-600 whitespace-normal`};
  grid-column: span 3 / span 3;
  /* Hide scrollbar for Chrome, Safari and Opera */
  &::-webkit-scrollbar {
    display: none !important;
  }
  /* Hide scrollbar for IE, Edge and Firefox */
  & {
    -ms-overflow-style: none !important; /* IE and Edge */
    scrollbar-width: none !important; /* Firefox */
  }
`;
