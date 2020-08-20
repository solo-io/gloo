import styled from '@emotion/styled';
import {
  SoloFormInput,
  SoloFormStringsList,
  SoloFormTypeahead,
  SoloFormSelect
} from 'Components/Common/Form/SoloFormField';
import { SoloButton } from 'Components/Common/SoloButton';
import { Formik } from 'formik';
import * as React from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { useHistory, useLocation } from 'react-router';
import { configAPI } from 'store/config/api';
import { createVirtualService } from 'store/virtualServices/actions';
import useSWR, { mutate } from 'swr';
import * as yup from 'yup';
import { virtualServiceAPI } from 'store/virtualServices/api';
import { AppState } from 'store';
import { guardByLicense } from 'store/config/actions';
import { Select } from 'antd';

const Footer = styled.div`
  display: flex;
  flex-direction: row;
  justify-content: flex-end;
`;

const InputContainer = styled.div`
  display: grid;
  grid-template-columns: 1fr 1fr;
  grid-gap: 10px;
`;

let initialValues = {
  virtualServiceName: '',
  displayName: '',
  namespace: '',
  domainsList: [] as string[]
};

const validationSchema = yup.object().shape({
  virtualServiceName: yup
    .string()
    .required('Name is required')
    .matches(/[a-z0-9]+[-.a-z0-9]*[a-z0-9]/, 'Name is invalid'),
  displayName: yup.string(),
  domainsList: yup
    .array()
    .of(yup.string())
    .required('At least one domain must be specified.'),
  namespace: yup
    .string()
    .required('Namespace is required')
    .matches(/[a-z0-9]+[-.a-z0-9]*[a-z0-9]/, 'Namespace is invalid')
});

interface Props {
  toggleModal: React.Dispatch<React.SetStateAction<boolean>>;
}

export const CreateVirtualServiceForm = (props: Props) => {
  let history = useHistory();
  const licenseError = useSelector((state: AppState) => state.modal.error);
  let location = useLocation();
  const { data: virtualServicesList, error } = useSWR(
    'listVirtualServices',
    virtualServiceAPI.listVirtualServices
  );
  const { data: namespacesList, error: listNamespacesError } = useSWR(
    'listNamespaces',
    configAPI.listNamespaces
  );
  const { data: podNamespace, error: podNamespaceError } = useSWR(
    'getPodNamespace',
    configAPI.getPodNamespace
  );

  const dispatch = useDispatch();
  if (!podNamespace) {
    return <div>Loading...</div>;
  }
  // this is to match the value displayed by the typeahead
  initialValues.namespace = podNamespace;

  const handleCreateVirtualService = (values: typeof initialValues) => {
    let { namespace, virtualServiceName, displayName, domainsList } = values;
    dispatch(
      createVirtualService({
        inputV2: {
          displayName: {
            value: displayName
          },
          domains: {
            valuesList: domainsList.filter(val => val !== '')
          },
          ref: {
            name: virtualServiceName,
            namespace
          }
        }
      })
    );

    setTimeout(() => {
      if (
        virtualServicesList?.find(
          vsD =>
            vsD.virtualService?.metadata?.name === values.virtualServiceName
        ) !== undefined
      ) {
        history.push({
          pathname: `${location.pathname}${values.namespace}/${values.virtualServiceName}`
        });
      }
    }, 500);
    props.toggleModal(s => !s);
  };

  return (
    <Formik
      initialValues={initialValues}
      onSubmit={handleCreateVirtualService}
      validationSchema={validationSchema}>
      {({ values, isSubmitting, handleSubmit }) => (
        <div>
          <InputContainer>
            <div>
              <SoloFormInput
                name='virtualServiceName'
                title='Virtual Service Name'
                placeholder='Virtual Service Name'
              />
            </div>
            <div>
              <SoloFormInput
                name='displayName'
                title='Display Name'
                placeholder='Display Name'
              />
            </div>
            <div>
              <SoloFormSelect
                label='Domains'
                name='domainsList'
                mode='tags'
                tokenSeparators={[',']}
                style={{ width: '100%' }}
                placeholder='Please select a domain'
                defaultValue={values.domainsList || []}>
                {values.domainsList.map((domain, index) => (
                  <Select.Option key={domain}>
                    <div
                      key={domain}
                      className='flex items-center mb-2 text-sm text-blue-600'>
                      {domain}
                    </div>
                  </Select.Option>
                ))}
              </SoloFormSelect>
            </div>
            <div>
              <SoloFormTypeahead
                name='namespace'
                title='Virtual Service Namespace'
                defaultValue={values.namespace}
                presetOptions={(namespacesList ?? []).map(ns => {
                  return { value: ns };
                })}
              />
            </div>
          </InputContainer>
          <Footer>
            <SoloButton
              onClick={handleSubmit}
              text='Create Virtual Service'
              disabled={isSubmitting}
            />
          </Footer>
        </div>
      )}
    </Formik>
  );
};
