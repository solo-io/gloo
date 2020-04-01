import React from 'react';
import {
  Tabs,
  TabList,
  Tab,
  TabPanels,
  TabPanel,
  TabPanelProps
} from '@reach/tabs';
import { ReactComponent as VewIcon } from 'assets/view-icon-blue.svg';
import { css } from '@emotion/core';
import { Formik } from 'formik';
import {
  SectionContainer,
  SectionHeader,
  SectionSubHeader
} from '../apis/CreateAPIModal';
import {
  SoloFormInput,
  SoloFormCheckbox
} from 'Components/Common/Form/SoloFormField';
import {
  SoloButtonStyledComponent,
  SoloCancelButton
} from 'Styles/CommonEmotions/button';
import { SoloTransfer } from 'Components/Common/SoloTransfer';
import useSWR from 'swr';
import { apiDocApi, portalApi, userApi } from '../api';
import { Loading } from 'Components/Common/DisplayOnly/Loading';
import { ObjectRef } from 'proto/dev-portal/api/dev-portal/v1/common_pb';
import { User } from 'proto/dev-portal/api/grpc/admin/user_pb';

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
  const [showPassword, setShowPassword] = React.useState(true);
  return (
    <SectionContainer>
      <SectionHeader> Create User</SectionHeader>
      <div className='p-3 mb-2 text-gray-700 bg-gray-100 rounded-lg'>
        Create a new API to expose as business capabilities
      </div>
      <div className='grid grid-flow-col grid-cols-5 grid-rows-3 gap-2'>
        {/* <div className='grid grid-flow-col grid-cols-2 grid-rows-3 gap-2'> */}

        <div className='col-span-3 mr-4'>
          <SoloFormInput
            name='name'
            title='Name'
            placeholder='Username goes here'
            hideError
          />
        </div>
        <div className='col-span-3 mr-4'>
          <SoloFormInput
            name='email'
            title='Email'
            placeholder='email@domain.com'
            hideError
          />
        </div>
        <div className='relative col-span-3 mr-4'>
          <span
            className='absolute cursor-pointer bottom-2 right-4'
            onClick={() => setShowPassword(!showPassword)}>
            <VewIcon />
          </span>
          <SoloFormInput
            type={showPassword ? 'password' : 'text'}
            name='password'
            title='Password'
            placeholder='type password here'
            hideError
          />
        </div>
      </div>
    </SectionContainer>
  );
};

type CreateUserValues = {
  name: string;
  email: string;
  password: string;
  chosenAPIs: ObjectRef.AsObject[];
  chosenPortals: ObjectRef.AsObject[];
  chosenGroups: ObjectRef.AsObject[];
};

export const CreateUserModal: React.FC<{ onClose: () => void }> = props => {
  const { data: portalsList, error: portalsListError } = useSWR(
    'listPortals',
    portalApi.listPortals
  );
  const { data: apiDocsList, error: apiDocsError } = useSWR(
    'listApiDocs',
    apiDocApi.listApiDocs
  );
  const [tabIndex, setTabIndex] = React.useState(0);
  const handleTabsChange = (index: number) => {
    setTabIndex(index);
  };
  const handleCreateUser = async (values: CreateUserValues) => {
    const {
      chosenAPIs,
      name,
      chosenPortals,
      email,
      chosenGroups,
      password
    } = values;
    let newUser = new User().toObject();

    await userApi.createUser({
      user: {
        ...newUser!,
        metadata: {
          ...newUser.metadata!,
          name,
          namespace: 'gloo-system'
        },
        spec: {
          email,
          username: name
        }
      },
      password,
      portalsList: chosenPortals,
      groupsList: chosenGroups,
      apiDocsList: chosenAPIs,
      userOnly: false
    });
    props.onClose();
  };

  if (!apiDocsList || !portalsList) {
    return <Loading center>Loading...</Loading>;
  }
  return (
    <>
      <div
        css={css`
          width: 750px;
        `}
        className='bg-white rounded-lg shadow '>
        <Formik<CreateUserValues>
          initialValues={{
            name: '',
            email: '',
            password: '',
            chosenGroups: [] as ObjectRef.AsObject[],
            chosenAPIs: [] as ObjectRef.AsObject[],
            chosenPortals: [] as ObjectRef.AsObject[]
          }}
          onSubmit={handleCreateUser}>
          {formik => (
            <>
              {/* <pre>{JSON.stringify(formik.values, null, 2)}</pre> */}
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
                  <StyledTab>APIs</StyledTab>
                  <StyledTab>Portals</StyledTab>
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
                    <SectionContainer>
                      <SectionHeader>Create a User: APIs</SectionHeader>
                      <div className='p-3 mb-2 text-gray-700 bg-gray-100 rounded-lg'>
                        Select the APIs you'd like to make available to this
                        user
                      </div>
                      <SoloTransfer
                        allOptionsListName='Available APIs'
                        allOptions={apiDocsList
                          .sort((a, b) =>
                            a.metadata?.name === b.metadata?.name
                              ? 0
                              : a.metadata!.name > b.metadata!.name
                              ? 1
                              : -1
                          )
                          .map(apiDoc => {
                            return {
                              name: apiDoc.metadata?.name!,
                              namespace: apiDoc.metadata?.namespace!
                            };
                          })}
                        chosenOptionsListName='Selected APIs'
                        chosenOptions={formik.values.chosenAPIs.map(api => {
                          return { name: api.name, namespace: api.namespace };
                        })}
                        onChange={newChosenOptions => {
                          console.log('newChosenOptions', newChosenOptions);
                          formik.setFieldValue('chosenAPIs', newChosenOptions);
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
                      <SectionHeader>Create a User: Portal</SectionHeader>
                      <div className='p-3 mb-2 text-gray-700 bg-gray-100 rounded-lg'>
                        Select the portals you'd like to make available to this
                        user
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
                    </SectionContainer>
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
                          onClick={formik.handleSubmit}>
                          Create User
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
