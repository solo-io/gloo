import css from '@emotion/css';
import { colors } from 'Styles/colors';
import { soloConstants } from '../constants';

export const CardCSS = css`
  position: relative;
  width: 100%;
  padding: ${soloConstants.largeBuffer}px ${soloConstants.smallBuffer}px
    ${soloConstants.smallBuffer}px;
  background: white;
  border-radius: ${soloConstants.radius}px;
  box-shadow: 0px 4px 9px ${colors.boxShadow};
  box-sizing: border-box;
`;
