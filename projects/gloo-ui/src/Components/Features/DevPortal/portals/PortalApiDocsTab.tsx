import React from 'react';
import { SoloInput } from 'Components/Common/SoloInput';
import { ReactComponent as EditIcon } from 'assets/edit-pencil.svg';
import useSWR from 'swr';
import { apiDocApi } from '../api';
import { format } from 'timeago.js';
import { SoloModal } from 'Components/Common/SoloModal';
import { CreateAPIModal } from '../apis/CreateAPIModal';
import { ReactComponent as GreenPlus } from 'assets/small-green-plus.svg';
import { ReactComponent as CodeIcon } from 'assets/code-icon.svg';

import { Portal } from '@solo-io/dev-portal-grpc/dev-portal/api/grpc/admin/portal_pb';
import { ApiDoc } from '@solo-io/dev-portal-grpc/dev-portal/api/grpc/admin/apidoc_pb';
import { secondsToString } from '../util';
import { AddApiModal } from './AddApiModal';
import { ObjectRef } from '@solo-io/dev-portal-grpc/dev-portal/api/dev-portal/v1/common_pb';
import { ConfirmationModal } from 'Components/Common/ConfirmationModal';

type PortalApiDocsTabProps = {
  portal: Portal.AsObject;
};
export const PortalApiDocsTab = ({ portal }: PortalApiDocsTabProps) => {
  const { data: apiDocsList, error: apiDocsError } = useSWR(
    `listApiDocs${portal.metadata?.name}${portal.metadata?.namespace}`,
    () =>
      apiDocApi.listApiDocs({
        portalsList: [
          {
            name: portal.metadata!.name,
            namespace: portal.metadata!.namespace
          }
        ]
      })
  );
  const getApiDocFromRef = (apiDocRef: ObjectRef.AsObject) => {
    let apiDocObj = apiDocsList?.find(
      apiDoc =>
        apiDoc.metadata?.name === apiDocRef.name &&
        apiDoc.metadata.namespace === apiDocRef.namespace
    );
    return apiDocObj;
  };

  const [APISearchTerm, setAPISearchTerm] = React.useState('');
  const [showCreateApiModal, setShowCreateApiModal] = React.useState(false);

  const [showConfirmApiDocDelete, setShowConfirmApiDocDelete] = React.useState(
    false
  );
  const [apiDocToDelete, setApiDocToDelete] = React.useState<ApiDoc.AsObject>();

  const [filteredAPIs, setFilteredAPIs] = React.useState<
    ObjectRef.AsObject[]
  >();
  React.useEffect(() => {
    if (APISearchTerm !== '' && !!portal.status?.apiDocsList) {
      setFilteredAPIs(
        portal.status?.apiDocsList.filter(user =>
          user.name.toLowerCase().includes(APISearchTerm)
        )
      );
    } else {
      setFilteredAPIs(undefined);
    }
  }, [APISearchTerm]);
  const attemptDeleteApiDoc = (apiDoc: ApiDoc.AsObject) => {
    setShowConfirmApiDocDelete(true);
    setApiDocToDelete(apiDoc);
  };

  const deleteApiDoc = async () => {
    await apiDocApi.deleteApiDoc({
      name: apiDocToDelete?.metadata?.name!,
      namespace: apiDocToDelete?.metadata?.namespace!
    });
    closeConfirmModal();
  };

  const closeConfirmModal = () => {
    setShowConfirmApiDocDelete(false);
  };
  return (
    <div className='relative flex flex-col p-4 border border-gray-300 rounded-lg'>
      <span
        onClick={() => setShowCreateApiModal(true)}
        className='absolute top-0 right-0 flex items-center mt-2 mr-2 text-green-400 cursor-pointer hover:text-green-300'>
        <GreenPlus className='mr-1 fill-current' />
        <span className='text-gray-700'> Add/Remove an API</span>
      </span>
      <div className='w-1/3 m-4'>
        <SoloInput
          placeholder='Search by API name...'
          value={APISearchTerm}
          onChange={e => setAPISearchTerm(e.target.value)}
        />
      </div>
      <div className='flex flex-col'>
        <div className='py-2 -my-2 overflow-x-auto sm:-mx-6 sm:px-6 lg:-mx-8 lg:px-8'>
          <div className='inline-block min-w-full overflow-hidden align-middle border-b border-gray-200 shadow sm:rounded-lg'>
            <table className='min-w-full'>
              <thead className='bg-gray-300 '>
                <tr>
                  <th className='px-6 py-3 text-sm font-medium leading-4 tracking-wider text-left text-gray-800 capitalize border-b border-gray-200 bg-gray-50'>
                    API Name
                  </th>
                  <th className='px-6 py-3 text-sm font-medium leading-4 tracking-wider text-left text-gray-800 capitalize border-b border-gray-200 bg-gray-50'>
                    Description
                  </th>

                  <th className='px-6 py-3 text-sm font-medium leading-4 tracking-wider text-left text-gray-800 capitalize border-b border-gray-200 bg-gray-50'>
                    Modified
                  </th>

                  <th className='px-6 py-3 text-sm font-medium leading-4 tracking-wider text-left text-gray-800 capitalize border-b border-gray-200 bg-gray-50'>
                    Actions
                  </th>
                </tr>
              </thead>
              <tbody className='bg-white'>
                {(!!filteredAPIs ? filteredAPIs : portal.status?.apiDocsList)
                  ?.sort((a, b) =>
                    a.name === b.name ? 0 : a.name > b.name ? 1 : -1
                  )
                  .map(apiDocRef => {
                    const doc = getApiDocFromRef(apiDocRef);

                    return (
                      <tr
                        key={`${doc?.metadata?.namespace}.${doc?.metadata?.name}`}>
                        <td className='px-6 py-4 whitespace-no-wrap border-b border-gray-200'>
                          <div className='text-sm leading-5 text-gray-900'>
                            <span className='flex items-center capitalize'>
                              {doc?.status?.displayName}
                            </span>
                          </div>
                        </td>
                        <td className='px-6 py-4 border-b border-gray-200'>
                          <div className='text-sm leading-5 text-gray-900'>
                            <span className='flex items-center capitalize'>
                              {doc?.status?.description}
                            </span>
                          </div>
                        </td>
                        <td className='px-6 py-4 whitespace-no-wrap border-b border-gray-200'>
                          <div className='text-sm leading-5 text-gray-900'>
                            <span className='flex items-center '>
                              {format(
                                secondsToString(
                                  doc?.status?.modifiedDate?.seconds
                                )
                              )}
                            </span>
                          </div>
                        </td>
                        <td className='px-6 py-4 text-sm font-medium leading-5 text-right whitespace-no-wrap border-b border-gray-200'>
                          <span className='flex items-center'>
                            <div
                              className='flex items-center justify-center w-4 h-4 text-gray-700 bg-gray-400 rounded-full cursor-pointer'
                              onClick={() => attemptDeleteApiDoc(doc!)}>
                              x
                            </div>
                          </span>
                        </td>
                      </tr>
                    );
                  })}
              </tbody>
            </table>
            {(!!filteredAPIs ? filteredAPIs : portal.status?.apiDocsList)
              ?.length === 0 && (
              <div className='w-full m-auto'>
                <div className='flex flex-col items-center justify-center w-full h-full py-4 mr-32 bg-white rounded-lg shadow-lg md:flex-row'>
                  <div className='mr-6 text-blue-600'>
                    <CodeIcon className='w-20 h-20 fill-current ' />
                  </div>
                  <div className='flex flex-col h-full'>
                    <p className='h-auto my-6 text-lg font-medium text-gray-800 '>
                      There are no APIs to display!{' '}
                    </p>
                  </div>
                </div>
              </div>
            )}
          </div>
        </div>
      </div>
      <SoloModal visible={showCreateApiModal} width={750} noPadding={true}>
        <AddApiModal onClose={() => setShowCreateApiModal(false)} />
      </SoloModal>
      <ConfirmationModal
        visible={showConfirmApiDocDelete}
        confirmationTopic='delete this API'
        confirmText='Delete'
        goForIt={deleteApiDoc}
        cancel={closeConfirmModal}
        isNegative={true}
      />
    </div>
  );
};
