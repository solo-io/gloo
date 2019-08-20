import * as React from 'react';
/** @jsx jsx */
import { jsx } from '@emotion/core';
import styled from '@emotion/styled/macro';
import { colors, TableActionCircle, soloConstants } from 'Styles';

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
  color: ${colors.seaBlue};
`;
const Content = styled.div`
  display: inline;
`;

interface Props {
  content: string;
}

export const YamlDisplayer = (props: Props) => {
  const lines = props.content.split('\n');

  return (
    <Container>
      <Displayer>
        {lines.map((line, ind) => {
          if (!line.length) {
            return null;
          }

          const sections = line.split(':');
          const heading = sections.splice(0, 1);

          return (
            <pre key={heading[0] + ind}>
              <Heading>{heading}:</Heading>{' '}
              <Content>{sections.join(':')}</Content>
            </pre>
          );
        })}
      </Displayer>
    </Container>
  );
};
