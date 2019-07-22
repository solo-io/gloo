import {
  SoloFormCheckbox,
  SoloFormInput,
  SoloFormStringsList
} from 'Components/Common/Form/SoloFormField';
import {
  Footer,
  InputRow,
  SoloFormTemplate
} from 'Components/Common/Form/SoloFormTemplate';
import { SoloButton } from 'Components/Common/SoloButton';
import { Field, Formik, FormikProps } from 'formik';
import * as React from 'react';
import * as yup from 'yup';

interface ConsulVauesType {
  consulServiceName: string;
  consulServiceTagsList: string[];
  consulConnectEnabled: boolean;
  consulDataCentersList: string[];
}

export const consulInitialValues: ConsulVauesType = {
  consulServiceName: '',
  consulServiceTagsList: [''],
  consulConnectEnabled: false,
  consulDataCentersList: ['']
};

interface Props {
  parentForm: FormikProps<ConsulVauesType>;
}

export const consulValidationSchema = yup.object().shape({
  consulServiceName: yup.string(),
  consulServiceTagsList: yup.array().of(yup.string()),
  consulConnectEnabled: yup.boolean(),
  consulDataCentersList: yup.array().of(yup.string())
});

export const ConsulUpstreamForm: React.FC<Props> = ({ parentForm }) => {
  return (
    <Formik<ConsulVauesType>
      validationSchema={consulValidationSchema}
      initialValues={consulInitialValues}
      onSubmit={values => {
        parentForm.setFieldValue(
          'consulConnectEnabled',
          values.consulConnectEnabled
        );
        parentForm.setFieldValue(
          'consulDataCentersList',
          values.consulDataCentersList
        );
        parentForm.setFieldValue('consulServiceName', values.consulServiceName);
        parentForm.setFieldValue(
          'consulServiceTagsList',
          values.consulServiceTagsList
        );

        parentForm.submitForm();
      }}>
      {({ values, handleSubmit }) => (
        <SoloFormTemplate formHeader='Consul Upstream Settings'>
          <InputRow>
            <Field
              name='consulServiceName'
              title='Service Name'
              placeholder='Service Name'
              component={SoloFormInput}
            />
            <Field
              name='consulConnectEnabled'
              title='Enable Consul Connect'
              component={SoloFormCheckbox}
            />
          </InputRow>
          <InputRow>
            <SoloFormTemplate formHeader='Service Tags'>
              <Field
                name='consulServiceTagsList'
                title='Consul Service Tags'
                createNewPromptText='Service Tags'
                component={SoloFormStringsList}
              />
            </SoloFormTemplate>
          </InputRow>
          <InputRow>
            <SoloFormTemplate formHeader='Data Centers'>
              <Field
                name='consulDataCentersList'
                title='Data Centers'
                createNewPromptText='Data Centers'
                component={SoloFormStringsList}
              />
            </SoloFormTemplate>
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
