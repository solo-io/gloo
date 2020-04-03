import { css } from '@emotion/core';
import {
  Tab,
  TabList,
  TabPanel,
  TabPanelProps,
  TabPanels,
  Tabs
} from '@reach/tabs';
import { ReactComponent as NoImageIcon } from 'assets/no-image-placeholder.svg';
import { Loading } from 'Components/Common/DisplayOnly/Loading';
import {
  SoloFormInput,
  SoloFormTextarea,
  SoloFormStringsList
} from 'Components/Common/Form/SoloFormField';
import { SoloTransfer, ListItemType } from 'Components/Common/SoloTransfer';
import { Formik, useFormikContext } from 'formik';
import { ObjectRef } from 'proto/dev-portal/api/dev-portal/v1/common_pb';
import { Portal } from 'proto/dev-portal/api/grpc/admin/portal_pb';
import React from 'react';
import ImageUploader from 'react-images-upload';
import { colors } from 'Styles';
import {
  SoloButtonStyledComponent,
  SoloCancelButton
} from 'Styles/CommonEmotions/button';
import useSWR from 'swr';
import { apiDocApi, groupApi, portalApi, userApi } from '../api';
import {
  SectionContainer,
  SectionHeader,
  SectionSubHeader
} from '../apis/CreateAPIModal';
import { configAPI } from 'store/config/api';

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
      <div className='p-3 text-gray-700 bg-gray-100 rounded-lg'>
        Create a new Portal to expose business APIs
      </div>
      <div className='grid w-full grid-cols-2 mt-2'>
        <div className='mr-4 '>
          <SoloFormInput
            name='displayName'
            title='Name'
            placeholder='Portal name goes here'
          />
        </div>
        <div>
          <SoloFormStringsList
            name='domainsList'
            label='Portal Domain(s)'
            createNewPromptText='Domain'
            hideError
          />
        </div>
        <div className='col-span-2 mt-2'>
          <SoloFormTextarea
            name='description'
            title='Description'
            placeholder='Portal description goes here'
            hideError
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
    formik.setFieldValue('banner', files[0]);

    setImages([...images, ...pictures]);
  };
  return (
    <SectionContainer>
      <SectionHeader>Create a Portal: Add Imagery</SectionHeader>
      <div className='p-3 mb-2 text-gray-700 bg-gray-100 rounded-lg'>
        Add a "Hero" image associated with your Portal
      </div>
      <div className='flex flex-col items-center p-4 pb-0 mr-4 bg-gray-100 border border-gray-400 rounded-lg'>
        <NoImageIcon />
        <ImageUploader
          label='Max 1.5MiB'
          css={css`
            .fileContainer {
              background: ${colors.januaryGrey};
              border-radius: 8px;
              text-align: center;

              box-shadow: none;
              padding: 0px;
            }
            .uploadPictureContainer {
              margin: 0;
            }
          `}
          withPreview
          singleImage
          withIcon={false}
          buttonStyles={{
            borderRadius: '8px',
            backgroundColor: colors.blue600,
            fontSize: '16px',
            fontWeight: 400
          }}
          buttonText='Upload an Image'
          onChange={onDrop}
          imgExtension={['.jpg', '.png', '.jpeg', 'ico']}
          maxFileSize={1572864}
        />
      </div>
    </SectionContainer>
  );
};

const BrandingSection = () => {
  const formik = useFormikContext();
  const [images, setImages] = React.useState<string[]>([]);

  const onDropPrimaryLogo = (files: File[], pictures: string[]) => {
    formik.setFieldValue('primaryLogo', files[0]);
    setImages([...images, ...pictures]);
  };
  const onDropFavicon = (files: File[], pictures: string[]) => {
    formik.setFieldValue('favicon', files[0]);
    setImages([...images, ...pictures]);
  };
  return (
    <SectionContainer>
      <SectionHeader>Create a Portal: Branding Logos</SectionHeader>
      <div className='p-3 mb-2 text-gray-700 bg-gray-100 rounded-lg'>
        Add company logos to give a branded experience
      </div>
      <div className='grid grid-cols-2'>
        <div className='flex flex-col items-center'>
          <SectionSubHeader>Primary Logo</SectionSubHeader>
          <div className='flex flex-col items-center p-4 pb-0 mr-4 bg-gray-100 border border-gray-400 rounded-lg'>
            <NoImageIcon />
            <ImageUploader
              label='Max 1.5MiB'
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
              singleImage
              withIcon={false}
              buttonStyles={{
                borderRadius: '8px',
                backgroundColor: colors.blue600,
                fontSize: '16px',
                fontWeight: 400
              }}
              buttonText='Upload an Image'
              onChange={onDropPrimaryLogo}
              imgExtension={['.jpg', '.png', '.jpeg', 'ico']}
              maxFileSize={1572864}
            />
          </div>
        </div>
        <div className='flex flex-col items-center'>
          <SectionSubHeader>Favicon</SectionSubHeader>

          <div className='flex flex-col items-center p-4 pb-0 mr-4 bg-gray-100 border border-gray-400 rounded-lg'>
            <NoImageIcon />
            <ImageUploader
              label='Max 1.5MiB'
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
              singleImage
              withIcon={false}
              buttonStyles={{
                borderRadius: '8px',
                backgroundColor: colors.blue600,
                fontSize: '16px',
                fontWeight: 400
              }}
              buttonText='Upload an Image'
              onChange={onDropFavicon}
              imgExtension={['.jpg', '.png', '.jpeg', 'ico']}
              maxFileSize={1572864}
            />
          </div>
        </div>
      </div>
    </SectionContainer>
  );
};

