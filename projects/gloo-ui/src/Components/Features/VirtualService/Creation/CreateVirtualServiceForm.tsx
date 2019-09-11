import styled from '@emotion/styled';
import {
  SoloFormInput,
  SoloFormTypeahead
} from 'Components/Common/Form/SoloFormField';
import { SoloButton } from 'Components/Common/SoloButton';
import { Formik } from 'formik';
import * as React from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { AppState } from 'store';
import { createVirtualService } from 'store/virtualServices/actions';
import * as yup from 'yup';
import { withRouter, RouteComponentProps } from 'react-router-dom';

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
  namespace: ''
};

const validationSchema = yup.object().shape({
  virtualServiceName: yup
    .string()
    .required('Name is required')
    .matches(/[a-z0-9]+[-.a-z0-9]*[a-z0-9]/, 'Name is invalid'),
  displayName: yup.string(),
  namespace: yup
    .string()
    .required('Namespace is required')
    .matches(/[a-z0-9]+[-.a-z0-9]*[a-z0-9]/, 'Namespace is invalid')
});

interface Props extends RouteComponentProps {
  toggleModal: React.Dispatch<React.SetStateAction<boolean>>;
}

export const CreateVirtualServiceForm = withRouter((props: Props) => {
  const {
    config: { namespacesList, namespace: podNamespace }
  } = useSelector((state: AppState) => state);
  const dispatch = useDispatch();
  // this is to match the value displayed by the typeahead
  initialValues.namespace = podNamespace;

  const handleCreateVirtualService = (values: typeof initialValues) => {
    let { namespace, virtualServiceName, displayName } = values;

    dispatch(
      createVirtualService({
        inputV2: {
          displayName: {
            value: displayName
          },
          ref: {
            name: virtualServiceName,
            namespace
          }
        }
      })
    );
    setTimeout(() => {
      props.history.push({
        pathname: `${props.match.path}${values.namespace}/${values.virtualServiceName}`
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
});
