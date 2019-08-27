import * as React from 'react';
/** @jsx jsx */
import { jsx } from '@emotion/core';
import styled from '@emotion/styled/macro';
import { colors, soloConstants } from 'Styles';
import Highlight, { defaultProps } from 'prism-react-renderer';
import theme from 'prism-react-renderer/themes/github';

const Container = styled<'div', { whiteBacked?: boolean }>('div')`
  padding: ${soloConstants.smallBuffer}px 5px ${soloConstants.smallBuffer}px
    ${soloConstants.smallBuffer}px;
  border: 1px solid ${colors.mayGrey};
  border-radius: ${soloConstants.smallRadius}px;
  background: ${props => (props.whiteBacked ? 'white' : colors.januaryGrey)};
  .prism-code {
    background: ${props =>
      props.whiteBacked ? 'white' : colors.januaryGrey} !important;
  }
`;

export const Pre = styled.pre`
  text-align: left;

  & .token-line {
    line-height: 1.3em;
    height: 1.3em;
  }
`;

export const LineNo = styled.span`
  display: inline-block;
  width: 2em;
  user-select: none;
  opacity: 0.3;
`;

const ourTheme = { ...theme, backgroundColor: 'transparent' };

interface Props {
  content: string;
  isJson?: boolean;
  whiteBacked?: boolean;
}

export const ConfigDisplayer = React.memo((props: Props) => {
  return (
    <Container whiteBacked={props.whiteBacked}>
      <Highlight
        {...defaultProps}
        theme={ourTheme}
        code={props.content}
        language={props.isJson ? 'json' : 'yaml'}>
        {({ className, style, tokens, getLineProps, getTokenProps }) => (
          <Pre className={className} style={style}>
            {tokens.map((line, i) => (
              <div {...getLineProps({ line, key: i })}>
                <LineNo>{i + 1}</LineNo>
                {line.map((token, key) => (
                  <span {...getTokenProps({ token, key })} />
                ))}
              </div>
            ))}
          </Pre>
        )}
      </Highlight>
    </Container>
  );
});
