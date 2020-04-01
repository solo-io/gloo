import { css } from '@emotion/core';
import {
  Tab,
  TabList,
  TabPanel,
  TabPanelProps,
  TabPanels,
  Tabs
} from '@reach/tabs';
import { ReactComponent as CodeIcon } from 'assets/code-icon.svg';
import { ReactComponent as EditIcon } from 'assets/edit-pencil.svg';
import { ReactComponent as ExternalLinkIcon } from 'assets/external-link-icon.svg';
import { ReactComponent as PlaceholderPortal } from 'assets/placeholder-portal.svg';
import { Breadcrumb } from 'Components/Common/Breadcrumb';
import { ConfirmationModal } from 'Components/Common/ConfirmationModal';
import { SectionCard } from 'Components/Common/SectionCard';
import { SoloInput } from 'Components/Common/SoloInput';
import { ErrorBoundary } from 'Components/Features/Errors/ErrorBoundary';
import { State } from 'proto/dev-portal/api/dev-portal/v1/common_pb';
import React from 'react';
import { useHistory, useParams } from 'react-router';
import { SoloNegativeButton } from 'Styles/CommonEmotions/button';
import useSWR from 'swr';
import { apiDocApi } from '../api';
import { ActiveTabCss, TabCss } from '../portals/PortalDetails';
import { formatHealthStatus } from '../portals/PortalsListing';
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

