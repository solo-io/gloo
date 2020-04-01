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
import { CreateAPIModal, SectionSubHeader } from '../apis/CreateAPIModal';
import { ReactComponent as GreenPlus } from 'assets/small-green-plus.svg';
import {
  SoloNegativeButton,
  SoloCancelButton,
  SoloButtonStyledComponent
} from 'Styles/CommonEmotions/button';
import { ConfirmationModal } from 'Components/Common/ConfirmationModal';
import { PortalApiDocsTab } from './PortalApiDocsTab';
import { PortalGroupsTab } from './PortalGroupsTab';
import {
  SoloFormDropdown,
  SoloFormInput,
  SoloFormTextarea
} from 'Components/Common/Form/SoloFormField';
import { ChromePicker } from 'react-color';
import { ColorPicker } from './ColorPicker';
import { Formik } from 'formik';
import { State } from 'proto/dev-portal/api/dev-portal/v1/common_pb';

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
type UpdatePortalValues = {
  displayName: string;
  description: string;
  domainsList: string[];
  primaryColor: string;
  secondaryColor: string;
  backgroundColor: string;
  defaultTextColor: string;
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
  const [editMode, setEditMode] = React.useState(false);

  const toggleEditMode = () => {
    setEditMode(!editMode);
  };
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

  const handleUpdatePortal = async (values: UpdatePortalValues) => {
    const {
      backgroundColor,
      primaryColor,
      secondaryColor,
      displayName,
      domainsList,
      description,
      defaultTextColor
    } = values;

    //@ts-ignore
    await portalApi.updatePortal({
      portal: {
        ...portal,
        //@ts-ignore
        metadata: {
          ...portal.metadata,
          name: portal.metadata?.name!,
          namespace: portal.metadata?.namespace!
        },
        spec: {
          ...portal.spec,
          domainsList: domainsList,
          // @ts-ignore
          customStyling: {
            backgroundColor,
            primaryColor,
            defaultTextColor,
            secondaryColor
          },
          description,
          displayName
        }
      },
      portalOnly: true
    });
  };

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
            }
          ]}
          healthMessage={'Portal Status'}
          onClose={() => history.push(`/dev-portal/`)}>
          <Formik
            onSubmit={handleUpdatePortal}
            initialValues={{
              displayName: '',
              description: '',
              domainsList: [],
              primaryColor:
                portal.spec?.customStyling?.primaryColor || '#2196C9',
              secondaryColor:
                portal.spec?.customStyling?.secondaryColor || '#253E58',
              backgroundColor:
                portal.spec?.customStyling?.backgroundColor || '#F9F9F9',

              defaultTextColor:
                portal.spec?.customStyling?.defaultTextColor || '#35393B'
            }}>
            {formik => (
              <div>
                <div className='relative flex items-center'>
                  <div className='w-64 max-h-72'>
                    {portal.spec?.banner?.inlineBytes ? (
                      <img
                        className='object-cover max-h-72'
                        src={`data:image/gif;base64,${portal.spec?.banner?.inlineBytes}`}></img>
                    ) : (
                      <PlaceholderPortal className='w-56 rounded-lg ' />
                    )}
                  </div>
                  <div className='grid w-full grid-cols-2 ml-2 h-36'>
                    <div>
                      <span className='font-medium text-gray-900'>
                        Portal Display Name
                      </span>
                      {editMode ? (
                        <SoloFormInput
                          name='displayName'
                          placeholder={portal?.spec?.displayName}
                          hideError
                        />
                      ) : (
                        <div>{portal?.spec?.displayName}</div>
                      )}
                    </div>
                    <div>
                      <span className='font-medium text-gray-900'>
                        Portal Domains
                      </span>
                      {portal.spec?.domainsList.map((domain, index) => (
                        <div
                          key={domain}
                          className='flex items-center mb-2 text-sm text-blue-600'>
                          <span>
                            <ExternalLinkIcon className='w-4 h-4 ' />
                          </span>
                          {editMode ? (
                            <SoloFormInput
                              hideError
                              name={`domainsList.${index}`}
                              placeholder={domain}
                            />
                          ) : (
                            <div>{domain}</div>
                          )}
                        </div>
                      ))}
                    </div>
                    <span
                      onClick={toggleEditMode}
                      className='absolute top-0 right-0 flex items-center'>
                      <span className='mr-2'> Edit</span>
                      <span className='flex items-center justify-center w-6 h-6 mr-3 text-gray-700 bg-gray-400 rounded-full cursor-pointer'>
                        <EditIcon className='w-3 h-3' />
                      </span>
                    </span>
                    <div className='col-span-2 '>
                      <span className='font-medium text-gray-900'>
                        Description
                      </span>
                      {editMode ? (
                        <SoloFormTextarea
                          rows={3}
                          name='description'
                          hideError
                          placeholder={portal.spec?.description}
                        />
                      ) : (
                        <div className='break-words '>
                          {portal.spec?.description}
                        </div>
                      )}
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
                    <TabPanel className=' focus:outline-none'>
                      <div className='relative flex p-4 border border-gray-300 rounded-lg'>
                        {/* background */}
                        <div>
                          <svg
                            className='border border-gray-600 rounded-md'
                            xmlns='http://www.w3.org/2000/svg'
                            width='581'
                            height='396'>
                            <g fill={formik.values.backgroundColor}>
                              <rect
                                width='581'
                                height='396'
                                stroke='none'
                                rx='4'
                              />
                              <rect
                                width='580'
                                height='395'
                                x='.5'
                                y='.5'
                                fill='none'
                                rx='3.5'
                              />
                              <svg
                                xmlns='http://www.w3.org/2000/svg'
                                viewBox='0 150 484.222 15'>
                                <g transform='translate(-38.384 -10)'>
                                  <g
                                    fill='#6e7477'
                                    transform='translate(38.384 11)'>
                                    <circle cx='7' cy='7' r='7' />
                                    <path d='M23.313 9.956a3.491 3.491 0 003.07-1.73l-1.148-.572a2.284 2.284 0 01-1.922 1.111 2.624 2.624 0 01-2.637-2.761 2.611 2.611 0 012.637-2.76 2.258 2.258 0 011.922 1.111l1.134-.6a3.426 3.426 0 00-3.056-1.703 3.848 3.848 0 00-4.012 3.952 3.855 3.855 0 004.012 3.952zm6.53 0a2.772 2.772 0 002.838-2.912 2.757 2.757 0 00-2.841-2.9 2.757 2.757 0 00-2.841 2.9 2.772 2.772 0 002.841 2.912zm0-1.065a1.663 1.663 0 01-1.595-1.847 1.608 1.608 0 113.185 0 1.656 1.656 0 01-1.593 1.847zm11.982.928v-4a1.463 1.463 0 00-1.595-1.675 2.333 2.333 0 00-1.879 1.008 1.445 1.445 0 00-1.512-1.008 2.451 2.451 0 00-1.81.882v-.74h-1.2v5.533h1.2V5.944a1.7 1.7 0 011.284-.73c.653 0 .9.4.9 1v3.605h1.2V5.944a1.647 1.647 0 011.29-.73c.653 0 .916.4.916 1v3.605zm4.456.137c1.432 0 2.451-1.088 2.451-2.91s-1.019-2.902-2.453-2.902a2.208 2.208 0 00-1.776.9v-.758h-1.2v7.641h1.2V9.051a2.22 2.22 0 001.776.905zm-.358-1.065a1.806 1.806 0 01-1.42-.756V5.944a1.8 1.8 0 011.42-.733 1.627 1.627 0 011.57 1.833 1.638 1.638 0 01-1.57 1.847zm8.523.928V6.144c0-1.478-1.077-1.993-2.314-1.993a3.241 3.241 0 00-2.325.894l.5.836a2.215 2.215 0 011.638-.71c.756 0 1.294.389 1.294 1.031v.822a2.313 2.313 0 00-1.8-.71 1.775 1.775 0 00-1.913 1.81 1.842 1.842 0 001.913 1.833 2.329 2.329 0 001.8-.745v.607zm-2.52-.687a1.06 1.06 0 01-1.191-1 1.067 1.067 0 011.191-1 1.651 1.651 0 011.317.573v.848a1.651 1.651 0 01-1.317.578zm8.947.687V5.924a1.6 1.6 0 00-1.787-1.78 2.638 2.638 0 00-1.959.882v-.74h-1.2v5.533h1.2V5.944a1.885 1.885 0 011.42-.733c.676 0 1.123.275 1.123 1.146v3.46zm1.249 2.165a3.3 3.3 0 00.676.08 1.839 1.839 0 001.89-1.249l2.669-6.53h-1.284l-1.615 4.148-1.615-4.147h-1.294l2.28 5.6-.275.63a.755.755 0 01-.8.481 1.257 1.257 0 01-.458-.092zm13.437-2.165v-1.18h-3.368V2.178h-1.34v7.641zm3.4.137a2.772 2.772 0 002.843-2.912 2.842 2.842 0 10-5.682 0 2.772 2.772 0 002.841 2.912zm0-1.065a1.663 1.663 0 01-1.59-1.847 1.608 1.608 0 113.185 0 1.656 1.656 0 01-1.593 1.847zm6.243 3.173c1.352 0 2.841-.538 2.841-2.532V4.286h-1.2v.768a2.172 2.172 0 00-1.776-.9c-1.432-.01-2.454 1.037-2.454 2.835 0 1.821 1.042 2.841 2.451 2.841a2.225 2.225 0 001.776-.928v.63a1.466 1.466 0 01-1.638 1.546 2.183 2.183 0 01-1.764-.71l-.561.871a3.166 3.166 0 002.327.825zm.218-3.311a1.568 1.568 0 01-1.569-1.764 1.577 1.577 0 011.571-1.775 1.836 1.836 0 011.42.733V8.02a1.836 1.836 0 01-1.42.733zm6.61 1.2a2.772 2.772 0 002.842-2.909 2.842 2.842 0 10-5.682 0 2.772 2.772 0 002.841 2.912zm0-1.065a1.663 1.663 0 01-1.591-1.844 1.608 1.608 0 113.185 0 1.656 1.656 0 01-1.593 1.847z' />
                                  </g>
                                  <path
                                    fill={formik.values.secondaryColor}
                                    d='M437.018 14h27.929v7h-27.929zm-38 0h27.929v7h-27.929zm-38 0h27.929v7h-27.929z'
                                  />
                                  <rect
                                    width='43.222'
                                    height='14'
                                    fill={formik.values.primaryColor}
                                    rx='4'
                                    transform='translate(479.385 10)'
                                  />
                                </g>
                              </svg>
                              <svg
                                xmlns='http://www.w3.org/2000/svg'
                                viewBox='5 20 300 50'
                                height='250'>
                                <g transform='translate(-674.136 2)'>
                                  <path
                                    fill='#c0cbd3'
                                    d='M674.137-2h312v208h-312z'
                                  />
                                  <circle
                                    fill='#f1f3f5'
                                    cx='12.94'
                                    cy='12.94'
                                    r='12.94'
                                    transform='translate(751.736 74.304)'
                                  />
                                  <path
                                    fill='#f1f3f5'
                                    d='M674.136 128.54s15.24-20.013 32.664-19.778c15.444.208 25.949 6.947 42.773 21.189s30.517 24.98 44.789 6.207S838 71.996 851.205 62.065c13.062-9.823 33.768-19.475 61.827-3.54 33.5 19.025 73.393 78.34 73.393 78.34v69.327H674.136z'
                                  />
                                </g>
                              </svg>
                              <svg
                                xmlns='http://www.w3.org/2000/svg'
                                viewBox='-40 80 581 86.498'>
                                <path
                                  d='M0 37.655h483.344v8.931H0zM0 56.919h483.344v8.931H0zM0 77.567h127.264v8.931H0zM4.7 16.8v-6.25h4.425a5.06 5.06 0 005.425-5.2A5.052 5.052 0 009.125.125h-7.35V16.8zm4.025-8.825H4.7V2.7h4.025a2.644 2.644 0 110 5.275zM21.85 17.1a6.05 6.05 0 006.2-6.35 6.017 6.017 0 00-6.2-6.325 6.017 6.017 0 00-6.2 6.325 6.05 6.05 0 006.2 6.35zm0-2.325c-2.225 0-3.475-1.875-3.475-4.025 0-2.125 1.25-4 3.475-4 2.25 0 3.475 1.875 3.475 4 0 2.15-1.225 4.025-3.475 4.025zM33.175 16.8V8.575A4.122 4.122 0 0136.3 6.95a3.623 3.623 0 01.8.075v-2.6a5.24 5.24 0 00-3.925 2.05v-1.75H30.55V16.8zm9.525.3a3.556 3.556 0 002.45-.775l-.625-2a1.633 1.633 0 01-1.15.45c-.75 0-1.15-.625-1.15-1.45V7h2.45V4.725h-2.45v-3.3H39.6v3.3h-2V7h2v6.975a2.782 2.782 0 003.1 3.125zm14.175-.3V8.775c0-3.225-2.35-4.35-5.05-4.35a7.074 7.074 0 00-5.075 1.95l1.1 1.825a4.835 4.835 0 013.575-1.55c1.65 0 2.825.85 2.825 2.25v1.8a5.047 5.047 0 00-3.925-1.55 3.874 3.874 0 00-4.175 3.95 4.02 4.02 0 004.175 4 5.083 5.083 0 003.925-1.625V16.8zm-5.5-1.5a2.314 2.314 0 01-2.6-2.175 2.329 2.329 0 012.6-2.175 3.6 3.6 0 012.875 1.25v1.85a3.6 3.6 0 01-2.875 1.25zm11.35 1.5V.125H60.1V16.8zm16.475 0V2.7h5.05V.125H71.225V2.7h5.05v14.1zm8.6-13.55a1.625 1.625 0 10-1.625-1.625A1.622 1.622 0 0087.8 3.25zm1.3 13.55V4.725h-2.625V16.8zm6.875.3a3.556 3.556 0 002.45-.775l-.625-2a1.633 1.633 0 01-1.15.45c-.75 0-1.15-.625-1.15-1.45V7h2.45V4.725H95.5v-3.3h-2.625v3.3h-2V7h2v6.975a2.782 2.782 0 003.1 3.125zm6.775-.3V.125h-2.625V16.8zm8.825.3a7.143 7.143 0 004.9-1.775l-1.2-1.725a5.077 5.077 0 01-3.45 1.35A3.651 3.651 0 01108 11.625h9.3v-.65c0-3.8-2.3-6.55-5.925-6.55a6.1 6.1 0 00-6.125 6.325 6.057 6.057 0 006.325 6.35zm3.2-7.35h-6.8a3.3 3.3 0 013.375-3.175 3.25 3.25 0 013.425 3.175z'
                                  fill='#fff'
                                />
                              </svg>

                              <svg
                                xmlns='http://www.w3.org/2000/svg'
                                width='600'
                                viewBox='-30 -50 581 132.492'>
                                <g transform='translate(-38.384 -203)'>
                                  <g
                                    fill='#fff'
                                    stroke='#d4d8de'
                                    transform='translate(38.384 203)'>
                                    <rect
                                      width='230.977'
                                      height='132.492'
                                      stroke='none'
                                      rx='7'
                                    />
                                    <rect
                                      width='229.977'
                                      height='131.492'
                                      x='.5'
                                      y='.5'
                                      fill='none'
                                      rx='6.5'
                                    />
                                  </g>
                                  <path
                                    fill={formik.values.defaultTextColor}
                                    d='M56.814 228.285a5.045 5.045 0 003.836-1.708v-3.486h-4.424v1.442h2.8v1.442a3.351 3.351 0 01-2.212.84 3.207 3.207 0 01-3.22-3.374 3.191 3.191 0 013.22-3.374 3.014 3.014 0 012.422 1.246l1.33-.77a4.369 4.369 0 00-3.752-1.932 4.694 4.694 0 00-4.9 4.83 4.714 4.714 0 004.9 4.844zm8.54-.014a4 4 0 002.744-.994l-.672-.966a2.843 2.843 0 01-1.932.756 2.044 2.044 0 01-2.142-1.864h5.208v-.364c0-2.128-1.288-3.668-3.318-3.668a3.418 3.418 0 00-3.43 3.542 3.392 3.392 0 003.542 3.558zm1.792-4.116h-3.808a1.85 1.85 0 011.888-1.778 1.82 1.82 0 011.92 1.778zm4.816 4.116a1.991 1.991 0 001.372-.434l-.35-1.12a.914.914 0 01-.644.252c-.42 0-.644-.35-.644-.812v-3.542h1.372v-1.274h-1.372v-1.848h-1.47v1.848h-1.12v1.274h1.12v3.906a1.558 1.558 0 001.736 1.75zm4.364 0a1.991 1.991 0 001.372-.434l-.35-1.12a.914.914 0 01-.644.252c-.42 0-.644-.35-.644-.812v-3.542h1.376v-1.274h-1.372v-1.848h-1.47v1.848h-1.12v1.274h1.12v3.906a1.558 1.558 0 001.732 1.75zm3.07-7.756a.908.908 0 00.91-.91.908.908 0 00-.91-.91.908.908 0 00-.91.91.908.908 0 00.91.91zm.73 7.588v-6.762h-1.47v6.762zm7.854 0v-4.76a1.949 1.949 0 00-2.184-2.17 3.224 3.224 0 00-2.394 1.078v-.91h-1.476v6.762h1.47v-4.732a2.3 2.3 0 011.736-.9c.826 0 1.372.336 1.372 1.4v4.232zm4.564 2.744c1.652 0 3.472-.658 3.472-3.094v-6.412h-1.47v.938a2.655 2.655 0 00-2.17-1.106c-1.75 0-3 1.274-3 3.472 0 2.226 1.274 3.472 3 3.472a2.719 2.719 0 002.17-1.134v.77a1.791 1.791 0 01-2 1.89 2.668 2.668 0 01-2.16-.868l-.686 1.064a3.869 3.869 0 002.842 1.008zm.264-4.044a1.916 1.916 0 01-1.918-2.156 1.928 1.928 0 011.918-2.17 2.244 2.244 0 011.736.9v2.526a2.244 2.244 0 01-1.736.9zm11.774 1.468c2.506 0 3.612-1.344 3.612-2.9 0-3.472-5.432-2.394-5.432-4.144 0-.686.616-1.162 1.568-1.162a3.806 3.806 0 012.716 1.08l.924-1.218a4.833 4.833 0 00-3.486-1.316c-2.058 0-3.4 1.19-3.4 2.744 0 3.43 5.432 2.212 5.432 4.172 0 .63-.518 1.288-1.862 1.288a4.012 4.012 0 01-2.954-1.3l-.924 1.274a5.1 5.1 0 003.806 1.482zm6.8 0a1.991 1.991 0 001.372-.434l-.35-1.12a.914.914 0 01-.644.252c-.42 0-.644-.35-.644-.812v-3.542h1.372v-1.274h-1.368v-1.848h-1.47v1.848h-1.124v1.274h1.12v3.906a1.558 1.558 0 001.74 1.75zm7.944-.168v-4.494c0-1.806-1.318-2.436-2.83-2.436a3.961 3.961 0 00-2.842 1.092l.616 1.022a2.707 2.707 0 012-.868c.924 0 1.582.476 1.582 1.26v1.008a2.826 2.826 0 00-2.2-.868 2.169 2.169 0 00-2.338 2.212 2.251 2.251 0 002.342 2.24 2.847 2.847 0 002.2-.91v.742zm-3.08-.84a1.3 1.3 0 01-1.456-1.218 1.3 1.3 0 011.456-1.218 2.018 2.018 0 011.61.7v1.036a2.018 2.018 0 01-1.612.7zm6.354.84v-4.606a2.308 2.308 0 011.75-.91 2.029 2.029 0 01.448.042v-1.456a2.934 2.934 0 00-2.2 1.148v-.98h-1.472v6.762zm5.334.168a1.991 1.991 0 001.372-.434l-.35-1.12a.914.914 0 01-.644.252c-.42 0-.644-.35-.644-.812v-3.542h1.372v-1.274h-1.372v-1.848h-1.47v1.848h-1.12v1.274h1.12v3.906a1.558 1.558 0 001.736 1.75zm5.46 0a4 4 0 002.744-.994l-.672-.966a2.843 2.843 0 01-1.932.756 2.044 2.044 0 01-2.142-1.864h5.208v-.364c0-2.128-1.288-3.668-3.318-3.668a3.418 3.418 0 00-3.43 3.542 3.392 3.392 0 003.542 3.558zm1.792-4.116h-3.808a1.85 1.85 0 011.89-1.778 1.82 1.82 0 011.918 1.778zm9 3.948v-9.338h-1.47v3.514a2.636 2.636 0 00-2.17-1.106c-1.736 0-3 1.358-3 3.542 0 2.24 1.274 3.556 3 3.556a2.718 2.718 0 002.17-1.092v.924zm-3.204-1.134a1.994 1.994 0 01-1.918-2.254 1.989 1.989 0 011.918-2.24 2.194 2.194 0 011.736.91v2.674a2.194 2.194 0 01-1.736.91z'
                                  />
                                  <path
                                    fill='#d4d8de'
                                    d='M51.017 243.049h205.479v8H51.017zm0 17.256h205.479v8H51.017zm0 18.496h114v8h-114z'
                                  />
                                  <g
                                    fill='#fff'
                                    stroke='#d4d8de'
                                    transform='translate(291.384 203)'>
                                    <rect
                                      width='231.426'
                                      height='132.492'
                                      stroke='none'
                                      rx='7'
                                    />
                                    <rect
                                      width='230.426'
                                      height='131.492'
                                      x='.5'
                                      y='.5'
                                      fill='none'
                                      rx='6.5'
                                    />
                                  </g>
                                  <path
                                    fill={formik.values.defaultTextColor}
                                    d='M313.734 228.103l-3.668-9.338h-2.04l-3.672 9.338h1.862l.686-1.806h4.284l.686 1.806zm-3.008-3.248h-3.362l1.68-4.452zm5.556 3.248v-3.5h2.478a2.834 2.834 0 003.038-2.912 2.829 2.829 0 00-3.038-2.926h-4.116v9.338zm2.254-4.942h-2.254v-2.958h2.254a1.461 1.461 0 011.582 1.484 1.458 1.458 0 01-1.582 1.474zm6.174 4.942v-9.338h-1.638v9.338zm4.172.168c1.82 0 2.814-.91 2.814-2.114 0-2.688-4.088-1.792-4.088-2.982 0-.476.476-.84 1.246-.84a2.758 2.758 0 012 .812l.616-1.036a3.863 3.863 0 00-2.618-.938c-1.708 0-2.66.938-2.66 2.044 0 2.6 4.088 1.652 4.088 2.982 0 .532-.462.9-1.344.9a3.473 3.473 0 01-2.282-.938l-.672 1.05a4.1 4.1 0 002.9 1.06z'
                                  />
                                  <path
                                    fill='#d4d8de'
                                    d='M304.017 243.049h205.929v8H304.017zm0 17.256h205.929v8H304.017zm0 18.496h114v8h-114z'
                                  />
                                  <rect
                                    width='52.222'
                                    height='20'
                                    fill={formik.values.primaryColor}
                                    rx='4'
                                    transform='translate(49.375 297)'
                                  />
                                  <rect
                                    width='52.222'
                                    height='20'
                                    fill={formik.values.primaryColor}
                                    rx='4'
                                    transform='translate(303.824 297)'
                                  />
                                </g>
                              </svg>
                            </g>
                          </svg>
                        </div>
                        <div className='grid w-full grid-cols-2 gap-2 ml-4'>
                          <div className='flex flex-col items-center justify-start '>
                            <SectionSubHeader>Primary Logo</SectionSubHeader>
                            <img
                              className='object-cover h-12'
                              src={`data:image/gif;base64,${portal.spec?.primaryLogo?.inlineBytes}`}></img>
                          </div>
                          <div className='flex flex-col items-center justify-start '>
                            <SectionSubHeader>Favicon</SectionSubHeader>
                            <img
                              className='object-cover h-12'
                              src={`data:image/gif;base64,${portal.spec?.favicon?.inlineBytes}`}></img>
                          </div>

                          <div className='flex flex-col items-start justify-start w-full'>
                            <SectionSubHeader>Primary Color </SectionSubHeader>
                            <ColorPicker
                              name='primaryColor'
                              initialColor={formik.values.primaryColor}
                            />
                          </div>

                          <div className='flex flex-col items-start justify-start w-full'>
                            <SectionSubHeader>Secondary Color</SectionSubHeader>
                            <ColorPicker
                              name='secondaryColor'
                              initialColor={formik.values.secondaryColor}
                            />
                          </div>

                          <div className='flex flex-col items-start justify-start w-full'>
                            <SectionSubHeader>
                              Background Color
                            </SectionSubHeader>
                            <ColorPicker
                              name='backgroundColor'
                              initialColor={formik.values.backgroundColor}
                            />
                          </div>
                          <div className='flex flex-col items-start justify-start w-full'>
                            <SectionSubHeader>
                              Default Text Color
                            </SectionSubHeader>
                            <ColorPicker
                              name='defaultTextColor'
                              initialColor={formik.values.defaultTextColor}
                            />
                          </div>
                        </div>
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
                      <PortalGroupsTab portal={portal} />
                    </TabPanel>
                  </TabPanels>
                </Tabs>
                <div className='flex justify-end justify-between items-bottom '>
                  <div>
                    <SoloButtonStyledComponent
                      className='mr-2'
                      disabled={!formik.dirty}
                      onClick={formik.handleSubmit}>
                      Publish Changes
                    </SoloButtonStyledComponent>
                    <SoloCancelButton
                      disabled={!formik.dirty}
                      onClick={() => formik.resetForm()}>
                      cancel
                    </SoloCancelButton>
                  </div>
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
            )}
          </Formik>
        </SectionCard>
      </div>
    </ErrorBoundary>
  );
};
