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
import { APIUsersTab } from './APIUsers';
import { APIGroupsTab } from './APIGroups';

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
            {apiDoc?.status?.state !== State.SUCCEEDED && (
              <div className='flex items-center p-2 mb-2 text-yellow-500 bg-yellow-100 border border-yellow-500 rounded-lg '>
                <div className='flex items-center justify-center w-4 h-4 mr-2 text-white text-yellow-500 bg-orange-100 border border-yellow-500 rounded-full'>
                  !
                </div>{' '}
                {apiDoc.status?.reason}
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
                  {/* {apiDoc..spec?.domainsList.map((domain, index) => (
                    <div
                      key={domain}
                      className='flex items-center mb-2 text-sm text-blue-600'>
                      <span>
                        <ExternalLinkIcon className='w-4 h-4 ' />
                      </span>
                    
                        <div>{domain}</div>
                     </div>
                  ))} */}
                </div>
                <span className='absolute top-0 right-0 flex items-center'>
                  <span className='mr-2'> Edit</span>
                  <span className='flex items-center justify-center w-6 h-6 mr-3 text-gray-700 bg-gray-400 rounded-full cursor-pointer'>
                    <EditIcon className='w-3 h-3' />
                  </span>
                </span>
                <div className='col-span-2 '>
                  <span className='font-medium text-gray-900'>Description</span>
                  <div>{apiDoc.status?.description}</div>
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
                  <APIUsersTab apiDoc={apiDoc} />
                </TabPanel>
                <TabPanel className='focus:outline-none'>
                  <APIGroupsTab apiDoc={apiDoc} />
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
