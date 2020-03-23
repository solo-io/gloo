import React from 'react';
import {
  Tabs,
  TabList,
  Tab,
  TabPanels,
  TabPanel,
  TabPanelProps
} from '@reach/tabs';
import { css } from '@emotion/core';

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

export function CreateAPIModal() {
  const [tabIndex, setTabIndex] = React.useState(0);

  const handleTabsChange = (index: number) => {
    setTabIndex(index);
  };
  return (
    <>
      <div className='bg-white rounded-lg shadow '>
        <Tabs
          className='bg-blue-600 rounded-lg h-96'
          index={tabIndex}
          onChange={handleTabsChange}
          css={css`
            display: grid;
            grid-template-columns: 190px 1fr;
          `}>
          <TabList className='flex flex-col mt-6'>
            <StyledTab>General</StyledTab>
            <StyledTab>Imagery</StyledTab>
            <StyledTab>Portals</StyledTab>
            <StyledTab>Access</StyledTab>
            <StyledTab>Spec</StyledTab>
          </TabList>

          <TabPanels className='bg-white rounded-r-lg'>
            <TabPanel className='focus:outline-none'>
              <div className='relative flex flex-col'>
                <div className='text-lg '>Create an API</div>
              </div>
            </TabPanel>
            <TabPanel className='focus:outline-none'>
              <div className='relative flex flex-col '>Add Imagery</div>
            </TabPanel>
            <TabPanel className='focus:outline-none'>
              <div className='relative flex flex-col '>Portals</div>
            </TabPanel>
            <TabPanel className='focus:outline-none'>
              <div className='relative flex flex-col '>
                Users and Groups Access
              </div>
            </TabPanel>
            <TabPanel className='focus:outline-none'>
              <div className='relative flex flex-col '>Spec</div>
            </TabPanel>
          </TabPanels>
        </Tabs>
      </div>
    </>
  );
}
