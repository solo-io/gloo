import {
  SoloFormInput,
  SoloFormTypeahead
} from 'Components/Common/Form/SoloFormField';
import { SoloFormTemplate } from 'Components/Common/Form/SoloFormTemplate';
import { Field } from 'formik';
import { NamespacesContext } from 'GlooIApp';
import * as React from 'react';
import * as yup from 'yup';

// TODO combine with main initial values
export const kubeInitialValues = {
  serviceName: '',
  serviceNamespace: 'gloo-system',
  servicePort: ''
};

interface Props {}

// TODO: figure out which fields are required
export const kubeValidationSchema = yup.object().shape({
  serviceName: yup.string(),
  serviceNamespace: yup.string(),
  servicePort: yup.string()
});

export const KubeUpstreamForm: React.FC<Props> = () => {
  const namespaces = React.useContext(NamespacesContext);

  return (
    <SoloFormTemplate formHeader='Kube Upstream Settings'>
      <Field
        name='serviceName'
        title='Service Name'
        placeholder='Service Name'
        component={SoloFormInput}
      />
      <Field
        name='serviceNamespace'
        title='Service Namespace'
        defaultValue='gloo-system'
        presetOptions={namespaces}
        component={SoloFormTypeahead}
      />
      <Field
        name='servicePort'
        title='Service Port'
        placeholder='Service Port'
        component={SoloFormInput}
      />
    </SoloFormTemplate>
  );
};
