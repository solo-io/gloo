import * as React from 'react';
/** @jsx jsx */
import { jsx } from '@emotion/core';
import styled from '@emotion/styled/macro';
import { colors } from 'Styles';
import Highlight, { defaultProps } from 'prism-react-renderer';
import theme from 'prism-react-renderer/themes/github';

export const Pre = styled.pre`
  text-align: left;
  margin: 1em 0;
  padding: 0.5em;

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

const Heading = styled.div`
  display: inline;
  color: ${colors.seaBlue};
`;
const Content = styled.div`
  display: inline;
`;

interface Props {
  content: string;
  isJson?: boolean;
}

export const YamlDisplayer = React.memo((props: Props) => {
  return (
    <Highlight
      {...defaultProps}
      theme={theme}
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
  );
});
