import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import { ResolverTypeSection } from './ResolverTypeSection';
import { Formik, Form } from 'formik';
import * as ace from 'ace-builds/src-noconflict/ace';

ace.config.set(
  'basePath',
  'https://cdnjs.cloudflare.com/ajax/libs/ace/1.4.13/'
);

it('renders a Resolver Type Section', async () => {
  render(
    <Formik
      initialValues={{
        resolverConfig: 'foo',
      }}
      onSubmit={jest.fn()}>
      <Form>
        <ResolverTypeSection isEdit={false} />
      </Form>
    </Formik>
  );
  await waitFor(() =>
    expect(screen.getByTestId('resolver-type-section')).toBeInTheDocument()
  );
});
