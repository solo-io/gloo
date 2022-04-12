import VisualEditor from 'Components/Common/VisualEditor';
import * as React from 'react';
import styled from '@emotion/styled/macro';

const StyledContainer = styled.div`
  margin-top: 20px;
  margin-bottom: 20px;
`;

interface UpdateApiEditorProps {
  setGraphqlSchema: (value: string) => any;
  graphqlSchema: string;
}

export const UpdateApiEditor = (props: UpdateApiEditorProps) => {
  const { setGraphqlSchema, graphqlSchema } = props;

  return (
    <>
      <StyledContainer>
        <VisualEditor
          theme='chrome'
          name='graphqlEditor'
          style={{
            width: '100%',
            maxHeight: '36vh',
            cursor: 'text',
          }}
          onChange={(newValue, _e) => {
            setGraphqlSchema(newValue);
            // Change values.
          }}
          focus={true}
          fontSize={16}
          showPrintMargin={false}
          showGutter={true}
          highlightActiveLine={true}
          defaultValue={graphqlSchema || ''}
          value={graphqlSchema}
          readOnly={false}
          setOptions={{
            highlightGutterLine: true,
            showGutter: true,
            fontFamily: 'monospace',
            enableBasicAutocompletion: true,
            enableLiveAutocompletion: true,
            showLineNumbers: true,
            tabSize: 2,
          }}
          mode='graphqlschema'
        />
      </StyledContainer>
    </>
  );
};
