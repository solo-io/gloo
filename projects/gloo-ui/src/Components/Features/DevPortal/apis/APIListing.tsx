import { ReactComponent as CodeIcon } from 'assets/code-icon.svg';
import { ReactComponent as PlaceholderPortal } from 'assets/placeholder-portal.svg';
import { ReactComponent as GreenPlus } from 'assets/small-green-plus.svg';
import { ReactComponent as UserIcon } from 'assets/user-icon.svg';
import { Loading } from 'Components/Common/DisplayOnly/Loading';
import { HealthIndicator } from 'Components/Common/HealthIndicator';
import { SoloModal } from 'Components/Common/SoloModal';
import { ApiDoc } from '@solo-io/dev-portal-grpc/dev-portal/api/grpc/admin/apidoc_pb';
import { Status } from 'proto/solo-kit/api/v1/status_pb';
import React from 'react';
import { useHistory } from 'react-router';
import useSWR from 'swr';
import { apiDocApi, portalApi, userApi } from '../api';
import { NoDataPanel } from '../DevPortal';
import { formatHealthStatus } from '../portals/PortalsListing';
import { CreateAPIModal } from './CreateAPIModal';

export const APIListing = () => {
  const { data: apiDocsList, error: apiDocsError } = useSWR(
    'listApiDocs',
    apiDocApi.listApiDocs
  );

  const [showCreateApiModal, setShowCreateApiModal] = React.useState(false);

  if (!apiDocsList) {
    return <Loading center>Loading...</Loading>;
  }
  return (
    <div className='container relative mx-auto '>
      <span
        onClick={() => setShowCreateApiModal(true)}
        className='absolute top-0 right-0 flex items-center -mt-8 text-green-400 cursor-pointer hover:text-green-300'>
        <GreenPlus className='mr-1 fill-current' />
        <span className='text-gray-700'> Create an API</span>
      </span>
      {apiDocsList.length === 0 ? (
        <NoDataPanel
          missingContentText='There are no APIs to display'
          helpText='Create a API to publish and share with developers.'
          identifier='apis-page'
        />
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
      <SoloModal visible={showCreateApiModal} width={750} noPadding={true}>
        <CreateAPIModal onClose={() => setShowCreateApiModal(false)} />
      </SoloModal>
    </div>
  );
};

const APIItem: React.FC<{ apiDoc: ApiDoc.AsObject }> = props => {
  const { apiDoc } = props;
  const { data: usersList, error: usersError } = useSWR(
    `listUsers${apiDoc.metadata?.name}${apiDoc.metadata?.namespace}`,
    () =>
      userApi.listUsers({
        apiDocsList: [
          {
            name: apiDoc.metadata!.name,
            namespace: apiDoc.metadata!.namespace
          }
        ],
        portalsList: [],
        groupsList: []
      })
  );
  const history = useHistory();

  const { data: portalsList, error: portalListError } = useSWR(
    `listPortalsApiDoc${apiDoc.metadata?.namespace}.${apiDoc.metadata?.name}`,
    portalApi.listPortals
  );

  const filteredPortalList = portalsList?.filter(portal =>
    portal.status?.apiDocsList.some(
      apiDocRef =>
        apiDocRef.name === apiDoc.metadata?.name &&
        apiDocRef.namespace === apiDoc.metadata.namespace
    )
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
        {apiDoc.spec?.image?.inlineBytes ? (
          <img
            className='object-contain h-40'
            src={`data:image/gif;base64,${apiDoc.spec?.image?.inlineBytes}`}></img>
        ) : (
          <PlaceholderPortal className='h-40 rounded-lg w-60 ' />
        )}
      </div>
      <div className='flex flex-col justify-around w-full h-40 ml-4'>
        <div className='text-lg text-gray-900 '>
          {apiDoc.status?.displayName || apiDoc.metadata?.name}
        </div>
        <div className=''>{apiDoc.status?.description}</div>
        <div className='flex items-center justify-between '>
          <div className='font-medium text-gray-900 capitalize'>
            published in:{' '}
            {(filteredPortalList || [])
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
                  <UserIcon className='w-6 h-6 fill-current' />
                </span>
                <span className='px-1 font-semibold'>
                  {usersList?.length || 0}
                </span>
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
    </div>
  );
};
