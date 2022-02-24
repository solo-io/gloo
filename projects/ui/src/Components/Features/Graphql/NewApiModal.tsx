import * as React from 'react';
import styled from '@emotion/styled/macro';
import { SoloModal } from 'Components/Common/SoloModal';
import { SectionCard } from 'Components/Common/SectionCard';
import { SoloInput } from 'Components/Common/SoloInput';
import { colors } from 'Styles/colors';
import { SoloButtonStyledComponent } from 'Styles/StyledComponents/button';
import {
  SoloFormFileUpload,
  SoloFormInput,
} from 'Components/Common/SoloFormComponents';
import { Formik } from 'formik';
import { GraphQLFileLoader } from '@graphql-tools/graphql-file-loader';
import gql from 'graphql-tag';
import { graphqlApi } from 'API/graphql';
import {
  useIsGlooFedEnabled,
  useListGlooInstances,
  useListGraphqlSchemas,
} from 'API/hooks';
import { GraphqlSchema } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/graphql_pb';
import { useNavigate, useParams } from 'react-router';
import { GlooInstance } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/glooinstance_pb';

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

type CreateApiValues = {
  name: string;
  schemaString: string;
  uploadedSchema: File;
};

export const NewApiModal = (props: NewApiModalProps) => {
  const { name = '', namespace = '' } = useParams();

  const { showNewModal, toggleNewModal } = props;
  const [_schemaFile, setSchemaFile] = React.useState<File>();
  const {
    data: graphqlSchemas,
    error: graphqlSchemaError,
    mutate,
  } = useListGraphqlSchemas();
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

    let createdGraphqlSchema = await graphqlApi.createGraphqlSchema({
      graphqlSchemaRef: {
        name,
        namespace: glooInstance?.metadata?.namespace!,
        clusterName: glooInstance?.spec?.cluster!,
      },
      spec: {
        executableSchema: {
          schemaDefinition: schemaString,
        },
      },
    });
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
                    } catch (error) {
                      // TODO replace with real validation
                      return { isValid: true, errorMessage: error as string };
                    }
                  }
                  return { isValid: true, errorMessage: '' };
                }}
              />
            </InputWrapper>
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
