import styled from '@emotion/styled/macro';
import { graphqlApi } from 'API/graphql';
import {
  useIsGlooFedEnabled,
  useListGlooInstances,
  useListGraphqlSchemas,
} from 'API/hooks';
import { DataError } from 'Components/Common/DataError';
import {
  ErrorText,
  SoloFormFileUpload,
  SoloFormInput,
} from 'Components/Common/SoloFormComponents';
import { SoloModal } from 'Components/Common/SoloModal';
import { Formik, useFormikContext } from 'formik';
import gql from 'graphql-tag';
import { GlooInstance } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/glooinstance_pb';
import * as React from 'react';
import { useNavigate, useParams } from 'react-router';
import { colors } from 'Styles/colors';
import { SoloButtonStyledComponent } from 'Styles/StyledComponents/button';
import { nameValidationSchema } from 'utils';
import * as yup from 'yup';

export interface NewApiModalProps {
  showNewModal: boolean;
  toggleNewModal: () => any;
}

const ModalContent = styled.div`
  padding: 25px 20px;
`;
const Title = styled.div`
  display: flex;
  font-size: 22px;
  line-height: 26px;
  font-weight: 500;
  margin-bottom: 20px;

  svg {
    margin-left: 8px;
  }
`;

const InputWrapper = styled.div`
  padding: 10px 0;
`;

const Footer = styled.footer`
  display: flex;
  flex-direction: row-reverse;
`;

const Button = styled.button`
  background-color: ${colors.seaBlue};
  color: white;
  padding: 15px;
  border: none;
  &:hover {
    cursor: pointer;
  }
`;

const validationSchema = yup.object().shape({
  name: nameValidationSchema.required('The API Environment must have a name.'),
});

type CreateApiValues = {
  name: string;
  schemaString: string;
  uploadedSchema: File;
};

export const NewApiModal = (props: NewApiModalProps) => {
  const { name = '', namespace = '' } = useParams();
  const [errorMessage, setErrorMessage] = React.useState('');
  const { showNewModal, toggleNewModal } = props;
  const { mutate } = useListGraphqlSchemas();
  const { data: glooInstances, error: instancesError } = useListGlooInstances();
  const { data: glooFedCheckResponse, error: glooFedCheckError } =
    useIsGlooFedEnabled();
  const isGlooFedEnabled = glooFedCheckResponse?.enabled;
  const [glooInstance, setGlooInstance] =
    React.useState<GlooInstance.AsObject>();

  React.useEffect(() => {
    if (glooInstances) {
      if (glooInstances.length === 1) {
        setGlooInstance(glooInstances[0]);
      } else {
        setGlooInstance(
          glooInstances.find(
            instance =>
              instance.metadata?.name === name &&
              instance.metadata?.namespace === namespace
          )
        );
      }
    } else {
      setGlooInstance(undefined);
    }
  }, [name, namespace, glooInstances]);

  const navigate = useNavigate();

  const createApi = async (values: CreateApiValues) => {
    let { uploadedSchema, name = '', schemaString } = values;
    mutate(async graphqlSchemas => {
      if (graphqlSchemas) {
        return [
          ...graphqlSchemas,
          {
            status: { state: 0 },
            metadata: {
              uid: 1,
              name,
              namespace: glooInstance?.metadata?.namespace,
              clusterName: glooInstance?.spec?.cluster,
            },
          } as any,
        ];
      }
    }, false);

    let createdGraphqlSchema = await graphqlApi
      .createGraphqlSchema({
        graphqlSchemaRef: {
          name,
          namespace: glooInstance?.metadata?.namespace!,
          clusterName: glooInstance?.spec?.cluster!,
        },
        spec: {
          executableSchema: {
            schemaDefinition: schemaString,
          },
	      allowedQueryHashesList: [],
        },
      })
      .catch(err => {
        // Catch any errors on the backend the frontend can't catch.
        setErrorMessage(err.message);
      });
    if (!createdGraphqlSchema) {
      return;
    }
    toggleNewModal();
    mutate();

    navigate(
      isGlooFedEnabled
        ? `/gloo-instances/${createdGraphqlSchema.glooInstance?.namespace}/${
            createdGraphqlSchema.glooInstance?.name
          }/apis/${glooInstance?.spec?.cluster!}/${
            createdGraphqlSchema.metadata?.namespace
          }/${createdGraphqlSchema.metadata?.name}/`
        : `/gloo-instances/${createdGraphqlSchema.glooInstance?.namespace}/${createdGraphqlSchema.glooInstance?.name}/apis/${createdGraphqlSchema.metadata?.namespace}/${createdGraphqlSchema.metadata?.name}/`
    );
  };
  return (
    <SoloModal visible={showNewModal} width={600} onClose={toggleNewModal}>
      <Formik
        validationSchema={validationSchema}
        initialValues={{
          uploadedSchema: undefined as unknown as File,
          name: '',
          schemaString: '',
        }}
        onSubmit={createApi}>
        {formik => (
          <ModalContent>
            <Title>Create new GraphQL API</Title>
            <InputWrapper>
              <SoloFormInput name='name' title='Name' />
            </InputWrapper>
            <InputWrapper>
              <SoloFormFileUpload
                name='uploadedSchema'
                title='Schema'
                buttonLabel='Upload Schema'
                fileType='.graphql,.gql'
                onRemoveFile={() => {
                  setErrorMessage('');
                }}
                validateFile={file => {
                  let schema = '';
                  if (file) {
                    const reader = new FileReader();
                    reader.onload = e => {
                      if (typeof e.target?.result === 'string') {
                        schema = e.target?.result;
                        formik.setFieldValue('schemaString', schema);
                      }
                    };

                    reader.readAsText(file!);

                    try {
                      let query = gql`
                        ${reader}
                      `;
                      setErrorMessage('');
                      formik.setFieldError('uploadedSchema', '');
                    } catch (error: any) {
                      setErrorMessage(error);
                      formik.setFieldError('uploadedSchema', error);

                      // TODO replace with real validation
                      return { isValid: true, errorMessage: error as string };
                    }
                  }

                  return { isValid: true, errorMessage: '' };
                }}
              />
            </InputWrapper>
            {!!errorMessage && <DataError error={errorMessage as any} />}
            <Footer>
              <SoloButtonStyledComponent onClick={formik.handleSubmit as any}>
                Create API
              </SoloButtonStyledComponent>
            </Footer>
          </ModalContent>
        )}
      </Formik>
    </SoloModal>
  );
};
