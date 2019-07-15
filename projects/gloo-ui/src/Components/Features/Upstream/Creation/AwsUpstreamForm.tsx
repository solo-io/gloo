import { useListSecrets } from 'Api';
import { ReactComponent as CloseX } from 'assets/close-x.svg';
import { ReactComponent as GreenPlus } from 'assets/small-green-plus.svg';
import {
  SoloFormInput,
  SoloFormTypeahead
} from 'Components/Common/Form/SoloFormField';
import {
  InputRow,
  SoloFormTemplate
} from 'Components/Common/Form/SoloFormTemplate';
import { SoloTypeahead } from 'Components/Common/SoloTypeahead';
import { ErrorText } from 'Components/Features/VirtualService/Details/ExtAuthForm';
import { Field, FieldArray, FieldArrayRenderProps, FieldProps } from 'formik';
import { NamespacesContext } from 'GlooIApp';
import { LambdaFunctionSpec } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins/aws/aws_pb';
import { ListSecretsRequest } from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/secret_pb';
import * as React from 'react';
import { AWS_REGIONS } from 'utils/upstreamHelpers';
import * as yup from 'yup';

/* ------------------------------ Upstream Spec ----------------------------- */
/*
  region: string,
  secretRef?: ResourceRef: {name: string, namespace: string},
  lambdaFunctionsList: Array<LambdaFunctionSpec: {logicalName: string,
    lambdaFunctionName: string,
    qualifier: string,}>,
*/

// TODO combine with main initial values
export const awsInitialValues = {
  awsRegion: 'us-east-1',
  awsSecretRefNamespace: 'gloo-system',
  awsSecretRefName: '',
  awsLambdaFunctionsList: [
    {
      logicalName: '',
      lambdaFunctionName: '',
      qualifier: ''
    }
  ]
};

interface Props {}

export const awsValidationSchema = yup.object().shape({
  awsRegion: yup.string(),
  awsSecretRefNamespace: yup.string(),
  awsSecretRefName: yup.string(),
  awsLambdaFunctionsList: yup.array().of(
    yup.object().shape({
      logicalName: yup.string(),
      lambdaFunctionName: yup.string(),
      qualifier: yup.string()
    })
  )
});

export const AwsUpstreamForm: React.FC<Props> = () => {
  const namespaces = React.useContext(NamespacesContext);

  const awsRegions = AWS_REGIONS.map(item => item.name);
  const [selectedNS, setSelectedNS] = React.useState('');
  const listSecretsReq = new ListSecretsRequest();

  listSecretsReq.setNamespacesList(namespaces);

  const { data: secretsListData } = useListSecrets(listSecretsReq);

  const [secretsFound, setSecretsFound] = React.useState<string[]>(() =>
    secretsListData
      ? secretsListData.secretsList
          .filter(secret => !!secret.aws && secret.aws)
          .map(secret => secret.metadata!.name)
      : []
  );

  React.useEffect(() => {
    setSecretsFound(
      secretsListData
        ? secretsListData.secretsList
            .filter(
              secret =>
                !!secret.aws &&
                secret.aws &&
                secret.metadata!.namespace === selectedNS
            )
            .map(secret => secret.metadata!.name)
        : []
    );
  }, [selectedNS]);

  return (
    <SoloFormTemplate formHeader='AWS Upstream Settings'>
      <InputRow>
        <div>
          <Field
            name='awsRegion'
            title='Region'
            presetOptions={awsRegions}
            component={SoloFormTypeahead}
          />
        </div>
        <div>
          <Field
            name='awsSecretRefNamespace'
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
                    form.setFieldValue('awsSecretRefName', '');
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
            name='awsSecretRefName'
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
      <SoloFormTemplate formHeader='Lambda Functions'>
        <FieldArray name='awsLambdaFunctionsList' render={LambdaFunctions} />
      </SoloFormTemplate>
    </SoloFormTemplate>
  );
};

export const LambdaFunctions: React.FC<FieldArrayRenderProps> = ({
  form,
  remove,
  insert,
  name
}) => (
  <React.Fragment>
    <InputRow>
      <Field
        name='awsLambdaFunctionsList[0].logicalName'
        title='Logical Name'
        placeholder='Logical Name'
        component={SoloFormInput}
      />
      <Field
        name='awsLambdaFunctionsList[0].lambdaFunctionName'
        title='Lambda Function Name'
        placeholder='Lambda Function Name'
        component={SoloFormInput}
      />
      <Field
        name='awsLambdaFunctionsList[0].qualifier'
        title='Qualifier'
        placeholder='Qualifier'
        component={SoloFormInput}
      />
      <GreenPlus
        style={{ alignSelf: 'center' }}
        onClick={() =>
          insert(0, {
            logicalName: '',
            lambdaFunctionName: '',
            qualifier: ''
          })
        }
      />
    </InputRow>
    <InputRow>
      {form.values.awsLambdaFunctionsList.map(
        (lambda: LambdaFunctionSpec.AsObject, index: number) => {
          return (
            <div key={lambda.logicalName}>
              <div>{lambda.logicalName}</div>
              <div>{lambda.lambdaFunctionName}</div>
              <div>{lambda.qualifier}</div>
              <CloseX onClick={() => remove(index)} />
            </div>
          );
        }
      )}
    </InputRow>
  </React.Fragment>
);
