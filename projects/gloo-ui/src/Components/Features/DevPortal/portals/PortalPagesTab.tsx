import React from 'react';
import { ReactComponent as EditPencilIcon } from 'assets/edit-pencil.svg';
import { ReactComponent as DeleteXIcon } from 'assets/small-grey-x.svg';
import { ReactComponent as Plus } from 'assets/small-green-plus.svg';
import { ReactComponent as PlaceholderPortal } from 'assets/placeholder-portal.svg';
import { ReactComponent as PortalPageIcon } from 'assets/portal-page-icon.svg';
import { SoloInput } from 'Components/Common/SoloInput';
import useSWR, { trigger } from 'swr';
import { useParams, useHistory, useLocation } from 'react-router';
import { SoloModal } from 'Components/Common/SoloModal';
import { CreatePageModal } from './CreatePageModal';
import { ConfirmationModal } from 'Components/Common/ConfirmationModal';
import { portalApi } from '../api';

export const PortalPagesTab = () => {
  const routerLocation = useLocation();
  const routerHistory = useHistory();
  const { portalname, portalnamespace } = useParams();

  const { data: portal, error: portalListError } = useSWR(
    !!portalname && !!portalnamespace
      ? ['getPortal', portalname, portalnamespace]
      : null,
    (key, name, namespace) => portalApi.getPortalWithAssets({ name, namespace })
  );

  const [pagesSearchTerm, setPagesSearchTerm] = React.useState('');
  const [createPageModalOpen, setCreatePageModalOpen] = React.useState(false);
  const [pageAttemptingToDelete, setPageAttemptingToDelete] = React.useState<
    string
  >();

  const openCreatePage = () => {
    setCreatePageModalOpen(true);
  };
  const closeCreatePage = () => {
    setCreatePageModalOpen(false);
  };

  const attemptDeletion = (name: string) => {
    setPageAttemptingToDelete(name);
  };
  const cancelDeletion = () => {
    setPageAttemptingToDelete(undefined);
  };

  const finishDeletion = () => {
    portalApi
      .deletePortalPage(
        { name: portalname!, namespace: portalnamespace! },
        pageAttemptingToDelete!
      )
      .then(portal => {
        trigger(['getPortal', portalname, portalnamespace]);

        setPageAttemptingToDelete(undefined);
      });
    /*.catch(err => {
        setErrorMessage(err);

        setTimeout(() => {
          setErrorMessage('');
        }, 10000);
      });*/
  };

  const filteredList = portal?.spec?.staticPagesList.filter(page =>
    page.name.includes(pagesSearchTerm)
  );

  return (
    <div className='relative flex flex-col p-4 border border-gray-300 rounded-lg'>
      <span
        className='absolute flex font-normal cursor-pointer top-4 right-4'
        onClick={openCreatePage}>
        <span className='text-green-400 hover:text-green-300'>
          <Plus className='w-5 h-5 mr-2 fill-current' />
        </span>
        Add a Page
      </span>
      {!!portal?.spec?.staticPagesList.length && (
        <div className='w-1/3 m-4'>
          <SoloInput
            placeholder='Search by page name...'
            value={pagesSearchTerm}
            onChange={e => setPagesSearchTerm(e.target.value)}
          />
        </div>
      )}
      {!!filteredList?.length ? (
        <div className='flex flex-col'>
          <div className='py-2 -my-2 overflow-x-auto sm:-mx-6 sm:px-6 lg:-mx-8 lg:px-8'>
            <div className='inline-block min-w-full overflow-hidden align-middle border-b border-gray-200 shadow sm:rounded-lg'>
              <table className='min-w-full'>
                <thead className='bg-gray-300 '>
                  <tr>
                    <th className='px-6 py-3 text-sm font-medium leading-4 tracking-wider text-left text-gray-800 capitalize border-b border-gray-200 bg-gray-50'>
                      Page Name
                    </th>
                    <th className='px-6 py-3 text-sm font-medium leading-4 tracking-wider text-left text-gray-800 capitalize border-b border-gray-200 bg-gray-50'>
                      Navigation Link
                    </th>
                    <th className='px-6 py-3 text-sm font-medium leading-4 tracking-wider text-left text-gray-800 capitalize border-b border-gray-200 bg-gray-50'>
                      Page Url Path
                    </th>
                    <th className='px-6 py-3 text-sm font-medium leading-4 tracking-wider text-left text-gray-800 capitalize border-b border-gray-200 bg-gray-50'>
                      Homepage?
                    </th>

                    <th className='px-6 py-3 text-sm font-medium leading-4 tracking-wider text-left text-gray-800 capitalize border-b border-gray-200 bg-gray-50'>
                      Actions
                    </th>
                  </tr>
                </thead>
                <tbody className='bg-white'>
                  {filteredList!
                    .sort((a, b) =>
                      a.name === b.name ? 0 : a.name > b.name ? 1 : -1
                    )
                    .map(page => (
                      <tr key={page.name}>
                        <td className='px-6 py-4 whitespace-no-wrap border-b border-gray-200'>
                          <div className='text-sm leading-5 text-gray-900'>
                            <span className='flex items-center capitalize'>
                              {page.name}
                            </span>
                          </div>
                        </td>
                        <td className='px-6 py-4 whitespace-no-wrap border-b border-gray-200'>
                          <div className='text-sm leading-5 text-gray-900'>
                            <span className='flex items-center capitalize'>
                              {page.navigationLinkName}
                            </span>
                          </div>
                        </td>
                        <td className='px-6 py-4 whitespace-no-wrap border-b border-gray-200'>
                          <div className='text-sm leading-5 text-gray-900'>
                            <span className='flex items-center capitalize'>
                              {page.path}
                            </span>
                          </div>
                        </td>
                        <td className='px-6 py-4 whitespace-no-wrap border-b border-gray-200'>
                          <div className='text-sm leading-5 text-gray-900'>
                            <span className='flex items-center capitalize'>
                              {page.displayOnHomepage ? 'yes' : ''}
                            </span>
                          </div>
                        </td>
                        <td className='px-6 py-4 text-sm font-medium leading-5 text-right whitespace-no-wrap border-b border-gray-200'>
                          <span className='flex items-center'>
                            <div className='flex items-center justify-center w-4 h-4 mr-3 bg-gray-400 rounded-full cursor-pointer'>
                              <EditPencilIcon
                                className='w-2 h-3 fill-current'
                                onClick={() =>
                                  routerHistory.push({
                                    pathname: `${routerLocation.pathname}/page-editor/${page.name}`
                                  })
                                }
                              />
                            </div>

                            <div className='flex items-center justify-center w-4 h-4 bg-gray-400 rounded-full cursor-pointer'>
                              <div
                                className='flex items-center justify-center w-4 h-4 text-gray-700 bg-gray-400 rounded-full cursor-pointer'
                                onClick={() => attemptDeletion(page.name)}>
                                x
                              </div>
                            </div>
                            {/* )} */}
                          </span>
                        </td>
                      </tr>
                    ))}
                </tbody>
              </table>
            </div>
          </div>
        </div>
      ) : (
        <div className=''>
          <div className='flex items-center w-full h-full bg-white rounded-lg md:flex-row'>
            <div className='mr-6 text-gray-500'>
              <PortalPageIcon className='w-32 rounded-lg fill-current' />
            </div>
            <div className='flex flex-col h-full'>
              <p className='h-auto mb-6 text-lg font-medium'>
                {!!portal?.spec?.staticPagesList.length
                  ? 'No pages match this search'
                  : `${portalname} has no Pages currently`}
                .
              </p>
              <p className='text-base font-normal text-gray-700 '>
                {!!portal?.spec?.staticPagesList.length ? (
                  <>
                    {' '}
                    Want to{' '}
                    <span
                      className='text-blue-600 cursor-pointer'
                      onClick={openCreatePage}>
                      add it
                    </span>
                    ?
                  </>
                ) : (
                  <>
                    Get started by{' '}
                    <span
                      className='text-blue-600 cursor-pointer'
                      onClick={openCreatePage}>
                      Adding a Page
                    </span>
                    .
                  </>
                )}
              </p>
            </div>
          </div>
        </div>
      )}

      <SoloModal
        visible={createPageModalOpen}
        width={625}
        onClose={closeCreatePage}>
        <CreatePageModal onClose={closeCreatePage} />
      </SoloModal>
      <ConfirmationModal
        visible={!!pageAttemptingToDelete}
        confirmationTopic={`delete ${pageAttemptingToDelete}`}
        confirmText='Delete'
        goForIt={finishDeletion}
        cancel={cancelDeletion}
        isNegative={true}
      />
    </div>
  );
};
