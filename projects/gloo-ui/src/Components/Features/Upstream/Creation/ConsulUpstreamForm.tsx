import {
  SoloFormCheckbox,
  SoloFormInput
} from 'Components/Common/Form/SoloFormField';
import {
  SoloFormTemplate,
  InputRow,
  Footer
} from 'Components/Common/Form/SoloFormTemplate';
import { Field, FormikProps, Formik, FieldArray } from 'formik';
import * as React from 'react';
import * as yup from 'yup';
import { SoloButton } from 'Components/Common/SoloButton';

interface ConsulVauesType {
  consulServiceName: string;
  consulServiceTagsList: string[];
  consulConnectEnabled: boolean;
  consulDataCenter: string;
}

export const consulInitialValues: ConsulVauesType = {
  consulServiceName: '',
  consulServiceTagsList: [],
  consulConnectEnabled: false,
  consulDataCenter: ''
};

interface Props {
  parentForm: FormikProps<ConsulVauesType>;
}

export const consulValidationSchema = yup.object().shape({
  consulServiceName: yup.string(),
  consulServiceTagsList: yup.array().of(yup.string()),
  consulConnectEnabled: yup.boolean(),
  consulDataCenter: yup.string()
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
        parentForm.setFieldValue('consulDataCenter', values.consulDataCenter);
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
            <div>
              <Field
                name='consulServiceName'
                title='Service Name'
                placeholder='Service Name'
                component={SoloFormInput}
              />
            </div>
            <div>
              <Field
                name='consulConnectEnabled'
                title='Enable Consul Connect'
                component={SoloFormCheckbox}
              />
            </div>
            <div>
              <Field
                name='consulDataCenter'
                title='Data Center'
                placeholder='Data Center'
                component={SoloFormInput}
              />
            </div>
            {/* TODO: Use String Cards List component */}
            <FieldArray
              name='consulServiceTagsList'
              render={({ form, remove, insert, name }) => (
                <React.Fragment>
                  <Field
                    name={`consulServiceTagsList.[0]`}
                    title='Consul Service Tags'
                    placeholder='Service Tags'
                    component={SoloFormInput}
                  />
                  <div>
                    {values.consulServiceTagsList.map(tag => (
                      <div key={tag}>{tag} </div>
                    ))}
                  </div>
                </React.Fragment>
              )}
            />
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
