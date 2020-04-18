import React from 'react';
import { useParams, useHistory } from 'react-router';
import { SectionCard } from 'Components/Common/SectionCard';
import { ReactComponent as NoApiKey } from 'assets/no-api-key-icon.svg';
import { ReactComponent as KeyIcon } from 'assets/key-on-ring.svg';
import { SoloInput } from 'Components/Common/SoloInput';
import useSWR from 'swr';
import { apiKeyApi } from '../api';
import { Loading } from 'Components/Common/DisplayOnly/Loading';
import { ConfirmationModal } from 'Components/Common/ConfirmationModal';
import { ApiKey } from '@solo-io/dev-portal-grpc/dev-portal/api/grpc/admin/api_key_pb';

export const APIKeys = () => {
  const { data: apiKeyList, error: apiKeysError } = useSWR(
    'listApiKeys',
    apiKeyApi.listApiKeys
  );
  const history = useHistory();
  const [apiKeySearchTerm, setApiKeySearchTerm] = React.useState('');
  const [attemptingDelete, setAttemptingDelete] = React.useState(false);
  const [apiKeyToDelete, setApiKeyToDelete] = React.useState<ApiKey.AsObject>();

  const [filteredAPIKeys, setFilteredAPIKeys] = React.useState<
    ApiKey.AsObject[]
  >();
  React.useEffect(() => {
    if (apiKeySearchTerm !== '' && !!apiKeyList) {
      setFilteredAPIKeys(
        apiKeyList.filter(apiKey =>
          apiKey.metadata?.name.toLowerCase().includes(apiKeySearchTerm)
        )
      );
    } else {
      setFilteredAPIKeys(undefined);
    }
  }, [apiKeySearchTerm]);
  const attemptDeleteApiKey = (apiKey: ApiKey.AsObject) => {
    setAttemptingDelete(true);
    setApiKeyToDelete(apiKey);
  };
  const cancelDeletion = () => {
    setAttemptingDelete(false);
  };

  const deleteApiKey = async () => {
    await apiKeyApi.deleteApiKey({
      name: apiKeyToDelete?.metadata?.name!,
      namespace: apiKeyToDelete?.metadata?.namespace!
    });
    setAttemptingDelete(false);
    history.push('/dev-portal/api-keys');
  };
  if (!apiKeyList) {
    return <Loading center>Loading...</Loading>;
  }
  return (
    <SectionCard
      cardName={'API Keys'}
      logoIcon={
        <span className='text-blue-500'>
          <KeyIcon className='fill-current' />
        </span>
      }>
      <div className='relative flex flex-col p-2 rounded-lg'>
        <div className='w-full mb-4'>
          <SoloInput
            placeholder='Search API keys by name...'
            value={apiKeySearchTerm}
            onChange={e => setApiKeySearchTerm(e.target.value)}
          />
        </div>
        <div className='flex flex-col'>
          <div className='py-2 -my-2 overflow-x-auto sm:-mx-6 sm:px-6 lg:-mx-8 lg:px-8'>
            <div className='inline-block min-w-full overflow-hidden align-middle border-b border-gray-200 shadow sm:rounded-lg'>
              <table className='min-w-full'>
                <thead className='bg-gray-300 '>
                  <tr>
                    <th className='px-6 py-3 text-sm font-medium leading-4 tracking-wider text-left text-gray-800 capitalize border-b border-gray-200 bg-gray-50'>
                      Secret
                    </th>
                    <th className='px-6 py-3 text-sm font-medium leading-4 tracking-wider text-left text-gray-800 capitalize border-b border-gray-200 bg-gray-50'>
                      Key
                    </th>

                    <th className='px-6 py-3 text-sm font-medium leading-4 tracking-wider text-left text-gray-800 capitalize border-b border-gray-200 bg-gray-50'>
                      User
                    </th>
                    <th className='px-6 py-3 text-sm font-medium leading-4 tracking-wider text-left text-gray-800 capitalize border-b border-gray-200 bg-gray-50'>
                      API Key Scope
                    </th>
                    <th className='px-6 py-3 text-sm font-medium leading-4 tracking-wider text-left text-gray-800 capitalize border-b border-gray-200 bg-gray-50'>
                      Labels
                    </th>

                    <th className='px-6 py-3 text-sm font-medium leading-4 tracking-wider text-left text-gray-800 capitalize border-b border-gray-200 bg-gray-50'>
                      Actions
                    </th>
                  </tr>
                </thead>
                <tbody className='bg-white'>
                  {(!!filteredAPIKeys ? filteredAPIKeys : apiKeyList)
                    .sort((a, b) =>
                      a.metadata?.name === b.metadata?.name
                        ? 0
                        : a.metadata!.name > b.metadata!.name
                        ? 1
                        : -1
                    )
                    .map(apiKey => (
                      <tr key={apiKey.metadata?.uid}>
                        <td className='max-w-xs px-6 py-4 border-b border-gray-200'>
                          <div className='text-sm leading-5 text-gray-700'>
                            <span className='flex items-center '>
                              {apiKey.metadata?.name}
                            </span>
                          </div>
                        </td>
                        <td className='max-w-xs px-6 py-4 border-b border-gray-200'>
                          <div className='text-sm leading-5 text-gray-700'>
                            <span className='flex items-center truncate '>
                              {apiKey.value}
                            </span>
                          </div>
                        </td>
                        <td className='max-w-xs px-6 py-4 border-b border-gray-200'>
                          <div className='text-sm leading-5 text-gray-700'>
                            <span className='flex items-center '>
                              {apiKey.user?.name}
                            </span>
                          </div>
                        </td>
                        <td className='max-w-xs px-6 py-4 border-b border-gray-200'>
                          <div className='text-sm leading-5 text-gray-700'>
                            <span className='flex items-center '>
                              {apiKey.keyScope?.name}
                            </span>
                          </div>
                        </td>
                        <td className='max-w-xs px-6 py-4 border-b border-gray-200'>
                          <div className='text-sm leading-5 text-gray-700'>
                            <span className='flex items-center '>
                              {apiKey.metadata?.labelsMap}
                            </span>
                          </div>
                        </td>
                        <td className='max-w-xs px-6 py-4 text-sm font-medium leading-5 text-right border-b border-gray-200'>
                          <span className='flex items-center'>
                            {/* <div className='flex items-center justify-center w-4 h-4 mr-3 text-gray-700 bg-gray-400 rounded-full cursor-pointer'>
                                <EditIcon className='w-2 h-3 fill-current' />
                              </div> */}

                            <div
                              className='flex items-center justify-center w-4 h-4 text-gray-700 bg-gray-400 rounded-full cursor-pointer'
                              onClick={() => attemptDeleteApiKey(apiKey)}>
                              x
                            </div>
                          </span>
                        </td>
                      </tr>
                    ))}
                </tbody>
              </table>
              {(!!filteredAPIKeys ? filteredAPIKeys : apiKeyList).length ===
                0 && (
                <div className='w-full m-auto'>
                  <div className='flex flex-col items-center justify-center w-full h-full py-4 mr-32 bg-white rounded-lg shadow-lg md:flex-row'>
                    <div className='mr-6'>
                      <NoApiKey />
                    </div>
                    <div className='flex flex-col h-full'>
                      <p className='h-auto my-6 text-lg font-medium text-gray-800 '>
                        There are no API keys to display!{' '}
                      </p>
                    </div>
                  </div>
                </div>
              )}
              <ConfirmationModal
                visible={attemptingDelete}
                confirmationTopic='delete this API key'
                confirmText='Delete'
                goForIt={deleteApiKey}
                cancel={cancelDeletion}
                isNegative={true}
              />
            </div>
          </div>
        </div>
      </div>
    </SectionCard>
  );
};
