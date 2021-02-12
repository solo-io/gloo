import React from 'react';
import ReactDOM from 'react-dom';
import GlooFedApp from './GlooFedApp';
import * as serviceWorker from './serviceWorker';
import { SWRConfig } from 'swr';

ReactDOM.render(
  <SWRConfig
    value={{
      shouldRetryOnError: false,
      dedupingInterval: 2000,
    }}>
    <GlooFedApp />
  </SWRConfig>,
  document.getElementById('root')
);

if (module.hot) {
  module.hot.accept('./GlooFedApp', () => {
    const NextApp = require('./GlooFedApp').default;
    ReactDOM.render(
      <SWRConfig
        value={{
          shouldRetryOnError: false,
          dedupingInterval: 2000,
        }}>
        <NextApp />
      </SWRConfig>,
      document.getElementById('root')
    );
  });
}

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
serviceWorker.unregister();
