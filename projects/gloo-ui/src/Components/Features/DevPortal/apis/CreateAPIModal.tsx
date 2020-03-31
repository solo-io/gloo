import React from 'react';
import {
  Tabs,
  TabList,
  Tab,
  TabPanels,
  TabPanel,
  TabPanelProps
} from '@reach/tabs';
import { css } from '@emotion/core';
import { Formik, useFormikContext } from 'formik';
import {
  SoloFormInput,
  SoloFormTextarea
} from 'Components/Common/Form/SoloFormField';
import ImageUploader from 'react-images-upload';
import { colors } from 'Styles';
import { ReactComponent as NoImageIcon } from 'assets/no-image-placeholder.svg';
import styled from '@emotion/styled';
import tw from 'tailwind.macro';
import { ReactComponent as GreenPlus } from 'assets/small-green-plus.svg';
import { ReactComponent as NoSelectedList } from 'assets/no-selected-list.svg';
import {
  SoloButtonStyledComponent,
  SoloCancelButton
} from 'Styles/CommonEmotions/button';
import { SoloTransfer } from 'Components/Common/SoloTransfer';
import useSWR from 'swr';
import { portalApi, userApi, groupApi, apiDocApi } from '../api';
import { Loading } from 'Components/Common/DisplayOnly/Loading';
import { ObjectRef } from 'proto/dev-portal/api/dev-portal/v1/common_pb';
import { Upload, Button } from 'antd';
import { ApiDoc } from 'proto/dev-portal/api/grpc/admin/apidoc_pb';

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

const GeneralSection = () => {
  return (
    <SectionContainer>
      <SectionHeader>Create an API</SectionHeader>
      <div className='p-3 mb-2 text-gray-700 bg-gray-100 rounded-lg'>
        Create a new API to expose as business capabilities
      </div>
      <div className='grid grid-cols-2 '>
        <div className='mr-4'>
          <SoloFormInput
            name='name'
            title='Name'
            placeholder='API title goes here'
            hideError
          />
        </div>
        <div>
          <SoloFormInput
            name='displayName'
            title='Display Name'
            placeholder='Display name goes here'
            hideError
          />
        </div>
        <div className='col-span-2 mt-2'>
          <SoloFormTextarea
            name='description'
            title='Description'
            placeholder='API description goes here'
            hideError
          />
        </div>
      </div>
    </SectionContainer>
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
          imgExtension={['.jpg', '.gif', '.png', '.gif', '.jpeg']}
          maxFileSize={5242880}
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
        Create an API by specification document (Open API Spec, Proto, etc)
      </div>
      <div className='mt-2 mb-1 font-medium text-gray-800'>
        Paste Swagger url
      </div>
      <div>
        <SoloFormInput name='swaggerUrl' hideError />
      </div>
      <div className='mt-2 mb-1 font-medium text-gray-800'>
        Upload Swagger File
      </div>
      <div>
        <Upload
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
            Upload a spec
          </Button>
        </Upload>
      </div>
    </SectionContainer>
  );
};

type CreateApiDocValues = {
  name: string;
  displayName: string;
  description: string;
  swaggerUrl: string;
  uploadedSwagger: File;
  image: File;
  chosenPortals: ObjectRef.AsObject[];
  chosenUsers: ObjectRef.AsObject[];
  chosenGroups: ObjectRef.AsObject[];
};

export const CreateAPIModal: React.FC<{ onClose: () => void }> = props => {
  const { data: portalsList, error: portalsListError } = useSWR(
    'listPortals',
    portalApi.listPortals
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

    let imageBuffer = await image.arrayBuffer();
    let imageUint8Arr = new Uint8Array(imageBuffer);

    let swaggerBuffer = await uploadedSwagger.arrayBuffer();
    let swaggerUint8Array = new Uint8Array(swaggerBuffer);

    //@ts-ignore
    await apiDocApi.createApiDoc({
      apidoc: {
        ...newApiDoc,
        metadata: {
          ...newApiDoc.metadata!,
          name: name || displayName,
          namespace: 'gloo-system'
        },
        spec: {
          dataSource: {
            ...newApiDoc.spec?.dataSource!
            // fetchUrl: swaggerUrl
            // inlineBytes: swaggerUint8Array
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

  if (!portalsList || !userList || !groupList) {
    return <Loading center>Loading...</Loading>;
  }

  return (
    <>
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
            chosenPortals: [] as ObjectRef.AsObject[],
            chosenUsers: [] as ObjectRef.AsObject[],
            chosenGroups: [] as ObjectRef.AsObject[]
          }}
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
                  <StyledTab>General</StyledTab>
                  <StyledTab>Imagery</StyledTab>
                  <StyledTab>Portals</StyledTab>
                  <StyledTab>User Access</StyledTab>
                  <StyledTab>Group Access</StyledTab>
                  <StyledTab>Spec</StyledTab>
                </TabList>

                <TabPanels className='bg-white rounded-r-lg'>
                  <TabPanel className='relative flex flex-col justify-between h-full focus:outline-none'>
                    <GeneralSection />
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
                              namespace: portal.metadata?.namespace!
                            };
                          })}
                        chosenOptionsListName='Selected Portal'
                        chosenOptions={formik.values.chosenPortals.map(
                          portal => {
                            return {
                              name: portal.name,
                              namespace: portal.namespace
                            };
                          }
                        )}
                        onChange={newChosenOptions => {
                          console.log('newChosenOptions', newChosenOptions);
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
                              namespace: user.metadata?.namespace!
                            };
                          })}
                        chosenOptionsListName='Selected Users'
                        chosenOptions={formik.values.chosenUsers.map(user => {
                          return { name: user.name, namespace: user.namespace };
                        })}
                        onChange={newChosenOptions => {
                          console.log('newChosenOptions', newChosenOptions);
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
                              namespace: group.metadata?.namespace!
                            };
                          })}
                        chosenOptionsListName='Selected Groups'
                        chosenOptions={formik.values.chosenGroups.map(user => {
                          return { name: user.name, namespace: user.namespace };
                        })}
                        onChange={newChosenOptions => {
                          console.log('newChosenOptions', newChosenOptions);
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
                          onClick={() => setTabIndex(tabIndex + 1)}>
                          Next Step
                        </SoloButtonStyledComponent>
                      </div>
                    </div>
                  </TabPanel>
                  <TabPanel className='relative flex flex-col justify-between h-full focus:outline-none'>
                    <SpecSection />
                    <div className='flex items-end justify-between h-full px-6 mb-4 '>
                      <button
                        className='text-blue-500 cursor-pointer'
                        onClick={props.onClose}>
                        cancel
                      </button>
                      <div>
                        <SoloCancelButton
                          onClick={() => handleTabsChange(4)}
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
    </>
  );
};
