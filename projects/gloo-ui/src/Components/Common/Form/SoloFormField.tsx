import styled from '@emotion/styled/macro';
import { useListSecrets } from 'Api';
import { useField, useFormikContext } from 'formik';
import { NamespacesContext } from 'GlooIApp';
import { ListSecretsRequest } from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/secret_pb';
import React from 'react';
import { colors } from 'Styles';
import {
  createUpstreamId,
  getIcon,
  getUpstreamType,
  parseUpstreamId
} from 'utils/helpers';
import { MultipartStringCardsList } from '../MultipartStringCardsList';
import { SoloCheckbox } from '../SoloCheckbox';
import { DropdownProps, SoloDropdown, OptionType } from '../SoloDropdown';
import { SoloInput } from '../SoloInput';
import { SoloMultiSelect } from '../SoloMultiSelect';
import { SoloTypeahead, TypeaheadProps } from '../SoloTypeahead';
import { StringCardsList } from '../StringCardsList';
import { VirtualService } from 'proto/github.com/solo-io/gloo/projects/gateway/api/v1/virtual_service_pb';

export const ErrorText = styled<'div', { errorExists?: boolean }>('div')`
  color: ${colors.grapefruitOrange};
  visibility: ${props => (props.errorExists ? 'visible' : 'hidden')};
  min-height: 19px;
`;
// TODO: find best way to type the components

export const SoloFormInput = ({ ...props }) => {
  const [field, meta] = useField(props.name);

  return (
    <React.Fragment>
      <SoloInput
        {...field}
        {...props}
        error={!!meta.error && meta.touched}
        title={props.title}
        value={field.value}
        onChange={field.onChange}
      />
      <ErrorText errorExists={!!meta.error && meta.touched}>
        {meta.error}
      </ErrorText>
    </React.Fragment>
  );
};

interface FormTypeaheadProps extends TypeaheadProps {
  name: string;
}
export const SoloFormTypeahead: React.FC<FormTypeaheadProps> = ({
  ...props
}) => {
  const [field, meta] = useField(props.name);
  const form = useFormikContext<any>();
  return (
    <React.Fragment>
      <SoloTypeahead
        {...props}
        {...field}
        title={props.title}
        presetOptions={props.presetOptions}
        onChange={value => form.setFieldValue(props.name, value)}
      />
      <ErrorText errorExists={!!meta.error}>{meta.error}</ErrorText>
    </React.Fragment>
  );
};

export const SoloFormDropdown = (props: any) => {
  const [field, meta] = useField(props.name);
  const form = useFormikContext<any>();
  return (
    <React.Fragment>
      <SoloDropdown
        {...field}
        {...props}
        error={!!meta.error && meta.touched}
        onChange={value => form.setFieldValue(field.name, value)}
        onBlur={value => form.setFieldValue(field.name, value)}
      />
      <ErrorText errorExists={!!meta.error && meta.touched}>
        {meta.error}
      </ErrorText>
    </React.Fragment>
  );
};

export const SoloFormCheckbox = (props: any) => {
  const [field, meta] = useField(props.name);
  const form = useFormikContext<any>();
  return (
    <React.Fragment>
      <SoloCheckbox
        {...props}
        {...field}
        error={!!meta.error && meta.touched}
        checked={!!field.value}
        onChange={value => form.setFieldValue(field.name, value.target.checked)}
        label
      />
      <ErrorText errorExists={!!meta.touched && !!meta.error}>
        {meta.error}
      </ErrorText>
    </React.Fragment>
  );
};

export const SoloFormMultiselect = (props: any) => {
  const [field, meta] = useField(props.name);
  const form = useFormikContext<any>();
  return (
    <React.Fragment>
      <SoloMultiSelect
        {...field}
        {...props}
        error={!!meta.error && meta.touched}
        values={Object.keys(field.value).filter(key => field.value[key])}
        onChange={newValues => {
          const newFieldValues = { ...field.value };
          Object.keys(newFieldValues).forEach(key => {
            newFieldValues[key] = false;
          });
          for (let val of newValues) {
            newFieldValues[val] = true;
          }

          form.setFieldValue(field.name, newFieldValues);
        }}
      />
      <ErrorText errorExists={!!meta.touched && !!meta.error}>
        {meta.error}
      </ErrorText>{' '}
    </React.Fragment>
  );
};

