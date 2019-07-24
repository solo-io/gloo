import styled from '@emotion/styled/macro';
import { FieldProps, Field, FormikProps } from 'formik';
import React from 'react';
import { colors } from 'Styles';
import { DropdownProps, SoloDropdown } from '../SoloDropdown';
import { InputProps, SoloInput } from '../SoloInput';
import { SoloTypeahead, TypeaheadProps } from '../SoloTypeahead';
import { SoloCheckbox, CheckboxProps } from '../SoloCheckbox';
import { SoloMultiSelect, MultiselectProps } from '../SoloMultiSelect';
import {
  MultipartStringCardsList,
  MultipartStringCardsProps
} from '../MultipartStringCardsList';
import { createUpstreamId, parseUpstreamId } from 'utils/helpers';
import { NamespacesContext } from 'GlooIApp';
import { ListSecretsRequest } from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/secret_pb';
import { useListSecrets } from 'Api';
import { StringCardsListProps, StringCardsList } from '../StringCardsList';

export const ErrorText = styled<'div', { errorExists?: boolean }>('div')`
  color: ${colors.grapefruitOrange};
  visibility: ${props => (props.errorExists ? 'visible' : 'hidden')};
  min-height: 19px;
`;

// TODO: make these wrappers generic to avoid repetition
export const SoloFormInput: React.FC<FieldProps & InputProps> = ({
  error,
  field,
  form: { errors, ...form },
  ...rest
}) => {
  return (
    <React.Fragment>
      <SoloInput error={!!errors[field.name]} {...field} {...rest} />
      <ErrorText errorExists={!!errors}>{errors[field.name]}</ErrorText>
    </React.Fragment>
  );
};

export const SoloFormTypeahead: React.FC<FieldProps & TypeaheadProps> = ({
  field,
  form: { errors, setFieldValue, ...form },
  ...rest
}) => {
  return (
    <React.Fragment>
      <SoloTypeahead
        {...rest}
        {...field}
        onChange={value => setFieldValue(field.name, value)}
      />
      <ErrorText errorExists={!!errors}>{errors[field.name]}</ErrorText>
    </React.Fragment>
  );
};

export const SoloFormDropdown: React.FC<FieldProps & DropdownProps> = ({
  field,
  form: { errors, setFieldValue, ...form },
  ...rest
}) => {
  return (
    <React.Fragment>
      <SoloDropdown
        {...field}
        {...rest}
        onChange={value => setFieldValue(field.name, value)}
      />
      <ErrorText errorExists={!!errors}>{errors[field.name]}</ErrorText>
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
    rest.value && rest.value.metadata
      ? rest.value.metadata.name
      : field.value && field.value.metadata
      ? field.value.metadata.name
      : undefined;

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
      <ErrorText errorExists={!!errors}>{errors[field.name]}</ErrorText>
    </React.Fragment>
  );
};

export const SoloFormMultiselect: React.FC<FieldProps & MultiselectProps> = ({
  field,
  form: { errors, setFieldValue, ...form },
  ...rest
}) => {
  return (
    <React.Fragment>
      <SoloMultiSelect
        {...field}
        {...rest}
        values={
          !!rest.values
            ? rest.values
            : Object.keys(field.value).filter(key => field.value[key])
        }
        onChange={newValues => {
          const newFieldValues = !!rest.values
            ? { ...rest.values }
            : { ...field.value };
          Object.keys(newFieldValues).forEach(key => {
            newFieldValues[key] = false;
          });
          for (let val of newValues) {
            newFieldValues[val] = true;
          }

          setFieldValue(field.name, newFieldValues);
        }}
      />
      <ErrorText errorExists={!!errors}>{errors[field.name]}</ErrorText>
    </React.Fragment>
  );
};

export const SoloFormCheckbox: React.FC<FieldProps & CheckboxProps> = ({
  field,
  form: { errors, setFieldValue, ...form },
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
      <ErrorText errorExists={!!errors}>{errors[field.name]}</ErrorText>
    </React.Fragment>
  );
};

export const SoloFormMultipartStringCardsList: React.FC<
  FieldProps & MultipartStringCardsProps
