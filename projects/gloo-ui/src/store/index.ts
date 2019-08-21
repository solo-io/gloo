import { createStore, combineReducers, applyMiddleware } from 'redux';
import thunkMiddleware from 'redux-thunk';
import { composeWithDevTools } from 'redux-devtools-extension';
import { upstreamsReducer } from './upstreams/reducers';
import { virtualServicesReducer } from './virtualServices/reducers';
import { loadingBarReducer } from 'react-redux-loading-bar';
import { secretsReducer } from './secrets/reducers';
import { normalizrMiddleware } from './requests';
import { configReducer } from './config/reducers';

const rootReducer = combineReducers({
  upstreams: upstreamsReducer,
  virtualServices: virtualServicesReducer,
  secrets: secretsReducer,
  config: configReducer,
  loadingBar: loadingBarReducer
});

export type AppState = ReturnType<typeof rootReducer>;

export function configureStore() {
  const middlewares = [thunkMiddleware /*normalizrMiddleware */];
  const middleWareEnhancer = applyMiddleware(...middlewares);

  const store = createStore(
    rootReducer,
    composeWithDevTools(middleWareEnhancer)
  );

  return store;
}
