import React from 'react';
import { SoloInput } from 'Components/Common/SoloInput';
import { ReactComponent as EditIcon } from 'assets/edit-pencil.svg';
import useSWR from 'swr';
import { apiDocApi } from '../api';
import { format } from 'timeago.js';
import { SoloModal } from 'Components/Common/SoloModal';
import { CreateAPIModal } from '../apis/CreateAPIModal';
import { ReactComponent as GreenPlus } from 'assets/small-green-plus.svg';
import { Portal } from 'proto/dev-portal/api/grpc/admin/portal_pb';
import { ApiDoc } from 'proto/dev-portal/api/grpc/admin/apidoc_pb';
import { secondsToString } from '../util';

type PortalApiDocsTabProps = {
  portal: Portal.AsObject;
};
export const PortalApiDocsTab = ({ portal }: PortalApiDocsTabProps) => {
  const { data: apiDocsList, error: apiDocsError } = useSWR(
    'listApiDocs',
    apiDocApi.listApiDocs
  );

  const [APISearchTerm, setAPISearchTerm] = React.useState('');
  const [showCreateApiModal, setShowCreateApiModal] = React.useState(false);

  const getApiDoc = (
    namespace: string,
    name: string
  ): ApiDoc.AsObject | undefined => {
    if (!apiDocsList) {
      return undefined;
    }

    return apiDocsList.find(
      doc => doc.metadata?.namespace == namespace && doc.metadata.name == name
    );
  };

  return (
    <div className='relative flex flex-col p-4 border border-gray-300 rounded-lg'>
      <span
        onClick={() => setShowCreateApiModal(true)}
        className='absolute top-0 right-0 flex items-center mt-2 mr-2 text-green-400 cursor-pointer hover:text-green-300'>
        <GreenPlus className='mr-1 fill-current' />
        <span className='text-gray-700'> Create an API</span>
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
                {portal.status?.apiDocsList.map(ref => {
                  const doc = getApiDoc(ref.namespace, ref.name);
                  if (!doc) {
                    return;
                  }
                  return (
                    <tr
                      key={`${doc.metadata?.namespace}.${doc.metadata?.name}`}>
                      <td className='px-6 py-4 whitespace-no-wrap border-b border-gray-200'>
                        <div className='text-sm leading-5 text-gray-900'>
                          <span className='flex items-center capitalize'>
                            {doc?.status?.displayName}
                          </span>
                        </div>
                      </td>
                      <td className='px-6 py-4 whitespace-no-wrap border-b border-gray-200'>
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
                          <div className='flex items-center justify-center w-4 h-4 mr-3 text-gray-700 bg-gray-400 rounded-full cursor-pointer'>
                            <EditIcon className='w-2 h-3 fill-current' />
                          </div>
                          <div
                            className='flex items-center justify-center w-4 h-4 text-gray-700 bg-gray-400 rounded-full cursor-pointer'
                            onClick={() => {}}>
                            x
                          </div>
                        </span>
                      </td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
          </div>
        </div>
      </div>
      <SoloModal visible={showCreateApiModal} width={750} noPadding={true}>
        <CreateAPIModal onClose={() => setShowCreateApiModal(false)} />
      </SoloModal>
    </div>
  );
};