interface MetadataBasedDropdownProps extends DropdownProps {
  value: any;
  options: any[];
  name: string;
  onChange?: (newValue: any) => any;
}
export const SoloFormMetadataBasedDropdown: React.FC<
  MetadataBasedDropdownProps
> = ({ ...props }) => {
  const [field, meta] = useField(props.name);
  const form = useFormikContext<any>();
  const usedOptions = props.options
    .sort((optionA, optionB) => {
      const nameA = optionA.metadata.name;
      const nameB = optionB.metadata.name;

      const nameOrder = nameA === nameB ? 0 : nameA < nameB ? -1 : 1;

      if (!!optionA.upstreamSpec) {
        const typeA = getUpstreamType(optionA);
        const typeB = getUpstreamType(optionB);

        return typeA < typeB ? -1 : typeA > typeB ? 1 : nameOrder;
      }

      return nameOrder;
    })
    .map(option => {
      return {
        key: createUpstreamId(option.metadata!), // the same as virtual service's currently
        displayValue: option.metadata.name,
        value: createUpstreamId(option.metadata!),
        icon: !!option.upstreamSpec
          ? getIcon(getUpstreamType(option))
          : undefined
      };
    });

  const usedValue =
    props.value && props.value.metadata
      ? props.value.metadata.name
      : field.value && field.value.metadata
      ? field.value.metadata.name
      : undefined;

  const setNewValue = (newValueId: any) => {
    const { name, namespace } = parseUpstreamId(newValueId);
    const optionChosen = props.options.find(
      option =>
        option.metadata.name === name && option.metadata.namespace === namespace
    );

    if (props.onChange) {
      props.onChange(optionChosen);
    }

    form.setFieldValue(field.name, optionChosen);
  };

  return (
    <React.Fragment>
      <SoloDropdown
        {...field}
        {...props}
        options={usedOptions}
        value={usedValue}
        onChange={setNewValue}
      />
      <ErrorText errorExists={!!meta.error && meta.touched}>
        {meta.error}
      </ErrorText>
    </React.Fragment>
  );
};

interface VirtualServiceTypeaheadProps extends TypeaheadProps {
  value: VirtualService.AsObject | undefined;
  options: VirtualService.AsObject[];
  name: string;
  onChange?: (newValue: any) => any;
}
export const SoloFormVirtualServiceTypeahead: React.FC<
  VirtualServiceTypeaheadProps
> = ({ ...props }) => {
  const [field, meta] = useField(props.name);
  const form = useFormikContext<any>();
  const usedOptions = props.options
    .sort((optionA, optionB) => {
      const nameA = optionA.metadata!.name;
      const nameB = optionB.metadata!.name;

      return nameA === nameB ? 0 : nameA < nameB ? -1 : 1;
    })
    .map(option => {
      return {
        key: createUpstreamId(option.metadata!), // the same as virtual service's currently
        displayValue: option.metadata!.name,
        value: createUpstreamId(option.metadata!)
      };
    });

  const usedValue =
    props.value && props.value.metadata
      ? props.value.metadata.name
      : field.value && field.value.metadata
      ? field.value.metadata.name
      : undefined;

  const setNewValue = (newValueId: any) => {
    const { name, namespace } = parseUpstreamId(newValueId);

    let optionChosen;
    if (!!namespace) {
      optionChosen = props.options.find(
        option =>
          option.metadata!.name === name &&
          option.metadata!.namespace === namespace
      );
    } else {
      const tempVirtualService = new VirtualService().toObject();

      // @ts-ignore
      tempVirtualService.metadata = {
        name: newValueId,
        namespace: 'gloo-system'
      };
    }
    console.log(optionChosen);

    if (props.onChange) {
      props.onChange(optionChosen);
    }

    form.setFieldValue(field.name, optionChosen);
  };

  return (
    <React.Fragment>
      <SoloTypeahead
        {...field}
        {...props}
        presetOptions={usedOptions}
        defaultValue={usedValue}
        onChange={setNewValue}
      />
      <ErrorText errorExists={!!meta.error && meta.touched}>
        {meta.error}
      </ErrorText>
    </React.Fragment>
  );
};

