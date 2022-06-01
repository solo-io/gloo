import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import { ResolverConfigSection } from '../ResolverConfigSection';
import { Formik, Form } from 'formik';
import * as ace from 'ace-builds/src-noconflict/ace';
import 'brace/ext/language_tools';
import { AppSettingsProvider } from 'Components/Context/AppSettingsContext';

ace.config.set(
  'basePath',
  'https://cdnjs.cloudflare.com/ajax/libs/ace/1.4.13/'
);
// useAppSettings
jest.mock('Components/Context/AppSettingsContext', () => {
  const actual = jest.requireActual('Components/Context/AppSettingsContext');
  return {
    ...actual,
    useAppSettings: () => {
      return {
        appSettings: {
          keyboardHandler: '',
        },
      };
    },
  };
});

describe('ResolverConfigSection', () => {
  it('renders a Resolver Config Section', async () => {
    render(
      <AppSettingsProvider>
        <Formik
          initialValues={{
            resolverConfig: 'foo',
            listOfResolvers: [],
          }}
          onSubmit={jest.fn()}>
          <Form>
            <ResolverConfigSection
              onCancel={jest.fn()}
              formik={
                {
                  values: {
                    upstream: '',
                    listOfResolvers: [],
                    resolverConfig: 'testing',
                    resolverType: 'Mock',
                  },
                  errors: [],
                  isValid: true,
                } as any
              }
              globalWarningMessage=''
            />
          </Form>
        </Formik>
      </AppSettingsProvider>
    );
    await waitFor(() =>
      expect(screen.getByTestId('resolver-config-section')).toBeInTheDocument()
    );
  });
});
