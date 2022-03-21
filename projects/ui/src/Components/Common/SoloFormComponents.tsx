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
import { Input, InputProps, Label, SoloInput } from './SoloInput';
import { SoloTextarea, TextareaProps } from './SoloTextarea';
import { SoloToggleSwitch, SoloToggleSwitchProps } from './SoloToggleSwitch';
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
  onRemoveFile?: () => void;
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
    onRemoveFile,
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
        } col-span-2 gap-2 mb-2`}>
        <div
          className={`${
            horizontal ? 'flex items-center' : ''
          } col-span-2 -space-y-px bg-white rounded-md`}>
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

                if (!validateFile(file)?.isValid) {
                  setHasError(true);

                  setFieldError(name, validateFile(file)?.errorMessage);
                  setFileError('File Error');
                }
                if (!!onRemoveFile) {
                  onRemoveFile();
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
                      setFieldError(name, undefined);
                      setHasError(false);
                      setFileError('');
                    }}>
                    {'Remove File'}
                  </RemoveFileButton>
                )}
                <small
                  className={` flex-1 block w-full min-w-0 border-r-0 border-gray-300 rounded-none  sm:text-sm w-2/3 p-2 flex-1 block w-full min-w-0 text-base text-gray-500 border ${
                    values[name] !== undefined
                      ? 'border-l-0 rounded-l-none '
                      : 'rounded-l-md'
                  } border-r-0 border-gray-300 rounded-lg rounded-r-none`}>
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
                  }}>
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

export type SoloFormRadioOption = {
  displayValue: string;
  value: string;
  subHeader?: string;
};
export type SoloFormRadioProps = {
  name: string; // the name of this field in Formik
  isUpdate?: boolean; // whether its an update or create
  title?: string; // display name of the field
  options: SoloFormRadioOption[];
  horizontal?: boolean;
  titleAbove?: boolean;
  onChange?: (result: any) => void;
};

export const SoloFormRadio = <Values extends FormikValues>(
  props: SoloFormRadioProps
) => {
  const { name, isUpdate, title, options, horizontal, titleAbove, onChange } =
    props;
  const { values, setFieldValue } = useFormikContext<Values>();

  return (
    <>
      {title && titleAbove ? (
        <label className='mb-2 text-base font-medium'>{title}</label>
      ) : null}
      <div
        className={`${
          horizontal ? 'flex items-center' : 'grid grid-cols-2 '
        } col-span-2 gap-2 mb-2 cursor-pointer`}>
        {title && titleAbove ? null : (
          <label className='text-base font-medium '>{title}</label>
        )}
        <div
          className={`${
            horizontal ? 'flex items-center' : ''
          } col-span-2  bg-white rounded-md`}>
          {options.map((option, index) => {
            let isFirst = index === 0;
            let isLast = index === options.length - 1;
            return (
              <div
                key={option.displayValue}
                onClick={() => {
                  if (onChange) {
                    onChange(option.value);
                  }
                  if (isUpdate) return;
                  setFieldValue(name, option.value);
                }}
                className={`relative flex p-5 border ${
                  isUpdate ? '' : hoverStyles
                } ${focusStyles} ${
                  isUpdate && values[name] === option.value
                    ? 'bg-gray-400 border-gray-200'
                    : values[name] === option.value
                    ? ' border-blue-300gloo bg-blue-150gloo z-10 '
                    : 'border-gray-200'
                } ${
                  isFirst
                    ? horizontal
                      ? 'rounded-tl-md rounded-bl-md'
                      : 'rounded-tl-md rounded-tr-md'
                    : isLast
                    ? horizontal
                      ? 'rounded-tr-md rounded-br-md'
                      : 'rounded-bl-md rounded-br-md'
                    : ''
                } `}>
                <div className='flex items-center h-5'>
                  <input
                    type='radio'
                    readOnly
                    disabled={isUpdate}
                    className={`'w-4 h-4 border-gray-300 ${
                      isUpdate
                        ? 'text-gray-600 focus:ring-gray-600 cursor-not-allowed'
                        : 'text-blue-600gloo focus:ring-blue-600gloo cursor-pointer'
                    } '`}
                    checked={values[name] === option.value}
                  />
                </div>
                <label
                  htmlFor='settings-option-0'
                  className={`flex flex-col ml-3 ${
                    isUpdate ? 'cursor-not-allowed' : 'cursor-pointer'
                  }`}>
                  <span
                    className={`block text-sm font-medium ${
                      values[name] === option.value
                        ? ' text-blue-700gloo'
                        : 'text-gray-900'
                    } `}>
                    {option.displayValue}
                  </span>
                  <span className='block text-sm text-blue-700gloo'>
                    {option.subHeader ?? ''}
                  </span>
                </label>
              </div>
            );
          })}
        </div>
      </div>
    </>
  );
};

interface SoloFormInputProps extends Partial<InputProps> {
  name: string;
  title?: string;
  hideError?: boolean;
}

export const SoloFormInput = (props: SoloFormInputProps) => {
  const { hideError = false } = props;
  const form = useFormikContext<any>();
  const field = form.getFieldProps(props.name);
  const meta = form.getFieldMeta(props.name);

  return (
    <div>
      <SoloInput
        className='focus:border-blue-400gloo hover:border-blue-400gloo focus:outline-none focus:ring-4 focus:ring-blue-500gloo focus:ring-opacity-10'
        {...field}
        {...props}
        error={!!meta.error}
        title={props.title}
        value={field.value}
        onChange={field.onChange}
      />
      {hideError ? null : (
        <ErrorText errorExists={!!meta.error}>{meta.error}</ErrorText>
      )}
    </div>
  );
};

interface SoloFormToggleProps extends Partial<SoloToggleSwitchProps> {
  name: string;
  label?: string;
  hideError?: boolean;
  disabled?: boolean;
  finishOnChange?: (val: boolean) => void;
  testId?: string;
}
export const SoloFormToggle = (props: SoloFormToggleProps) => {
  const { hideError, finishOnChange } = props;
  const form = useFormikContext<any>();
  const field = form.getFieldProps(props.name);
  const meta = form.getFieldMeta(props.name);

  return (
    <>
      {props.label && <Label>{props.label}</Label>}
      <SoloToggleSwitch
        {...props}
        {...field}
        checked={!!field.value}
        onChange={value => {
          form.setFieldValue(field.name, value);
          typeof finishOnChange === 'function' && finishOnChange(value);
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
