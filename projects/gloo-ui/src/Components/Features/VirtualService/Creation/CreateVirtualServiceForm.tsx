import styled from '@emotion/styled';
import { useCreateVirtualService } from 'Api/useVirtualServiceClient';
import { SoloButton } from 'Components/Common/SoloButton';
import { SoloInput } from 'Components/Common/SoloInput';
import { SoloTypeahead } from 'Components/Common/SoloTypeahead';
import { useFormValidation } from 'Hooks/useFormValidation';
import { ResourceRef } from 'proto/github.com/solo-io/solo-kit/api/v1/ref_pb';
import {
  CreateVirtualServiceRequest,
  VirtualServiceInput
} from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/virtualservice_pb';
import * as React from 'react';
import { useSelector } from 'react-redux';
import { AppState } from 'store';
import { ErrorText } from '../Details/ExtAuthForm';

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
  addDomain: '',
  namespace: ''
};

const validate = (values: typeof initialValues) => {
  let errors = {} as typeof initialValues;
  if (!values.virtualServiceName) {
    errors.virtualServiceName = 'Name is required';
  }
  if (values.virtualServiceName.toLowerCase() !== values.virtualServiceName) {
    errors.virtualServiceName = 'Letters in a name may only be lower-case';
  }
  if (values.virtualServiceName.length >= 254) {
    errors.virtualServiceName = 'Names must be 253 characters or less';
  }

  if (!values.namespace) {
    errors.namespace = 'Namespace is required';
  }

  if (!!values.displayName && values.displayName.length > 500) {
    errors.displayName = 'Display name cannot be longer than 500 characters';
  }

  return errors;
};

interface Props {
  onCompletion?: (succeeded?: { namespace: string; name: string }) => any;
}

export const CreateVirtualServiceForm = (props: Props) => {
  const {
    config: { namespacesList, namespace: podNamespace }
  } = useSelector((state: AppState) => state);
  // this is to match the value displayed by the typeahead
  initialValues.namespace = podNamespace;
  const {
    handleSubmit,
    handleChange,
    handleBlur,
    values,
    errors,
    isSubmitting,
    isDifferent
  } = useFormValidation(initialValues, validate, createVirtualService);

  const { refetch: makeRequest } = useCreateVirtualService(null);

  function createVirtualService(values: typeof initialValues) {
    let vsRequest = new CreateVirtualServiceRequest();
    let vsInput = new VirtualServiceInput();

    let vsRef = new ResourceRef();
    vsRef.setName(values.virtualServiceName);
    vsRef.setNamespace(values.namespace);
    vsInput.setRef(vsRef);

    vsInput.setDisplayName(values.displayName);

    vsRequest.setInput(vsInput);
    makeRequest(vsRequest);

    if (!!props.onCompletion) {
      props.onCompletion({
        name: values.virtualServiceName,
        namespace: values.namespace
      });
    }
  }

  const isSubmittable = () => {
    return isDifferent && !Object.keys(errors).length && !isSubmitting;
  };

  return (
    <div>
      <InputContainer>
        <div>
          <SoloInput
            title='Virtual Service Name'
            name='virtualServiceName'
            value={values.virtualServiceName}
            placeholder={'Virtual Service Name'}
            onChange={handleChange}
            onBlur={handleBlur}
          />
          {errors && <ErrorText>{errors.virtualServiceName}</ErrorText>}
        </div>
        <div>
          <SoloInput
            title='Display Name'
            name='displayName'
            value={values.displayName}
            placeholder={'Display Name'}
            onChange={handleChange}
            onBlur={handleBlur}
          />
          {errors && <ErrorText>{errors.displayName}</ErrorText>}
        </div>
        <div>
          <SoloTypeahead
            title='Virtual Service Namespace'
            defaultValue={values.namespace}
            onChange={e => handleChange(e, 'namespace')}
            presetOptions={namespacesList.map(ns => {
              return { value: ns };
            })}
          />
          {errors && <ErrorText>{errors.namespace}</ErrorText>}
        </div>
      </InputContainer>
      <Footer>
        <SoloButton
          onClick={handleSubmit}
          text='Create Virtual Service'
          disabled={!isSubmittable()}
        />
      </Footer>
    </div>
  );
};
