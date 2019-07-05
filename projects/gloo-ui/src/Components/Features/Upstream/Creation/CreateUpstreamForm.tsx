import * as React from 'react';
/** @jsx jsx */
import { jsx } from '@emotion/core';

import styled from '@emotion/styled/macro';
import { colors } from 'Styles';
import { SoloInput } from 'Components/Common/SoloInput';
import { useFormValidation } from 'Hooks/useFormValidation';

import { ErrorText } from '../../VirtualService/Details/ExtAuthForm';
import { SoloTypeahead } from 'Components/Common/SoloTypeahead';
import { SoloButton } from 'Components/Common/SoloButton';
import {
  CreateUpstreamRequest,
  UpstreamInput
} from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/upstream_pb';
import { ResourceRef } from 'proto/github.com/solo-io/solo-kit/api/v1/ref_pb';
import { UpstreamSpec as KubeUpstreamSpec } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins/kubernetes/kubernetes_pb';
import { UpstreamSpec as AwsUpstreamSpec } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins/aws/aws_pb';
import { UpstreamSpec as AzureUpstreamSpec } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins/azure/azure_pb';
import { UpstreamSpec as StaticUpstreamSpec } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins/static/static_pb';

import { useCreateUpstream } from 'Api';
import { NamespacesContext } from 'GlooIApp';
import { SoloDropdown } from 'Components/Common/SoloDropdown';
import { UPSTREAM_TYPES } from 'utils/helpers';

interface Props {}

let initialValues = {
  name: '',
  type: '',
  namespace: ''
};

const FormContainer = styled.div`
  display: flex;
  flex-direction: column;
`;

const InputContainer = styled.div`
  display: flex;
  flex-direction: row;
  justify-content: space-around;
  padding: 15px 15px 0;
`;

const Footer = styled.div`
  display: flex;
  flex-direction: row;
  justify-content: flex-end;
`;

const validate = (values: typeof initialValues) => {
  let errors = {} as typeof initialValues;
  if (!values.name) {
    errors.name = 'Name is required';
  }
  if (!values.namespace) {
    errors.namespace = 'Namespace is required';
  }
  if (!values.type) {
    errors.type = 'Type is required';
  }
  return errors;
};

export const CreateUpstreamForm = (props: Props) => {
  const {
    handleSubmit,
    handleChange,
    handleBlur,
    values,
    errors,
    isSubmitting,
    isDifferent
  } = useFormValidation(initialValues, validate, createUpstream);

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
      case 'AWS':
        let awsSpec = new AwsUpstreamSpec();
        usInput.setAws(awsSpec);
        break;
      case 'Azure':
        let azureSpec = new AzureUpstreamSpec();
        usInput.setAzure(azureSpec);
        break;
      case 'Kubernetes':
        let kubeSpec = new KubeUpstreamSpec();
        usInput.setKube(kubeSpec);
        break;
      case 'Static':
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
    <FormContainer>
      <InputContainer>
        <div>
          <SoloInput
            title='Upstream Name'
            name='name'
            value={values.name}
            placeholder={'Upstream Name'}
            onChange={handleChange}
            onBlur={handleBlur}
          />
          {errors && <ErrorText>{errors.name}</ErrorText>}
        </div>

        <div>
          <SoloDropdown
            title={'Upstream Type'}
            placeholder='Type'
            value={values.type}
            options={UPSTREAM_TYPES}
            onChange={e => handleChange(e, 'type')}
            onBlur={handleBlur}
          />
          {errors && <ErrorText>{errors.type}</ErrorText>}
        </div>
        <div>
          <SoloTypeahead
            title='Upstream Namespace'
            defaultValue={values.namespace}
            onChange={e => handleChange(e, 'namespace')}
            presetOptions={namespaces}
          />
          {errors && <ErrorText>{errors.namespace}</ErrorText>}
        </div>
      </InputContainer>
      <Footer>
        <SoloButton
          onClick={handleSubmit}
          text='Create Upstream'
          disabled={isSubmitting}
        />
      </Footer>
    </FormContainer>
  );
};
