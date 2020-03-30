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
import { Formik, useField, useFormikContext } from 'formik';
import {
  SectionContainer,
  SectionHeader,
  SectionSubHeader
} from '../apis/CreateAPIModal';
import {
  SoloFormInput,
  SoloFormTextarea
} from 'Components/Common/Form/SoloFormField';
import ImageUploader from 'react-images-upload';
import { colors } from 'Styles';
import { ReactComponent as NoImageIcon } from 'assets/no-image-placeholder.svg';

import { ReactComponent as NoSelectedList } from 'assets/no-selected-list.svg';
import { SoloButtonStyledComponent } from 'Styles/CommonEmotions/button';
import { portalApi, portalMessageFromObject } from '../api';
import { Portal } from 'proto/dev-portal/api/grpc/admin/portal_pb';
import {
  PortalSpec,
  PortalStatus
} from 'proto/dev-portal/api/dev-portal/v1/portal_pb';

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
      <SectionHeader>Create an Portal</SectionHeader>
      <div className='grid grid-cols-2 '>
        <div className='mr-4'>
          <SoloFormInput
            name='name'
            title='Name'
            placeholder='Portal name goes here'
          />
        </div>
        <div>
          <SoloFormInput
            name='portalAddress'
            title='Portal Address'
            placeholder='https://subdomain.domain.io'
          />
        </div>
        <div className='col-span-2 '>
          <SoloFormTextarea
            name='description'
            title='Description'
            placeholder='Portal description goes here'
          />
        </div>
      </div>
    </SectionContainer>
  );
};

