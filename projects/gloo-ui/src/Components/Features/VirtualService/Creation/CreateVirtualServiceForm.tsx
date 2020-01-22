import styled from '@emotion/styled';
import {
  SoloFormInput,
  SoloFormStringsList,
  SoloFormTypeahead
} from 'Components/Common/Form/SoloFormField';
import { SoloButton } from 'Components/Common/SoloButton';
import { Formik } from 'formik';
import * as React from 'react';
import { useDispatch } from 'react-redux';
import { useHistory, useLocation } from 'react-router';
import { configAPI } from 'store/config/api';
import { createVirtualService } from 'store/virtualServices/actions';
import useSWR from 'swr';
import * as yup from 'yup';

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
  domainsList: ['']
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
  let location = useLocation();

  const { data: namespacesList, error: listNamespacesError } = useSWR(
    'listNamespaces',
    configAPI.listNamespaces
  );
  const { data: podNamespace, error: podNamespaceError } = useSWR(
    'getPodNamespace',
    configAPI.getPodNamespace
  );

  const dispatch = useDispatch();
  if (!podNamespace || !namespacesList) {
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
      history.push({
        pathname: `${location.pathname}${values.namespace}/${values.virtualServiceName}`
      });
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
              <SoloFormStringsList
                name='domainsList'
                label='Domains'
                createNewPromptText='Specify a domain'
              />
            </div>
            <div>
              <SoloFormTypeahead
                name='namespace'
                title='Virtual Service Namespace'
                defaultValue={values.namespace}
                presetOptions={namespacesList.map(ns => {
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
