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

export const AwsUpstreamForm: React.FC<Props> = () => {
  const awsRegions = AWS_REGIONS.map(item => item.name);

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
      </SoloFormTemplate>
    </React.Fragment>
  );
};
