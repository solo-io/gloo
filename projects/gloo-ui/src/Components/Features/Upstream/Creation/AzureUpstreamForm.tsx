import {
  SoloFormInput,
  SoloFormTypeahead
} from 'Components/Common/Form/SoloFormField';
import {
  SoloFormTemplate,
  InputRow
} from 'Components/Common/Form/SoloFormTemplate';
import { Field, FieldProps, FieldArray, FieldArrayRenderProps } from 'formik';
import { NamespacesContext } from 'GlooIApp';
import * as React from 'react';
import * as yup from 'yup';
import { ListSecretsRequest } from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/secret_pb';
import { useListSecrets } from 'Api';
import { SoloTypeahead } from 'Components/Common/SoloTypeahead';
import { ErrorText } from 'Components/Features/VirtualService/Details/ExtAuthForm';
import { ReactComponent as CloseX } from 'assets/close-x.svg';
import { ReactComponent as GreenPlus } from 'assets/small-green-plus.svg';
import { SoloFormDropdown } from '../../../Common/Form/SoloFormField';
import { AZURE_AUTH_LEVELS } from 'utils/azureHelpers';
import { UpstreamSpec as AzureUpstreamSpec } from '../../../../proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins/azure/azure_pb';
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
  azureSecretRefNamespace: '',
  azureSecretRefName: '',
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
  const namespaces = React.useContext(NamespacesContext);
  const [selectedNS, setSelectedNS] = React.useState('');
  const listSecretsReq = new ListSecretsRequest();

  listSecretsReq.setNamespacesList(namespaces);

  const { data: secretsListData } = useListSecrets(listSecretsReq);

  const [secretsFound, setSecretsFound] = React.useState<string[]>(() =>
    secretsListData
      ? secretsListData.secretsList
          .filter(secret => !!secret.azure && secret.azure)
          .map(secret => secret.metadata!.name)
      : []
  );

  React.useEffect(() => {
    setSecretsFound(
      secretsListData
        ? secretsListData.secretsList
            .filter(
              secret =>
                !!secret.azure &&
                secret.azure &&
                secret.metadata!.namespace === selectedNS
            )
            .map(secret => secret.metadata!.name)
        : []
    );
  }, [selectedNS]);

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
            name='azureSecretRefNamespace'
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
                    form.setFieldValue('azureSecretRefName', '');
                  }}
                />
                {form.errors && (
                  <ErrorText>{form.errors[field.name]}</ErrorText>
                )}
              </div>
            )}
          />
        </div>
        <div>
          <Field
            name='azureSecretRefName'
            render={({ form, field }: FieldProps) => (
              <div>
                <SoloTypeahead
                  {...field}
                  title='Secret Ref Name'
                  disabled={secretsFound.length === 0}
                  presetOptions={secretsFound}
                  defaultValue='Secret...'
                  onChange={value => form.setFieldValue(field.name, value)}
                />
                {form.errors && (
                  <ErrorText>{form.errors[field.name]}</ErrorText>
                )}
              </div>
            )}
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
