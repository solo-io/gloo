import React from 'react';
import {
  Tabs,
  TabList,
  Tab,
  TabPanels,
  TabPanel,
  TabPanelProps
} from '@reach/tabs';
import { css } from '@emotion/core';
import { ReactComponent as StacksIcon } from 'assets/app-icon.svg';
import {
  KeyScopeStatus,
  KeyScope
} from '@solo-io/dev-portal-grpc/dev-portal/api/dev-portal/v1/portal_pb';
import { Formik } from 'formik';
import { ObjectRef } from '@solo-io/dev-portal-grpc/dev-portal/api/dev-portal/v1/common_pb';
import useSWR from 'swr';
import { portalApi, apiDocApi, apiKeyScopeApi } from '../api';
import {
  SoloFormInput,
  SoloFormDropdown,
  SoloFormTextarea
} from 'Components/Common/Form/SoloFormField';
import { SoloTransfer, ListItemType } from 'Components/Common/SoloTransfer';
import {
  SoloButtonStyledComponent,
  SoloCancelButton
} from 'Styles/CommonEmotions/button';
import { Portal } from '@solo-io/dev-portal-grpc/dev-portal/api/grpc/admin/portal_pb';
import { ApiKeyScope } from '@solo-io/dev-portal-grpc/dev-portal/api/grpc/admin/api_key_scope_pb';
import * as yup from 'yup';

const StyledTab = (
  props: {
    disabled?: boolean | undefined;
  } & TabPanelProps & {
      isSelected?: boolean | undefined;
    }
) => {
  const { isSelected, children } = props;
  return (
    <Tab
      {...props}
      className={`p-1 text-left w-48 text-white  pl-6 mb-2 focus:outline-none ${
        isSelected
          ? ' bg-blue-500 border-r-8 border-blue-300  '
          : 'bg-blue-600 '
      }`}>
      {children}
    </Tab>
  );
};

const validationSchema = yup.object().shape({
  name: yup.string().required('Display name is required.'),
  assignedPortal: yup.string().required('Assigned portal is required.')
});

interface InitialKeyScopeValuesType {
  name: string;
  assignedPortal: string;
  description: string;
  chosenAPIs: ListItemType[];
}

