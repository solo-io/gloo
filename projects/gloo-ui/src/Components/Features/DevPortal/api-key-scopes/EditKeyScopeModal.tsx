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
import { KeyScopeStatus } from 'proto/dev-portal/api/dev-portal/v1/portal_pb';
import { Formik } from 'formik';
import { ObjectRef } from 'proto/dev-portal/api/dev-portal/v1/common_pb';
import useSWR from 'swr';
import { DevPortalApi } from '../api';
import {
  SoloFormInput,
  SoloFormDropdown,
  SoloFormTextarea
} from 'Components/Common/Form/SoloFormField';
import { SoloTransfer } from 'Components/Common/SoloTransfer';
import {
  SoloButtonStyledComponent,
  SoloCancelButton
} from 'Styles/CommonEmotions/button';
import { Portal } from 'proto/dev-portal/api/grpc/admin/portal_pb';

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

interface InitialKeyScopeValuesType {
  name: string;
  assignedPortal: string;
  description: string;
  chosenAPIs: ObjectRef.AsObject[];
}

interface EditKeyScopeModalProps {
  createNotEdit?: boolean;
  onEdit: (newScopeData: KeyScopeStatus.AsObject) => any;
  onCancel: () => any;
}
export function EditKeyScopeModal(props: EditKeyScopeModalProps) {
  const { data: portalsList, error: getApiKeyDocsError } = useSWR(
    'listPortals',
    DevPortalApi.listPortals
  );

  const [tabIndex, setTabIndex] = React.useState(0);
  const [portalSelected, setPortalSelected] = React.useState<Portal.AsObject>();

  const textTopic = !!props.createNotEdit ? 'Create' : 'Edit';

  if (!portalsList) {
    return <React.Fragment />;
  }

  const handleTabsChange = (index: number, portalChosenUid?: string) => {
    if (!!portalSelected && portalSelected?.metadata?.uid !== portalChosenUid) {
      setPortalSelected(
        portalsList.find(portal => portal.metadata?.uid === portalChosenUid)!
      );
    }
    setTabIndex(index);
  };

  const initialValues: InitialKeyScopeValuesType = {
    name: '',
    assignedPortal: !!portalsList?.length
      ? portalsList[0]?.metadata?.uid || ''
      : '',
    description: '',
    chosenAPIs: []
  };

  const updateAPIsList = (newPortalUid: string) => {};
  const attemptCreate = (values: InitialKeyScopeValuesType) => {};

  const { onCancel, onEdit } = props;
  return (
    <div className='bg-white rounded-lg shadow '>
      <Formik<InitialKeyScopeValuesType>
        initialValues={initialValues}
        onSubmit={attemptCreate}>
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
              <div className='pt-5 pr-6 pb-6 pl-8'>
                <TabPanel className='focus:outline-none'>
                  <div className='relative flex flex-col'>
                    <div className='text-lg flex items-center font-medium'>
                      {textTopic} an API Key Scope{' '}
                      <StacksIcon className='w-8 h-8 ml-2' />
                    </div>

                    <div className='rounded-lg bg-gray-100 p-3 mt-3 text-gray-700'>
                      Lorem Ipsum
                    </div>

                    <div className='mt-4 grid grid-cols-2 gap-4'>
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
                          testId={`edit-key-scope-portal-id`}
                          name='assignedPortal'
                          placeholder='Select'
                          title='Assigned Portal'
                          hideError={true}
                          options={portalsList.map(portal => {
                            return {
                              value: portal.metadata?.uid,
                              key: portal.metadata?.uid,
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
                <TabPanel className='focus:outline-none'>
                  <div className='relative flex flex-col '>
                    <div className='text-lg flex items-center font-medium'>
                      {textTopic} an API Key Scope: APIs{' '}
                      <StacksIcon className='w-8 h-8  ml-2' />
                    </div>

                    <div className='rounded-lg bg-gray-100 p-3 mt-3 text-gray-700'>
                      Lorem Ipsum
                    </div>

                    <div className='mt-4'>
                      {!!portalSelected?.status?.apiDocsList ? (
                        <SoloTransfer
                          allOptionsListName='Available APIs'
                          allOptions={portalSelected?.status?.apiDocsList.map(
                            apiDoc => {
                              return {
                                value: apiDoc
                              };
                            }
                          )}
                          chosenOptionsListName='Selected APIs'
                          chosenOptions={values.chosenAPIs.map(api => {
                            return { value: api.name + api.namespace };
                          })}
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

                    <div className='mt-8 flex items-center justify-between'>
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
                          Create Scope
                        </SoloButtonStyledComponent>
                      </div>
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
