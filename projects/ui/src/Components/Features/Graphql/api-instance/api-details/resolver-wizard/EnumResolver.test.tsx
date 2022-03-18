import React from 'react';
import { render, screen } from '@testing-library/react';
import { EnumResolver } from './EnumResolver';

describe('EnumResolver', () => {
  it('renders an Enum Resolver', () => {
    render(<EnumResolver resolverType='resolver' fields={[]} />);
    expect(screen.getByTestId('enum-resolver')).toBeInTheDocument();
  });
})
