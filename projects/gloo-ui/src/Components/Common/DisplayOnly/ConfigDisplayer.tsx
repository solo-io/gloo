import styled from '@emotion/styled';
import { Spin } from 'antd';
import { ReactComponent as EditIcon } from 'assets/edit-pencil.svg';
import Highlight, { defaultProps } from 'prism-react-renderer';
import theme from 'prism-react-renderer/themes/github';
import * as React from 'react';
import Editor from 'react-simple-code-editor';
import { colors, soloConstants } from 'Styles';
import { SoloCancelButton } from 'Styles/CommonEmotions/button';
import { SoloButton } from '../SoloButton';

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

  > div {
    > textarea {
      outline: none !important;
      border: 1px solid ${colors.mayGrey} !important;
    }
    > pre {
      overflow: visible;
    }
  }
`;

type EditPencilHolderProps = { inEditingMode?: boolean };
const EditPencilHolder = styled.div`
  ${(props: EditPencilHolderProps) =>
    props.inEditingMode
      ? `background: ${colors.septemberGrey}; 
      .aEditPencil { fill: ${colors.januaryGrey}; }`
      : `background: transparent; 
      .aEditPencil { fill: ${colors.septemberGrey}; }`};
  position: absolute;
  top: 0;
  right: 0;
  text-align: center;
  padding: 8px 8px 4px;
  cursor: pointer;
  z-index: 2;
`;
const EditPencil = styled(EditIcon)`
  width: 20px;
  height: 20px;
`;

const EditingActionsContainer = styled.div`
  position: absolute;
  top: 0;
  right: 0;
  padding: 8px;
  background: white;
  border-top-right-radius: ${soloConstants.smallRadius}px;
  z-index: 2;

  background: ${(props: ContainerProps) =>
    props.whiteBacked ? 'white' : colors.januaryGrey};
`;

const CancelButton = styled(SoloCancelButton)`
  margin-right: 8px;
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
    ...ourTheme,
    fontFamily:
      "'SFMono-Regular', Consolas, 'Liberation Mono', Menlo, Courier, monospace"
  }
};

type ConfigDisplayState = {
  isError: boolean;
  isLoading: boolean;
  isEditing: boolean;
};
let initialState = {
  isError: false,
  isLoading: false,
  isEditing: false
};

type ConfigDisplayAction =
  | { type: 'EDIT_MODE' }
  | { type: 'CHANGE_START' }
  | { type: 'CHANGE_SUCCESS' }
  | { type: 'CHANGE_ERROR' };

export function configDisplayReducer(
  state: ConfigDisplayState,
  action: ConfigDisplayAction
): ConfigDisplayState {
  switch (action.type) {
    case 'EDIT_MODE':
      return {
        ...state,
        isError: false,
        isEditing: !state.isEditing
      };
    case 'CHANGE_START':
      return {
        ...state,
        isLoading: true
      };
    case 'CHANGE_SUCCESS':
      return {
        ...state,
        isLoading: false,
        isError: false,
        isEditing: false
      };
    case 'CHANGE_ERROR':
      return {
        ...state,
        isLoading: false,
        isError: true,
        isEditing: true
      };
    default:
      return state;
  }
}
interface Props {
  content: string;
  isJson?: boolean;
  whiteBacked?: boolean;
  asEditor?: boolean;
  saveEdits?: (newContent: string) => void;
  yamlError?: boolean;
}

export const ConfigDisplayer = React.memo((props: Props) => {
  const [configState, configDispatch] = React.useReducer(
    configDisplayReducer,
    initialState
  );

  const [editingContent, setEditingContent] = React.useState(props.content);

  React.useEffect(() => {
    if (editingContent !== props.content) {
      setEditingContent(props.content);
    }
    configDispatch({ type: 'CHANGE_SUCCESS' });
  }, [props.content]);

  React.useEffect(() => {
    if (props.yamlError == true) {
      configDispatch({ type: 'CHANGE_ERROR' });
    }
  }, [props.yamlError]);

  const onContentChange = (code: string): void => {
    setEditingContent(code);
  };

  const saveEdits = () => {
    configDispatch({ type: 'CHANGE_START' });
    if (props.saveEdits) {
      props.saveEdits(editingContent);

      if (configState.isError === false) {
        configDispatch({ type: 'EDIT_MODE' });
      }
    }
  };

  const cancelEdits = (): void => {
    setEditingContent(props.content);
    configDispatch({ type: 'EDIT_MODE' });
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
          props.asEditor && configState.isEditing ? (
            <>
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
            </>
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
    <Spin spinning={configState.isLoading}>
      <Container whiteBacked={props.whiteBacked}>
        {props.asEditor && (
          <>
            {configState.isEditing ? (
              <EditingActionsContainer whiteBacked={props.whiteBacked}>
                <CancelButton onClick={cancelEdits}>Reset</CancelButton>
                <SoloButton
                  text='Apply'
                  onClick={saveEdits}
                  disabled={props.content === editingContent}></SoloButton>
              </EditingActionsContainer>
            ) : (
              <EditPencilHolder
                inEditingMode={configState.isEditing}
                onClick={() => configDispatch({ type: 'EDIT_MODE' })}>
                <EditPencil />
              </EditPencilHolder>
            )}
          </>
        )}
        {props.asEditor && configState.isEditing ? (
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
    </Spin>
  );
});
