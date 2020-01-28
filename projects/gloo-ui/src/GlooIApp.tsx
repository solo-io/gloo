import { Global } from '@emotion/core';
import styled from '@emotion/styled';
import { SuccessModal } from 'Components/Common/DisplayOnly/SuccessModal';
import { ErrorBoundary } from 'Components/Features/Errors/ErrorBoundary';
import * as React from 'react';
import { hot } from 'react-hot-loader/root';
import { useDispatch, useSelector } from 'react-redux';
import { BrowserRouter } from 'react-router-dom';
import { AppState } from 'store';
import { getIsLicenseValid } from 'store/config/actions';
import { Content } from './Components/Structure/Content';
import { Footer } from './Components/Structure/Footer';
import { MainMenu } from './Components/Structure/MainMenu';
import { globalStyles } from './Styles';
import './Styles/styles.css';

const AppContainer = styled.div`
  display: grid;
  min-height: 100vh;
  grid-template-rows: 55px 1fr 62px;
`;

const App = () => {
  const dispatch = useDispatch();

  React.useEffect(() => {
    dispatch(getIsLicenseValid());
  }, []);

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
