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
  chosenGroups: ObjectRef.AsObject[];
};

export const AddGroupToAPIModal = (props: AddUserProps) => {
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

  const { data: portalsList, error: portalsError } = useSWR(`listPortals`, () =>
    portalApi.listPortals()
  );

  const { data: allGroupsList, error: allGroupsError } = useSWR(
    'listGroups',
    groupApi.listGroups
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
      groupsList: values.chosenGroups,
      portalsList:
        filteredPortalList?.map(portal => {
          return {
            name: portal.metadata?.name!,
            namespace: portal.metadata?.namespace!
          };
        }) || [],
      usersList:
        usersList?.map(user => {
          return {
            name: user.metadata?.name!,
            namespace: user.metadata?.namespace!
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

  if (!groupsList || !allGroupsList) {
    return <Loading center>Loading...</Loading>;
  }
  return (
    <Formik<AddUserValues>
      initialValues={{
        chosenGroups:
          groupsList.map(group => {
            return {
              name: group.metadata?.name!,
              namespace: group.metadata?.namespace!,
              displayValue: group.spec?.displayName
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
            <SectionHeader>Add Groups</SectionHeader>
            <div className='p-3 mb-2 text-gray-700 bg-gray-100 rounded-lg'>
              Select the groups to which you'd like to grant access to this API
            </div>

            <SoloTransfer
              allOptionsListName='Available Groups'
              allOptions={allGroupsList
                .sort((a, b) =>
                  a.metadata?.name === b.metadata?.name
                    ? 0
                    : a.metadata!.name > b.metadata!.name
                    ? 1
                    : -1
                )
                .map(group => {
                  return {
                    name: group.metadata?.name!,
                    namespace: group.metadata?.namespace!,
                    displayValue: group.spec?.displayName
                  };
                })}
              chosenOptionsListName='Selected Groups'
              chosenOptions={values.chosenGroups}
              onChange={newChosenOptions => {
                setFieldValue('chosenGroups', newChosenOptions);
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
                Add Groups
              </SoloButtonStyledComponent>
            </div>
          </div>
        </div>
      )}
    </Formik>
  );
};
