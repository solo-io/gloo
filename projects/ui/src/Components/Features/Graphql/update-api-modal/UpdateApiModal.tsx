import { graphqlConfigApi } from 'API/graphql';
import { useListGraphqlApis } from 'API/hooks';
import { TabList, TabPanel, TabPanels } from '@reach/tabs';
import { DataError } from 'Components/Common/DataError';
import styled from '@emotion/styled';
import {
  SoloModal,
  StyledModalTab,
  StyledModalTabs,
} from 'Components/Common/SoloModal';
import { Formik } from 'formik';
import { ClusterObjectRef } from 'proto/github.com/solo-io/skv2/api/core/v1/core_pb';
import React, { useState } from 'react';
import {
  SoloButtonStyledComponent,
  SoloCancelButton,
} from 'Styles/StyledComponents/button';
import { UpdateApiFile } from './UpdateApiFile';
import * as styles from './UpdateApiModal.style';
import { UpdateApiEditor } from './UpdateApiEditor';
import { colors } from 'Styles/colors';
import { GraphqlDiff } from './GraphqlDiff';
import { ValidateSchemaDefinitionRequest } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/graphql_pb';

enum UploadType {
  FILE = 'FILE',
  TEXT = 'TEXT',
}

const StyledSoloModal = styled(SoloModal)`
  height: 100%;
`;

const UpdatedStyleTabs = styled(StyledModalTabs)`
  height: 100%;
`;

const FooterContainer = styled.div`
  position: absolute;
  bottom: 20px;
  right: 0;
`;

export const UpdateApiModal: React.FC<{
  show: boolean;
  onClose: () => any;
  apiRef: ClusterObjectRef.AsObject;
  schemaString: string;
}> = ({ show, onClose, apiRef, schemaString }) => {
  // Api
  const { mutate } = useListGraphqlApis();

  // State
  const [tabIndex, setTabIndex] = React.useState(0);
  const [errorMessage, setErrorMessage] = useState('');
  const [warningMessage, setWarningMessage] = useState('');
  const [loading, setLoading] = useState(false);
  const [uploadType, setUploadType] = useState<UploadType>(UploadType.FILE);

  const pageOptions: {
    displayName: string;
    id: string | number;
  }[] = [
    {
      displayName: 'File Upload',
      id: UploadType.FILE,
    },
    {
      displayName: 'Text Editor',
      id: UploadType.TEXT,
    },
  ];

  const handleRadioChange = (newId?: string | number) => {
    const found = pageOptions.find(p => p.id === newId);
    if (found) {
      setUploadType(found.id as UploadType);
    }
  };

  // TODO:  Add an onclick listener for resetting errors and checking them.

  const initialValues = {
    schemaString,
  };

  const validateApi = async (schemaString: string) => {
    const validationRequest = new ValidateSchemaDefinitionRequest();
    const validationObj = validationRequest.toObject();
    await graphqlConfigApi
      .validateSchema({
        ...validationObj,
        ...{
          schemaDefinition: schemaString,
          apiRef,
        },
      })
      .catch(err => {
        // Catch any errors on the backend the frontend can't catch.
        setWarningMessage(err.message);
        throw err;
      });
  };

  // TODO:  This can be extracted, and it should.
  const createApi = async ({ schemaString }: typeof initialValues) => {
    // Only executable APIs have uploaded schemas.
    // TODO: Add in the validation right here.
    setLoading(true);

    let createdGraphqlApi = await validateApi(schemaString)
      .then(_res => {
        return graphqlConfigApi.updateGraphqlApi({
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
        });
      })
      .catch(err => {
        // Catch any errors on the backend the frontend can't catch.
        setErrorMessage(err.message);
      })
      .finally(() => {
        setLoading(false);
      });
    if (!createdGraphqlApi) {
      return;
    }
    mutate(async graphqlApis => {
      const foundIndex =
        graphqlApis?.findIndex(g => {
          return (
            g.metadata?.name === apiRef.name &&
            g.metadata?.namespace === apiRef.namespace
          );
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
      ];
    }, false);

    mutate();
    onClose();
  };

  return (
    <StyledSoloModal visible={show} width={800} onClose={onClose}>
      <Formik initialValues={initialValues} onSubmit={createApi}>
        {formik => {
          const setSchema = (newValue: string) => {
            formik.setFieldValue('schemaString', newValue);
          };
          const handleReset = () => {
            formik.setFieldValue('schemaString', schemaString);
          };
          const isDisabled =
            Boolean(warningMessage) ||
            loading ||
            Boolean(errorMessage) ||
            schemaString === formik.values.schemaString;

          return (
            <>
              <styles.ModalContent className='h-full'>
                <styles.InputWrapper>
                  <UpdatedStyleTabs
                    style={{ backgroundColor: colors.oceanBlue }}
                    className='grid rounded-lg grid-cols-[150px_1fr] absolute top-0 left-0 w-full h-full'
                    index={tabIndex}
                    onChange={setTabIndex}>
                    <TabList className='flex flex-col mt-6'>
                      {/* --- SIDEBAR --- */}
                      <StyledModalTab
                        data-testid='update-schema-string-tab'
                        isCompleted={
                          schemaString !== formik.values.schemaString
                        }>
                        Update Schema
                      </StyledModalTab>
                      <StyledModalTab data-testid='update-schema-confirm-changes-tab'>
                        Confirm Changes
                      </StyledModalTab>
                    </TabList>

                    <TabPanels className='bg-white rounded-r-lg flex flex-col h-full p-6'>
                      {/* --- Update Schema --- */}
                      <TabPanel className='mt-3 relative flex flex-col h-full pb-4 focus:outline-none'>
                        <UpdateApiFile
                          setLoading={setLoading}
                          setErrorMessage={setErrorMessage}
                          setWarningMessage={setWarningMessage}
                          apiRef={apiRef}
                        />
                        <UpdateApiEditor
                          setGraphqlSchema={setSchema}
                          graphqlSchema={formik.values.schemaString}
                        />
                        <FooterContainer>
                          <styles.Footer>
                            <SoloButtonStyledComponent
                              disabled={
                                schemaString === formik.values.schemaString
                              }
                              onClick={() => {
                                setTabIndex(1);
                                validateApi(formik.values.schemaString);
                              }}>
                              Next
                            </SoloButtonStyledComponent>
                            <SoloCancelButton onClick={handleReset as any}>
                              Reset Schema
                            </SoloCancelButton>
                          </styles.Footer>
                        </FooterContainer>
                      </TabPanel>
                      <TabPanel className='mt-3 relative flex flex-col h-full pb-4 focus:outline-none'>
                        <GraphqlDiff
                          validateApi={validateApi}
                          warningMessage={warningMessage}
                          originalSchemaString={schemaString}
                          newSchemaString={formik.values.schemaString}
                          setWarningMessage={setWarningMessage}
                        />
                        <FooterContainer>
                          <styles.Footer>
                            <SoloButtonStyledComponent
                              disabled={isDisabled}
                              onClick={formik.handleSubmit as any}>
                              Update Schema
                            </SoloButtonStyledComponent>
                          </styles.Footer>
                        </FooterContainer>
                      </TabPanel>
                    </TabPanels>
                  </UpdatedStyleTabs>

                  {!!errorMessage && <DataError error={errorMessage as any} />}
                </styles.InputWrapper>
              </styles.ModalContent>
            </>
          );
        }}
      </Formik>
    </StyledSoloModal>
  );
};
