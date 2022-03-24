import { Alert } from 'antd';
import { graphqlConfigApi } from 'API/graphql';
import {
  useIsGlooFedEnabled,
  useListGraphqlApis,
  usePageGlooInstance,
} from 'API/hooks';
import { DataError } from 'Components/Common/DataError';
import {
  SoloFormFileUpload,
  SoloFormInput,
} from 'Components/Common/SoloFormComponents';
import { SoloModal } from 'Components/Common/SoloModal';
import { SoloRadioGroup } from 'Components/Common/SoloRadioGroup';
import { Formik } from 'formik';
import { ValidateSchemaDefinitionRequest } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/graphql_pb';
import React, { useMemo, useState } from 'react';
import { useNavigate } from 'react-router';
import { SoloButtonStyledComponent } from 'Styles/StyledComponents/button';
import { nameValidationSchema } from 'utils';
import * as yup from 'yup';
import * as styles from './NewApiModal.style';

export const NewApiModal: React.FC<{
  show: boolean;
  onClose: () => any;
}> = ({ show, onClose }) => {
  const navigate = useNavigate();
  const [glooInstance] = usePageGlooInstance();

  // Api
  const { mutate } = useListGraphqlApis();
  const isGlooFedEnabled = useIsGlooFedEnabled().data?.enabled;

  // State
  const [errorMessage, setErrorMessage] = useState('');
  const [warningMessage, setWarningMessage] = useState('');

  const initialValues = {
    name: '',
    apiType: 'executable' as 'gateway' | 'executable',
    schemaString: '',
    uploadedSchema: undefined as unknown as File,
  };
  const validationSchema = useMemo(() => {
    return yup.object().shape({
      name: nameValidationSchema.required(
        'The API Environment must have a name.'
      ),
      apiType: yup.string().matches(/^(gateway)|(executable)$/),
      uploadedSchema: yup.string().test((item, ctx) => {
        if (ctx.parent.apiType === 'executable') return !!item && !errorMessage;
        else return true;
      }),
    });
  }, [errorMessage]);

  const createApi = async ({
    uploadedSchema,
    name = '',
    apiType,
    schemaString,
  }: typeof initialValues) => {
    // Only executable APIs have uploaded schemas.
    if (apiType !== 'executable') schemaString = '';
    if (apiType === 'gateway') {
      alert('Creating Gateway GraphQL APIs is not currently supported.');
      return;
    }

    let createdGraphqlApi = await graphqlConfigApi
      .createGraphqlApi({
        graphqlApiRef: {
          name,
          namespace: glooInstance?.metadata?.namespace!,
          clusterName: glooInstance?.spec?.cluster!,
        },
        spec: {
          executableSchema: {
            schemaDefinition: schemaString,
            executor: {
              //@ts-ignore
              local: {
                enableIntrospection: true,
              },
            },
          },
          allowedQueryHashesList: [],
        },
      })
      .catch(err => {
        // Catch any errors on the backend the frontend can't catch.
        setErrorMessage(err.message);
      });
    if (!createdGraphqlApi) {
      return;
    }
    mutate(
      async graphqlApis => [
        ...(graphqlApis ?? []),
        {
          status: { state: 0 },
          metadata: {
            uid: 1,
            name,
            namespace: glooInstance?.metadata?.namespace,
            clusterName: glooInstance?.spec?.cluster,
          },
        } as any,
      ],
      false
    );

    onClose();
    mutate();

    navigate(
      isGlooFedEnabled
        ? `/gloo-instances/${createdGraphqlApi.glooInstance?.namespace}/${
            createdGraphqlApi.glooInstance?.name
          }/apis/${glooInstance?.spec?.cluster!}/${
            createdGraphqlApi.metadata?.namespace
          }/${createdGraphqlApi.metadata?.name}/`
        : `/gloo-instances/${createdGraphqlApi.glooInstance?.namespace}/${createdGraphqlApi.glooInstance?.name}/apis/${createdGraphqlApi.metadata?.namespace}/${createdGraphqlApi.metadata?.name}/`
    );
  };

  return (
    <SoloModal visible={show} width={600} onClose={onClose}>
      <Formik
        validationSchema={validationSchema}
        initialValues={initialValues}
        onSubmit={createApi}>
        {formik => {
          const { values } = formik;
          return (
            <styles.ModalContent>
              <styles.Title>Create new GraphQL API</styles.Title>
              {Boolean(warningMessage) && (
                <styles.StyledWarning className='p-2 text-orange-400 border border-orange-400 mb-5'>
                  {warningMessage}
                </styles.StyledWarning>
              )}
              <styles.InputWrapper>
                <SoloFormInput name='name' title='Name' />

                <div className='grid grid-cols-[min-content_auto]'>
                  <SoloRadioGroup
                    className='pb-3 pr-[3rem]'
                    title='Type'
                    options={[
                      {
                        id: 'executable',
                        displayName: 'Executable',
                      },
                      {
                        id: 'gateway',
                        displayName: 'Gateway',
                      },
                    ]}
                    forceAChoice={true}
                    currentSelection={values.apiType}
                    onChange={function (
                      idSelected: string | number | undefined
                    ) {
                      formik.setFieldValue('apiType', idSelected);
                    }}
                  />
                  <div className='mt-[1.75rem]'>
                    <Alert
                      type='info'
                      showIcon
                      message={
                        values.apiType[0].toUpperCase() +
                        values.apiType.slice(1) +
                        ' GraphQL API'
                      }
                      description={
                        values.apiType === 'executable' ? (
                          <div>
                            Create an executable GraphQL API from REST or gRPC
                            services.
                          </div>
                        ) : (
                          <div>
                            Create a gateway GraphQL API by stitching together
                            other GraphQL APIs.
                          </div>
                        )
                      }
                    />
                  </div>
                </div>

                {values.apiType === 'executable' && (
                  <SoloFormFileUpload
                    name='uploadedSchema'
                    title='Schema'
                    buttonLabel='Upload Schema'
                    fileType='.graphql,.gql'
                    onRemoveFile={() => {
                      setErrorMessage('');
                      setWarningMessage('');
                    }}
                    validateFile={file => {
                      return new Promise((resolve, reject) => {
                        formik.setFieldError('uploadedSchema', undefined);
                        setErrorMessage('');
                        setWarningMessage('');
                        let schema = '';
                        if (file) {
                          const reader = new FileReader();
                          reader.onload = e => {
                            if (typeof e.target?.result === 'string') {
                              schema = e.target?.result;
                              formik.setFieldValue('schemaString', schema);
                              const request =
                                new ValidateSchemaDefinitionRequest().toObject();
                              request.schemaDefinition = schema;
                              graphqlConfigApi
                                .validateSchema(request)
                                .then(() => {
                                  resolve({ isValid: true, errorMessage: '' });
                                })
                                .catch(err => {
                                  setWarningMessage(err.message);
                                  reject({
                                    isValid: false,
                                    errorMessage: err.message,
                                  });
                                });
                            }
                          };

                          reader.readAsText(file);
                        } else {
                          resolve({ isValid: true, errorMessage: '' });
                        }
                      });
                    }}
                  />
                )}
                {!!errorMessage && <DataError error={errorMessage as any} />}
              </styles.InputWrapper>
              <styles.Footer>
                <SoloButtonStyledComponent
                  disabled={!formik.dirty || !formik.isValid}
                  onClick={formik.handleSubmit as any}>
                  Create API
                </SoloButtonStyledComponent>
              </styles.Footer>
            </styles.ModalContent>
          );
        }}
      </Formik>
    </SoloModal>
  );
};
