import * as React from 'react';
import styled from '@emotion/styled/macro';
import AceEditor, { IAceEditorProps } from 'react-ace';
import { di } from 'react-magnetic-di/macro';
/*
  These imports are needed for syntax highlighting and snippets. DO NOT REMOVE.
*/
import 'ace-builds/src-noconflict/ext-language_tools';
import 'ace-builds/src-noconflict/ext-searchbox';
import 'ace-builds/src-noconflict/mode-yaml';
import 'ace-builds/src-noconflict/mode-html';
import 'ace-builds/src-noconflict/mode-graphqlschema';
import 'ace-builds/src-noconflict/snippets/yaml';
import 'ace-builds/src-noconflict/snippets/graphqlschema';
import 'ace-builds/src-noconflict/theme-chrome';
import 'ace-builds/webpack-resolver';
import { colors } from 'Styles/colors';
import { SoloDropdown } from './SoloDropdown';
import { useAppSettings } from 'Components/Context/AppSettingsContext';

export const Label = styled.label`
  display: block;
  color: ${colors.novemberGrey};
  font-size: 16px;
  margin-bottom: 10px;
  font-weight: 500;
`;

const StyledVisualEditor = styled.div`
  position: relative;
  border: 1px solid ${colors.aprilGrey};
  .ant-select {
    opacity: 0;
    transition: 0.1s opacity;
  }
  &:focus-within .ant-select,
  &:hover .ant-select {
    opacity: 1;
  }
`;

const StyledEditorSettings = styled.div`
  position: absolute;
  right: 0px;
  bottom: 0px;
`;

const StyledAceEditor = styled(AceEditor)`
  .ace_editor span,
  .ace_editor textarea {
    font-size: 16px;
    font-family: 'monospace';
  }
`;

export interface SoloFormVisualEditorProps extends IAceEditorProps {
  name: string; // the name of this field in Formik
  title?: string; // display name of the field
}

const VisualEditor = (props: SoloFormVisualEditorProps) => {
  di(useAppSettings);
  const { name, title, value, ...rest } = props;
  const { appSettings, onAppSettingsChange } = useAppSettings();
  const { keyboardHandler } = appSettings;

  return (
    <StyledVisualEditor>
      {title && <Label>{title}</Label>}

      <StyledAceEditor
        keyboardHandler={keyboardHandler}
        mode={rest.mode ?? 'yaml'}
        theme='chrome'
        name={name ?? title}
        style={{
          maxWidth: '40vw',
          maxHeight: '25vh',
          // cursor: 'text',
        }}
        onChange={rest.onChange}
        focus={true}
        onInput={rest.onInput}
        fontSize={14}
        showPrintMargin={false}
        showGutter={true}
        highlightActiveLine={true}
        value={value}
        readOnly={false}
        setOptions={{
          highlightGutterLine: true,
          showGutter: true,
          cursorStyle: 'wide',
          fontFamily: 'monospace',
          enableBasicAutocompletion: true,
          enableLiveAutocompletion: true,
          showLineNumbers: true,
          tabSize: 2,
        }}
        {...rest}
      />

      <StyledEditorSettings>
        <SoloDropdown
          value={keyboardHandler}
          options={[
            {
              key: 'keyboard',
              value: 'keyboard',
              displayValue: 'Editor: Standard',
            },
            { key: 'vim', value: 'vim', displayValue: 'Editor: Vim' },
            { key: 'emacs', value: 'emacs', displayValue: 'Editor: Emacs' },
          ]}
          onChange={newKeyboardHandler => {
            onAppSettingsChange({
              ...appSettings,
              keyboardHandler: newKeyboardHandler as string,
            });
          }}
        />
      </StyledEditorSettings>
    </StyledVisualEditor>
  );
};

export { VisualEditor as default };
