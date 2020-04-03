import { css } from '@emotion/core';
import {
  Tab,
  TabList,
  TabPanel,
  TabPanelProps,
  TabPanels,
  Tabs
} from '@reach/tabs';
import { ReactComponent as CodeIcon } from 'assets/code-icon.svg';
import { ReactComponent as EditIcon } from 'assets/edit-pencil.svg';
import { ReactComponent as ExternalLinkIcon } from 'assets/external-link-icon.svg';
import { ReactComponent as PlaceholderPortal } from 'assets/placeholder-portal.svg';
import { Breadcrumb } from 'Components/Common/Breadcrumb';
import { ConfirmationModal } from 'Components/Common/ConfirmationModal';
import { SectionCard } from 'Components/Common/SectionCard';
import { SoloInput } from 'Components/Common/SoloInput';
import { ErrorBoundary } from 'Components/Features/Errors/ErrorBoundary';
import {
  State,
  DataSource
} from 'proto/dev-portal/api/dev-portal/v1/common_pb';
import React from 'react';
import { ReactComponent as NoImageIcon } from 'assets/no-image-placeholder.svg';

import { useHistory, useParams } from 'react-router';
import {
  SoloNegativeButton,
  SoloButtonStyledComponent,
  SoloCancelButton
} from 'Styles/CommonEmotions/button';
import useSWR from 'swr';
import { apiDocApi, portalApi } from '../api';
import { ActiveTabCss, TabCss } from '../portals/PortalDetails';
import { formatHealthStatus } from '../portals/PortalsListing';
import { Loading } from 'Components/Common/DisplayOnly/Loading';
import { APIUsersTab } from './APIUsers';
import { APIGroupsTab } from './APIGroups';
import {
  SoloFormInput,
  SoloFormTextarea
} from 'Components/Common/Form/SoloFormField';
import { Formik } from 'formik';
import { format } from 'timeago.js';
import { secondsToString } from '../util';
import { colors } from 'Styles';
import ImageUploader from 'react-images-upload';
import { SoloButton } from 'Components/Common/SoloButton';

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

