import React from 'react';
import { SoloInput } from 'Components/Common/SoloInput';
import { ReactComponent as EditIcon } from 'assets/edit-pencil.svg';
import useSWR from 'swr';
import { groupApi } from '../api';
import { format } from 'timeago.js';
import { SoloModal } from 'Components/Common/SoloModal';
import { ReactComponent as GreenPlus } from 'assets/small-green-plus.svg';
import { Portal } from 'proto/dev-portal/api/grpc/admin/portal_pb';
import { Group } from 'proto/dev-portal/api/grpc/admin/group_pb';
import { CreateUserModal } from '../users/CreateUserModal';

type PortalGroupsTabProps = {
  portal: Portal.AsObject;
};
export const PortalGroupsTab = ({ portal }: PortalGroupsTabProps) => {
  const { data: groupsList, error: groupsError } = useSWR(
    `listGroups${portal.metadata?.name}${portal.metadata?.namespace}`,
    () =>
      groupApi.listGroups({
        portalsList: [
          { name: portal.metadata!.name, namespace: portal.metadata!.namespace }
        ],
        apiDocsList: []
      })
  );

  const [groupSearchTerm, setGroupSearchTerm] = React.useState('');
  const [showCreateUserModal, setShowCreateUserModal] = React.useState(false);

  return (
    <div className='relative flex flex-col p-4 border border-gray-300 rounded-lg'>
      <span
        onClick={() => setShowCreateUserModal(true)}
        className='absolute top-0 right-0 flex items-center mt-2 mr-2 text-green-400 cursor-pointer hover:text-green-300'>
        <GreenPlus className='mr-1 fill-current' />
        <span className='text-gray-700'> Create a Group</span>
      </span>
      <div className='w-1/3 m-4'>
        <SoloInput
          placeholder='Search by API name...'
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
                    Actions
                  </th>
                </tr>
              </thead>
              <tbody className='bg-white'>
                {!!groupsList &&
                  groupsList!.map(group => {
                    return (
                      <tr>
                        <td className='px-6 py-4 whitespace-no-wrap border-b border-gray-200'>
                          <div className='text-sm leading-5 text-gray-900'>
                            <span className='flex items-center capitalize'>
                              {group?.spec?.displayName || group.metadata?.name}
                            </span>
                          </div>
                        </td>
                        <td className='px-6 py-4 whitespace-no-wrap border-b border-gray-200'>
                          <div className='text-sm leading-5 text-gray-900'>
                            <span className='flex items-center capitalize'>
                              {group?.spec?.description}
                            </span>
                          </div>
                        </td>
                        <td className='px-6 py-4 whitespace-no-wrap border-b border-gray-200'>
                          <div className='text-sm leading-5 text-gray-900'>
                            <span className='flex items-center capitalize'>
                              {group?.status?.usersList
                                .map(u => u.name)
                                .join(', ')}
                            </span>
                          </div>
                        </td>
                        <td className='px-6 py-4 text-sm font-medium leading-5 text-right whitespace-no-wrap border-b border-gray-200'>
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
          </div>
        </div>
      </div>
      <SoloModal visible={showCreateUserModal} width={750} noPadding={true}>
        <CreateUserModal onClose={() => setShowCreateUserModal(false)} />
      </SoloModal>
    </div>
  );
};
