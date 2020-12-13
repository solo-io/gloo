import {
  SoloFormInput,
  SoloFormSecretRefInput
} from 'Components/Common/Form/SoloFormField';
import {
  InputRow,
  SoloFormTemplate
} from 'Components/Common/Form/SoloFormTemplate';
import { ResourceRef } from 'proto/github.com/solo-io/solo-kit/api/v1/ref_pb';
import * as React from 'react';
import * as yup from 'yup';

export interface AzureValuesType {
  azureFunctionAppName: string;
  azureSecretRef: ResourceRef.AsObject;
}

export const azureInitialValues: AzureValuesType = {
  azureFunctionAppName: '',
  azureSecretRef: {
    name: '',
    namespace: 'gloo-system'
  }
};

interface Props {}

export const azureValidationSchema = yup.object().shape({
  azureFunctionAppName: yup.string(),
  azureSecretRefNamespace: yup.string(),
  azureSecretRefName: yup.string()
});

export const AzureUpstreamForm: React.FC<Props> = () => {
  return (
    <SoloFormTemplate formHeader='Azure Upstream Settings'>
      <InputRow>
        <div>
          <SoloFormInput
            name='azureFunctionAppName'
            title='Function App Name'
            placeholder='Function App Name'
          />
        </div>
        <div>
          <SoloFormSecretRefInput name='azureSecretRef' type='azure' />
        </div>
      </InputRow>
    </SoloFormTemplate>
  );
};
