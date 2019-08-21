/** @jsx jsx */
import { jsx } from '@emotion/core';
import * as ReactDOM from 'react-dom';
import { GlooIApp } from './GlooIApp';
import * as serviceWorker from './serviceWorker';
import './fontFace.css';
import { configureStore } from './store';

import { Provider } from 'react-redux';

const store = configureStore();

ReactDOM.render(
  <Provider store={store}>
    <GlooIApp />
  </Provider>,
  document.getElementById('root')
);

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
serviceWorker.unregister();
