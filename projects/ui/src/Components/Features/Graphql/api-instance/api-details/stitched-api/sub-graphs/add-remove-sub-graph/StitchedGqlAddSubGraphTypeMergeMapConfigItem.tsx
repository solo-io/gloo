import VisualEditor from 'Components/Common/VisualEditor';
import React from 'react';
import { supportedDefinitionTypes } from 'utils/graphql-helpers';

const StitchedGqlAddSubGraphTypeMergeConfigItem: React.FC<{
  schemaDefinitions: supportedDefinitionTypes[];
  typeMergeConfig: string;
  onTypeMergeConfigChange(newConfig: string): void;
}> = ({ typeMergeConfig, onTypeMergeConfigChange }) => {
  return (
    <VisualEditor
      mode='yaml'
      theme='chrome'
      name='resolverConfiguration'
      style={{
        width: '100%',
        height: '200px',
      }}
      onChange={newValue => onTypeMergeConfigChange(newValue)}
      fontSize={16}
      showPrintMargin={false}
      showGutter={true}
      highlightActiveLine={true}
      value={typeMergeConfig}
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
    />
  );
};

export default StitchedGqlAddSubGraphTypeMergeConfigItem;
