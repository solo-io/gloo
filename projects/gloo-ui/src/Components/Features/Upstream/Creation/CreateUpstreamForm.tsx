import styled from '@emotion/styled/macro';
import { useCreateUpstream } from 'Api';
import {
  SoloFormDropdown,
  SoloFormInput,
  SoloFormTypeahead
} from 'Components/Common/Form/SoloFormField';
import {
  SoloFormTemplate,
  InputRow
} from 'Components/Common/Form/SoloFormTemplate';
import { Field, Formik } from 'formik';
import { NamespacesContext } from 'GlooIApp';
import {
  UpstreamSpec as AwsUpstreamSpec,
  LambdaFunctionSpec
} from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins/aws/aws_pb';
import { UpstreamSpec as AzureUpstreamSpec } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins/azure/azure_pb';
import { UpstreamSpec as ConsulUpstreamSpec } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins/consul/consul_pb';
import { UpstreamSpec as KubeUpstreamSpec } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins/kubernetes/kubernetes_pb';
import { ServiceSpec } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins/service_spec_pb';
import { UpstreamSpec as StaticUpstreamSpec } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins/static/static_pb';
import { ResourceRef } from 'proto/github.com/solo-io/solo-kit/api/v1/ref_pb';
import {
  CreateUpstreamRequest,
  UpstreamInput
} from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/upstream_pb';
import * as React from 'react';
import { UPSTREAM_SPEC_TYPES, UPSTREAM_TYPES } from 'utils/upstreamHelpers';
import * as yup from 'yup';
import { awsInitialValues, AwsUpstreamForm } from './AwsUpstreamForm';
import { azureInitialValues, AzureUpstreamForm } from './AzureUpstreamForm';
import { consulInitialValues, ConsulUpstreamForm } from './ConsulUpstreamForm';
import { kubeInitialValues, KubeUpstreamForm } from './KubeUpstreamForm';
import { staticInitialValues, StaticUpstreamForm } from './StaticUpstreamForm';

interface Props {}

const FormContainer = styled.div`
  display: flex;
  flex-direction: column;
`;

// TODO: better way to include all initial values?
export const initialValues = {
  name: '',
  type: '',
  namespace: 'gloo-system',
  ...awsInitialValues,
  ...kubeInitialValues,
  ...staticInitialValues,
  ...azureInitialValues,
  ...consulInitialValues
};

// TODO combine validation schemas
const validationSchema = yup.object().shape({
  name: yup.string(),
  namespace: yup.string(),
  type: yup.string()
});

