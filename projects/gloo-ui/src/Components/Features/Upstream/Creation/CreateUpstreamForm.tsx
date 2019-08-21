import styled from '@emotion/styled/macro';
import {
  SoloFormDropdown,
  SoloFormInput,
  SoloFormTypeahead
} from 'Components/Common/Form/SoloFormField';
import {
  SoloFormTemplate,
  InputRow,
  Footer
} from 'Components/Common/Form/SoloFormTemplate';
import { Formik } from 'formik';

import * as React from 'react';
import { UPSTREAM_SPEC_TYPES, UPSTREAM_TYPES } from 'utils/upstreamHelpers';
import * as yup from 'yup';
import { awsInitialValues, AwsUpstreamForm } from './AwsUpstreamForm';
import { azureInitialValues, AzureUpstreamForm } from './AzureUpstreamForm';
import {
  consulInitialValues,
  ConsulUpstreamForm,
  consulValidationSchema
} from './ConsulUpstreamForm';
import { kubeInitialValues, KubeUpstreamForm } from './KubeUpstreamForm';
import { staticInitialValues, StaticUpstreamForm } from './StaticUpstreamForm';
import { SoloButton } from 'Components/Common/SoloButton';
import { withRouter } from 'react-router-dom';
import { RouteComponentProps } from 'react-router';

import { useDispatch, useSelector } from 'react-redux';
import { createUpstream } from 'store/upstreams/actions';
import { AppState } from 'store';
interface Props {
  toggleModal: React.Dispatch<React.SetStateAction<boolean>>;
}

const FormContainer = styled.div`
  display: flex;
  flex-direction: column;
`;

// TODO combine validation schemas
const validationSchema = yup.object().shape({
  name: yup
    .string()
    .required('Upstream name is required')
    .min(2, `Name can't be that short`),
  namespace: yup.string().required('Namespace is required'),
  type: yup.string().required('Must specify an upstream type'),
  awsRegion: yup.string().when('type', {
    is: type => type === 'AWS',
    then: yup.string().required(),
    otherwise: yup.string()
  }),
  awsSecretRef: yup.object().shape({
    name: yup.string().when('type', {
      is: type => type === 'AWS',
      then: yup.string().required(),
      otherwise: yup.string()
    }),
    namespace: yup.string().when('type', {
      is: type => type === 'AWS',
      then: yup.string().required(),
      otherwise: yup.string()
    })
  }),
  staticHostList: yup.array().of(
    yup.object().shape({
      addr: yup.string().min(1, 'Invalid host address'),
      port: yup.number().min(10, 'Invalid port number')
    })
  )
});

const CreateUpstreamFormC: React.FC<Props & RouteComponentProps> = props => {
  const {
    config: { namespace, namespacesList }
  } = useSelector((state: AppState) => state);
  const dispatch = useDispatch();
  const initialValues = {
    name: '',
    type: '',
    namespace,
    ...awsInitialValues,
    ...kubeInitialValues,
    ...staticInitialValues,
    ...azureInitialValues,
    ...consulInitialValues,
    awsSecretRef: {
      ...awsInitialValues.awsSecretRef,
      namespace
    },
    azureSecretRef: {
      ...azureInitialValues.azureSecretRef,
      namespace
    },
    kubeServiceNamespace: namespace
  };

  // grpc request
  async function handleCreateUpstream(values: typeof initialValues) {
    const { name, namespace } = values;
    const ref = { name, namespace };
    if (values.type === UPSTREAM_SPEC_TYPES.AWS) {
      const { awsRegion: region, awsSecretRef: secretRef } = values;

      const aws = {
        region,
        secretRef,
        lambdaFunctionsList: []
      };

      dispatch(createUpstream({ input: { ref, aws } }));
    } else if (values.type === UPSTREAM_SPEC_TYPES.AZURE) {
      const {
        azureFunctionAppName: functionAppName,
        azureSecretRef: secretRef
      } = values;
      const azure = {
        ref,
        functionAppName,
        secretRef,
        functionsList: []
      };
      dispatch(createUpstream({ input: { ref, azure } }));
    } else if (values.type === UPSTREAM_SPEC_TYPES.KUBE) {
      const {
        kubeServiceName: serviceName,
        kubeServiceNamespace: serviceNamespace,
        kubeServicePort: servicePort
      } = values;
      const kube = {
        serviceName,
        serviceNamespace,
        servicePort,
        selectorMap: []
      };
      dispatch(createUpstream({ input: { ref, kube } }));
    } else if (values.type === UPSTREAM_SPEC_TYPES.STATIC) {
      const { staticUseTls: useTls } = values;
      let hostsList = values.staticHostList.map(h => {
        return {
          addr: h.name,
          port: +h.value
        };
      });
      const pb_static = {
        useTls,
        hostsList
      };
      dispatch(createUpstream({ input: { ref, pb_static } }));
    } else if (values.type === UPSTREAM_SPEC_TYPES.CONSUL) {
      const {
        consulConnectEnabled,
        consulDataCentersList,
        consulServiceName,
        consulServiceTagsList
      } = values;
      let consul = {
        connectEnabled: consulConnectEnabled,
        dataCentersList: consulDataCentersList,
        serviceName: consulServiceName,
        serviceTagsList: consulServiceTagsList
      };
      dispatch(createUpstream({ input: { ref, consul } }));
    }

    props.toggleModal(s => !s);
    props.history.push('/upstreams', { showSuccess: true });
  }

  return (
    <Formik
      initialValues={initialValues}
      validationSchema={validationSchema}
      onSubmit={handleCreateUpstream}>
      {formik => (
        <FormContainer>
          <SoloFormTemplate>
            <InputRow>
              <div>
                <SoloFormInput
                  name='name'
                  title='Upstream Name'
                  placeholder='Upstream Name'
                />
              </div>
              <div>
                <SoloFormDropdown
                  name='type'
                  title='Upstream Type'
                  placeholder='Type'
                  options={UPSTREAM_TYPES}
                />
              </div>
              <div>
                <SoloFormTypeahead
                  name='namespace'
                  title='Upstream Namespace'
                  defaultValue={namespace}
                  presetOptions={namespacesList.map(ns => {
                    return { value: ns };
                  })}
                />
              </div>
            </InputRow>
          </SoloFormTemplate>
          {formik.values.type === UPSTREAM_SPEC_TYPES.AWS && (
            <AwsUpstreamForm />
          )}
          {formik.values.type === UPSTREAM_SPEC_TYPES.KUBE && (
            <KubeUpstreamForm />
          )}
          {formik.values.type === UPSTREAM_SPEC_TYPES.STATIC && (
            <StaticUpstreamForm />
          )}
          {formik.values.type === UPSTREAM_SPEC_TYPES.AZURE && (
            <AzureUpstreamForm />
          )}
          {formik.values.type === UPSTREAM_SPEC_TYPES.CONSUL && (
            <ConsulUpstreamForm />
          )}

          <Footer>
            <SoloButton
              onClick={() => formik.handleSubmit()}
              text='Create Upstream'
              disabled={formik.isSubmitting}
            />
          </Footer>
        </FormContainer>
      )}
    </Formik>
  );
};

export const CreateUpstreamForm = withRouter(CreateUpstreamFormC);
