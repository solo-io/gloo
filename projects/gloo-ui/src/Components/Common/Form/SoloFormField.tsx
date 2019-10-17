import styled from '@emotion/styled';
import { Select } from 'antd';
import { useField, useFormikContext } from 'formik';
import _ from 'lodash';
import { VirtualService } from 'proto/github.com/solo-io/gloo/projects/gateway/api/v1/virtual_service_pb';
import React from 'react';
import { shallowEqual, useSelector } from 'react-redux';
import { AppState } from 'store';
import { colors } from 'Styles';
import {
  createUpstreamId,
  getIcon,
  getUpstreamType,
  groupBy,
  parseUpstreamId
} from 'utils/helpers';
import {
  MultipartStringCardsList,
  MultipartStringCardsProps
} from '../MultipartStringCardsList';
import { SoloCheckbox } from '../SoloCheckbox';
import {
  DropdownProps,
  SoloDropdown,
  SoloDropdownBlock
} from '../SoloDropdown';
import { DurationProps, SoloDurationEditor } from '../SoloDurationEditor';
import { Label, SoloInput } from '../SoloInput';
import { SoloMultiSelect } from '../SoloMultiSelect';
import { SoloTypeahead, TypeaheadProps } from '../SoloTypeahead';
import { StringCardsList } from '../StringCardsList';
import { RouteTable } from 'proto/github.com/solo-io/gloo/projects/gateway/api/v1/route_table_pb';
import { RouteParent } from 'Components/Features/VirtualService/RouteTableDetails';

const { Option, OptGroup } = Select;

type ErrorTextProps = { errorExists?: boolean };
export const ErrorText = styled.div`
  color: ${colors.grapefruitOrange};
  visibility: ${(props: ErrorTextProps) =>
    props.errorExists ? 'visible' : 'hidden'};
  min-height: 19px;
`;

export const SoloFormInput = ({ ...props }) => {
  const [field, meta] = useField(props.name);

  return (
    <>
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
    </>
  );
};

interface FormDurationProps extends DurationProps {
  name: string;
}
export const SoloFormDurationEditor: React.FC<FormDurationProps> = ({
  ...props
}) => {
  const [field, meta] = useField(props.name);
  const form = useFormikContext<any>();

  return (
    <>
      <SoloDurationEditor
        {...field}
        {...props}
        error={!!meta.error && meta.touched}
        title={props.title}
        value={field.value}
        onChange={newDurationValues =>
          form.setFieldValue(field.name, newDurationValues)
        }
        onBlur={newDurationValues =>
          form.setFieldValue(field.name, newDurationValues)
        }
      />
      <ErrorText errorExists={!!meta.error && meta.touched}>
        {meta.error}
      </ErrorText>
    </>
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
    <>
      <SoloTypeahead
        {...props}
        {...field}
        title={props.title}
        presetOptions={props.presetOptions}
        onChange={value => form.setFieldValue(props.name, value)}
      />
      <ErrorText errorExists={!!meta.error}>{meta.error}</ErrorText>
    </>
  );
};

export const SoloFormDropdown = (props: any) => {
  const [field, meta] = useField(props.name);
  const form = useFormikContext<any>();
  return (
    <>
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
    </>
  );
};

export const SoloFormCheckbox = (props: any) => {
  const [field, meta] = useField(props.name);
  const form = useFormikContext<any>();
  return (
    <>
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
    </>
  );
};

export const SoloFormMultiselect = (props: any) => {
  const [field, meta] = useField(props.name);
  const form = useFormikContext<any>();
  return (
    <>
      <SoloMultiSelect
        {...field}
        {...props}
        error={!!meta.error && meta.touched}
        values={field.value}
        onChange={(value: string[]) => form.setFieldValue(field.name, value)}
      />
      <ErrorText errorExists={!!meta.touched && !!meta.error}>
        {meta.error}
      </ErrorText>{' '}
    </>
  );
};

interface MetadataBasedDropdownProps extends DropdownProps {
  value: any;
  options: any[];
  name: string;
  onChange?: (newValue: any) => void;
}
export const SoloFormMetadataBasedDropdown: React.FC<
  MetadataBasedDropdownProps
