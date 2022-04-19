import { Alert } from 'antd';
import { graphqlConfigApi } from 'API/graphql';
import { useIsGlooFedEnabled, useListGraphqlApis } from 'API/hooks';
import { DataError } from 'Components/Common/DataError';
import {
  SoloFormFileUpload,
  SoloFormInput,
} from 'Components/Common/SoloFormComponents';
import { SoloModal } from 'Components/Common/SoloModal';
import { SoloRadioGroup } from 'Components/Common/SoloRadioGroup';
import { Formik } from 'formik';
import { GlooInstance } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/glooinstance_pb';
import {
  CreateGraphqlApiRequest,
  ValidateSchemaDefinitionRequest,
} from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/graphql_pb';
import React, { useMemo, useState } from 'react';
import { useNavigate } from 'react-router';
import { SoloButtonStyledComponent } from 'Styles/StyledComponents/button';
import { nameValidationSchema } from 'utils';
import { makeGraphqlApiLink } from 'utils/graphql-helpers';
import * as yup from 'yup';
import * as styles from './NewApiModal.style';

export const NewApiModal: React.FC<{
  glooInstance: GlooInstance.AsObject;
  show: boolean;
  onClose(): void;
}> = ({ glooInstance, show, onClose }) => {
  const navigate = useNavigate();

  // Api
  const { mutate } = useListGraphqlApis();
  const isGlooFedEnabled = useIsGlooFedEnabled().data?.enabled;

  // State
  const [errorMessage, setErrorMessage] = useState('');
  const [warningMessage, setWarningMessage] = useState('');

  const initialValues = {
    name: '',
    apiType: 'executable' as 'stitched' | 'executable',
    schemaString: '',
    uploadedSchema: undefined as unknown as File,
  };
  const validationSchema = useMemo(() => {
    return yup.object().shape({
      name: nameValidationSchema.required(
        'The API Environment must have a name.'
      ),
      apiType: yup.string().matches(/^(stitched)|(executable)$/),
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
    const newApiSchema: CreateGraphqlApiRequest.AsObject = {
      graphqlApiRef: {
        name,
        namespace: glooInstance?.metadata?.namespace!,
        clusterName: glooInstance?.spec?.cluster!,
      },
      spec: {
        allowedQueryHashesList: [],
      },
    };
    if (!apiType) return;
    else if (apiType === 'stitched') {
      newApiSchema.spec!.stitchedSchema = { subschemasList: [] };
    } else if (apiType === 'executable') {
      newApiSchema.spec!.executableSchema = {
        schemaDefinition: schemaString,
        executor: {
          //@ts-ignore
          local: {
            enableIntrospection: true,
          },
        },
      };
    }

    let createdGraphqlApi = await graphqlConfigApi
      .createGraphqlApi(newApiSchema)
      .catch(err => {
        // Catch any errors on the backend the frontend can't catch.
        setErrorMessage(err.message);
      });
    if (!createdGraphqlApi) {
      return;
    }
    mutate(
      graphqlApis => [
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
      makeGraphqlApiLink(
        createdGraphqlApi.metadata?.name,
        createdGraphqlApi.metadata?.namespace,
        glooInstance?.spec?.cluster,
        glooInstance?.metadata?.name,
        glooInstance?.metadata?.namespace,
        isGlooFedEnabled
      )
    );
  };

  return (
    <SoloModal
      data-testid='new-api-modal-modal'
      visible={show}
      width={600}
      onClose={onClose}>
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
                <styles.StyledWarning
                  data-testid='new-api-modal-warning-message'
                  className='p-2 text-orange-400 border border-orange-400 mb-5'>
                  {warningMessage}
                </styles.StyledWarning>
              )}
              <styles.InputWrapper>
                <SoloFormInput
                  data-testid='new-api-modal-name'
                  name='name'
                  title='Name'
                />

                <div className='grid grid-cols-[min-content_auto]'>
                  <SoloRadioGroup
                    data-testid='new-api-modal-apitype'
                    className='pb-3 pr-[3rem]'
                    title='Type'
                    options={[
                      {
                        id: 'executable',
                        displayName: 'Executable',
                      },
                      {
                        id: 'stitched',
                        displayName: 'Stitched',
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
                            Create a stitched GraphQL API by combining other
                            GraphQL APIs.
                          </div>
                        )
                      }
                    />
                  </div>
                </div>

                {values.apiType === 'executable' && (
                  <SoloFormFileUpload
                    data-testid='new-api-modal-file-upload'
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
                                  resolve({
                                    isValid: true,
                                    errorMessage: '',
                                  });
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
                  data-testid='new-api-modal-submit'
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