export const SoloFormMultipartStringCardsList: React.FC<
  { name: string } & any
> = ({ name, ...props }) => {
  const [field, meta] = useField(name);
  const form = useFormikContext<any>();
  return (
    <React.Fragment>
      <MultipartStringCardsList
        {...field}
        {...props}
        error={!!meta.error && meta.touched}
        values={field.value}
        valueDeleted={indexDeleted => {
          form.setFieldValue(
            field.name,
            [...field.value].splice(indexDeleted, 1)
          );
        }}
        createNew={newPair => {
          let newList = [...field.value];
          newList.push({
            value: newPair.newValue,
            name: newPair.newName
          });
          form.setFieldValue(field.name, newList);
        }}
      />
      <ErrorText errorExists={!!meta.error && meta.touched}>
        {meta.error}
      </ErrorText>
    </React.Fragment>
  );
};

export const SoloFormSecretRefInput: React.FC<{
  type: string;
  asColumn?: boolean;
  name: string;
}> = props => {
  const { name, type } = props;
  const [field, meta] = useField(name);
  const form = useFormikContext<any>();

  const [namespaceField, namespaceMeta] = useField(`${field.name}.namespace`);
  const [nameField, nameMeta] = useField(`${field.name}.name`);

  const namespaces = React.useContext(NamespacesContext);
  const [selectedNS, setSelectedNS] = React.useState(namespaceField.value);
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
      <div>
        <SoloTypeahead
          {...namespaceField}
          title='Secret Ref Namespace'
          presetOptions={namespaces.map(ns => {
            return { value: ns };
          })}
          onChange={value => {
            form.setFieldValue(`${field.name}.namespace`, value);

            setSelectedNS(value);
            if (noSecrets) {
              form.setFieldError(
                `${field.name}.namespace`,
                'No secrets found on this namespace'
              );
            }
            form.setFieldTouched(`${field.name}.name`);
            form.setFieldValue(`${field.name}.name`, '');
          }}
        />
        <ErrorText errorExists={!!namespaceMeta.error && namespaceMeta.touched}>
          {namespaceMeta.error}
        </ErrorText>
      </div>
      <div>
        <SoloTypeahead
          {...nameField}
          title='Secret Ref Name'
          disabled={secretsFound.length === 0}
          presetOptions={secretsFound.map(sF => {
            return { value: sF };
          })}
          defaultValue='Secret...'
          onChange={value => {
            form.setFieldValue(`${field.name}.name`, value);
          }}
        />
        <ErrorText errorExists={!!nameMeta.error && nameMeta.touched}>
          {nameMeta.error}
        </ErrorText>
      </div>
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

export const SoloFormStringsList: React.FC<any> = ({
  createNewPromptText,
  ...props
}) => {
  const [field, meta] = useField(props.name);
  const form = useFormikContext<any>();

  const removeValue = (index: number) => {
    form.setFieldValue(field.name, form.values[field.name].splice(index, 1));
  };
  const addValue = (value: string) => {
    form.setFieldValue(field.name, form.values[field.name].concat(value));
  };

  return (
    <React.Fragment>
      <StringCardsList
        {...props}
        {...field}
        error={!!meta.error && meta.touched}
        values={form.values[field.name].slice(1)}
        valueDeleted={removeValue}
        createNew={addValue}
        createNewPromptText={createNewPromptText}
      />
      <ErrorText errorExists={!!meta.error && meta.touched}>
        {meta.error}
      </ErrorText>
    </React.Fragment>
  );
};
