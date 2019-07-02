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

type Action = {
  type: string;
  payload: {};
};

type State = {};
const initialState: State = {};
export const StoreContext = React.createContext({
  state: initialState,
  dispatch: {} as React.Dispatch<Action>
});

//@ts-ignore
export const useStore = () => React.useContext(StoreContext);

const reducer: React.Reducer<State, Action> = (state, action) => {
  switch (action.type) {
    default:
      return state;
  }
};

const AppContainer = styled.div`
  display: grid;
  min-height: 100vh;
  grid-template-rows: 55px 1fr 62px;
`;

export const GlooIApp = () => {
  const [state, dispatch] = React.useReducer(reducer, initialState);

  return (
    <GlooEContext.Provider value={initialGlooEContext}>
      <StoreContext.Provider value={{ state, dispatch }}>
        <BrowserRouter>
          <Global styles={globalStyles} />
          <AppContainer>
            <MainMenu />
            <Content />
            <Footer />
          </AppContainer>
        </BrowserRouter>
      </StoreContext.Provider>
    </GlooEContext.Provider>
  );
};
