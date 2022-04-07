import React from 'react';
import { render, screen } from '@testing-library/react';
import { SchemaEnumDefinition } from './SchemaEnumDefinition';

describe('SchemaEnumDefinition', () => {
  it('renders an Enum Resolver', () => {
    render(<SchemaEnumDefinition resolverType='resolver' values={[]} />);
    expect(screen.getByTestId('enum-resolver')).toBeInTheDocument();
  });
});
