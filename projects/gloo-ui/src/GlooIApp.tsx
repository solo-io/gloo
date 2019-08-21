import * as React from 'react';
/** @jsx jsx */
import { jsx } from '@emotion/core';
import styled from '@emotion/styled/macro';
import { BrowserRouter } from 'react-router-dom';
import { MainMenu } from './Components/Structure/MainMenu';
import { Content } from './Components/Structure/Content';
import { Global } from '@emotion/core';
import { globalStyles } from './Styles';
import { Footer } from './Components/Structure/Footer';
import { GlooEContext, initialGlooEContext } from 'Api';
import './Styles/styles.css';

import { useDispatch, useSelector } from 'react-redux';
import { listUpstreams } from 'store/upstreams/actions';
import { listVirtualServices } from 'store/virtualServices/actions';
import {
  listNamespaces,
  getSettings,
  getPodNamespace,
  getIsLicenseValid,
  getVersion
} from 'store/config/actions';
import { AppState } from 'store';
import { listSecrets } from 'store/secrets/actions';

const AppContainer = styled.div`
  display: grid;
  min-height: 100vh;
  grid-template-rows: 55px 1fr 62px;
`;

export const GlooIApp = () => {
  const dispatch = useDispatch();

  // TODO: make a generalized action in reducer
  React.useEffect(() => {
    dispatch(listNamespaces());
    dispatch(getSettings());
    dispatch(getPodNamespace());
    dispatch(getIsLicenseValid());
    dispatch(getVersion());
  }, []);

  const { namespacesList } = useSelector((store: AppState) => store.config);

  React.useEffect(() => {
    if (namespacesList) {
      dispatch(listUpstreams({ namespacesList }));
      dispatch(listVirtualServices({ namespacesList }));
      dispatch(listSecrets({ namespacesList }));
    }
  }, [namespacesList.length]);

  return (
    <GlooEContext.Provider value={initialGlooEContext}>
      <BrowserRouter>
        <Global styles={globalStyles} />
        <AppContainer>
          <MainMenu />
          <Content />
          <Footer />
        </AppContainer>
      </BrowserRouter>
    </GlooEContext.Provider>
  );
};