type UpdateApiValues = {
  displayName: string;
  description: string;
  image: File;
};
export const APIDetails = () => {
  const { apiname, apinamespace } = useParams();
  const { data: apiDoc, error: apiDocError } = useSWR(
    !!apiname && !!apinamespace ? ['getApiDoc', apiname, apinamespace] : null,
    (key, name, namespace) =>
      apiDocApi.getApiDoc({ apidoc: { name, namespace }, withassets: true })
  );
  const { data: portalsList, error: portalListError } = useSWR(
    'listPortals',
    portalApi.listPortals
  );

  const history = useHistory();
  const [tabIndex, setTabIndex] = React.useState(0);
  const [APISearchTerm, setAPISearchTerm] = React.useState('');
  const [attemptingDelete, setAttemptingDelete] = React.useState(false);

  const attemptDeleteApiDoc = () => {
    setAttemptingDelete(true);
  };

  const cancelDeletion = () => {
    setAttemptingDelete(false);
  };

  const deleteApi = async () => {
    await apiDocApi.deleteApiDoc({
      name: apiDoc?.metadata?.name!,
      namespace: apiDoc?.metadata?.namespace!
    });
    setAttemptingDelete(false);
    history.push('/dev-portal/apis');
  };

  const handleTabsChange = (index: number) => {
    setTabIndex(index);
  };

  const filteredPortalList = portalsList?.filter(portal =>
    portal.status?.apiDocsList.some(
      apiDocRef =>
        apiDocRef.name === apiDoc?.metadata?.name &&
        apiDocRef.namespace === apiDoc?.metadata.namespace
    )
  );

  const goToEdit = () => {
    history.push(`/dev-portal/apis/${apinamespace}/${apiname}/edit`);
  };

  if (!apiDoc) {
    return <Loading center>Loading...</Loading>;
  }
  return (
    <ErrorBoundary
      fallback={<div>There was an error with the Dev Portal section</div>}>
      <div>
        <Breadcrumb />
        <SectionCard
          cardName={apiDoc.status?.displayName || apiname || 'API'}
          logoIcon={
            <span className='text-blue-500'>
              <CodeIcon className='fill-current' />
            </span>
          }
          health={formatHealthStatus(apiDoc?.status?.state)}
          headerSecondaryInformation={[
            {
              title: 'Modified',
              value: format(
                secondsToString(apiDoc?.status?.modifiedDate?.seconds)
              )
            }
          ]}
          healthMessage={'API Status'}
          onClose={() => history.push(`/dev-portal/`)}>
          <div>
            {apiDoc?.status?.state !== State.SUCCEEDED && (
              <div className='flex items-center p-2 mb-2 text-yellow-500 bg-yellow-100 border border-yellow-500 rounded-lg '>
                <div className='flex items-center justify-center w-4 h-4 mr-2 text-white text-yellow-500 bg-orange-100 border border-yellow-500 rounded-full'>
                  !
                </div>{' '}
                {apiDoc.status?.reason}
              </div>
            )}
            <div className='relative flex items-center'>
              <div className=' max-h-72'>
                {apiDoc.spec?.image?.inlineBytes ? (
                  <img
                    className='object-cover max-h-72'
                    src={`data:image/gif;base64,${apiDoc.spec?.image?.inlineBytes}`}></img>
                ) : (
                  <PlaceholderPortal className='w-56 rounded-lg ' />
                )}
              </div>
              <div className='grid w-full grid-cols-2 ml-2 h-36'>
                <div>
                  <span className='font-medium text-gray-900'>
                    Display Name
                  </span>

                  <div>
                    {apiDoc?.status?.displayName || apiDoc?.metadata?.name}
                  </div>
                </div>
                <div>
                  <span className='font-medium text-gray-900'>
                    Published In
                  </span>
                  <div className='grid w-1/2 grid-flow-col grid-flow-col-dense grid-cols-2'>
                    {(filteredPortalList || [])
                      .sort((a, b) =>
                        a.metadata?.name === b.metadata?.name
                          ? 0
                          : a.metadata!.name > b.metadata!.name
                          ? 1
                          : -1
                      )
                      .map((portal, index) => (
                        <div
                          key={portal.metadata?.uid}
                          className='flex items-center mb-2 text-sm text-blue-600'>
                          <div>{portal.spec?.displayName}</div>
                        </div>
                      ))}
                  </div>
                </div>

                <div className='col-span-2 '>
                  <span className='font-medium text-gray-900'>Description</span>

                  <div className='break-words '>
                    {apiDoc?.status?.description}
                  </div>
                </div>
              </div>
            </div>
            <Tabs
              index={tabIndex}
              className='mt-6 mb-4 border-none rounded-lg'
              onChange={handleTabsChange}>
              <TabList className='flex items-start ml-4 '>
                <StyledTab>Users</StyledTab>
                <StyledTab>Groups</StyledTab>
              </TabList>
              <TabPanels
                css={css`
                  margin-top: -1px;
                `}>
                <TabPanel className='focus:outline-none'>
                  <APIUsersTab apiDoc={apiDoc} />
                </TabPanel>
                <TabPanel className='focus:outline-none'>
                  <APIGroupsTab apiDoc={apiDoc} />
                </TabPanel>
              </TabPanels>
            </Tabs>
            <div className='flex justify-end justify-between items-bottom'>
              <SoloButton text='Open API Editor' onClick={goToEdit} />

              <SoloNegativeButton onClick={attemptDeleteApiDoc}>
                Delete API
              </SoloNegativeButton>
            </div>
            <ConfirmationModal
              visible={attemptingDelete}
              confirmationTopic='delete this API'
              confirmText='Delete'
              goForIt={deleteApi}
              cancel={cancelDeletion}
              isNegative={true}
            />
          </div>
        </SectionCard>
      </div>
    </ErrorBoundary>
  );
};
