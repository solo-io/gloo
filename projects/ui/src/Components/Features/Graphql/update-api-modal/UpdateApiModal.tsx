import { Alert } from 'antd';
import { getGraphqlApiPb, graphqlConfigApi } from 'API/graphql';
import {
  useIsGlooFedEnabled,
  useListGraphqlApis,
  usePageGlooInstance,
} from 'API/hooks';
import { DataError } from 'Components/Common/DataError';
import { SoloFormFileUpload } from 'Components/Common/SoloFormComponents';
import { SoloModal } from 'Components/Common/SoloModal';
import { Formik } from 'formik';
import { ClusterObjectRef } from 'proto/github.com/solo-io/skv2/api/core/v1/core_pb';
import { ValidateSchemaDefinitionRequest } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/graphql_pb';
import React, { useState } from 'react';
import { SoloButtonStyledComponent } from 'Styles/StyledComponents/button';
import * as styles from './UpdateApiModal.style';

export const UpdateApiModal: React.FC<{
  show: boolean;
  onClose: () => any;
  apiRef: ClusterObjectRef.AsObject;
  schemaString: string;
}> = ({ show, onClose, apiRef, schemaString }) => {
  // Api
  const { mutate } = useListGraphqlApis();

  // State
  const [errorMessage, setErrorMessage] = useState('');
  const [warningMessage, setWarningMessage] = useState('');
  const [loading, setLoading] = useState(false);

  const initialValues = {
    schemaString,
  };

  const createApi = async ({ schemaString }: typeof initialValues) => {
    // Only executable APIs have uploaded schemas.
    setLoading(true);
    let createdGraphqlApi = await graphqlConfigApi
      .updateGraphqlApi({
        graphqlApiRef: apiRef,
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
      }).finally(() => {
        setLoading(false);
      });
    if (!createdGraphqlApi) {
      return;
    }
    mutate(
      async graphqlApis => {
        const foundIndex = graphqlApis?.findIndex((g) => {
          return g.metadata?.name === apiRef.name && g.metadata?.namespace === apiRef.namespace;
        }) ?? -1;
        if (foundIndex > -1) {
          (graphqlApis as any)[foundIndex] = createdGraphqlApi;
        }
        return [
        ...(graphqlApis ?? []),
        {
          status: { state: 0 },
          metadata: {
            uid: 1,
            ...apiRef,
          },
        } as any,
      ]},
      false
    );

    mutate();
    onClose();
  };

  return (
    <SoloModal visible={show} width={600} onClose={onClose}>
      <Formik initialValues={initialValues} onSubmit={createApi}>
        {formik => {
          const isDisabled = !formik.dirty || !formik.isValid || Boolean(warningMessage) || loading || Boolean(errorMessage);
          return (
            <styles.ModalContent>
              <styles.Title>Update Schema</styles.Title>
              {Boolean(warningMessage) && (
                <styles.StyledWarning className='p-2 text-orange-400 border border-orange-400 mb-5'>
                  {warningMessage}
                </styles.StyledWarning>
              )}
              <styles.InputWrapper>
                <SoloFormFileUpload
                  name='uploadedSchema'
                  title='Schema'
                  buttonLabel='Upload Schema'
                  fileType='.graphql,.gql'
                  onRemoveFile={() => {
                    setLoading(false);
                    setErrorMessage('');
                    setWarningMessage('');
                  }}
                  validateFile={file => {
                    return new Promise((resolve, reject) => {
                      setLoading(true);
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

                            return getGraphqlApiPb(apiRef)
                              .then(res => {
                                request.spec = res.toObject().spec;
                                request.spec!.executableSchema!.schemaDefinition =
                                  schema;
                                return graphqlConfigApi
                                  .validateSchema({
                                    ...request,
                                    apiRef,
                                  })
                                  .then(() => {
                                    resolve({
                                      isValid: true,
                                      errorMessage: '',
                                    });
                                  })
                                  .catch(validationError => {
                                    throw validationError;
                                  });
                              })
                              .catch(err => {
                                setWarningMessage(err.message);
                                reject({
                                  isValid: false,
                                  errorMessage: err.message,
                                });
                              }).finally(() => {
                                setLoading(false);
                              });
                          }
                        };
                        reader.readAsText(file);
                      } else {
                        setLoading(false);
                        resolve({ isValid: true, errorMessage: '' });
                      }
                    });
                  }}
                />
                {!!errorMessage && <DataError error={errorMessage as any} />}
              </styles.InputWrapper>
              <styles.Footer>
                <SoloButtonStyledComponent
                  disabled={isDisabled}
                  onClick={formik.handleSubmit as any}>
                  Update Schema
                </SoloButtonStyledComponent>
              </styles.Footer>
            </styles.ModalContent>
          );
        }}
      </Formik>
    </SoloModal>
  );
};
