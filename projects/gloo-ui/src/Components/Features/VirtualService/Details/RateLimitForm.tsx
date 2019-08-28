import styled from '@emotion/styled';
import {
  SoloFormDropdown,
  SoloFormInput
} from 'Components/Common/Form/SoloFormField';
import { InputRow } from 'Components/Common/Form/SoloFormTemplate';
import { SoloButton } from 'Components/Common/SoloButton';
import { SoloCheckbox } from 'Components/Common/SoloCheckbox';
import { Label } from 'Components/Common/SoloInput';
import { Formik, FormikErrors } from 'formik';
import {
  IngressRateLimit,
  RateLimit
} from 'proto/github.com/solo-io/solo-projects/projects/gloo/api/v1/plugins/ratelimit/ratelimit_pb';
import * as React from 'react';
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
  rates: IngressRateLimit.AsObject | undefined;
  rateLimitsChanged: (newRateLimits: IngressRateLimit.AsObject) => any;
}

export const RateLimitForm = (props: Props) => {
  const [applyAuthorizedLimits, setApplyAuthorizedLimits] = React.useState(
    !!props.rates && !!props.rates.authorizedLimits
  );
  const [applyAnonymousLimits, setApplyAnonymousLimits] = React.useState(
    !!props.rates && !!props.rates.anonymousLimits
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

  const updateRateLimit = (values: ValuesType) => {
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

    props.rateLimitsChanged(newRateLimits);
  };

  const isDirty = (formIsDirty: boolean) => {
    return (
      formIsDirty ||
      applyAuthorizedLimits !==
        (!!props.rates && !!props.rates.authorizedLimits) ||
      applyAnonymousLimits !== (!!props.rates && !!props.rates.anonymousLimits)
    );
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
        values,
        setFieldValue
      }) => {
        return (
          <FormContainer>
            <FormItem>
              <StrongLabel>
                <SoloCheckbox
                  title={'Apply Authorized Limits '}
                  checked={applyAuthorizedLimits}
                  onChange={() => {
                    if (
                      !applyAuthorizedLimits &&
                      values.authLimitNumber === undefined
                    ) {
                      setFieldValue('authLimitNumber', 1);
                    }
                    setApplyAuthorizedLimits(s => !s);
                  }}
                />
              </StrongLabel>
              {applyAuthorizedLimits && (
                <InputRow>
                  <div>
                    <SoloFormInput
                      name='authLimitNumber'
                      title=''
                      placeholder='##'
                    />
                  </div>
                  <PerText>per</PerText>
                  <div style={{ minWidth: '50%' }}>
                    <SoloFormDropdown
                      name='authLimitTimeUnit'
                      title=''
                      options={timeOptions}
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
                  onChange={() => {
                    if (
                      !applyAnonymousLimits &&
                      values.anonLimitNumber === undefined
                    ) {
                      setFieldValue('anonLimitNumber', 1);
                    }
                    setApplyAnonymousLimits(s => !s);
                  }}
                />
              </StrongLabel>
              {applyAnonymousLimits && (
                <InputRow>
                  <div>
                    <SoloFormInput
                      name='anonLimitNumber'
                      title=''
                      placeholder='##'
                    />
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
              )}
            </FormItem>
            <FormFooter>
              <SmallSoloNegativeButton onClick={handleReset} disabled={!dirty}>
                Clear
              </SmallSoloNegativeButton>
              <SoloButton
                onClick={handleSubmit}
                text='Submit'
                disabled={
                  isSubmitting || invalid(values, errors) || !isDirty(dirty)
                }
              />
            </FormFooter>
          </FormContainer>
        );
      }}
    </Formik>
  );
};
