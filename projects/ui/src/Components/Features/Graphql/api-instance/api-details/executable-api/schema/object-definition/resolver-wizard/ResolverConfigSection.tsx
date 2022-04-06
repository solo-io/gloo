import * as React from 'react';
import { ResolverWizardFormProps } from './ResolverWizard';
import { useFormikContext } from 'formik';
import YAML from 'yaml';
import styled from '@emotion/styled/macro';
import VisualEditor from 'Components/Common/VisualEditor';
import { SoloFormDropdown } from 'Components/Common/SoloFormComponents';
import { SoloCancelButton } from 'Styles/StyledComponents/button';
import { useParams } from 'react-router';
import { getResolverFromConfig } from './converters';

export const EditorContainer = styled.div<{ editMode: boolean }>`
  .ace_cursor {
    opacity: ${props => (props.editMode ? 1 : 0)};
  }
  cursor: ${props => (props.editMode ? 'text' : 'default')};
`;

const SpacedValuesContainer = styled.div`
  margin: 10px 0;
`;

const InnerValues = styled.div`
  margin: 5px 0;
`;

type ResolverConfigSectionProps = {
  warningMessage: string;
};

export const getDefaultConfigFromType = (
  resolverType: ResolverWizardFormProps['resolverType']
) => {
  YAML.scalarOptions.null.nullStr = '';
  return resolverType === 'REST'
    ? YAML.stringify(
        YAML.parse(`
          request:
            headers:
              :method:
              :path:
            queryParams:
            body:
          response:
            resultRoot:
            setters:`),
        {
          simpleKeys: true,
        }
      )
    : YAML.stringify(
        YAML.parse(`requestTransform:
         serviceName:
         methodName:
         requestMetadata:
         outgoingMessageJson:
  `),
        {
          simpleKeys: true,
        }
      );
};

export const ResolverConfigSection = ({
  warningMessage,
}: ResolverConfigSectionProps) => {
  const { setFieldValue, values, dirty, errors } =
    useFormikContext<ResolverWizardFormProps>();
  const [isValid, setIsValid] = React.useState(false);
  const [errorMessage, setErrorMessage] = React.useState('');
  const [demoConfig, setDemoConfig] = React.useState('');

  const { graphqlApiName = '', graphqlApiNamespace = '' } = useParams();

  React.useEffect(() => {
    setErrorMessage(warningMessage);
  }, [warningMessage, setErrorMessage]);

  const [selectedName, setSelectedName] = React.useState<string>();
  const resolverOptions = values.listOfResolvers
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
    });

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

  React.useEffect(() => {
    const newDemo = getDefaultConfigFromType(values.resolverType);
    setDemoConfig(newDemo);
  }, [graphqlApiName, graphqlApiNamespace, values.resolverType]);

  return (
    <div data-testid='resolver-config-section' className='px-6 pb-0 '>
      <div className='mb-2 '>
        <div>
          <EditorContainer editMode={true}>
            <div className=''>
              <div className='' style={{ height: 'min-content' }}>
                {isValid && errorMessage.length === 0 ? (
                  <div
                    className={`${
                      isValid ? 'opacity-100' : 'opacity-0'
                    } h-10 text-center`}>
                    <div
                      style={{ backgroundColor: '#f2fef2' }}
                      className='p-2 text-green-400 border border-green-400 '>
                      <div className='font-medium '>Valid</div>
                    </div>
                  </div>
                ) : (
                  <div
                    className={`${
                      errorMessage.length > 0 ? 'opacity-100' : '  opacity-0'
                    }`}>
                    <div
                      style={{ backgroundColor: '#FEF2F2' }}
                      className='p-2 text-orange-400 border border-orange-400 mb-5'>
                      <div className='font-medium '>
                        {errorMessage?.split(',')[0]}
                      </div>
                      <ul className='pl-2 list-disc'>
                        {errorMessage?.split(',')[1]}
                      </ul>
                    </div>
                  </div>
                )}
              </div>
            </div>
            <div className='flex flex-col w-full '>
              <>
                <SpacedValuesContainer>
                  <InnerValues className='p-2'>
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
                  </InnerValues>
                </SpacedValuesContainer>
                <VisualEditor
                  mode='yaml'
                  theme='chrome'
                  name='resolverConfiguration'
                  style={{
                    width: '100%',
                    height: '30vh',
                    maxHeight: '400px',
                    minHeight: '200px',
                    cursor: 'text',
                  }}
                  onChange={(newValue, e) => {
                    setFieldValue('resolverConfig', newValue);
                  }}
                  focus={true}
                  onInput={() => {
                    setIsValid(false);
                  }}
                  fontSize={16}
                  showPrintMargin={false}
                  showGutter={true}
                  highlightActiveLine={true}
                  defaultValue={values.resolverConfig || ''}
                  value={values.resolverConfig}
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
                <div className='flex gap-3 mt-2'>
                  <SoloCancelButton
                    disabled={!dirty}
                    onClick={() => {
                      setFieldValue('resolverConfig', demoConfig);
                      setErrorMessage('');
                    }}>
                    Reset
                  </SoloCancelButton>
                </div>
              </>
            </div>
          </EditorContainer>
        </div>
      </div>
    </div>
  );
};
