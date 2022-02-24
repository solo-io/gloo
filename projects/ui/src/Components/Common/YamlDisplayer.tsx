import React, { ReactNode, useState } from 'react';
import styled from '@emotion/styled';
import { colors } from 'Styles/colors';
import Highlight, { defaultProps } from 'prism-react-renderer';
import theme from 'prism-react-renderer/themes/github';
import { ReactComponent as DocumentIcon } from 'assets/document.svg';
import { copyTextToClipboard } from 'utils';

const Container = styled.div`
  position: relative;
  padding: 0 5px 0 calc(2em + 18px);
  border: 1px solid ${colors.mayGrey};
  border-radius: 8px;
  background: ${colors.januaryGrey};
  max-height: 1000px;
  overflow-y: scroll;

  .prism-code {
    background: ${colors.januaryGrey} !important;
  }

  > div {
    > textarea {
      outline: none !important;

      border: 1px solid ${colors.mayGrey} !important;
    }
    > pre {
      overflow-x: auto;
      overflow-wrap: normal;
    }
  }
`;

const DescriptionContainer = styled.div`
  color: ${colors.septemberGrey};
  font-size: 14px;
  margin-left: -2em;
  margin-top: 18px;
`;

const YamlContainer = styled.div`
  margin-top: 10px;
`;

const Pre = styled.pre`
  text-align: left;

  & .token-line {
    line-height: 1.3em;
    height: 1.3em;
  }

  > div {
    > span {
      font-family: monospace;
    }
  }
`;

const LineNo = styled.span`
  position: absolute;
  left: 18px;
  display: inline-block;
  width: 2em;
  user-select: none;
  opacity: 0.3;
  pointer-events: none;
`;

const displayTheme = {
  ...theme,
  backgroundColor: 'transparent',
  overflow: 'initial',
};

type CopyButtonProps = {
  copySuccessful: boolean | 'inactive';
  oneLine: boolean;
};
const CopyButton = styled.div<CopyButtonProps>`
  position: absolute;

  ${(props: CopyButtonProps) =>
    props.oneLine
      ? `top: 0;
        right: 0;
        height: 44px;
        width: 44px;`
      : `top: 13px;
        right: 13px;
        height: 36px;
        width: 36px;`}

  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  border-radius: 8px;
  transition: background 0.7s ease-out;

  ${(props: CopyButtonProps) => {
    if (props.copySuccessful === 'inactive') {
      return `background: ${colors.seaBlue};`;
    } else {
      return `background: ${
        props.copySuccessful ? colors.forestGreen : colors.pumpkinOrange
      };
      transition: background 0.2s ease-in;`;
    }
  }}

  svg {
    width: 16px;

    * {
      fill: white;
    }
  }
`;

type Props = {
  description?: ReactNode;
  contentString: string;
  copyable?: boolean;
};

const YamlDisplayer = ({ contentString, description, copyable }: Props) => {
  const [attemptedCopy, setAttemptedCopy] = useState<boolean | 'inactive'>(
    'inactive'
  );

  const attemptCopyToClipboard = () => {
    copyTextToClipboard(contentString)
      .then(() => {
        setAttemptedCopy(true);
        setTimeout(() => {
          setAttemptedCopy('inactive');
        }, 500);
      })
      .catch(() => {
        setAttemptedCopy(false);
      });
  };

  return (
    <Container className='YamlDisplayerContainer'>
      {description && (
        <DescriptionContainer>{description}</DescriptionContainer>
      )}

      <YamlContainer>
        <Highlight
          {...defaultProps}
          theme={displayTheme}
          code={contentString}
          language='yaml'>
          {({ className, style, tokens, getLineProps, getTokenProps }) => (
            <>
              {copyable && tokens.length > 0 && (
                <CopyButton
                  copySuccessful={attemptedCopy}
                  oneLine={tokens.length === 1}
                  onClick={attemptCopyToClipboard}>
                  <DocumentIcon />
                </CopyButton>
              )}
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
            </>
          )}
        </Highlight>
      </YamlContainer>
    </Container>
  );
};

export default YamlDisplayer;
