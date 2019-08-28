import css from '@emotion/css/macro';
import styled from '@emotion/styled';
import { soloConstants } from 'Styles/constants';
import { colors } from '../colors';

export const ErrorTextEmo = css`
  position: relative;
  font-size: 24px;
  color: ${colors.juneGrey};
  background: white;
  padding: 45px 15px 30px;

  :before {
    content: 'Sorry!';
    position: absolute;
    top: 10px;
    left: 0;
    right: 0;
    padding: 3px 15px;
    font-size: 14px;
    font-weight: bold;
    color: white;
    background: ${colors.oceanBlue};
  }
`;

export const ErrorTextFocusEmo = css`
  color: ${colors.pumpkinOrange};
`;

const ErrorMessageCSS = css`
  background-color: ${colors.tangerineOrange};
  border: 1px solid ${colors.grapefruitOrange};
  color: ${colors.pumpkinOrange};
  border-radius: ${soloConstants.smallRadius}px;
  padding: 10px 24px 13px;
  font-size: 16px;
  line-height: 16px;
  text-align: left;
  opacity: 1;

  span {
    font-weight: bold;
  }
`;

type LoadingProps = { loading?: boolean };
export const ErrorMessage = styled.div`
  ${ErrorMessageCSS};

  ${(props: LoadingProps) =>
    // @ts-ignore
    props.loading ? `opacity: .4;` : ''}
`;

export const SmallErrorMessage = styled.div`
  ${ErrorMessageCSS};
  padding: 7px 15px;

  ${(props: LoadingProps) =>
    // @ts-ignore
    props.loading ? `opacity: .4;` : ''}
`;
