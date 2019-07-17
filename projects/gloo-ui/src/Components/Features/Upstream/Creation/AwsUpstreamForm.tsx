import { ReactComponent as CloseX } from 'assets/close-x.svg';
import { ReactComponent as GreenPlus } from 'assets/small-green-plus.svg';
import {
  SoloFormInput,
  SoloFormTypeahead,
  SoloSecretRefInput
} from 'Components/Common/Form/SoloFormField';
import {
  Footer,
  InputRow,
  SoloFormTemplate
} from 'Components/Common/Form/SoloFormTemplate';
import { SoloButton } from 'Components/Common/SoloButton';
import {
  Field,
  FieldArray,
  FieldArrayRenderProps,
  Formik,
  FormikProps
} from 'formik';
import { LambdaFunctionSpec } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins/aws/aws_pb';
import { ResourceRef } from 'proto/github.com/solo-io/solo-kit/api/v1/ref_pb';
import * as React from 'react';
import { AWS_REGIONS } from 'utils/upstreamHelpers';
import * as yup from 'yup';

interface AwsValuesType {
  awsRegion: string;
  awsSecretRef: ResourceRef.AsObject;
  awsLambdaFunctionsList: LambdaFunctionSpec.AsObject[];
}

export const awsInitialValues: AwsValuesType = {
  awsRegion: 'us-east-1',
  awsSecretRef: {
    name: '',
    namespace: 'gloo-system'
  },
  awsLambdaFunctionsList: [
    {
      logicalName: '',
      lambdaFunctionName: '',
      qualifier: ''
    }
  ]
};

interface Props {
  parentForm: FormikProps<AwsValuesType>;
}

export const awsValidationSchema = yup.object().shape({
  awsRegion: yup.string(),
  awsSecretRefNamespace: yup.string(),
  awsSecretRefName: yup.string(),
  awsLambdaFunctionsList: yup.array().of(
    yup.object().shape({
      logicalName: yup.string(),
      lambdaFunctionName: yup.string(),
      qualifier: yup.string()
    })
  )
});

export const AwsUpstreamForm: React.FC<Props> = ({ parentForm }) => {
  const awsRegions = AWS_REGIONS.map(item => item.name);

  return (
    <Formik<AwsValuesType>
      validationSchema={awsValidationSchema}
      initialValues={awsInitialValues}
      onSubmit={() => parentForm.submitForm()}>
      <SoloFormTemplate formHeader='AWS Upstream Settings'>
        <InputRow>
          <div>
            <Field
              name='awsRegion'
              title='Region'
              presetOptions={awsRegions}
              component={SoloFormTypeahead}
            />
          </div>
          <Field
            name='awsSecretRef'
            type='aws'
            component={SoloSecretRefInput}
          />
        </InputRow>
        <SoloFormTemplate formHeader='Lambda Functions'>
          <FieldArray name='awsLambdaFunctionsList' render={LambdaFunctions} />
        </SoloFormTemplate>
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

export const LambdaFunctions: React.FC<FieldArrayRenderProps> = ({
  form,
  remove,
  insert,
  name
}) => (
  <React.Fragment>
    <InputRow>
      <Field
        name='awsLambdaFunctionsList[0].logicalName'
        title='Logical Name'
        placeholder='Logical Name'
        component={SoloFormInput}
      />
      <Field
        name='awsLambdaFunctionsList[0].lambdaFunctionName'
        title='Lambda Function Name'
        placeholder='Lambda Function Name'
        component={SoloFormInput}
      />
      <Field
        name='awsLambdaFunctionsList[0].qualifier'
        title='Qualifier'
        placeholder='Qualifier'
        component={SoloFormInput}
      />
      <GreenPlus
        style={{ alignSelf: 'center' }}
        onClick={() =>
          insert(0, {
            logicalName: '',
            lambdaFunctionName: '',
            qualifier: ''
          })
        }
      />
    </InputRow>
    <InputRow>
      {form.values.awsLambdaFunctionsList.map(
        (lambda: LambdaFunctionSpec.AsObject, index: number) => {
          return (
            <div key={lambda.logicalName}>
              <div>{lambda.logicalName}</div>
              <div>{lambda.lambdaFunctionName}</div>
              <div>{lambda.qualifier}</div>
              <CloseX onClick={() => remove(index)} />
            </div>
          );
        }
      )}
    </InputRow>
  </React.Fragment>
);
