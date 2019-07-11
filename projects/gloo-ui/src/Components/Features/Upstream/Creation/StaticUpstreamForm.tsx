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
  hostList: [
    {
      addr: '',
      port: ''
    }
  ],
  useTls: false,
  servicePort: ''
};

interface Props {}

// TODO: figure out which fields are required
export const staticValidationSchema = yup.object().shape({
  serviceName: yup.string(),
  serviceNamespace: yup.string(),
  servicePort: yup.string()
});
// hostlist
//host ={addr:'', port:''}
export const StaticUpstreamForm: React.FC<Props> = () => {
  return (
    <SoloFormTemplate formHeader='Static Upstream Settings'>
      <Field
        name='serviceName'
        title='Service Name'
        placeholder='Service Name'
        component={SoloFormInput}
      />
      <Field name='useTls' title='Use Tls' component={SoloFormCheckbox} />
      <Field
        name='servicePort'
        title='Service Port'
        placeholder='Service Port'
        component={SoloFormInput}
      />
      <FieldArray name='hostList' render={SoloFormArrayField} />
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
        name={`hostList[${0}].addr`}
        title='Host Address'
        placeholder='Host Address'
        component={SoloFormInput}
      />
      <Field
        name={`hostList[${0}].port`}
        title='Host Port'
        placeholder='Host Port'
        component={SoloFormInput}
      />
      <SoloButton
        text='add to list'
        onClick={() => insert(0, { addr: '', port: '' })}
      />
    </div>
    {form.values.hostList.map((host: any, index: any) => {
      return (
        <div key={`${name}-${index}`}>
          {
            <StringCard>
              {` ${host.addr || ''}:${host.port || ''}`}
              <div onClick={() => remove(index)}>x</div>
            </StringCard>
          }
        </div>
      );
    })}
  </div>
);
