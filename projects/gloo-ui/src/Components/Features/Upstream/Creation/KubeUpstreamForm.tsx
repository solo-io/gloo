import {
  SoloFormInput,
  SoloFormTypeahead
} from 'Components/Common/Form/SoloFormField';
import {
  SoloFormTemplate,
  InputRow
} from 'Components/Common/Form/SoloFormTemplate';
import { Field } from 'formik';
import { NamespacesContext } from 'GlooIApp';
import * as React from 'react';
import * as yup from 'yup';

/* ------------------------------ Upstream Spec ----------------------------- */
/* serviceName: string,
  serviceNamespace: string,
  servicePort: number,
  selectorMap: Array < [string, string] >,
  serviceSpec?: {
    rest?: {
      transformationsMap: Map<string, TransformationTemplate : {
        advancedTemplates: boolean,
        extractorsMap: Array<[string, Extraction: {
          header: string,
          regex: string,
          subgroup: number,
        }]>,
        headersMap: Array<[string, InjaTemplate: {text: string}]>,
        body?: InjaTemplate: {text: string},
        passthrough?: Passthrough: ???,
        mergeExtractorsToBody?: MergeExtractorsToBody: ????,
      }>
    swaggerInfo?: {url: string, inline: string, }
    }
    grpc?: {
      descriptors: Uint8Array | string,
      grpcServicesList: Array<ServiceSpec.GrpcService: {
        packageName: string,
        serviceName: string,
        functionNamesList: Array<string>,}>,}
  }
},
  subsetSpec ?: SubsetSpec : {selectorsList: Array<Selector: { keysList: string[]}> },

*/
// TODO combine with main initial values
export const kubeInitialValues = {
  kubeServiceName: '',
  kubeServiceNamespace: 'gloo-system',
  kubeServicePort: ''
};

interface Props {}

// TODO: figure out which fields are required
export const kubeValidationSchema = yup.object().shape({
  kubeServiceName: yup.string(),
  kubeServiceNamespace: yup.string(),
  kubeServicePort: yup.string()
});

export const KubeUpstreamForm: React.FC<Props> = () => {
  const namespaces = React.useContext(NamespacesContext);

  return (
    <SoloFormTemplate formHeader='Kubernetes Upstream Settings'>
      <InputRow>
        <Field
          name='kubeServiceName'
          title='Service Name'
          placeholder='Service Name'
          component={SoloFormInput}
        />
        <Field
          name='kubeServiceNamespace'
          title='Service Namespace'
          defaultValue='gloo-system'
          presetOptions={namespaces}
          component={SoloFormTypeahead}
        />
        <Field
          name='kubeServicePort'
          title='Service Port'
          placeholder='Service Port'
          component={SoloFormInput}
        />
      </InputRow>
    </SoloFormTemplate>
  );
};
