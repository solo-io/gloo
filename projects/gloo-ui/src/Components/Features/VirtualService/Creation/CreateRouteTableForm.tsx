import styled from '@emotion/styled';
import {
  SoloFormInput,
  SoloFormTypeahead
} from 'Components/Common/Form/SoloFormField';
import { SoloButton } from 'Components/Common/SoloButton';
import { Formik } from 'formik';
import { Metadata } from 'proto/solo-kit/api/v1/metadata_pb';
import * as React from 'react';
import { useDispatch } from 'react-redux';
import { useHistory, useLocation } from 'react-router';
import { configAPI } from 'store/config/api';
import { createRouteTable } from 'store/routeTables/actions';
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
  routeTableName: '',
  namespace: ''
};

const validationSchema = yup.object().shape({
  routeTableName: yup
    .string()
    .required('Name is required')
    .matches(/[a-z0-9]+[-.a-z0-9]*[a-z0-9]/, 'Name is invalid'),
  namespace: yup
    .string()
    .required('Namespace is required')
    .matches(/[a-z0-9]+[-.a-z0-9]*[a-z0-9]/, 'Namespace is invalid')
});

interface Props {
  toggleModal: React.Dispatch<React.SetStateAction<boolean>>;
}

export const CreateRouteTableForm = (props: Props) => {
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

  const handleCreateRouteTable = (values: typeof initialValues) => {
    let { namespace, routeTableName } = values;
    let newMetadata = new Metadata().toObject();
    dispatch(
      createRouteTable({
        routeTable: {
          metadata: {
            ...newMetadata,
            name: routeTableName,
            namespace
          },
          routesList: []
        }
      })
    );
    setTimeout(() => {
      history.push({
        pathname: `/routetables/${values.namespace}/${values.routeTableName}`
      });
    }, 500);
    props.toggleModal(s => !s);
  };

  return (
    <Formik
      initialValues={initialValues}
      onSubmit={handleCreateRouteTable}
      validationSchema={validationSchema}>
      {({ values, isSubmitting, handleSubmit }) => (
        <div>
          <InputContainer>
            <div>
              <SoloFormInput
                name='routeTableName'
                title='Route Table Name'
                placeholder='Route Table Name'
              />
            </div>

            <div>
              <SoloFormTypeahead
                name='namespace'
                title='Route Table Namespace'
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
              text='Create Route Table'
              disabled={isSubmitting}
            />
          </Footer>
        </div>
      )}
    </Formik>
  );
};
