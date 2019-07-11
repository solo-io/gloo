import {
  SoloFormCheckbox,
  SoloFormInput
} from 'Components/Common/Form/SoloFormField';
import { SoloFormTemplate } from 'Components/Common/Form/SoloFormTemplate';
import { Field } from 'formik';
import * as React from 'react';
import * as yup from 'yup';

// TODO combine with main initial values
export const consulInitialValues = {
  consulServiceName: '',
  // TODO: decide on best way to display lists
  consulServiceTagsList: '',
  consulConnectEnabled: false,
  consulDataCenter: ''
};

interface Props {}

export const consulValidationSchema = yup.object().shape({
  consulServiceName: yup.string(),
  consulServiceTagsList: yup.string(),
  consulConnectEnabled: yup.boolean(),
  consulDataCenter: yup.string()
});

export const ConsulUpstreamForm: React.FC<Props> = () => {
  return (
    <SoloFormTemplate formHeader='Consul Upstream Settings'>
      <Field
        name='consulServiceName'
        title='Service Name'
        placeholder='Service Name'
        component={SoloFormInput}
      />
      <Field
        name='consulConnectEnabled'
        title='Enable Consul Connect'
        component={SoloFormCheckbox}
      />
      <Field
        name='consulDataCenter'
        title='Data Center'
        placeholder='Data Center'
        component={SoloFormInput}
      />
    </SoloFormTemplate>
  );
};
