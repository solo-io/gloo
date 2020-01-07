import {
  SoloAWSSecretsList,
  SoloFormTypeahead
} from 'Components/Common/Form/SoloFormField';
import {
  InputRow,
  SoloFormTemplate
} from 'Components/Common/Form/SoloFormTemplate';
import { SoloButton } from 'Components/Common/SoloButton';
import { useField } from 'formik';
import { ResourceRef } from 'proto/solo-kit/api/v1/ref_pb';
import * as React from 'react';
import { AWS_REGIONS } from 'utils/upstreamHelpers';
import * as yup from 'yup';
import { useHistory } from 'react-router';

export interface AwsValuesType {
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

export const awsValidationSchema = yup.object().shape({
  awsRegion: yup.string(),
  awsSecretRef: yup.object().shape({
    name: yup.string(),
    namespace: yup.string()
  })
});

export const AwsUpstreamForm = () => {
  let history = useHistory();
  const awsRegions = AWS_REGIONS.map(item => item.name);
  const [_, meta] = useField('awsSecretRef');

  return (
    <>
      <SoloFormTemplate formHeader='AWS Upstream Settings'>
        <InputRow style={{ justifyContent: 'spaceAround' }}>
          <div>
            <SoloFormTypeahead
              testId='aws-region'
              name='awsRegion'
              title='Region'
              presetOptions={awsRegions.map(region => {
                return { value: region };
              })}
            />
          </div>
          <div>
            <SoloAWSSecretsList
              testId='aws-secret'
              name='awsSecretRef'
              type='aws'
            />
          </div>
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
    </>
  );
};
