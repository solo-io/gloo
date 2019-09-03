import React from 'react';
import { render } from '@testing-library/react';
import { configureStore, AppState } from 'store';
import { Provider } from 'react-redux';
import { Router } from 'react-router-dom';
import { createMemoryHistory } from 'history';

function renderWithWrappers(
  ui: React.ReactNode,
  { store = configureStore() } = {},
  {
    //@ts-ignore
    route = '/',
    history = createMemoryHistory({
      initialEntries: [route]
    })
  } = {}
) {
  return {
    ...render(
      <Provider store={store}>
        <Router history={history}>{ui}</Router>
      </Provider>
    ),
    store
  };
}

export function renderWithRouter(
  ui: React.ReactNode,
  {
    //@ts-ignore
    route = '/',
    history = createMemoryHistory({
      initialEntries: [route]
    })
  } = {}
) {
  return {
    ...render(<Router history={history}>{ui}</Router>)
  };
}
export * from '@testing-library/react';
export { renderWithWrappers as render };
