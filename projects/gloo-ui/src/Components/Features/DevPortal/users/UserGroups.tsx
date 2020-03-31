import React from 'react';
import { useParams, useHistory } from 'react-router';
import { Breadcrumb } from 'Components/Common/Breadcrumb';
import { SectionCard } from 'Components/Common/SectionCard';
import { ReactComponent as UserIcon } from 'assets/user-icon.svg';
import { healthConstants, colors, soloConstants } from 'Styles';
import { css } from '@emotion/core';
import {
  Tabs,
  TabList,
  Tab,
  TabPanels,
  TabPanel,
  TabsProps,
  TabPanelProps
} from '@reach/tabs';
import { SoloInput } from 'Components/Common/SoloInput';
import { ReactComponent as EditIcon } from 'assets/edit-pencil.svg';
import { ReactComponent as NoUser } from 'assets/no-user-icon.svg';

import { TabCss, ActiveTabCss } from '../portals/PortalDetails';
import { ReactComponent as GreenPlus } from 'assets/small-green-plus.svg';
import { CreateUserModal } from './CreateUserModal';
import { SoloModal } from 'Components/Common/SoloModal';
import { CreateGroupModal } from './CreateGroupModal';
import useSWR from 'swr';
import { userApi, groupApi } from '../api';
import { Loading } from 'Components/Common/DisplayOnly/Loading';

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
      css={css`
        ${TabCss}
        ${isSelected ? ActiveTabCss : ''}
      `}
      className='border rounded-lg focus:outline-none'>
      {children}
    </Tab>
  );
};
export const UserGroups = () => {
  const { data: userList, error: userError } = useSWR(
    'listUsers',
    userApi.listUsers
  );

  const { data: groupList, error: groupError } = useSWR(
    'listGroups',
    groupApi.listGroups
  );

  const [tabIndex, setTabIndex] = React.useState(0);
  const [userSearchTerm, setUserSearchTerm] = React.useState('');
  const [groupSearchTerm, setGroupSearchTerm] = React.useState('');
  const [showCreateUserModal, setShowCreateUserModal] = React.useState(false);
  const [showCreateGroupModal, setShowCreateGroupModal] = React.useState(false);

  const handleTabsChange = (index: number) => {
    setTabIndex(index);
  };
  if (!userList || !groupList) {
    return <Loading center>Loading...</Loading>;
  }
  return (
    <div className='relative'>
      <div className='absolute top-0 right-0 flex items-center-mt-8'>
        <span
          onClick={() => setShowCreateUserModal(true)}
          className='flex items-center text-green-400 cursor-pointer hover:text-green-300'>
          <GreenPlus className='mr-1 fill-current' />
          <span className='mr-2 text-gray-700'> Create a User</span>
        </span>
        <span
          onClick={() => setShowCreateGroupModal(true)}
          className='flex items-center text-green-400 cursor-pointer hover:text-green-300'>
          <GreenPlus className='mr-1 fill-current' />
          <span className='text-gray-700'> Create a Group</span>
        </span>
      </div>
      <Tabs
        index={tabIndex}
        className='mb-4 border-none rounded-lg '
        onChange={handleTabsChange}>
        <TabList className='flex items-start ml-4 '>
          <StyledTab>Users</StyledTab>
          <StyledTab>Groups</StyledTab>
        </TabList>
        <TabPanels
          className='bg-white rounded-lg '
          css={css`
            margin-top: -1px;
          `}>
          <TabPanel className='focus:outline-none'>
            <div className='relative flex flex-col p-4 border border-gray-300 rounded-lg'>
              <div className='w-full mb-4'>
                <SoloInput
                  placeholder='Search by member by user name or email...'
                  value={userSearchTerm}
                  onChange={e => setUserSearchTerm(e.target.value)}
                />
              </div>
              <div className='flex flex-col'>
                <div className='py-2 -my-2 overflow-x-auto sm:-mx-6 sm:px-6 lg:-mx-8 lg:px-8'>
                  <div className='inline-block min-w-full overflow-hidden align-middle border-b border-gray-200 shadow sm:rounded-lg'>
                    <table className='min-w-full'>
                      <thead className='bg-gray-300 '>
                        <tr>
                          <th className='px-6 py-3 text-sm font-medium leading-4 tracking-wider text-left text-gray-800 capitalize border-b border-gray-200 bg-gray-50'>
                            User Name
                          </th>
                          <th className='px-6 py-3 text-sm font-medium leading-4 tracking-wider text-left text-gray-800 capitalize border-b border-gray-200 bg-gray-50'>
                            Email
                          </th>
                          <th className='px-6 py-3 text-sm font-medium leading-4 tracking-wider text-left text-gray-800 capitalize border-b border-gray-200 bg-gray-50'>
                            Groups
                          </th>
                          <th className='px-6 py-3 text-sm font-medium leading-4 tracking-wider text-left text-gray-800 capitalize border-b border-gray-200 bg-gray-50'>
                            API Access
                          </th>
                          <th className='px-6 py-3 text-sm font-medium leading-4 tracking-wider text-left text-gray-800 capitalize border-b border-gray-200 bg-gray-50'>
                            Portal Access
                          </th>

                          <th className='px-6 py-3 text-sm font-medium leading-4 tracking-wider text-left text-gray-800 capitalize border-b border-gray-200 bg-gray-50'>
                            Actions
                          </th>
                        </tr>
                      </thead>
                      <tbody className='bg-white'>
                        {userList
                          .sort((a, b) =>
                            a.metadata?.name === b.metadata?.name
                              ? 0
                              : a.metadata!.name > b.metadata!.name
                              ? 1
                              : -1
                          )
                          .map(user => {
                            return (
                              <tr key={user.metadata?.uid}>
                                <td className='max-w-xs px-6 py-4 border-b border-gray-200'>
                                  <div className='flex items-center'>
                                    <div className='flex-shrink-0 w-8 h-8 text-blue-400'>
                                      <div className='flex items-center justify-center w-8 h-8 text-white bg-blue-600 rounded-full'>
                                        <UserIcon className='w-6 h-6 fill-current' />
                                      </div>
                                    </div>
                                    <div className='ml-4'>
                                      <div className='text-sm font-medium leading-5'>
                                        {user.spec?.username}
                                      </div>
                                    </div>
                                  </div>
                                </td>
                                <td className='max-w-xs px-6 py-4 border-b border-gray-200'>
                                  <div className='text-sm leading-5 text-gray-700'>
                                    <span className='flex items-center '>
                                      {user.spec?.email}
                                    </span>
                                  </div>
                                </td>
                                <td className='max-w-xs px-6 py-4 border-b border-gray-200'>
                                  <div className='text-sm leading-5 text-gray-700'>
                                    <span className='flex items-center '>
                                      {user.spec?.email}
                                    </span>
                                  </div>
                                </td>
                                <td className='max-w-xs px-6 py-4 border-b border-gray-200'>
                                  <div className='text-sm leading-5 text-gray-700'>
                                    <span className='flex items-center '>
                                      View (4)
                                    </span>
                                  </div>
                                </td>
                                <td className='max-w-xs px-6 py-4 border-b border-gray-200'>
                                  <div className='text-sm leading-5 text-gray-700'>
                                    <span className='flex items-center '>
                                      View (2)
                                    </span>
                                  </div>
                                </td>
                                <td className='max-w-xs px-6 py-4 text-sm font-medium leading-5 text-right border-b border-gray-200'>
                                  <span className='flex items-center'>
                                    <div className='flex items-center justify-center w-4 h-4 mr-3 text-gray-700 bg-gray-400 rounded-full cursor-pointer'>
                                      <EditIcon className='w-2 h-3 fill-current' />
                                    </div>

                                    <div
                                      className='flex items-center justify-center w-4 h-4 text-gray-700 bg-gray-400 rounded-full cursor-pointer'
                                      onClick={() => {}}>
                                      x
                                    </div>
                                  </span>
                                </td>
                              </tr>
                            );
                          })}
                      </tbody>
                    </table>
                    {userList.length === 0 && (
                      <div className='w-full m-auto'>
                        <div className='flex flex-col items-center justify-center w-full h-full py-4 mr-32 bg-white rounded-lg shadow-lg md:flex-row'>
                          <div className='mr-6'>
                            <NoUser />
                          </div>
                          <div className='flex flex-col h-full'>
                            <p className='h-auto my-6 text-lg font-medium text-gray-800 '>
                              There are no Users to display!{' '}
                            </p>
                          </div>
                        </div>
                      </div>
                    )}
                  </div>
                </div>
              </div>
            </div>
          </TabPanel>
          <TabPanel className='focus:outline-none'>
            <div className='relative flex flex-col p-4 border border-gray-300 rounded-lg'>
              <div className='w-full mb-4'>
                <SoloInput
                  placeholder='Search group by user name or email...'
                  value={groupSearchTerm}
                  onChange={e => setGroupSearchTerm(e.target.value)}
                />
              </div>
              <div className='flex flex-col'>
                <div className='py-2 -my-2 overflow-x-auto sm:-mx-6 sm:px-6 lg:-mx-8 lg:px-8'>
                  <div className='inline-block min-w-full overflow-hidden align-middle border-b border-gray-200 shadow sm:rounded-lg'>
                    <table className='min-w-full'>
                      <thead className='bg-gray-300 '>
                        <tr>
                          <th className='px-6 py-3 text-sm font-medium leading-4 tracking-wider text-left text-gray-800 capitalize border-b border-gray-200 bg-gray-50'>
                            Group Name
                          </th>
                          <th className='px-6 py-3 text-sm font-medium leading-4 tracking-wider text-left text-gray-800 capitalize border-b border-gray-200 bg-gray-50'>
                            Description
                          </th>

                          <th className='px-6 py-3 text-sm font-medium leading-4 tracking-wider text-left text-gray-800 capitalize border-b border-gray-200 bg-gray-50'>
                            Members
                          </th>
                          <th className='px-6 py-3 text-sm font-medium leading-4 tracking-wider text-left text-gray-800 capitalize border-b border-gray-200 bg-gray-50'>
                            API Access
                          </th>
                          <th className='px-6 py-3 text-sm font-medium leading-4 tracking-wider text-left text-gray-800 capitalize border-b border-gray-200 bg-gray-50'>
                            Portal Access
                          </th>

                          <th className='px-6 py-3 text-sm font-medium leading-4 tracking-wider text-left text-gray-800 capitalize border-b border-gray-200 bg-gray-50'>
                            Actions
                          </th>
                        </tr>
                      </thead>
                      <tbody className='bg-white'>
                        {groupList
                          .sort((a, b) =>
                            a.metadata?.name === b.metadata?.name
                              ? 0
                              : a.metadata!.name > b.metadata!.name
                              ? 1
                              : -1
                          )
                          .map(group => (
                            <tr key={group.metadata?.uid}>
                              <td className='max-w-xs px-6 py-4 border-b border-gray-200'>
                                <div className='flex items-center'>
                                  <div className='flex-shrink-0 w-8 h-8 text-blue-400'>
                                    <div className='flex items-center justify-center w-8 h-8 text-white bg-blue-600 rounded-full'>
                                      <UserIcon className='w-6 h-6 fill-current' />
                                    </div>
                                  </div>
                                  <div className='ml-4'>
                                    <div className='text-sm font-medium leading-5'>
                                      {group.spec?.displayName}
                                    </div>
                                  </div>
                                </div>
                              </td>
                              <td className='max-w-xs px-6 py-4 border-b border-gray-200'>
                                <div className='text-sm leading-5 text-gray-700'>
                                  <span className='flex items-center '>
                                    {group.spec?.description}
                                  </span>
                                </div>
                              </td>
                              <td className='max-w-xs px-6 py-4 border-b border-gray-200'>
                                <div className='text-sm leading-5 text-gray-700'>
                                  <span className='flex items-center '>
                                    {group.status?.usersList.map(
                                      user => user.name
                                    )}
                                  </span>
                                </div>
                              </td>
                              <td className='max-w-xs px-6 py-4 border-b border-gray-200'>
                                <div className='text-sm leading-5 text-gray-700'>
                                  <span className='flex items-center '>
                                    api access
                                  </span>
                                </div>
                              </td>
                              <td className='max-w-xs px-6 py-4 border-b border-gray-200'>
                                <div className='text-sm leading-5 text-gray-700'>
                                  <span className='flex items-center '>
                                    portals access
                                  </span>
                                </div>
                              </td>
                              <td className='max-w-xs px-6 py-4 text-sm font-medium leading-5 text-right border-b border-gray-200'>
                                <span className='flex items-center'>
                                  <div className='flex items-center justify-center w-4 h-4 mr-3 text-gray-700 bg-gray-400 rounded-full cursor-pointer'>
                                    <EditIcon className='w-2 h-3 fill-current' />
                                  </div>

                                  <div
                                    className='flex items-center justify-center w-4 h-4 text-gray-700 bg-gray-400 rounded-full cursor-pointer'
                                    onClick={() => {}}>
                                    x
                                  </div>
                                </span>
                              </td>
                            </tr>
                          ))}
                      </tbody>
                    </table>
                    {groupList.length === 0 && (
                      <div className='w-full m-auto'>
                        <div className='flex flex-col items-center justify-center w-full h-full py-4 mr-32 bg-white rounded-lg shadow-lg md:flex-row'>
                          <div className='mr-6'>
                            <NoUser />
                          </div>
                          <div className='flex flex-col h-full'>
                            <p className='h-auto my-6 text-lg font-medium text-gray-800 '>
                              There are no Groups to display!{' '}
                            </p>
                          </div>
                        </div>
                      </div>
                    )}
                  </div>
                </div>
              </div>
            </div>
          </TabPanel>
        </TabPanels>
      </Tabs>
      <SoloModal visible={showCreateUserModal} width={750} noPadding={true}>
        <CreateUserModal onClose={() => setShowCreateUserModal(false)} />
      </SoloModal>
      <SoloModal visible={showCreateGroupModal} width={750} noPadding={true}>
        <CreateGroupModal onClose={() => setShowCreateGroupModal(false)} />
      </SoloModal>
    </div>
  );
};
