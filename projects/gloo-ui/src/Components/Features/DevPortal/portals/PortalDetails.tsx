import React from 'react';
import { useParams, useHistory } from 'react-router';
import { Breadcrumb } from 'Components/Common/Breadcrumb';
import { SectionCard } from 'Components/Common/SectionCard';
import { ReactComponent as CodeIcon } from 'assets/code-icon.svg';
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
import { ReactComponent as PlaceholderPortal } from 'assets/placeholder-portal.svg';
import { ReactComponent as ExternalLinkIcon } from 'assets/external-link-icon.svg';
import { ErrorBoundary } from 'Components/Features/Errors/ErrorBoundary';
import { PortalPagesTab } from './PortalPagesTab';
import useSWR from 'swr';
import { portalApi } from '../api';
import { formatHealthStatus } from './PortalsListing';
import { Loading } from 'Components/Common/DisplayOnly/Loading';

export const TabCss = css`
  line-height: 40px;
  width: 80px;
  text-align: center;
  color: ${colors.septemberGrey};
  background: ${colors.februaryGrey};
  border: 1px solid ${colors.marchGrey};
  border-radius: ${soloConstants.radius}px ${soloConstants.radius}px 0 0;
  cursor: pointer;
  margin-right: 3px;
`;

export const ActiveTabCss = css`
  border-bottom: 1px solid white;
  color: ${colors.seaBlue};
  background: white;
  z-index: 2;
  cursor: default;
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
      css={css`
        ${TabCss}
        ${isSelected ? ActiveTabCss : ''}
      `}
      className='border rounded-lg focus:outline-none'>
      {children}
    </Tab>
  );
};

export const PortalDetails = () => {
  const { portalname, portalnamespace } = useParams();
  const { data: portal, error: portalListError } = useSWR(
    !!portalname && !!portalnamespace
      ? ['getPortal', portalname, portalnamespace]
      : null,
    (key, name, namespace) => portalApi.getPortal({ name, namespace })
  );

  const history = useHistory();
  const [tabIndex, setTabIndex] = React.useState(0);
  const [APISearchTerm, setAPISearchTerm] = React.useState('');

  const handleTabsChange = (index: number) => {
    setTabIndex(index);
  };
  if (!portal) {
    return <Loading center>Loading...</Loading>;
  }

  return (
    <ErrorBoundary
      fallback={<div>There was an error with the Dev Portal section</div>}>
      <div>
        <Breadcrumb />
        <SectionCard
          cardName={portalname || 'portal'}
          logoIcon={
            <span className='text-blue-500'>
              <CodeIcon className='fill-current' />
            </span>
          }
          health={formatHealthStatus(portal?.status?.state)}
          headerSecondaryInformation={[
            {
              title: 'Modified',
              value: 'Feb 26, 2020'
            }
          ]}
          healthMessage={'Portal Status'}
          onClose={() => history.push(`/dev-portal/`)}>
          <div>
            <div className='relative flex items-center'>
              <div>
                <PlaceholderPortal className='w-56 rounded-lg ' />
              </div>
              <div className='grid w-full grid-cols-2 ml-2 h-36'>
                <div>
                  <span className='font-medium text-gray-900'>
                    Portal Display Name
                  </span>
                  <div>{portal?.spec?.displayName}</div>
                </div>
                <div>
                  <span className='font-medium text-gray-900'>
                    Portal Address
                  </span>
                  <div className='flex items-center mb-2 text-sm text-blue-600'>
                    <span>
                      <ExternalLinkIcon className='w-4 h-4 ' />
                    </span>
                    https://production.subdomain.gloo.io
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
                  <div>{}</div>
                </div>
              </div>
            </div>
            <Tabs
              index={tabIndex}
              className='mt-6 mb-4 border-none rounded-lg'
              onChange={handleTabsChange}>
              <TabList className='flex items-start ml-4 '>
                <StyledTab>Theme</StyledTab>
                <StyledTab>Pages</StyledTab>
                <StyledTab>APIs</StyledTab>
                <StyledTab>Users</StyledTab>
                <StyledTab>Groups</StyledTab>
              </TabList>
              <TabPanels
                css={css`
                  margin-top: -1px;
                `}>
                <TabPanel className='focus:outline-none'>
                  <div className='relative flex flex-col p-4 border border-gray-300 rounded-lg'>
                    Theme Section
                  </div>
                </TabPanel>
                <TabPanel className='focus:outline-none'>
                  <PortalPagesTab />
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
                <TabPanel className='focus:outline-none'>
                  <div className='relative flex flex-col p-4 border border-gray-300 rounded-lg'>
                    Users Section
                  </div>
                </TabPanel>
                <TabPanel className='focus:outline-none'>
                  <div className='relative flex flex-col p-4 border border-gray-300 rounded-lg'>
                    Groups Section
                  </div>
                </TabPanel>
              </TabPanels>
            </Tabs>
            <button>Delete Portal</button>
          </div>
        </SectionCard>
      </div>
    </ErrorBoundary>
  );
};
