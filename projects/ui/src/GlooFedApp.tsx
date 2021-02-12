import React from 'react';
import { Global } from '@emotion/core';
import styled from '@emotion/styled';
import { globalStyles } from './Styles/globalStyles';
import { Footer } from 'Components/SiteStructure/Footer';
import { Content } from 'Components/SiteStructure/Content';
import { MainMenu } from 'Components/SiteStructure/MainMenu';
import { BrowserRouter } from 'react-router-dom';
import './Styles/styles.css';

const AppContainer = styled.div`
  display: grid;
  height: 100vh;
  grid-template-rows: 55px 1fr 62px;
`;

function GlooFedApp() {
  return (
    <BrowserRouter>
      <Global styles={globalStyles} />
      <AppContainer>
        <MainMenu />

        <Content />

        <Footer />
      </AppContainer>
    </BrowserRouter>
  );
}

export default GlooFedApp;