> = ({ field, form: { errors, setFieldValue }, ...rest }) => {
  return (
    <React.Fragment>
      <MultipartStringCardsList
        {...field}
        {...rest}
        values={rest.values || field.value}
        valueDeleted={indexDeleted => {
          setFieldValue(field.name, [...field.value].splice(indexDeleted, 1));
        }}
        createNew={newPair => {
          let newList = !!rest.values ? [...rest.values] : [...field.value];
          newList.push({
            value: newPair.newValue,
            name: newPair.newName
          });
          setFieldValue(field.name, newList);
        }}
      />
      <ErrorText errorExists={!!errors}>{errors[field.name]}</ErrorText>
    </React.Fragment>
  );
};

export const SoloSecretRefInput: React.FC<
  FieldProps &
    TypeaheadProps & {
      type: string;
      asColumn: boolean;
      parentForm: FormikProps<any>;
    }
> = ({ parentForm, field: parentField, form, type, asColumn, ...rest }) => {
  const namespaces = React.useContext(NamespacesContext);
  const [selectedNS, setSelectedNS] = React.useState('');
  const listSecretsRequest = new ListSecretsRequest();
  const [noSecrets, setNoSecrets] = React.useState(false);
  React.useEffect(() => {
    listSecretsRequest.setNamespacesList(namespaces);
  }, [namespaces]);

  const { data: secretsListData } = useListSecrets(listSecretsRequest);

  const [secretsFound, setSecretsFound] = React.useState(
    secretsListData
      ? secretsListData.secretsList
          .filter(secret => {
            // TODO: are these the only forms requiring a secret ref?
            if (type === 'aws') return !!secret.aws;
            if (type === 'azure') return !!secret.azure;
          })
          .filter(secret => secret.metadata!.namespace === selectedNS)
          .map(secret => secret.metadata!.name)
      : []
  );

  React.useEffect(() => {
    setSecretsFound(
      secretsListData
        ? secretsListData.secretsList
            .filter(secret => {
              if (type === 'aws') return !!secret.aws;
              if (type === 'azure') return !!secret.azure;
            })
            .filter(secret => secret.metadata!.namespace === selectedNS)

            .map(secret => secret.metadata!.name)
        : []
    );
    if (secretsListData && secretsFound.length === 0) {
      setNoSecrets(true);
    }
  }, [selectedNS]);

  return (
    <React.Fragment>
      <Field
        name={`${parentField.name}.namespace`}
        render={({ form, field }: FieldProps) => (
          <div>
            <SoloTypeahead
              {...field}
              title='Secret Ref Namespace'
              defaultValue='gloo-system'
              presetOptions={namespaces}
              onChange={value => {
                form.setFieldValue(field.name, value);
                setSelectedNS(value);
                form.setFieldTouched(`${parentField.name}.name`);
                form.setFieldValue(`${parentField.name}.name`, '');
                if (secretsFound.length === 0) {
                  setNoSecrets(true);
                  parentForm.setFieldError(field.name, 'No secrets found');
                }
              }}
            />
            <ErrorText errorExists={noSecrets}>
              {form.errors[field.name]}
            </ErrorText>
          </div>
        )}
      />

      <Field
        name={`${parentField.name}.name`}
        render={({ form, field }: FieldProps) => (
          <div>
            <SoloTypeahead
              {...field}
              title='Secret Ref Name'
              disabled={secretsFound.length === 0}
              presetOptions={secretsFound}
              defaultValue='Secret...'
              onChange={value => {
                form.setFieldValue(`${parentField.name}.name`, value);
              }}
            />
            {form.errors && (
              <ErrorText>{form.errors[`${parentField.name}.name`]}</ErrorText>
            )}
          </div>
        )}
      />
    </React.Fragment>
  );
};

export const TableFormWrapper: React.FC = props => {
  return (
    <React.Fragment>
      {React.Children.map(props.children, child => (
        <td>{child}</td>
      ))}
    </React.Fragment>
  );
};

export const SoloFormStringsList: React.FC<
  FieldProps & StringCardsListProps
> = ({
  field,
  form: { errors, setFieldValue, ...form },
  createNewPromptText,
  ...rest
}) => {
  const removeValue = (index: number) => {
    setFieldValue(field.name, form.values[field.name].splice(index, 1));
  };
  const addValue = (value: string) => {
    setFieldValue(field.name, form.values[field.name].concat(value));
  };

  return (
    <React.Fragment>
      <StringCardsList
        {...rest}
        {...field}
        values={form.values[field.name].slice(1)}
        valueDeleted={removeValue}
        createNew={addValue}
        createNewPromptText={createNewPromptText}
      />
      <ErrorText errorExists={!!errors}>{errors[field.name]}</ErrorText>
    </React.Fragment>
  );
};
