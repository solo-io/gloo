import VisualEditor from 'Components/Common/VisualEditor';
import React, { useEffect, useState } from 'react';
import WarningMessage from '../../../executable-api/WarningMessage';
import { typeMergeConfigChangedFromDefault } from './StitchedGqlTypeMergeMapConfig';
import { getPreMarshalledConfig } from './StitchedGqlTypeMergeMapHelpers';

const StitchedGqlTypeMergeConfigItem: React.FC<{
  typeMergeConfig: string;
  onTypeMergeConfigChange(newConfig: string): void;
}> = ({ typeMergeConfig, onTypeMergeConfigChange }) => {
  const [warningMessage, setWarningMessage] = useState('');
  useEffect(() => {
    try {
      getPreMarshalledConfig(typeMergeConfig);
      setWarningMessage('');
    } catch (e: any) {
      const newWarningMessage =
        !typeMergeConfigChangedFromDefault(typeMergeConfig) ||
        !typeMergeConfig?.trim()
          ? ''
          : e?.message ?? e;

      setWarningMessage(newWarningMessage);
    }
  }, [typeMergeConfig, setWarningMessage]);

  return (
    <>
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
      <WarningMessage className='pt-2 pb-2' message={warningMessage} />
    </>
  );
};

export default StitchedGqlTypeMergeConfigItem;