> = React.memo(props => {
  const [field, meta] = useField(props.name);
  const form = useFormikContext<any>();
  const usedOptions = React.useMemo(
    () =>
      props.options
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
        }),
    [props.options.length]
  );

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
    <>
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
    </>
  );
}, shallowEqual);

interface VirtualServiceTypeaheadProps extends TypeaheadProps {
  value: RouteParent | undefined;
  options: RouteParent[];
  name: string;
  onChange?: (newValue: any) => any;
}
export const SoloFormVirtualServiceTypeahead: React.FC<
  VirtualServiceTypeaheadProps
> = ({ ...props }) => {
  const {
    config: { namespace: podNamespace }
  } = useSelector((state: AppState) => state);
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
        namespace: podNamespace
      };
    }

    if (props.onChange) {
      props.onChange(optionChosen);
    }

    form.setFieldValue(field.name, optionChosen);
  };

  return (
    <>
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
    </>
  );
};

export const SoloFormMultipartStringCardsList: React.FC<
  { name: string } & MultipartStringCardsProps
> = ({ name, ...props }) => {
  const [field, meta] = useField(name);
  const form = useFormikContext<any>();

  return (
    <>
      <MultipartStringCardsList
        {...field}
        {...props}
        values={field.value}
        valueDeleted={indexDeleted => {
          const newArr = [...field.value];
          newArr.splice(indexDeleted, 1);

          form.setFieldValue(field.name, newArr);
        }}
        createNew={newPair => {
          let newList = [...field.value];
          const newObj: any = {};
          newObj[props.nameSlotTitle || 'name'] = newPair.newName;
          newObj[props.valueSlotTitle || 'value'] = newPair.newValue;
          if (newPair.newBool !== undefined) {
            newObj[props.boolSlotTitle || 'regex'] = newPair.newBool;
          }

          newList.push(newObj);
          form.setFieldValue(field.name, newList);
        }}
      />
      <ErrorText errorExists={!!meta.error && meta.touched}>
        {meta.error}
      </ErrorText>
    </>
  );
};

