import * as React from 'react';
import styled from '@emotion/styled/macro';
import { Field, Formik, FormikValues, FormikErrors } from 'formik';
import * as yup from 'yup';
import { colors } from 'Styles';
import { Label, SoloInput } from 'Components/Common/SoloInput';
import { SoloNegativeButton } from 'Styles/CommonEmotions/button';
import { SoloButton } from 'Components/Common/SoloButton';
import {
  IngressRateLimit,
  RateLimit
} from 'proto/github.com/solo-io/solo-projects/projects/gloo/api/v1/plugins/ratelimit/ratelimit_pb';
import { InputRow } from 'Components/Common/Form/SoloFormTemplate';
import {
  SoloFormInput,
  SoloFormDropdown
} from 'Components/Common/Form/SoloFormField';
import { SoloCheckbox } from 'Components/Common/SoloCheckbox';

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
  authLimitTimeUnit: RateLimit.Unit | undefined;
  anonLimitNumber: number | undefined;
  anonLimitTimeUnit: RateLimit.Unit | undefined;
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
    .test('Valid Number', 'Greater than 0', val => val > 0),
  authLimitTimeUnit: yup.number(),
  anonLimitNumber: yup
    .number()
    .test('Valid Number', 'Greater than 0', val => val > 0),
  anonLimitTimeUnit: yup.number()
});

interface Props {
  rates: IngressRateLimit.AsObject | undefined;
  rateLimitsChanged: (newRateLimits: IngressRateLimit.AsObject) => any;
}

export const RateLimitForm = (props: Props) => {
  const [applyAuthorizedLimits, setApplyAuthorizedLimits] = React.useState(
    !!props.rates && !!props.rates.authorizedLimits
  );
  const [applyAnonymousLimits, setApplyAnonymousLimits] = React.useState(
    !!props.rates && !!props.rates.authorizedLimits
  );

  const initialValues: ValuesType = { ...defaultValues };

  if (!!props.rates) {
    if (!!props.rates.authorizedLimits) {
      initialValues.authLimitNumber =
        props.rates.authorizedLimits.requestsPerUnit;
      initialValues.authLimitTimeUnit = props.rates.authorizedLimits.unit;
    }
    if (!!props.rates.anonymousLimits) {
      initialValues.anonLimitNumber =
        props.rates.anonymousLimits.requestsPerUnit;
      initialValues.anonLimitTimeUnit = props.rates.anonymousLimits.unit;
    }
  }

  const invalid = (values: ValuesType, errors: FormikErrors<ValuesType>) => {
    let isInvalid = false;

    if (applyAnonymousLimits) {
      isInvalid = !!errors.anonLimitNumber || !values.anonLimitTimeUnit;
    }
    if (applyAuthorizedLimits) {
      isInvalid =
        isInvalid || !!errors.authLimitNumber || !values.authLimitTimeUnit;
    }

    return isInvalid;
  };

  const updateRateLimit = (values: typeof defaultValues) => {
    const newRateLimits = new IngressRateLimit().toObject();

    if (applyAnonymousLimits) {
      newRateLimits.anonymousLimits = {
        unit: values.anonLimitTimeUnit!,
        requestsPerUnit: values.anonLimitNumber!
      };
    }
    if (applyAuthorizedLimits) {
      newRateLimits.authorizedLimits = {
        unit: values.authLimitTimeUnit!,
        requestsPerUnit: values.authLimitNumber!
      };
    }
    console.log(newRateLimits);

    props.rateLimitsChanged(newRateLimits);
  };

  return (
    <Formik
      initialValues={initialValues}
      validationSchema={validationSchema}
      onSubmit={updateRateLimit}>
      {({
        isSubmitting,
        handleSubmit,
        isValid,
        errors,
        handleReset,
        dirty,
        values
      }) => {
        return (
          <FormContainer>
            <FormItem>
              <StrongLabel>
                <SoloCheckbox
                  title={'Apply Authorized Limits '}
                  checked={applyAuthorizedLimits}
                  onChange={() => setApplyAuthorizedLimits(s => !s)}
                />
              </StrongLabel>
              {applyAuthorizedLimits && (
                <InputRow>
                  <div>
                    <Field
                      name='authLimitNumber'
                      title=''
                      placeholder='##'
                      component={SoloFormInput}
                    />
                  </div>
                  <PerText>per</PerText>
                  <div style={{ minWidth: '50%' }}>
                    <Field
                      name='authLimitTimeUnit'
                      title=''
                      options={timeOptions}
                      component={SoloFormDropdown}
                    />
                  </div>
                </InputRow>
              )}
            </FormItem>
            <FormItem>
              <StrongLabel>
                <SoloCheckbox
                  title={'Apply Anonymous Limits '}
                  checked={applyAnonymousLimits}
                  onChange={() => setApplyAnonymousLimits(s => !s)}
                />
              </StrongLabel>
              {applyAnonymousLimits && (
                <InputRow>
                  <div>
                    <Field
                      name='anonLimitNumber'
                      title=''
                      placeholder='##'
                      component={SoloFormInput}
                    />
                  </div>
                  <PerText>per</PerText>
                  <div style={{ minWidth: '50%' }}>
                    <Field
                      name='anonLimitTimeUnit'
                      title=''
                      options={timeOptions}
                      component={SoloFormDropdown}
                    />
                  </div>
                </InputRow>
              )}
            </FormItem>
            <FormFooter>
              <SmallSoloNegativeButton onClick={handleReset} disabled={!dirty}>
                Clear
              </SmallSoloNegativeButton>
              <SoloButton
                onClick={handleSubmit}
                text='Submit'
                disabled={isSubmitting || invalid(values, errors) || !dirty}
              />
            </FormFooter>
          </FormContainer>
        );
      }}
    </Formik>
  );
};
