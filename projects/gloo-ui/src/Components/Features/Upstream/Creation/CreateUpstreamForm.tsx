import styled from '@emotion/styled';
import {
  SoloFormDropdown,
  SoloFormInput,
  SoloFormTypeahead
} from 'Components/Common/Form/SoloFormField';
import {
  Footer,
  InputRow,
  SoloFormTemplate
} from 'Components/Common/Form/SoloFormTemplate';
import { SoloButton } from 'Components/Common/SoloButton';
import { Formik } from 'formik';
import * as React from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { useHistory } from 'react-router';
import { AppState } from 'store';
import { createUpstream } from 'store/upstreams/actions';
import { UPSTREAM_SPEC_TYPES, UPSTREAM_TYPES } from 'utils/upstreamHelpers';
import * as yup from 'yup';
import { awsInitialValues, AwsUpstreamForm } from './AwsUpstreamForm';
import { azureInitialValues, AzureUpstreamForm } from './AzureUpstreamForm';
import { consulInitialValues, ConsulUpstreamForm } from './ConsulUpstreamForm';
import { kubeInitialValues, KubeUpstreamForm } from './KubeUpstreamForm';
import { staticInitialValues, StaticUpstreamForm } from './StaticUpstreamForm';
import { Upstream } from 'proto/gloo/projects/gloo/api/v1/upstream_pb';
import { CreateUpstreamRequest } from 'proto/solo-projects/projects/grpcserver/api/v1/upstream_pb';

interface Props {
  toggleModal: React.Dispatch<React.SetStateAction<boolean>>;
}

const FormContainer = styled.div`
  display: flex;
  flex-direction: column;
`;

const FormRow = styled.div`
  display: grid;
  grid-template-columns: 1fr 1fr 1fr;
  grid-gap: 15px;
`;

// TODO combine validation schemas
const validationSchema = yup.object().shape({
  name: yup
    .string()
    .required('Upstream name is required')
    .max(254, `Names must be 254 characters or shorter`)
    .test(
      'only lowercase',
      'Names may only be lower-case',
      val => val && val.toLowerCase() === val
    )
    .test(
      'start and end with letters and/or numbers',
      'Names must start and end with alphanumerics',
      val => {
        if (!val) {
          return false;
        }
        const firstCharVal = val[0].charCodeAt(0);
        const lastCharVal = val[val.length - 1].charCodeAt(0);

        if (
          ((firstCharVal >= 48 && firstCharVal <= 57) ||
            (firstCharVal >= 97 && firstCharVal <= 122)) &&
          ((lastCharVal >= 48 && lastCharVal <= 57) ||
            (lastCharVal >= 97 && lastCharVal <= 122))
        ) {
          return true;
        }

        return false;
      }
    )
    .test(
      'Regex test',
      'Must consist of lower case alphanumeric characters, " - " or ".", and must start and end with an alphanumeric character',
      val => {
        if (!val) {
          return false;
        }
        if (val.length < 2) return true;

        const regexTest = /^[a-z0-9]+[-.a-z0-9]*[a-z0-9]{1}$/;
        return !!val.match(regexTest) && val.match(regexTest)[0] === val;
      }
    ),
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
  staticHostList: yup
    .array()
    .of(
      yup.object().shape({
        addr: yup.string().min(1, 'Invalid host address'),
        port: yup.number().min(10, 'Invalid port number')
      })
    )
    .when('type', {
      is: type => type === 'Static',
      then: yup
        .array()
        .of(
          yup.object().shape({
            addr: yup.string().min(1, 'Invalid host address'),
            port: yup.number().min(10, 'Invalid port number')
          })
        )
        .required('You need to specify at least one host'),
      otherwise: yup.array().of(
        yup.object().shape({
          addr: yup.string().min(1, 'Invalid host address'),
          port: yup.number().min(10, 'Invalid port number')
        })
      )
    })
});

export const CreateUpstreamForm: React.FC<Props> = props => {
  let history = useHistory();
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

    let initialUpstream = new CreateUpstreamRequest().toObject()?.upstreamInput;
    let initialUpstreamSpec = new Upstream().toObject();
    let initialUpstreamInput = {
      ...initialUpstream,
      metadata: {
        ...initialUpstream?.metadata!,
        name: values.name,
        namespace: values.namespace
      }
    };
    if (values.type === UPSTREAM_SPEC_TYPES.AWS) {
      const { awsRegion: region, awsSecretRef: secretRef } = values;
      const aws: Upstream.AsObject = {
        ...initialUpstreamSpec,
        aws: {
          region,
          secretRef,
          lambdaFunctionsList: []
        }
      };

      dispatch(
        createUpstream({
          upstreamInput: {
            ...aws,
            ...initialUpstreamInput
          }
        })
      );
    } else if (values.type === UPSTREAM_SPEC_TYPES.AZURE) {
      const {
        azureFunctionAppName: functionAppName,
        azureSecretRef: secretRef
      } = values;
      const azure: Upstream.AsObject = {
        ...initialUpstreamSpec,
        azure: {
          functionAppName,
          secretRef,
          functionsList: []
        }
      };

      dispatch(
        createUpstream({
          upstreamInput: {
            ...azure,
            ...initialUpstreamInput
          }
        })
      );
    } else if (values.type === UPSTREAM_SPEC_TYPES.KUBE) {
      const {
        kubeServiceName: serviceName,
        kubeServiceNamespace: serviceNamespace,
        kubeServicePort: servicePort
      } = values;
      const kube: Upstream.AsObject = {
        ...initialUpstreamSpec,
        kube: {
          serviceName,
          serviceNamespace,
          servicePort,
          selectorMap: []
        }
      };
      dispatch(
        createUpstream({
          upstreamInput: {
            ...kube,
            ...initialUpstreamInput
          }
        })
      );
    } else if (values.type === UPSTREAM_SPEC_TYPES.STATIC) {
      const { staticUseTls: useTls } = values;
      let hostsList = values.staticHostList.map(h => {
        return {
          addr: h.name,
          port: +h.value
        };
      });
      const pb_static: Upstream.AsObject = {
        ...initialUpstreamSpec,

        pb_static: {
          useTls,
          hostsList
        }
      };
      dispatch(
        createUpstream({
          upstreamInput: {
            ...pb_static,
            ...initialUpstreamInput
          }
        })
      );
    } else if (values.type === UPSTREAM_SPEC_TYPES.CONSUL) {
      const {
        consulConnectEnabled,
        consulDataCentersList,
        consulServiceName,
        consulServiceTagsList
      } = values;
      let consul: Upstream.AsObject = {
        ...initialUpstreamSpec,
        consul: {
          connectEnabled: consulConnectEnabled,
          dataCentersList: consulDataCentersList,
          serviceName: consulServiceName,
          serviceTagsList: consulServiceTagsList
        }
      };
      dispatch(
        createUpstream({
          upstreamInput: {
            ...consul,
            ...initialUpstreamInput
          }
        })
      );
    }

    props.toggleModal(s => !s);
    history.push('/upstreams');
  }

  return (
    <Formik
      initialValues={initialValues}
      validationSchema={validationSchema}
      onSubmit={handleCreateUpstream}>
      {formik => (
        <FormContainer>
          <SoloFormTemplate>
            <FormRow>
              <div>
                <SoloFormInput
                  name='name'
                  title='Upstream Name'
                  placeholder='Upstream Name'
                />
              </div>
              <div>
                <SoloFormDropdown
                  testId='upstream-type'
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
            </FormRow>
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
