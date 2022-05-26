import styled from '@emotion/styled/macro';
import { useGetConsoleOptions } from 'API/hooks';
import { SoloFormDropdown } from 'Components/Common/SoloFormComponents';
import VisualEditor from 'Components/Common/VisualEditor';
import { useFormikContext } from 'formik';
import * as React from 'react';
import { useEffect, useMemo, useState } from 'react';
import {
  SoloButtonStyledComponent,
  SoloCancelButton,
} from 'Styles/StyledComponents/button';
import { Spacer } from 'Styles/StyledComponents/spacer';
import YAML from 'yaml';
import WarningMessage from '../../executable-api/WarningMessage';
import { getResolverFromConfig } from './converters';
import { ResolverWizardFormProps } from './ResolverWizard';
import * as ResolverWizardStyles from './ResolverWizard.styles';

export const EditorContainer = styled.div<{ editMode: boolean }>`
  .ace_cursor {
    opacity: ${props => (props.editMode ? 1 : 0)};
  }
  cursor: ${props => (props.editMode ? 'text' : 'default')};
`;

export const getDefaultConfigFromType = (
  resolverType: ResolverWizardFormProps['resolverType']
) => {
  YAML.scalarOptions.null.nullStr = '';
  if (resolverType === 'gRPC')
    return YAML.stringify(
      YAML.parse(`requestTransform:
         serviceName:
         methodName:
         requestMetadata:
         outgoingMessageJson:
  `),
      { simpleKeys: true }
    );
  if (resolverType === 'Mock')
    return YAML.stringify(
      YAML.parse(`syncResponse:
  `),
      { simpleKeys: true }
    );
  // Default: resolverType==='REST'
  return YAML.stringify(
    YAML.parse(`
          request:
            headers:
              :method:
              :path:
            queryParams:
            body:
          response:
            resultRoot:
            setters:
  `),
    { simpleKeys: true }
  );
};

export const ResolverConfigSection: React.FC<{
  onCancel(): void;
  submitDisabled: boolean;
  warningMessage: string;
}> = ({ onCancel, submitDisabled, warningMessage }) => {
  const { readonly } = useGetConsoleOptions();
  const { setFieldValue, values, dirty, handleSubmit } =
    useFormikContext<ResolverWizardFormProps>();
  const [errorMessage, setErrorMessage] = React.useState('');

  useEffect(() => {
    setErrorMessage(warningMessage);
  }, [warningMessage, setErrorMessage]);

  const [selectedName, setSelectedName] = useState<string>();
  const resolverOptions = useMemo(
    () =>
      values.listOfResolvers
        .filter(([_rName, rObj]) => {
          let type = '';
          if (rObj.grpcResolver) {
            type = 'gRPC';
          } else if (rObj.restResolver) {
            type = 'REST';
          }
          return type === values.resolverType;
        })
        .map(([rName]) => {
          return {
            key: rName,
            value: rName,
          };
        }),
    [values.listOfResolvers, values.resolverType]
  );

  const onResolverCopy = (copyName: any) => {
    const resolver = values.listOfResolvers.find(([rName]) => {
      return rName === copyName;
    });
    if (resolver) {
      const [_rName, newResolver] = resolver;
      setSelectedName(_rName);
      const stringifiedResolver = getResolverFromConfig(newResolver);
      setFieldValue('resolverConfig', stringifiedResolver);
    }
  };

  const resetResolverConfig = () => {
    setFieldValue(
      'resolverConfig',
      getDefaultConfigFromType(values.resolverType)
    );
    setErrorMessage('');
  };

  return (
    <>
      <div data-testid='resolver-config-section'>
        <Spacer mb={2}>
          <div>
            <EditorContainer editMode={true}>
              <Spacer px={6} mt='-20px'>
                {!!warningMessage ? (
                  <WarningMessage message={warningMessage} />
                ) : (
                  <Spacer pt='1.5rem' />
                )}
              </Spacer>
              <div className='flex flex-col w-full'>
                <Spacer my={3} mx={6}>
                  {resolverOptions.length > 0 && (
                    <div
                      data-testid='create-resolver-from-config'
                      className='grid grid-cols-2 gap-4 '>
                      <div>
                        <SoloFormDropdown
                          searchable={true}
                          name='resolverCopy'
                          title='Create Resolver From Config'
                          value={selectedName}
                          onChange={onResolverCopy}
                          options={resolverOptions}
                        />
                      </div>
                    </div>
                  )}
                </Spacer>
                <VisualEditor
                  data-testid='resolve-config-editor'
                  mode='yaml'
                  theme='chrome'
                  name='resolverConfiguration'
                  style={{
                    width: '100%',
                    height: '35vh',
                    maxHeight: '450px',
                    minHeight: '350px',
                    cursor: 'text',
                  }}
                  defaultValue={values.resolverConfig || ''}
                  value={values.resolverConfig}
                  onChange={(newValue, e) => {
                    setFieldValue('resolverConfig', newValue);
                  }}
                  focus={true}
                  fontSize={16}
                  showPrintMargin={false}
                  showGutter={true}
                  highlightActiveLine={true}
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
                <Spacer mt={2} px={3}>
                  <SoloCancelButton
                    disabled={!dirty}
                    onClick={resetResolverConfig}>
                    Reset
                  </SoloCancelButton>
                </Spacer>
              </div>
            </EditorContainer>
          </div>
        </Spacer>
      </div>
      <Spacer px={6} className='flex items-center justify-between'>
        <ResolverWizardStyles.IconButton onClick={onCancel}>
          Cancel
        </ResolverWizardStyles.IconButton>
        {!readonly && (
          <SoloButtonStyledComponent
            data-testid='resolver-wizard-submit'
            onClick={handleSubmit as any}
            disabled={submitDisabled}>
            Submit
          </SoloButtonStyledComponent>
        )}
      </Spacer>
    </>
  );
};
