import React from 'react';
import {
  Tabs,
  TabList,
  Tab,
  TabPanels,
  TabPanel,
  TabsProps,
  TabPanelProps
} from '@reach/tabs';
import { css } from '@emotion/core';
import { ReactComponent as DevPortalIcon } from 'assets/developer-portal-icon.svg';
import { ReactComponent as PlaceholderPortalTile } from 'assets/portal-tile.svg';
import { ReactComponent as PlaceholderAPITile } from 'assets/api-tile.svg';
import { SoloButton } from 'Components/Common/SoloButton';
import { InverseButtonCSS } from 'Styles/CommonEmotions/button';
import { colors } from 'Styles';
import { Breadcrumb } from 'Components/Common/Breadcrumb';
import { PortalsListing } from './portals/PortalsListing';
import { CreateAPIModal } from './apis/CreateAPIModal';
import { UserGroups } from './users/UserGroups';
import { APIListing } from './apis/APIListing';
import { useHistory } from 'react-router-dom';
import { useLocation } from 'react-router';
import { ErrorBoundary } from '../Errors/ErrorBoundary';

export const EmptyPortalsPanel: React.FC<{ itemName: string }> = props => {
  const { itemName = 'Portal' } = props;
  return (
    <>
      <div className='w-full '>
        <div className='flex items-center justify-center mx-4 mb-2 bg-white rounded-lg shadow'>
          <span className='text-blue-500 '>
            <DevPortalIcon className='w-32 h-32 fill-current' />
          </span>
          <div className='ml-8'>
            <div className='text-lg font-medium text-gray-900'>
              {` There are no ${itemName}s to display.`}
            </div>
            <div>
              <a>Create a {itemName}</a> to publish your APIs to and share with
              developers.
            </div>
          </div>
        </div>
        {props.children}
      </div>
    </>
  );
};

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
      className={`p-1 w-48 text-left pl-4 border rounded-lg mb-4 focus:outline-none ${
        isSelected
          ? 'text-blue-600 bg-blue-100 border-blue-600 '
          : 'text-gray-700 bg-white border-gray-600'
      }`}>
      {children}
    </Tab>
  );
};

const routesMap: { [key: string]: number } = {
  '/dev-portal/portals': 0,
  '/dev-portal/apis': 1,
  '/dev-portal/users': 2,
  '/dev-portal/api-key-scopes': 3,
  '/dev-portal/api-keys': 4
};

export function DevPortal() {
  const history = useHistory();
  const location = useLocation();
  const [tabIndex, setTabIndex] = React.useState(routesMap[location.pathname]);

  const handleTabsChange = (index: number) => {
    setTabIndex(index);
    let newRoute = Object.entries(routesMap).find(
      ([path, mapping]) => mapping === index
    )!;
    history.push(newRoute[0]);
  };

  React.useEffect(() => {
    if (tabIndex !== routesMap[location.pathname]) {
      setTabIndex(routesMap[location.pathname]);
    }
  }, [location.pathname]);

  return (
    <ErrorBoundary
      fallback={<div>There was an error with the Dev Portal section</div>}>
      <Breadcrumb />
      <Tabs
        index={tabIndex}
        onChange={handleTabsChange}
        css={css`
          display: grid;
          grid-template-columns: 190px 1fr;
        `}>
        <TabList className='flex flex-col items-start'>
          <StyledTab>Portals</StyledTab>
          <StyledTab>APIs</StyledTab>
          <StyledTab>Users & Groups</StyledTab>
          <StyledTab>API Key Scopes</StyledTab>
          <StyledTab>API Keys</StyledTab>
        </TabList>

        <TabPanels
          css={css`
            margin-left: 30px;
          `}>
          <TabPanel className='focus:outline-none'>
            <div className='relative flex flex-col '>
              <span className='absolute top-0 right-0 -mt-8'>
                <span></span> create a portal
              </span>
              <PortalsListing />
            </div>
          </TabPanel>
          <TabPanel className='focus:outline-none'>
            <div className='relative flex flex-col '>
              <span className='absolute top-0 right-0 -mt-8'>
                <span></span> create an API
              </span>
              <APIListing />
            </div>
          </TabPanel>
          <TabPanel className='focus:outline-none'>
            <div className='relative flex flex-col '>
              <span className='absolute top-0 right-0 -mt-8'>
                <span> Create a User</span> <span>Create a Group</span>
              </span>
              <UserGroups />
              {/* <EmptyPortalsPanel itemName=''></EmptyPortalsPanel> */}
            </div>
          </TabPanel>
          <TabPanel className='focus:outline-none'>
            <div className='relative flex flex-col '>
              <span className='absolute top-0 right-0 -mt-8'>
                <span></span> create an API Key Scope
              </span>
              <EmptyPortalsPanel itemName=''></EmptyPortalsPanel>
            </div>
          </TabPanel>
          <TabPanel className='focus:outline-none'>
            <div className='relative flex flex-col '>
              <span className='absolute top-0 right-0 -mt-8'>
                <span></span> create an API Key
              </span>
              <EmptyPortalsPanel itemName=''></EmptyPortalsPanel>
            </div>
          </TabPanel>
        </TabPanels>
      </Tabs>
    </ErrorBoundary>
  );
}
