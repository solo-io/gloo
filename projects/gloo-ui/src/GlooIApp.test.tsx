import React from 'react';
import { render, waitForElement, fireEvent } from 'test-utils';
import { GlooIApp } from './GlooIApp';

xdescribe('<App />', () => {
  it('renders without errors', async () => {
    const { container, debug } = render(<GlooIApp />);
    expect(container).toBeDefined();
  });
  it('can navigate to virtual services page', async () => {
    const { getByText, debug, getByTestId } = render(<GlooIApp />);
    fireEvent.click(getByTestId('virtual-services-navlink'));
    const createVSButton = getByText('Create Virtual Service');
    expect(createVSButton).toBeDefined();
  });
});
