import {
  SoloFormCheckbox,
  SoloFormInput
} from 'Components/Common/Form/SoloFormField';
import { SoloFormTemplate } from 'Components/Common/Form/SoloFormTemplate';
import { SoloButton } from 'Components/Common/SoloButton';
import { Field, FieldArray, FieldArrayRenderProps } from 'formik';
import * as React from 'react';
import * as yup from 'yup';
import { StringCard } from 'Components/Common/StringCardsList';

// TODO combine with main initial values
export const staticInitialValues = {
  staticServiceName: '',
  staticHostList: [
    {
      addr: '',
      port: ''
    }
  ],
  staticUseTls: false,
  staticServicePort: ''
};

interface Props {}

// TODO: figure out which fields are required
export const staticValidationSchema = yup.object().shape({
  staticHostList: yup.string(),
  staticServicePort: yup.string(),
  staticServiceName: yup.string()
});

export const StaticUpstreamForm: React.FC<Props> = () => {
  return (
    <SoloFormTemplate formHeader='Static Upstream Settings'>
      <Field
        name='staticServiceName'
        title='Service Name'
        placeholder='Service Name'
        component={SoloFormInput}
      />
      <Field name='staticUseTls' title='Use Tls' component={SoloFormCheckbox} />
      <Field
        name='staticServicePort'
        title='Service Port'
        placeholder='Service Port'
        component={SoloFormInput}
      />
      <FieldArray name='staticHostList' render={SoloFormArrayField} />
    </SoloFormTemplate>
  );
};

// TODO: make this component generic once we know how to display input values
export const SoloFormArrayField: React.FC<FieldArrayRenderProps> = ({
  form,
  remove,
  insert,
  name
}) => (
  <div>
    <div key={`${name}-${0}`}>
      <Field
        name={`kubeHostList[${0}].addr`}
        title='Host Address'
        placeholder='Host Address'
        component={SoloFormInput}
      />
      <Field
        name={`kubeHostList[${0}].port`}
        title='Host Port'
        placeholder='Host Port'
        component={SoloFormInput}
      />
      <SoloButton
        text='add to list'
        onClick={() => insert(0, { addr: '', port: '' })}
      />
    </div>
    {form.values.kubeHostList.map((host: any, index: any) => {
      return (
        <StringCard key={`${name}-${index}`}>
          {` ${host.addr || ''}:${host.port || ''}`}
          <div onClick={() => remove(index)}>x</div>
        </StringCard>
      );
    })}
  </div>
);
