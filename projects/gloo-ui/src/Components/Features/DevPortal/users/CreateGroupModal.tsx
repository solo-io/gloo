import React from 'react';

import {
  Tabs,
  TabList,
  Tab,
  TabPanels,
  TabPanel,
  TabPanelProps
} from '@reach/tabs';
import { ReactComponent as GreenPlus } from 'assets/small-green-plus.svg';
import { css } from '@emotion/core';
import { Formik } from 'formik';
import { SectionContainer, SectionHeader } from '../apis/CreateAPIModal';
import {
  SoloFormInput,
  SoloFormTextarea
} from 'Components/Common/Form/SoloFormField';
import {
  SoloButtonStyledComponent,
  SoloCancelButton
} from 'Styles/CommonEmotions/button';
import { SoloTransfer, ListItemType } from 'Components/Common/SoloTransfer';
import useSWR from 'swr';
import { userApi, apiDocApi, portalApi, groupApi } from '../api';
import { Loading } from 'Components/Common/DisplayOnly/Loading';
import { ObjectRef } from 'proto/dev-portal/api/dev-portal/v1/common_pb';
import { Group } from 'proto/dev-portal/api/grpc/admin/group_pb';
import { configAPI } from 'store/config/api';

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

type GeneralSectionProps = {
  isEdit: boolean;
};

const GeneralSection = ({ isEdit }: GeneralSectionProps) => {
  return (
    <SectionContainer>
      <SectionHeader> {isEdit ? 'Edit' : 'Create a'} Group</SectionHeader>
      <div className='p-3 mb-2 text-gray-700 bg-gray-100 rounded-lg'>
        {isEdit ? 'Edit a' : 'Create a new'} group of users that can be treated
        as one entity to atomically assign permissions for APIs and Portals to
        many users
      </div>
      <div className='flex flex-col items-center w-full'>
        <div className='w-full mb-4 mr-4'>
          <SoloFormInput
            name='name'
            title='Group Name'
            placeholder='Name of the group'
            hideError
          />
        </div>
        <div className='w-full mb-4 mr-4'>
          <SoloFormTextarea
            name='description'
            title='Description'
            placeholder='Description of the group'
            hideError
          />
        </div>
      </div>
    </SectionContainer>
  );
};

type CreateGroupValues = {
  name: string;
  description: string;
  chosenAPIs: ListItemType[];
  chosenPortals: ListItemType[];
  chosenUsers: ListItemType[];
};

