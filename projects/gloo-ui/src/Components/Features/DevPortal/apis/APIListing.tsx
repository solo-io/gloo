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
import { apiDocApi, portalApi } from '../api';
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
  console.log('apiDocsList', apiDocsList);
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

  const { data: portalsList, error: portalListError } = useSWR(
    `listPortalsApiDoc${apiDoc.metadata?.namespace}.${apiDoc.metadata?.name}`,
    portalApi.listPortals
  );

  return (
    <div
      onClick={() =>
        history.push(
          `/dev-portal/apis/${apiDoc.metadata?.namespace}/${apiDoc.metadata?.name}`
        )
      }
      className='relative flex mb-4 bg-white rounded-lg shadow cursor-pointer'>
      <span className='absolute top-0 right-0 flex items-center mt-3 mr-8 text-base font-medium text-gray-900'>
        Publish Status
        <HealthIndicator
          healthStatus={formatHealthStatus(apiDoc.status?.state)}
        />
      </span>
      <div className='items-center flex-none w-40 h-40 overflow-hidden text-center bg-cover rounded-l lg:rounded-t-none lg:rounded-l'>
        {apiDoc.spec?.image ? (
          <img
            className='object-cover max-h-72'
            src={`data:image/gif;base64,${apiDoc.spec?.image?.inlineBytes}`}></img>
        ) : (
          <PlaceholderPortal className='w-56 rounded-lg ' />
        )}
      </div>
      <div className='flex flex-col justify-around w-full h-40 ml-4'>
        <div className='text-lg text-gray-900 '>{apiDoc.metadata?.name}</div>
        <div className=''>{apiDoc.status?.description}</div>
        <div className='flex items-center justify-between '>
          <div className='font-medium text-gray-900 capitalize'>
            published in:{' '}
            {(portalsList || [])
              ?.map(p => p.spec!.displayName || p.metadata!.name)
              .sort((a, b) => (a === b ? 0 : a > b ? 1 : -1))
              .join(', ')}
          </div>
          <div className='flex items-center justify-between'>
            <div>
              {formatHealthStatus(apiDoc.status?.state) ===
                Status.State.PENDING && (
                <div className='flex items-center justify-center w-4 h-4 text-white text-orange-700 bg-orange-100 border border-orange-700 rounded-full'>
                  !
                </div>
              )}
            </div>
            <div className='flex items-center'>
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
    </div>
  );
};
