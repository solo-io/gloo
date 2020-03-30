import React from 'react';

import {
  Tabs,
  TabList,
  Tab,
  TabPanels,
  TabPanel,
  TabPanelProps
} from '@reach/tabs';
import { ReactComponent as GreenPlus } from 'assets/small-green-plus.svg';
import { css } from '@emotion/core';
import { Formik } from 'formik';
import { SectionContainer, SectionHeader } from '../apis/CreateAPIModal';
import { SoloFormInput } from 'Components/Common/Form/SoloFormField';
import { SoloButtonStyledComponent } from 'Styles/CommonEmotions/button';

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

const GeneralSection = () => {
  return (
    <SectionContainer>
      <SectionHeader> Create User</SectionHeader>
      <div className='grid grid-cols-2'>
        <div>
          <div className='mr-4'>
            <SoloFormInput
              name='name'
              title='Name'
              placeholder='Username goes here'
              hideError
            />
          </div>
          <div className='mr-4'>
            <SoloFormInput
              name='email'
              title='Email'
              placeholder='email@domain.com'
              hideError
            />
          </div>
          <div className='mr-4'>
            <SoloFormInput
              name='password'
              title='Password'
              placeholder='type password here'
              hideError
            />
          </div>
        </div>
        <div></div>
      </div>
    </SectionContainer>
  );
};

const ApiSection = () => {
  return <div>API Section</div>;
};

const PortalsSection = () => {
  return <div>Portals Section</div>;
};

export const CreateGroupModal: React.FC<{ onClose: () => void }> = props => {
  const [tabIndex, setTabIndex] = React.useState(0);
  const handleTabsChange = (index: number) => {
    setTabIndex(index);
  };
  return (
    <>
      <div
        css={css`
          width: 750px;
        `}
        className='bg-white rounded-lg shadow '>
        <Formik
          initialValues={{
            name: '',
            email: '',
            password: '',
            group: '',
            options: ''
          }}
          onSubmit={() => {}}>
          <>
            <Tabs
              className='bg-blue-600 rounded-lg h-96'
              index={tabIndex}
              onChange={handleTabsChange}
              css={css`
                display: grid;
                height: 450px;
                grid-template-columns: 190px 1fr;
              `}>
              <TabList className='flex flex-col mt-6'>
                <StyledTab>General</StyledTab>
                <StyledTab>APIs</StyledTab>
                <StyledTab>Portals</StyledTab>
              </TabList>

              <TabPanels className='bg-white rounded-r-lg'>
                <TabPanel className='flex flex-col justify-between h-full focus:outline-none'>
                  <GeneralSection />
                  <div className='flex items-end justify-between h-full px-6 mb-4 '>
                    <button
                      className='text-blue-500 cursor-pointer'
                      onClick={props.onClose}>
                      cancel
                    </button>
                    <SoloButtonStyledComponent
                      onClick={() => setTabIndex(tabIndex + 1)}>
                      Next Step
                    </SoloButtonStyledComponent>
                  </div>
                </TabPanel>
                <TabPanel className='flex flex-col justify-between h-full focus:outline-none'>
                  <ApiSection />
                  <div className='flex items-end justify-between h-full px-6 mb-4 '>
                    <button
                      className='text-blue-500 cursor-pointer'
                      onClick={props.onClose}>
                      cancel
                    </button>
                    <SoloButtonStyledComponent
                      onClick={() => setTabIndex(tabIndex + 1)}>
                      Next Step
                    </SoloButtonStyledComponent>
                  </div>
                </TabPanel>
                <TabPanel className='flex flex-col justify-between h-full focus:outline-none'>
                  <PortalsSection />
                  <div className='flex items-end justify-between h-full px-6 mb-4 '>
                    <button
                      className='text-blue-500 cursor-pointer'
                      onClick={props.onClose}>
                      cancel
                    </button>
                    <SoloButtonStyledComponent onClick={() => setTabIndex(0)}>
                      Create Portal
                    </SoloButtonStyledComponent>
                  </div>
                </TabPanel>
              </TabPanels>
            </Tabs>
          </>
        </Formik>
      </div>
    </>
  );
};
