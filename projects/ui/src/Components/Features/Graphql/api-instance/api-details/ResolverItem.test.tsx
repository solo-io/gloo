import * as React from 'react';
import { ResolverItem } from './ResolverItem';
import { render, screen } from '@testing-library/react';

describe('ResolverItem', () => {
  it('Shoud be able to render', () => {
    render(
      <ResolverItem
        resolverType='resolver'
        fields={[]}
        handleResolverConfigModal={jest.fn()}
      />
    );
    expect(screen.getByTestId('resolver-item')).toBeInTheDocument();
  });
});