type CreatePortalValues = {
  displayName: string;
  description: string;
  domainsList: string[];
  banner: File;
  favicon: File;
  primaryLogo: File;
  chosenAPIs: ListItemType[];
  chosenUsers: ListItemType[];
  chosenGroups: ListItemType[];
};
export const CreatePortalModal: React.FC<{ onClose: () => void }> = props => {
  const { data: apiDocsList, error: apiDocsError } = useSWR(
    'listApiDocs',
    apiDocApi.listApiDocs
  );
  const { data: userList, error: userError } = useSWR(
    'listUsers',
    userApi.listUsers
  );

  const { data: groupList, error: groupError } = useSWR(
    'listGroups',
    groupApi.listGroups
  );
  const { data: podNamespace, error: podNamespaceError } = useSWR(
    'getPodNamespace',
    configAPI.getPodNamespace
  );
  const [tabIndex, setTabIndex] = React.useState(0);

  const handleTabsChange = (index: number) => {
    setTabIndex(index);
  };

  const handleCreatePortal = async (values: CreatePortalValues) => {
    const {
      banner,
      description,
      displayName,
      domainsList,
      favicon,
      primaryLogo
    } = values;
    let newPortal = new Portal().toObject();

    let bannerUint8Arr = new Uint8Array();
    if (banner !== undefined) {
      let bannerBuffer = await banner.arrayBuffer();
      bannerUint8Arr = new Uint8Array(bannerBuffer);
    }

    let faviconUint8Arr = new Uint8Array();
    if (favicon !== undefined) {
      let faviconBuffer = await favicon.arrayBuffer();
      faviconUint8Arr = new Uint8Array(faviconBuffer);
    }

    let primaryLogoUint8Arr = new Uint8Array();
    if (primaryLogo !== undefined) {
      let primaryLogoBuffer = await primaryLogo.arrayBuffer();
      primaryLogoUint8Arr = new Uint8Array(primaryLogoBuffer);
    }

    //@ts-ignore
    await portalApi.createPortal({
      portal: {
        ...newPortal!,
        //@ts-ignore
        metadata: {
          // ...newPortal.metadata!,
          namespace: podNamespace!,
          annotationsMap: [],
          labelsMap: [],
          resourceVersion: '',
          uid: '',
          creationTimestamp: { nanos: 0, seconds: 0 }
        },
        spec: {
          // ...newPortal.spec!,
          domainsList,
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
            inlineBytes: bannerUint8Arr
          },
          favicon: {
            ...newPortal.spec?.favicon!,
            inlineBytes: faviconUint8Arr
          },
          primaryLogo: {
            ...newPortal.spec?.primaryLogo!,
            inlineBytes: primaryLogoUint8Arr
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
      },
      apiDocsList: values.chosenAPIs,
      groupsList: values.chosenGroups,
      usersList: values.chosenUsers
    });
    props.onClose();
  };

  if (!apiDocsList || !userList || !groupList) {
    return <Loading center>Loading</Loading>;
  }
  return (
    <>
      <div
        css={css`
          width: 750px;
        `}
        className='bg-white rounded-lg shadow '>
        <Formik<CreatePortalValues>
          initialValues={{
            displayName: '',
            description: '',
            domainsList: ([] as unknown) as string[],
            banner: (undefined as unknown) as File,
            favicon: (undefined as unknown) as File,
            primaryLogo: (undefined as unknown) as File,
            chosenAPIs: [] as ListItemType[],
            chosenUsers: [] as ListItemType[],
            chosenGroups: [] as ListItemType[]
          }}
          onSubmit={handleCreatePortal}>
          {formik => (
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
                  <StyledTab>Branding</StyledTab>
                  <StyledTab>APIs</StyledTab>
                  <StyledTab>User Access</StyledTab>
                  <StyledTab>Group Access</StyledTab>
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
                      <div>
                        <SoloButtonStyledComponent
                          onClick={() => setTabIndex(tabIndex + 1)}>
                          Next Step
                        </SoloButtonStyledComponent>
                      </div>
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
                      <div>
                        <SoloCancelButton
                          onClick={() => handleTabsChange(0)}
                          className='mr-2'>
                          Back
                        </SoloCancelButton>
                        <SoloButtonStyledComponent
                          onClick={() => setTabIndex(tabIndex + 1)}>
                          Next Step
                        </SoloButtonStyledComponent>
                      </div>
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
                      <div>
                        <SoloCancelButton
                          onClick={() => handleTabsChange(1)}
                          className='mr-2'>
                          Back
                        </SoloCancelButton>
                        <SoloButtonStyledComponent
                          onClick={() => setTabIndex(tabIndex + 1)}>
                          Next Step
                        </SoloButtonStyledComponent>
                      </div>
                    </div>
                  </TabPanel>
                  <TabPanel className='relative flex flex-col justify-between h-full focus:outline-none'>
                    <SectionContainer>
                      <SectionHeader>Create a Portal: APIs</SectionHeader>
                      <div className='p-3 mb-2 text-gray-700 bg-gray-100 rounded-lg'>
                        Select the APIs you'd like to make available through
                        this portal
                      </div>

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
                              name: apiDoc.metadata?.name!,
                              namespace: apiDoc.metadata?.namespace!,
                              displayValue: apiDoc.status?.displayName
                            };
                          })}
                        chosenOptionsListName='Selected APIs'
                        chosenOptions={formik.values.chosenAPIs}
                        onChange={newChosenOptions => {
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
                      <div>
                        <SoloCancelButton
                          onClick={() => handleTabsChange(2)}
                          className='mr-2'>
                          Back
                        </SoloCancelButton>
                        <SoloButtonStyledComponent
                          onClick={() => setTabIndex(tabIndex + 1)}>
                          Next Step
                        </SoloButtonStyledComponent>
                      </div>
                    </div>
                  </TabPanel>
                  <TabPanel className='relative flex flex-col justify-between h-full focus:outline-none'>
                    <SectionContainer>
                      <SectionHeader>Create an API: User Access</SectionHeader>
                      <div className='p-3 mb-2 text-gray-700 bg-gray-100 rounded-lg'>
                        Select the users and groups to which you'd like to grant
                        access to this portal
                      </div>

                      <SoloTransfer
                        allOptionsListName='Available Users'
                        allOptions={userList
                          .sort((a, b) =>
                            a.metadata?.name === b.metadata?.name
                              ? 0
                              : a.metadata!.name > b.metadata!.name
                              ? 1
                              : -1
                          )
                          .map(user => {
                            return {
                              name: user.metadata?.name!,
                              namespace: user.metadata?.namespace!,
                              displayValue: user.spec?.username
                            };
                          })}
                        chosenOptionsListName='Selected Users'
                        chosenOptions={formik.values.chosenUsers}
                        onChange={newChosenOptions => {
                          formik.setFieldValue('chosenUsers', newChosenOptions);
                        }}
                      />
                    </SectionContainer>
                    <div className='flex items-end justify-between h-full px-6 mb-4 '>
                      <button
                        className='text-blue-500 cursor-pointer'
                        onClick={props.onClose}>
                        cancel
                      </button>
                      <div>
                        <SoloCancelButton
                          onClick={() => handleTabsChange(3)}
                          className='mr-2'>
                          Back
                        </SoloCancelButton>
                        <SoloButtonStyledComponent
                          onClick={() => setTabIndex(tabIndex + 1)}>
                          Next Step
                        </SoloButtonStyledComponent>
                      </div>
                    </div>
                  </TabPanel>
                  <TabPanel className='relative flex flex-col justify-between h-full focus:outline-none'>
                    <SectionContainer>
                      <SectionHeader>Create an API: Group Access</SectionHeader>
                      <div className='p-3 mb-2 text-gray-700 bg-gray-100 rounded-lg'>
                        Select the users and groups to which you'd like to grant
                        access to this portal
                      </div>

                      <SoloTransfer
                        allOptionsListName='Available Groups'
                        allOptions={groupList
                          .sort((a, b) =>
                            a.metadata?.name === b.metadata?.name
                              ? 0
                              : a.metadata!.name > b.metadata!.name
                              ? 1
                              : -1
                          )
                          .map(group => {
                            return {
                              name: group.metadata?.name!,
                              namespace: group.metadata?.namespace!,
                              displayValue: group?.spec?.displayName
                            };
                          })}
                        chosenOptionsListName='Selected Groups'
                        chosenOptions={formik.values.chosenGroups}
                        onChange={newChosenOptions => {
                          formik.setFieldValue(
                            'chosenGroups',
                            newChosenOptions
                          );
                        }}
                      />
                    </SectionContainer>
                    <div className='flex items-end justify-between h-full px-6 mb-4 '>
                      <button
                        className='text-blue-500 cursor-pointer'
                        onClick={props.onClose}>
                        cancel
                      </button>
                      <div>
                        <SoloCancelButton
                          onClick={() => handleTabsChange(4)}
                          className='mr-2'>
                          Back
                        </SoloCancelButton>
                        <SoloButtonStyledComponent
                          onClick={formik.handleSubmit}>
                          Create Portal
                        </SoloButtonStyledComponent>
                      </div>
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