interface EditKeyScopeModalProps {
  existingKeyScope?: ApiKeyScope.AsObject;
  createNotEdit?: boolean;
  onEdit: (newScopeData: KeyScopeStatus.AsObject) => any;
  onCancel: () => any;
}
export function EditKeyScopeModal(props: EditKeyScopeModalProps) {
  const { data: portalsList, error: getApiKeyDocsError } = useSWR(
    'listPortals',
    portalApi.listPortals
  );
  const { data: apiDocsList, error: apiDocsError } = useSWR(
    'listApiDocs',
    apiDocApi.listApiDocs
  );

  const [tabIndex, setTabIndex] = React.useState(0);
  const [portalSelected, setPortalSelected] = React.useState<Portal.AsObject>();

  const textTopic = !!props.createNotEdit ? 'Create' : 'Edit';

  if (!portalsList) {
    return <React.Fragment />;
  }

  const handleTabsChange = (index: number, portalNamespaceName?: string) => {
    setPortalSelected(
      portalsList.find(
        portal =>
          getPortalId(portal.metadata!.namespace, portal.metadata!.name) ===
          portalNamespaceName
      )!
    );

    setTabIndex(index);
  };

  const getPortalId = (namespace: string, name: string) =>
    `${namespace}.${name}`;

  const getApiDocFromRef = (apiDocRef: ObjectRef.AsObject) => {
    let apiDocObj = apiDocsList?.find(
      apiDoc =>
        apiDoc.metadata?.name === apiDocRef.name &&
        apiDoc.metadata.namespace === apiDocRef.namespace
    );
    return apiDocObj;
  };
  const initialValues: InitialKeyScopeValuesType = {
    name:
      props.existingKeyScope?.spec?.displayName ||
      props.existingKeyScope?.spec?.name ||
      '',
    assignedPortal: props.existingKeyScope?.portal
      ? getPortalId(
          props.existingKeyScope.portal.namespace,
          props.existingKeyScope.portal.name
        )
      : '',
    description: props.existingKeyScope?.spec?.description || '',
    chosenAPIs: props.existingKeyScope?.status?.accessibleApiDocsList || []
  };

  const updateAPIsList = (newPortalUid: string) => {};
  const attemptWrite = (values: InitialKeyScopeValuesType) => {
    if (props.createNotEdit) {
      apiKeyScopeApi.createKeyScope({
        apiKeyScope: {
          portal: {
            namespace: values.assignedPortal.split('.')[0],
            name: values.assignedPortal.split('.')[1]
          },
          spec: {
            displayName: values.name,
            description: values.description,
            name: values.name,
            namespace: values.assignedPortal.split('.')[0]
          }
        },
        apiDocsList: values.chosenAPIs,
        apiKeyScopeOnly: false
      });
    } else {
      apiKeyScopeApi.updateKeyScope({
        apiKeyScope: {
          portal: {
            namespace: values.assignedPortal.split('.')[0],
            name: values.assignedPortal.split('.')[1]
          },
          spec: {
            displayName: values.name,
            description: values.description,
            name: props.existingKeyScope!.spec!.name,
            namespace: props.existingKeyScope!.spec!.namespace
          }
        },
        apiDocsList: values.chosenAPIs,
        apiKeyScopeOnly: false
      });
    }

    props.onCancel();
  };

  const { onCancel, onEdit } = props;
  return (
    <div className='bg-white rounded-lg shadow '>
      <Formik<InitialKeyScopeValuesType>
        initialValues={initialValues}
        validationSchema={validationSchema}
        onSubmit={attemptWrite}>
        {({ setFieldValue, handleSubmit, values }) => (
          <Tabs
            className='bg-blue-600 rounded-lg'
            style={{ minHeight: '400px' }}
            index={tabIndex}
            onChange={newIndex =>
              handleTabsChange(newIndex, values.assignedPortal)
            }
            css={css`
              display: grid;
              grid-template-columns: 190px 1fr;
            `}>
            <TabList className='flex flex-col mt-6'>
              <StyledTab>General</StyledTab>
              <StyledTab>APIs</StyledTab>
            </TabList>

            <TabPanels className='bg-white rounded-r-lg'>
              <div className='h-full pt-5 pb-6 pl-8 pr-6'>
                <TabPanel className='focus:outline-none'>
                  <div className='relative flex flex-col'>
                    <div className='flex items-center text-lg font-medium text-gray-900'>
                      {textTopic} an API Key Scope{' '}
                      <StacksIcon className='w-8 h-8 ml-2 ' />
                    </div>

                    <div className='p-3 mt-3 text-gray-700 bg-gray-100 rounded-lg'>
                      Create a new rule for access keys to APIs for a particular
                      portal
                    </div>

                    <div className='grid grid-cols-2 gap-4 mt-4'>
                      <div className='mb-4'>
                        <SoloFormInput
                          testId={`edit-key-scope-name`}
                          name='name'
                          placeholder='App Name'
                          title='API Key Scope Name'
                          hideError={true}
                        />
                      </div>
                      <div className='mb-4'>
                        <SoloFormDropdown
                          disabled={!props.createNotEdit}
                          testId={`edit-key-scope-portal-id`}
                          name='assignedPortal'
                          placeholder='Select'
                          title='Assigned Portal'
                          hideError={true}
                          options={portalsList
                            .sort((a, b) =>
                              a.metadata?.name === b.metadata?.name
                                ? 0
                                : a.metadata!.name > b.metadata!.name
                                ? 1
                                : -1
                            )
                            .map(portal => {
                              return {
                                value: getPortalId(
                                  portal.metadata!.namespace,
                                  portal.metadata!.name
                                ),
                                key: getPortalId(
                                  portal.metadata!.namespace,
                                  portal.metadata!.name
                                ),
                                displayValue: portal.spec?.displayName
                              };
                            })}
                        />
                      </div>
                    </div>
                    <div>
                      <SoloFormTextarea
                        testId={`edit-key-scope-description`}
                        name='description'
                        placeholder='This is a description of the test scope'
                        title='Description'
                        hideError={true}
                      />
                    </div>

                    <div className='flex items-center justify-between'>
                      <div
                        onClick={onCancel}
                        className='text-blue-500 cursor-pointer'>
                        Cancel
                      </div>

                      <SoloButtonStyledComponent
                        onClick={() =>
                          handleTabsChange(1, values.assignedPortal)
                        }>
                        Next Step
                      </SoloButtonStyledComponent>
                    </div>
                  </div>
                </TabPanel>
                <TabPanel className='relative flex flex-col justify-between h-full focus:outline-none'>
                  <div className='w-full h-full'>
                    <div className='flex items-center text-lg font-medium'>
                      {textTopic} an API Key Scope: APIs{' '}
                      <StacksIcon className='w-8 h-8 ml-2' />
                    </div>

                    <div className='p-3 mt-3 text-gray-700 bg-gray-100 rounded-lg'>
                      Choose selected APIs to which to apply this key scope
                    </div>

                    <div className='mt-4'>
                      {!!apiDocsList ? (
                        <SoloTransfer
                          allOptionsListName='Available APIs'
                          allOptions={
                            portalSelected?.status?.apiDocsList.map(
                              apiDocRef => {
                                let apiDocObj = getApiDocFromRef(apiDocRef)!;
                                return {
                                  name: apiDocObj?.metadata?.name!,
                                  namespace: apiDocObj?.metadata?.namespace!,
                                  displayValue: apiDocObj?.status?.displayName
                                };
                              }
                            ) || []
                          }
                          chosenOptionsListName='Selected APIs'
                          chosenOptions={values.chosenAPIs}
                          onChange={newChosenOptions => {
                            setFieldValue('chosenAPIs', newChosenOptions);
                          }}
                        />
                      ) : (
                        <div className='font-normal'>
                          No APIs are available on this portal to assign.
                        </div>
                      )}
                    </div>
                  </div>

                  <div className='flex items-center justify-between mt-16'>
                    <div
                      onClick={onCancel}
                      className='text-blue-500 cursor-pointer'>
                      Cancel
                    </div>

                    <div>
                      <SoloCancelButton
                        onClick={() => handleTabsChange(0)}
                        className='mr-2'>
                        Back
                      </SoloCancelButton>
                      <SoloButtonStyledComponent onClick={handleSubmit}>
                        {props.createNotEdit
                          ? 'Create Scope'
                          : 'Publish Changes'}
                      </SoloButtonStyledComponent>
                    </div>
                  </div>
                </TabPanel>
              </div>
            </TabPanels>
          </Tabs>
        )}
      </Formik>
    </div>
  );
}
