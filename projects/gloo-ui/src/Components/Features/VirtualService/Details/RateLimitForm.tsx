import styled from '@emotion/styled';
import {
  SoloFormDropdown,
  SoloFormInput
} from 'Components/Common/Form/SoloFormField';
import { InputRow } from 'Components/Common/Form/SoloFormTemplate';
import { SoloButton } from 'Components/Common/SoloButton';
import { Label } from 'Components/Common/SoloInput';
import { Formik } from 'formik';
import { RateLimit } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/enterprise/plugins/ratelimit/ratelimit_pb';
import { RateLimitPlugin } from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/virtualservice_pb';
import * as React from 'react';
import { useDispatch } from 'react-redux';
import { useParams } from 'react-router';
import { updateRateLimit } from 'store/virtualServices/actions';
import { colors } from 'Styles';
import { SoloNegativeButton } from 'Styles/CommonEmotions/button';
import * as yup from 'yup';

const FormContainer = styled.div`
  display: grid;
  grid-template-columns: 1fr 1fr;
  grid-gap: 8px;
`;

const FormItem = styled.div`
  display: flex;
  flex-direction: column;
`;

const StrongLabel = styled(Label)`
  font-weight: 600;
  margin-bottom: 0;
`;

const PerText = styled.div`
  line-height: 36px;
  margin: 0 8px;
  color: ${colors.septemberGrey};
`;

const FormFooter = styled.div`
  grid-column: 2;
  display: grid;
  grid-template-columns: 90px 125px;
  grid-template-rows: 33px;
  grid-gap: 5px;
  justify-content: right;
`;

const SmallSoloNegativeButton = styled(SoloNegativeButton)`
  min-width: 0;
`;

interface ValuesType {
  authLimitNumber: number | undefined;
  authLimitTimeUnit: RateLimit.UnitMap[keyof RateLimit.UnitMap] | undefined;
  anonLimitNumber: number | undefined;
  anonLimitTimeUnit: RateLimit.UnitMap[keyof RateLimit.UnitMap] | undefined;
}
let defaultValues: ValuesType = {
  authLimitNumber: undefined,
  authLimitTimeUnit: undefined,
  anonLimitNumber: undefined,
  anonLimitTimeUnit: undefined
};

export const timeOptions = [
  {
    key: RateLimit.Unit.DAY,
    value: RateLimit.Unit.DAY,
    displayValue: 'Day'
  },
  {
    key: RateLimit.Unit.HOUR,
    value: RateLimit.Unit.HOUR,
    displayValue: 'Hour'
  },
  {
    key: RateLimit.Unit.MINUTE,
    value: RateLimit.Unit.MINUTE,
    displayValue: 'Minute'
  },
  {
    key: RateLimit.Unit.SECOND,
    value: RateLimit.Unit.SECOND,
    displayValue: 'Second'
  }
];

const validationSchema = yup.object().shape({
  authLimitNumber: yup
    .number()
    .nullable(true)
    .test(
      'Valid Number',
      'Must be positive or zero',
      val => val === undefined || val > -1
    ),
  authLimitTimeUnit: yup.number().nullable(),
  anonLimitNumber: yup
    .number()
    .nullable(true)
    .test(
      'Valid Number',
      'Must be positive or zero',
      val => val === undefined || val > -1
    ),
  anonLimitTimeUnit: yup.number().nullable()
});

interface Props {
  rateLimits?: RateLimitPlugin.AsObject;
}

export const RateLimitForm = (props: Props) => {
  let { virtualservicename, virtualservicenamespace } = useParams();
  const dispatch = useDispatch();

  const initialValues: ValuesType = { ...defaultValues };

  if (!!props.rateLimits && props.rateLimits.value) {
    if (!!props.rateLimits.value!.authorizedLimits) {
      initialValues.authLimitNumber = props.rateLimits.value!.authorizedLimits.requestsPerUnit;
      initialValues.authLimitTimeUnit = props.rateLimits.value!.authorizedLimits.unit;
    }
    if (!!props.rateLimits.value!.anonymousLimits) {
      initialValues.anonLimitNumber = props.rateLimits.value!.anonymousLimits.requestsPerUnit;
      initialValues.anonLimitTimeUnit = props.rateLimits.value!.anonymousLimits.unit;
    }
  }

  const handleUpdateRateLimit = (values: ValuesType) => {
    dispatch(
      updateRateLimit({
        ref: {
          name: virtualservicename!,
          namespace: virtualservicenamespace!
        },
        rateLimit: {
          anonymousLimits: {
            requestsPerUnit: values.anonLimitNumber || 0,
            unit: values.anonLimitTimeUnit || 0
          },
          authorizedLimits: {
            requestsPerUnit: values.authLimitNumber || 0,
            unit: values.authLimitTimeUnit || 0
          }
        }
      })
    );
  };

  return (
    <Formik
      initialValues={initialValues}
      validationSchema={validationSchema}
      onSubmit={handleUpdateRateLimit}>
      {({
        isSubmitting,
        handleSubmit,
        isValid,
        errors,
        handleReset,
        dirty,
        values,
        setFieldValue
      }) => {
        return (
          <FormContainer>
            <FormItem>
              <StrongLabel>Authorized Requests</StrongLabel>

              <InputRow>
                <div>
                  <SoloFormInput name='authLimitNumber' />
                </div>
                <PerText>per</PerText>
                <div style={{ minWidth: '50%' }}>
                  <SoloFormDropdown
                    name='authLimitTimeUnit'
                    options={timeOptions}
                  />
                </div>
              </InputRow>
            </FormItem>
            <FormItem>
              <StrongLabel>Anonymous Requests</StrongLabel>

              <InputRow>
                <div>
                  <SoloFormInput name='anonLimitNumber' title='' />
                </div>
                <PerText>per</PerText>
                <div style={{ minWidth: '50%' }}>
                  <SoloFormDropdown
                    name='anonLimitTimeUnit'
                    title=''
                    options={timeOptions}
                  />
                </div>
              </InputRow>
            </FormItem>
            <FormFooter>
              <SmallSoloNegativeButton onClick={handleReset} disabled={!dirty}>
                Clear
              </SmallSoloNegativeButton>
              <SoloButton
                onClick={handleSubmit}
                text='Submit'
                disabled={isSubmitting}
              />
            </FormFooter>
          </FormContainer>
        );
      }}
    </Formik>
  );
};
