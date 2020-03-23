import React from 'react';
import { EmptyPortalsPanel } from '../DevPortal';
import { ReactComponent as PlaceholderPortalTile } from 'assets/portal-tile.svg';
import { useHistory } from 'react-router';
import { ReactComponent as PlaceholderPortal } from 'assets/placeholder-portal.svg';
import { ReactComponent as UserIcon } from 'assets/user-icon.svg';
import { ReactComponent as CodeIcon } from 'assets/code-icon.svg';
import { HealthIndicator } from 'Components/Common/HealthIndicator';
import { css } from '@emotion/core';

export const APIListing = () => {
  let isEmpty = false;
  return (
    <div className='container mx-auto'>
      {isEmpty ? (
        <EmptyPortalsPanel itemName='API'>
          <PlaceholderPortalTile /> <PlaceholderPortalTile />
        </EmptyPortalsPanel>
      ) : (
        <>
          <APIItem />
          <APIItem />
        </>
      )}
    </div>
  );
};

const APIItem = () => {
  const history = useHistory();
  return (
    <div
      onClick={() => history.push(`/dev-portal/apis/${'APIName'}`)}
      className='relative flex mb-4 bg-white rounded-lg shadow cursor-pointer'>
      <PlaceholderPortal className='w-1/3 rounded-l-lg h-1/2 ' />

      <div className='flex flex-col ml-4 '>
        <div className='mb-2 text-lg text-gray-800'>API Name</div>
        <span className='absolute top-0 right-0 flex items-center mt-3 mr-8 text-base font-medium text-gray-900'>
          Published
          <HealthIndicator healthStatus={1} />
        </span>
        <div className='my-2'>
          Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam
          nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat,
          sed diam voluptua. At vero eos et accusam et justo duo dolores et ea
          rebum.
        </div>
        <div className='flex items-center justify-between mt-4'>
          <div>published in</div>
          <div className='flex items-center'>
            <div className='flex items-center justify-center w-4 h-4 text-white text-orange-700 bg-orange-100 border border-orange-700 rounded-full'>
              !
            </div>
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
