import * as React from 'react';
import { ResolverWizardFormProps } from '../ResolverWizard';
import { useFormikContext } from 'formik';
import YAML from 'yaml';
import { Resolution } from 'proto/github.com/solo-io/solo-apis/api/gloo/graphql.gloo/v1alpha1/graphql_pb';
import { graphqlApi } from 'API/graphql';
import { ValidateResolverYamlRequest } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/graphql_pb';
import styled from '@emotion/styled/macro';
import YamlEditor from 'Components/Common/YamlEditor';
import {
  SoloButtonStyledComponent,
  SoloCancelButton,
} from 'Styles/StyledComponents/button';

export const EditorContainer = styled.div<{ editMode: boolean }>`
  .ace_cursor {
    opacity: ${props => (props.editMode ? 1 : 0)};
  }
  cursor: ${props => (props.editMode ? 'text' : 'default')};
`;
type ResolverConfigSectionProps = {
  isEdit: boolean;
  resolverConfig: string;
  existingResolverConfig?: Resolution.AsObject;
};

let demoConfig = `restResolver:
request:
    headers:
    :method:
    :path:
    queryParams:
    body:
response:
    resultRoot:
    setters:`;

export const ResolverConfigSection = ({
  isEdit,
  existingResolverConfig,
}: ResolverConfigSectionProps) => {
  const { setFieldValue, values, dirty, errors } =
    useFormikContext<ResolverWizardFormProps>();
  const [isValid, setIsValid] = React.useState(false);
  const [errorModal, setErrorModal] = React.useState(false);
  const [errorMessage, setErrorMessage] = React.useState('');

  // remove `null` and empty fields
  // change headersMap,
  // queryParamsMap
  // rsponse settersMap

  function resolverConfigToDisplay(
    resolverConfig: Resolution.AsObject
  ): string {
    let simpleConfig = { ...resolverConfig } as Partial<Resolution.AsObject>;

    if (values.resolverType === 'REST') {
      delete simpleConfig?.restResolver?.upstreamRef;
      delete simpleConfig?.grpcResolver;
      if (!simpleConfig?.restResolver?.spanName) {
        //@ts-ignore
        delete simpleConfig?.restResolver?.spanName;
      }

      if (
        Object.keys(simpleConfig?.restResolver?.request ?? {})?.length === 0 ||
        !simpleConfig?.restResolver?.request
      ) {
        delete simpleConfig?.restResolver?.request;
      } else {
        //     request?: {
        //  headersMap: Array<[string, string]>,
        //  queryParamsMap: Array<[string, string]>,
        //  body?: google_protobuf_struct_pb.Value.AsObject,
        // },
        if (
          simpleConfig?.restResolver?.request?.queryParamsMap?.length === 0 ||
          !simpleConfig?.restResolver?.request?.queryParamsMap
        ) {
          // @ts-ignore
          delete simpleConfig?.restResolver?.request?.queryParamsMap;
        } else {
          let qParams =
            Object.fromEntries(
              simpleConfig.restResolver?.request?.queryParamsMap
            ) ?? [];
          //@ts-ignore
          simpleConfig.restResolver.request = {
            ...simpleConfig?.restResolver?.request,
            //@ts-ignore
            qParams,
          };
          //@ts-ignore
          delete simpleConfig?.restResolver?.request?.queryParamsMap;
        }

        // headers
        if (
          simpleConfig?.restResolver?.request?.headersMap?.length === 0 ||
          !simpleConfig?.restResolver?.request?.headersMap
        ) {
          // @ts-ignore
          delete simpleConfig?.restResolver?.request?.queryParamsMap;
        } else {
          let headers =
            Object.fromEntries(
              simpleConfig.restResolver?.request?.headersMap
            ) ?? [];
          //@ts-ignore
          simpleConfig.restResolver.request = {
            ...simpleConfig?.restResolver?.request,
            //@ts-ignore
            headers,
          };
          //@ts-ignore
          delete simpleConfig?.restResolver?.request?.headersMap;
        }

        // body
        if (!simpleConfig?.restResolver?.request?.body) {
          // @ts-ignore
          delete simpleConfig?.restResolver?.request?.body;
        } else {
          // TODO: parse body
        }
      }

      if (
        Object.keys(simpleConfig?.restResolver?.response ?? {})?.length === 0 ||
        !simpleConfig?.restResolver?.response
      ) {
        delete simpleConfig?.restResolver?.response;
      }
    } else {
      delete simpleConfig?.restResolver;
      delete simpleConfig?.grpcResolver?.upstreamRef;
      if (!simpleConfig?.grpcResolver?.spanName) {
        //@ts-ignore
        delete simpleConfig?.grpcResolver?.spanName;
      }
      if (
        Object.keys(simpleConfig?.grpcResolver?.requestTransform ?? {})
          ?.length === 0 ||
        !simpleConfig?.grpcResolver?.requestTransform
      ) {
        delete simpleConfig?.grpcResolver?.requestTransform;
      } else {
        // outgoingMessageJson?: google_protobuf_struct_pb.Value.AsObject,
        // serviceName: string,
        // methodName: string,
        // requestMetadataMap: Array<[string, string]>
        if (
          simpleConfig?.grpcResolver?.requestTransform?.requestMetadataMap
            ?.length === 0 ||
          !simpleConfig?.grpcResolver?.requestTransform?.requestMetadataMap
        ) {
          // @ts-ignore
          delete simpleConfig?.grpcResolver?.requestTransform
            ?.requestMetadataMap;
        } else {
          let requestMetadata =
            Object.fromEntries(
              simpleConfig.grpcResolver?.requestTransform?.requestMetadataMap
            ) ?? [];
          //@ts-ignore
          simpleConfig.grpcResolver.requestTransform = {
            ...simpleConfig?.grpcResolver?.requestTransform,
            //@ts-ignore
            requestMetadata,
          };
          //@ts-ignore
          delete simpleConfig?.grpcResolver?.requestTransform
            ?.requestMetadataMap;
        }
      }
    }

    if (!simpleConfig?.statPrefix?.value) {
      delete simpleConfig?.statPrefix;
    }

    if (Object.keys(resolverConfig?.restResolver ?? {})?.length > 0) {
    } else if (Object.keys(resolverConfig?.grpcResolver ?? {})?.length > 0) {
    }

    return YAML.stringify(simpleConfig);
  }

  React.useEffect(() => {
    if (existingResolverConfig) {
      // this needs to be parsed because it shows fields we don't care about
      let stringifiedResolverConfig = resolverConfigToDisplay(
        existingResolverConfig
      );

      setFieldValue('resolverConfig', stringifiedResolverConfig);
    } else {
      setFieldValue('resolverConfig', demoConfig);
    }
  }, [!!existingResolverConfig]);

  const validateResolverSchema = async (resolver: string) => {
    setIsValid(!isValid);
    try {
      let res = await graphqlApi.validateResolverYaml({
        yaml: resolver,
        resolverType:
          values.resolverType === 'REST'
            ? ValidateResolverYamlRequest.ResolverType.REST_RESOLVER
            : ValidateResolverYamlRequest.ResolverType.REST_RESOLVER,
      });
      setIsValid(true);
      setErrorMessage('');
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

  return (
    <div data-testid='resolver-config-section' className='h-full p-6 pb-0 '>
      <div
        className={'flex items-center mb-2 text-lg font-medium text-gray-800'}>
        {isEdit ? 'Edit' : 'Configure'} Resolver{' '}
      </div>
      <div className=''>
        <div className='mb-2 '>
          <div>
            <EditorContainer editMode={true}>
              <div className=''>
                <div className='' style={{ height: 'min-content' }}>
                  {isValid ? (
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
                      } h-10`}>
                      <div
                        style={{ backgroundColor: '#FEF2F2' }}
                        className='p-2 text-orange-400 border border-orange-400 '>
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
                  <YamlEditor
                    mode='yaml'
                    theme='chrome'
                    name='resolverConfiguration'
                    style={{
                      width: '100%',
                      maxHeight: '36vh',
                      cursor: 'text',
                    }}
                    onChange={e => {
                      setFieldValue('resolverConfig', e);
                    }}
                    focus={true}
                    onInput={() => {
                      setIsValid(false);
                    }}
                    fontSize={16}
                    showPrintMargin={false}
                    showGutter={true}
                    highlightActiveLine={true}
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
    </div>
  );
};
