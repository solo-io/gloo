import { css } from '@emotion/core';
import styled from '@emotion/styled';
import {
  Tab,
  TabList,
  TabPanel,
  TabPanelProps,
  TabPanels,
  Tabs
} from '@reach/tabs';
import { Button, Upload } from 'antd';
import { ReactComponent as NoImageIcon } from 'assets/no-image-placeholder.svg';
import { Loading } from 'Components/Common/DisplayOnly/Loading';
import { SoloFormInput } from 'Components/Common/Form/SoloFormField';
import { SoloTransfer, ListItemType } from 'Components/Common/SoloTransfer';
import { Formik, useFormikContext } from 'formik';
import { ObjectRef } from 'proto/dev-portal/api/dev-portal/v1/common_pb';
import { ApiDoc } from 'proto/dev-portal/api/grpc/admin/apidoc_pb';
import React from 'react';
import ImageUploader from 'react-images-upload';
import { configAPI } from 'store/config/api';
import { colors } from 'Styles';
import {
  SoloButtonStyledComponent,
  SoloCancelButton
} from 'Styles/CommonEmotions/button';
import useSWR from 'swr';
import tw from 'tailwind.macro';
import { apiDocApi, groupApi, portalApi, userApi } from '../api';
import { ErrorBoundary } from 'Components/Features/Errors/ErrorBoundary';
import * as Yup from 'yup';

export const SectionContainer = styled.div`
  ${tw`w-full h-full p-6 pb-0`}
`;

export const SectionHeader = styled.div`
  ${tw`mb-4 text-lg font-medium text-gray-800`}
`;

export const SectionSubHeader = styled.div`
  ${tw`mb-1 font-medium text-gray-800`}
`;

const StyledTab = (
  props: {
    disabled?: boolean | undefined;
  } & TabPanelProps & {
      isSelected?: boolean | undefined;
    }
) => {
  const { isSelected, children } = props;
  return (
    <Tab
      {...props}
      className={`p-1 text-left w-48 text-white  pl-6 mb-2 focus:outline-none ${
        isSelected
          ? ' bg-blue-500 border-r-8 border-blue-300  '
          : 'bg-blue-600 '
      }`}>
      {children}
    </Tab>
  );
};

const ImagerySection = () => {
  const formik = useFormikContext();

  const [images, setImages] = React.useState<any[]>([]);

  const onDrop = async (files: File[], pictures: string[]) => {
    formik.setFieldValue('image', files[0]);

    setImages([...images, ...pictures]);
  };
  return (
    <SectionContainer>
      <SectionHeader>Create an API: Add Imagery</SectionHeader>
      <div className='p-3 mb-2 text-gray-700 bg-gray-100 rounded-lg'>
        Add an image associated with your API
      </div>
      <div className='flex flex-col items-center p-4 pb-0 mr-4 bg-gray-100 border border-gray-400 rounded-lg'>
        <NoImageIcon />
        <ImageUploader
          label='Max 1.5MiB'
          css={css`
            .fileContainer {
              background: ${colors.januaryGrey};
              border-radius: 8px;
              box-shadow: none;
              padding: 0px;
            }
            .uploadPictureContainer {
              margin: 0;
            }
          `}
          withPreview
          withIcon={false}
          buttonStyles={{
            borderRadius: '8px',
            backgroundColor: colors.blue600,
            fontSize: '16px',
            fontWeight: 400
          }}
          buttonText='Upload an Image'
          onChange={onDrop}
          imgExtension={['.jpg', '.png', '.jpeg', 'ico']}
          maxFileSize={1572864}
        />
      </div>
    </SectionContainer>
  );
};

