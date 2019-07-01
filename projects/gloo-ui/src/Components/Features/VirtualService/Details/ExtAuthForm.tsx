import * as React from 'react';
import styled from '@emotion/styled/macro';
import { SoloInput } from 'Components/Common/SoloInput';
import { SoloButton } from 'Components/Common/SoloButton';
import { useFormValidation } from 'Hooks/useFormValidation';
import { colors } from 'Styles';
import { SoloNegativeButton } from 'Styles/CommonEmotions/button';

const FormContainer = styled.div`
  display: grid;
  grid-template-columns: 1fr 1fr;
  grid-template-rows: 1fr 1fr 1fr;
  padding-top: 10px;
  grid-gap: 8px;
`;

const ErrorText = styled.div`
  color: ${colors.grapefruitOrange};
`;

const initialValues = {
  clientId: '',
  callbackPath: '',
  issuerURL: '',
  appURL: '',
  secretRefNamespace: '',
  secretRefName: ''
};

const validate = (values: typeof initialValues) => {
  let errors = {} as typeof initialValues;
  if (!values.clientId) {
    errors.clientId = 'Need a client ID ';
  }
  if (!values.callbackPath) {
    errors.callbackPath = 'Need a callback path ';
  }
  if (!values.issuerURL) {
    errors.issuerURL = 'Need an issuer url ';
  }

  return errors;
};

export const ExtAuthForm = () => {
  const {
    handleSubmit,
    handleChange,
    handleBlur,
    values,
    errors,
    isSubmitting,
    isDifferent
  } = useFormValidation(initialValues, validate, updateExtAuth);

  function updateExtAuth() {
    console.log('do grpc request');
  }

  return (
    <div>
      <FormContainer>
        <div>
          <SoloInput
            title='Client ID'
            name='clientId'
            value={values.clientId}
            placeholder={'example'}
            onChange={handleChange}
            onBlur={handleBlur}
          />
          {errors && <ErrorText>{errors.clientId}</ErrorText>}
        </div>
        <div>
          <SoloInput
            title='Callback Path'
            name='callbackPath'
            value={values.callbackPath}
            placeholder={'example'}
            onChange={handleChange}
            onBlur={handleBlur}
          />
          {errors && <ErrorText>{errors.callbackPath}</ErrorText>}
        </div>
        <div>
          <SoloInput
            title='Issuer URL'
            name='issuerURL'
            value={values.issuerURL}
            placeholder={'example'}
            onChange={handleChange}
            onBlur={handleBlur}
          />
          {errors && <ErrorText>{errors.issuerURL}</ErrorText>}
        </div>
        <div>
          <SoloInput
            title='App URL'
            name='appURL'
            value={values.appURL}
            placeholder={'example'}
            onChange={handleChange}
            onBlur={handleBlur}
          />
          {errors && <ErrorText>{errors.appURL}</ErrorText>}
        </div>
        <div>
          <SoloInput
            title='Secret Ref Namespace'
            name='secretRefNamespace'
            value={values.secretRefNamespace}
            placeholder={'Secret Ref Namespace'}
            onChange={handleChange}
            onBlur={handleBlur}
          />
          {errors && <ErrorText>{errors.secretRefNamespace}</ErrorText>}
        </div>
        <div>
          <SoloInput
            title='Secret Ref Name'
            name='secretRefName'
            value={values.secretRefName}
            placeholder={'example'}
            onChange={handleChange}
            onBlur={handleBlur}
          />
          {errors && <ErrorText>{errors.secretRefName}</ErrorText>}
        </div>
        <SoloNegativeButton>Clear</SoloNegativeButton>
        <SoloButton
          onClick={handleSubmit!}
          text='Submit'
          disabled={isSubmitting}
        />
      </FormContainer>
    </div>
  );
};
