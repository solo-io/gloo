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
import { portalApi, portalMessageFromObject, apiDocApi } from '../api';
import { Portal } from 'proto/dev-portal/api/grpc/admin/portal_pb';
import {
  PortalSpec,
  PortalStatus
} from 'proto/dev-portal/api/dev-portal/v1/portal_pb';
import {
  StateMap,
  State,
  ObjectRef,
  DataSource
} from 'proto/dev-portal/api/dev-portal/v1/common_pb';
import { SoloTransfer } from 'Components/Common/SoloTransfer';
import useSWR from 'swr';
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
            name='displayName'
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

  const onDrop = async (files: File[], pictures: string[]) => {
    console.log('files', files);
    console.log('pictures', pictures);
    let buff = await files[0].arrayBuffer();
    let newDataSource = new DataSource();
    let uint8Arr = new Uint8Array(buff);
    newDataSource.setInlineBytes(uint8Arr);
    formik.setFieldValue('banner', files[0]);
    // formik.setFieldValue('banner', files[0].arrayBuffer);

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
    formik.setFieldValue(
      'primaryLogo',
      pictures[0].split(';')[2].split(',')[1]
    );
    setImages([...images, ...pictures]);
  };
  const onDropFavicon = (files: File[], pictures: string[]) => {
    formik.setFieldValue('favicon', pictures[0]);
    setImages([...images, ...pictures]);
  };
  return (
    <SectionContainer>
      <SectionHeader>Create a Portal: Branding Logos</SectionHeader>
      <div className='grid grid-cols-2'>
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

const AccessSection = () => {
  return (
    <SectionContainer>
      <SectionHeader>Create a Portal: Access</SectionHeader>
    </SectionContainer>
  );
};

export const CreatePortalModal: React.FC<{ onClose: () => void }> = props => {
  const { data: apiDocsList, error: apiDocsError } = useSWR(
    'listApiDocs',
    apiDocApi.listApiDocs
  );
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
    description: string;
    name: string;
    banner: File;
    favicon: string;
    primaryLogo: string;
  }) => {
    const {
      banner,
      description,
      displayName,
      name,
      favicon,
      primaryLogo
    } = values;
    let newPortal = new Portal().toObject();
    let newPortalStatus = new PortalStatus().toObject();
    let buff = await values.banner.arrayBuffer();
    let uint8Arr = new Uint8Array(buff);

    //@ts-ignore
    await portalApi.createPortal({
      portal: {
        ...newPortal!,
        metadata: {
          // ...newPortal.metadata!,
          name: displayName,
          namespace: 'gloo-system',
          annotationsMap: [],
          labelsMap: [],
          resourceVersion: '',
          uid: '',
          creationTimestamp: { nanos: 0, seconds: 0 }
        },
        spec: {
          // ...newPortal.spec!,
          domainsList: ['localhost:3001'],
          keyScopesList: [],
          staticPagesList: [],
          customStyling: {
            backgroundColor: '',
            buttonColorOverride: '',
            navigationLinksColorOverride: '',
            primaryColor: '',
            defaultTextColor: '',
            secondaryColor: ''
          },
          publishApiDocs: { matchLabelsMap: [] },
          description,
          displayName,
          banner: {
            ...newPortal.spec?.banner!,
            inlineBytes: uint8Arr
          },
          favicon: {
            ...newPortal.spec?.favicon!,
            inlineString: favicon
          },
          primaryLogo: {
            ...newPortal.spec?.primaryLogo!,
            inlineString: primaryLogo
          }
        },
        status: {
          // ...newPortal.status!,
          apiDocsList: [],
          keyScopesList: [],
          observedGeneration: (undefined as unknown) as number,
          publishUrl: (undefined as unknown) as string,
          reason: (undefined as unknown) as string,
          state: (undefined as unknown) as any
        }
      }
    });
    props.onClose();
  };

  if (!apiDocsList) {
    return <Loading center>Loading</Loading>;
  }
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
            description: '',
            name: '',
            banner: (undefined as unknown) as File,
            favicon: '',
            primaryLogo: '',
            chosenAPIs: [] as ObjectRef.AsObject[]
          }}
          onSubmit={handleCreatePortal}>
          {formik => (
            <>
              {/* <pre>{JSON.stringify(formik.values, null, 2)}</pre> */}
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
                  <TabPanel className='flex flex-col justify-between h-full focus:outline-none'>
                    <BrandingSection />

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
                    <SectionContainer>
                      <SectionHeader>Create a Portal: APIs</SectionHeader>
                      <SoloTransfer
                        allOptionsListName='Available APIs'
                        allOptions={apiDocsList
                          .sort((a, b) =>
                            a.metadata?.name === b.metadata?.name
                              ? 0
                              : a.metadata!.name > b.metadata!.name
                              ? 1
                              : -1
                          )
                          .map(apiDoc => {
                            return {
                              value: apiDoc.metadata?.name!,
                              displayValue: apiDoc.metadata?.name!
                            };
                          })}
                        chosenOptionsListName='Selected APIs'
                        chosenOptions={formik.values.chosenAPIs.map(api => {
                          return { value: api.name + api.namespace };
                        })}
                        onChange={newChosenOptions => {
                          console.log('newChosenOptions', newChosenOptions);
                          formik.setFieldValue('chosenAPIs', newChosenOptions);
                        }}
                      />
                    </SectionContainer>
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
                    <AccessSection />
                    <div className='flex items-end justify-between h-full px-6 mb-4 '>
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
