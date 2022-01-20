import * as React from 'react';
import styled from '@emotion/styled/macro';
import AceEditor, { IAceEditorProps } from 'react-ace';
/* 
  These imports are needed for syntax highlighting and snippets. DO NOT REMOVE.
*/
import 'ace-builds/src-noconflict/ext-language_tools';
import 'ace-builds/src-noconflict/ext-searchbox';
import 'ace-builds/src-noconflict/mode-yaml';
import 'ace-builds/src-noconflict/mode-html';
import 'ace-builds/src-noconflict/snippets/yaml';
import 'ace-builds/src-noconflict/theme-chrome';
import 'ace-builds/webpack-resolver';
import { colors } from 'Styles/colors';

export const Label = styled.label`
  display: block;
  color: ${colors.novemberGrey};
  font-size: 16px;
  margin-bottom: 10px;
  font-weight: 500;
`;

export interface SoloFormYamlEditorProps extends IAceEditorProps {
  name: string; // the name of this field in Formik
  title?: string; // display name of the field
}

const YamlEditor = (props: SoloFormYamlEditorProps) => {
  const { name, title, ...rest } = props;

  return (
    <div>
      {title && <Label>{title}</Label>}

      <AceEditor
        mode='yaml'
        theme='chrome'
        name='gatewayConfiguration'
        style={{
          maxWidth: '40vw',
          maxHeight: '25vh',
          cursor: 'text',
        }}
        onChange={rest.onChange}
        focus={true}
        onInput={rest.onInput}
        fontSize={14}
        showPrintMargin={false}
        showGutter={true}
        highlightActiveLine={true}
        value={rest.value}
        readOnly={false}
        setOptions={{
          highlightGutterLine: true,
          showGutter: true,
          enableBasicAutocompletion: true,
          enableLiveAutocompletion: true,
          showLineNumbers: true,
          tabSize: 2,
        }}
        {...rest}
      />
    </div>
  );
};

export { YamlEditor as default };
