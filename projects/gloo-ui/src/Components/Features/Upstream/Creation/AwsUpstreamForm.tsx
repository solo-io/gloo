import {
  SoloFormInput,
  SoloFormTypeahead
} from 'Components/Common/Form/SoloFormField';
import { SoloFormTemplate } from 'Components/Common/Form/SoloFormTemplate';
import { Field } from 'formik';
import { NamespacesContext } from 'GlooIApp';
import * as React from 'react';
import { AWS_REGIONS } from 'utils/upstreamHelpers';
import * as yup from 'yup';

// TODO combine with main initial values
export const awsInitialValues = {
  region: 'us-east-1',
  secretRefNamespace: '',
  secretRefName: ''
};

interface Props {}

export const awsValidationSchema = yup.object().shape({
  region: yup.string(),
  secretRefNamespace: yup.string(),
  secretRefName: yup.string()
});

export const AwsUpstreamForm: React.FC<Props> = () => {
  const namespaces = React.useContext(NamespacesContext);

  const awsRegions = AWS_REGIONS.map(item => item.name);

  return (
    <SoloFormTemplate formHeader='AWS Upstream Settings'>
      <Field
        name='region'
        title='Region'
        presetOptions={awsRegions}
        component={SoloFormTypeahead}
      />
      <Field
        name='secretRefNamespace'
        title='Secret Ref Namespace'
        presetOptions={namespaces}
        component={SoloFormTypeahead}
      />
      <Field
        name='secretRefName'
        title='Secret Ref Name'
        placeholder='Secret Ref Name'
        component={SoloFormInput}
      />
    </SoloFormTemplate>
  );
};