export const SoloFormSecretRefInput: React.FC<{
  type: string;
  asColumn?: boolean;
  name: string;
}> = props => {
  const {
    config: { namespacesList }
  } = useSelector((state: AppState) => state);
  const secretsList = useSelector(
    (state: AppState) => state.secrets.secretsList
  );
  const { name, type } = props;
  const [field, meta] = useField(name);
  const form = useFormikContext<any>();

  const [namespaceField, namespaceMeta] = useField(`${field.name}.namespace`);
  const [nameField, nameMeta] = useField(`${field.name}.name`);

  const [selectedNS, setSelectedNS] = React.useState(namespaceField.value);
  const [noSecrets, setNoSecrets] = React.useState(false);

  const [secretsFound, setSecretsFound] = React.useState(
    secretsList.length
      ? secretsList
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
      secretsList
        ? secretsList
            .filter(secret => {
              if (type === 'aws') return !!secret.aws;
              if (type === 'azure') return !!secret.azure;
            })
            .filter(secret => secret.metadata!.namespace === selectedNS)

            .map(secret => secret.metadata!.name)
        : []
    );
    if (secretsList && secretsFound.length === 0) {
      setNoSecrets(true);
    }
  }, [selectedNS]);

  return (
    <>
      <div>
        <SoloTypeahead
          {...namespaceField}
          title='Secret Ref Namespace'
          presetOptions={namespacesList.map(ns => {
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
    </>
  );
};

const compareFn = (
  prevProps: Readonly<{ type: string; name: string }>,
  newProps: Readonly<{ type: string; name: string }>
) => {
  return _.isEqual(prevProps, newProps);
};
export const SoloAWSSecretsList: React.FC<{
  type: string;
  name: string;
}> = React.memo(props => {
  const secretsList = useSelector(
    (state: AppState) => state.secrets.secretsList
  );
  const { name, type } = props;
  const [field, meta] = useField(name);
  const form = useFormikContext<any>();

  const groupedSecrets = Array.from(
    groupBy(secretsList, secret => secret.metadata!.namespace).entries()
  );
  const awsSecretsList = groupedSecrets.filter(([ns, secrets]) =>
    secrets.filter(s => !!s.aws)
  );

  return (
    <div style={{ width: '100%' }}>
      {name && <Label>AWS Secret</Label>}
      <SoloDropdownBlock
        style={{ width: '200px' }}
        onChange={(value: any) => {
          let [name, namespace] = value.split('::');
          let selectedSecret = secretsList.find(
            s =>
              s.metadata!.name === name && s.metadata!.namespace === namespace
          );

          form.setFieldValue(
            `${field.name}.name`,
            selectedSecret!.metadata!.name
          );
          form.setFieldValue(
            `${field.name}.namespace`,
            selectedSecret!.metadata!.namespace
          );
        }}>
        {awsSecretsList.map(([namespace, secrets]) => {
          return (
            <OptGroup key={namespace} label={namespace}>
              {secrets.map(s => (
                <Option key={`${s.metadata!.name}::${s.metadata!.namespace}`}>
                  {s.metadata!.name}
                </Option>
              ))}
            </OptGroup>
          );
        })}
      </SoloDropdownBlock>
      <ErrorText errorExists={!!meta.error && meta.touched}>
        {meta.error}
      </ErrorText>
    </div>
  );
}, compareFn);

const compare = (
  prevProps: Readonly<{ name: string; title: string; defaultValue: string }>,
  newProps: Readonly<{ name: string; title: string; defaultValue: string }>
) => {
  return _.isEqual(prevProps, newProps);
};
type RouteParentType = RouteTable.AsObject | VirtualService.AsObject;
export const SoloRouteParentDropdown: React.FC<{
  name: string;
  title: string;
  defaultValue: string;
}> = React.memo(props => {
  const routeTablesList = useSelector((state: AppState) =>
    state.routeTables.routeTablesList.map(rtd => !!rtd && rtd.routeTable!)
  );
  const groupedRouteTables = React.useMemo(
    () => Array.from(groupBy(routeTablesList, () => 'Route Table')),
    [routeTablesList.length]
  );

  const virtualServicesList = useSelector((state: AppState) =>
    state.virtualServices.virtualServicesList.map(
      vsd => !!vsd && vsd.virtualService!
    )
  );
  const groupedVirtualServices = React.useMemo(
    () => Array.from(groupBy(virtualServicesList, () => 'Virtual Service')),
    [virtualServicesList.length]
  );

  const routeParents: [string, RouteParentType[]][] = [
    ...groupedRouteTables,
    ...groupedVirtualServices
  ];

  const { name, title } = props;
  const [field, meta] = useField(name);
  const form = useFormikContext<any>();

  return (
    <div style={{ width: '100%' }}>
      {title && <Label>{title}</Label>}
      <SoloDropdownBlock
        disabled={!!props.defaultValue}
        value={props.defaultValue}
        onChange={(value: any) => {

          let [name, namespace] = value.split('::');
          // find the correct one in the list
          let selectedRouteParent;
          form.setFieldValue(`${field.name}.name`, name);
          // set it in the form
          form.setFieldValue(`${field.name}.namespace`, namespace);
        }}>
        {routeParents.map(([namespace, secrets]) => {
          return (
            <OptGroup key={namespace} label={namespace}>
              {secrets.map(s => (
                <Option key={`${s.metadata!.name}::${s.metadata!.namespace}`}>
                  {s.metadata!.name}
                </Option>
              ))}
            </OptGroup>
          );
        })}
      </SoloDropdownBlock>
      <ErrorText errorExists={!!meta.error && meta.touched}>
        {meta.error}
      </ErrorText>
    </div>
  );
}, compare);

export const TableFormWrapper: React.FC = props => {
  return (
    <>
      {React.Children.map(props.children, child => (
        <td>{child}</td>
      ))}
    </>
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
    <>
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
    </>
  );
};