const ImagerySection = () => {
  const formik = useFormikContext();
  const [images, setImages] = React.useState<string[]>([]);

  const onDrop = (files: File[], pictures: string[]) => {
    console.log('files', files);
    console.log('pictures', pictures);
    formik.setFieldValue('banner', pictures[0]);
    setImages([...images, ...pictures]);
  };
  return (
    <SectionContainer>
      <SectionHeader>Create a Portal: Add Imagery</SectionHeader>
      <div className='flex flex-col items-center p-4 pb-0 mr-4 bg-gray-100 border border-gray-400 rounded-lg'>
        <NoImageIcon />
        <ImageUploader
          css={css`
            .fileContainer {
              background: ${colors.januaryGrey};
              border-radius: 8px;
              text-align: center;

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

const BrandingSection = () => {
  const formik = useFormikContext();
  const [images, setImages] = React.useState<string[]>([]);

  const onDropPrimaryLogo = (files: File[], pictures: string[]) => {
    formik.setFieldValue('primaryLogo', pictures[0]);
    setImages([...images, ...pictures]);
  };
  const onDropFavicon = (files: File[], pictures: string[]) => {
    formik.setFieldValue('favicon', pictures[0]);
    setImages([...images, ...pictures]);
  };
  return (
    <SectionContainer>
      <SectionHeader>Create a Portal: Branding Logos</SectionHeader>
      <div className='flex items-center'>
        <div className='flex flex-col items-center'>
          <SectionSubHeader>Primary Logo</SectionSubHeader>
          <div className='flex flex-col items-center p-4 pb-0 mr-4 bg-gray-100 border border-gray-400 rounded-lg'>
            <NoImageIcon />
            <ImageUploader
              css={css`
                .fileContainer {
                  background: ${colors.januaryGrey};
                  border-radius: 8px;
                  box-shadow: none;
                  text-align: center;
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
              onChange={onDropPrimaryLogo}
              imgExtension={['.jpg', '.png']}
              maxFileSize={5242880}
            />
          </div>
        </div>
        <div className='flex flex-col items-center'>
          <SectionSubHeader>Favicon</SectionSubHeader>

          <div className='flex flex-col items-center p-4 pb-0 mr-4 bg-gray-100 border border-gray-400 rounded-lg'>
            <NoImageIcon />
            <ImageUploader
              css={css`
                .fileContainer {
                  background: ${colors.januaryGrey};
                  border-radius: 8px;
                  box-shadow: none;
                  text-align: center;

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
              onChange={onDropFavicon}
              imgExtension={['.jpg', '.png']}
              maxFileSize={5242880}
            />
          </div>
        </div>
      </div>
    </SectionContainer>
  );
};

const APIsSection = () => {
  return (
    <SectionContainer>
      <SectionHeader>Create a Portal: APIs</SectionHeader>
    </SectionContainer>
  );
};

const AccessSection = () => {
  return (
    <SectionContainer>
      <SectionHeader>Create a Portal: Access</SectionHeader>
    </SectionContainer>
  );
};

export const CreatePortalModal: React.FC<{ onClose: () => void }> = props => {
  const [tabIndex, setTabIndex] = React.useState(0);

  const handleTabsChange = (index: number) => {
    setTabIndex(index);
  };

  // displayName: string,
  //   description: string,
  //   domainsList: Array<string>,
  //   primaryLogo?: dev_portal_api_dev_portal_v1_common_pb.DataSource.AsObject, asdf
  //   favicon?: dev_portal_api_dev_portal_v1_common_pb.DataSource.AsObject, asdf
  //   banner?: dev_portal_api_dev_portal_v1_common_pb.DataSource.AsObject, b
  //   customStyling?: CustomStyling.AsObject,
  //   staticPagesList: Array<StaticPage.AsObject>,
  //   publishApiDocs?: dev_portal_api_dev_portal_v1_common_pb.Selector.AsObject,
  //   keyScopesList: Array < KeyScope.AsObject >,

  const handleCreatePortal = async (values: {
    displayName: string;
    version: string;
    description: string;
    name: string;
    namespace: string;
    banner: string;
    favicon: string;
    primaryLogo: string;
  }) => {
    const {
      banner,
      description,
      displayName,
      name,
      namespace,
      favicon,
      primaryLogo
    } = values;
    let newPortal = new Portal().toObject();
    let newPortalStatus = new PortalStatus().toObject();

    await portalApi.createPortal({
      portal: {
        metadata: {
          name,
          namespace,
          ...newPortal.metadata!
        },
        spec: {
          publishApiDocs: {
            matchLabelsMap: [['', '']],
            ...newPortal.spec?.publishApiDocs
          },
          keyScopesList: {},
          domainsList: [],
          staticPagesList: [],
          favicon: {
            inlineString: favicon!,
            ...newPortal.spec?.favicon!
          },
          primaryLogo: {
            inlineString: primaryLogo!,
            ...newPortal.spec?.primaryLogo!
          },
          displayName,
          description,
          customStyling: {
            backgroundColor: '',
            buttonColorOverride: '',
            defaultTextColor: '',
            navigationLinksColorOverride: '',
            primaryColor: '',
            secondaryColor: ''
          },
          banner: {
            inlineString: banner,
            ...newPortal.spec?.banner!
          },
          ...newPortal.spec!
        },
        status: {
          ...newPortalStatus,
          state: 0,
          reason: '',
          publishUrl: '',
          keyScopesList: [],
          observedGeneration: 0,
          apiDocsList: []
        },
        ...newPortal!
      }
    } as any);
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
            namespace: '',
            name: '',
            banner: '',
            favicon: '',
            primaryLogo: ''
          }}
          onSubmit={handleCreatePortal}>
          {formik => (
            <>
              <pre>{JSON.stringify(formik.values, null, 2)}</pre>
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
                  <StyledTab>Branding</StyledTab>
                  <StyledTab>APIs</StyledTab>
                  <StyledTab>Access</StyledTab>
                </TabList>

                <TabPanels className='bg-white rounded-r-lg'>
                  <TabPanel className='relative focus:outline-none'>
                    <GeneralSection />

                    <div className='flex items-center justify-between px-6 '>
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
                  <TabPanel className='focus:outline-none'>
                    <ImagerySection />
                    <div className='flex items-center justify-between px-6 '>
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
                  <TabPanel className='focus:outline-none'>
                    <BrandingSection />
                    <div className='flex items-center justify-between px-6 '>
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
                  <TabPanel className='focus:outline-none'>
                    <APIsSection />
                    <div className='flex items-center justify-between px-6 '>
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
                  <TabPanel className='focus:outline-none'>
                    <AccessSection />
                    <div className='flex items-center justify-between px-6 '>
                      <button
                        className='text-blue-500 cursor-pointer'
                        onClick={props.onClose}>
                        cancel
                      </button>
                      <SoloButtonStyledComponent onClick={formik.handleSubmit}>
                        Create Portal
                      </SoloButtonStyledComponent>
                    </div>
                  </TabPanel>
                </TabPanels>
              </Tabs>
            </>
          )}
        </Formik>
      </div>
    </>
  );
};