export const CreateUpstreamForm = (props: Props) => {
  const namespaces = React.useContext(NamespacesContext);

  const { refetch: makeRequest } = useCreateUpstream(null);

  const LambdaFunctionList = React.useRef<LambdaFunctionSpec[]>([]);
  // grpc request
  function createUpstream(values: typeof initialValues) {
    const newUpstreamReq = new CreateUpstreamRequest();
    const usInput = new UpstreamInput();

    const usResourceRef = new ResourceRef();

    usResourceRef.setName(values.name);
    usResourceRef.setNamespace(values.namespace);

    usInput.setRef(usResourceRef);

    //TODO: set up correct upstream spec
    // TODO: validation for specific fields
    switch (values.type) {
      case UPSTREAM_SPEC_TYPES.AWS:
        const awsSpec = new AwsUpstreamSpec();
        awsSpec.setRegion(values.awsRegion);
        const awsSecretRef = new ResourceRef();
        awsSecretRef.setName(values.awsSecretRef.name);
        awsSecretRef.setNamespace(values.awsSecretRef.namespace);
        awsSpec.setSecretRef(awsSecretRef);
        values.awsLambdaFunctionsList.forEach(lambda => {
          const newLambdaList = new LambdaFunctionSpec();
          newLambdaList.setLambdaFunctionName(lambda.lambdaFunctionName);
          newLambdaList.setLogicalName(lambda.logicalName);
          newLambdaList.setQualifier(lambda.qualifier);
          LambdaFunctionList.current.push(newLambdaList);
        });
        console.log(values.awsLambdaFunctionsList);
        awsSpec.setLambdaFunctionsList(LambdaFunctionList.current);

        usInput.setAws(awsSpec);
        break;
      case UPSTREAM_SPEC_TYPES.AZURE:
        const azureSpec = new AzureUpstreamSpec();
        const azureSecretRef = new ResourceRef();
        azureSecretRef.setName(values.azureSecretRef.name);
        azureSecretRef.setNamespace(values.azureSecretRef.namespace);
        azureSpec.setSecretRef(azureSecretRef);
        const azureFnList = values.azureFunctionsList.map(azureFn => {
          const azureFnSpec = new AzureUpstreamSpec.FunctionSpec();
          azureFnSpec.setFunctionName(azureFn.functionName);
          azureFnSpec.setAuthLevel(azureFn.authLevel);
          return azureFnSpec;
        });
        azureSpec.setFunctionsList(azureFnList);
        azureSpec.setFunctionAppName(values.azureFunctionAppName);
        usInput.setAzure(azureSpec);
        break;
      case UPSTREAM_SPEC_TYPES.KUBE:
        const kubeSpec = new KubeUpstreamSpec();
        kubeSpec.setServiceName(values.kubeServiceName);
        kubeSpec.setServiceNamespace(values.kubeServiceNamespace);
        kubeSpec.setServicePort(values.kubeServicePort);
        usInput.setKube(kubeSpec);
        break;
      case UPSTREAM_SPEC_TYPES.STATIC:
        const staticSpec = new StaticUpstreamSpec();
        staticSpec.setUseTls(values.staticUseTls);
        usInput.setStatic(staticSpec);
        break;
      case UPSTREAM_SPEC_TYPES.CONSUL:
        const consulSpec = new ConsulUpstreamSpec();
        consulSpec.setServiceName(values.consulServiceName);
        consulSpec.setServiceTagsList(values.consulServiceTagsList);
        consulSpec.setConnectEnabled(values.consulConnectEnabled);
        consulSpec.setDataCentersList([values.consulDataCenter]);
        const consulServiceSpec = new ServiceSpec();
        consulSpec.setServiceSpec(consulServiceSpec);
        usInput.setConsul(consulSpec);
      default:
        break;
    }

    newUpstreamReq.setInput(usInput);
    makeRequest(newUpstreamReq);
  }

  return (
    <Formik
      initialValues={initialValues}
      validationSchema={validationSchema}
      onSubmit={createUpstream}>
      {formik => (
        <FormContainer>
          <SoloFormTemplate>
            <InputRow>
              <div>
                <Field
                  name='name'
                  title='Upstream Name'
                  placeholder='Upstream Name'
                  component={SoloFormInput}
                />
              </div>
              <div>
                <Field
                  name='type'
                  title='Upstream Type'
                  placeholder='Type'
                  options={UPSTREAM_TYPES}
                  component={SoloFormDropdown}
                />
              </div>
              <div>
                <Field
                  name='namespace'
                  title='Upstream Namespace'
                  defaultValue='gloo-system'
                  presetOptions={namespaces}
                  component={SoloFormTypeahead}
                />
              </div>
            </InputRow>
          </SoloFormTemplate>
          {formik.values.type === UPSTREAM_SPEC_TYPES.AWS && (
            <AwsUpstreamForm parentForm={formik} />
          )}
          {formik.values.type === UPSTREAM_SPEC_TYPES.KUBE && (
            <KubeUpstreamForm parentForm={formik} />
          )}
          {formik.values.type === UPSTREAM_SPEC_TYPES.STATIC && (
            <StaticUpstreamForm parentForm={formik} />
          )}
          {formik.values.type === UPSTREAM_SPEC_TYPES.AZURE && (
            <AzureUpstreamForm parentForm={formik} />
          )}
          {formik.values.type === UPSTREAM_SPEC_TYPES.CONSUL && (
            <ConsulUpstreamForm parentForm={formik} />
          )}
        </FormContainer>
      )}
    </Formik>
  );
};