export const APIDetails = () => {
  const { apiname, apinamespace } = useParams();
  const { data: apiDoc, error: apiDocError } = useSWR(
    !!apiname && !!apinamespace ? ['getApiDoc', apiname, apinamespace] : null,
    (key, name, namespace) =>
      apiDocApi.getApiDoc({ apidoc: { name, namespace }, withassets: true })
  );
  const history = useHistory();
  const [tabIndex, setTabIndex] = React.useState(0);
  const [APISearchTerm, setAPISearchTerm] = React.useState('');
  const [attemptingDelete, setAttemptingDelete] = React.useState(false);

  const attemptDeleteApiDoc = () => {
    setAttemptingDelete(true);
  };

  const cancelDeletion = () => {
    setAttemptingDelete(false);
  };

  const deleteApi = async () => {
    await apiDocApi.deleteApiDoc({
      name: apiDoc?.metadata?.name!,
      namespace: apiDoc?.metadata?.namespace!
    });
    setAttemptingDelete(false);
    history.push('/dev-portal/portals');
  };

  const handleTabsChange = (index: number) => {
    setTabIndex(index);
  };

  if (!apiDoc) {
    return <Loading center>Loading...</Loading>;
  }
  return (
    <ErrorBoundary
      fallback={<div>There was an error with the Dev Portal section</div>}>
      <div>
        <Breadcrumb />
        <SectionCard
          cardName={apiname || 'API'}
          logoIcon={
            <span className='text-blue-500'>
              <CodeIcon className='fill-current' />
            </span>
          }
          health={formatHealthStatus(apiDoc?.status?.state)}
          headerSecondaryInformation={[
            {
              title: 'Modified',
              value: 'Feb 26, 2020'
            }
          ]}
          healthMessage={'Portal Status'}
          onClose={() => history.push(`/dev-portal/`)}>
          <div>
            {(apiDoc?.status?.state === State.PENDING ||
              apiDoc?.status?.state === State.PROCESSING) && (
              <div className='p-2 text-yellow-500 bg-yellow-100 rounded-lg '>
                Updates are pending publication
              </div>
            )}
            <div className='relative flex items-center'>
              <div className=' max-h-72'>
                {apiDoc.spec?.image ? (
                  <img
                    className='object-cover max-h-72'
                    src={`data:image/gif;base64,${apiDoc.spec?.image?.inlineBytes}`}></img>
                ) : (
                  <PlaceholderPortal className='w-56 rounded-lg ' />
                )}
              </div>
              <div className='grid w-full grid-cols-2 ml-2 h-36'>
                <div>
                  <span className='font-medium text-gray-900'>
                    Display Name
                  </span>
                  <div>
                    {apiDoc?.status?.displayName || apiDoc?.metadata?.name}
                  </div>
                </div>
                <div>
                  <span className='font-medium text-gray-900'>
                    Published In
                  </span>
                  <div className='flex items-center mb-2 text-sm text-blue-600'>
                    <span>
                      <ExternalLinkIcon className='w-4 h-4 ' />
                    </span>
                    {apiDoc.status?.basePath}
                  </div>
                </div>
                <span className='absolute top-0 right-0 flex items-center'>
                  <span className='mr-2'> Edit</span>
                  <span className='flex items-center justify-center w-6 h-6 mr-3 text-gray-700 bg-gray-400 rounded-full cursor-pointer'>
                    <EditIcon className='w-3 h-3' />
                  </span>
                </span>
                <div className='col-span-2 '>
                  <span className='font-medium text-gray-900'>Description</span>
                  <div>
                    Updated - Lorem ipsum dolor sit amet, consetetur sadipscing
                    elitr, sed diam nonumy eirmod tempor invidunt ut labore et
                    dolore magna aliquyam erat, sed diam voluptua. At vero eos
                    et accusam et justo duo dolores et ea rebum. Stet clita kasd
                    gubergren, no sea takimata sanctus est Lorem ipsum dolor sit
                    amet.
                  </div>
                </div>
              </div>
            </div>
            <Tabs
              index={tabIndex}
              className='mt-6 mb-4 border-none rounded-lg'
              onChange={handleTabsChange}>
              <TabList className='flex items-start ml-4 '>
                <StyledTab>Users</StyledTab>
                <StyledTab>Groups</StyledTab>
              </TabList>
              <TabPanels
                css={css`
                  margin-top: -1px;
                `}>
                <TabPanel className='focus:outline-none'>
                  <div className='relative flex flex-col p-4 border border-gray-300 rounded-lg'>
                    <span className='absolute top-0 right-0 p-4 '>
                      <span></span> add an API
                    </span>
                    <div className='w-1/3 m-4'>
                      <SoloInput
                        placeholder='Search by API name...'
                        value={APISearchTerm}
                        onChange={e => setAPISearchTerm(e.target.value)}
                      />
                    </div>
                    <div className='flex flex-col'>
                      <div className='py-2 -my-2 overflow-x-auto sm:-mx-6 sm:px-6 lg:-mx-8 lg:px-8'>
                        <div className='inline-block min-w-full overflow-hidden align-middle border-b border-gray-200 shadow sm:rounded-lg'>
                          <table className='min-w-full'>
                            <thead className='bg-gray-300 '>
                              <tr>
                                <th className='px-6 py-3 text-sm font-medium leading-4 tracking-wider text-left text-gray-800 capitalize border-b border-gray-200 bg-gray-50'>
                                  API Name
                                </th>
                                <th className='px-6 py-3 text-sm font-medium leading-4 tracking-wider text-left text-gray-800 capitalize border-b border-gray-200 bg-gray-50'>
                                  Description
                                </th>

                                <th className='px-6 py-3 text-sm font-medium leading-4 tracking-wider text-left text-gray-800 capitalize border-b border-gray-200 bg-gray-50'>
                                  Modified
                                </th>
                                <th className='px-6 py-3 text-sm font-medium leading-4 tracking-wider text-left text-gray-800 capitalize border-b border-gray-200 bg-gray-50'>
                                  Status
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
                                <td className='px-6 py-4 whitespace-no-wrap border-b border-gray-200'>
                                  <div className='text-sm leading-5 text-gray-900'>
                                    <span className='flex items-center capitalize'>
                                      Getting Started
                                    </span>
                                  </div>
                                </td>
                                <td className='px-6 py-4 whitespace-no-wrap border-b border-gray-200'>
                                  <div className='text-sm leading-5 text-gray-900'>
                                    <span className='flex items-center capitalize'>
                                      Getting Started
                                    </span>
                                  </div>
                                </td>
                                <td className='px-6 py-4 whitespace-no-wrap border-b border-gray-200'>
                                  <div className='text-sm leading-5 text-gray-900'>
                                    <span className='flex items-center '>
                                      /getting-started
                                    </span>
                                  </div>
                                </td>
                                <td className='px-6 py-4 whitespace-no-wrap border-b border-gray-200'>
                                  <div className='text-sm leading-5 text-gray-900'>
                                    <span className='flex items-center capitalize'>
                                      modified
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
                </TabPanel>
                <TabPanel className='focus:outline-none'>
                  <div className='relative flex flex-col p-4 border border-gray-300 rounded-lg'>
                    <span className='absolute top-0 right-0 p-4 '>
                      <span></span> add an API
                    </span>
                    <div className='w-1/3 m-4'>
                      <SoloInput
                        placeholder='Search by API name...'
                        value={APISearchTerm}
                        onChange={e => setAPISearchTerm(e.target.value)}
                      />
                    </div>
                    <div className='flex flex-col'>
                      <div className='py-2 -my-2 overflow-x-auto sm:-mx-6 sm:px-6 lg:-mx-8 lg:px-8'>
                        <div className='inline-block min-w-full overflow-hidden align-middle border-b border-gray-200 shadow sm:rounded-lg'>
                          <table className='min-w-full'>
                            <thead className='bg-gray-300 '>
                              <tr>
                                <th className='px-6 py-3 text-sm font-medium leading-4 tracking-wider text-left text-gray-800 capitalize border-b border-gray-200 bg-gray-50'>
                                  API Name
                                </th>
                                <th className='px-6 py-3 text-sm font-medium leading-4 tracking-wider text-left text-gray-800 capitalize border-b border-gray-200 bg-gray-50'>
                                  Description
                                </th>

                                <th className='px-6 py-3 text-sm font-medium leading-4 tracking-wider text-left text-gray-800 capitalize border-b border-gray-200 bg-gray-50'>
                                  Modified
                                </th>
                                <th className='px-6 py-3 text-sm font-medium leading-4 tracking-wider text-left text-gray-800 capitalize border-b border-gray-200 bg-gray-50'>
                                  Status
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
                                <td className='px-6 py-4 whitespace-no-wrap border-b border-gray-200'>
                                  <div className='text-sm leading-5 text-gray-900'>
                                    <span className='flex items-center capitalize'>
                                      Getting Started
                                    </span>
                                  </div>
                                </td>
                                <td className='px-6 py-4 whitespace-no-wrap border-b border-gray-200'>
                                  <div className='text-sm leading-5 text-gray-900'>
                                    <span className='flex items-center capitalize'>
                                      Getting Started
                                    </span>
                                  </div>
                                </td>
                                <td className='px-6 py-4 whitespace-no-wrap border-b border-gray-200'>
                                  <div className='text-sm leading-5 text-gray-900'>
                                    <span className='flex items-center '>
                                      /getting-started
                                    </span>
                                  </div>
                                </td>
                                <td className='px-6 py-4 whitespace-no-wrap border-b border-gray-200'>
                                  <div className='text-sm leading-5 text-gray-900'>
                                    <span className='flex items-center capitalize'>
                                      modified
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
                </TabPanel>
              </TabPanels>
            </Tabs>
            <div className='flex justify-end items-bottom'>
              <SoloNegativeButton onClick={attemptDeleteApiDoc}>
                Delete API
              </SoloNegativeButton>
            </div>
            <ConfirmationModal
              visible={attemptingDelete}
              confirmationTopic='delete this API'
              confirmText='Delete'
              goForIt={deleteApi}
              cancel={cancelDeletion}
              isNegative={true}
            />
          </div>
        </SectionCard>
      </div>
    </ErrorBoundary>
  );
};
