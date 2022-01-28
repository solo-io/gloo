import * as React from 'react';
import styled from '@emotion/styled/macro';
import { SoloModal } from 'Components/Common/SoloModal';
import { SectionCard } from 'Components/Common/SectionCard';
import { SoloInput } from 'Components/Common/SoloInput';
import { colors } from 'Styles/colors';
import { SoloButtonStyledComponent } from 'Styles/StyledComponents/button';
import { SoloFormFileUpload } from 'Components/Common/SoloFormComponents';
import { Formik } from 'formik';
import { GraphQLFileLoader } from '@graphql-tools/graphql-file-loader';
import gql from 'graphql-tag';

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
  uploadedSchema: File;
};

export const NewApiModal = (props: NewApiModalProps) => {
  const { showNewModal, toggleNewModal } = props;
  const [name, setName] = React.useState<string>('');
  const [_schemaFile, setSchemaFile] = React.useState<File>();

  const changeName = (e: React.ChangeEvent<HTMLInputElement>) => {
    setName(e.target.value);
  };

  // Check .graphql files as well.
  const changeSchema = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      setSchemaFile(file);
    }
  };

  const createApi = async (values: CreateApiValues) => {};

  return (
    <SoloModal visible={showNewModal} width={600} onClose={toggleNewModal}>
      <Formik
        initialValues={{
          uploadedSchema: undefined as unknown as File,
        }}
        onSubmit={createApi}
      >
        {formik => (
          <ModalContent>
            <Title>Create new GraphQL API</Title>
            <InputWrapper>
              <SoloInput title='Name' onChange={changeName} value={name} />
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
function loadDocuments(
  uploadedSchema: File,
  arg1: {
    // load from a single schema file
    loaders: GraphQLFileLoader[];
  }
) {
  throw new Error('Function not implemented.');
}
