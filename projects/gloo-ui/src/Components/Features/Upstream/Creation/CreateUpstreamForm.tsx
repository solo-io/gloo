import styled from '@emotion/styled/macro';
import { useCreateUpstream } from 'Api';
import { SoloButton } from 'Components/Common/SoloButton';
import {
  SoloFormDropdown,
  SoloFormInput,
  SoloFormTypeahead
} from 'Components/Common/SoloFormField';
import { Field, Formik } from 'formik';
import { NamespacesContext } from 'GlooIApp';
import { UpstreamSpec as AwsUpstreamSpec } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins/aws/aws_pb';
import { UpstreamSpec as AzureUpstreamSpec } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins/azure/azure_pb';
import { UpstreamSpec as KubeUpstreamSpec } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins/kubernetes/kubernetes_pb';
import { UpstreamSpec as StaticUpstreamSpec } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins/static/static_pb';
import { ResourceRef } from 'proto/github.com/solo-io/solo-kit/api/v1/ref_pb';
import {
  CreateUpstreamRequest,
  UpstreamInput
} from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/upstream_pb';
import * as React from 'react';
import { UPSTREAM_SPEC_TYPES, UPSTREAM_TYPES } from 'utils/upstreamHelpers';
import * as yup from 'yup';
import { AwsUpstreamForm } from './AwsUpstreamForm';

interface Props {}

const FormContainer = styled.div`
  display: flex;
  flex-direction: column;
`;

export const InputContainer = styled.div`
  display: flex;
  flex-direction: row;
  justify-content: space-around;
  padding: 10px;
`;

const InputItem = styled.div`
  display: flex;
  flex-direction: column;
`;

const Footer = styled.div`
  display: flex;
  flex-direction: row;
  justify-content: flex-end;
`;

// TODO: combine initial values
let initialValues = {
  name: '',
  type: '',
  namespace: 'gloo-system'
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

  // grpc request
  function createUpstream(values: typeof initialValues) {
    let newUpstreamReq = new CreateUpstreamRequest();
    let usInput = new UpstreamInput();

    let usResourceRef = new ResourceRef();

    usResourceRef.setName(values.name);
    usResourceRef.setNamespace(values.namespace);

    usInput.setRef(usResourceRef);

    //TODO: set up correct upstream spec
    switch (values.type) {
      case UPSTREAM_SPEC_TYPES.AWS:
        let awsSpec = new AwsUpstreamSpec();
        usInput.setAws(awsSpec);
        break;
      case UPSTREAM_SPEC_TYPES.AZURE:
        let azureSpec = new AzureUpstreamSpec();
        usInput.setAzure(azureSpec);
        break;
      case UPSTREAM_SPEC_TYPES.KUBE:
        let kubeSpec = new KubeUpstreamSpec();
        usInput.setKube(kubeSpec);
        break;
      case UPSTREAM_SPEC_TYPES.STATIC:
        let staticSpec = new StaticUpstreamSpec();
        usInput.setStatic(staticSpec);
        break;
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
      {({ values, isSubmitting, handleSubmit }) => (
        <FormContainer>
          <InputContainer>
            <InputItem>
              <Field
                name='name'
                title='Upstream Name'
                placeholder='Upstream Name'
                component={SoloFormInput}
              />
            </InputItem>
            <InputItem>
              <Field
                name='type'
                title='Upstream Type'
                options={UPSTREAM_TYPES}
                component={SoloFormDropdown}
              />
            </InputItem>
            <InputItem>
              <Field
                name='namespace'
                title='Upstream Namespace'
                defaultValue='gloo-system'
                presetOptions={namespaces}
                component={SoloFormTypeahead}
              />
            </InputItem>
          </InputContainer>
          {values.type === UPSTREAM_SPEC_TYPES.AWS && <AwsUpstreamForm />}
          <Footer>
            <pre>{JSON.stringify(values, null, 2)}</pre>
            <SoloButton
              onClick={handleSubmit}
              text='Create Upstream'
              disabled={isSubmitting}
            />
          </Footer>
        </FormContainer>
      )}
    </Formik>
  );
};
