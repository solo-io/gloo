import { ObjectRef } from '@solo-io/dev-portal-grpc/dev-portal/api/dev-portal/v1/common_pb';
import { Loading } from 'Components/Common/DisplayOnly/Loading';
import { ListItemType, SoloTransfer } from 'Components/Common/SoloTransfer';
import { Formik } from 'formik';
import React from 'react';
import { useParams } from 'react-router';
import { SoloButtonStyledComponent } from 'Styles/CommonEmotions/button';
import useSWR from 'swr';
import { apiDocApi, groupApi, portalApi, userApi } from '../api';
import { SectionContainer, SectionHeader } from '../apis/CreateAPIModal';

interface AddUserProps {
  onClose: () => void;
}

type AddUserValues = {
  chosenUsers: ObjectRef.AsObject[];
};

export const AddUserToAPIModal = (props: AddUserProps) => {
  const { apiname, apinamespace } = useParams();
  const { data: apiDoc, error: apiDocError } = useSWR(
    !!apiname && !!apinamespace ? ['getApiDoc', apiname, apinamespace] : null,
    (key, name, namespace) =>
      apiDocApi.getApiDoc({ apidoc: { name, namespace }, withassets: true })
  );

  const { data: usersList, error: usersError } = useSWR(
    `listUsers${apiname!}${apinamespace!}`,
    () =>
      userApi.listUsers({
        portalsList: [],
        apiDocsList: [{ name: apiname!, namespace: apinamespace! }],
        groupsList: []
      })
  );
  const { data: groupsList, error: groupsError } = useSWR(
    `listGroups${apiname}${apinamespace}`,
    () =>
      groupApi.listGroups({
        portalsList: [],
        apiDocsList: [{ name: apiname!, namespace: apinamespace! }]
      })
  );

  const {
    data: portalsList,
    error: portalsError
  } = useSWR(`listApiDocs${apiname}${apinamespace}`, () =>
    portalApi.listPortals()
  );

  const { data: allUsersList, error: allUsersError } = useSWR(
    'listUsers',
    userApi.listUsers
  );

  const filteredPortalList = portalsList?.filter(portal =>
    portal.status?.apiDocsList.some(
      apiDocRef =>
        apiDocRef.name === apiDoc?.metadata?.name &&
        apiDocRef.namespace === apiDoc?.metadata.namespace
    )
  );
  const [errorMessage, setErrorMessage] = React.useState('');

  const addApi = async (values: AddUserValues) => {
    //@ts-ignore
    await apiDocApi.updateApiDoc({
      usersList: values.chosenUsers,
      portalsList:
        filteredPortalList?.map(portal => {
          return {
            name: portal.metadata?.name!,
            namespace: portal.metadata?.namespace!
          };
        }) || [],
      groupsList:
        groupsList?.map(group => {
          return {
            name: group.metadata?.name!,
            namespace: group.metadata?.namespace!
          };
        }) || [],
      //@ts-ignore
      apidoc: {
        //@ts-ignore
        metadata: {
          name: apiDoc?.metadata?.name!,
          namespace: apiDoc?.metadata?.namespace!
        }
      },
      apiDocOnly: false
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
              Select the users to which you'd like to grant access to this API
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