export const CreateGroupModal: React.FC<{
  onClose: () => void;
  existingGroup?: Group.AsObject;
}> = props => {
  const { data: portalsList, error: portalsListError } = useSWR(
    'listPortals',
    portalApi.listPortals
  );

  const { data: userList, error: userError } = useSWR(
    'listUsers',
    userApi.listUsers
  );

  const { data: podNamespace, error: podNamespaceError } = useSWR(
    'getPodNamespace',
    configAPI.getPodNamespace
  );

  const { data: apiDocsList, error: apiDocsError } = useSWR(
    'listApiDocs',
    apiDocApi.listApiDocs
  );

  const [tabIndex, setTabIndex] = React.useState(0);
  const handleTabsChange = (index: number) => {
    setTabIndex(index);
  };

  const handleWriteGroup = async (values: CreateGroupValues) => {
    const {
      chosenAPIs,
      chosenPortals,
      chosenUsers,
      name,
      description
    } = values;
    if (!props.existingGroup) {
      let newGroup = new Group().toObject();
      //@ts-ignore
      await groupApi.createGroup({
        apiDocsList: chosenAPIs,
        portalsList: chosenPortals,
        usersList: chosenUsers,
        group: {
          ...newGroup!,
          metadata: {
            ...newGroup.metadata!,
            namespace: podNamespace!
          },
          spec: {
            description,
            displayName: name
          },
          status: {
            ...newGroup.status!
          }
        }
      });
    } else {
      groupApi.updateGroup(
        {
          namespace: props.existingGroup.metadata!.namespace,
          name: props.existingGroup.metadata!.name
        },
        name,
        description,
        chosenAPIs,
        chosenUsers,
        chosenPortals
      );
    }
    props.onClose();
  };

  if (!userList || !apiDocsList || !portalsList) {
    return <Loading center>Loading...</Loading>;
  }

  const getInitialApis = (): ListItemType[] => {
    if (!props.existingGroup || !props.existingGroup.status?.accessLevel) {
      return [];
    }

    return props.existingGroup.status.accessLevel.apiDocsList.map(
      ({ namespace, name }) => {
        let matchingDoc = apiDocsList.find(
          doc =>
            doc.metadata?.namespace === namespace && doc.metadata.name === name
        );
        return {
          namespace,
          name,
          displayValue: matchingDoc?.status?.displayName
        };
      }
    );
  };

  const getInitialPortals = (): ListItemType[] => {
    if (!props.existingGroup || !props.existingGroup.status?.accessLevel) {
      return [];
    }

    return props.existingGroup.status.accessLevel.portalsList.map(
      ({ namespace, name }) => {
        let matchingPortal = portalsList.find(
          portal =>
            portal.metadata?.namespace === namespace &&
            portal.metadata.name === name
        );

        return {
          namespace,
          name,
          displayValue: matchingPortal?.spec?.displayName
        };
      }
    );
  };

  const getInitialUsers = (): ListItemType[] => {
    if (!props.existingGroup || !props.existingGroup.status?.accessLevel) {
      return [];
    }

    return props.existingGroup.status.usersList.map(({ namespace, name }) => {
      let matchingUser = userList.find(
        user =>
          user.metadata?.namespace === namespace && user.metadata.name === name
      );

      return {
        namespace,
        name,
        displayValue: matchingUser?.spec?.username
      };
    });
  };

  return (
    <>
      <div
        css={css`
          width: 750px;
        `}
        className='bg-white rounded-lg shadow '>
        <Formik<CreateGroupValues>
          initialValues={{
            name: props.existingGroup?.spec?.displayName || '',
            description: props.existingGroup?.spec?.description || '',
            chosenAPIs: getInitialApis(),
            chosenPortals: getInitialPortals(),
            chosenUsers: getInitialUsers()
          }}
          onSubmit={handleWriteGroup}>
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
                  <StyledTab>Users</StyledTab>

                  <StyledTab>APIs</StyledTab>
                  <StyledTab>Portals</StyledTab>
                </TabList>

                <TabPanels className='bg-white rounded-r-lg'>
                  <TabPanel className='flex flex-col justify-between h-full focus:outline-none'>
                    <GeneralSection isEdit={!!props.existingGroup} />
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
                  <TabPanel className='flex flex-col justify-between h-full focus:outline-none'>
                    <SectionContainer>
                      <SectionHeader>
                        {!!props.existingGroup
                          ? 'Update Group'
                          : 'Create a Group'}
                        : Users
                      </SectionHeader>
                      <div className='p-3 mb-2 text-gray-700 bg-gray-100 rounded-lg'>
                        Select the users that are included in this group
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
                        chosenOptionsListName='Selected User'
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
                  <TabPanel className='flex flex-col justify-between h-full focus:outline-none'>
                    <SectionContainer>
                      <SectionHeader>
                        {!!props.existingGroup
                          ? 'Update Group'
                          : 'Create a Group'}
                        : APIs
                      </SectionHeader>
                      <div className='p-3 mb-2 text-gray-700 bg-gray-100 rounded-lg'>
                        Select the APIs you'd like to make available to this
                        group
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
                              namespace: apiDoc.metadata?.namespace!,
                              displayValue: apiDoc.status?.displayName
                            };
                          })}
                        chosenOptionsListName='Selected APIs'
                        chosenOptions={formik.values.chosenAPIs}
                        onChange={newChosenOptions => {
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
                  <TabPanel className='flex flex-col justify-between h-full focus:outline-none'>
                    <SectionContainer>
                      <SectionHeader>
                        {!!props.existingGroup
                          ? 'Update Group'
                          : 'Create a Group'}
                        : Portal
                      </SectionHeader>
                      <div className='p-3 mb-2 text-gray-700 bg-gray-100 rounded-lg'>
                        Select the portals you'd like to make available to this
                        group
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
                          onClick={formik.handleSubmit}>
                          {!!props.existingGroup
                            ? 'Update Group'
                            : 'Create Group'}
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
