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
import { AWS_REGIONS } from 'utils/upstreamHelpers';
import * as yup from 'yup';

// TODO combine with main initial values
export const awsInitialValues = {
  awsRegion: 'us-east-1',
  awsSecretRefNamespace: '',
  awsSecretRefName: ''
};

interface Props {}

export const awsValidationSchema = yup.object().shape({
  awsRegion: yup.string(),
  awsSecretRefNamespace: yup.string(),
  awsSecretRefName: yup.string()
});

export const AwsUpstreamForm: React.FC<Props> = () => {
  const namespaces = React.useContext(NamespacesContext);

  const awsRegions = AWS_REGIONS.map(item => item.name);

  return (
    <SoloFormTemplate formHeader='AWS Upstream Settings'>
      <InputRow>
        <Field
          name='awsRegion'
          title='Region'
          presetOptions={awsRegions}
          component={SoloFormTypeahead}
        />
        <Field
          name='awsSecretRefNamespace'
          title='Secret Ref Namespace'
          presetOptions={namespaces}
          component={SoloFormTypeahead}
        />
        <Field
          name='awsSecretRefName'
          title='Secret Ref Name'
          placeholder='Secret Ref Name'
          component={SoloFormInput}
        />
      </InputRow>
    </SoloFormTemplate>
  );
};
