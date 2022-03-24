import React from 'react';
import { render, screen } from '@testing-library/react';
import { ExeGqlEnumDefinition } from './ExeGqlEnumDefinition';

describe('ExeGqlEnumDefinition', () => {
  it('renders an Enum Resolver', () => {
    render(<ExeGqlEnumDefinition resolverType='resolver' values={[]} />);
    expect(screen.getByTestId('enum-resolver')).toBeInTheDocument();
  });
});
