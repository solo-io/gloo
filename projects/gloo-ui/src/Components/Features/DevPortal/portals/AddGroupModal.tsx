import { ReactComponent as PortalPageIcon } from 'assets/portal-page-icon.svg';
import { Loading } from 'Components/Common/DisplayOnly/Loading';
import { SoloTransfer, ListItemType } from 'Components/Common/SoloTransfer';
import { Formik } from 'formik';
import { ObjectRef } from '@solo-io/dev-portal-grpc/dev-portal/api/dev-portal/v1/common_pb';
import React from 'react';
import { useParams } from 'react-router';
import { SoloButtonStyledComponent } from 'Styles/CommonEmotions/button';
import useSWR from 'swr';
import { apiDocApi, portalApi, userApi, groupApi } from '../api';
import { SectionContainer, SectionHeader } from '../apis/CreateAPIModal';

interface AddGroupProps {
  onClose: () => void;
}

type AddUserValues = {
  chosenGroups: ObjectRef.AsObject[];
};

export const AddGroupModal = (props: AddGroupProps) => {
  const { portalname, portalnamespace } = useParams();
  const { data: portal, error: portalListError } = useSWR(
    !!portalname && !!portalnamespace
      ? ['getPortal', portalname, portalnamespace]
      : null,
    (key, name, namespace) => portalApi.getPortalWithAssets({ name, namespace })
  );

  const { data: groupsList, error: groupsError } = useSWR(
    `listGroups${portalname!}${portalnamespace!}`,
    () =>
      groupApi.listGroups({
        portalsList: [{ name: portalname!, namespace: portalnamespace! }],
        apiDocsList: []
      })
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
  const { data: allGroupsList, error: allGroupsError } = useSWR(
    'listGroups',
    groupApi.listGroups
  );

  const [errorMessage, setErrorMessage] = React.useState('');

  const addApi = async (values: AddUserValues) => {
    //@ts-ignore
    await portalApi.updatePortal({
      groupsList: values.chosenGroups,
      usersList:
        usersList?.map(user => {
          return {
            name: user.metadata?.name!,
            namespace: user.metadata?.namespace!
          };
        }) || [],
      apiDocsList:
        apiDocsList?.map(apiDOc => {
          return {
            name: apiDOc.metadata?.name!,
            namespace: apiDOc.metadata?.namespace!
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
              Select the groups to which you'd like to grant access to this
              portal
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
                .map(user => {
                  return {
                    name: user.metadata?.name!,
                    namespace: user.metadata?.namespace!,
                    displayValue: user.spec?.displayName
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
