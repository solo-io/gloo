import styled from '@emotion/styled';
import Highlight, { defaultProps } from 'prism-react-renderer';
import Editor from 'react-simple-code-editor';
import theme from 'prism-react-renderer/themes/github';
import * as React from 'react';
import { colors, soloConstants } from 'Styles';

type ContainerProps = { whiteBacked?: boolean };
const Container = styled.div`
  position: relative;
  padding: ${soloConstants.smallBuffer}px 5px ${soloConstants.smallBuffer}px
    calc(2em + ${soloConstants.smallBuffer}px);
  border: 1px solid ${colors.mayGrey};
  border-radius: ${soloConstants.smallRadius}px;
  background: ${(props: ContainerProps) =>
    props.whiteBacked ? 'white' : colors.januaryGrey};
  .prism-code {
    background: ${(props: ContainerProps) =>
      props.whiteBacked ? 'white' : colors.januaryGrey} !important;
  }

  > div > pre {
    overflow: visible;
  }
`;

export const Pre = styled.pre`
  text-align: left;

  & .token-line {
    line-height: 1.3em;
    height: 1.3em;
  }
`;

type LineNoProps = {
  edited?: boolean;
  editable?: boolean;
};
export const LineNo = styled.span`
  position: absolute;
  left: ${(props: LineNoProps) =>
    props.editable ? '-2em' : `${soloConstants.smallBuffer}px`};
  display: inline-block;
  width: 2em;
  user-select: none;
  opacity: 0.3;
  pointer-events: none;

  ${(props: LineNoProps) =>
    props.edited
      ? `background: ${colors.novemberGrey}; color: ${colors.januaryGrey};`
      : ''};
`;

const ourTheme = {
  ...theme,
  backgroundColor: 'transparent',
  overflow: 'initial'
};

const styles = {
  root: {
    ...ourTheme
  }
};

interface Props {
  content: string;
  isJson?: boolean;
  whiteBacked?: boolean;
  asEditor?: boolean;
}

export const ConfigDisplayer = React.memo((props: Props) => {
  const [editingContent, setEditingContent] = React.useState(props.content);

  React.useEffect(() => {
    if (editingContent !== props.content) {
      setEditingContent(props.content);
    }
  }, [props.content]);

  const onContentChange = (code: string): void => {
    setEditingContent(code);
  };

  const highlight = (code: string): React.ReactNode => {
    const originalContentLines = props.content.split('\n');

    return (
      <Highlight
        {...defaultProps}
        theme={ourTheme}
        code={code}
        language={props.isJson ? 'json' : 'yaml'}>
        {({ className, style, tokens, getLineProps, getTokenProps }) =>
          props.asEditor ? (
            <React.Fragment>
              {tokens.map((line, i) => {
                return (
                  <div {...getLineProps({ line, key: i })}>
                    <LineNo
                      editable
                      edited={
                        line.map(line => line.content).join('') !==
                        originalContentLines[i]
                      }>
                      {i + 1}
                    </LineNo>
                    {line.map((token, key) => (
                      <span {...getTokenProps({ token, key })} />
                    ))}
                  </div>
                );
              })}
            </React.Fragment>
          ) : (
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
          )
        }
      </Highlight>
    );
  };

  return (
    <Container whiteBacked={props.whiteBacked}>
      {props.asEditor ? (
        <Editor
          value={editingContent}
          onValueChange={onContentChange}
          highlight={highlight}
          padding={10}
          style={styles.root}
        />
      ) : (
        highlight(props.content)
      )}
    </Container>
  );
});
