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
import { DevPortalApi } from '../api';
import { Portal } from 'proto/dev-portal/api/grpc/admin/portal_pb';
import { PortalStatus } from 'proto/dev-portal/api/dev-portal/v1/portal_pb';
import { Status } from 'proto/solo-kit/api/v1/status_pb';

function formatHealthStatus(
  status: PortalStatus.StateMap[keyof PortalStatus.StateMap] | undefined
): Status.StateMap[keyof Status.StateMap] {
  if (
    status === PortalStatus.State.PENDING ||
    status === PortalStatus.State.PROCESSING
  ) {
    return Status.State.PENDING;
  } else if (
    status === PortalStatus.State.FAILED ||
    status === PortalStatus.State.INVALID
  ) {
    return Status.State.REJECTED;
  }
  return Status.State.ACCEPTED;
}

export const PortalsListing = () => {
  let isEmpty = false;
  const { data: portalsList, error: portalListError } = useSWR(
    'listPortals',
    DevPortalApi.listPortals,
    { refreshInterval: 0 }
  );
  if (!portalsList) {
    return <div>Loading...</div>;
  }
  console.log('portalsList', portalsList);
  return (
    <div className='container mx-auto'>
      {isEmpty ? (
        <EmptyPortalsPanel itemName='Portal'>
          <PlaceholderPortalTile /> <PlaceholderPortalTile />
        </EmptyPortalsPanel>
      ) : (
        <>
          {portalsList.map(portal => (
            <PortalItem key={portal.metadata?.uid} portal={portal} />
          ))}
        </>
      )}
    </div>
  );
};

const PortalItem: React.FC<{ portal: Portal.AsObject }> = props => {
  const { portal } = props;
  const history = useHistory();
  return (
    <div
      className='w-full max-w-md rounded-lg shadow lg:max-w-full lg:flex'
      onClick={() =>
        history.push(`/dev-portal/portals/${portal.metadata?.name}`)
      }>
      <div
        className='flex-none h-48 overflow-hidden text-center bg-cover rounded-l lg:h-auto lg:w-56 lg:rounded-t-none lg:rounded-l'
        title='Woman holding a mug'>
        <PlaceholderPortal className='rounded-l-lg ' />
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
            {' '}
            {portal.spec?.displayName}
          </div>

          <p className='flex items-center py-1 text-base text-blue-600'>
            <ExternalLinkIcon className='w-4 h-4 mr-1' />
            https://production.subdomain.gloo.io
          </p>
          <p className='text-base text-gray-700'>{portal.spec?.description}</p>
        </div>
        <div className='text-sm text-gray-600 '>Modified: Feb 26, 2020</div>
        <div className='flex items-center justify-between'>
          <div className='pb-2'>
            <CompanyLogo className='w-1/2 h-1/2' />
          </div>
          <div className='flex items-center '>
            <div className='flex items-center px-4'>
              <span className='text-blue-600'>
                <UserIcon className='w-6 h-6 fill-current' />
              </span>
              <span className='px-1 font-semibold'>10</span>
              <span>Users</span>
            </div>
            <div className='flex items-center px-4'>
              <span className='text-blue-600'>
                <CodeIcon className='w-6 h-6 fill-current' />
              </span>
              <span className='px-1 font-semibold'>5</span>
              <span>APIs</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};
