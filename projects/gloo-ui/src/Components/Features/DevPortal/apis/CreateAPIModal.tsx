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
import { Formik } from 'formik';
import {
  SoloFormInput,
  SoloFormTextarea
} from 'Components/Common/Form/SoloFormField';
import ImageUploader from 'react-images-upload';
import { colors } from 'Styles';
import { ReactComponent as NoImageIcon } from 'assets/no-image-placeholder.svg';
import styled from '@emotion/styled';
import tw from 'tailwind.macro';
import { ReactComponent as GreenPlus } from 'assets/small-green-plus.svg';
import { ReactComponent as NoSelectedList } from 'assets/no-selected-list.svg';
import { SoloButtonStyledComponent } from 'Styles/CommonEmotions/button';

export const SectionContainer = styled.div`
  ${tw`w-full h-full p-6 pb-0`}
`;

export const SectionHeader = styled.div`
  ${tw`mb-6 text-lg font-medium text-gray-800`}
`;

export const SectionSubHeader = styled.div`
  ${tw`mb-1 font-medium text-gray-800`}
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
      <SectionHeader>Create an API</SectionHeader>
      <div className='grid grid-cols-2 '>
        <div className='mr-4'>
          <SoloFormInput name='name' title='Name' />
        </div>
        <div>
          <SoloFormInput name='displayName' title='Display Name' />
        </div>
        <div className='col-span-2 '>
          <SoloFormTextarea name='description' title='Description' />
        </div>
      </div>
    </SectionContainer>
  );
};

const ImagerySection = () => {
  const [images, setImages] = React.useState<any[]>([]);

  const onDrop = (image: any) => {
    console.log('image', image);
    setImages([...images, image]);
  };
  return (
    <SectionContainer>
      <SectionHeader>Create an API: Add Imagery</SectionHeader>
      <div className='flex flex-col items-center p-4 pb-0 mr-4 bg-gray-100 border border-gray-400 rounded-lg'>
        <NoImageIcon />
        <ImageUploader
          css={css`
            .fileContainer {
              background: ${colors.januaryGrey};
              border-radius: 8px;
              box-shadow: none;
              padding: 0px;
            }
          `}
          withPreview
          withIcon={false}
          buttonStyles={{
            borderRadius: '8px',
            backgroundColor: colors.blue600,
            fontSize: '16px',
            fontWeight: 400
          }}
          buttonText='Upload an Image'
          onChange={onDrop}
          imgExtension={['.jpg', '.gif', '.png', '.gif', '.jpeg']}
          maxFileSize={5242880}
        />
      </div>
    </SectionContainer>
  );
};

const PortalsSection = () => {
  return (
    <SectionContainer>
      <SectionHeader>Portals</SectionHeader>
      <div className='grid grid-cols-2 col-gap-6'>
        <div>
          <div className='mb-1 font-medium text-gray-800'>
            Available Portals
          </div>
          <div className='w-full h-40 p-2 border border-gray-400 rounded-lg'>
            <div className='flex items-center justify-between'>
              Production Portal
              <span className='ml-1 text-green-400 cursor-pointer hover:text-green-300'>
                <GreenPlus className='w-4 h-4 fill-current' />
              </span>
            </div>
            <div>Production Portal</div>
          </div>
        </div>
        <div>
          <div className='mb-1 font-medium text-gray-800'>Selected Portals</div>
          <div className='w-full h-40 border border-gray-400 rounded-lg'>
            <div className='flex flex-col items-center justify-center w-full h-full bg-gray-100 rounded-lg'>
              <NoSelectedList className='w-12 h-12' />
              <div className='mt-2 text-gray-500'>Nothing Selected</div>
            </div>
            {/* <div className='flex items-center'>
              Production Portal
              <span className='ml-1 text-green-400 hover:text-green-300'>
                <GreenPlus className='w-4 h-4 fill-current' />
              </span>
            </div>
            <div>Production Portal</div> */}
          </div>
        </div>
      </div>
    </SectionContainer>
  );
};

const UserSection = () => {
  return (
    <SectionContainer>
      <SectionHeader>Users and Group Access</SectionHeader>
    </SectionContainer>
  );
};

const SpecSection = () => {
  return (
    <SectionContainer>
      <SectionHeader>Spec</SectionHeader>
      <div className='p-2 bg-gray-100 '>
        Upload or paste your swagger specs code below. Your API documentation
        will be automatically generated from the specs. You will be able to
        preview and modify your API docs from the API details page.
      </div>
      <div className='mb-1 font-medium text-gray-800'>Upload Swagger File</div>
      <div>file input</div>
      <div className='mb-1 font-medium text-gray-800'>Paste Swagger url</div>
      <div>file input</div>
    </SectionContainer>
  );
};

export const CreateAPIModal: React.FC<{ onClose: () => void }> = props => {
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
            displayName: '',
            version: '',
            description: '',
            name: ''
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
                <StyledTab>Imagery</StyledTab>
                <StyledTab>Portals</StyledTab>
                <StyledTab>Access</StyledTab>
                <StyledTab>Spec</StyledTab>
              </TabList>

              <TabPanels className='bg-white rounded-r-lg'>
                <TabPanel className='relative flex flex-col justify-between h-full focus:outline-none'>
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
                <TabPanel className='relative flex flex-col justify-between h-full focus:outline-none'>
                  <ImagerySection />
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
                <TabPanel className='relative flex flex-col justify-between h-full focus:outline-none'>
                  <PortalsSection />
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
                <TabPanel className='relative flex flex-col justify-between h-full focus:outline-none'>
                  <UserSection />
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
                <TabPanel className='relative flex flex-col justify-between h-full focus:outline-none'>
                  <SpecSection />
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
