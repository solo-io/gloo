import React from 'react';
import { FieldProps } from 'formik';
import { InputProps, SoloInput } from './SoloInput';
import { SoloDropdown, DropdownProps } from './SoloDropdown';
import { SoloTypeahead, TypeaheadProps } from './SoloTypeahead';
import styled from '@emotion/styled/macro';
import { colors } from 'Styles';

const ErrorText = styled.div`
  color: ${colors.grapefruitOrange};
`;
// TODO: make these wrappers generic to avoid repetition
export const SoloFormInput: React.FC<FieldProps & InputProps> = ({
  error,
  field,
  form: { errors },
  ...rest
}) => {
  return (
    <React.Fragment>
      <SoloInput error={!!errors[field.name]} {...rest} {...field} />
      {errors && <ErrorText>{errors[field.name]}</ErrorText>}
    </React.Fragment>
  );
};

export const SoloFormTypeahead: React.FC<FieldProps & TypeaheadProps> = ({
  field,
  form: { errors, setFieldValue },
  ...rest
}) => {
  return (
    <React.Fragment>
      <SoloTypeahead
        {...rest}
        {...field}
        onChange={value => setFieldValue(field.name, value)}
      />
      {errors && <ErrorText>{errors[field.name]}</ErrorText>}
    </React.Fragment>
  );
};

export const SoloFormDropdown: React.FC<FieldProps & DropdownProps> = ({
  field,
  form: { errors, setFieldValue },
  ...rest
}) => {
  return (
    <React.Fragment>
      <SoloDropdown
        {...rest}
        {...field}
        onChange={value => setFieldValue(field.name, value)}
      />
      {errors && <ErrorText>{errors[field.name]}</ErrorText>}
    </React.Fragment>
  );
};