const SpecSection = () => {
  const formik = useFormikContext();

  return (
    <SectionContainer>
      <SectionHeader>Create an API: Specs</SectionHeader>
      <div className='p-2 bg-gray-100 '>
        Create an API by providing an Open API document
      </div>
      <div className='mt-2 mb-1 font-medium text-gray-800'>
        Paste OpenAPI document URL
      </div>
      <div>
        <SoloFormInput name='swaggerUrl' hideError />
      </div>
      <div className='mt-2 mb-1 font-medium text-gray-800'>
        Upload OpenAPI document
      </div>
      <div>
        <Upload
          action={'https://www.mocky.io/v2/5cc8019d300000980a055e76'}
          name='uploadedSwagger'
          onChange={info => {
            if (info.file.status === 'done') {
              formik.setFieldValue('uploadedSwagger', info.file.originFileObj);
            }
          }}>
          <Button
            css={css`
              width: 100%;
              border-radius: 9px;
            `}>
            Upload a document{' '}
          </Button>
        </Upload>
      </div>
    </SectionContainer>
  );
};

const validationSchema = Yup.object().shape(
  {
    swaggerUrl: Yup.string().when('uploadedSwagger', {
      is: undefined,
      then: Yup.string().required('This field is required.'),
      otherwise: Yup.string()
    })
  },
  [['swaggerUrl', 'uploadedSwagger']]
);

type CreateApiDocValues = {
  name: string;
  displayName: string;
  description: string;
  swaggerUrl: string;
  uploadedSwagger: File;
  image: File;
  chosenPortals: ListItemType[];
  chosenUsers: ListItemType[];
  chosenGroups: ListItemType[];
};

