import styled from '@emotion/styled/macro';
import {
  SoloFormInput,
  SoloFormTypeahead
} from 'Components/Common/Form/SoloFormField';
import { SoloButton } from 'Components/Common/SoloButton';
import { Formik, FormikErrors } from 'formik';
import { NamespacesContext } from 'GlooIApp';
import { OAuth } from 'proto/github.com/solo-io/solo-projects/projects/gloo/api/v1/plugins/extauth/extauth_pb';
import * as React from 'react';
import { colors } from 'Styles';
import { SoloNegativeButton } from 'Styles/CommonEmotions/button';
import * as yup from 'yup';

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
  clientId: yup.string().required('A client ID is required.'),
  secretRefName: yup.string(),
  secretRefNamespace: yup.string(),
  issuerUrl: yup
    .string()
    .url()
    .required('An issuer URL is required.'),
  appUrl: yup
    .string()
    .url()
    .required('An app URL is required.'),
  callbackPath: yup.string().required('A callback path is required.')
});

interface Props {
  externalAuth: OAuth.AsObject | undefined;
  externalAuthChanged: (newExternalAuth: OAuth.AsObject) => any;
}

export const ExtAuthForm = (props: Props) => {
  const { externalAuth, externalAuthChanged } = props;

  const namespaces = React.useContext(NamespacesContext);

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

    externalAuthChanged(newExternalAuth);
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
              <SoloFormInput
                name='clientId'
                title='Client ID'
                placeholder='myclientid'
              />
            </div>
            <div>
              <SoloFormInput
                name='callbackPath'
                title='Callback Path'
                placeholder='/my/callback/path/'
              />
            </div>
            <div>
              <SoloFormInput
                name='issuerUrl'
                title='Issuer URL'
                placeholder='myclientidtheissuer.com'
              />
            </div>
            <div>
              <SoloFormInput
                name='appUrl'
                title='App URL'
                placeholder='myapp.com'
              />
            </div>
            <div>
              <SoloFormInput
                name='secretRefName'
                title='Secret Ref Name'
                placeholder='myoauthsecret'
              />
            </div>
            <div>
              <SoloFormTypeahead
                name='secretRefNamespace'
                title='Secret Ref Namespace'
                defaultValue={namespaces.defaultNamespace}
                presetOptions={namespaces.namespacesList.map(ns => {
                  return { value: ns };
                })}
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
