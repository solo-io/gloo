import React from 'react';
import { useParams, useHistory } from 'react-router';
import { SectionCard } from 'Components/Common/SectionCard';
import { ReactComponent as EditIcon } from 'assets/edit-pencil.svg';
import { ReactComponent as KeyIcon } from 'assets/key-on-ring.svg';
import { SoloInput } from 'Components/Common/SoloInput';

export const APIKeys = () => {
  const history = useHistory();
  const [apiKeySearchTerm, setApiKeySearchTerm] = React.useState('');

  return (
    <SectionCard
      cardName={'API Keys'}
      logoIcon={
        <span className='text-blue-500'>
          <KeyIcon className='fill-current' />
        </span>
      }
      onClose={() => history.push(`/dev-portal/`)}>
      <div className='relative flex flex-col p-2 rounded-lg'>
        <div className='w-full mb-4'>
          <SoloInput
            placeholder='Search group by user name or email...'
            value={apiKeySearchTerm}
            onChange={e => setApiKeySearchTerm(e.target.value)}
          />
        </div>
        <div className='flex flex-col'>
          <div className='py-2 -my-2 overflow-x-auto sm:-mx-6 sm:px-6 lg:-mx-8 lg:px-8'>
            <div className='inline-block min-w-full overflow-hidden align-middle border-b border-gray-200 shadow sm:rounded-lg'>
              <table className='min-w-full'>
                <thead className='bg-gray-300 '>
                  <tr>
                    <th className='px-6 py-3 text-sm font-medium leading-4 tracking-wider text-left text-gray-800 capitalize border-b border-gray-200 bg-gray-50'>
                      Secret
                    </th>
                    <th className='px-6 py-3 text-sm font-medium leading-4 tracking-wider text-left text-gray-800 capitalize border-b border-gray-200 bg-gray-50'>
                      Key
                    </th>

                    <th className='px-6 py-3 text-sm font-medium leading-4 tracking-wider text-left text-gray-800 capitalize border-b border-gray-200 bg-gray-50'>
                      User
                    </th>
                    <th className='px-6 py-3 text-sm font-medium leading-4 tracking-wider text-left text-gray-800 capitalize border-b border-gray-200 bg-gray-50'>
                      API Key Scope
                    </th>
                    <th className='px-6 py-3 text-sm font-medium leading-4 tracking-wider text-left text-gray-800 capitalize border-b border-gray-200 bg-gray-50'>
                      Labels
                    </th>

                    <th className='px-6 py-3 text-sm font-medium leading-4 tracking-wider text-left text-gray-800 capitalize border-b border-gray-200 bg-gray-50'>
                      Actions
                    </th>
                  </tr>
                </thead>
                <tbody className='bg-white'>
                  {/* {(filteredUsers
                              ? filteredUsers
                              : organizationMembersList
                            )?.map(orgMember => ( */}
                  <tr>
                    <td className='max-w-xs px-6 py-4 border-b border-gray-200'>
                      <div className='text-sm leading-5 text-gray-700'>
                        <span className='flex items-center '>
                          Getting Started
                        </span>
                      </div>
                    </td>
                    <td className='max-w-xs px-6 py-4 border-b border-gray-200'>
                      <div className='text-sm leading-5 text-gray-700'>
                        <span className='flex items-center '>
                          Getting Started
                        </span>
                      </div>
                    </td>
                    <td className='max-w-xs px-6 py-4 border-b border-gray-200'>
                      <div className='text-sm leading-5 text-gray-700'>
                        <span className='flex items-center '>
                          Getting Started
                        </span>
                      </div>
                    </td>
                    <td className='max-w-xs px-6 py-4 border-b border-gray-200'>
                      <div className='text-sm leading-5 text-gray-700'>
                        <span className='flex items-center '>
                          /getting-started{' '}
                        </span>
                      </div>
                    </td>
                    <td className='max-w-xs px-6 py-4 border-b border-gray-200'>
                      <div className='text-sm leading-5 text-gray-700'>
                        <span className='flex items-center '>modified</span>
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
                        {/* )} */}
                      </span>
                    </td>
                  </tr>
                </tbody>
              </table>
              {/* empty state */}
              {/* <div className='w-full m-auto'>
                          <div className='flex flex-col items-center justify-center w-full h-full py-4 mr-32 bg-white rounded-lg shadow-lg md:flex-row'>
                            <div className='mr-6'>
                              <NoRepositories />
                            </div>
                            <div className='flex flex-col h-full'>
                              <p className='h-auto my-6 text-lg font-medium text-gray-800 '>
                                There are no matching members in this
                                organization.
                              </p>
                              <p className='text-base font-normal text-gray-700 '>
                                Not finding what you're looking for? If you have
                                access, try switching organizations via the top
                                left dropdown.
                              </p>
                              <p className='py-2 text-base font-normal text-gray-700 '>
                                Please contact your organizations admin for more
                                details.
                              </p>
                            </div>
                          </div>
                        </div> */}
              {/* empty state */}
            </div>
          </div>
        </div>
      </div>
    </SectionCard>
  );
};
