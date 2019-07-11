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
export const azureInitialValues = {
  azureFunctionAppName: '',
  azureSecretRefNamespace: '',
  azureSecretRefName: ''
};

interface Props {}

export const azureValidationSchema = yup.object().shape({
  azureFunctionAppName: yup.string(),
  azureSecretRefNamespace: yup.string(),
  azureSecretRefName: yup.string()
});

export const AzureUpstreamForm: React.FC<Props> = () => {
  const namespaces = React.useContext(NamespacesContext);

  return (
    <SoloFormTemplate formHeader='AWS Upstream Settings'>
      <Field
        name='azureFunctionAppName'
        title='Function App Name'
        placeholder='Function App Name'
        component={SoloFormInput}
      />
      <Field
        name='azureSecretRefNamespace'
        title='Secret Ref Namespace'
        presetOptions={namespaces}
        component={SoloFormTypeahead}
      />
      <Field
        name='azureSecretRefName'
        title='Secret Ref Name'
        placeholder='Secret Ref Name'
        component={SoloFormInput}
      />
    </SoloFormTemplate>
  );
};
