import React from 'react';
import { EmptyPortalsPanel } from '../DevPortal';
import { ReactComponent as PlaceholderPortalTile } from 'assets/portal-tile.svg';
import { useHistory } from 'react-router';
import { ReactComponent as PlaceholderPortal } from 'assets/placeholder-portal.svg';
import { ReactComponent as UserIcon } from 'assets/user-icon.svg';
import { ReactComponent as CodeIcon } from 'assets/code-icon.svg';
import { HealthIndicator } from 'Components/Common/HealthIndicator';
import { ReactComponent as GreenPlus } from 'assets/small-green-plus.svg';
import { CreateAPIModal } from './CreateAPIModal';
import { SoloModal } from 'Components/Common/SoloModal';
import { apiDocApi } from '../api';
import useSWR from 'swr';
import { ApiDoc } from 'proto/dev-portal/api/grpc/admin/apidoc_pb';
import { ApiDocStatus } from 'proto/dev-portal/api/dev-portal/v1/apidoc_pb';
import { formatHealthStatus } from '../portals/PortalsListing';
import { format } from 'timeago.js';
import { Status } from 'proto/solo-kit/api/v1/status_pb';

export const APIListing = () => {
  const { data: apiDocsList, error: apiDocsError } = useSWR(
    'listApiDocs',
    apiDocApi.listApiDocs
  );
  let isEmpty = false;
  const [showCreateApiModal, setShowCreateApiModal] = React.useState(false);

  if (!apiDocsList) {
    return <div>Loading...</div>;
  }
  return (
    <>
      <div className='container relative mx-auto '>
        <span
          onClick={() => setShowCreateApiModal(true)}
          className='absolute top-0 right-0 flex items-center -mt-8 text-green-400 cursor-pointer hover:text-green-300'>
          <GreenPlus className='mr-1 fill-current' />
          <span className='text-gray-700'> Create an API</span>
        </span>
        {isEmpty ? (
          <EmptyPortalsPanel itemName='API'>
            <PlaceholderPortalTile /> <PlaceholderPortalTile />
          </EmptyPortalsPanel>
        ) : (
          <>
            {apiDocsList
              .sort((a, b) =>
                a.metadata?.name === b.metadata?.name
                  ? 0
                  : a.metadata!.name > b.metadata!.name
                  ? 1
                  : -1
              )
              .map(apiDoc => (
                <APIItem key={apiDoc.metadata?.uid} apiDoc={apiDoc} />
              ))}
          </>
        )}
      </div>
      <SoloModal visible={showCreateApiModal} width={750} noPadding={true}>
        <CreateAPIModal onClose={() => setShowCreateApiModal(false)} />
      </SoloModal>
    </>
  );
};

const APIItem: React.FC<{ apiDoc: ApiDoc.AsObject }> = props => {
  const { apiDoc } = props;
  const history = useHistory();
  return (
    <div
      onClick={() => history.push(`/dev-portal/apis/${apiDoc.metadata?.name}`)}
      className='relative flex mb-4 bg-white rounded-lg shadow cursor-pointer'>
      <div className='flex-none h-32 overflow-hidden text-center bg-cover rounded-l lg:h-auto lg:w-56 lg:rounded-t-none lg:rounded-l'>
        <PlaceholderPortal className='rounded-l-lg ' />
      </div>
      <div className='flex flex-col ml-4 '>
        <div className='mb-2 text-lg text-gray-900'>
          {apiDoc.metadata?.name}
        </div>
        <span className='absolute top-0 right-0 flex items-center mt-3 mr-8 text-base font-medium text-gray-900'>
          Publish Status
          <HealthIndicator
            healthStatus={formatHealthStatus(apiDoc.status?.state)}
          />
        </span>
        <div className='my-2'>
          {' '}
          Nisi cupidatat commodo id incididunt. Laboris officia anim velit
          deserunt pariatur Lorem amet culpa sint velit amet fugiat occaecat. Do
          cillum ad cillum sint excepteur ea aute sint. Eiusmod sit eiusmod
          pariatur cupidatat laboris ullamco veniam minim. Aliquip minim ipsum
          voluptate labore eiusmod quis deserunt. Anim proident ad minim
          excepteur dolor reprehenderit. Qui ullamco commodo in laboris.{' '}
        </div>
        <div className='text-sm text-gray-600 '>
          Modified:
          {format(apiDoc.status?.modifiedDate?.seconds!, 'en_US')}
        </div>
        <div className='flex items-center justify-between mt-4'>
          <div className='font-medium text-gray-900 capitalize'>
            published in
          </div>
          <div className='flex items-center'>
            {formatHealthStatus(apiDoc.status?.state) ===
              Status.State.PENDING && (
              <div className='flex items-center justify-center w-4 h-4 text-white text-orange-700 bg-orange-100 border border-orange-700 rounded-full'>
                !
              </div>
            )}

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
              <span className='px-1 font-semibold'>
                {apiDoc.status?.numberOfEndpoints}
              </span>
              <span>Endpoints</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};
