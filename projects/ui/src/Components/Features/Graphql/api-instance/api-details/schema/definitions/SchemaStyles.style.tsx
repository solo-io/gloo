import { css } from '@emotion/react';
import styled from '@emotion/styled';
import { colors } from 'Styles/colors';

export namespace SchemaStyles {
  export const Field = styled.div<{ spaceY?: number }>(
    ({ spaceY = 1 }) => css`
      margin-top: ${spaceY / 2}rem;
      margin-left: 1rem;
      padding-top: ${spaceY / 2}rem;
      padding-bottom: ${spaceY}rem;
      border-bottom: 2px solid ${colors.februaryGrey};
      &:last-child {
        border-bottom: none;
      }
    `
  );

  export const Description = styled('div')`
    padding-top: 0.25rem;
    padding-bottom: 0.25rem;
    color: ${colors.juneGrey};
    font-style: italic;
  `;

  export const SelectionSet = styled.div<{ outerPadding: boolean }>(
    ({ outerPadding }) => `
    margin-left: 0.5rem;
    padding-top: ${outerPadding ? '1rem' : '0'};
    padding-bottom: ${outerPadding ? '1rem' : '0'};
  `
  );
}
