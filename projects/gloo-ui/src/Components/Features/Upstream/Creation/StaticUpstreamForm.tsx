import {
  SoloFormCheckbox,
  SoloFormInput
} from 'Components/Common/Form/SoloFormField';
import {
  Footer,
  InputRow,
  SoloFormTemplate
} from 'Components/Common/Form/SoloFormTemplate';
import { SoloButton } from 'Components/Common/SoloButton';
import { StringCard } from 'Components/Common/StringCardsList';
import {
  Field,
  FieldArray,
  FieldArrayRenderProps,
  Formik,
  FormikProps
} from 'formik';
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
  staticHostList: [
    {
      addr: '',
      port: 8080
    }
  ],
  staticUseTls: false,
  staticServicePort: ''
};

interface Props {
  parentForm: FormikProps<StaticValuesType>;
}
// TODO: figure out which fields are required
export const staticValidationSchema = yup.object().shape({
  staticHostList: yup.string(),
  staticServicePort: yup.number(),
  staticServiceName: yup.string()
});

export const StaticUpstreamForm: React.FC<Props> = ({ parentForm }) => {
  return (
    <Formik<StaticValuesType>
      validationSchema={staticValidationSchema}
      initialValues={staticInitialValues}
      onSubmit={() => parentForm.submitForm()}>
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
        <InputRow>
          <FieldArray name='staticHostList' render={SoloFormArrayField} />
        </InputRow>
        <Footer>
          <SoloButton
            onClick={parentForm.handleSubmit}
            text='Create Upstream'
            disabled={parentForm.isSubmitting}
          />
        </Footer>
      </SoloFormTemplate>
    </Formik>
  );
};

// TODO: make this component generic once we know how to display input values
export const SoloFormArrayField: React.FC<FieldArrayRenderProps> = ({
  form,
  remove,
  insert,
  name
}) => (
  <div>
    <div key={`${name}-${0}`}>
      <Field
        name={`staticHostList[${0}].addr`}
        title='Host Address'
        placeholder='Host Address'
        component={SoloFormInput}
      />
      <Field
        name={`staticHostList[${0}].port`}
        title='Host Port'
        placeholder='Host Port'
        component={SoloFormInput}
      />
      <SoloButton
        text='add to list'
        onClick={() => insert(0, { addr: '', port: '' })}
      />
    </div>
    {form.values.staticHostList.map((host: any, index: any) => {
      return (
        <StringCard key={`${name}-${index}`}>
          {` ${host.addr || ''}:${host.port || ''}`}
          <div onClick={() => remove(index)}>x</div>
        </StringCard>
      );
    })}
  </div>
);
