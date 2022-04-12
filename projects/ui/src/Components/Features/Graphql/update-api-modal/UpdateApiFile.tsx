import { useFormikContext } from 'formik';
import * as React from 'react';
import { SoloFormFileUpload } from 'Components/Common/SoloFormComponents';
import { ValidateSchemaDefinitionRequest } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/graphql_pb';
import { getGraphqlApiPb, graphqlConfigApi } from 'API/graphql';
import { ClusterObjectRef } from 'proto/github.com/solo-io/skv2/api/core/v1/core_pb';
import styled from '@emotion/styled';

const FileContainer = styled.div`
  margin-top: 20px;
  margin-bottom: 20px;
`;

export interface UpdateApiFileProps {
  setLoading: (isLoading: boolean) => any;
  setErrorMessage: (errorMessage: string) => any;
  setWarningMessage: (warningMessage: string) => any;
  apiRef: ClusterObjectRef.AsObject;
}

export const UpdateApiFile = (props: UpdateApiFileProps) => {
  const { setLoading, setErrorMessage, setWarningMessage, apiRef } = props;
  const formik = useFormikContext();
  return (
    <FileContainer>
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
                      request.spec!.executableSchema!.schemaDefinition = schema;
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
                    })
                    .finally(() => {
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
    </FileContainer>
  );
};
