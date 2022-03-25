import * as React from 'react';
import { ResolverWizardFormProps } from './ResolverWizard';
import { useFormikContext } from 'formik';
import YAML from 'yaml';
import { Resolution } from 'proto/github.com/solo-io/solo-apis/api/gloo/graphql.gloo/v1alpha1/graphql_pb';
import { graphqlConfigApi } from 'API/graphql';
import { ValidateResolverYamlRequest } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/graphql_pb';
import styled from '@emotion/styled/macro';
import VisualEditor from 'Components/Common/VisualEditor';
import {
  SoloButtonStyledComponent,
  SoloCancelButton,
} from 'Styles/StyledComponents/button';
import { useParams } from 'react-router';
import isEqual from 'lodash/isEqual';

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
  isEdit: boolean;
  warningMessage: string;
};

export const getDefaultConfigFromType = (
  name: string,
  namespace: string,
  resolverType: ResolverWizardFormProps['resolverType']
) => {
  return resolverType === 'REST'
    ? `
  restResolver:
    upstreamRef:
      name: ${name}
      namespace: ${namespace}
    request:
      headers:
        :method: GET
        :path: /api/v1/products
      queryParams:
      body:
    response:
      resultRoot: author
      setters:
        numStars: '5'
        reviewer: '1'
    `.trimEnd()
    : `grpcResolver:
    upstreamRef:
      name: ${name}
      namespace: ${namespace}
    requestTransform:
      serviceName: my-service
      methodName: my-method
    spanName: hello
  `.trimEnd();
};

export const ResolverConfigSection = ({
  isEdit,
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

  const validateResolverSchema = async (resolver: string) => {
    const resolverObj = YAML.parse(resolver);
    if (!resolverObj) {
      setIsValid(true);
      return;
    }
    delete resolverObj?.restResolver?.request?.headersMap;
    delete resolverObj?.restResolver?.request?.queryParamsMap;
    delete resolverObj?.restResolver?.response?.settersMap;
    const resolverType =
      values.resolverType === 'REST'
        ? ValidateResolverYamlRequest.ResolverType.REST_RESOLVER
        : ValidateResolverYamlRequest.ResolverType.GRPC_RESOLVER;
    let parsed = {};
    if (
      resolverType === ValidateResolverYamlRequest.ResolverType.REST_RESOLVER
    ) {
      parsed = resolverObj.restResolver;
    } else {
      parsed = resolverObj.grpcResolver;
    }
    YAML.scalarOptions.null.nullStr = '';
    const yaml = YAML.stringify(parsed);
    try {
      await graphqlConfigApi
        .validateResolverYaml({
          yaml,
          resolverType,
        })
        .then(resp => {
          setIsValid(true);
        })
        .catch(err => {
          setErrorMessage(err.message);
          setIsValid(false);
        });
    } catch (err: any) {
      let [_, conversionError] = err?.message?.split(
        'failed to convert options YAML to JSON: yaml:'
      ) as [string, string];
      let [__, yamlError] = err?.message?.split(' invalid options YAML:') as [
        string,
        string
      ];
      if (conversionError) {
        setIsValid(false);
        setErrorMessage(`Error on ${conversionError}`);
      } else if (yamlError) {
        setIsValid(false);
        setErrorMessage(
          `Error: ${yamlError?.substring(yamlError.indexOf('):') + 2) ?? ''}`
        );
      }
    }
  };
  React.useEffect(() => {
    const newDemo = getDefaultConfigFromType(
      graphqlApiName,
      graphqlApiNamespace,
      values.resolverType
    );
    setDemoConfig(newDemo);
  }, [graphqlApiName, graphqlApiNamespace, values.resolverType]);

  return (
    <div data-testid='resolver-config-section' className='h-full p-6 pb-0 '>
      <div
        className={'flex items-center mb-2 text-lg font-medium text-gray-800'}>
        {isEdit ? 'Edit' : 'Configure'} Resolver{' '}
      </div>

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

              <div className='mt-2'></div>
            </div>
            <div className='flex flex-col w-full '>
              <>
                <SpacedValuesContainer>
                  <InnerValues className='p-2 border'>
                    <p>
                      <b>
                        Editing the resolver type or upstream values in this
                        editor will have no effect.
                      </b>
                    </p>
                    <p>
                      <b>
                        Please go back to the appropriate step on the wizard to
                        make a change to these values.
                      </b>
                    </p>
                  </InnerValues>
                </SpacedValuesContainer>
                <VisualEditor
                  mode='yaml'
                  theme='chrome'
                  name='resolverConfiguration'
                  style={{
                    width: '100%',
                    maxHeight: '36vh',
                    cursor: 'text',
                  }}
                  onChange={(newValue, e) => {
                    // TODO:  Could add in a check here for a change to those values.
                    try {
                      const resolverObj = YAML.parse(newValue);
                      if (
                        values?.resolverType === 'REST' &&
                        !resolverObj?.restResolver
                      ) {
                        setErrorMessage('Cannot edit the restResolver here.');
                        e.preventDefault();
                        return;
                      }
                      if (
                        values?.resolverType === 'gRPC' &&
                        !resolverObj?.grpcResolver
                      ) {
                        setErrorMessage('Cannot edit the grpcResolver here.');
                        e.preventDefault();
                        return;
                      }
                      let joinedName = '';
                      if (values.resolverType === 'gRPC') {
                        joinedName = `${resolverObj?.grpcResolver?.upstreamRef?.name}::${resolverObj?.grpcResolver?.upstreamRef?.namespace}`;
                      } else if (values.resolverType === 'REST') {
                        joinedName = `${resolverObj?.restResolver?.upstreamRef?.name}::${resolverObj?.restResolver?.upstreamRef?.namespace}`;
                      }
                      if (joinedName !== values.upstream) {
                        setErrorMessage(
                          'Cannot edit the upstream references here.'
                        );
                        return;
                      }
                      setErrorMessage('');
                    } catch (err: any) {
                      console.error('error on parse change', err);
                    }
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
                  <SoloButtonStyledComponent
                    data-testid='save-route-options-changes-button '
                    disabled={!dirty || isValid}
                    onClick={() =>
                      validateResolverSchema(values.resolverConfig ?? '')
                    }>
                    Validate
                  </SoloButtonStyledComponent>

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
