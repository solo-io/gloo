import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import { UpstreamSection } from './UpstreamSection';
import { Formik, Form } from 'formik';
import * as ace from 'ace-builds/src-noconflict/ace';

jest.mock('API/hooks', () => {
  return {
    useListUpstreams: () => {
      return {
        data: [],
      };
    },
  };
});

ace.config.set(
  'basePath',
  'https://cdnjs.cloudflare.com/ajax/libs/ace/1.4.13/'
);

describe('UpstreamSection', () => {
  it('renders an Upstream Section', async () => {
    render(
      <Formik
        initialValues={{
          resolverConfig: 'foo',
        }}
        onSubmit={jest.fn()}>
        <Form>
          <UpstreamSection
            onCancel={jest.fn()}
            onNextClicked={jest.fn()}
            nextButtonDisabled={false}
            existingUpstreamId={''}
            setWarningMessage={jest.fn()}
          />
        </Form>
      </Formik>
    );
    await waitFor(() =>
      expect(screen.getByTestId('upstream-section')).toBeInTheDocument()
    );
  });
});
