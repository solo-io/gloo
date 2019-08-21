import * as React from 'react';
/** @jsx jsx */
import { jsx } from '@emotion/core';
import styled from '@emotion/styled/macro';
import { colors, soloConstants } from 'Styles';

const Container = styled.div`
  padding: ${soloConstants.smallBuffer}px 5px ${soloConstants.smallBuffer}px
    ${soloConstants.smallBuffer}px;
  border: 1px solid ${colors.mayGrey};
  border-radius: ${soloConstants.smallRadius}px;
  background: ${colors.januaryGrey};
`;

const Displayer = styled.div`
  display: block;
  max-height: 75vh;
  overflow-y: auto;
  color: ${colors.septemberGrey};

  pre {
    margin: 0;
    font-family: 'Proxima Nova', 'Open Sans', 'Helvetica', 'Arial', 'sans-serif';
  }
`;

const Heading = styled.div`
  display: inline;
`;
const Content = styled<'div', { type: 'number' | 'text' | 'other' }>('div')`
  display: inline;
  color: ${props =>
    props.type === 'text'
      ? colors.seaBlue
      : props.type === 'number'
      ? colors.grapefruitOrange
      : 'inherit'};
`;

interface Props {
  content: string;
}

export const JsonDisplayer = (props: Props) => {
  const lines = /*JSON.stringify(JSON.parse(props.content)).split('\n');*/ props.content.split(
    '\n'
  );

  return (
    <Container>
      <Displayer>
        {lines.map((line, ind) => {
          if (!line.length) {
            return null;
          }

          const sections = line.split(':');
          const heading = sections.splice(0, 1);
          const sectionsType =
            !sections[0] || !sections[0].length
              ? 'other'
              : sections[0].charAt(1) === '"'
              ? 'text'
              : Number.isInteger(parseInt(sections[0].charAt(1)))
              ? 'number'
              : 'other';

          return (
            <pre key={heading[0] + ind}>
              <Heading>
                {heading}
                {!!sections.length ? ':' : ''}
              </Heading>{' '}
              {!!sections.length && (
                <Content type={sectionsType}>{sections.join(':')}</Content>
              )}
            </pre>
          );
        })}
      </Displayer>
    </Container>
  );
};
