import { ReactComponent as CloseX } from 'assets/close-x.svg';
import { ReactComponent as GreenPlus } from 'assets/small-green-plus.svg';
import {
  SoloFormInput,
  SoloSecretRefInput
} from 'Components/Common/Form/SoloFormField';
import {
  InputRow,
  SoloFormTemplate
} from 'Components/Common/Form/SoloFormTemplate';
import { Field, FieldArray, FieldArrayRenderProps } from 'formik';
import * as React from 'react';
import { AZURE_AUTH_LEVELS } from 'utils/azureHelpers';
import * as yup from 'yup';
import { UpstreamSpec as AzureUpstreamSpec } from '../../../../proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins/azure/azure_pb';
import { SoloFormDropdown } from '../../../Common/Form/SoloFormField';
/* ------------------------------ Upstream Spec ----------------------------- */
/*
functionAppName: string,
secretRef?: ResourceRef :{ name:string, namespace: string},
functionsList: Array<UpstreamSpec.FunctionSpec: {
  functionName: string,
  authLevel: UpstreamSpec.FunctionSpec.AuthLevel: {ANONYMOUS = 0,FUNCTION = 1, ADMIN = 2,},
  
}>,
*/
// TODO combine with main initial values
export const azureInitialValues = {
  azureFunctionAppName: '',
  azureSecretRef: {
    name: '',
    namespace: 'gloo-system'
  },
  azureFunctionsList: [
    {
      functionName: '',
      authLevel: AzureUpstreamSpec.FunctionSpec.AuthLevel.FUNCTION
    }
  ]
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
        <Field
          name='azureFunctionAppName'
          title='Function App Name'
          placeholder='Function App Name'
          component={SoloFormInput}
        />

        <div>
          <Field
            name='azureSecretRef'
            type='azure'
            component={SoloSecretRefInput}
          />
        </div>
      </InputRow>
      <SoloFormTemplate formHeader='Azure Functions'>
        <FieldArray name='azureFunctionsList' render={AzureFunctions} />
      </SoloFormTemplate>
    </SoloFormTemplate>
  );
};

export const AzureFunctions: React.FC<FieldArrayRenderProps> = ({
  form,
  remove,
  insert,
  name
}) => (
  <React.Fragment>
    <InputRow>
      <div>
        <Field
          name='azureFunctionsList[0].functionName'
          title='Function Name'
          placeholder='Function Name'
          component={SoloFormInput}
        />
      </div>
      <div>
        <Field
          name='azureFunctionsList[0].authLevel'
          title='Function Auth Level'
          placeholder='Function Auth Level'
          options={AZURE_AUTH_LEVELS}
          component={SoloFormDropdown}
        />
      </div>
      <GreenPlus
        style={{ alignSelf: 'center' }}
        onClick={() =>
          insert(0, {
            functionName: '',
            authLevel: ''
          })
        }
      />
    </InputRow>
    <InputRow>
      {form.values.azureFunctionsList.map(
        (azureFn: AzureUpstreamSpec.FunctionSpec.AsObject, index: number) => {
          return (
            <div key={azureFn.functionName}>
              <div>{azureFn.functionName}</div>
              <div>{azureFn.authLevel}</div>
              <CloseX onClick={() => remove(index)} />
            </div>
          );
        }
      )}
    </InputRow>
  </React.Fragment>
);
