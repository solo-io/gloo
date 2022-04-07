import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import { ResolverTypeSection } from './ResolverTypeSection';
import { Formik, Form } from 'formik';
import * as ace from 'ace-builds/src-noconflict/ace';
import { Resolution } from 'proto/github.com/solo-io/solo-apis/api/gloo/graphql.gloo/v1beta1/graphql_pb';

ace.config.set(
  'basePath',
  'https://cdnjs.cloudflare.com/ajax/libs/ace/1.4.13/'
);

describe('ResolverTypeSection', () => {
  it('renders a Resolver Type Section', async () => {
    const listOfResolvers = [['resolver', new Resolution().toObject()]];
    render(
      <Formik
        initialValues={{
          resolverConfig: 'foo',
          listOfResolvers,
        }}
        onSubmit={jest.fn()}>
        <Form>
          <ResolverTypeSection />
        </Form>
      </Formik>
    );
    await waitFor(() => {
      expect(screen.getByTestId('resolver-type-section')).toBeInTheDocument();
      expect(
        screen.getByTestId('create-resolver-from-config')
      ).toBeInTheDocument();
    });
  });
});
