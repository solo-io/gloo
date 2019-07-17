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
      onSubmit={() => parentForm.submitForm()}>
      {({ values }) => (
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
            <Field
              name='consulDataCenter'
              title='Data Center'
              placeholder='Data Center'
              component={SoloFormInput}
            />
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
              onClick={parentForm.handleSubmit}
              text='Create Upstream'
              disabled={parentForm.isSubmitting}
            />
          </Footer>
        </SoloFormTemplate>
      )}
    </Formik>
  );
};
