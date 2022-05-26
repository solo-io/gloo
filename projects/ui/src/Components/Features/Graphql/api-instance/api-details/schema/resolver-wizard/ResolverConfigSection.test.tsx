import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import { ResolverConfigSection } from './ResolverConfigSection';
import { Formik, Form } from 'formik';
import * as ace from 'ace-builds/src-noconflict/ace';
import 'brace/ext/language_tools';

ace.config.set(
  'basePath',
  'https://cdnjs.cloudflare.com/ajax/libs/ace/1.4.13/'
);

describe('ResolverConfigSection', () => {
  it('renders a Resolver Config Section', async () => {
    render(
      <Formik
        initialValues={{
          resolverConfig: 'foo',
        }}
        onSubmit={jest.fn()}>
        <Form>
          <ResolverConfigSection
            onCancel={jest.fn()}
            submitDisabled={false}
            warningMessage=''
          />
        </Form>
      </Formik>
    );
    await waitFor(() =>
      expect(screen.getByTestId('resolver-config-section')).toBeInTheDocument()
    );
  });
});
