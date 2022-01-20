import { useTheme } from '@emotion/react';
import styled from '@emotion/styled';
import { Select } from 'antd';
import { CheckboxChangeEvent } from 'antd/lib/checkbox';
import { SelectProps, SelectValue } from 'antd/lib/select';
import { ReactComponent as NoImageIcon } from 'assets/no-image-placeholder.svg';
import { ReactComponent as PlaceholderPortal } from 'assets/placeholder-portal.svg';
import { ReactComponent as RemoveIcon } from 'assets/remove.svg';
import {
  Field,
  FieldProps,
  FormikContextType,
  FormikValues,
  useFormikContext,
} from 'formik';
import React from 'react';
import { colors } from 'Styles/colors';

import { SoloCheckbox, CheckboxProps } from './SoloCheckbox';
import { SoloDropdown, DropdownProps } from './SoloDropdown';
import { SoloTextarea, TextareaProps } from './SoloTextarea';
type ErrorTextProps = { errorExists?: boolean };

// focus rings for form elements
export const focusStyles =
  'focus:border-blue-400gloo focus:outline-none focus:ring-4 focus:ring-blue-500gloo focus:ring-opacity-10';

export const hoverStyles = 'hover:border-blue-400gloo';

export const ErrorText = React.memo(
  styled.div`
    color: ${colors.grapefruitOrange};
    visibility: ${(props: ErrorTextProps) =>
      props.errorExists ? 'visible' : 'hidden'};
    height: 1rem;
    overflow: auto;
  `,
  // We want to allow updating when either error state changes or error message changes
  (prev, curr) =>
    prev.errorExists === curr.errorExists && prev.children === curr.children
);

interface SoloFormDropdownProps extends Partial<DropdownProps> {
  name: string;
  hideError?: boolean;
}
export const SoloFormDropdown = (props: SoloFormDropdownProps) => {
  const { hideError } = props;
  const form = useFormikContext<any>();
  const field = form.getFieldProps(props.name);
  const meta = form.getFieldMeta(props.name);
  return (
    <>
      <SoloDropdown
        {...field}
        {...props}
        error={!!meta.error && meta.touched}
        onChange={(newVal: SelectValue) => {
          form.setFieldValue(field.name, newVal as string);
          props.onChange && props.onChange(newVal);
        }}
      />
      {hideError ? null : (
        <ErrorText errorExists={!!meta.error && meta.touched}>
          {meta.error}
        </ErrorText>
      )}
    </>
  );
};
