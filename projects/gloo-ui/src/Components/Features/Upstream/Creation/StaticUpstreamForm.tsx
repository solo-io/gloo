import {
  SoloFormCheckbox,
  SoloFormInput,
  SoloFormMultipartStringCardsList
} from 'Components/Common/Form/SoloFormField';
import {
  Footer,
  InputRow,
  SoloFormTemplate
} from 'Components/Common/Form/SoloFormTemplate';
import { SoloButton } from 'Components/Common/SoloButton';
import { Field, Formik, FormikProps } from 'formik';
import { Host } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins/static/static_pb';
import * as React from 'react';
import * as yup from 'yup';

// TODO: handle service spec
interface StaticValuesType {
  staticServiceName: string;
  staticHostList: Host.AsObject[];
  staticUseTls: boolean;
  staticServicePort: string;
}

export const staticInitialValues: StaticValuesType = {
  staticServiceName: '',
  staticHostList: [],
  staticUseTls: false,
  staticServicePort: ''
};

interface Props {
  parentForm: FormikProps<StaticValuesType>;
}
// TODO: figure out which fields are required
export const staticValidationSchema = yup.object().shape({
  staticServicePort: yup.number(),
  staticServiceName: yup.string()
});

export const StaticUpstreamForm: React.FC<Props> = ({ parentForm }) => {
  return (
    <Formik<StaticValuesType>
      validationSchema={staticValidationSchema}
      initialValues={staticInitialValues}
      onSubmit={values => {
        parentForm.setFieldValue('staticHostList', values.staticHostList);
        parentForm.setFieldValue('staticServiceName', values.staticServiceName);
        parentForm.setFieldValue('staticServicePort', values.staticServicePort);
        parentForm.setFieldValue('staticUseTls', values.staticUseTls);

        parentForm.submitForm();
      }}>
      {({ handleSubmit }) => (
        <SoloFormTemplate formHeader='Static Upstream Settings'>
          <InputRow>
            <Field
              name='staticServiceName'
              title='Service Name'
              placeholder='Service Name'
              component={SoloFormInput}
            />
            <Field
              name='staticUseTls'
              title='Use Tls'
              component={SoloFormCheckbox}
            />
            <Field
              name='staticServicePort'
              title='Service Port'
              placeholder='Service Port'
              type='number'
              component={SoloFormInput}
            />
          </InputRow>
          <SoloFormTemplate formHeader='Host List'>
            <InputRow>
              <Field
                name='staticHostList'
                createNewNamePromptText={'Address...'}
                createNewValuePromptText={'Port...'}
                component={SoloFormMultipartStringCardsList}
              />
            </InputRow>
          </SoloFormTemplate>

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
