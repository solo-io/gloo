import * as React from 'react';
import styled from '@emotion/styled/macro';
import { useFormValidation } from 'Hooks/useFormValidation';
import { colors } from 'Styles';
import { Label, SoloInput } from 'Components/Common/SoloInput';
import { SoloNegativeButton } from 'Styles/CommonEmotions/button';
import { SoloButton } from 'Components/Common/SoloButton';

const ErrorText = styled.div`
  color: ${colors.grapefruitOrange};
`;

const FormContainer = styled.div`
  display: grid;
  grid-template-columns: 1fr 1fr;
  grid-gap: 8px;
`;

const FormItem = styled.div`
  display: flex;
  flex-direction: column;
`;

const InputContainer = styled.div`
  display: flex;
  flex-direction: row;
  align-items: baseline;
`;

const FormFooter = styled.div`
  grid-column: 2;
  display: grid;
  grid-template-columns: 1fr 1fr;
  grid-template-rows: 33px;
  grid-gap: 5px;
`;

let initialValues = {
  authLimNumber: 0,
  authLimTimeUnit: '',
  anonLimNumber: 0,
  anonLimTimeUnit: ''
};

const validate = (values: typeof initialValues) => {
  let errors = {} as typeof initialValues;

  return errors;
};

export const RateLimitForm = () => {
  const {
    handleSubmit,
    handleChange,
    handleBlur,
    values,
    errors,
    isSubmitting,
    isDifferent
  } = useFormValidation(initialValues, validate, updateRateLimit);

  function updateRateLimit() {
    console.log('grpc request');
  }

  return (
    <FormContainer>
      <FormItem>
        <Label>Authorized Limits</Label>
        <InputContainer>
          <SoloInput
            name='authLimNumber'
            value={values.authLimNumber}
            placeholder={'example'}
            onChange={handleChange}
            onBlur={handleBlur}
          />
          {errors && <ErrorText>{errors.authLimNumber}</ErrorText>}
          <div>per</div>
          <SoloInput
            name='authLimTimeUnit'
            value={values.authLimTimeUnit}
            placeholder={'example'}
            onChange={handleChange}
            onBlur={handleBlur}
          />
          {errors && <ErrorText>{errors.authLimTimeUnit}</ErrorText>}
        </InputContainer>
      </FormItem>
      <FormItem>
        <Label>Anonymous Limits</Label>
        <InputContainer>
          <SoloInput
            name='anonLimNumber'
            value={values.anonLimNumber}
            placeholder={'example'}
            onChange={handleChange}
            onBlur={handleBlur}
          />
          {errors && <ErrorText>{errors.anonLimNumber}</ErrorText>}
          per
          <SoloInput
            name='anonLimTimeUnit'
            value={values.anonLimTimeUnit}
            placeholder={'example'}
            onChange={handleChange}
            onBlur={handleBlur}
          />
          {errors && <ErrorText>{errors.anonLimTimeUnit}</ErrorText>}
        </InputContainer>
      </FormItem>
      <FormFooter>
        <SoloNegativeButton>Clear</SoloNegativeButton>
        <SoloButton
          onClick={handleSubmit}
          text='Submit'
          disabled={isSubmitting}
        />
      </FormFooter>
    </FormContainer>
  );
};