export const CreateAPIModal: React.FC<{ onClose: () => void }> = props => {
  const { data: portalsList, error: portalsListError } = useSWR(
    'listPortals',
    portalApi.listPortals
  );
  const { data: podNamespace, error: podNamespaceError } = useSWR(
    'getPodNamespace',
    configAPI.getPodNamespace
  );
  const { data: userList, error: userError } = useSWR(
    'listUsers',
    userApi.listUsers
  );

  const { data: groupList, error: groupError } = useSWR(
    'listGroups',
    groupApi.listGroups
  );
  const [tabIndex, setTabIndex] = React.useState(0);

  const handleTabsChange = (index: number) => {
    setTabIndex(index);
  };

  const handleCreateApiDoc = async (values: CreateApiDocValues) => {
    const {
      name,
      description,
      displayName,
      image,
      chosenPortals,
      chosenGroups,
      chosenUsers,
      swaggerUrl,
      uploadedSwagger
    } = values;

    let newApiDoc = new ApiDoc().toObject();

    let imageUint8Arr = new Uint8Array();
    if (image !== undefined) {
      let imageBuffer = await image.arrayBuffer();
      imageUint8Arr = new Uint8Array(imageBuffer);
    }

    let swaggerUint8Array = new Uint8Array();
    if (uploadedSwagger !== undefined) {
      let swaggerBuffer = await uploadedSwagger.arrayBuffer();
      swaggerUint8Array = new Uint8Array(swaggerBuffer);
    }

    //@ts-ignore
    await apiDocApi.createApiDoc({
      apidoc: {
        ...newApiDoc,
        metadata: {
          ...newApiDoc.metadata!,
          name: name || displayName,
          namespace: podNamespace!
        },
        spec: {
          dataSource: {
            ...newApiDoc.spec?.dataSource!,
            fetchUrl: swaggerUrl,
            inlineBytes: swaggerUint8Array
          },
          image: {
            ...newApiDoc.spec?.image!,
            inlineBytes: imageUint8Arr
          }
        },
        status: {
          ...newApiDoc.status!,
          description,
          displayName
        }
      },

      portalsList: chosenPortals,
      groupsList: chosenGroups,
      usersList: chosenUsers
    });
    props.onClose();
  };

  if (!portalsList || !userList || !groupList || !podNamespace) {
    return <Loading center>Loading...</Loading>;
  }

  return (
    <ErrorBoundary fallback={<div>There was an error.</div>}>
      <div
        css={css`
          width: 750px;
        `}
        className='bg-white rounded-lg shadow '>
        <Formik<CreateApiDocValues>
          initialValues={{
            displayName: '',
            description: '',
            swaggerUrl: '',
            uploadedSwagger: (undefined as unknown) as File,
            image: (undefined as unknown) as File,
            name: '',
            chosenPortals: [] as ListItemType[],
            chosenUsers: [] as ListItemType[],
            chosenGroups: [] as ListItemType[]
          }}
          validationSchema={validationSchema}
          onSubmit={handleCreateApiDoc}>
          {formik => (
            <>
              <Tabs
                className='bg-blue-600 rounded-lg h-96'
                index={tabIndex}
                onChange={handleTabsChange}
                css={css`
                  display: grid;
                  height: 450px;
                  grid-template-columns: 190px 1fr;
                `}>
                <TabList className='flex flex-col mt-6'>
                  <StyledTab>Spec</StyledTab>
                  <StyledTab>Imagery</StyledTab>
                  <StyledTab>Portals</StyledTab>
                  <StyledTab>User Access</StyledTab>
                  <StyledTab>Group Access</StyledTab>
                </TabList>

                <TabPanels className='bg-white rounded-r-lg'>
                  <TabPanel className='relative flex flex-col justify-between h-full focus:outline-none'>
                    <SpecSection />
                    <div className='flex items-end justify-between h-full px-6 mb-4 '>
                      <button
                        className='text-blue-500 cursor-pointer'
                        onClick={props.onClose}>
                        cancel
                      </button>
                      <SoloButtonStyledComponent
                        onClick={() => setTabIndex(tabIndex + 1)}>
                        Next Step
                      </SoloButtonStyledComponent>
                    </div>
                  </TabPanel>
                  <TabPanel className='relative flex flex-col justify-between h-full focus:outline-none'>
                    <ImagerySection />
                    <div className='flex items-end justify-between h-full px-6 mb-4 '>
                      <button
                        className='text-blue-500 cursor-pointer'
                        onClick={props.onClose}>
                        cancel
                      </button>
                      <div>
                        <SoloCancelButton
                          onClick={() => handleTabsChange(0)}
                          className='mr-2'>
                          Back
                        </SoloCancelButton>
                        <SoloButtonStyledComponent
                          onClick={() => setTabIndex(tabIndex + 1)}>
                          Next Step
                        </SoloButtonStyledComponent>
                      </div>
                    </div>
                  </TabPanel>
                  <TabPanel className='relative flex flex-col justify-between h-full focus:outline-none'>
                    <SectionContainer>
                      <SectionHeader>Create an API: Portal</SectionHeader>
                      <div className='p-3 mb-2 text-gray-700 bg-gray-100 rounded-lg'>
                        Select the portals to which you'd like to publish this
                        API
                      </div>
                      <SoloTransfer
                        allOptionsListName='Available Portals'
                        allOptions={portalsList
                          .sort((a, b) =>
                            a.metadata?.name === b.metadata?.name
                              ? 0
                              : a.metadata!.name > b.metadata!.name
                              ? 1
                              : -1
                          )
                          .map(portal => {
                            return {
                              name: portal.metadata?.name!,
                              namespace: portal.metadata?.namespace!,
                              displayValue: portal.spec?.displayName
                            };
                          })}
                        chosenOptionsListName='Selected Portal'
                        chosenOptions={formik.values.chosenPortals}
                        onChange={newChosenOptions => {
                          formik.setFieldValue(
                            'chosenPortals',
                            newChosenOptions
                          );
                        }}
                      />
                    </SectionContainer>{' '}
                    <div className='flex items-end justify-between h-full px-6 mb-4 '>
                      <button
                        className='text-blue-500 cursor-pointer'
                        onClick={props.onClose}>
                        cancel
                      </button>
                      <div>
                        <SoloCancelButton
                          onClick={() => handleTabsChange(1)}
                          className='mr-2'>
                          Back
                        </SoloCancelButton>
                        <SoloButtonStyledComponent
                          onClick={() => setTabIndex(tabIndex + 1)}>
                          Next Step
                        </SoloButtonStyledComponent>
                      </div>
                    </div>
                  </TabPanel>
                  <TabPanel className='relative flex flex-col justify-between h-full focus:outline-none'>
                    <SectionContainer>
                      <SectionHeader>Create an API: User Access</SectionHeader>
                      <div className='p-3 mb-2 text-gray-700 bg-gray-100 rounded-lg'>
                        Select the users and groups to which you'd like to grant
                        access to this API
                      </div>

                      <SoloTransfer
                        allOptionsListName='Available Users'
                        allOptions={userList
                          .sort((a, b) =>
                            a.metadata?.name === b.metadata?.name
                              ? 0
                              : a.metadata!.name > b.metadata!.name
                              ? 1
                              : -1
                          )
                          .map(user => {
                            return {
                              name: user.metadata?.name!,
                              namespace: user.metadata?.namespace!,
                              displayValue: user.spec?.username
                            };
                          })}
                        chosenOptionsListName='Selected Users'
                        chosenOptions={formik.values.chosenUsers}
                        onChange={newChosenOptions => {
                          formik.setFieldValue('chosenUsers', newChosenOptions);
                        }}
                      />
                    </SectionContainer>
                    <div className='flex items-end justify-between h-full px-6 mb-4 '>
                      <button
                        className='text-blue-500 cursor-pointer'
                        onClick={props.onClose}>
                        cancel
                      </button>
                      <div>
                        <SoloCancelButton
                          onClick={() => handleTabsChange(2)}
                          className='mr-2'>
                          Back
                        </SoloCancelButton>
                        <SoloButtonStyledComponent
                          onClick={() => setTabIndex(tabIndex + 1)}>
                          Next Step
                        </SoloButtonStyledComponent>
                      </div>
                    </div>
                  </TabPanel>
                  <TabPanel className='relative flex flex-col justify-between h-full focus:outline-none'>
                    <SectionContainer>
                      <SectionHeader>Create an API: Group Access</SectionHeader>
                      <div className='p-3 mb-2 text-gray-700 bg-gray-100 rounded-lg'>
                        Select the users and groups to which you'd like to grant
                        access to this API
                      </div>

                      <SoloTransfer
                        allOptionsListName='Available Groups'
                        allOptions={groupList
                          .sort((a, b) =>
                            a.metadata?.name === b.metadata?.name
                              ? 0
                              : a.metadata!.name > b.metadata!.name
                              ? 1
                              : -1
                          )
                          .map(group => {
                            return {
                              name: group.metadata?.name!,
                              namespace: group.metadata?.namespace!,
                              displayValue: group.spec?.displayName
                            };
                          })}
                        chosenOptionsListName='Selected Groups'
                        chosenOptions={formik.values.chosenGroups}
                        onChange={newChosenOptions => {
                          formik.setFieldValue(
                            'chosenGroups',
                            newChosenOptions
                          );
                        }}
                      />
                    </SectionContainer>
                    <div className='flex items-end justify-between h-full px-6 mb-4 '>
                      <button
                        className='text-blue-500 cursor-pointer'
                        onClick={props.onClose}>
                        cancel
                      </button>
                      <div>
                        <SoloCancelButton
                          onClick={() => handleTabsChange(3)}
                          className='mr-2'>
                          Back
                        </SoloCancelButton>
                        <SoloButtonStyledComponent
                          onClick={formik.handleSubmit}>
                          Create API
                        </SoloButtonStyledComponent>
                      </div>
                    </div>
                  </TabPanel>
                </TabPanels>
              </Tabs>
            </>
          )}
        </Formik>
      </div>
    </ErrorBoundary>
  );
};
