import {
  SoloFormSecretRefInput,
  SoloFormTypeahead
} from 'Components/Common/Form/SoloFormField';
import {
  InputRow,
  SoloFormTemplate
} from 'Components/Common/Form/SoloFormTemplate';
import { ResourceRef } from 'proto/github.com/solo-io/solo-kit/api/v1/ref_pb';
import * as React from 'react';
import { AWS_REGIONS } from 'utils/upstreamHelpers';
import * as yup from 'yup';
import { SoloButton } from 'Components/Common/SoloButton';
import { withRouter, RouterProps } from 'react-router';
import { useFormikContext, useField } from 'formik';
interface AwsValuesType {
  awsRegion: string;
  awsSecretRef: ResourceRef.AsObject;
}

export const awsInitialValues: AwsValuesType = {
  awsRegion: 'us-east-1',
  awsSecretRef: {
    name: '',
    namespace: 'gloo-system'
  }
};

interface Props {}

export const awsValidationSchema = yup.object().shape({
  awsRegion: yup.string(),
  awsSecretRef: yup.object().shape({
    name: yup.string().required('You need to provide a secret'),
    namespace: yup.string()
  })
});

const AwsUpstreamFormComponent: React.FC<Props & RouterProps> = ({
  history
}) => {
  const awsRegions = AWS_REGIONS.map(item => item.name);
  const [_, meta] = useField('awsSecretRef');

  return (
    <React.Fragment>
      <SoloFormTemplate formHeader='AWS Upstream Settings'>
        <InputRow>
          <div>
            <SoloFormTypeahead
              name='awsRegion'
              title='Region'
              presetOptions={awsRegions}
            />
          </div>
          <SoloFormSecretRefInput name='awsSecretRef' type='aws' />
        </InputRow>
        <InputRow>
          {!!meta.error && !!meta.touched && (
            <SoloButton
              text='Create a secret first'
              onClick={() => history.push('/settings/secrets')}
            />
          )}
        </InputRow>
      </SoloFormTemplate>
    </React.Fragment>
  );
};

export const AwsUpstreamForm = withRouter(AwsUpstreamFormComponent);
