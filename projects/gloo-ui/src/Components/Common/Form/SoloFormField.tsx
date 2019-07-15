import styled from '@emotion/styled/macro';
import {
  FieldProps,
  FieldArrayRenderProps,
  Form,
  FieldArray,
  Field
} from 'formik';
import React from 'react';
import { colors } from 'Styles';
import { DropdownProps, SoloDropdown } from '../SoloDropdown';
import { InputProps, SoloInput } from '../SoloInput';
import { SoloTypeahead, TypeaheadProps } from '../SoloTypeahead';
import { SoloCheckbox, CheckboxProps } from '../SoloCheckbox';
import { staticInitialValues } from 'Components/Features/Upstream/Creation/StaticUpstreamForm';
import { SoloButton } from '../SoloButton';
import { SoloMultiSelect } from '../SoloMultiSelect';
import { MultipartStringCardsList } from '../MultipartStringCardsList';
import { createUpstreamId, parseUpstreamId } from 'utils/helpers';
const ErrorText = styled.div`
  color: ${colors.grapefruitOrange};
`;
// TODO: make these wrappers generic to avoid repetition
export const SoloFormInput: React.FC<FieldProps & InputProps> = ({
  error,
  field,
  form: { errors },
  ...rest
}) => (
  <React.Fragment>
    <SoloInput error={!!errors[field.name]} {...field} {...rest} />
    {errors && <ErrorText>{errors[field.name]}</ErrorText>}
  </React.Fragment>
);

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
        {...field}
        {...rest}
        onChange={value => setFieldValue(field.name, value)}
      />
      {errors && <ErrorText>{errors[field.name]}</ErrorText>}
    </React.Fragment>
  );
};

interface MetadataBasedDropdownProps extends DropdownProps {
  value: any;
  options: any[];
}

export const SoloFormMetadataBasedDropdown: React.FC<
  FieldProps & MetadataBasedDropdownProps
> = ({ field, form: { errors, setFieldValue }, ...rest }) => {
  const usedOptions = rest.options.map(option => {
    return {
      key: createUpstreamId(option.metadata!), // the same as virtual service's currently
      displayValue: option.metadata.name,
      value: createUpstreamId(option.metadata!)
    };
  });

  const usedValue =
    rest.value && rest.value.metadata ? rest.value.metadata.name : undefined;

  const setNewValue = (newValueId: any) => {
    const { name, namespace } = parseUpstreamId(newValueId);
    const optionChosen = rest.options.find(
      option =>
        option.metadata.name === name && option.metadata.namespace === namespace
    );

    setFieldValue(field.name, optionChosen);
  };

  return (
    <React.Fragment>
      <SoloDropdown
        {...field}
        {...rest}
        options={usedOptions}
        value={usedValue}
        onChange={setNewValue}
      />
      {errors && <ErrorText>{errors[field.name]}</ErrorText>}
    </React.Fragment>
  );
};

export const SoloFormMultiselect: React.FC<FieldProps & DropdownProps> = ({
  field,
  form: { errors, setFieldValue },
  ...rest
}) => {
  return (
    <React.Fragment>
      <SoloMultiSelect
        {...field}
        {...rest}
        values={Object.keys(field.value).filter(key => field.value[key])}
        onChange={newValues => {
          const newFieldValues = { ...field.value };
          Object.keys(newFieldValues).forEach(key => {
            newFieldValues[key] = false;
          });
          for (let val of newValues) {
            newFieldValues[val] = true;
          }

          setFieldValue(field.name, newFieldValues);
        }}
      />
      {errors && <ErrorText>{errors[field.name]}</ErrorText>}
    </React.Fragment>
  );
};

export const SoloFormCheckbox: React.FC<FieldProps & CheckboxProps> = ({
  field,
  form: { errors, setFieldValue },
  ...rest
}) => {
  return (
    <React.Fragment>
      <SoloCheckbox
        {...rest}
        {...field}
        checked={!!field.value}
        onChange={value => setFieldValue(field.name, value.target.checked)}
        label
      />
      {errors && <ErrorText>{errors[field.name]}</ErrorText>}
    </React.Fragment>
  );
};

export const SoloFormMultipartStringCardsList: React.FC<
  FieldProps & DropdownProps
> = ({ field, form: { errors, setFieldValue }, ...rest }) => {
  return (
    <React.Fragment>
      <MultipartStringCardsList
        {...field}
        {...rest}
        values={field.value}
        valueDeleted={indexDeleted => {
          console.log(indexDeleted);
          setFieldValue(field.name, [...field.value].splice(indexDeleted, 1));
        }}
        createNew={newPair => {
          let newList = [...field.value];
          newList.push({
            value: newPair.newValue,
            name: newPair.newName
          });
          setFieldValue(field.name, newList);
        }}
      />
      {errors && <ErrorText>{errors[field.name]}</ErrorText>}
    </React.Fragment>
  );
};
