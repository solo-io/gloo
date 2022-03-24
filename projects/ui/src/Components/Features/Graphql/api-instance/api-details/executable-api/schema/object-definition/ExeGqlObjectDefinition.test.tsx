import * as React from 'react';
import { ExeGqlObjectDefinition } from './ExeGqlObjectDefinition';
import { render, screen } from '@testing-library/react';

describe('ExeGqlObjectDefinition', () => {
  it('Shoud be able to render', () => {
    render(
      <ExeGqlObjectDefinition
        apiRef={{ name: '', namespace: '', clusterName: '' }}
        resolverType='resolver'
        onReturnTypeClicked={t => console.log(t)}
        schemaDefinitions={[]}
        fields={[]}
      />
    );
    expect(screen.getByTestId('resolver-item')).toBeInTheDocument();
  });
});
