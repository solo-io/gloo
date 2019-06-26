import * as React from 'react';
import * as ReactDOM from 'react-dom';
import { GlooIApp } from './GlooIApp';

it('renders without crashing', () => {
  const div = document.createElement('div');
  ReactDOM.render(<GlooIApp />, div);
  ReactDOM.unmountComponentAtNode(div);
});
