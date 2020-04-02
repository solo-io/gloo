import { ReactComponent as PortalPageIcon } from 'assets/portal-page-icon.svg';
import { Loading } from 'Components/Common/DisplayOnly/Loading';
import { SoloTransfer, ListItemType } from 'Components/Common/SoloTransfer';
import { Formik } from 'formik';
import { ObjectRef } from 'proto/dev-portal/api/dev-portal/v1/common_pb';
import React from 'react';
import { useParams } from 'react-router';
import { SoloButtonStyledComponent } from 'Styles/CommonEmotions/button';
import useSWR from 'swr';
import { apiDocApi, portalApi, userApi, groupApi } from '../api';
import { SectionContainer, SectionHeader } from '../apis/CreateAPIModal';

interface AddUserProps {
  onClose: () => void;
}

type AddUserValues = {
  chosenUsers: ObjectRef.AsObject[];
};

export const AddUserModal = (props: AddUserProps) => {
  const { portalname, portalnamespace } = useParams();
  const { data: portal, error: portalListError } = useSWR(
    !!portalname && !!portalnamespace
      ? ['getPortal', portalname, portalnamespace]
      : null,
    (key, name, namespace) => portalApi.getPortalWithAssets({ name, namespace })
  );

  const { data: usersList, error: usersError } = useSWR(
    `listUsers${portalname!}${portalnamespace!}`,
    () =>
      userApi.listUsers({
        portalsList: [{ name: portalname!, namespace: portalnamespace! }],
        apiDocsList: [],
        groupsList: []
      })
  );
  const { data: groupsList, error: groupsError } = useSWR(
    `listGroups${portalname}${portalnamespace}`,
    () =>
      groupApi.listGroups({
        portalsList: [
          {
            name: portalname!,
            namespace: portalnamespace!
          }
        ],
        apiDocsList: []
      })
  );
  const { data: apiDocsList, error: apiDocsError } = useSWR(
    `listApiDocs${portalname}${portalnamespace}`,
    () =>
      apiDocApi.listApiDocs({
        portalsList: [
          {
            name: portalname!,
            namespace: portalnamespace!
          }
        ]
      })
  );

  const { data: allUsersList, error: allUsersError } = useSWR(
    'listUsers',
    userApi.listUsers
  );

  const [errorMessage, setErrorMessage] = React.useState('');

  const addApi = async (values: AddUserValues) => {
    //@ts-ignore
    await portalApi.updatePortal({
      usersList: values.chosenUsers,
      apiDocsList:
        apiDocsList?.map(apiDoc => {
          return {
            name: apiDoc.metadata?.name!,
            namespace: apiDoc.metadata?.namespace!
          };
        }) || [],
      groupsList:
        groupsList?.map(group => {
          return {
            name: group.metadata?.name!,
            namespace: group.metadata?.namespace!
          };
        }) || [],

      portal: {
        //@ts-ignore
        metadata: {
          name: portal?.metadata?.name!,
          namespace: portal?.metadata?.namespace!
        }
      },
      portalOnly: false
    });

    props.onClose();
  };

  if (!usersList || !allUsersList) {
    return <Loading center>Loading...</Loading>;
  }
  return (
    <Formik<AddUserValues>
      initialValues={{
        chosenUsers:
          usersList.map(user => {
            return {
              name: user.metadata?.name!,
              namespace: user.metadata?.namespace!,
              displayValue: user.spec?.username
            };
          }) || ([] as ListItemType[])
      }}
      onSubmit={addApi}>
      {({ handleSubmit, setFieldValue, values }) => (
        <div className='flex flex-col h-full pt-4 '>
          {!!errorMessage.length && (
            <div className='p-4 text-orange-600 bg-orange-200'>
              {errorMessage}
            </div>
          )}

          <SectionContainer className='mb-8'>
            <SectionHeader>Add Users</SectionHeader>
            <div className='p-3 mb-2 text-gray-700 bg-gray-100 rounded-lg'>
              Select the users to which you'd like to grant access to this
              portal
            </div>

            <SoloTransfer
              allOptionsListName='Available Users'
              allOptions={allUsersList
                .sort((a, b) =>
                  a.metadata?.name === b.metadata?.name
                    ? 0
                    : a.metadata!.name > b.metadata!.name
                    ? 1
                    : -1
                )
                .map(user => {
                  return {
                    name: user.metadata?.name!,
                    namespace: user.metadata?.namespace!,
                    displayValue: user.spec?.username
                  };
                })}
              chosenOptionsListName='Selected Users'
              chosenOptions={values.chosenUsers}
              onChange={newChosenOptions => {
                setFieldValue('chosenUsers', newChosenOptions);
              }}
            />
          </SectionContainer>
          <div className='flex items-end justify-between h-full px-6 mb-4 '>
            <button
              className='text-blue-500 cursor-pointer'
              onClick={props.onClose}>
              cancel
            </button>
            <div>
              <SoloButtonStyledComponent onClick={handleSubmit}>
                Add Users
              </SoloButtonStyledComponent>
            </div>
          </div>
        </div>
      )}
    </Formik>
  );
};
