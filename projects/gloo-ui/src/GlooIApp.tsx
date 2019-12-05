import { Global } from '@emotion/core';
import styled from '@emotion/styled';
import { useInterval } from 'Hooks/useInterval';
import * as React from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { BrowserRouter } from 'react-router-dom';
import { AppState } from 'store';
import {
  getIsLicenseValid,
  getPodNamespace,
  getSettings,
  getVersion,
  listNamespaces
} from 'store/config/actions';
import { listEnvoyDetails } from 'store/envoy/actions';
import { listGateways } from 'store/gateway/actions';
import { listProxies } from 'store/proxy/actions';
import { listSecrets } from 'store/secrets/actions';
import { listUpstreams } from 'store/upstreams/actions';
import { listVirtualServices } from 'store/virtualServices/actions';
import { Content } from './Components/Structure/Content';
import { Footer } from './Components/Structure/Footer';
import { MainMenu } from './Components/Structure/MainMenu';
import { globalStyles } from './Styles';
import './Styles/styles.css';
import { hot } from 'react-hot-loader/root';
import { SuccessModal } from 'Components/Common/DisplayOnly/SuccessModal';
import { listRouteTables } from 'store/routeTables/actions';
import { ErrorBoundary } from 'Components/Features/Errors/ErrorBoundary';

const AppContainer = styled.div`
  display: grid;
  min-height: 100vh;
  grid-template-rows: 55px 1fr 62px;
`;

const App = () => {
  const dispatch = useDispatch();

  React.useEffect(() => {
    dispatch(listNamespaces());
    dispatch(getSettings());
    dispatch(getPodNamespace());
    dispatch(getIsLicenseValid());
    dispatch(getVersion());
    dispatch(listEnvoyDetails());
    dispatch(listUpstreams());
    dispatch(listVirtualServices());
    dispatch(listSecrets());
    dispatch(listGateways());
    dispatch(listProxies());
    dispatch(listRouteTables());
  }, []);

  useInterval(() => {
    dispatch(listVirtualServices());
    dispatch(listSecrets());
    dispatch(listGateways());
    dispatch(listProxies());
    dispatch(listNamespaces());
    dispatch(getSettings());
    dispatch(getPodNamespace());
    dispatch(getVersion());
    dispatch(listEnvoyDetails());
    dispatch(listRouteTables());
  }, 3000);

  useInterval(() => {
    dispatch(getIsLicenseValid());
    dispatch(listUpstreams());
  }, 60000);
  const showModal = useSelector((state: AppState) => state.modal.showModal);
  const modalMessage = useSelector((state: AppState) => state.modal.message);
  return (
    <BrowserRouter>
      <SuccessModal visible={!!showModal} successMessage={modalMessage} />
      <Global styles={globalStyles} />
      <AppContainer>
        <ErrorBoundary fallback={<div> there was an error</div>}>
          <MainMenu />
          <Content />
          <Footer />
        </ErrorBoundary>
      </AppContainer>
    </BrowserRouter>
  );
};
export const GlooIApp = hot(App);
