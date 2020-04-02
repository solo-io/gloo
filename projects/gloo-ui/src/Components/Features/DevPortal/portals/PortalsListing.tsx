import { Breadcrumb } from 'Components/Common/Breadcrumb';
import React from 'react';
import { EmptyPortalsPanel } from '../DevPortal';
import { ReactComponent as PlaceholderPortalTile } from 'assets/portal-tile.svg';
import { ReactComponent as PlaceholderPortal } from 'assets/placeholder-portal.svg';
import { ReactComponent as ExternalLinkIcon } from 'assets/external-link-icon.svg';
import { ReactComponent as CompanyLogo } from 'assets/company-logo.svg';
import { ReactComponent as UserIcon } from 'assets/user-icon.svg';
import { ReactComponent as CodeIcon } from 'assets/code-icon.svg';
import { useHistory } from 'react-router';
import { HealthIndicator } from 'Components/Common/HealthIndicator';
import { css } from '@emotion/core';
import useSWR from 'swr';
import { portalApi, userApi, apiDocApi } from '../api';
import { Portal } from 'proto/dev-portal/api/grpc/admin/portal_pb';
import { PortalStatus } from 'proto/dev-portal/api/dev-portal/v1/portal_pb';
import { Status } from 'proto/solo-kit/api/v1/status_pb';
import { format } from 'timeago.js';
import { CreatePortalModal } from './CreatePortalModal';
import { SoloModal } from 'Components/Common/SoloModal';
import { ReactComponent as GreenPlus } from 'assets/small-green-plus.svg';
import { StateMap, State } from 'proto/dev-portal/api/dev-portal/v1/common_pb';
import { secondsToString } from '../util';

export function formatHealthStatus(
  status: StateMap[keyof StateMap] | undefined
): Status.StateMap[keyof Status.StateMap] {
  if (status === State.PENDING || status === State.PROCESSING) {
    return Status.State.PENDING;
  } else if (status === State.FAILED || status === State.INVALID) {
    return Status.State.REJECTED;
  }
  return Status.State.ACCEPTED;
}

export const PortalsListing = () => {
  const { data: portalsList, error: portalListError } = useSWR(
    'listPortals',
    portalApi.listPortals
  );
  const [showCreatePortalModal, setShowCreatePortalModal] = React.useState(
    false
  );

  if (!portalsList) {
    return <div>Loading...</div>;
  }
  return (
    <div className='container mx-auto'>
      <span
        onClick={() => setShowCreatePortalModal(true)}
        className='absolute top-0 right-0 flex items-center -mt-8 text-green-400 cursor-pointer hover:text-green-300'>
        <GreenPlus className='mr-1 fill-current' />
        <span className='text-gray-700'> Create a Portal</span>
      </span>
      {portalsList.length === 0 ? (
        <EmptyPortalsPanel itemName='Portal'>
          <PlaceholderPortalTile /> <PlaceholderPortalTile />
        </EmptyPortalsPanel>
      ) : (
        <>
          {portalsList
            .sort((a, b) =>
              a.metadata?.name === b.metadata?.name
                ? 0
                : a.metadata!.name > b.metadata!.name
                ? 1
                : -1
            )
            .map(portal => (
              <PortalItem key={portal.metadata?.uid} portal={portal} />
            ))}
        </>
      )}
      <SoloModal visible={showCreatePortalModal} width={750} noPadding={true}>
        <CreatePortalModal onClose={() => setShowCreatePortalModal(false)} />
      </SoloModal>
    </div>
  );
};

const PortalItem: React.FC<{ portal: Portal.AsObject }> = props => {
  const { portal } = props;
  const history = useHistory();

  const { data: usersList, error: usersListError } = useSWR(
    `listUsers${portal.metadata?.name}${portal.metadata?.namespace}`,
    () =>
      userApi.listUsers({
        portalsList: [
          { namespace: portal.metadata!.namespace, name: portal.metadata!.name }
        ],
        groupsList: [],
        apiDocsList: []
      }),
    { refreshInterval: 0 }
  );
  const { data: apiDocsList, error: apiDocsListError } = useSWR(
    `listApiDocs${portal.metadata?.name}${portal.metadata?.namespace}`,
    () =>
      apiDocApi.listApiDocs({
        portalsList: [
          { namespace: portal.metadata!.namespace, name: portal.metadata!.name }
        ]
      }),
    { refreshInterval: 0 }
  );

  return (
    <div
      className='w-full max-w-md mb-4 rounded-lg shadow lg:max-w-full lg:flex'
      onClick={() =>
        history.push(
          `/dev-portal/portals/${portal.metadata?.namespace}/${portal.metadata?.name}`
        )
      }>
      <div className='flex-none h-48 overflow-hidden text-center bg-cover rounded-l lg:h-auto lg:w-56 lg:rounded-t-none lg:rounded-l'>
        {portal.spec?.banner?.inlineBytes ? (
          <img
            className='object-cover h-48'
            src={`data:image/gif;base64,${portal.spec?.banner?.inlineBytes}`}></img>
        ) : (
          <PlaceholderPortal className='rounded-lg w-72 ' />
        )}
      </div>
      <div className='relative flex flex-col justify-between w-full p-2 leading-normal bg-white rounded-r'>
        <span className='absolute top-0 right-0 flex items-center mt-3 mr-8 text-base font-medium text-gray-900'>
          Portal Status
          <HealthIndicator
            healthStatus={formatHealthStatus(portal.status?.state)}
          />
        </span>
        <div className='mb-4'>
          <div className='text-xl text-gray-900 '>
            {portal.spec?.displayName}
          </div>

          {portal.spec?.domainsList.map(domain => (
            <p
              className='flex items-center py-1 text-base text-blue-600'
              key={domain}>
              <ExternalLinkIcon className='w-4 h-4 mr-1' /> {domain}
            </p>
          ))}
          <p className='text-base text-gray-700 break-all'>
            {portal.spec?.description}
          </p>
        </div>
        <div className='flex items-center justify-end'>
          <div className='flex items-center '>
            <div className='flex items-center px-4'>
              <span className='text-blue-600'>
                <UserIcon className='w-6 h-6 fill-current' />
              </span>
              <span className='px-1 font-semibold'>
                {!!usersList ? usersList.length : '...'}
              </span>
              <span>
                {`User${!!usersList && usersList.length === 1 ? '' : 's'}`}
              </span>
            </div>
            <div className='flex items-center px-4'>
              <span className='text-blue-600'>
                <CodeIcon className='w-6 h-6 fill-current' />
              </span>
              <span className='px-1 font-semibold'>
                {!!apiDocsList ? apiDocsList.length : '...'}
              </span>
              <span>{`API${
                !!apiDocsList && apiDocsList.length === 1 ? '' : 's'
              }`}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};
