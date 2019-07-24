import * as React from 'react';
import styled from '@emotion/styled/macro';
import * as yup from 'yup';
import { Field, Formik, FormikValues, FormikErrors } from 'formik';
import { SoloInput } from 'Components/Common/SoloInput';
import { SoloButton } from 'Components/Common/SoloButton';
import { colors } from 'Styles';
import { SoloNegativeButton } from 'Styles/CommonEmotions/button';
import { OAuth } from 'proto/github.com/solo-io/solo-projects/projects/gloo/api/v1/plugins/extauth/extauth_pb';
import FormItem from 'antd/lib/form/FormItem';
import { SoloFormInput } from 'Components/Common/Form/SoloFormField';

const FormContainer = styled.div`
  display: grid;
  grid-template-columns: 1fr 1fr;
  grid-template-rows: 1fr 1fr 1fr;
  padding-top: 10px;
  grid-gap: 8px;
`;

export const ErrorText = styled.div`
  color: ${colors.grapefruitOrange};
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
  clientId: string;
  callbackPath: string;
  issuerUrl: string;
  appUrl: string;
  secretRefNamespace: string | undefined;
  secretRefName: string | undefined;
}
const defaultValues: ValuesType = {
  clientId: '',
  callbackPath: '',
  issuerUrl: '',
  appUrl: '',
  secretRefNamespace: '',
  secretRefName: ''
};

const validationSchema = yup.object().shape({
  clientId: yup.string().required(),
  secretRefName: yup.string(),
  secretRefNamespace: yup.string(),
  issuerUrl: yup.string().required(),
  appUrl: yup.string().required(),
  callbackPath: yup.string().required()
});

interface Props {
  externalAuth: OAuth.AsObject | undefined;
  externalAuthChanged: (newExternalAuth: OAuth.AsObject) => any;
}

export const ExtAuthForm = (props: Props) => {
  const { externalAuth, externalAuthChanged } = props;

  const initialValues: ValuesType = { ...defaultValues, ...externalAuth };

  const invalid = (values: ValuesType, errors: FormikErrors<ValuesType>) => {
    let isInvalid = Object.keys(errors).length;

    return !!isInvalid;
  };

  const updateExtAuth = (values: ValuesType) => {
    const {
      clientId,
      callbackPath,
      appUrl,
      issuerUrl,
      secretRefName,
      secretRefNamespace
    } = values;
    let newExternalAuth = new OAuth().toObject();

    newExternalAuth = {
      clientId,
      callbackPath,
      appUrl,
      issuerUrl
    };

    if (!!secretRefName && !!secretRefNamespace) {
      newExternalAuth = {
        ...newExternalAuth,
        clientSecretRef: {
          name: secretRefName,
          namespace: secretRefNamespace
        }
      };
    }

    props.externalAuthChanged(newExternalAuth);
  };

  return (
    <Formik
      initialValues={initialValues}
      validationSchema={validationSchema}
      onSubmit={updateExtAuth}>
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
            <div>
              <Field
                name='clientId'
                title='Client ID'
                placeholder='myclientid'
                component={SoloFormInput}
              />
            </div>
            <div>
              <Field
                name='callbackPath'
                title='Callback Path'
                placeholder='/my/callback/path/'
                component={SoloFormInput}
              />
            </div>
            <div>
              <Field
                name='issuerUrl'
                title='Issuer URL'
                placeholder='myclientidtheissuer.com'
                component={SoloFormInput}
              />
            </div>
            <div>
              <Field
                name='appUrl'
                title='App URL'
                placeholder='myapp.com'
                component={SoloFormInput}
              />
            </div>
            <div>
              <Field
                name='secretRefName'
                title='Secret Ref Name'
                placeholder='myoauthsecret'
                component={SoloFormInput}
              />
            </div>
            <div>
              <Field
                name='secretRefNamespace'
                title='Secret Ref Namespace'
                placeholder='gloo-system'
                component={SoloFormInput}
              />
            </div>
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
