import { css } from '@emotion/core';
import {
  Tab,
  TabList,
  TabPanel,
  TabPanelProps,
  TabPanels,
  Tabs
} from '@reach/tabs';
import { ReactComponent as DevPortalIcon } from 'assets/developer-portal-icon.svg';
import { Breadcrumb } from 'Components/Common/Breadcrumb';
import { APIKeyScopes } from 'Components/Features/DevPortal/api-key-scopes/ApiKeyScopes';
import React from 'react';
import { useLocation } from 'react-router';
import { useHistory } from 'react-router-dom';
import { ErrorBoundary } from '../Errors/ErrorBoundary';
import { APIKeys } from './api-keys/ApiKeys';
import { APIListing } from './apis/APIListing';
import { PortalsListing } from './portals/PortalsListing';
import { UserGroups } from './users/UserGroups';

export const NoDataPanel: React.FC<{
  missingContentText: string;
  helpText: string;
  identifier: string;
}> = ({ missingContentText, helpText, identifier, ...restProps }) => {
  return (
    <>
      <div className='w-full '>
        <div className='flex items-center justify-center mx-4 mb-2 bg-white rounded-lg shadow'>
          <span className='text-blue-500 '>
            <DevPortalIcon className='w-32 h-32 fill-current' />
          </span>
          <div className='ml-8'>
            <div className='text-lg font-medium text-gray-900'>
              {missingContentText}
            </div>
            <div>{helpText}</div>
          </div>
        </div>
        <svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 1090 254.712'>
          <defs>
            <clipPath id={`portal-tile-${identifier}-c`}>
              <path
                d='M9 0h216v208H9a9 9 0 01-9-9V9a9 9 0 019-9z'
                fill='#d4d8de'
              />
            </clipPath>
            <filter
              id={`portal-tile-${identifier}-a`}
              width='1090'
              height='254.712'
              x='0'
              y='0'
              filterUnits='userSpaceOnUse'>
              <feOffset />
              <feGaussianBlur result='b' stdDeviation='7.5' />
              <feFlood floodOpacity='.063' />
              <feComposite in2='b' operator='in' />
              <feComposite in='SourceGraphic' />
            </filter>
          </defs>
          <g opacity='.502'>
            <g filter={`url(#portal-tile-${identifier}-a)`}>
              <g fill='#fff' stroke='#d4d8de' transform='translate(22.5 22.5)'>
                <rect width='1045' height='209.712' stroke='none' rx='10' />
                <rect
                  width='1044'
                  height='208.712'
                  x='.5'
                  y='.5'
                  fill='none'
                  rx='9.5'
                />
              </g>
            </g>
            <path
              d='M274.944 107.811H1047v12.317H274.944zm0 26.569H1047v12.317H274.944zm646-86.569H1021v12.317H920.944z'
              fill='#d4d8de'
            />
            <path fill='#2196c9' d='M274.5 77.5h272v15h-272z' />
            <path d='M274.5 47.5h272v15h-272z' fill='#6e7477' />
            <path d='M865.5 188.5h81v15h-81z' fill='#d4d8de' />
            <circle
              cx='9'
              cy='9'
              r='9'
              fill='#d4d8de'
              transform='translate(1029 44.5)'
            />
            <path
              d='M966.5 188.5h81v15h-81zM33.456 23.5H248.5v209H33.456a9.978 9.978 0 01-9.956-10v-189a9.978 9.978 0 019.956-10z'
              fill='#d4d8de'
            />
            <g
              clipPath={`url(#portal-tile-${identifier}-c)`}
              transform='translate(23.5 23.5)'>
              <g transform='translate(-719.137 2)'>
                <path fill='#c0cbd3' d='M674.137-2h312v208h-312z' />
                <circle
                  cx='12.94'
                  cy='12.94'
                  r='12.94'
                  fill='#f1f3f5'
                  transform='translate(751.736 74.304)'
                />
                <path
                  d='M674.136 128.54s15.24-20.013 32.664-19.778c15.444.208 25.949 6.947 42.773 21.189s30.517 24.98 44.789 6.207S838 71.996 851.205 62.065c13.062-9.823 33.768-19.475 61.827-3.54 33.5 19.025 73.393 78.34 73.393 78.34v69.327H674.136z'
                  fill='#f1f3f5'
                />
              </g>
            </g>
            <g>
              <path d='M327.95 180.38H432v22.317H327.95z' fill='#d4d8de' />
              <path
                d='M294.53 209.237a14.637 14.637 0 01-14.32-24.8 14.62 14.62 0 003.971 13.383c.183.183.369.36.56.533.054.251.115.5.183.756a14.617 14.617 0 009.606 10.128zm8.319.22a14.547 14.547 0 01-8.318-.22 14.619 14.619 0 009.6-10.128c.069-.251.129-.5.181-.754.193-.173.379-.35.564-.535a14.62 14.62 0 003.969-13.383 14.637 14.637 0 01-6 25.02z'
                fill='#adb3bc'
              />
              <path
                d='M304.315 198.356c-.052.251-.113.5-.18.754a9.619 9.619 0 11-19.21 0 14.54 14.54 0 01-.183-.756 14.635 14.635 0 0019.573.002z'
                fill='#6e7477'
              />
              <path
                d='M308.847 184.437a14.621 14.621 0 00-13.574-3.253c-.251.068-.5.14-.743.222-.245-.08-.494-.153-.745-.222a14.625 14.625 0 00-13.576 3.253 14.637 14.637 0 0128.637 0z'
                fill='#adb3bc'
              />
              <path
                d='M284.74 198.354c-.19-.173-.377-.35-.56-.533a9.619 9.619 0 119.606-16.636c.25.068.5.142.744.222a14.636 14.636 0 00-9.79 16.947zm20.139-.533c-.185.185-.371.362-.564.535a14.635 14.635 0 00-9.784-16.949c.244-.082.492-.154.743-.222a9.618 9.618 0 119.605 16.636z'
                fill='#6e7477'
              />
            </g>
          </g>
        </svg>
        <svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 1090 254.712'>
          <g opacity='.502'>
            <g filter={`url(#portal-tile-${identifier}-a)`}>
              <g fill='#fff' stroke='#d4d8de' transform='translate(22.5 22.5)'>
                <rect width='1045' height='209.712' stroke='none' rx='10' />
                <rect
                  width='1044'
                  height='208.712'
                  x='.5'
                  y='.5'
                  fill='none'
                  rx='9.5'
                />
              </g>
            </g>
            <path
              d='M274.944 107.811H1047v12.317H274.944zm0 26.569H1047v12.317H274.944zm646-86.569H1021v12.317H920.944z'
              fill='#d4d8de'
            />
            <path fill='#2196c9' d='M274.5 77.5h272v15h-272z' />
            <path d='M274.5 47.5h272v15h-272z' fill='#6e7477' />
            <path d='M865.5 188.5h81v15h-81z' fill='#d4d8de' />
            <circle
              cx='9'
              cy='9'
              r='9'
              fill='#d4d8de'
              transform='translate(1029 44.5)'
            />
            <path
              d='M966.5 188.5h81v15h-81zM33.456 23.5H248.5v209H33.456a9.978 9.978 0 01-9.956-10v-189a9.978 9.978 0 019.956-10z'
              fill='#d4d8de'
            />
            <g
              clipPath={`url(#portal-tile-${identifier}-c)`}
              transform='translate(23.5 23.5)'>
              <g transform='translate(-719.137 2)'>
                <path fill='#c0cbd3' d='M674.137-2h312v208h-312z' />
                <circle
                  cx='12.94'
                  cy='12.94'
                  r='12.94'
                  fill='#f1f3f5'
                  transform='translate(751.736 74.304)'
                />
                <path
                  d='M674.136 128.54s15.24-20.013 32.664-19.778c15.444.208 25.949 6.947 42.773 21.189s30.517 24.98 44.789 6.207S838 71.996 851.205 62.065c13.062-9.823 33.768-19.475 61.827-3.54 33.5 19.025 73.393 78.34 73.393 78.34v69.327H674.136z'
                  fill='#f1f3f5'
                />
              </g>
            </g>
            <g>
              <path d='M327.95 180.38H432v22.317H327.95z' fill='#d4d8de' />
              <path
                d='M294.53 209.237a14.637 14.637 0 01-14.32-24.8 14.62 14.62 0 003.971 13.383c.183.183.369.36.56.533.054.251.115.5.183.756a14.617 14.617 0 009.606 10.128zm8.319.22a14.547 14.547 0 01-8.318-.22 14.619 14.619 0 009.6-10.128c.069-.251.129-.5.181-.754.193-.173.379-.35.564-.535a14.62 14.62 0 003.969-13.383 14.637 14.637 0 01-6 25.02z'
                fill='#adb3bc'
              />
              <path
                d='M304.315 198.356c-.052.251-.113.5-.18.754a9.619 9.619 0 11-19.21 0 14.54 14.54 0 01-.183-.756 14.635 14.635 0 0019.573.002z'
                fill='#6e7477'
              />
              <path
                d='M308.847 184.437a14.621 14.621 0 00-13.574-3.253c-.251.068-.5.14-.743.222-.245-.08-.494-.153-.745-.222a14.625 14.625 0 00-13.576 3.253 14.637 14.637 0 0128.637 0z'
                fill='#adb3bc'
              />
              <path
                d='M284.74 198.354c-.19-.173-.377-.35-.56-.533a9.619 9.619 0 119.606-16.636c.25.068.5.142.744.222a14.636 14.636 0 00-9.79 16.947zm20.139-.533c-.185.185-.371.362-.564.535a14.635 14.635 0 00-9.784-16.949c.244-.082.492-.154.743-.222a9.618 9.618 0 119.605 16.636z'
                fill='#6e7477'
              />
            </g>
          </g>
        </svg>
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
              <PortalsListing />
            </div>
          </TabPanel>
          <TabPanel className='focus:outline-none'>
            <div className='relative flex flex-col '>
              <APIListing />
            </div>
          </TabPanel>
          <TabPanel className='focus:outline-none'>
            <div className='relative flex flex-col '>
              <UserGroups />
            </div>
          </TabPanel>
          <TabPanel className='focus:outline-none'>
            <div className='relative flex flex-col '>
              <APIKeyScopes />
            </div>
          </TabPanel>
          <TabPanel className='focus:outline-none'>
            <div className='relative flex flex-col '>
              <APIKeys />
            </div>
          </TabPanel>
        </TabPanels>
      </Tabs>
    </ErrorBoundary>
  );
}
