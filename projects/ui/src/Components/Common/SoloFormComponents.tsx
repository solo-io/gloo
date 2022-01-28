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
import React, { useRef } from 'react';
import { colors } from 'Styles/colors';
import tw from 'twin.macro';

import { SoloCheckbox, CheckboxProps } from './SoloCheckbox';
import { SoloDropdown, DropdownProps } from './SoloDropdown';
import { Input, Label } from './SoloInput';
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

export const ButtonWithErrorState = styled.button<{
  hasError: boolean;
  isFileUploaded?: boolean;
  disabled?: boolean;
}>`
  ${tw`w-1/3 p-2 text-white border-l-0 rounded-lg rounded-l-none focus:outline-none`}
  ${props =>
    props.hasError
      ? `background-color:  #B05464;`
      : props.isFileUploaded
      ? `background-color: #0DCE93;
                            border-color: #0DCE93;`
      : props.disabled
      ? `background-color: ${colors.mayGrey};`
      : `background-color: ${colors.seaBlue}; 
                            border-color: ${colors.seaBlue};`}
`;

export const RemoveFileButton = styled.button<{
  hasError: boolean;
}>`
  ${tw`w-1/5 p-2 text-white border border-r-0 rounded-lg rounded-r-none focus:outline-none`}
  background-color: ${props => (props.hasError ? '#B05464' : colors.seaBlue)};
  border-color: ${props => (props.hasError ? '#B05464' : colors.seaBlue)};
`;
export type SoloFormFileUploadProps = {
  name: string; // the name of this field in Formik
  isUpdate?: boolean; // whether its an update or create
  title?: string; // display name of the field
  horizontal?: boolean;
  titleAbove?: boolean;
  isDisabled?: boolean;
  buttonLabel?: string;
  fileType?: string;
  validateFile?: (file?: File) => {
    isValid: boolean;
    errorMessage: string;
  };
};

export const SoloFormFileUpload = <Values extends FormikValues>(
  props: SoloFormFileUploadProps
) => {
  const {
    name,
    title,
    horizontal,
    titleAbove,
    isDisabled,
    buttonLabel = 'Upload File',
    fileType = 'application/json,application/x-yaml,text/*',
    validateFile = () => ({ isValid: true, errorMessage: '' }),
  } = props;
  const { values, setFieldValue, setFieldError } = useFormikContext<Values>();
  let fileInput = useRef<HTMLInputElement>(null);
  const [hasError, setHasError] = React.useState(false);

  const [fileError, setFileError] = React.useState('');

  return (
    <>
      {title && titleAbove ? (
        <label className='mt-2 mb-1 text-base font-medium text-gray-800'>
          {title}
        </label>
      ) : null}
      <div
        className={`${
          horizontal ? 'flex items-center' : 'grid grid-cols-2 '
        } col-span-2 gap-2 mb-2`}
      >
        <div
          className={`${
            horizontal ? 'flex items-center' : ''
          } col-span-2 -space-y-px bg-white rounded-md`}
        >
          {title && titleAbove ? null : <Label>{title}</Label>}
          <Input
            type='file'
            style={{ display: 'none' }}
            ref={fileInput}
            accept={fileType}
            value={''}
            onError={e => {
              setHasError(true);
              setFileError('There was an error with file upload');
            }}
            disabled={isDisabled}
            onChange={event => {
              if (event.currentTarget?.files) {
                let file = event.currentTarget.files[0];
                setFieldValue(name, event.currentTarget.files[0]);
                setHasError(false);
                console.log('fileType.split(', ')', fileType.split(','));
                console.log('file.type', file);

                if (!validateFile(file)?.isValid) {
                  setHasError(true);

                  setFieldError(name, validateFile(file)?.errorMessage);
                  setFileError('File Error');
                }
              }
            }}
          />
          <div className='w-full mb-6 sm:grid sm:grid-cols-3 sm:gap-4 sm:items-start sm:border-gray-200 sm:pt-2'>
            <div className='mt-1 sm:mt-0 sm:col-span-3'>
              <div className='flex rounded-md '>
                {values[name] !== undefined && (
                  <RemoveFileButton
                    hasError={hasError}
                    type='button'
                    onClick={() => {
                      setFieldValue(name, undefined as unknown as File);
                      setHasError(false);
                      setFileError('');
                    }}
                  >
                    {'Remove File'}
                  </RemoveFileButton>
                )}
                <small
                  className={` flex-1 block w-full min-w-0 border-r-0 border-gray-300 rounded-none  sm:text-sm w-2/3 p-2 flex-1 block w-full min-w-0 text-base text-gray-500 border ${
                    values[name] !== undefined
                      ? 'border-l-0 rounded-l-none '
                      : 'rounded-l-md'
                  } border-r-0 border-gray-300 rounded-lg rounded-r-none`}
                >
                  {values[name]
                    ? values[name].name && !hasError
                      ? values[name].name
                      : fileError ?? 'Error'
                    : 'No file chosen'}
                </small>

                <ButtonWithErrorState
                  hasError={hasError}
                  isFileUploaded={!!values[name]}
                  type='button'
                  disabled={isDisabled}
                  onClick={e => {
                    fileInput?.current?.click();
                  }}
                >
                  {hasError
                    ? 'Upload Error'
                    : !!values[name]
                    ? 'Upload successful'
                    : buttonLabel}
                </ButtonWithErrorState>
              </div>
            </div>
          </div>
        </div>
      </div>
    </>
  );
};
