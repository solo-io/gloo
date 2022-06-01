import { default as SoloAddButton } from '../SoloAddButton';

import React from 'react';
import { render, screen } from '@testing-library/react';

describe('SoloAddButton', () => {
  it('Should render', async () => {
    render(<SoloAddButton data-testid='solo-add-button' />);
    const domEl = await screen.queryByTestId('solo-add-button');
    expect(domEl).not.toBe(null);
  });
});
