import {
  SoloFormInput,
  SoloFormTypeahead
} from 'Components/Common/Form/SoloFormField';
import {
  Footer,
  InputRow,
  SoloFormTemplate
} from 'Components/Common/Form/SoloFormTemplate';
import { SoloButton } from 'Components/Common/SoloButton';
import { Field, Formik, FormikProps } from 'formik';
import { NamespacesContext } from 'GlooIApp';
import * as React from 'react';
import * as yup from 'yup';

// TODO: handle servicespec and subset spec
interface KubeValuesType {
  kubeServiceName: string;
  kubeServiceNamespace: string;
  kubeServicePort: number;
}

export const kubeInitialValues: KubeValuesType = {
  kubeServiceName: '',
  kubeServiceNamespace: 'gloo-system',
  kubeServicePort: 8080
};

interface Props {
  parentForm: FormikProps<KubeValuesType>;
}

// TODO: figure out which fields are required
export const kubeValidationSchema = yup.object().shape({
  kubeServiceName: yup.string(),
  kubeServiceNamespace: yup.string(),
  kubeServicePort: yup.number()
});

export const KubeUpstreamForm: React.FC<Props> = ({ parentForm }) => {
  const namespaces = React.useContext(NamespacesContext);

  return (
    <Formik<KubeValuesType>
      validationSchema={kubeValidationSchema}
      initialValues={kubeInitialValues}
      onSubmit={values => {
        parentForm.setFieldValue('kubeServiceName', values.kubeServiceName);
        parentForm.setFieldValue(
          'kubeServiceNamespace',
          values.kubeServiceNamespace
        );
        parentForm.setFieldValue('kubeServicePort', values.kubeServicePort);
        parentForm.submitForm();
      }}>
      {({ handleSubmit }) => (
        <SoloFormTemplate formHeader='Kubernetes Upstream Settings'>
          <InputRow>
            <div>
              <Field
                name='kubeServiceName'
                title='Service Name'
                placeholder='Service Name'
                component={SoloFormInput}
              />
            </div>
            <div>
              <Field
                name='kubeServiceNamespace'
                title='Service Namespace'
                defaultValue='gloo-system'
                presetOptions={namespaces}
                component={SoloFormTypeahead}
              />
            </div>
            <div>
              <Field
                name='kubeServicePort'
                title='Service Port'
                placeholder='Service Port'
                type='number'
                component={SoloFormInput}
              />
            </div>
          </InputRow>
          <Footer>
            <SoloButton
              onClick={handleSubmit}
              text='Create Upstream'
              disabled={parentForm.isSubmitting}
            />
          </Footer>
        </SoloFormTemplate>
      )}
    </Formik>
  );
};
