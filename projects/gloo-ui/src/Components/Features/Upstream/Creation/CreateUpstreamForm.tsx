import * as React from 'react';
/** @jsx jsx */
import { jsx } from '@emotion/core';

import styled from '@emotion/styled/macro';
import { withRouter, RouteComponentProps } from 'react-router';
import { colors } from 'Styles';
import { SoloInput } from 'Components/Common/SoloInput';
import { useFormValidation } from 'Hooks/useFormValidation';

import { ErrorText } from '../../VirtualService/Details/ExtAuthForm';
import { SoloTypeahead } from '../../../Common/SoloTypeahead';
import { SoloButton } from 'Components/Common/SoloButton';

interface Props {}

let initialValues = {
  name: '',
  type: '',
  namespace: ''
};

const FormContainer = styled.div`
  display: flex;
  flex-direction: column;
`;

const InputContainer = styled.div`
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  padding: 15px;
`;

const Footer = styled.div`
  display: flex;
  flex-direction: row;
  justify-content: flex-end;
  padding-top: 15px;
`;

const validate = () => {
  return {
    name: '',
    type: '',
    namespace: ''
  };
};
export const CreateUpstreamForm = (props: Props) => {
  const {
    handleSubmit,
    handleChange,
    handleBlur,
    values,
    errors,
    isSubmitting,
    isDifferent
  } = useFormValidation(initialValues, validate, createUpstream);

  // grpc request
  function createUpstream() {}

  return (
    <FormContainer>
      <InputContainer>
        <div>
          <SoloInput
            title='Upstream Name'
            name='name'
            value={values.name}
            placeholder={'example'}
            onChange={handleChange}
            onBlur={handleBlur}
          />
          {errors && <ErrorText>{errors.name}</ErrorText>}
        </div>
        <div>
          <SoloInput
            title='Upstream Type'
            name='type'
            value={values.type}
            placeholder={'example'}
            onChange={handleChange}
            onBlur={handleBlur}
          />
          {errors && <ErrorText>{errors.type}</ErrorText>}
        </div>
        <div>
          <SoloInput
            title='Upstream Namespace'
            name='namespace'
            value={values.namespace}
            placeholder={'example'}
            onChange={handleChange}
            onBlur={handleBlur}
          />
          {errors && <ErrorText>{errors.namespace}</ErrorText>}
        </div>
      </InputContainer>
      <Footer>
        <SoloButton
          onClick={handleSubmit!}
          text='Create Upstream'
          disabled={isSubmitting}
        />
      </Footer>
    </FormContainer>
  );
};
