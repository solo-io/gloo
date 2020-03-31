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
import { PortalUsersTab } from './PortalUsersTab';
import useSWR from 'swr';
import { portalApi } from '../api';
import { formatHealthStatus } from './PortalsListing';
import { Loading } from 'Components/Common/DisplayOnly/Loading';
import { format } from 'timeago.js';
import { SoloModal } from 'Components/Common/SoloModal';
import { CreateAPIModal } from '../apis/CreateAPIModal';
import { ReactComponent as GreenPlus } from 'assets/small-green-plus.svg';
import { SoloNegativeButton } from 'Styles/CommonEmotions/button';
import { ConfirmationModal } from 'Components/Common/ConfirmationModal';
import { PortalApiDocsTab } from './PortalApiDocsTab';
import { PortalGroupsTab } from './PortalGroupsTab';

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
    (key, name, namespace) => portalApi.getPortalWithAssets({ name, namespace })
  );

  const history = useHistory();
  const [tabIndex, setTabIndex] = React.useState(0);
  const [attemptingDelete, setAttemptingDelete] = React.useState(false);

  const attemptDeletePortal = () => {
    setAttemptingDelete(true);
  };
  const cancelDeletion = () => {
    setAttemptingDelete(false);
  };
  const deletePortal = async () => {
    await portalApi.deletePortal({
      name: portal?.metadata?.name!,
      namespace: portal?.metadata?.namespace!
    });
    setAttemptingDelete(false);
    history.push('/dev-portal/portals');
  };
  const handleTabsChange = (index: number) => {
    setTabIndex(index);
  };
  if (!portal) {
    return <Loading center>Loading...</Loading>;
  }
  const domainsList = portal.spec?.domainsList.map(domain => {
    return {
      value: domain
    };
  });
  console.log('portal', portal);
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
              value: format(
                portal.metadata?.creationTimestamp?.seconds!,
                'en_US'
              )
            },
            ...domainsList!
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
                  <div className='break-words '>{portal.spec?.description}</div>
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
                  <PortalApiDocsTab portal={portal} />
                </TabPanel>
                <TabPanel className='focus:outline-none'>
                  <PortalUsersTab portal={portal} />
                </TabPanel>
                <TabPanel className='focus:outline-none'>
                  <div className='relative flex flex-col p-4 border border-gray-300 rounded-lg'>
                    <PortalGroupsTab portal={portal} />
                  </div>
                </TabPanel>
              </TabPanels>
            </Tabs>
            <div className='flex justify-end items-bottom'>
              <SoloNegativeButton onClick={attemptDeletePortal}>
                Delete Portal
              </SoloNegativeButton>
            </div>
            <ConfirmationModal
              visible={attemptingDelete}
              confirmationTopic='delete this portal'
              confirmText='Delete'
              goForIt={deletePortal}
              cancel={cancelDeletion}
              isNegative={true}
            />
          </div>
        </SectionCard>
      </div>
    </ErrorBoundary>
  );
};
